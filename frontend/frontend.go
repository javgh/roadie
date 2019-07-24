package frontend

import (
	"fmt"

	"gitlab.com/NebulousLabs/Sia/types"

	"github.com/javgh/roadie/blockchain/ethereum"
	"github.com/javgh/roadie/blockchain/sia"
	"github.com/javgh/roadie/trader"
)

type (
	ConsoleFrontend struct {
		similarityPercentage int64
		useExchangeRate      bool
		exchangeRate         trader.ExchangeRate
	}

	AutoAcceptFrontend struct{}

	Frontend interface {
		ApproveOffer(siacoin types.Currency, offer trader.Offer, binding bool) (bool, error)
		CheckSimilarity(a trader.Offer, b trader.Offer) bool
	}
)

func NewConsoleFrontend(similarityPercentage int64, useExchangeRate bool) *ConsoleFrontend {
	frontend := ConsoleFrontend{
		similarityPercentage: similarityPercentage,
		useExchangeRate:      useExchangeRate,
		exchangeRate:         trader.NewExchangeRate(),
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
