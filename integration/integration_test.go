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
	"github.com/javgh/roadie/rpc"
	"github.com/javgh/roadie/trader"
)

type (
	mockTrader struct{}
)

const (
	bindingOfferLifetime = 1 * time.Minute
	serverNetwork        = "tcp"
	serverAddress        = "localhost:9000"
)

var (
	oneSiacoin         = types.SiacoinPrecision
	finney             = big.NewInt(1e15)
	defaultAntiSpamFee = big.NewInt(1e14)
)

func (mt *mockTrader) PrepareNonBindingOffer(siacoin types.Currency, minerFee types.Currency) (*trader.Offer, error) {
	offer := trader.Offer{
		Msg:         "integration test offer",
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
	mockTrader := mockTrader{}
	blacklist := bob.NewBlacklist()

	newAtomicSwap := func(now time.Time) *bob.AtomicSwap {
		return bob.NewAtomicSwap(&mockTrader, ethChain, siaChain, blacklist, now)
	}
	bobServer, err := rpc.NewBobServer(serverNetwork, serverAddress, "", "", newAtomicSwap)
	if err != nil {
		t.Fatal(err)
	}

	err = bobServer.Serve()
	if err != nil {
		t.Fatal(err)
	}
}

func client(t *testing.T, ethChain ethereum.Blockchain, siaChain sia.Blockchain) {
	frontend := alice.AutoAcceptFrontend{}

	roadieClient, err := rpc.Dial(serverAddress)
	if err != nil {
		t.Fatal(err)
	}

	err = alice.PerformSwap(oneSiacoin, frontend, ethChain, siaChain, roadieClient)
	if err != nil {
		t.Fatal(err)
	}
}
