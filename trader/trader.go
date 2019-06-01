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
)

type (
	Offer struct {
		Msg         string
		Available   bool
		Ether       big.Int
		AntiSpamFee big.Int
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
)

var (
	ErrParsingFailed        = errors.New("unable to parse exchange rate")
	ErrExchangeRateNotFound = errors.New("requested exchange rate not found")

	exchangeRateExpiration, _ = time.ParseDuration("2m")
	exchangeRateInterval, _   = time.ParseDuration("1m")
	httpTimeout, _            = time.ParseDuration("20s")
)

func NewExchangeRate() ExchangeRate {
	cache := cache.New(exchangeRateExpiration, exchangeRateInterval)
	client := &http.Client{Timeout: httpTimeout}
	return ExchangeRate{cache: cache, client: client}
}

func (r *ExchangeRate) Fetch(id string) (*big.Float, error) {
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
			usd, ok := new(big.Float).SetString(rate.USD)
			if !ok {
				return nil, ErrParsingFailed
			}

			return usd, nil
		}
	}

	return nil, ErrExchangeRateNotFound
}

func FormatUSD(usd *big.Float) string {
	return fmt.Sprintf("%s USD", usd.Text('f', formatUSDPrecision))
}
