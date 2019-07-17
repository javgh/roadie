package main

import (
	"log"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"

	"github.com/javgh/roadie/blockchain/ethereum"
	"github.com/javgh/roadie/blockchain/sia"
	"github.com/javgh/roadie/bob"
	"github.com/javgh/roadie/config"
	"github.com/javgh/roadie/rpc"
	"github.com/javgh/roadie/trader"
)

const (
	serverNetwork         = "tcp"
	registryCheckInterval = 12 * time.Hour
	serverCheckInterval   = time.Hour
)

var (
	siaDaemonAddress     = "localhost:9980"
	siaPasswordFile      = config.PrependHomeDirectory(".sia/apipassword")
	siaDryRun            = false
	useGanache           = false
	serverAddress        = "localhost:9000"
	externalAddress      = "localhost:9000"
	certFile             = ""
	keyFile              = ""
	jsonRPCEndpoint      = config.PrependHomeDirectory(".ethereum/geth.ipc")
	keystoreFile         = config.PrependConfigDirectory("keystore")
	maxGasPriceInGwei    = int64(21)
	boostIntervalSeconds = int64(90)

	gwei                = big.NewInt(1e9)
	defaultAntiSpamFee  = big.NewInt(1e14)
	registryEntryMaxAge = big.NewInt(14 * 24 * 60 * 60) // 14 days in seconds
	contractAddress     = common.HexToAddress("0x8CeF4dDFFcad47Ead5389A60ca9771EEe33Fd460")
)

func serve(cmd *cobra.Command, args []string) {
	var err error
	var ethChain ethereum.Blockchain
	if useGanache {
		ethChain, err = ethereum.NewGanacheBlockchain()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = ethereum.EnsureKeystoreExists(keystoreFile)
		if err != nil {
			log.Fatal(err)
		}

		maxGasPrice := new(big.Int).Mul(big.NewInt(maxGasPriceInGwei), gwei)
		boostInterval := time.Duration(boostIntervalSeconds) * time.Second
		ethChain, err = ethereum.NewLocalNodeBlockchain(
			jsonRPCEndpoint, keystoreFile, &contractAddress, *maxGasPrice, boostInterval)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = ethChain.CheckSmartContract()
	if err != nil {
		log.Fatal(err)
	}

	siaPassword, err := config.ReadPasswordFile(siaPasswordFile)
	if err != nil {
		log.Fatal(err)
	}

	lnSiaChain, err := sia.NewLocalNodeBlockchain(siaDaemonAddress, siaPassword)
	if err != nil {
		log.Fatal(err)
	}

	var siaChain sia.Blockchain = lnSiaChain
	if siaDryRun {
		siaChain = sia.NewDryRunBlockchain(lnSiaChain)
	}

	trader := trader.NewFixedPremiumTrader(nil, *defaultAntiSpamFee, ethChain, siaChain)
	blacklist := bob.NewBlacklist()

	newAtomicSwap := func(now time.Time) *bob.AtomicSwap {
		return bob.NewAtomicSwap(&trader, ethChain, siaChain, blacklist, now)
	}
	bobServer, err := rpc.NewBobServer(serverNetwork, serverAddress, certFile, keyFile, externalAddress, newAtomicSwap)
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR1)
	go func() {
		for {
			<-c
			log.Println("Report requested")
			bobServer.Report()
		}
	}()

	err = bobServer.Register(*registryEntryMaxAge, ethChain)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			time.Sleep(registryCheckInterval)
			err2 := bobServer.Register(*registryEntryMaxAge, ethChain)
			if err2 != nil {
				log.Printf("Error while attempting to re-register: %s\n", err2)
			}
		}
	}()

	go func() {
		for {
			time.Sleep(serverCheckInterval)
			err3 := bobServer.Check(time.Now())
			if err3 != nil {
				log.Printf("Error while running check: %s\n", err3)
			}
		}
	}()

	err = bobServer.Serve()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	cmdServe := &cobra.Command{
		Use:   "serve",
		Short: "Start and register a server to offer atomic swaps",
		Run:   serve,
	}
	cmdServe.Flags().StringVarP(&serverAddress, "listen", "l", serverAddress, "interface and port to listen on")
	cmdServe.Flags().StringVarP(&certFile, "cert", "c", certFile, "path to certificate (or omit to disable encryption)")
	cmdServe.Flags().StringVarP(&keyFile, "key", "k", certFile, "path to certificate key (or omit to disable encryption)")
	cmdServe.Flags().StringVarP(&externalAddress, "addr", "a", externalAddress, "external server address (host and port to register with the smart contract)")

	rootCmd := &cobra.Command{Use: "roadie-bob"}
	rootCmd.AddCommand(cmdServe)
	rootCmd.PersistentFlags().StringVar(&siaPasswordFile, "sia-password-file", siaPasswordFile, "path to Sia API password file")
	rootCmd.PersistentFlags().StringVar(&siaDaemonAddress, "sia-daemon", siaDaemonAddress, "host and port of Sia daemon")
	rootCmd.PersistentFlags().BoolVar(&siaDryRun, "sia-dry-run", siaDryRun, "do not actually broadcast Sia transactions")
	rootCmd.PersistentFlags().BoolVarP(&useGanache, "ganache", "g", useGanache, "use Ganache as Ethereum node (expected at 127.0.0.1:8545)")
	rootCmd.PersistentFlags().StringVar(&jsonRPCEndpoint, "ethereum-node", jsonRPCEndpoint, "IPC socket/pipe to Ethereum node")
	rootCmd.PersistentFlags().Int64Var(&maxGasPriceInGwei, "max-gas-price", maxGasPriceInGwei, "maximum amount (in Gwei) when boosting the gas price")
	rootCmd.PersistentFlags().Int64Var(&boostIntervalSeconds, "boost-interval", boostIntervalSeconds, "seconds to wait for a transaction to confirm before boosting gas price")

	rootCmd.Execute()
}
