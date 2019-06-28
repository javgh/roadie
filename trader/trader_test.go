package trader

import (
	"math/big"
	"testing"

	"gitlab.com/NebulousLabs/Sia/types"

	"github.com/javgh/roadie/blockchain/ethereum"
	"github.com/javgh/roadie/blockchain/sia"
)

var (
	antiSpamFee = big.NewInt(1e14)
	minerFee    = types.SiacoinPrecision
)

func TestTrader(t *testing.T) {
	ethChain, err := ethereum.NewSimulatedBlockchain()
	if err != nil {
		t.Fatal(err)
	}

	siaChain, err := sia.NewSimulatedBlockchain()
	if err != nil {
		t.Fatal(err)
	}

	trader := NewFixedPremiumTrader(nil, *antiSpamFee, ethChain, siaChain)

	t.Run("TooSmall", func(t *testing.T) {
		offer, err := trader.PrepareNonBindingOffer(types.ZeroCurrency, minerFee)
		if err != nil {
			t.Fatal(err)
		}

		if offer.Available {
			t.Errorf("Expected no offer for 0 SC")
		}
	})

	t.Run("TooLarge", func(t *testing.T) {
		siacoin := types.SiacoinPrecision.Mul64(100000000)
		offer, err := trader.PrepareNonBindingOffer(siacoin, minerFee)
		if err != nil {
			t.Fatal(err)
		}

		if offer.Available {
			t.Errorf("Expected no offer for 100000000 SC")
		}
	})

	t.Run("SmallAmount", func(t *testing.T) {
		siacoin := types.SiacoinPrecision.Mul64(1000)
		offer, err := trader.PrepareNonBindingOffer(siacoin, minerFee)
		if err != nil {
			t.Fatal(err)
		}

		if !offer.Available {
			t.Errorf("Expected offer for 1000 SC")
		}

		if offer.Ether.Sign() != 1 {
			t.Errorf("Expected Ether amount to be positive")
		}
	})
}