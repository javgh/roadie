package frontend

import (
	"fmt"
	"math/big"

	"gitlab.com/NebulousLabs/Sia/types"

	"github.com/javgh/roadie/blockchain/ethereum"
	"github.com/javgh/roadie/blockchain/sia"
	"github.com/javgh/roadie/trader"
)

type (
	ConsoleFrontend struct {
		similarityPercentage int64
		useExchangeRate      bool
		exchangeRate         Fetcher
	}

	AutoAcceptFrontend struct{}

	RuleBasedFrontend struct {
		absDiffRule  float64
		relDiffRule  float64
		exchangeRate Fetcher
	}

	Frontend interface {
		ApproveOffer(siacoin types.Currency, offer trader.Offer, binding bool) (bool, error)
		CheckSimilarity(a trader.Offer, b trader.Offer) bool
	}

	Fetcher interface {
		Fetch(id string) (*big.Rat, error)
	}
)

func NewConsoleFrontend(similarityPercentage int64, useExchangeRate bool, exchangeRate Fetcher) *ConsoleFrontend {
	frontend := ConsoleFrontend{
		similarityPercentage: similarityPercentage,
		useExchangeRate:      useExchangeRate,
		exchangeRate:         exchangeRate,
	}
	return &frontend
}

func (f *ConsoleFrontend) ApproveOffer(siacoin types.Currency, offer trader.Offer, binding bool) (bool, error) {
	if !offer.Available {
		return false, nil
	}

	antiSpamFeeUSDSegment := ""
	etherUSDSegment := ""
	siacoinUSDSegment := ""
	if f.useExchangeRate {
		usdEther, err := f.exchangeRate.Fetch("ethereum")
		if err != nil {
			return false, err
		}

		usdSiacoin, err := f.exchangeRate.Fetch("siacoin")
		if err != nil {
			return false, err
		}

		antiSpamFeeUSD := ethereum.ApplyRate(&offer.AntiSpamFee, usdEther)
		etherUSD := ethereum.ApplyRate(&offer.Ether, usdEther)
		siacoinUSD := sia.ApplyRate(siacoin, usdSiacoin)

		antiSpamFeeUSDSegment = fmt.Sprintf(" (~ %s)", trader.FormatUSD(antiSpamFeeUSD))
		etherUSDSegment = fmt.Sprintf(" (~ %s)", trader.FormatUSD(etherUSD))
		siacoinUSDSegment = fmt.Sprintf(" (~ %s)", trader.FormatUSD(siacoinUSD))
	}

	fmt.Printf("Best offer received:\n")
	if !binding {
		fmt.Printf("Burn: %s%s\n", ethereum.FormatEther(&offer.AntiSpamFee), antiSpamFeeUSDSegment)
	}
	fmt.Printf("Give: %s%s\n", ethereum.FormatEther(&offer.Ether), etherUSDSegment)
	fmt.Printf("Get : %s%s\n", siacoin.HumanString(), siacoinUSDSegment)
	fmt.Printf("\nThe offer contains the following message:\n")
	fmt.Printf("-----BEGIN MESSAGE-----\n")
	fmt.Println(offer.Msg)
	fmt.Printf("-----END MESSAGE-----\n\n")

	if !binding {
		fmt.Printf("Note that this offer is non-binding. To continue, you will need to burn\n")
		fmt.Printf("the listed anti-spam fee to receive a binding offer. Should the binding offer\n")
		fmt.Printf("be different, you will be prompted again, but the anti-spam fee is non-refundable.\n\n")
	} else {
		fmt.Printf("The other party has indicated that this offer is binding and that they\n")
		fmt.Printf("are ready to proceed with the swap.\n\n")
	}

	fmt.Printf("Press ENTER to continue and accept the offer or CTRL+C to cancel. >")

	var in string
	fmt.Scanln(&in)
	fmt.Println()

	return true, nil
}

func (f *ConsoleFrontend) CheckSimilarity(a trader.Offer, b trader.Offer) bool {
	return trader.CheckSimilarity(a, b, f.similarityPercentage)
}

func (f AutoAcceptFrontend) ApproveOffer(siacoin types.Currency, offer trader.Offer, binding bool) (bool, error) {
	return true, nil
}

func (f AutoAcceptFrontend) CheckSimilarity(a trader.Offer, b trader.Offer) bool {
	return true
}

func NewRuleBasedFrontend(absDiffRule float64, relDiffRule float64, exchangeRate Fetcher) *RuleBasedFrontend {
	frontend := RuleBasedFrontend{
		absDiffRule:  absDiffRule,
		relDiffRule:  relDiffRule,
		exchangeRate: exchangeRate,
	}
	return &frontend
}

func (f *RuleBasedFrontend) ApproveOffer(siacoin types.Currency, offer trader.Offer, binding bool) (bool, error) {
	if !offer.Available {
		return false, nil
	}

	usdEther, err := f.exchangeRate.Fetch("ethereum")
	if err != nil {
		return false, err
	}

	usdSiacoin, err := f.exchangeRate.Fetch("siacoin")
	if err != nil {
		return false, err
	}

	etherTotal := new(big.Int).Add(&offer.Ether, &offer.AntiSpamFee)

	etherTotalUSD := ethereum.ApplyRate(etherTotal, usdEther)
	siacoinUSD := sia.ApplyRate(siacoin, usdSiacoin)

	if f.absDiffRule != 0 {
		absDiff := new(big.Rat).Sub(etherTotalUSD, siacoinUSD)

		absDiffRuleAsRat := new(big.Rat).SetFloat64(f.absDiffRule)
		if absDiff.Cmp(absDiffRuleAsRat) != 1 {
			return true, nil
		}
	}

	if f.relDiffRule != 0 {
		relDiff := new(big.Rat).Quo(etherTotalUSD, siacoinUSD)
		relDiff.Sub(relDiff, new(big.Rat).SetInt64(1))
		relDiff.Mul(relDiff, new(big.Rat).SetInt64(100))

		relDiffRuleAsRat := new(big.Rat).SetFloat64(f.relDiffRule)
		if relDiff.Cmp(relDiffRuleAsRat) != 1 {
			return true, nil
		}
	}

	return false, nil
}

func (f *RuleBasedFrontend) CheckSimilarity(a trader.Offer, b trader.Offer) bool {
	return false
}
