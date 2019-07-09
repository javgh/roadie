package integration

import (
	"math/big"
	"testing"
	"time"

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
	serverNetwork = "tcp"
	serverAddress = "localhost:9000"
)

var (
	oneSiacoin         = types.SiacoinPrecision
	defaultAntiSpamFee = big.NewInt(1e14)
)

func TestIntegration(t *testing.T) {
	ethChain, err := ethereum.NewSimulatedBlockchain()
	if err != nil {
		t.Fatal(err)
	}

	siaChain, err := sia.NewSimulatedBlockchain()
	if err != nil {
		t.Fatal(err)
	}

	go server(t, ethChain, siaChain)
	time.Sleep(2 * time.Second)
	client(t, ethChain, siaChain)
}

func server(t *testing.T, ethChain ethereum.Blockchain, siaChain sia.Blockchain) {
	trader := trader.NewFixedPremiumTrader(nil, *defaultAntiSpamFee, ethChain, siaChain)
	blacklist := bob.NewBlacklist()

	newAtomicSwap := func(now time.Time) *bob.AtomicSwap {
		return bob.NewAtomicSwap(&trader, ethChain, siaChain, blacklist, now)
	}
	bobServer, err := rpc.NewBobServer(serverNetwork, serverAddress, "", "", serverAddress, newAtomicSwap)
	if err != nil {
		t.Fatal(err)
	}

	err = bobServer.Serve()
	if err != nil {
		t.Fatal(err)
	}
}

func client(t *testing.T, ethChain ethereum.Blockchain, siaChain sia.Blockchain) {
	frontend := frontend.AutoAcceptFrontend{}

	roadieClient, err := rpc.Dial(serverAddress, []byte{})
	if err != nil {
		t.Fatal(err)
	}

	err = alice.PerformSwap(oneSiacoin, frontend, ethChain, siaChain, roadieClient)
	if err != nil {
		t.Fatal(err)
	}
}
