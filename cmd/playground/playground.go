package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/HyperspaceApp/ed25519"
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
	defaultClientAddress  = "localhost:9980"
	defaultPasswordFile   = ".sia/apipassword"
	bindingOfferLifetime  = 1 * time.Minute
	antiSpamConfirmations = 10
	depositConfirmations  = 10
	jsonRPCEndpoint       = ".ethereum/geth.ipc"
	jsonRPCKeystoreFile   = ".config/roadie/keystore"
	boostInterval         = 90 * time.Second
)

var (
	oneSiacoin         = types.SiacoinPrecision
	gwei               = big.NewInt(1e9)
	finney             = big.NewInt(1e15)
	defaultAntiSpamFee = big.NewInt(1e14)
	mockWalletAddress  = common.HexToAddress("0x0000000000000000000000000000000000000000")
	contractAddress    = common.HexToAddress("0x799DF2482f589663d7754451de3FfeF4CAA439c8")
	maxGasPrice        = new(big.Int).Mul(big.NewInt(21), gwei)
)

type (
	mockTrader struct{}
)

func (mt *mockTrader) PrepareNonBindingOffer(siacoin types.Currency, minerFee types.Currency) (*trader.Offer, error) {
	offer := trader.Offer{
		Msg:         "playground offer",
		Available:   true,
		Ether:       *finney,
		AntiSpamFee: *defaultAntiSpamFee,
	}
	return &offer, nil
}

func (mt *mockTrader) PrepareBindingOffer(siacoin types.Currency, minerFee types.Currency,
	now time.Time) (*trader.Offer, *time.Time, error) {
	offer, err := mt.PrepareNonBindingOffer(siacoin, minerFee)
	if err != nil {
		return nil, nil, err
	}

	deadline := now.Add(bindingOfferLifetime)
	return offer, &deadline, nil
}

func (mt *mockTrader) PauseOrderPreparation(now time.Time) {
}

func (mt *mockTrader) ResumeOrderPreparation() {
}

type mockChain struct {
	adaptorPrivKey ed25519.Adaptor
}

func (mc *mockChain) BurnAntiSpamFee(antiSpamID big.Int, antiSpamFee big.Int) error {
	return nil
}

func (mc *mockChain) CheckAntiSpamConfirmations(antiSpamID big.Int, antiSpamFee big.Int) (int64, error) {
	return antiSpamConfirmations, nil
}

func (mc *mockChain) DepositEther(
	recipient common.Address, adaptorPubKey ed25519.CurvePoint, ether big.Int, antiSpamID big.Int) error {
	return nil
}

func (mc *mockChain) CheckDepositConfirmations(
	recipient common.Address, adaptorPubKey ed25519.CurvePoint, ether big.Int, antiSpamID big.Int) (int64, error) {
	return depositConfirmations, nil
}

func (mc *mockChain) ClaimDeposit(adaptorPrivKey ed25519.Adaptor, antiSpamID big.Int) error {
	mc.adaptorPrivKey = adaptorPrivKey
	return nil
}

func (mc *mockChain) LookupAdaptorPrivKey(adaptorPubKey ed25519.CurvePoint) (bool, *ed25519.Adaptor, error) {
	return true, &mc.adaptorPrivKey, nil
}

func (mc *mockChain) WalletAddress() common.Address {
	return mockWalletAddress
}

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

	//go server(ethChain, &drSiaChain)
	//time.Sleep(2 * time.Second)
	//client(ethChain, &drSiaChain)

	trader := trader.NewFixedPremiumTrader(nil, *defaultAntiSpamFee, ethChain, &drSiaChain)

	offer, err := trader.PrepareNonBindingOffer(oneSiacoin, oneSiacoin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(offer.Msg)
	fmt.Println(ethereum.FormatEther(&offer.Ether))
}

func server(ethChain ethereum.Blockchain, siaChain sia.Blockchain) {
	mockTrader := mockTrader{}
	blacklist := bob.NewBlacklist()

	newAtomicSwap := func(now time.Time) *bob.AtomicSwap {
		return bob.NewAtomicSwap(&mockTrader, ethChain, siaChain, blacklist, now)
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
