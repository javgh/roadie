package trader

import (
	"math/big"
	"time"

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
)
