package trader

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

	now := time.Now()
	trader := NewFixedPremiumTrader(nil, *antiSpamFee, ethChain, siaChain)

	t.Run("TooSmall", func(t *testing.T) {
		offer, err := trader.PrepareNonBindingOffer(types.ZeroCurrency, minerFee, now)
		if err != nil {
			t.Fatal(err)
		}

		assert.False(t, offer.Available, "expected no offer for 0 SC")
	})

	t.Run("TooLarge", func(t *testing.T) {
		siacoin := types.SiacoinPrecision.Mul64(100000000)
		offer, err := trader.PrepareNonBindingOffer(siacoin, minerFee, now)
		if err != nil {
			t.Fatal(err)
		}

		assert.False(t, offer.Available, "expected no offer for 100000000 SC")
	})

	t.Run("SmallAmount", func(t *testing.T) {
		siacoin := types.SiacoinPrecision.Mul64(1000)
		offer, err := trader.PrepareNonBindingOffer(siacoin, minerFee, now)
		if err != nil {
			t.Fatal(err)
		}

		assert.True(t, offer.Available, "expected offer for 1000 SC")
		assert.Equal(t, 1, offer.Ether.Sign(), "expected Ether amount to be positive")
	})

	t.Run("PausedOrderPreparation", func(t *testing.T) {
		trader.PauseOrderPreparation(now)

		siacoin := types.SiacoinPrecision.Mul64(1000)
		offer, err := trader.PrepareNonBindingOffer(siacoin, minerFee, now)
		if err != nil {
			t.Fatal(err)
		}

		assert.False(t, offer.Available, "expected offer preparation to be paused")
	})

	t.Run("ResumedOrderPreparation", func(t *testing.T) {
		trader.ResumeOrderPreparation()

		siacoin := types.SiacoinPrecision.Mul64(1000)
		offer, err := trader.PrepareNonBindingOffer(siacoin, minerFee, now)
		if err != nil {
			t.Fatal(err)
		}

		assert.True(t, offer.Available, "expected offer preparation to be resumed")
	})

	t.Run("ResumeAfterDeadline", func(t *testing.T) {
		trader.PauseOrderPreparation(now)
		later := now.Add(2 * time.Minute)

		siacoin := types.SiacoinPrecision.Mul64(1000)
		offer, err := trader.PrepareNonBindingOffer(siacoin, minerFee, later)
		if err != nil {
			t.Fatal(err)
		}

		assert.True(t, offer.Available, "expected offer preparation to be resumed")
	})
}

func TestCheckSimilarity(t *testing.T) {
	a := Offer{Available: false}
	b := Offer{Available: false}
	assert.True(t, CheckSimilarity(a, b, 0), "two unavailable offers are similar")

	b.Available = true
	assert.False(t, CheckSimilarity(a, b, 0), "if only one offer is available, they are not similar")

	a.Available = true
	assert.True(t, CheckSimilarity(a, b, 0), "two zero offers are similar")

	one := big.NewInt(1)
	a.Ether = *one
	assert.False(t, CheckSimilarity(a, b, 0), "if only one offer is zero, they are not similar")

	b.Ether = *one
	assert.True(t, CheckSimilarity(a, b, 0), "offers with the same amount are similar")

	two := big.NewInt(2)
	b.Ether = *two
	assert.False(t, CheckSimilarity(a, b, 0), "if one offer is double the other, they are not similar")

	large1 := big.NewInt(100)
	a.Ether = *large1
	large2 := big.NewInt(101)
	b.Ether = *large2
	assert.False(t, CheckSimilarity(a, b, 0), "offers with different amounts are not similar with zero tolerance")
	assert.True(t, CheckSimilarity(a, b, 2), "offers are similar if within the specified tolerance")

	large3 := big.NewInt(103)
	b.Ether = *large3
	assert.False(t, CheckSimilarity(a, b, 2), "offers are not similar if outside the specified tolerance")
}
