package bob

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/HyperspaceApp/ed25519"
	"github.com/patrickmn/go-cache"
	"github.com/satori/go.uuid"
	"gitlab.com/NebulousLabs/Sia/types"

	"github.com/javgh/roadie/blockchain/ethereum"
	"github.com/javgh/roadie/blockchain/sia"
	"github.com/javgh/roadie/keypair"
	"github.com/javgh/roadie/trader"
)

type (
	state int

	AtomicSwap struct {
		id             uuid.UUID
		state          state
		deadline       *time.Time
		siacoin        types.Currency
		ether          big.Int
		antiSpamFee    big.Int
		antiSpamID     big.Int
		bobKeypair     keypair.Keypair
		alicePubKey    ed25519.PublicKey
		jointPubKey    ed25519.PublicKey
		jointPrimeKeys []ed25519.PublicKey
		fundingTx      types.Transaction
		refundTx       types.Transaction
		adaptorPrivKey ed25519.Adaptor
		adaptorPubKey  ed25519.CurvePoint
		trader         trader.Trader
		ethChain       ethereum.Blockchain
		siaChain       sia.Blockchain
		blacklist      Blacklist
	}

	Blacklist struct {
		cache *cache.Cache
	}

	RefundDetails struct {
		BobPubKey           ed25519.PublicKey
		FundingOutputID     types.SiacoinOutputID
		BobRefundUnlockHash types.UnlockHash
		Timelock            types.BlockHeight
		BobRefundNoncePoint ed25519.CurvePoint
	}

	AdaptorDetails struct {
		BobClaimNoncePoint ed25519.CurvePoint
		AdaptorPubKey      ed25519.CurvePoint
		AdaptorSigBob      []byte
	}
)

const (
	stateInitialized state = iota
	stateMadeNonBindingOffer
	stateMadeBindingOffer
	stateOfferAccepted
	stateFunded
	stateProvidedAdaptorDetails
	stateCompleted
	stateAborted

	timelockOffset = types.BlockHeight(0)
)

var (
	ErrWrongState          = errors.New("atomic swap is in a state where this action is not permitted")
	ErrOfferExpired        = errors.New("offer has expired")
	ErrAntiSpamNotDetected = errors.New("no sufficient anti spam payment detected")
	ErrAntiSpamReused      = errors.New("new anti spam payment required")
	ErrInvalidRefundSig    = errors.New("unable to build valid refund transaction")
	ErrInvalidDeposit      = errors.New("no suitable deposit recognized")

	oneSiacoin             = types.SiacoinPrecision
	defaultMinerFee        = oneSiacoin
	atomicSwapLifetime, _  = time.ParseDuration("6h")
	blacklistExpiration, _ = time.ParseDuration("6h")
	blacklistInterval, _   = time.ParseDuration("1h")
)

func NewBlacklist() Blacklist {
	c := cache.New(blacklistExpiration, blacklistInterval)
	return Blacklist{cache: c}
}

func (b *Blacklist) add(id big.Int) {
	b.cache.Set(id.String(), true, cache.DefaultExpiration)
}

func (b *Blacklist) contains(id big.Int) bool {
	_, ok := b.cache.Get(id.String())
	return ok
}

func NewAtomicSwap(trader trader.Trader, ethChain ethereum.Blockchain, siaChain sia.Blockchain,
	blacklist Blacklist, now time.Time) AtomicSwap {
	id := uuid.NewV4()
	deadline := now.Add(atomicSwapLifetime)
	return AtomicSwap{
		id:        id,
		state:     stateInitialized,
		deadline:  &deadline,
		trader:    trader,
		ethChain:  ethChain,
		siaChain:  siaChain,
		blacklist: blacklist,
	}
}

func (s *AtomicSwap) RequestNonBindingOffer(
	siacoin types.Currency) (msg string, available bool, ether *big.Int, antiSpamFee *big.Int, err error) {
	if s.state != stateInitialized {
		return ErrWrongState.Error(), false, nil, nil, ErrWrongState
	}

	msg, available, ether, antiSpamFee, err = s.trader.PrepareNonBindingOffer(siacoin, defaultMinerFee)
	if err != nil {
		return err.Error(), false, nil, nil, err
	}

	s.antiSpamFee = *antiSpamFee
	s.state = stateMadeNonBindingOffer
	return msg, available, ether, antiSpamFee, nil
}

func (s *AtomicSwap) RequestBindingOffer(siacoin types.Currency, antiSpamID big.Int,
	now time.Time) (msg string, available bool, ether *big.Int, err error) {
	if s.state != stateMadeNonBindingOffer {
		return ErrWrongState.Error(), false, nil, ErrWrongState
	}

	msg, available, ether, deadline, err := s.trader.PrepareBindingOffer(siacoin, defaultMinerFee, now)
	if err != nil {
		return err.Error(), false, nil, err
	}

	if !available {
		return msg, available, ether, nil
	}

	if s.blacklist.contains(antiSpamID) {
		return ErrAntiSpamReused.Error(), false, nil, ErrAntiSpamReused
	}

	ok, err := s.ethChain.VerifyAntiSpamPayment(antiSpamID, s.antiSpamFee)
	if err != nil {
		return err.Error(), false, nil, err
	}

	if !ok {
		return ErrAntiSpamNotDetected.Error(), false, nil, ErrAntiSpamNotDetected
	}

	s.blacklist.add(antiSpamID)
	s.trader.PauseOrderPreparation(now)

	s.siacoin = siacoin
	s.ether = *ether
	s.antiSpamID = antiSpamID
	s.deadline = deadline
	s.state = stateMadeBindingOffer
	return msg, available, ether, nil
}

func (s *AtomicSwap) AcceptOffer(alicePubKey ed25519.PublicKey, now time.Time) (*RefundDetails, error) {
	if s.state != stateMadeBindingOffer {
		return nil, ErrWrongState
	}

	if time.Now().After(*s.deadline) {
		return nil, ErrOfferExpired
	}
	newDeadline := now.Add(atomicSwapLifetime)
	s.deadline = &newDeadline

	bobKeypair, err := keypair.Generate()
	if err != nil {
		return nil, err
	}

	s.alicePubKey = alicePubKey
	s.bobKeypair = bobKeypair

	usableOutputs, err := s.siaChain.FetchUsableOutputs()
	if err != nil {
		return nil, err
	}

	walletUnlockHash, err := s.siaChain.NextWalletUnlockHash()
	if err != nil {
		return nil, err
	}

	s.jointPubKey, s.jointPrimeKeys, err = ed25519.GenerateJointKey(
		[]ed25519.PublicKey{s.alicePubKey, s.bobKeypair.PubKey})
	if err != nil {
		return nil, err
	}
	jointUnlockConditions := sia.PubKeyUnlockConditions(s.jointPubKey)
	jointUnlockHash := jointUnlockConditions.UnlockHash()

	value := s.siacoin.Add(defaultMinerFee)

	fundingTx, err := sia.BuildFundingTransaction(
		usableOutputs, *walletUnlockHash, jointUnlockHash, value, defaultMinerFee)
	if err != nil {
		return nil, err
	}
	s.fundingTx = *fundingTx

	walletUnlockHash2, err := s.siaChain.NextWalletUnlockHash()
	if err != nil {
		return nil, err
	}

	height, err := s.siaChain.Height()
	if err != nil {
		return nil, err
	}
	timelock := *height + timelockOffset

	s.refundTx = sia.BuildRefundTransaction(
		fundingTx.SiacoinOutputID(0), jointUnlockConditions, *walletUnlockHash2, s.siacoin, defaultMinerFee, timelock)
	refundSigHash := sia.WholeSigHash(s.refundTx, *height)
	bobRefundNoncePoint := ed25519.GenerateNoncePoint(s.bobKeypair.PrivKey, refundSigHash)

	refundDetails := RefundDetails{
		BobPubKey:           s.bobKeypair.PubKey,
		FundingOutputID:     fundingTx.SiacoinOutputID(0),
		BobRefundUnlockHash: *walletUnlockHash2,
		Timelock:            timelock,
		BobRefundNoncePoint: bobRefundNoncePoint,
	}

	s.state = stateOfferAccepted
	return &refundDetails, nil
}

func (s *AtomicSwap) EnableFunding(aliceRefundNoncePoint ed25519.CurvePoint, refundSigAlice []byte) (*types.TransactionID, error) {
	if s.state != stateOfferAccepted {
		return nil, ErrWrongState
	}

	height, err := s.siaChain.Height()
	if err != nil {
		return nil, err
	}

	refundSigHash := sia.WholeSigHash(s.refundTx, *height)
	bobRefundNoncePoint := ed25519.GenerateNoncePoint(s.bobKeypair.PrivKey, refundSigHash)
	refundSigBob, err := keypair.JointSignBob(s.bobKeypair, s.alicePubKey,
		[]ed25519.CurvePoint{aliceRefundNoncePoint, bobRefundNoncePoint}, refundSigHash)
	if err != nil {
		return nil, err
	}
	refundSig := ed25519.AddSignature(refundSigAlice, refundSigBob)

	refundSigOk := ed25519.Verify(s.jointPubKey, refundSigHash, refundSig)
	if !refundSigOk {
		return nil, ErrInvalidRefundSig
	}
	s.refundTx = sia.AddSignature(s.refundTx, refundSig)

	fundingTxSigned, err := s.siaChain.WalletSign(s.fundingTx)
	if err != nil {
		return nil, err
	}
	s.fundingTx = *fundingTxSigned
	s.siaChain.BroadcastTransaction(s.fundingTx)
	fundingTxID := s.fundingTx.ID()

	s.state = stateFunded
	return &fundingTxID, nil
}

func (s *AtomicSwap) RequestAdaptorDetails(aliceClaimUnlockHash types.UnlockHash,
	aliceClaimNoncePoint ed25519.CurvePoint) (*AdaptorDetails, error) {
	if s.state != stateFunded {
		return nil, ErrWrongState
	}

	height, err := s.siaChain.Height()
	if err != nil {
		return nil, err
	}

	jointUnlockConditions := sia.PubKeyUnlockConditions(s.jointPubKey)
	claimTx := sia.BuildClaimTransaction(
		s.fundingTx.SiacoinOutputID(0), jointUnlockConditions, aliceClaimUnlockHash,
		s.siacoin, defaultMinerFee)
	claimSigHash := sia.WholeSigHash(claimTx, *height)

	s.adaptorPrivKey, s.adaptorPubKey, err = ed25519.GenerateAdaptor(rand.Reader)
	if err != nil {
		return nil, err
	}

	bobClaimNoncePoint := ed25519.GenerateNoncePoint(s.bobKeypair.PrivKey, claimSigHash)
	noncePoints := []ed25519.CurvePoint{aliceClaimNoncePoint, bobClaimNoncePoint}
	adaptorSigBob, err := keypair.JointSignWithAdaptorBob(
		s.bobKeypair, s.alicePubKey, noncePoints, s.adaptorPubKey, claimSigHash)
	if err != nil {
		return nil, err
	}

	adaptorDetails := AdaptorDetails{
		BobClaimNoncePoint: bobClaimNoncePoint,
		AdaptorPubKey:      s.adaptorPubKey,
		AdaptorSigBob:      adaptorSigBob,
	}

	s.state = stateProvidedAdaptorDetails
	return &adaptorDetails, nil
}

func (s *AtomicSwap) AnnounceDeposit() error {
	if s.state != stateProvidedAdaptorDetails {
		return ErrWrongState
	}

	ok, err := s.ethChain.CheckDeposit(s.adaptorPubKey, s.ether, s.antiSpamID)
	if err != nil {
		return err
	}

	if !ok {
		return ErrInvalidDeposit
	}

	err = s.ethChain.ClaimDeposit(s.adaptorPubKey, s.adaptorPrivKey)
	if err != nil {
		return err
	}

	s.state = stateCompleted
	return nil
}

func (s *AtomicSwap) Rollback() {
	if s.state == stateFunded {
		fmt.Printf("Refund: %s\n", sia.EncodeTransaction(s.refundTx))
	}
	s.state = stateAborted
}
