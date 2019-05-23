package trader

import (
	"math/big"
	"time"

	"gitlab.com/NebulousLabs/Sia/types"
)

type Trader interface {
	PrepareNonBindingOffer(siacoin types.Currency,
		minerFee types.Currency) (msg string, available bool, ether *big.Int, antiSpamFee *big.Int, err error)
	PrepareBindingOffer(siacoin types.Currency, minerFee types.Currency,
		now time.Time) (msg string, available bool, ether *big.Int, deadline *time.Time, err error)
	PauseOrderPreparation(now time.Time)
	ResumeOrderPreparation()
}
