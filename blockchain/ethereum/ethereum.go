package ethereum

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"

	"github.com/HyperspaceApp/ed25519"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	contract "github.com/javgh/roadie/contract/hub"
)

const (
	formatEtherPrecision = 6
	smallGasLimit        = 100000
	mediumGasLimit       = 200000
	largeGasLimit        = 1500000
	ganacheEndpoint      = "http://127.0.0.1:8545"
	ganachePrivKey       = "a1d63a5f23ac9b62199e84d87fff196c603b61f6c42bddd0bcca9839d7449ba7"
)

var (
	ErrCasting = errors.New("error casting public key to ECDSA")

	oneEther = big.NewInt(1e18)
)

type (
	JSONRPCBlockchain struct {
		client        ethclient.Client
		privKey       ecdsa.PrivateKey
		walletAddress common.Address
		hub           contract.Hub
	}

	Blockchain interface {
		BurnAntiSpamFee(antiSpamID big.Int, antiSpamFee big.Int) error
		CheckAntiSpamConfirmations(antiSpamID big.Int, antiSpamFee big.Int) (int64, error)
		DepositEther(recipient common.Address, adaptorPubKey ed25519.CurvePoint, ether big.Int, antiSpamID big.Int) error
		CheckDepositConfirmations(recipient common.Address, adaptorPubKey ed25519.CurvePoint, ether big.Int, antiSpamID big.Int) (int64, error)
		ClaimDeposit(adaptorPrivKey ed25519.Adaptor, antiSpamID big.Int) error
		LookupAdaptorPrivKey(adaptorPubKey ed25519.CurvePoint) (bool, *ed25519.Adaptor, error)
		WalletAddress() common.Address
	}
)

func NewGanacheBlockchain() (*JSONRPCBlockchain, error) {
	client, err := ethclient.Dial(ganacheEndpoint)
	if err != nil {
		return nil, err
	}

	privKeyECDSA, err := crypto.HexToECDSA(ganachePrivKey)
	if err != nil {
		return nil, err
	}

	pubKey := privKeyECDSA.Public()
	pubKeyECDSA, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, ErrCasting
	}
	walletAddress := crypto.PubkeyToAddress(*pubKeyECDSA)

	auth := bind.NewKeyedTransactor(privKeyECDSA)
	_, _, hub, err := contract.DeployHub(auth, client)
	if err != nil {
		return nil, err
	}

	c := JSONRPCBlockchain{
		client:        *client,
		privKey:       *privKeyECDSA,
		walletAddress: walletAddress,
		hub:           *hub,
	}
	return &c, nil
}

func NewJSONRPCBlockchain(endpoint string, keystoreFile string, contractAddress *common.Address) (*JSONRPCBlockchain, error) {
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		return nil, err
	}

	json, err := ioutil.ReadFile(keystoreFile)
	if err != nil {
		return nil, err
	}

	key, err := keystore.DecryptKey(json, "")
	if err != nil {
		return nil, err
	}

	var hub *contract.Hub
	if contractAddress != nil {
		hub, err = contract.NewHub(*contractAddress, client)
		if err != nil {
			return nil, err
		}
	} else {
		auth := bind.NewKeyedTransactor(key.PrivateKey)

		//gwei := new(big.Int).Mul(big.NewInt(5), big.NewInt(1e9))
		//auth.Nonce = big.NewInt(0)
		//auth.GasPrice = gwei

		_, _, hub, err = contract.DeployHub(auth, client)
		if err != nil {
			return nil, err
		}
	}

	c := JSONRPCBlockchain{
		client:        *client,
		privKey:       *key.PrivateKey,
		walletAddress: key.Address,
		hub:           *hub,
	}
	return &c, nil
}

func (c *JSONRPCBlockchain) BurnAntiSpamFee(antiSpamID big.Int, antiSpamFee big.Int) error {
	hashedID := hash(antiSpamID)

	auth := bind.NewKeyedTransactor(&c.privKey)
	auth.Value = &antiSpamFee
	auth.GasLimit = smallGasLimit

	_, err := c.hub.BurnAntiSpamFee(auth, hashedID)
	return err
}

func (c *JSONRPCBlockchain) CheckAntiSpamConfirmations(antiSpamID big.Int, antiSpamFee big.Int) (int64, error) {
	confs, err := c.hub.CheckAntiSpamConfirmations(nil, &antiSpamID, &antiSpamFee)
	if err != nil {
		return 0, err
	}

	return confs.Int64(), nil
}
func (c *JSONRPCBlockchain) DepositEther(
	recipient common.Address, adaptorPubKey ed25519.CurvePoint, ether big.Int, antiSpamID big.Int) error {
	hashedID := hash(antiSpamID)

	adaptorPubKeyBytes := switchEndianness(adaptorPubKey[:])
	adaptorPubKeyBytes[0] &= 127 // clear sign bit
	adaptorPubKeyBigInt := new(big.Int).SetBytes(adaptorPubKeyBytes)

	auth := bind.NewKeyedTransactor(&c.privKey)
	auth.Value = &ether
	auth.GasLimit = mediumGasLimit

	_, err := c.hub.DepositEther(auth, recipient, adaptorPubKeyBigInt, hashedID)
	return err
}

func (c *JSONRPCBlockchain) CheckDepositConfirmations(
	recipient common.Address, adaptorPubKey ed25519.CurvePoint, ether big.Int, antiSpamID big.Int) (int64, error) {
	hashedID := hash(antiSpamID)

	adaptorPubKeyBytes := switchEndianness(adaptorPubKey[:])
	adaptorPubKeyBytes[0] &= 127 // clear sign bit
	adaptorPubKeyBigInt := new(big.Int).SetBytes(adaptorPubKeyBytes)

	confs, err := c.hub.CheckDepositConfirmations(nil, recipient, adaptorPubKeyBigInt, &ether, hashedID)
	if err != nil {
		return 0, err
	}

	return confs.Int64(), nil
}

func (c *JSONRPCBlockchain) ClaimDeposit(adaptorPrivKey ed25519.Adaptor, antiSpamID big.Int) error {
	adaptorPrivKeyBigInt := new(big.Int).SetBytes(switchEndianness(adaptorPrivKey[:]))

	auth := bind.NewKeyedTransactor(&c.privKey)
	auth.GasLimit = largeGasLimit

	_, err := c.hub.ClaimDeposit(auth, adaptorPrivKeyBigInt, &antiSpamID)
	return err
}

func (c *JSONRPCBlockchain) LookupAdaptorPrivKey(adaptorPubKey ed25519.CurvePoint) (bool, *ed25519.Adaptor, error) {
	adaptorPubKeyBytes := switchEndianness(adaptorPubKey[:])
	adaptorPubKeyBytes[0] &= 127 // clear sign bit
	adaptorPubKeyBigInt := new(big.Int).SetBytes(adaptorPubKeyBytes)

	adaptorPrivKeyBigInt, err := c.hub.AdaptorPrivKeys(nil, adaptorPubKeyBigInt)
	if err != nil {
		return false, nil, err
	}

	if adaptorPrivKeyBigInt.Cmp(big.NewInt(0)) == 0 {
		return false, nil, nil
	}

	adaptorPrivKey := ed25519.Adaptor(switchEndianness(adaptorPrivKeyBigInt.Bytes()))
	return true, &adaptorPrivKey, nil
}

func (c *JSONRPCBlockchain) WalletAddress() common.Address {
	return c.walletAddress
}

func switchEndianness(in []byte) []byte {
	out := make([]byte, len(in))
	for i := range in {
		out[i] = in[len(in)-1-i]
	}
	return out
}

func hash(id big.Int) [32]byte {
	return sha256.Sum256(math.PaddedBigBytes(&id, 32))
}

func FormatEther(ether *big.Int) string {
	r := big.NewRat(0, 1).SetFrac(ether, oneEther)
	return fmt.Sprintf("%s ETH", r.FloatString(formatEtherPrecision))
}

func ApplyRate(ether *big.Int, rate *big.Float) *big.Float {
	r := new(big.Rat).SetFrac(ether, oneEther)
	f := new(big.Float).Mul(new(big.Float).SetRat(r), rate)
	return f
}
