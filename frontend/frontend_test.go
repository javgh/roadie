package frontend

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/NebulousLabs/Sia/types"

	"github.com/javgh/roadie/trader"
)

type (
	MockExchangeRate struct{}
)

var (
	oneEther = big.NewInt(1e18)
)

func (r *MockExchangeRate) Fetch(id string) (*big.Rat, error) {
	if id == "ethereum" {
		return new(big.Rat).SetInt(oneEther), nil
	}

	return new(big.Rat).SetInt(types.SiacoinPrecision.Big()), nil
}

func assertApproveOffer(t *testing.T, frontend *RuleBasedFrontend,
	siacoin types.Currency, offer trader.Offer, value bool, msg string) {
	approved, err := frontend.ApproveOffer(siacoin, offer, false)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, value, approved, msg)
}

func TestRuleBasedFrontend(t *testing.T) {
	exchangeRate := &MockExchangeRate{}
	frontend := NewRuleBasedFrontend(1.0, 5.0, exchangeRate)

	siacoin := types.NewCurrency64(100)

	offer := trader.Offer{Available: false}
	assertApproveOffer(t, frontend, siacoin, offer, false, "should decline unavailable offer")

	offer.Available = true
	assertApproveOffer(t, frontend, siacoin, offer, true, "should approve offer where we pay nothing")

	ether := big.NewInt(55)
	antiSpamFee := big.NewInt(55)
	offer.Ether = *ether
	offer.AntiSpamFee = *antiSpamFee
	assertApproveOffer(t, frontend, siacoin, offer, false, "should decline expensive offer")

	siacoin = types.NewCurrency64(1)
	ether = big.NewInt(1)
	antiSpamFee = big.NewInt(1)
	offer.Ether = *ether
	offer.AntiSpamFee = *antiSpamFee
	assertApproveOffer(t, frontend, siacoin, offer, true, "should approve offer based on small absolute difference")

	siacoin = types.NewCurrency64(1)
	ether = big.NewInt(2)
	antiSpamFee = big.NewInt(1)
	offer.Ether = *ether
	offer.AntiSpamFee = *antiSpamFee
	assertApproveOffer(t, frontend, siacoin, offer, false, "should decline offer based on large absolute difference")

	siacoin = types.NewCurrency64(100)
	ether = big.NewInt(100)
	antiSpamFee = big.NewInt(5)
	offer.Ether = *ether
	offer.AntiSpamFee = *antiSpamFee
	assertApproveOffer(t, frontend, siacoin, offer, true, "should approve offer based on small relative difference")

	siacoin = types.NewCurrency64(100)
	ether = big.NewInt(100)
	antiSpamFee = big.NewInt(6)
	offer.Ether = *ether
	offer.AntiSpamFee = *antiSpamFee
	assertApproveOffer(t, frontend, siacoin, offer, false, "should decline offer based on large relative difference")

	siacoin = types.NewCurrency64(100)
	ether = big.NewInt(99)
	antiSpamFee = big.NewInt(0)
	offer.Ether = *ether
	offer.AntiSpamFee = *antiSpamFee
	assertApproveOffer(t, frontend, siacoin, offer, true, "should approve offer where we make money")
}
