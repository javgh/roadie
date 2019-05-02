package main

import (
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

func buildSimpleRefundTransaction(parentID types.SiacoinOutputID, parentUnlockConditions types.UnlockConditions,
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

func signRefundTransaction(tx types.Transaction, blockHeight types.BlockHeight,
	aliceKeypair Keypair, bobKeypair Keypair) types.Transaction {
	wholeSigHash := tx.SigHash(0, blockHeight)
	sig, _ := jointSign(aliceKeypair, bobKeypair, wholeSigHash[:])
	tx.TransactionSignatures[0].Signature = sig[:]
	return tx
}

func jointSign(aliceKeypair Keypair, bobKeypair Keypair, msg []byte) ([]byte, error) {
	aliceNoncePoint := ed25519.GenerateNoncePoint(aliceKeypair.PrivKey, msg)
	bobNoncePoint := ed25519.GenerateNoncePoint(bobKeypair.PrivKey, msg)
	noncePoints := []ed25519.CurvePoint{aliceNoncePoint, bobNoncePoint}

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

	jointPubKey, _, err := ed25519.GenerateJointKey([]ed25519.PublicKey{aliceKeypair.PubKey, bobKeypair.PubKey})
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

	result2, err := httpClient.WalletAddressGet()
	if err != nil {
		log.Fatal(err)
	}
	refundTx := buildSimpleRefundTransaction(tx.SiacoinOutputID(0), jointUnlockConditions,
		result2.Address, status.Height)

	refundTx = signRefundTransaction(refundTx, status.Height, aliceKeypair, bobKeypair)

	fmt.Printf("Funding tx encoded: %s\n", base64.StdEncoding.EncodeToString(encoding.Marshal(tx)))
	fmt.Printf("Refund tx encoded: %s\n", base64.StdEncoding.EncodeToString(encoding.Marshal(refundTx)))

	//err = broadcastTransaction(httpClient, tx)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//err = broadcastTransaction(httpClient, refundTx)
	//if err != nil {
	//	log.Fatal(err)
	//}
}
