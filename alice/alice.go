package alice

import (
	"fmt"

	"gitlab.com/NebulousLabs/Sia/types"

	"github.com/javgh/roadie/blockchain/ethereum"
	"github.com/javgh/roadie/blockchain/sia"
	"github.com/javgh/roadie/trader"
)

type (
	ConsoleFrontend struct {
		exchangeRate trader.ExchangeRate
	}

	Frontend interface {
		ApproveOffer(siacoin types.Currency, offer trader.Offer, binding bool) (bool, error)
	}
)

func NewConsoleFrontend() ConsoleFrontend {
	r := trader.NewExchangeRate()
	return ConsoleFrontend{exchangeRate: r}
}

func (f *ConsoleFrontend) ApproveOffer(siacoin types.Currency, offer trader.Offer, binding bool) (bool, error) {
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

	antiSpamFeeUSD := ethereum.ApplyRate(&offer.AntiSpamFee, usdEther)
	etherUSD := ethereum.ApplyRate(&offer.Ether, usdEther)
	siacoinUSD := sia.ApplyRate(siacoin, usdSiacoin)

	fmt.Printf("Best offer received:\n")
	fmt.Printf("Burn: %s (~ %s)\n", ethereum.FormatEther(&offer.AntiSpamFee), trader.FormatUSD(antiSpamFeeUSD))
	fmt.Printf("Give: %s (~ %s)\n", ethereum.FormatEther(&offer.Ether), trader.FormatUSD(etherUSD))
	fmt.Printf("Get : %s (~ %s)\n", siacoin.HumanString(), trader.FormatUSD(siacoinUSD))
	fmt.Printf("\nThe offer contains the following message:\n")
	fmt.Printf("-----BEGIN MESSAGE-----\n")
	fmt.Println(offer.Msg)
	fmt.Printf("-----END MESSAGE-----\n\n")
	fmt.Printf("USD amounts are based on data from CoinMarketCap.\n")

	return false, nil
}
