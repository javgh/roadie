package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/HyperspaceApp/ed25519"
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/NebulousLabs/Sia/types"

	"github.com/javgh/roadie/alice"
	"github.com/javgh/roadie/blockchain/ethereum"
	"github.com/javgh/roadie/blockchain/sia"
	"github.com/javgh/roadie/bob"
	"github.com/javgh/roadie/keypair"
	"github.com/javgh/roadie/trader"
)

const (
	defaultClientAddress  = "localhost:9980"
	defaultPasswordFile   = ".sia/apipassword"
	antiSpamConfirmations = 10
	depositConfirmations  = 10
	minTimelockOffset     = 1
	//minTimelockOffset     = types.BlockHeight(24 - 2) // 24 blocks (~ 4 hours) with some leeway
	fundingConfirmations = 1 //3
	jsonRPCEndpoint      = ".ethereum/geth.ipc"
	jsonRPCKeystoreFile  = ".config/roadie/keystore"
)

var (
	oneSiacoin              = types.SiacoinPrecision
	defaultMinerFee         = oneSiacoin
	finney                  = big.NewInt(1e15)
	defaultAntiSpamFee      = big.NewInt(1e14)
	maxAntiSpamID           = new(big.Int).Exp(big.NewInt(2), big.NewInt(64), nil)
	bindingOfferLifetime, _ = time.ParseDuration("1m")
	mockWalletAddress       = common.HexToAddress("0x0000000000000000000000000000000000000000")
	contractAddress         = common.HexToAddress("0x799DF2482f589663d7754451de3FfeF4CAA439c8")
)

type (
	mockTrader struct{}

	confirmationDisplay struct {
		current int64
		total   int64
	}
)

func (mt *mockTrader) PrepareNonBindingOffer(siacoin types.Currency, minerFee types.Currency) (*trader.Offer, error) {
	offer := trader.Offer{
		Msg:         "playground offer",
		Available:   true,
		Ether:       *finney,
		AntiSpamFee: *defaultAntiSpamFee,
	}
	return &offer, nil
}

func (mt *mockTrader) PrepareBindingOffer(siacoin types.Currency, minerFee types.Currency,
	now time.Time) (*trader.Offer, *time.Time, error) {
	offer, err := mt.PrepareNonBindingOffer(siacoin, minerFee)
	if err != nil {
		return nil, nil, err
	}

	deadline := now.Add(bindingOfferLifetime)
	return offer, &deadline, nil
}

func (mt *mockTrader) PauseOrderPreparation(now time.Time) {
}

func (mt *mockTrader) ResumeOrderPreparation() {
}

type mockChain struct {
	adaptorPrivKey ed25519.Adaptor
}

func (mc *mockChain) BurnAntiSpamFee(antiSpamID big.Int, antiSpamFee big.Int) error {
	return nil
}

func (mc *mockChain) CheckAntiSpamConfirmations(antiSpamID big.Int, antiSpamFee big.Int) (int64, error) {
	return antiSpamConfirmations, nil
}

func (mc *mockChain) DepositEther(
	recipient common.Address, adaptorPubKey ed25519.CurvePoint, ether big.Int, antiSpamID big.Int) error {
	return nil
}

func (mc *mockChain) CheckDepositConfirmations(
	recipient common.Address, adaptorPubKey ed25519.CurvePoint, ether big.Int, antiSpamID big.Int) (int64, error) {
	return depositConfirmations, nil
}

func (mc *mockChain) ClaimDeposit(adaptorPrivKey ed25519.Adaptor, antiSpamID big.Int) error {
	mc.adaptorPrivKey = adaptorPrivKey
	return nil
}

func (mc *mockChain) LookupAdaptorPrivKey(adaptorPubKey ed25519.CurvePoint) (bool, *ed25519.Adaptor, error) {
	return true, &mc.adaptorPrivKey, nil
}

func (mc *mockChain) WalletAddress() common.Address {
	return mockWalletAddress
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

func prependHomeDirectory(path string) string {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(currentUser.HomeDir, path)
}

func main() {
	passwordBytes, err := ioutil.ReadFile(prependHomeDirectory(defaultPasswordFile))
	if err != nil {
		log.Fatal(err)
	}
	password := strings.TrimSpace(string(passwordBytes))

	ethChain, err := ethereum.NewGanacheBlockchain()
	//endpoint := prependHomeDirectory(jsonRPCEndpoint)
	//keystoreFile := prependHomeDirectory(jsonRPCKeystoreFile)
	//ethChain, err := ethereum.NewJSONRPCBlockchain(endpoint, keystoreFile, &contractAddress)
	if err != nil {
		log.Fatal(err)
	}

	siaChain, err := sia.NewHTTPAPIBlockchain(defaultClientAddress, password)
	if err != nil {
		log.Fatal(err)
	}
	drSiaChain := sia.NewDryRunBlockchain(*siaChain)

	mockTrader := mockTrader{}
	blacklist := bob.NewBlacklist()
	atomicSwap := bob.NewAtomicSwap(&mockTrader, ethChain, &drSiaChain, blacklist, time.Now())

	nonBindingOffer, err := atomicSwap.RequestNonBindingOffer(oneSiacoin)
	if err != nil {
		log.Fatal(err)
	}

	consoleFrontend := alice.NewConsoleFrontend()
	_, err = consoleFrontend.ApproveOffer(oneSiacoin, *nonBindingOffer, false)
	if err != nil {
		log.Fatal(err)
	}

	antiSpamID, err := rand.Int(rand.Reader, maxAntiSpamID)
	if err != nil {
		log.Fatal(err)
	}
	ethChain.BurnAntiSpamFee(*antiSpamID, nonBindingOffer.AntiSpamFee)
	fmt.Printf("Burned anti-spam fee (id %s) and waiting for Ethereum confirmations.\n", antiSpamID.Text(10))
	confDisplay := confirmationDisplay{current: -1, total: antiSpamConfirmations}
	for {
		confs, err := ethChain.CheckAntiSpamConfirmations(*antiSpamID, nonBindingOffer.AntiSpamFee)
		if err != nil {
			log.Fatal(err)
		}

		confDisplay.show(confs)
		if confs < antiSpamConfirmations {
			time.Sleep(10 * time.Second)
		} else {
			fmt.Printf("\n")
			break
		}
	}

	bindingOffer, err := atomicSwap.RequestBindingOffer(*antiSpamID, time.Now())
	if err != nil {
		log.Fatal(err)
	}

	if bindingOffer.Ether.Cmp(&nonBindingOffer.Ether) != 0 {
		_, err = consoleFrontend.ApproveOffer(oneSiacoin, *bindingOffer, true)
	}

	aliceKeypair, err := keypair.Generate()
	if err != nil {
		log.Fatal(err)
	}

	refundDetails, err := atomicSwap.AcceptOffer(aliceKeypair.PubKey, time.Now())
	if err != nil {
		log.Fatal(err)
	}

	height, err := siaChain.Height()
	if err != nil {
		log.Fatal(err)
	}

	minTimelock := *height + minTimelockOffset
	if refundDetails.Timelock < minTimelock {
		log.Fatal("proposed timelock is too short")
	}

	jointPubKey, jointPrimeKeys, err := ed25519.GenerateJointKey(
		[]ed25519.PublicKey{aliceKeypair.PubKey, refundDetails.BobPubKey})
	if err != nil {
		log.Fatal(err)
	}
	jointUnlockConditions := sia.PubKeyUnlockConditions(jointPubKey)

	refundTx := sia.BuildRefundTransaction(
		refundDetails.FundingOutputID, jointUnlockConditions, refundDetails.BobRefundUnlockHash,
		oneSiacoin, defaultMinerFee, refundDetails.Timelock)

	refundSigHash := sia.WholeSigHash(refundTx, *height)
	aliceRefundNoncePoint := ed25519.GenerateNoncePoint(aliceKeypair.PrivKey, refundSigHash)
	refundSigAlice, err := keypair.JointSignAlice(aliceKeypair, refundDetails.BobPubKey,
		[]ed25519.CurvePoint{aliceRefundNoncePoint, refundDetails.BobRefundNoncePoint}, refundSigHash)
	if err != nil {
		log.Fatal(err)
	}

	fundingTxID, err := atomicSwap.EnableFunding(aliceRefundNoncePoint, refundSigAlice)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Waiting for Sia confirmations for funding transaction %s .\n", fundingTxID)
	confDisplay = confirmationDisplay{current: -1, total: fundingConfirmations}
	for {
		confs, err := drSiaChain.ConfsOfRecentUnlockHash(jointUnlockConditions.UnlockHash(), oneSiacoin.Add(defaultMinerFee))
		if err != nil {
			log.Fatal(err)
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
		log.Fatal(err)
	}
	claimTx := sia.BuildClaimTransaction(
		refundDetails.FundingOutputID, jointUnlockConditions, *aliceClaimUnlockHash,
		oneSiacoin, defaultMinerFee)
	claimSigHash := sia.WholeSigHash(claimTx, *height)
	aliceClaimNoncePoint := ed25519.GenerateNoncePoint(aliceKeypair.PrivKey, claimSigHash)

	adaptorDetails, err := atomicSwap.RequestAdaptorDetails(*aliceClaimUnlockHash, aliceClaimNoncePoint)
	if err != nil {
		log.Fatal(err)
	}

	adaptorSigOK := keypair.VerifyBobsAdaptorSignature(
		jointPrimeKeys, jointPubKey, []ed25519.CurvePoint{aliceClaimNoncePoint, adaptorDetails.BobClaimNoncePoint},
		adaptorDetails.AdaptorPubKey, claimSigHash, adaptorDetails.AdaptorSigBob)
	if !adaptorSigOK {
		log.Fatal("unable to verify adaptor signature")
	}

	ethChain.DepositEther(adaptorDetails.DepositRecipient, adaptorDetails.AdaptorPubKey, bindingOffer.Ether, *antiSpamID)
	fmt.Printf("Deposited payment and waiting for Ethereum confirmations.\n")
	confDisplay = confirmationDisplay{current: -1, total: depositConfirmations}
	for {
		confs, err := ethChain.CheckDepositConfirmations(
			adaptorDetails.DepositRecipient, adaptorDetails.AdaptorPubKey, bindingOffer.Ether, *antiSpamID)
		if err != nil {
			log.Fatal(err)
		}

		confDisplay.show(confs)
		if confs < depositConfirmations {
			time.Sleep(10 * time.Second)
		} else {
			fmt.Printf("\n\n")
			break
		}
	}

	err = atomicSwap.AnnounceDeposit()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Waiting for other party to claim deposit and reveal adaptor secret.\n")
	var ok bool
	var adaptorPrivKey *ed25519.Adaptor
	for {
		ok, adaptorPrivKey, err = ethChain.LookupAdaptorPrivKey(adaptorDetails.AdaptorPubKey)
		if err != nil {
			log.Fatal(err)
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
		log.Fatal(err)
	}

	claimSig := ed25519.AddSignature(adaptorSigAlice, adaptorDetails.AdaptorSigBob)
	claimSig = ed25519.AddSignature(claimSig, append(adaptorDetails.AdaptorPubKey, *adaptorPrivKey...))

	claimSigOK := ed25519.Verify(jointPubKey, claimSigHash, claimSig)
	if !claimSigOK {
		log.Fatal("unable to use adaptor secret to build a valid claim transaction - we were tricked somehow")
	}

	claimTx = sia.AddSignature(claimTx, claimSig)
	err = drSiaChain.BroadcastTransaction(claimTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Swap successfully completed with Sia claim transaction %s .\n", claimTx.ID())

	atomicSwap.Rollback()
}
