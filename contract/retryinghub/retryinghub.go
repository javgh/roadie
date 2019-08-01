package retryinghub

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jpillora/backoff"

	contract "github.com/javgh/roadie/contract/hub"
)

const (
	backoffMin     = 10 * time.Second
	backoffMax     = 90 * time.Second
	backoffFactor  = 2
	backoffJitter  = false
	boostFactorNum = 120 // boost in steps of 20 %
	boostFactorDen = 100
)

type (
	Backend interface {
		bind.ContractBackend
		NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
	}

	RetryingHub struct {
		maxGasPrice     big.Int
		boostInterval   time.Duration
		txCheckInterval time.Duration
		backend         Backend
		privKey         ecdsa.PrivateKey
		walletAddress   common.Address
		hub             *contract.Hub
	}

	ServerDetails struct {
		OK     bool
		Target string
		Cert   []byte
	}

	blockchainReader func() (interface{}, error)
	blockchainWriter func(auth *bind.TransactOpts) (*types.Transaction, error)
)

func New(maxGasPrice big.Int, boostInterval time.Duration, txCheckInterval time.Duration,
	backend Backend, privKey ecdsa.PrivateKey, walletAddress common.Address,
	hub *contract.Hub) RetryingHub {
	h := RetryingHub{
		maxGasPrice:     maxGasPrice,
		boostInterval:   boostInterval,
		txCheckInterval: txCheckInterval,
		backend:         backend,
		privKey:         privKey,
		walletAddress:   walletAddress,
		hub:             hub,
	}
	return h
}

func (h *RetryingHub) BurnAntiSpamFee(hashedID [32]byte, value *big.Int, gasLimit uint64) {
	h.robustWrite(func(auth *bind.TransactOpts) (*types.Transaction, error) {
		return h.hub.BurnAntiSpamFee(auth, hashedID)
	}, value, gasLimit)
}

func (h *RetryingHub) CheckAntiSpamConfirmations(id *big.Int, fee *big.Int) *big.Int {
	confs := robustRead(func() (interface{}, error) {
		return h.hub.CheckAntiSpamConfirmations(nil, id, fee)
	})
	return confs.(*big.Int)
}

func (h *RetryingHub) DepositEther(recipient common.Address,
	adaptorPubKey *big.Int, hashedAntiSpamID [32]byte, value *big.Int, gasLimit uint64) {
	h.robustWrite(func(auth *bind.TransactOpts) (*types.Transaction, error) {
		return h.hub.DepositEther(auth, recipient, adaptorPubKey, hashedAntiSpamID)
	}, value, gasLimit)
}

func (h *RetryingHub) CheckDepositConfirmations(recipient common.Address,
	adaptorPubKey *big.Int, value *big.Int, hashedAntiSpamID [32]byte) *big.Int {
	confs := robustRead(func() (interface{}, error) {
		return h.hub.CheckDepositConfirmations(nil, recipient, adaptorPubKey, value, hashedAntiSpamID)
	})
	return confs.(*big.Int)
}

func (h *RetryingHub) ClaimDeposit(adaptorPrivKey *big.Int, antiSpamID *big.Int,
	value *big.Int, gasLimit uint64) {
	h.robustWrite(func(auth *bind.TransactOpts) (*types.Transaction, error) {
		return h.hub.ClaimDeposit(auth, adaptorPrivKey, antiSpamID)
	}, value, gasLimit)
}

func (h *RetryingHub) AdaptorPrivKeys(adaptorPubKey *big.Int) *big.Int {
	adaptorPrivKey := robustRead(func() (interface{}, error) {
		return h.hub.AdaptorPrivKeys(nil, adaptorPubKey)
	})
	return adaptorPrivKey.(*big.Int)
}

func (h *RetryingHub) ReclaimDeposit(hashedID [32]byte, value *big.Int, gasLimit uint64) {
	h.robustWrite(func(auth *bind.TransactOpts) (*types.Transaction, error) {
		return h.hub.ReclaimDeposit(auth, hashedID)
	}, value, gasLimit)
}

func (h *RetryingHub) RegisterServer(target string, cert []byte, value *big.Int, gasLimit uint64) {
	h.robustWrite(func(auth *bind.TransactOpts) (*types.Transaction, error) {
		return h.hub.RegisterServer(auth, target, cert)
	}, value, gasLimit)
}

func (h *RetryingHub) FetchServer(maxAge *big.Int, offset *big.Int) ServerDetails {
	serverDetails := robustRead(func() (interface{}, error) {
		ok, target, cert, err := h.hub.FetchServer(nil, maxAge, offset)
		return ServerDetails{OK: ok, Target: target, Cert: cert}, err
	})
	return serverDetails.(ServerDetails)
}

func (h *RetryingHub) Version() string {
	version := robustRead(func() (interface{}, error) {
		return h.hub.Version(nil)
	})
	return version.(string)
}

func (h *RetryingHub) Deprecated() bool {
	deprecated := robustRead(func() (interface{}, error) {
		return h.hub.Deprecated(nil)
	})
	return deprecated.(bool)
}

func (h *RetryingHub) SuggestGasPrice() *big.Int {
	gasPrice := robustRead(func() (interface{}, error) {
		return h.backend.SuggestGasPrice(context.Background())
	})
	return gasPrice.(*big.Int)
}

func newBackoff() backoff.Backoff {
	b := backoff.Backoff{
		Min:    backoffMin,
		Max:    backoffMax,
		Factor: backoffFactor,
		Jitter: backoffJitter,
	}
	return b
}

func robustRead(reader blockchainReader) interface{} {
	b := newBackoff()

	for {
		result, err := reader()

		if err != nil {
			duration := b.Duration()
			fmt.Printf("%s - retrying in %s\n", err, duration)
			time.Sleep(duration)
		} else {
			return result
		}
	}
}

func (h *RetryingHub) robustWrite(writer blockchainWriter, value *big.Int, gasLimit uint64) {
	b := newBackoff()

	nonceBefore := robustRead(func() (interface{}, error) {
		return h.backend.NonceAt(context.Background(), h.walletAddress, nil)
	})
	nonce := new(big.Int).SetUint64(nonceBefore.(uint64))

	var gasPrice *big.Int
	for {
		auth := bind.NewKeyedTransactor(&h.privKey)
		auth.Value = value
		auth.GasLimit = gasLimit
		auth.Nonce = nonce
		auth.GasPrice = gasPrice

		var tx *types.Transaction
		var err error
		for {
			tx, err = writer(auth)

			if err != nil {
				duration := b.Duration()
				fmt.Printf("%s - retrying in %s\n", err, duration)
				time.Sleep(duration)
			} else {
				break
			}
		}

		boostDeadline := time.Now().Add(h.boostInterval)
		for {
			time.Sleep(h.txCheckInterval)

			nonceNow := robustRead(func() (interface{}, error) {
				return h.backend.NonceAt(context.Background(), h.walletAddress, nil)
			})

			if nonceBefore.(uint64) != nonceNow.(uint64) { // one of our transactions confirmed
				return
			}

			if time.Now().After(boostDeadline) {
				break
			}
		}

		fmt.Printf("Transaction is still pending - boosting gas price\n")
		gasPrice = new(big.Int).Div(
			new(big.Int).Mul(tx.GasPrice(), big.NewInt(boostFactorNum)),
			big.NewInt(boostFactorDen))
		if gasPrice.Cmp(&h.maxGasPrice) == 1 {
			gasPrice = &h.maxGasPrice
		}
	}
}
