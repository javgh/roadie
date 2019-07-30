package bob

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/NebulousLabs/Sia/types"

	"github.com/javgh/roadie/blockchain/ethereum"
	"github.com/javgh/roadie/blockchain/sia"
	"github.com/javgh/roadie/trader"
)

const (
	bindingOfferLifetime = 1 * time.Minute
)

var (
	antiSpamFee = big.NewInt(1e14)
	oneSiacoin  = types.SiacoinPrecision
)

func TestAtomicSwap(t *testing.T) {
	ethChain, err := ethereum.NewSimulatedBlockchain()
	if err != nil {
		t.Fatal(err)
	}

	siaChain, err := sia.NewSimulatedBlockchain()
	if err != nil {
		t.Fatal(err)
	}

	trader := trader.NewFixedPremiumTrader(nil, *antiSpamFee, ethChain, siaChain)
	blacklist := NewBlacklist()
	now := time.Now()

	t.Run("CriticalPhaseAndBlacklist", func(t *testing.T) {
		swap1 := NewAtomicSwap(&trader, ethChain, siaChain, blacklist, now)
		swap2 := NewAtomicSwap(&trader, ethChain, siaChain, blacklist, now)

		nonBindingOffer1, err := swap1.RequestNonBindingOffer(oneSiacoin, now)
		if err != nil {
			t.Fatal(err)
		}
		assert.True(t, nonBindingOffer1.Available, "should receive non-binding offer")

		nonBindingOffer2, err := swap2.RequestNonBindingOffer(oneSiacoin, now)
		if err != nil {
			t.Fatal(err)
		}
		assert.True(t, nonBindingOffer2.Available, "should receive non-binding offer")

		antiSpamID := big.NewInt(0)
		err = ethChain.BurnAntiSpamFee(*antiSpamID, nonBindingOffer1.AntiSpamFee)
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(4 * time.Second) // wait for confirmations

		bindingOffer1, err := swap1.RequestBindingOffer(*antiSpamID, now)
		if err != nil {
			t.Fatal(err)
		}
		assert.True(t, bindingOffer1.Available, "should receive binding offer")

		bindingOffer2, err := swap2.RequestBindingOffer(*antiSpamID, now)
		if err != nil {
			t.Fatal(err)
		}
		assert.False(t, bindingOffer2.Available, "should not have multiple binding offers in parallel")
		assert.Contains(t, bindingOffer2.Msg, "critical phase", "should inform of critical phase")

		later := now.Add(2 * bindingOfferLifetime)
		_, err = swap2.RequestBindingOffer(*antiSpamID, later)
		assert.Equal(t, ErrAntiSpamReused, err, "should not allow re-use of anti spam id")
	})
}
