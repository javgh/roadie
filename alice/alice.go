package alice

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/HyperspaceApp/ed25519"
	"gitlab.com/NebulousLabs/Sia/types"

	"github.com/javgh/roadie/blockchain/ethereum"
	"github.com/javgh/roadie/blockchain/sia"
	"github.com/javgh/roadie/frontend"
	"github.com/javgh/roadie/keypair"
	"github.com/javgh/roadie/rpc"
)

const (
	antiSpamConfirmations = 10
	depositConfirmations  = 10
	minTimelockOffset     = types.BlockHeight(24 - 2) // 24 blocks (~ 4 hours) with some leeway
)

var (
	ErrTimelockTooShort  = errors.New("proposed timelock is too short")
	ErrInvalidAdaptorSig = errors.New("unable to verify adaptor signature")
	ErrInvalidClaimSig   = errors.New(
		"unable to use adaptor secret to build a valid claim transaction - we were tricked somehow")

	maxAntiSpamID   = new(big.Int).Exp(big.NewInt(2), big.NewInt(64), nil)
	defaultMinerFee = types.SiacoinPrecision
)

type (
	confirmationDisplay struct {
		current int64
		total   int64
	}
)

func (d *confirmationDisplay) show(current int64) {
	if d.current == current {
		return
	}

	d.current = current
	fmt.Printf("%d/%d", d.current, d.total)
	if d.current < d.total {
		fmt.Printf(".. ")
	}
}

func PerformSwap(siacoin types.Currency, frontend frontend.Frontend, fundingConfirmations int64,
	ethChain ethereum.Blockchain, siaChain sia.Blockchain, roadieClient *rpc.Client) error {

	id, nonBindingOffer, err := roadieClient.RequestNonBindingOffer(siacoin)
	if err != nil {
		return err
	}

	_, err = frontend.ApproveOffer(siacoin, *nonBindingOffer, false)
	if err != nil {
		return err
	}

	antiSpamID, err := rand.Int(rand.Reader, maxAntiSpamID)
	if err != nil {
		return err
	}
	fmt.Printf("Burning anti-spam fee (id %s) and waiting for Ethereum confirmations.\n", antiSpamID)
	ethChain.BurnAntiSpamFee(*antiSpamID, nonBindingOffer.AntiSpamFee)
	confDisplay := confirmationDisplay{current: -1, total: antiSpamConfirmations}
	for {
		confs, err := ethChain.CheckAntiSpamConfirmations(*antiSpamID, nonBindingOffer.AntiSpamFee)
		if err != nil {
			return err
		}

		confDisplay.show(confs)
		if confs < antiSpamConfirmations {
			time.Sleep(10 * time.Second)
		} else {
			fmt.Printf("\n")
			break
		}
	}

	bindingOffer, err := roadieClient.RequestBindingOffer(*id, *antiSpamID)
	if err != nil {
		return err
	}

	if bindingOffer.Ether.Cmp(&nonBindingOffer.Ether) != 0 {
		_, err = frontend.ApproveOffer(siacoin, *bindingOffer, true)
		if err != nil {
			return err
		}
	}

	aliceKeypair, err := keypair.Generate()
	if err != nil {
		return err
	}

	refundDetails, err := roadieClient.AcceptOffer(*id, aliceKeypair.PubKey)
	if err != nil {
		return err
	}

	height, err := siaChain.Height()
	if err != nil {
		return err
	}

	minTimelock := *height + minTimelockOffset
	if refundDetails.Timelock < minTimelock {
		return ErrTimelockTooShort
	}

	jointPubKey, jointPrimeKeys, err := ed25519.GenerateJointKey(
		[]ed25519.PublicKey{aliceKeypair.PubKey, refundDetails.BobPubKey})
	if err != nil {
		return err
	}
	jointUnlockConditions := sia.PubKeyUnlockConditions(jointPubKey)

	refundTx := sia.BuildRefundTransaction(
		refundDetails.FundingOutputID, jointUnlockConditions, refundDetails.BobRefundUnlockHash,
		siacoin, defaultMinerFee, refundDetails.Timelock)

	refundSigHash := sia.WholeSigHash(refundTx, *height)
	aliceRefundNoncePoint := ed25519.GenerateNoncePoint(aliceKeypair.PrivKey, refundSigHash)
	refundSigAlice, err := keypair.JointSignAlice(aliceKeypair, refundDetails.BobPubKey,
		[]ed25519.CurvePoint{aliceRefundNoncePoint, refundDetails.BobRefundNoncePoint}, refundSigHash)
	if err != nil {
		return err
	}

	fundingTxID, err := roadieClient.EnableFunding(*id, aliceRefundNoncePoint, refundSigAlice)
	if err != nil {
		return err
	}

	fmt.Printf("\nWaiting for Sia confirmations for funding transaction %s .\n", fundingTxID)
	confDisplay = confirmationDisplay{current: -1, total: fundingConfirmations}
	for {
		confs, err := siaChain.ConfsOfRecentUnlockHash(jointUnlockConditions.UnlockHash(), siacoin.Add(defaultMinerFee))
		if err != nil {
			return err
		}

		confDisplay.show(confs)
		if confs < fundingConfirmations {
			time.Sleep(10 * time.Second)
		} else {
			fmt.Printf("\n\n")
			break
		}
	}

	aliceClaimUnlockHash, err := siaChain.NextWalletUnlockHash()
	if err != nil {
		return err
	}
	claimTx := sia.BuildClaimTransaction(
		refundDetails.FundingOutputID, jointUnlockConditions, *aliceClaimUnlockHash,
		siacoin, defaultMinerFee)
	claimSigHash := sia.WholeSigHash(claimTx, *height)
	aliceClaimNoncePoint := ed25519.GenerateNoncePoint(aliceKeypair.PrivKey, claimSigHash)

	adaptorDetails, err := roadieClient.RequestAdaptorDetails(*id, *aliceClaimUnlockHash, aliceClaimNoncePoint)
	if err != nil {
		return err
	}

	adaptorSigOK := keypair.VerifyBobsAdaptorSignature(
		jointPrimeKeys, jointPubKey, []ed25519.CurvePoint{aliceClaimNoncePoint, adaptorDetails.BobClaimNoncePoint},
		adaptorDetails.AdaptorPubKey, claimSigHash, adaptorDetails.AdaptorSigBob)
	if !adaptorSigOK {
		return ErrInvalidAdaptorSig
	}

	fmt.Printf("Depositing payment and waiting for Ethereum confirmations.\n")
	ethChain.DepositEther(adaptorDetails.DepositRecipient, adaptorDetails.AdaptorPubKey, bindingOffer.Ether, *antiSpamID)
	confDisplay = confirmationDisplay{current: -1, total: depositConfirmations}
	for {
		confs, err := ethChain.CheckDepositConfirmations(
			adaptorDetails.DepositRecipient, adaptorDetails.AdaptorPubKey, bindingOffer.Ether, *antiSpamID)
		if err != nil {
			return err
		}

		confDisplay.show(confs)
		if confs < depositConfirmations {
			time.Sleep(10 * time.Second)
		} else {
			fmt.Printf("\n\n")
			break
		}
	}

	fmt.Printf("Announcing deposit and waiting for other party to claim it and reveal adaptor secret.\n")
	err = roadieClient.AnnounceDeposit(*id)
	if err != nil {
		return err
	}

	var ok bool
	var adaptorPrivKey *ed25519.Adaptor
	for {
		ok, adaptorPrivKey, err = ethChain.LookupAdaptorPrivKey(adaptorDetails.AdaptorPubKey)
		if err != nil {
			return err
		}

		if !ok {
			time.Sleep(10 * time.Second)
		} else {
			break
		}
	}

	fmt.Printf("Using adaptor secret to build a valid claim transaction and broadcast it.\n")

	noncePoints := []ed25519.CurvePoint{aliceClaimNoncePoint, adaptorDetails.BobClaimNoncePoint}
	adaptorSigAlice, err := keypair.JointSignWithAdaptorAlice(
		aliceKeypair, refundDetails.BobPubKey, noncePoints, adaptorDetails.AdaptorPubKey, claimSigHash)
	if err != nil {
		return err
	}

	claimSig := ed25519.AddSignature(adaptorSigAlice, adaptorDetails.AdaptorSigBob)
	claimSig = ed25519.AddSignature(claimSig, append(adaptorDetails.AdaptorPubKey, *adaptorPrivKey...))

	claimSigOK := ed25519.Verify(jointPubKey, claimSigHash, claimSig)
	if !claimSigOK {
		return ErrInvalidClaimSig
	}

	claimTx = sia.AddSignature(claimTx, claimSig)
	err = siaChain.BroadcastTransaction(claimTx)
	if err != nil {
		return err
	}

	fmt.Printf("Swap completed successfully with Sia claim transaction %s .\n", claimTx.ID())
	return nil
}

func ReclaimDeposit(ethChain ethereum.Blockchain, antiSpamID big.Int) error {
	fmt.Printf("Reclaiming deposit with id %s.\n", &antiSpamID)
	return ethChain.ReclaimDeposit(antiSpamID)
}
