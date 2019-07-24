package ethereum

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/HyperspaceApp/ed25519"
	"github.com/blang/semver"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pborman/uuid"

	contract "github.com/javgh/roadie/contract/hub"
	"github.com/javgh/roadie/contract/retryinghub"
)

const (
	formatEtherPrecision     = 6
	formatGweiPrecision      = 1
	smallGasLimit            = 100000
	mediumGasLimit           = 200000
	largeGasLimit            = 1500000
	txCheckInterval          = 10 * time.Second
	ganacheEndpoint          = "http://127.0.0.1:8545"
	ganachePrivKey           = "a1d63a5f23ac9b62199e84d87fff196c603b61f6c42bddd0bcca9839d7449ba7"
	ganacheBoostInterval     = 5 * time.Second
	ganacheTxCheckInterval   = time.Second
	simulatedPrivKey         = "a1d63a5f23ac9b62199e84d87fff196c603b61f6c42bddd0bcca9839d7449ba7"
	simulatedBlockInterval   = 500 * time.Millisecond
	simulatedTxCheckInterval = time.Second
	requiredMajorVersion     = 0
)

var (
	ErrStillSyncing        = errors.New("Ethereum node is still syncing")
	ErrIncompatibleVersion = errors.New("smart contract has an incompatible version - please upgrade")
	ErrDeprecated          = errors.New("smart contract is marked as deprecated - please check for updates")
	ErrUnexpectedDirectory = errors.New("keystore location appears to be a directory")

	gwei               = big.NewInt(1e9)
	oneEther           = big.NewInt(1e18)
	ganacheMaxGasPrice = new(big.Int).Mul(big.NewInt(100), gwei)
	simulatedBalance   = new(big.Int).Mul(big.NewInt(100), oneEther)
	simulatedGasLimit  = uint64(10000000)
)

type (
	GethBlockchain struct {
		backend       bind.ContractBackend
		walletAddress common.Address
		retryingHub   retryinghub.RetryingHub
	}

	ServerDetails struct {
		Target string
		Cert   []byte
	}

	Blockchain interface {
		CheckSmartContract() error
		BurnAntiSpamFee(antiSpamID big.Int, antiSpamFee big.Int) error
		CheckAntiSpamConfirmations(antiSpamID big.Int, antiSpamFee big.Int) (int64, error)
		DepositEther(recipient common.Address, adaptorPubKey ed25519.CurvePoint, ether big.Int, antiSpamID big.Int) error
		CheckDepositConfirmations(recipient common.Address, adaptorPubKey ed25519.CurvePoint, ether big.Int, antiSpamID big.Int) (int64, error)
		ClaimDeposit(adaptorPrivKey ed25519.Adaptor, antiSpamID big.Int) error
		LookupAdaptorPrivKey(adaptorPubKey ed25519.CurvePoint) (bool, *ed25519.Adaptor, error)
		ReclaimDeposit(antiSpamID big.Int) error
		RegisterServer(target string, cert []byte) error
		FetchServers(maxAge big.Int) ([]ServerDetails, error)
		WalletAddress() common.Address
		SuggestGasPrice() (*big.Int, error)
	}
)

func NewGanacheBlockchain(contractAddress *common.Address) (*GethBlockchain, error) {
	client, err := ethclient.Dial(ganacheEndpoint)
	if err != nil {
		return nil, err
	}

	privKeyECDSA, err := crypto.HexToECDSA(ganachePrivKey)
	if err != nil {
		return nil, err
	}
	walletAddress := crypto.PubkeyToAddress(privKeyECDSA.PublicKey)

	var hub *contract.Hub
	if contractAddress != nil {
		hub, err = contract.NewHub(*contractAddress, client)
		if err != nil {
			return nil, err
		}
	} else {
		auth := bind.NewKeyedTransactor(privKeyECDSA)
		_, _, hub, err = contract.DeployHub(auth, client)
		if err != nil {
			return nil, err
		}
	}

	time.Sleep(1200 * time.Millisecond) // wait for contract to deploy

	retryingHub := retryinghub.New(
		*ganacheMaxGasPrice, ganacheBoostInterval, ganacheTxCheckInterval, client, *privKeyECDSA, walletAddress, hub)

	c := GethBlockchain{
		backend:       client,
		walletAddress: walletAddress,
		retryingHub:   retryingHub,
	}
	return &c, nil
}

func NewSimulatedBlockchain() (*GethBlockchain, error) {
	privKeyECDSA, err := crypto.HexToECDSA(simulatedPrivKey)
	if err != nil {
		return nil, err
	}
	walletAddress := crypto.PubkeyToAddress(privKeyECDSA.PublicKey)

	backend := backends.NewSimulatedBackend(
		core.GenesisAlloc{walletAddress: {Balance: simulatedBalance}}, simulatedGasLimit)
	go func() {
		for {
			time.Sleep(simulatedBlockInterval)
			backend.Commit()
		}
	}()

	auth := bind.NewKeyedTransactor(privKeyECDSA)
	_, _, hub, err := contract.DeployHub(auth, backend)
	if err != nil {
		return nil, err
	}
	backend.Commit()

	retryingHub := retryinghub.New(
		*ganacheMaxGasPrice, ganacheBoostInterval, ganacheTxCheckInterval, backend, *privKeyECDSA, walletAddress, hub)

	c := GethBlockchain{
		backend:       backend,
		walletAddress: walletAddress,
		retryingHub:   retryingHub,
	}
	return &c, nil
}

func NewLocalNodeBlockchain(endpoint string, keystoreFile string, contractAddress *common.Address,
	maxGasPrice big.Int, boostInterval time.Duration) (*GethBlockchain, error) {
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		return nil, err
	}

	syncProgress, err := client.SyncProgress(context.Background())
	if err != nil {
		return nil, err
	}
	if syncProgress != nil {
		return nil, ErrStillSyncing
	}

	json, err := ioutil.ReadFile(keystoreFile)
	if err != nil {
		return nil, err
	}

	key, err := keystore.DecryptKey(json, "")
	if err != nil {
		return nil, err
	}
	walletAddress := key.Address

	var hub *contract.Hub
	if contractAddress != nil {
		hub, err = contract.NewHub(*contractAddress, client)
		if err != nil {
			return nil, err
		}
	} else {
		auth := bind.NewKeyedTransactor(key.PrivateKey)
		_, _, hub, err = contract.DeployHub(auth, client)
		if err != nil {
			return nil, err
		}
	}

	retryingHub := retryinghub.New(
		maxGasPrice, boostInterval, txCheckInterval, client, *key.PrivateKey, walletAddress, hub)

	c := GethBlockchain{
		backend:       client,
		walletAddress: walletAddress,
		retryingHub:   retryingHub,
	}
	return &c, nil
}

func (c *GethBlockchain) CheckSmartContract() error {
	version := c.retryingHub.Version()
	semVersion, err := semver.Make(version)
	if err != nil {
		return err
	}
	if semVersion.Major != requiredMajorVersion {
		return ErrIncompatibleVersion
	}

	deprecated := c.retryingHub.Deprecated()
	if deprecated {
		return ErrDeprecated
	}

	return nil
}

func (c *GethBlockchain) BurnAntiSpamFee(antiSpamID big.Int, antiSpamFee big.Int) error {
	hashedID := hash(antiSpamID)
	c.retryingHub.BurnAntiSpamFee(hashedID, &antiSpamFee, smallGasLimit)
	return nil
}

func (c *GethBlockchain) CheckAntiSpamConfirmations(antiSpamID big.Int, antiSpamFee big.Int) (int64, error) {
	confs := c.retryingHub.CheckAntiSpamConfirmations(&antiSpamID, &antiSpamFee)
	return confs.Int64(), nil
}
func (c *GethBlockchain) DepositEther(
	recipient common.Address, adaptorPubKey ed25519.CurvePoint, ether big.Int, antiSpamID big.Int) error {
	hashedID := hash(antiSpamID)

	adaptorPubKeyBytes := switchEndianness(adaptorPubKey[:])
	adaptorPubKeyBytes[0] &= 127 // clear sign bit
	adaptorPubKeyBigInt := new(big.Int).SetBytes(adaptorPubKeyBytes)

	c.retryingHub.DepositEther(recipient, adaptorPubKeyBigInt, hashedID, &ether, mediumGasLimit)
	return nil
}

func (c *GethBlockchain) CheckDepositConfirmations(
	recipient common.Address, adaptorPubKey ed25519.CurvePoint, ether big.Int, antiSpamID big.Int) (int64, error) {
	hashedID := hash(antiSpamID)

	adaptorPubKeyBytes := switchEndianness(adaptorPubKey[:])
	adaptorPubKeyBytes[0] &= 127 // clear sign bit
	adaptorPubKeyBigInt := new(big.Int).SetBytes(adaptorPubKeyBytes)

	confs := c.retryingHub.CheckDepositConfirmations(recipient, adaptorPubKeyBigInt, &ether, hashedID)
	return confs.Int64(), nil
}

func (c *GethBlockchain) ClaimDeposit(adaptorPrivKey ed25519.Adaptor, antiSpamID big.Int) error {
	adaptorPrivKeyBigInt := new(big.Int).SetBytes(switchEndianness(adaptorPrivKey[:]))
	c.retryingHub.ClaimDeposit(adaptorPrivKeyBigInt, &antiSpamID, big.NewInt(0), largeGasLimit)
	return nil
}

func (c *GethBlockchain) LookupAdaptorPrivKey(adaptorPubKey ed25519.CurvePoint) (bool, *ed25519.Adaptor, error) {
	adaptorPubKeyBytes := switchEndianness(adaptorPubKey[:])
	adaptorPubKeyBytes[0] &= 127 // clear sign bit
	adaptorPubKeyBigInt := new(big.Int).SetBytes(adaptorPubKeyBytes)

	adaptorPrivKeyBigInt := c.retryingHub.AdaptorPrivKeys(adaptorPubKeyBigInt)
	if adaptorPrivKeyBigInt.Cmp(big.NewInt(0)) == 0 {
		return false, nil, nil
	}

	adaptorPrivKey := ed25519.Adaptor(switchEndianness(adaptorPrivKeyBigInt.Bytes()))
	return true, &adaptorPrivKey, nil
}

func (c *GethBlockchain) ReclaimDeposit(antiSpamID big.Int) error {
	hashedID := hash(antiSpamID)
	c.retryingHub.ReclaimDeposit(hashedID, big.NewInt(0), mediumGasLimit)
	return nil
}

func (c *GethBlockchain) RegisterServer(target string, cert []byte) error {
	c.retryingHub.RegisterServer(target, cert, big.NewInt(0), largeGasLimit)
	return nil
}

func (c *GethBlockchain) FetchServers(maxAge big.Int) ([]ServerDetails, error) {
	offset := big.NewInt(0)
	var serverDetails []ServerDetails

	for {
		moreDetails := c.retryingHub.FetchServer(&maxAge, offset)
		if !moreDetails.OK {
			break
		}

		serverDetails = append(serverDetails,
			ServerDetails{Target: moreDetails.Target, Cert: moreDetails.Cert})
		offset.Add(offset, big.NewInt(1))
	}

	return serverDetails, nil
}

func (c *GethBlockchain) WalletAddress() common.Address {
	return c.walletAddress
}

func (c *GethBlockchain) SuggestGasPrice() (*big.Int, error) {
	return c.backend.SuggestGasPrice(context.Background())
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
	r := new(big.Rat).SetFrac(ether, oneEther)
	return fmt.Sprintf("%s ETH", r.FloatString(formatEtherPrecision))
}

func FormatGwei(ether *big.Int) string {
	r := new(big.Rat).SetFrac(ether, gwei)
	return fmt.Sprintf("%s Gwei", r.FloatString(formatGweiPrecision))
}

func ApplyRate(ether *big.Int, rate *big.Rat) *big.Rat {
	etherRat := new(big.Rat).SetFrac(ether, oneEther)
	result := new(big.Rat).Mul(etherRat, rate)
	return result
}

// EnsureKeystoreExists tries to determine whether we already have a keystore.
// Otherwise it will create a fresh key. The key will be 'encrypted' with an
// empty password. This provides no protection, but will make the keystore
// compatible with other Ethereum wallets.
func EnsureKeystoreExists(path string) error {
	info, err := os.Stat(path)
	if err != nil && os.IsExist(err) {
		return err
	}

	if info != nil && info.IsDir() {
		return ErrUnexpectedDirectory
	}

	if info != nil {
		fmt.Printf("Using Ethereum keystore %s\n", path)
		return nil
	}

	fmt.Printf("Creating new Ethereum keystore %s\n", path)
	err = os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil {
		return err
	}

	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		return err
	}

	key := &keystore.Key{
		Id:         uuid.NewRandom(),
		Address:    crypto.PubkeyToAddress(privateKeyECDSA.PublicKey),
		PrivateKey: privateKeyECDSA,
	}
	json, err := keystore.EncryptKey(key, "", keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, json, 0600)
}
