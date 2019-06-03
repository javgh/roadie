package ethereum

import (
	"fmt"
	"math/big"

	"github.com/HyperspaceApp/ed25519"
)

const (
	formatEtherPrecision = 6
)

var (
	oneEther = big.NewInt(1e18)
)

type Blockchain interface {
	BurnAntiSpamFee(antiSpamID big.Int, antiSpamFee big.Int) error
	CheckAntiSpamConfirmations(antiSpamID big.Int, antiSpamFee big.Int) (int, error)
	DepositEther(adaptorPubKey ed25519.CurvePoint, ether big.Int, antiSpamID big.Int) error
	CheckDepositConfirmations(adaptorPubKey ed25519.CurvePoint, ether big.Int, antiSpamID big.Int) (int, error)
	ClaimDeposit(adaptorPubKey ed25519.CurvePoint, adaptorPrivKey ed25519.Adaptor) error
}

func FormatEther(ether *big.Int) string {
	r := big.NewRat(0, 1).SetFrac(ether, oneEther)
	return fmt.Sprintf("%s ETH", r.FloatString(formatEtherPrecision))
}

func ApplyRate(ether *big.Int, rate *big.Float) *big.Float {
	r := new(big.Rat).SetFrac(ether, oneEther)
	f := new(big.Float).Mul(new(big.Float).SetRat(r), rate)
	return f
}
