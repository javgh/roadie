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

func GenerateKeypair() Keypair {
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		log.Fatal(err)
	}

	return Keypair{PubKey: pubKey, PrivKey: privKey}
}

func (k Keypair) SiaPublicKey() types.SiaPublicKey {
	return types.SiaPublicKey{
		Algorithm: types.SignatureEd25519,
		Key:       k.PubKey[:],
	}
}

func (k Keypair) UnlockConditions() types.UnlockConditions {
	return types.UnlockConditions{
		PublicKeys:         []types.SiaPublicKey{k.SiaPublicKey()},
		SignaturesRequired: 1,
	}
}

func (k Keypair) UnlockHash() types.UnlockHash {
	return k.UnlockConditions().UnlockHash()
}

//func playgroundUnlockHash() types.UnlockHash {
//	_, publicKey1 := crypto.GenerateKeyPair()
//	siaPublicKey1 := types.Ed25519PublicKey(publicKey1)
//
//	_, publicKey2 := crypto.GenerateKeyPair()
//	siaPublicKey2 := types.Ed25519PublicKey(publicKey2)
//
//	unlockConditions := types.UnlockConditions{
//		PublicKeys:         []types.SiaPublicKey{siaPublicKey1, siaPublicKey2},
//		SignaturesRequired: 1,
//	}
//
//	return unlockConditions.UnlockHash()
//}

//func playgroundKeysAndUnlockConditions() (crypto.SecretKey, crypto.PublicKey, types.UnlockConditions) {
//	privKey, pubKey := crypto.GenerateKeyPair()
//	siaPubKey := types.Ed25519PublicKey(pubKey)
//	unlockConditions := types.UnlockConditions{
//		PublicKeys:         []types.SiaPublicKey{siaPubKey},
//		SignaturesRequired: 1,
//	}
//	return privKey, pubKey, unlockConditions
//}

func buildFundingTransaction(httpClient client.Client,
	destinationUnlockHash types.UnlockHash, value types.Currency) (types.Transaction, error) {
	tx := types.Transaction{}

	unspent, err := httpClient.WalletUnspentGet()
	if err != nil {
		return tx, err
	}

	change, err := httpClient.WalletAddressGet()
	if err != nil {
		return tx, err
	}

	sum := types.ZeroCurrency
	threshold := value.Add(defaultMinerFee)
	for _, output := range unspent.Outputs {
		if output.FundType != types.SpecifierSiacoinOutput {
			continue
		}

		result, err := httpClient.WalletUnlockConditionsGet(output.UnlockHash)
		if err != nil {
			return tx, err
		}

		input := types.SiacoinInput{
			ParentID:         types.SiacoinOutputID(output.ID),
			UnlockConditions: result.UnlockConditions,
		}
		tx.SiacoinInputs = append(tx.SiacoinInputs, input)

		signature := types.TransactionSignature{
			ParentID:      crypto.Hash(output.ID),
			CoveredFields: types.CoveredFields{WholeTransaction: true},
		}
		tx.TransactionSignatures = append(tx.TransactionSignatures, signature)

		sum = sum.Add(output.Value)
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
		UnlockHash: change.Address,
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
	keypair Keypair) types.Transaction {
	wholeSigHash := tx.SigHash(0, blockHeight)
	sig := ed25519.Sign(keypair.PrivKey, wholeSigHash[:])
	tx.TransactionSignatures[0].Signature = sig[:]
	return tx
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

	//playPrivKey, playPubKey, playUnlockConditions := playgroundKeysAndUnlockConditions()
	playKeypair := GenerateKeypair()
	fmt.Printf("Playground public key: %x\n", playKeypair.PubKey)
	fmt.Printf("Playground private key: %x\n", playKeypair.PrivKey)
	fmt.Printf("Playground unlock hash: %s\n", playKeypair.UnlockHash())

	//result, err := httpClient.WalletSiacoinsPost(oneSiacoin, unlockHash)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(result)

	//var hash crypto.Hash
	//hash.LoadString("de6a8f1dfa2db63156d9b0397d384694a8d99a82854aed358e974e1bfdb36436")
	//fmt.Println(hash)
	//ptx, err := httpClient.WalletTransactionGet(types.TransactionID(hash))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(ptx)

	tx, err := buildFundingTransaction(httpClient, playKeypair.UnlockHash(), twoSiacoins)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tx)

	result, err := httpClient.WalletSignPost(tx, []crypto.Hash{})
	tx = result.Transaction
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tx)

	fmt.Printf("wholeSigHash: %s\n", tx.SigHash(0, status.Height))
	fmt.Printf("funding output ID: %s\n", tx.SiacoinOutputID(0))

	result2, err := httpClient.WalletAddressGet()
	if err != nil {
		log.Fatal(err)
	}
	refundTx := buildSimpleRefundTransaction(tx.SiacoinOutputID(0), playKeypair.UnlockConditions(), result2.Address, status.Height)
	fmt.Println(refundTx)

	refundTx = signRefundTransaction(refundTx, status.Height, playKeypair)
	fmt.Println(refundTx)

	fmt.Printf("tx encoded: %s\n", base64.StdEncoding.EncodeToString(encoding.Marshal(tx)))
	fmt.Printf("refund tx encoded: %s\n", base64.StdEncoding.EncodeToString(encoding.Marshal(refundTx)))

	//err = broadcastTransaction(httpClient, tx)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//err = broadcastTransaction(httpClient, refundTx)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//for _, output := range unspent.Outputs {
	//	if output.FundType != types.SpecifierSiacoinOutput {
	//		continue
	//	}

	//	fmt.Println(output)
	//}

	//result, err := httpClient.WalletAddressGet()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("Fresh wallet address: %s\n", result.Address)

	// ed25519.GenerateKey(nil)
	// asSiaPublicKey -> construct SiaPublicKey
}
