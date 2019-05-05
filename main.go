package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os/user"
	"path/filepath"
	"strings"

	"gitlab.com/NebulousLabs/Sia/crypto"
	"gitlab.com/NebulousLabs/Sia/encoding"
	"gitlab.com/NebulousLabs/Sia/modules"
	"gitlab.com/NebulousLabs/Sia/node/api/client"
	"gitlab.com/NebulousLabs/Sia/types"

	"github.com/HyperspaceApp/ed25519"
)

const (
	defaultClientAddress = "localhost:9980"
	defaultPasswordFile  = ".sia/apipassword"
	timelockOffset       = types.BlockHeight(0)
)

var (
	oneSiacoin      = types.SiacoinPrecision
	twoSiacoins     = oneSiacoin.Mul64(2)
	defaultMinerFee = oneSiacoin
)

type Keypair struct {
	PubKey  ed25519.PublicKey
	PrivKey ed25519.PrivateKey
}

type UsableOutput struct {
	UnspentOutput    modules.UnspentOutput
	UnlockConditions types.UnlockConditions
}

func generateKeypair() (Keypair, error) {
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	return Keypair{PubKey: pubKey, PrivKey: privKey}, err
}

func (k Keypair) unlockConditions() types.UnlockConditions {
	return pubKeyUnlockConditions(k.PubKey)
}

func (k Keypair) unlockHash() types.UnlockHash {
	return k.unlockConditions().UnlockHash()
}

func pubKeyUnlockConditions(pubKey ed25519.PublicKey) types.UnlockConditions {
	siaPublicKey := types.SiaPublicKey{
		Algorithm: types.SignatureEd25519,
		Key:       pubKey[:],
	}
	return types.UnlockConditions{
		PublicKeys:         []types.SiaPublicKey{siaPublicKey},
		SignaturesRequired: 1,
	}
}

func fetchUsableOutputs(httpClient client.Client) ([]UsableOutput, error) {
	unspent, err := httpClient.WalletUnspentGet()
	if err != nil {
		return nil, err
	}

	var usableOutputs []UsableOutput
	for _, unspentOutput := range unspent.Outputs {
		if unspentOutput.FundType != types.SpecifierSiacoinOutput {
			continue
		}

		result, err := httpClient.WalletUnlockConditionsGet(unspentOutput.UnlockHash)
		if err != nil {
			return nil, err
		}

		usableOutputs = append(usableOutputs,
			UsableOutput{UnspentOutput: unspentOutput, UnlockConditions: result.UnlockConditions})
	}

	return usableOutputs, nil
}

func buildFundingTransaction(usableOutputs []UsableOutput, changeUnlockHash types.UnlockHash,
	destinationUnlockHash types.UnlockHash, value types.Currency) (types.Transaction, error) {
	tx := types.Transaction{}
	sum := types.ZeroCurrency
	threshold := value.Add(defaultMinerFee)
	for _, usableOutput := range usableOutputs {
		input := types.SiacoinInput{
			ParentID:         types.SiacoinOutputID(usableOutput.UnspentOutput.ID),
			UnlockConditions: usableOutput.UnlockConditions,
		}
		tx.SiacoinInputs = append(tx.SiacoinInputs, input)

		signature := types.TransactionSignature{
			ParentID:      crypto.Hash(usableOutput.UnspentOutput.ID),
			CoveredFields: types.CoveredFields{WholeTransaction: true},
		}
		tx.TransactionSignatures = append(tx.TransactionSignatures, signature)

		sum = sum.Add(usableOutput.UnspentOutput.Value)
		if sum.Cmp(threshold) == 1 {
			// only stop when strictly above the needed amount; this way we
			// always have a change output and have only one codepath
			break
		}
	}

	if sum.Cmp(threshold) != 1 {
		return tx, fmt.Errorf("Not enough siacoins to send %s", value.HumanString())
	}

	tx.SiacoinOutputs = []types.SiacoinOutput{{
		Value:      value,
		UnlockHash: destinationUnlockHash,
	}, {
		Value:      sum.Sub(value).Sub(defaultMinerFee),
		UnlockHash: changeUnlockHash,
	}}
	tx.MinerFees = []types.Currency{defaultMinerFee}

	return tx, nil
}

func buildRefundTransaction(parentID types.SiacoinOutputID, parentUnlockConditions types.UnlockConditions,
	destinationUnlockHash types.UnlockHash, blockHeight types.BlockHeight) types.Transaction {
	tx := types.Transaction{}
	tx.SiacoinInputs = []types.SiacoinInput{{
		ParentID:         parentID,
		UnlockConditions: parentUnlockConditions,
	}}
	tx.SiacoinOutputs = []types.SiacoinOutput{{
		Value:      oneSiacoin,
		UnlockHash: destinationUnlockHash,
	}}
	tx.MinerFees = []types.Currency{defaultMinerFee}
	tx.TransactionSignatures = []types.TransactionSignature{{
		ParentID:      crypto.Hash(parentID),
		Timelock:      blockHeight + timelockOffset,
		CoveredFields: types.CoveredFields{WholeTransaction: true},
	}}
	return tx
}

func buildClaimTransaction(parentID types.SiacoinOutputID, parentUnlockConditions types.UnlockConditions,
	destinationUnlockHash types.UnlockHash) types.Transaction {
	tx := types.Transaction{}
	tx.SiacoinInputs = []types.SiacoinInput{{
		ParentID:         parentID,
		UnlockConditions: parentUnlockConditions,
	}}
	tx.SiacoinOutputs = []types.SiacoinOutput{{
		Value:      oneSiacoin,
		UnlockHash: destinationUnlockHash,
	}}
	tx.MinerFees = []types.Currency{defaultMinerFee}
	tx.TransactionSignatures = []types.TransactionSignature{{
		ParentID:      crypto.Hash(parentID),
		CoveredFields: types.CoveredFields{WholeTransaction: true},
	}}
	return tx
}

func wholeSigHash(tx types.Transaction, blockHeight types.BlockHeight) []byte {
	sigHash := tx.SigHash(0, blockHeight)
	return sigHash[:]
}

func addSignature(tx types.Transaction, signature []byte) types.Transaction {
	tx.TransactionSignatures[0].Signature = signature
	return tx
}

func buildNoncePoints(aliceKeypair Keypair, bobKeypair Keypair, msg []byte) []ed25519.CurvePoint {
	aliceNoncePoint := ed25519.GenerateNoncePoint(aliceKeypair.PrivKey, msg)
	bobNoncePoint := ed25519.GenerateNoncePoint(bobKeypair.PrivKey, msg)
	return []ed25519.CurvePoint{aliceNoncePoint, bobNoncePoint}
}

func jointSign(aliceKeypair Keypair, bobKeypair Keypair, msg []byte) ([]byte, error) {
	noncePoints := buildNoncePoints(aliceKeypair, bobKeypair, msg)

	jointSigAlice, err := jointSignAlice(aliceKeypair, bobKeypair.PubKey, noncePoints, msg)
	if err != nil {
		return nil, err
	}

	jointSigBob, err := jointSignBob(bobKeypair, aliceKeypair.PubKey, noncePoints, msg)
	if err != nil {
		return nil, err
	}

	return ed25519.AddSignature(jointSigAlice, jointSigBob), nil
}

func jointSignAlice(aliceKeypair Keypair, bobPubKey ed25519.PublicKey,
	noncePoints []ed25519.CurvePoint, msg []byte) ([]byte, error) {
	pubKeys := []ed25519.PublicKey{aliceKeypair.PubKey, bobPubKey}
	aliceN := 0
	jointPrivateKey, err := ed25519.GenerateJointPrivateKey(
		pubKeys, aliceKeypair.PrivKey, aliceN)
	if err != nil {
		return nil, err
	}

	return ed25519.JointSign(
		aliceKeypair.PrivKey, jointPrivateKey, noncePoints, msg), nil
}

func jointSignBob(bobKeypair Keypair, alicePubKey ed25519.PublicKey,
	noncePoints []ed25519.CurvePoint, msg []byte) ([]byte, error) {
	pubKeys := []ed25519.PublicKey{alicePubKey, bobKeypair.PubKey}
	bobN := 1
	jointPrivateKey, err := ed25519.GenerateJointPrivateKey(
		pubKeys, bobKeypair.PrivKey, bobN)
	if err != nil {
		return nil, err
	}

	return ed25519.JointSign(
		bobKeypair.PrivKey, jointPrivateKey, noncePoints, msg), nil
}

func jointSignWithAdaptorBob(bobKeypair Keypair, alicePubKey ed25519.PublicKey,
	noncePoints []ed25519.CurvePoint, adaptorPubKey ed25519.CurvePoint, msg []byte) ([]byte, error) {
	pubKeys := []ed25519.PublicKey{alicePubKey, bobKeypair.PubKey}
	bobN := 1
	jointPrivateKey, err := ed25519.GenerateJointPrivateKey(
		pubKeys, bobKeypair.PrivKey, bobN)
	if err != nil {
		return nil, err
	}

	return ed25519.JointSignWithAdaptor(
		bobKeypair.PrivKey, jointPrivateKey, noncePoints[0], noncePoints[1], adaptorPubKey, msg), nil
}

func verifyAdaptorSignature(bobPrimeKey ed25519.PublicKey, jointPubKey ed25519.PublicKey,
	noncePoints []ed25519.CurvePoint, adaptorPubKey ed25519.CurvePoint, msg []byte, sig []byte) bool {
	return ed25519.VerifyAdaptorSignature(
		bobPrimeKey, jointPubKey, noncePoints[0], noncePoints[1], adaptorPubKey, msg, sig)
}

func jointSignWithAdaptorAlice(aliceKeypair Keypair, bobPubKey ed25519.PublicKey,
	noncePoints []ed25519.CurvePoint, adaptorPubKey ed25519.CurvePoint, msg []byte) ([]byte, error) {
	pubKeys := []ed25519.PublicKey{aliceKeypair.PubKey, bobPubKey}
	aliceN := 0
	jointPrivateKey, err := ed25519.GenerateJointPrivateKey(
		pubKeys, aliceKeypair.PrivKey, aliceN)
	if err != nil {
		return nil, err
	}

	return ed25519.JointSignWithAdaptor(
		aliceKeypair.PrivKey, jointPrivateKey, noncePoints[0], noncePoints[1], adaptorPubKey, msg), nil
}

func broadcastTransaction(httpClient client.Client, tx types.Transaction) error {
	return httpClient.TransactionPoolRawPost(tx, []types.Transaction{})
}

func prependHomeDirectory(path string) string {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(currentUser.HomeDir, path)
}

func main() {
	var httpClient client.Client

	httpClient.Address = defaultClientAddress
	pw, err := ioutil.ReadFile(prependHomeDirectory(defaultPasswordFile))
	if err != nil {
		log.Fatal(err)
	}
	httpClient.Password = strings.TrimSpace(string(pw))

	status, err := httpClient.WalletGet()
	if err != nil {
		log.Fatal(err)
	}

	if !status.Unlocked {
		log.Fatal("Please unlock wallet to continue")
	}

	fmt.Printf("Confirmed siacoin balance: %s\n", status.ConfirmedSiacoinBalance.HumanString())
	fmt.Printf("Height: %d\n", status.Height)

	aliceKeypair, err := generateKeypair()
	if err != nil {
		log.Fatal(err)
	}

	bobKeypair, err := generateKeypair()
	if err != nil {
		log.Fatal(err)
	}

	jointPubKey, primeKeys, err := ed25519.GenerateJointKey([]ed25519.PublicKey{aliceKeypair.PubKey, bobKeypair.PubKey})
	jointUnlockConditions := pubKeyUnlockConditions(jointPubKey)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Joint unlock hash: %s\n", jointUnlockConditions.UnlockHash())

	//msg := []byte("this is a test")
	//jointSig, err := jointSign(aliceKeypair, bobKeypair, msg)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(ed25519.Verify(jointPubKey, msg, jointSig))

	usableOutputs, err := fetchUsableOutputs(httpClient)
	if err != nil {
		log.Fatal(err)
	}

	change, err := httpClient.WalletAddressGet()
	if err != nil {
		log.Fatal(err)
	}

	tx, err := buildFundingTransaction(usableOutputs, change.Address, jointUnlockConditions.UnlockHash(), twoSiacoins)
	if err != nil {
		log.Fatal(err)
	}

	result, err := httpClient.WalletSignPost(tx, []crypto.Hash{})
	tx = result.Transaction
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Funding tx encoded: %s\n", base64.StdEncoding.EncodeToString(encoding.Marshal(tx)))

	result2, err := httpClient.WalletAddressGet()
	if err != nil {
		log.Fatal(err)
	}
	refundTx := buildRefundTransaction(tx.SiacoinOutputID(0), jointUnlockConditions, result2.Address, status.Height)
	refundTxSigHash := wholeSigHash(refundTx, status.Height)
	refundTxSig, err := jointSign(aliceKeypair, bobKeypair, refundTxSigHash)
	if err != nil {
		log.Fatal(err)
	}
	refundTx = addSignature(refundTx, refundTxSig)
	fmt.Printf("Refund tx encoded: %s\n", base64.StdEncoding.EncodeToString(encoding.Marshal(refundTx)))

	claimTx := buildClaimTransaction(tx.SiacoinOutputID(0), jointUnlockConditions, result2.Address)
	claimTxSighHash := wholeSigHash(claimTx, status.Height)

	adaptorPrivKey, adaptorPubKey, err := ed25519.GenerateAdaptor(rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	noncePoints := buildNoncePoints(aliceKeypair, bobKeypair, claimTxSighHash)
	adaptorSigBob, err := jointSignWithAdaptorBob(bobKeypair, aliceKeypair.PubKey, noncePoints, adaptorPubKey, claimTxSighHash)
	if err != nil {
		log.Fatal(err)
	}

	adaptorSigOk := verifyAdaptorSignature(primeKeys[1], jointPubKey, noncePoints, adaptorPubKey, claimTxSighHash, adaptorSigBob)
	fmt.Printf("Adaptor sig verified: %t\n", adaptorSigOk)

	adaptorSigAlice, err := jointSignWithAdaptorAlice(aliceKeypair, bobKeypair.PubKey, noncePoints, adaptorPubKey, claimTxSighHash)
	if err != nil {
		log.Fatal(err)
	}
	claimTxSig := ed25519.AddSignature(adaptorSigAlice, adaptorSigBob)
	claimTxSig = ed25519.AddSignature(claimTxSig, append(adaptorPubKey, adaptorPrivKey...))

	claimTxSigOk := ed25519.Verify(jointPubKey, claimTxSighHash, claimTxSig)
	fmt.Printf("Claim tx sig verified: %t\n", claimTxSigOk)

	claimTx = addSignature(claimTx, claimTxSig)
	fmt.Printf("Claim tx encoded: %s\n", base64.StdEncoding.EncodeToString(encoding.Marshal(claimTx)))

	//err = broadcastTransaction(httpClient, tx)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//err = broadcastTransaction(httpClient, refundTx)
	//if err != nil {
	//	log.Fatal(err)
	//}

	// Approach 1:
	// Alice: think of adaptor
	// Alice: lock ether behind adaptor
	// Alice: ed25519.JointSignWithAdaptor
	// Bob: VerifyAdaptorSignature, JointSignWithAdaptor
	// Alice: combine signatures, add adaptor, send tx
	// Bob: notice tx, extract adaptor, claim ether
	//
	// Approach 2:
	// Bob: think of adaptor, JointSignWithAdaptor
	// Alice: VerifyAdaptorSignature, lock ether behind adaptor
	// Bob: claim ether
	// Alice: see adaptor, JointSignWithAdaptor, combine signatures, add adaptor, send tx
}
