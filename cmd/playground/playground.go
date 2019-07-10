package main

import (
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/NebulousLabs/Sia/types"

	"github.com/javgh/roadie/alice"
	"github.com/javgh/roadie/blockchain/ethereum"
	"github.com/javgh/roadie/blockchain/sia"
	"github.com/javgh/roadie/bob"
	"github.com/javgh/roadie/frontend"
	"github.com/javgh/roadie/rpc"
	"github.com/javgh/roadie/trader"
)

const (
	defaultClientAddress  = "localhost:9980"
	defaultPasswordFile   = ".sia/apipassword"
	jsonRPCEndpoint       = ".ethereum/geth.ipc"
	jsonRPCKeystoreFile   = ".config/roadie/keystore"
	boostInterval         = 90 * time.Second
	serverCheckInterval   = time.Hour
	registryCheckInterval = 12 * time.Hour
)

var (
	oneSiacoin                    = types.SiacoinPrecision
	gwei                          = big.NewInt(1e9)
	finney                        = big.NewInt(1e15)
	defaultAntiSpamFee            = big.NewInt(1e14)
	contractAddress               = common.HexToAddress("0x799DF2482f589663d7754451de3FfeF4CAA439c8")
	maxGasPrice                   = new(big.Int).Mul(big.NewInt(21), gwei)
	registryEntryMaxAge           = big.NewInt(14 * 24 * 60 * 60) // 14 days in seconds
	registryEntryMaxAgeWithMargin = big.NewInt(15 * 24 * 60 * 60) // 15 days in seconds
)

func prependHomeDirectory(path string) string {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(currentUser.HomeDir, path)
}

func main() {
	//ethChain, err := ethereum.NewGanacheBlockchain()
	ethChain, err := ethereum.NewSimulatedBlockchain()
	//endpoint := prependHomeDirectory(jsonRPCEndpoint)
	//keystoreFile := prependHomeDirectory(jsonRPCKeystoreFile)
	//ethChain, err := ethereum.NewLocalNodeBlockchain(
	//	endpoint, keystoreFile, &contractAddress, *maxGasPrice, boostInterval)
	if err != nil {
		log.Fatal(err)
	}

	passwordBytes, err := ioutil.ReadFile(prependHomeDirectory(defaultPasswordFile))
	if err != nil {
		log.Fatal(err)
	}
	password := strings.TrimSpace(string(passwordBytes))

	//siaChain, err := sia.NewSimulatedBlockchain()
	siaChain, err := sia.NewLocalNodeBlockchain(defaultClientAddress, password)
	if err != nil {
		log.Fatal(err)
	}
	drSiaChain := sia.NewDryRunBlockchain(*siaChain)

	err = ethChain.CheckSmartContract()
	if err != nil {
		log.Fatal(err)
	}

	go server(ethChain, &drSiaChain)
	time.Sleep(time.Second)
	client(ethChain, &drSiaChain)
}

func server(ethChain ethereum.Blockchain, siaChain sia.Blockchain) {
	trader := trader.NewFixedPremiumTrader(nil, *defaultAntiSpamFee, ethChain, siaChain)
	blacklist := bob.NewBlacklist()

	newAtomicSwap := func(now time.Time) *bob.AtomicSwap {
		return bob.NewAtomicSwap(&trader, ethChain, siaChain, blacklist, now)
	}
	bobServer, err := rpc.NewBobServer(
		"tcp", "localhost:9000", "./testdata/server.crt", "./testdata/server.key", "localhost:9000", newAtomicSwap)
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

func client(ethChain ethereum.Blockchain, siaChain sia.Blockchain) {
	frontend := frontend.NewConsoleFrontend()
	//frontend := frontend.AutoAcceptFrontend{}

	serverDetails, err := ethChain.FetchServers(*registryEntryMaxAgeWithMargin)
	if len(serverDetails) == 0 {
		log.Fatal("no server available")
	}

	roadieClient, err := rpc.Dial(serverDetails[0].Target, serverDetails[0].Cert)
	if err != nil {
		log.Fatal(err)
	}

	err = alice.PerformSwap(oneSiacoin, frontend, ethChain, siaChain, roadieClient)
	if err != nil {
		log.Fatal(err)
	}
}
