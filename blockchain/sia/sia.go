package sia

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/HyperspaceApp/ed25519"
	"gitlab.com/NebulousLabs/Sia/crypto"
	"gitlab.com/NebulousLabs/Sia/encoding"
	"gitlab.com/NebulousLabs/Sia/modules"
	"gitlab.com/NebulousLabs/Sia/node/api/client"
	"gitlab.com/NebulousLabs/Sia/types"
)

type (
	HTTPAPIBlockchain struct {
		httpClient client.Client
	}

	NoBroadcastBlockchain struct {
		chain HTTPAPIBlockchain
	}

	Blockchain interface {
		FetchUsableOutputs() ([]UsableOutput, error)
		NextWalletUnlockHash() (*types.UnlockHash, error)
		Height() (*types.BlockHeight, error)
		WalletSign(tx types.Transaction) (*types.Transaction, error)
		BroadcastTransaction(tx types.Transaction) error
	}

	UsableOutput struct {
		UnspentOutput    modules.UnspentOutput
		UnlockConditions types.UnlockConditions
	}
)

var (
	ErrWalletLocked      = errors.New("wallet is locked")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

func NewHTTPAPIBlockchain(address string, password string) (*HTTPAPIBlockchain, error) {
	c := HTTPAPIBlockchain{}
	c.httpClient.Address = address
	c.httpClient.Password = password

	status, err := c.httpClient.WalletGet()
	if err != nil {
		return nil, err
	}

	if !status.Unlocked {
		return nil, ErrWalletLocked
	}

	return &c, nil
}

func (c *HTTPAPIBlockchain) FetchUsableOutputs() ([]UsableOutput, error) {
	unspent, err := c.httpClient.WalletUnspentGet()
	if err != nil {
		return nil, err
	}

	var usableOutputs []UsableOutput
	for _, unspentOutput := range unspent.Outputs {
		if unspentOutput.FundType != types.SpecifierSiacoinOutput {
			continue
		}

		result, err := c.httpClient.WalletUnlockConditionsGet(unspentOutput.UnlockHash)
		if err != nil {
			return nil, err
		}

		usableOutputs = append(usableOutputs,
			UsableOutput{UnspentOutput: unspentOutput, UnlockConditions: result.UnlockConditions})
	}

	return usableOutputs, nil
}

func (c *HTTPAPIBlockchain) NextWalletUnlockHash() (*types.UnlockHash, error) {
	result, err := c.httpClient.WalletAddressGet()
	if err != nil {
		return nil, err
	}

	return &result.Address, nil
}

func (c *HTTPAPIBlockchain) Height() (*types.BlockHeight, error) {
	status, err := c.httpClient.WalletGet()
	if err != nil {
		return nil, err
	}

	return &status.Height, nil
}

func (c *HTTPAPIBlockchain) WalletSign(tx types.Transaction) (*types.Transaction, error) {
	result, err := c.httpClient.WalletSignPost(tx, []crypto.Hash{})
	if err != nil {
		return nil, err
	}

	return &result.Transaction, nil
}

func (c *HTTPAPIBlockchain) BroadcastTransaction(tx types.Transaction) error {
	return c.httpClient.TransactionPoolRawPost(tx, []types.Transaction{})
}

func NewNoBroadcastBlockchain(chain HTTPAPIBlockchain) NoBroadcastBlockchain {
	return NoBroadcastBlockchain{chain: chain}
}

func (c *NoBroadcastBlockchain) FetchUsableOutputs() ([]UsableOutput, error) {
	return c.chain.FetchUsableOutputs()
}

func (c *NoBroadcastBlockchain) NextWalletUnlockHash() (*types.UnlockHash, error) {
	return c.chain.NextWalletUnlockHash()
}

func (c *NoBroadcastBlockchain) Height() (*types.BlockHeight, error) {
	return c.chain.Height()
}

func (c *NoBroadcastBlockchain) WalletSign(tx types.Transaction) (*types.Transaction, error) {
	return c.chain.WalletSign(tx)
}

func (c *NoBroadcastBlockchain) BroadcastTransaction(tx types.Transaction) error {
	fmt.Printf("Skipping broadcast for: %s\n", EncodeTransaction(tx))
	return nil
}

func PubKeyUnlockConditions(pubKey ed25519.PublicKey) types.UnlockConditions {
	siaPublicKey := types.SiaPublicKey{
		Algorithm: types.SignatureEd25519,
		Key:       pubKey[:],
	}
	return types.UnlockConditions{
		PublicKeys:         []types.SiaPublicKey{siaPublicKey},
		SignaturesRequired: 1,
	}
}

func BuildFundingTransaction(usableOutputs []UsableOutput, changeUnlockHash types.UnlockHash,
	destinationUnlockHash types.UnlockHash, value types.Currency, minerFee types.Currency) (*types.Transaction, error) {
	tx := types.Transaction{}
	sum := types.ZeroCurrency
	threshold := value.Add(minerFee)
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
		return nil, ErrInsufficientFunds
	}

	tx.SiacoinOutputs = []types.SiacoinOutput{{
		Value:      value,
		UnlockHash: destinationUnlockHash,
	}, {
		Value:      sum.Sub(value).Sub(minerFee),
		UnlockHash: changeUnlockHash,
	}}
	tx.MinerFees = []types.Currency{minerFee}

	return &tx, nil
}

func BuildRefundTransaction(parentID types.SiacoinOutputID, parentUnlockConditions types.UnlockConditions,
	destinationUnlockHash types.UnlockHash, value types.Currency, minerFee types.Currency,
	timelock types.BlockHeight) types.Transaction {
	tx := types.Transaction{}
	tx.SiacoinInputs = []types.SiacoinInput{{
		ParentID:         parentID,
		UnlockConditions: parentUnlockConditions,
	}}
	tx.SiacoinOutputs = []types.SiacoinOutput{{
		Value:      value,
		UnlockHash: destinationUnlockHash,
	}}
	tx.MinerFees = []types.Currency{minerFee}
	tx.TransactionSignatures = []types.TransactionSignature{{
		ParentID:      crypto.Hash(parentID),
		Timelock:      timelock,
		CoveredFields: types.CoveredFields{WholeTransaction: true},
	}}
	return tx
}

func BuildClaimTransaction(parentID types.SiacoinOutputID, parentUnlockConditions types.UnlockConditions,
	destinationUnlockHash types.UnlockHash, value types.Currency, minerFee types.Currency) types.Transaction {
	tx := types.Transaction{}
	tx.SiacoinInputs = []types.SiacoinInput{{
		ParentID:         parentID,
		UnlockConditions: parentUnlockConditions,
	}}
	tx.SiacoinOutputs = []types.SiacoinOutput{{
		Value:      value,
		UnlockHash: destinationUnlockHash,
	}}
	tx.MinerFees = []types.Currency{minerFee}
	tx.TransactionSignatures = []types.TransactionSignature{{
		ParentID:      crypto.Hash(parentID),
		CoveredFields: types.CoveredFields{WholeTransaction: true},
	}}
	return tx
}

func WholeSigHash(tx types.Transaction, blockHeight types.BlockHeight) []byte {
	sigHash := tx.SigHash(0, blockHeight)
	return sigHash[:]
}

func AddSignature(tx types.Transaction, signature []byte) types.Transaction {
	tx.TransactionSignatures[0].Signature = signature
	return tx
}

func EncodeTransaction(tx types.Transaction) string {
	return base64.StdEncoding.EncodeToString(encoding.Marshal(tx))
}
