package main

import (
	"io/ioutil"
	"log"
	"math/big"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/NebulousLabs/Sia/types"

	"github.com/javgh/roadie/alice"
	"github.com/javgh/roadie/blockchain/ethereum"
	"github.com/javgh/roadie/blockchain/sia"
	"github.com/javgh/roadie/bob"
	"github.com/javgh/roadie/rpc"
	"github.com/javgh/roadie/trader"
)

const (
	defaultClientAddress = "localhost:9980"
	defaultPasswordFile  = ".sia/apipassword"
	jsonRPCEndpoint      = ".ethereum/geth.ipc"
	jsonRPCKeystoreFile  = ".config/roadie/keystore"
	boostInterval        = 90 * time.Second
)

var (
	oneSiacoin         = types.SiacoinPrecision
	gwei               = big.NewInt(1e9)
	finney             = big.NewInt(1e15)
	defaultAntiSpamFee = big.NewInt(1e14)
	contractAddress    = common.HexToAddress("0x799DF2482f589663d7754451de3FfeF4CAA439c8")
	maxGasPrice        = new(big.Int).Mul(big.NewInt(21), gwei)
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

	go server(ethChain, &drSiaChain)
	time.Sleep(2 * time.Second)
	client(ethChain, &drSiaChain)
}

func server(ethChain ethereum.Blockchain, siaChain sia.Blockchain) {
	trader := trader.NewFixedPremiumTrader(nil, *defaultAntiSpamFee, ethChain, siaChain)
	blacklist := bob.NewBlacklist()

	newAtomicSwap := func(now time.Time) *bob.AtomicSwap {
		return bob.NewAtomicSwap(&trader, ethChain, siaChain, blacklist, now)
	}
	bobServer, err := rpc.NewBobServer("tcp", "localhost:9000", "", "", newAtomicSwap)
	if err != nil {
		log.Fatal(err)
	}

	err = bobServer.Serve()
	if err != nil {
		log.Fatal(err)
	}
}

func client(ethChain ethereum.Blockchain, siaChain sia.Blockchain) {
	frontend := alice.NewConsoleFrontend()
	//frontend := alice.AutoAcceptFrontend{}

	roadieClient, err := rpc.Dial("localhost:9000")
	if err != nil {
		log.Fatal(err)
	}

	err = alice.PerformSwap(oneSiacoin, frontend, ethChain, siaChain, roadieClient)
	if err != nil {
		log.Fatal(err)
	}
}
