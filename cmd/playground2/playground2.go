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
	"gitlab.com/NebulousLabs/Sia/types"

	"github.com/javgh/roadie/alice"
	"github.com/javgh/roadie/blockchain/sia"
	"github.com/javgh/roadie/bob"
	"github.com/javgh/roadie/keypair"
	"github.com/javgh/roadie/trader"
)

const (
	defaultClientAddress  = "localhost:9980"
	defaultPasswordFile   = ".sia/apipassword"
	antiSpamConfirmations = 10
)

var (
	oneSiacoin              = types.SiacoinPrecision
	defaultMinerFee         = oneSiacoin
	finney                  = big.NewInt(1e15)
	defaultAntiSpamFee      = finney
	maxAntiSpamID           = new(big.Int).Exp(big.NewInt(2), big.NewInt(64), nil)
	bindingOfferLifetime, _ = time.ParseDuration("1m")
)

type mockTrader struct{}

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

func (mc *mockChain) CheckAntiSpamConfirmations(antiSpamID big.Int, antiSpamFee big.Int) (int, error) {
	return antiSpamConfirmations, nil
}

func (mc *mockChain) CheckDeposit(adaptorPubKey ed25519.CurvePoint, ether big.Int, antiSpamID big.Int) (bool, error) {
	return true, nil
}

func (mc *mockChain) ClaimDeposit(adaptorPubKey ed25519.CurvePoint, adaptorPrivKey ed25519.Adaptor) error {
	mc.adaptorPrivKey = adaptorPrivKey
	return nil
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

	siaChain, err := sia.NewHTTPAPIBlockchain(defaultClientAddress, password)
	if err != nil {
		log.Fatal(err)
	}
	nbSiaChain := sia.NewNoBroadcastBlockchain(*siaChain)

	mockTrader := mockTrader{}
	mockChain := mockChain{}
	blacklist := bob.NewBlacklist()
	atomicSwap := bob.NewAtomicSwap(&mockTrader, &mockChain, &nbSiaChain, blacklist, time.Now())

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
	mockChain.BurnAntiSpamFee(*antiSpamID, nonBindingOffer.AntiSpamFee)
	fmt.Printf("Burned anti-spam fee (id %s) and waiting for confirmations.\n", antiSpamID.Text(10))
	for {
		confs, err := mockChain.CheckAntiSpamConfirmations(*antiSpamID, nonBindingOffer.AntiSpamFee)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%d/%d", confs, antiSpamConfirmations)
		if confs > antiSpamConfirmations {
			fmt.Printf("... ")
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
	fmt.Println(bindingOffer)

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

	jointPubKey, _ /* jointPrimeKeys */, err := ed25519.GenerateJointKey(
		[]ed25519.PublicKey{aliceKeypair.PubKey, refundDetails.BobPubKey})
	if err != nil {
		log.Fatal(err)
	}
	jointUnlockConditions := sia.PubKeyUnlockConditions(jointPubKey)

	refundTx := sia.BuildRefundTransaction(
		refundDetails.FundingOutputID, jointUnlockConditions, refundDetails.BobRefundUnlockHash,
		oneSiacoin, defaultMinerFee, refundDetails.Timelock)

	height, err := siaChain.Height()
	if err != nil {
		log.Fatal(err)
	}
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
	fmt.Println(fundingTxID)

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

	err = atomicSwap.AnnounceDeposit()
	if err != nil {
		log.Fatal(err)
	}

	noncePoints := []ed25519.CurvePoint{aliceClaimNoncePoint, adaptorDetails.BobClaimNoncePoint}
	adaptorSigAlice, err := keypair.JointSignWithAdaptorAlice(
		aliceKeypair, refundDetails.BobPubKey, noncePoints, adaptorDetails.AdaptorPubKey, claimSigHash)
	if err != nil {
		log.Fatal(err)
	}

	claimSig := ed25519.AddSignature(adaptorSigAlice, adaptorDetails.AdaptorSigBob)
	claimSig = ed25519.AddSignature(claimSig, append(adaptorDetails.AdaptorPubKey, mockChain.adaptorPrivKey...))

	claimSigOk := ed25519.Verify(jointPubKey, claimSigHash, claimSig)
	fmt.Println("Claim sig ok:", claimSigOk)

	claimTx = sia.AddSignature(claimTx, claimSig)
	fmt.Printf("Claim: %s\n", sia.EncodeTransaction(claimTx))

	atomicSwap.Rollback()
}
