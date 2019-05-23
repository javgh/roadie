package ethereum

import (
	"math/big"

	"github.com/HyperspaceApp/ed25519"
)

type Blockchain interface {
	VerifyAntiSpamPayment(antiSpamID big.Int, antiSpamFee big.Int) (bool, error)
	CheckDeposit(adaptorPubKey ed25519.CurvePoint, ether big.Int, antiSpamID big.Int) (bool, error)
	ClaimDeposit(adaptorPubKey ed25519.CurvePoint, adaptorPrivKey ed25519.Adaptor) error
}
