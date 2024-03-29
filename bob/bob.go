package bob

import (
	"crypto/rand"
	"errors"
	"math/big"
	"time"

	"github.com/HyperspaceApp/ed25519"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"gitlab.com/NebulousLabs/Sia/types"

	"github.com/javgh/roadie/blockchain/ethereum"
	"github.com/javgh/roadie/blockchain/sia"
	"github.com/javgh/roadie/keypair"
	"github.com/javgh/roadie/trader"
)

type (
	state int

	AtomicSwap struct {
		ID             uuid.UUID
		state          state
		deadline       time.Time
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
		DepositRecipient   common.Address
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
	stateRefunded
	stateAborted

	timelockOffset        = types.BlockHeight(24) // 24 blocks (~ 4 hours)
	antiSpamConfirmations = 8
	depositConfirmations  = 8
)

var (
	ErrWrongState          = errors.New("atomic swap is in a state where this action is not permitted")
	ErrOfferExpired        = errors.New("offer has expired")
	ErrAntiSpamNotDetected = errors.New("no sufficient anti spam payment detected")
	ErrAntiSpamReused      = errors.New("new anti spam payment required")
	ErrInvalidRefundSig    = errors.New("unable to build valid refund transaction")
	ErrInvalidDeposit      = errors.New("no suitable deposit recognized")

	defaultMinerFee        = types.SiacoinPrecision
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
	blacklist Blacklist, now time.Time) *AtomicSwap {
	id := uuid.Must(uuid.NewRandom())
	deadline := now.Add(atomicSwapLifetime)
	atomicSwap := AtomicSwap{
		ID:        id,
		state:     stateInitialized,
		deadline:  deadline,
		trader:    trader,
		ethChain:  ethChain,
		siaChain:  siaChain,
		blacklist: blacklist,
	}
	return &atomicSwap
}

func (s *AtomicSwap) RequestNonBindingOffer(siacoin types.Currency, now time.Time) (*trader.Offer, error) {
	if s.state != stateInitialized {
		return nil, ErrWrongState
	}

	offer, err := s.trader.PrepareNonBindingOffer(siacoin, defaultMinerFee, now)
	if err != nil {
		return nil, err
	}

	s.siacoin = siacoin
	s.antiSpamFee = offer.AntiSpamFee
	s.state = stateMadeNonBindingOffer
	return offer, nil
}

func (s *AtomicSwap) RequestBindingOffer(antiSpamID big.Int, now time.Time) (*trader.Offer, error) {
	if s.state != stateMadeNonBindingOffer {
		return nil, ErrWrongState
	}

	offer, deadline, err := s.trader.PrepareBindingOffer(s.siacoin, defaultMinerFee, now)
	if err != nil {
		return nil, err
	}

	if !offer.Available {
		return offer, nil
	}

	if s.blacklist.contains(antiSpamID) {
		return nil, ErrAntiSpamReused
	}

	confs, err := s.ethChain.CheckAntiSpamConfirmations(antiSpamID, s.antiSpamFee)
	if err != nil {
		return nil, err
	}

	if confs < antiSpamConfirmations {
		return nil, ErrAntiSpamNotDetected
	}

	s.blacklist.add(antiSpamID)
	s.trader.PauseOrderPreparation(now)

	s.ether = offer.Ether
	s.antiSpamID = antiSpamID
	s.deadline = *deadline
	s.state = stateMadeBindingOffer
	return offer, nil
}

func (s *AtomicSwap) AcceptOffer(alicePubKey ed25519.PublicKey, now time.Time) (*RefundDetails, error) {
	if s.state != stateMadeBindingOffer {
		return nil, ErrWrongState
	}

	if now.After(s.deadline) {
		return nil, ErrOfferExpired
	}
	s.deadline = now.Add(atomicSwapLifetime)

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

	refundSigOK := ed25519.Verify(s.jointPubKey, refundSigHash, refundSig)
	if !refundSigOK {
		return nil, ErrInvalidRefundSig
	}
	s.refundTx = sia.AddSignature(s.refundTx, refundSig)

	fundingTxSigned, err := s.siaChain.WalletSign(s.fundingTx)
	if err != nil {
		return nil, err
	}
	s.fundingTx = *fundingTxSigned
	err = s.siaChain.BroadcastTransaction(s.fundingTx)
	if err != nil {
		return nil, err
	}
	fundingTxID := s.fundingTx.ID()
	s.trader.ResumeOrderPreparation()

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

	depositRecipient := s.ethChain.WalletAddress()
	adaptorDetails := AdaptorDetails{
		BobClaimNoncePoint: bobClaimNoncePoint,
		AdaptorPubKey:      s.adaptorPubKey,
		AdaptorSigBob:      adaptorSigBob,
		DepositRecipient:   depositRecipient,
	}

	s.state = stateProvidedAdaptorDetails
	return &adaptorDetails, nil
}

func (s *AtomicSwap) AnnounceDeposit() error {
	if s.state != stateProvidedAdaptorDetails {
		return ErrWrongState
	}

	depositRecipient := s.ethChain.WalletAddress()
	confs, err := s.ethChain.CheckDepositConfirmations(depositRecipient, s.adaptorPubKey, s.ether, s.antiSpamID)
	if err != nil {
		return err
	}

	if confs < depositConfirmations {
		return ErrInvalidDeposit
	}

	err = s.ethChain.ClaimDeposit(s.adaptorPrivKey, s.antiSpamID)
	if err != nil {
		return err
	}

	s.state = stateCompleted
	return nil
}

func (s *AtomicSwap) Check(now time.Time) (noLongerNeeded bool, maybeRefundTxID *types.TransactionID, err error) {
	if now.After(s.deadline) {
		if s.state == stateInitialized || s.state == stateMadeNonBindingOffer ||
			s.state == stateMadeBindingOffer || s.state == stateOfferAccepted {
			s.state = stateAborted
		} else if s.state == stateFunded || s.state == stateProvidedAdaptorDetails {
			err = s.siaChain.BroadcastTransaction(s.refundTx)
			if err != nil {
				return false, nil, err
			}
			s.state = stateRefunded
			refundTxID := s.refundTx.ID()
			maybeRefundTxID = &refundTxID
		}
	}

	if now.After(s.deadline.Add(atomicSwapLifetime)) { // wait extra lifetime before it is safe to forget about it
		noLongerNeeded = true
	}

	return noLongerNeeded, maybeRefundTxID, nil
}

func (s *AtomicSwap) StateText() string {
	switch s.state {
	case stateInitialized:
		return "stateInitialized"
	case stateMadeNonBindingOffer:
		return "stateMadeNonBindingOffer"
	case stateMadeBindingOffer:
		return "stateMadeBindingOffer"
	case stateOfferAccepted:
		return "stateOfferAccepted"
	case stateFunded:
		return "stateFunded"
	case stateProvidedAdaptorDetails:
		return "stateProvidedAdaptorDetails"
	case stateCompleted:
		return "stateCompleted"
	case stateRefunded:
		return "stateRefunded"
	default:
		return "stateAborted"
	}
}

func (s *AtomicSwap) EncodedRefundTransaction() (string, bool) {
	if s.state == stateFunded || s.state == stateProvidedAdaptorDetails ||
		s.state == stateCompleted || s.state == stateRefunded {
		return sia.EncodeTransaction(s.refundTx), true
	}

	return "", false
}
