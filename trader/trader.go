package trader

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
	"gitlab.com/NebulousLabs/Sia/types"

	"github.com/javgh/roadie/blockchain/ethereum"
	"github.com/javgh/roadie/blockchain/sia"
)

type (
	Offer struct {
		Msg         string
		Available   bool
		Ether       big.Int
		AntiSpamFee big.Int
	}

	FixedPremiumTrader struct {
		premiumUSD   *big.Rat
		antiSpamFee  big.Int
		exchangeRate ExchangeRate
		paused       bool
		ethChain     ethereum.Blockchain
		siaChain     sia.Blockchain
	}

	Trader interface {
		PrepareNonBindingOffer(siacoin types.Currency, minerFee types.Currency) (offer *Offer, err error)
		PrepareBindingOffer(siacoin types.Currency, minerFee types.Currency,
			now time.Time) (offer *Offer, deadline *time.Time, err error)
		PauseOrderPreparation(now time.Time)
		ResumeOrderPreparation()
	}

	ExchangeRate struct {
		cache  *cache.Cache
		client *http.Client
	}

	Rate struct {
		ID  string
		USD string `json:"price_usd"`
	}
)

const (
	exchangeRateEndpoint = "https://api.coinmarketcap.com/v1/ticker/"
	formatUSDPrecision   = 4

	msgPaused = "The server is currently in a critical phase of another swap and not ready to make an offer.\n" +
		"Please try again later.\n"
	msgTooSmall = "The minimum amount is %s.\n"
	msgTooLarge = "Insufficient funds to make an offer.\n"
	msgOffer    = "Please note that this offer is also influenced by the current Ethereum gas price of %s.\n"

	gasEstimate = 500000
)

var (
	ErrParsingFailed        = errors.New("unable to parse exchange rate")
	ErrExchangeRateNotFound = errors.New("requested exchange rate not found")

	exchangeRateExpiration, _ = time.ParseDuration("2m")
	exchangeRateInterval, _   = time.ParseDuration("1m")
	httpTimeout, _            = time.ParseDuration("20s")

	minSiacoin = types.SiacoinPrecision
	oneEther   = big.NewRat(1e18, 1)
)

func NewFixedPremiumTrader(premiumUSD *big.Rat, antiSpamFee big.Int,
	ethChain ethereum.Blockchain, siaChain sia.Blockchain) FixedPremiumTrader {
	if premiumUSD == nil {
		premiumUSD = big.NewRat(0, 1)
	}

	return FixedPremiumTrader{
		premiumUSD:   premiumUSD,
		antiSpamFee:  antiSpamFee,
		exchangeRate: NewExchangeRate(),
		paused:       false,
		ethChain:     ethChain,
		siaChain:     siaChain,
	}
}

func (t *FixedPremiumTrader) PrepareNonBindingOffer(siacoin types.Currency, minerFee types.Currency) (*Offer, error) {
	offer := Offer{
		Msg:         "",
		Available:   false,
		Ether:       *big.NewInt(0),
		AntiSpamFee: t.antiSpamFee,
	}

	if t.paused {
		offer.Msg = msgPaused
		return &offer, nil
	}

	if siacoin.Cmp(minSiacoin) == -1 {
		offer.Msg = fmt.Sprintf(msgTooSmall, minSiacoin.HumanString())
		return &offer, nil
	}

	siacoinBalance, err := t.calculateSiacoinBalance()
	if err != nil {
		return nil, err
	}

	if siacoin.Cmp(*siacoinBalance) != -1 {
		offer.Msg = msgTooLarge
		return &offer, nil
	}

	usdEther, err := t.exchangeRate.Fetch("ethereum")
	if err != nil {
		return nil, err
	}

	usdSiacoin, err := t.exchangeRate.Fetch("siacoin")
	if err != nil {
		return nil, err
	}

	siacoinAndFees := siacoin.Add(minerFee).Add(minerFee)
	siacoinAndFeesUSD := sia.ApplyRate(siacoinAndFees, usdSiacoin)
	withPremiumUSD := new(big.Rat).Add(siacoinAndFeesUSD, t.premiumUSD)
	etherRat := new(big.Rat).Mul(new(big.Rat).Quo(withPremiumUSD, usdEther), oneEther)
	ether, _ := new(big.Float).SetRat(etherRat).Int(nil)

	gasPrice, err := t.ethChain.SuggestGasPrice()
	if err != nil {
		return nil, err
	}
	contractCost := new(big.Int).Mul(big.NewInt(gasEstimate), gasPrice)
	ether.Add(ether, contractCost)

	offer.Msg = fmt.Sprintf(msgOffer, ethereum.FormatGwei(gasPrice))
	offer.Available = true
	offer.Ether = *ether

	return &offer, nil
}

func (t *FixedPremiumTrader) calculateSiacoinBalance() (*types.Currency, error) {
	usableOutputs, err := t.siaChain.FetchUsableOutputs()
	if err != nil {
		return nil, err
	}

	balance := types.ZeroCurrency
	for _, usableOutput := range usableOutputs {
		balance = balance.Add(usableOutput.UnspentOutput.Value)
	}

	return &balance, nil
}

func NewExchangeRate() ExchangeRate {
	cache := cache.New(exchangeRateExpiration, exchangeRateInterval)
	client := &http.Client{Timeout: httpTimeout}
	return ExchangeRate{cache: cache, client: client}
}

func (r *ExchangeRate) Fetch(id string) (*big.Rat, error) {
	var rates []Rate

	entry, ok := r.cache.Get("rates")
	if ok {
		rates = entry.([]Rate)
	} else {
		resp, err := r.client.Get(exchangeRateEndpoint)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(body, &rates)
		if err != nil {
			return nil, err
		}

		r.cache.Set("rates", rates, cache.DefaultExpiration)
	}

	for _, rate := range rates {
		if rate.ID == id {
			usd, ok := new(big.Rat).SetString(rate.USD)
			if !ok {
				return nil, ErrParsingFailed
			}

			return usd, nil
		}
	}

	return nil, ErrExchangeRateNotFound
}

func FormatUSD(usd *big.Rat) string {
	return fmt.Sprintf("%s USD", usd.FloatString(formatUSDPrecision))
}
