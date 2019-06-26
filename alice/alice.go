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
	"github.com/javgh/roadie/keypair"
	"github.com/javgh/roadie/rpc"
	"github.com/javgh/roadie/trader"
)

const (
	antiSpamConfirmations = 10
	fundingConfirmations  = 1 //3
	depositConfirmations  = 10
	minTimelockOffset     = 1
	//minTimelockOffset     = types.BlockHeight(24 - 2) // 24 blocks (~ 4 hours) with some leeway
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
	ConsoleFrontend struct {
		exchangeRate trader.ExchangeRate
	}

	AutoAcceptFrontend struct{}

	Frontend interface {
		ApproveOffer(siacoin types.Currency, offer trader.Offer, binding bool) (bool, error)
	}

	confirmationDisplay struct {
		current int64
		total   int64
	}
)

func NewConsoleFrontend() *ConsoleFrontend {
	frontend := ConsoleFrontend{exchangeRate: trader.NewExchangeRate()}
	return &frontend
}

func (f *ConsoleFrontend) ApproveOffer(siacoin types.Currency, offer trader.Offer, binding bool) (bool, error) {
	if !offer.Available {
		return false, nil
	}

	usdEther, err := f.exchangeRate.Fetch("ethereum")
	if err != nil {
		return false, err
	}

	usdSiacoin, err := f.exchangeRate.Fetch("siacoin")
	if err != nil {
		return false, err
	}

	antiSpamFeeUSD := ethereum.ApplyRate(&offer.AntiSpamFee, usdEther)
	etherUSD := ethereum.ApplyRate(&offer.Ether, usdEther)
	siacoinUSD := sia.ApplyRate(siacoin, usdSiacoin)

	fmt.Printf("Best offer received:\n")
	if !binding {
		fmt.Printf("Burn: %s (~ %s)\n", ethereum.FormatEther(&offer.AntiSpamFee), trader.FormatUSD(antiSpamFeeUSD))
	}
	fmt.Printf("Give: %s (~ %s)\n", ethereum.FormatEther(&offer.Ether), trader.FormatUSD(etherUSD))
	fmt.Printf("Get : %s (~ %s)\n", siacoin.HumanString(), trader.FormatUSD(siacoinUSD))
	fmt.Printf("\nThe offer contains the following message:\n")
	fmt.Printf("-----BEGIN MESSAGE-----\n")
	fmt.Println(offer.Msg)
	fmt.Printf("-----END MESSAGE-----\n\n")
	fmt.Printf("USD amounts are based on data from CoinMarketCap.\n\n")

	if !binding {
		fmt.Printf("Note that this offer is non-binding. To continue, you will need to burn\n")
		fmt.Printf("the listed anti-spam fee to receive a binding offer. Should the binding offer\n")
		fmt.Printf("be different, you will be prompted again, but the anti-spam fee is non-refundable.\n\n")
	} else {
		fmt.Printf("The other party has indicated that this offer is binding and that they\n")
		fmt.Printf("are ready to proceed with the swap.\n\n")
	}

	fmt.Printf("Press ENTER to continue and accept the offer or CTRL+C to cancel. >")

	var in string
	fmt.Scanln(&in)
	fmt.Println()

	return true, nil
}

func (f AutoAcceptFrontend) ApproveOffer(siacoin types.Currency, offer trader.Offer, binding bool) (bool, error) {
	return true, nil
}

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

func PerformSwap(siacoin types.Currency, frontend Frontend,
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
	fmt.Printf("Burning anti-spam fee (id %s) and waiting for Ethereum confirmations.\n", antiSpamID.Text(10))
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
