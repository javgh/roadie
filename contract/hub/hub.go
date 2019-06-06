// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package hub

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// HubABI is the input ABI used to generate the binding from.
const HubABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"deposits\",\"outputs\":[{\"name\":\"recipient\",\"type\":\"address\"},{\"name\":\"adaptorPubKey\",\"type\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\"},{\"name\":\"adaptorPubKey\",\"type\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"hashedAntiSpamID\",\"type\":\"bytes32\"}],\"name\":\"checkDepositConfirmations\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"antiSpamFees\",\"outputs\":[{\"name\":\"fee\",\"type\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"},{\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"checkAntiSpamConfirmations\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"hashedID\",\"type\":\"bytes32\"}],\"name\":\"burnAntiSpamFee\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"hash\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\"},{\"name\":\"adaptorPubKey\",\"type\":\"uint256\"},{\"name\":\"hashedAntiSpamID\",\"type\":\"bytes32\"}],\"name\":\"depositEther\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"s\",\"type\":\"uint256\"}],\"name\":\"scalarmult\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"adaptorPrivKeys\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"adaptorPrivKey\",\"type\":\"uint256\"},{\"name\":\"antiSpamID\",\"type\":\"uint256\"}],\"name\":\"claimDeposit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// HubBin is the compiled bytecode used for deploying new contracts.
const HubBin = `608060405234801561001057600080fd5b50610aff806100206000396000f3fe6080604052600436106100915760003560e01c8063b189fd4c11610059578063b189fd4c146101d9578063b90d104d14610203578063e49cf91114610235578063e74db5a91461025f578063ea32a89e1461028957610091565b80633d4dff7b1461009657806357888e92146100f05780635a161ba51461014757806366db09c61461018a578063ab80cdc2146101ba575b600080fd5b3480156100a257600080fd5b506100c0600480360360208110156100b957600080fd5b50356102b9565b604080516001600160a01b0390951685526020850193909352838301919091526060830152519081900360800190f35b3480156100fc57600080fd5b506101356004803603608081101561011357600080fd5b506001600160a01b0381351690602081013590604081013590606001356102ed565b60408051918252519081900360200190f35b34801561015357600080fd5b506101716004803603602081101561016a57600080fd5b5035610372565b6040805192835260208301919091528051918290030190f35b34801561019657600080fd5b50610135600480360360408110156101ad57600080fd5b508035906020013561038b565b6101d7600480360360208110156101d057600080fd5b50356103da565b005b3480156101e557600080fd5b50610135600480360360208110156101fc57600080fd5b5035610427565b6101d76004803603606081101561021957600080fd5b506001600160a01b0381351690602081013590604001356104d6565b34801561024157600080fd5b506101716004803603602081101561025857600080fd5b5035610533565b34801561026b57600080fd5b506101356004803603602081101561028257600080fd5b5035610626565b34801561029557600080fd5b506101d7600480360360408110156102ac57600080fd5b5080359060200135610638565b600160208190526000918252604090912080549181015460028201546003909201546001600160a01b039093169290919084565b6000818152600160205260408120546001600160a01b0386811691161415806103285750600082815260016020819052604090912001548414155b80610343575060008281526001602052604090206002015483115b156103505750600061036a565b506000818152600160208190526040909120600301544303015b949350505050565b6000602081905290815260409020805460019091015482565b60008061039784610427565b6000818152602081905260409020549091508311156103ba5760009150506103d4565b600090815260208190526040902060019081015443030190505b92915050565b600081815260208190526040808220805434908101825543600190920191909155905181156108fc02919083818181858288f19350505050158015610423573d6000803e3d6000fd5b5050565b6000600282604051602001808281526020019150506040516020818303038152906040526040518082805190602001908083835b6020831061047a5780518252601f19909201916020918201910161045b565b51815160209384036101000a60001901801990921691161790526040519190930194509192505080830381855afa1580156104b9573d6000803e3d6000fd5b5050506040513d60208110156104ce57600080fd5b505192915050565b600081815260016020526040902060030154156104f257600080fd5b600090815260016020819052604090912080546001600160a01b0319166001600160a01b039490941693909317835582015534600282015543600390910155565b60008061053e610a64565b610546610a64565b7f216936d3cd6e53fec0a4e231fdd6dc5c692cc7609525a7b2c9562d608f25d51a82527f66666666666666666666666666666666666666666666666666666666666666586020808401919091526001604080850182905260008452918301819052908201525b84156105e25784600116600114156105cb576105c8818361072b565b90505b600185901c94506105db826108c1565b91506105ac565b60006105f18260400151610a04565b90506013600160ff1b03825182900982526013600160ff1b038183602001510960208301819052915194509092505050915091565b60026020526000908152604090205481565b600061064382610427565b6000818152600160205260409020549091506001600160a01b0316331461066957600080fd5b8261067357600080fd5b600061067e84610533565b915050806001600084815260200190815260200160002060010154146106a357600080fd5b60008181526002602081815260408084208890558584526001808352818520938401805485546001600160a01b031916865585830187905590869055600390940185905591849052808420848155909101839055519091339183156108fc0291849190818181858888f19350505050158015610723573d6000803e3d6000fd5b505050505050565b610733610a64565b61073b610a85565b6013600160ff1b03836040015185604001510981526013600160ff1b038151800960208201526013600160ff1b03835185510960408201526013600160ff1b03836020015185602001510960608201526013600160ff1b038082606001518360400151097f52036cee2b6ffe738cc740797779e89800700a4d4141d8ab75eb4dca135978a30960808201526013600160ff1b0381608001516013600160ff1b030382602001510860a08201526013600160ff1b03816080015182602001510860c08201526013600160ff1b038082606001516013600160ff1b03036013600160ff1b038061082557fe5b84604001516013600160ff1b03036013600160ff1b038061084257fe5b6013600160ff1b0360208a01518a51086013600160ff1b0360208c01518c51080908086013600160ff1b0360a08401518451090982526013600160ff1b038082604001518360600151086013600160ff1b0360c08401518451090960208301526013600160ff1b038160c001518260a001510960408301525092915050565b6108c9610a64565b6108d1610a85565b6013600160ff1b03602084015184510881526013600160ff1b038151800960208201526013600160ff1b038351800960408201526013600160ff1b03602084015180096060820181905260408201516013600160ff1b03908103608084018190529091900860a08201526013600160ff1b036040840151800960e08201526013600160ff1b03808260e001516002096013600160ff1b03038260a001510860c08201526013600160ff1b0360c08201516013600160ff1b0383606001516013600160ff1b03036013600160ff1b03806109a657fe5b85604001516013600160ff1b0303866020015108080982526013600160ff1b038082606001516013600160ff1b03038360800151088260a001510960208301526013600160ff1b038160c001518260a0015109604083015250919050565b60008060026013600160ff1b0303905060006013600160ff1b03905060405160208152602080820152602060408201528460608201528260808201528160a082015260208160c0836005600019fa610a5b57600080fd5b51949350505050565b60405180606001604052806000815260200160008152602001600081525090565b6040518061010001604052806000815260200160008152602001600081526020016000815260200160008152602001600081526020016000815260200160008152509056fea265627a7a723058201e240eca4113f4b04e3e6f175bc28987b8cc3c456081d6e92878b4601440dd3c64736f6c63430005090032`

// DeployHub deploys a new Ethereum contract, binding an instance of Hub to it.
func DeployHub(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Hub, error) {
	parsed, err := abi.JSON(strings.NewReader(HubABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(HubBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Hub{HubCaller: HubCaller{contract: contract}, HubTransactor: HubTransactor{contract: contract}, HubFilterer: HubFilterer{contract: contract}}, nil
}

// Hub is an auto generated Go binding around an Ethereum contract.
type Hub struct {
	HubCaller     // Read-only binding to the contract
	HubTransactor // Write-only binding to the contract
	HubFilterer   // Log filterer for contract events
}

// HubCaller is an auto generated read-only Go binding around an Ethereum contract.
type HubCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HubTransactor is an auto generated write-only Go binding around an Ethereum contract.
type HubTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HubFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type HubFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HubSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type HubSession struct {
	Contract     *Hub              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// HubCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type HubCallerSession struct {
	Contract *HubCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// HubTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type HubTransactorSession struct {
	Contract     *HubTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// HubRaw is an auto generated low-level Go binding around an Ethereum contract.
type HubRaw struct {
	Contract *Hub // Generic contract binding to access the raw methods on
}

// HubCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type HubCallerRaw struct {
	Contract *HubCaller // Generic read-only contract binding to access the raw methods on
}

// HubTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type HubTransactorRaw struct {
	Contract *HubTransactor // Generic write-only contract binding to access the raw methods on
}

// NewHub creates a new instance of Hub, bound to a specific deployed contract.
func NewHub(address common.Address, backend bind.ContractBackend) (*Hub, error) {
	contract, err := bindHub(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Hub{HubCaller: HubCaller{contract: contract}, HubTransactor: HubTransactor{contract: contract}, HubFilterer: HubFilterer{contract: contract}}, nil
}

// NewHubCaller creates a new read-only instance of Hub, bound to a specific deployed contract.
func NewHubCaller(address common.Address, caller bind.ContractCaller) (*HubCaller, error) {
	contract, err := bindHub(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &HubCaller{contract: contract}, nil
}

// NewHubTransactor creates a new write-only instance of Hub, bound to a specific deployed contract.
func NewHubTransactor(address common.Address, transactor bind.ContractTransactor) (*HubTransactor, error) {
	contract, err := bindHub(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &HubTransactor{contract: contract}, nil
}

// NewHubFilterer creates a new log filterer instance of Hub, bound to a specific deployed contract.
func NewHubFilterer(address common.Address, filterer bind.ContractFilterer) (*HubFilterer, error) {
	contract, err := bindHub(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &HubFilterer{contract: contract}, nil
}

// bindHub binds a generic wrapper to an already deployed contract.
func bindHub(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(HubABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Hub *HubRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Hub.Contract.HubCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Hub *HubRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Hub.Contract.HubTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Hub *HubRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Hub.Contract.HubTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Hub *HubCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Hub.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Hub *HubTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Hub.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Hub *HubTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Hub.Contract.contract.Transact(opts, method, params...)
}

// AdaptorPrivKeys is a free data retrieval call binding the contract method 0xe74db5a9.
//
// Solidity: function adaptorPrivKeys(uint256 ) constant returns(uint256)
func (_Hub *HubCaller) AdaptorPrivKeys(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Hub.contract.Call(opts, out, "adaptorPrivKeys", arg0)
	return *ret0, err
}

// AdaptorPrivKeys is a free data retrieval call binding the contract method 0xe74db5a9.
//
// Solidity: function adaptorPrivKeys(uint256 ) constant returns(uint256)
func (_Hub *HubSession) AdaptorPrivKeys(arg0 *big.Int) (*big.Int, error) {
	return _Hub.Contract.AdaptorPrivKeys(&_Hub.CallOpts, arg0)
}

// AdaptorPrivKeys is a free data retrieval call binding the contract method 0xe74db5a9.
//
// Solidity: function adaptorPrivKeys(uint256 ) constant returns(uint256)
func (_Hub *HubCallerSession) AdaptorPrivKeys(arg0 *big.Int) (*big.Int, error) {
	return _Hub.Contract.AdaptorPrivKeys(&_Hub.CallOpts, arg0)
}

// AntiSpamFees is a free data retrieval call binding the contract method 0x5a161ba5.
//
// Solidity: function antiSpamFees(bytes32 ) constant returns(uint256 fee, uint256 blockNumber)
func (_Hub *HubCaller) AntiSpamFees(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Fee         *big.Int
	BlockNumber *big.Int
}, error) {
	ret := new(struct {
		Fee         *big.Int
		BlockNumber *big.Int
	})
	out := ret
	err := _Hub.contract.Call(opts, out, "antiSpamFees", arg0)
	return *ret, err
}

// AntiSpamFees is a free data retrieval call binding the contract method 0x5a161ba5.
//
// Solidity: function antiSpamFees(bytes32 ) constant returns(uint256 fee, uint256 blockNumber)
func (_Hub *HubSession) AntiSpamFees(arg0 [32]byte) (struct {
	Fee         *big.Int
	BlockNumber *big.Int
}, error) {
	return _Hub.Contract.AntiSpamFees(&_Hub.CallOpts, arg0)
}

// AntiSpamFees is a free data retrieval call binding the contract method 0x5a161ba5.
//
// Solidity: function antiSpamFees(bytes32 ) constant returns(uint256 fee, uint256 blockNumber)
func (_Hub *HubCallerSession) AntiSpamFees(arg0 [32]byte) (struct {
	Fee         *big.Int
	BlockNumber *big.Int
}, error) {
	return _Hub.Contract.AntiSpamFees(&_Hub.CallOpts, arg0)
}

// CheckAntiSpamConfirmations is a free data retrieval call binding the contract method 0x66db09c6.
//
// Solidity: function checkAntiSpamConfirmations(uint256 id, uint256 fee) constant returns(uint256)
func (_Hub *HubCaller) CheckAntiSpamConfirmations(opts *bind.CallOpts, id *big.Int, fee *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Hub.contract.Call(opts, out, "checkAntiSpamConfirmations", id, fee)
	return *ret0, err
}

// CheckAntiSpamConfirmations is a free data retrieval call binding the contract method 0x66db09c6.
//
// Solidity: function checkAntiSpamConfirmations(uint256 id, uint256 fee) constant returns(uint256)
func (_Hub *HubSession) CheckAntiSpamConfirmations(id *big.Int, fee *big.Int) (*big.Int, error) {
	return _Hub.Contract.CheckAntiSpamConfirmations(&_Hub.CallOpts, id, fee)
}

// CheckAntiSpamConfirmations is a free data retrieval call binding the contract method 0x66db09c6.
//
// Solidity: function checkAntiSpamConfirmations(uint256 id, uint256 fee) constant returns(uint256)
func (_Hub *HubCallerSession) CheckAntiSpamConfirmations(id *big.Int, fee *big.Int) (*big.Int, error) {
	return _Hub.Contract.CheckAntiSpamConfirmations(&_Hub.CallOpts, id, fee)
}

// CheckDepositConfirmations is a free data retrieval call binding the contract method 0x57888e92.
//
// Solidity: function checkDepositConfirmations(address recipient, uint256 adaptorPubKey, uint256 value, bytes32 hashedAntiSpamID) constant returns(uint256)
func (_Hub *HubCaller) CheckDepositConfirmations(opts *bind.CallOpts, recipient common.Address, adaptorPubKey *big.Int, value *big.Int, hashedAntiSpamID [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Hub.contract.Call(opts, out, "checkDepositConfirmations", recipient, adaptorPubKey, value, hashedAntiSpamID)
	return *ret0, err
}

// CheckDepositConfirmations is a free data retrieval call binding the contract method 0x57888e92.
//
// Solidity: function checkDepositConfirmations(address recipient, uint256 adaptorPubKey, uint256 value, bytes32 hashedAntiSpamID) constant returns(uint256)
func (_Hub *HubSession) CheckDepositConfirmations(recipient common.Address, adaptorPubKey *big.Int, value *big.Int, hashedAntiSpamID [32]byte) (*big.Int, error) {
	return _Hub.Contract.CheckDepositConfirmations(&_Hub.CallOpts, recipient, adaptorPubKey, value, hashedAntiSpamID)
}

// CheckDepositConfirmations is a free data retrieval call binding the contract method 0x57888e92.
//
// Solidity: function checkDepositConfirmations(address recipient, uint256 adaptorPubKey, uint256 value, bytes32 hashedAntiSpamID) constant returns(uint256)
func (_Hub *HubCallerSession) CheckDepositConfirmations(recipient common.Address, adaptorPubKey *big.Int, value *big.Int, hashedAntiSpamID [32]byte) (*big.Int, error) {
	return _Hub.Contract.CheckDepositConfirmations(&_Hub.CallOpts, recipient, adaptorPubKey, value, hashedAntiSpamID)
}

// Deposits is a free data retrieval call binding the contract method 0x3d4dff7b.
//
// Solidity: function deposits(bytes32 ) constant returns(address recipient, uint256 adaptorPubKey, uint256 value, uint256 blockNumber)
func (_Hub *HubCaller) Deposits(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Recipient     common.Address
	AdaptorPubKey *big.Int
	Value         *big.Int
	BlockNumber   *big.Int
}, error) {
	ret := new(struct {
		Recipient     common.Address
		AdaptorPubKey *big.Int
		Value         *big.Int
		BlockNumber   *big.Int
	})
	out := ret
	err := _Hub.contract.Call(opts, out, "deposits", arg0)
	return *ret, err
}

// Deposits is a free data retrieval call binding the contract method 0x3d4dff7b.
//
// Solidity: function deposits(bytes32 ) constant returns(address recipient, uint256 adaptorPubKey, uint256 value, uint256 blockNumber)
func (_Hub *HubSession) Deposits(arg0 [32]byte) (struct {
	Recipient     common.Address
	AdaptorPubKey *big.Int
	Value         *big.Int
	BlockNumber   *big.Int
}, error) {
	return _Hub.Contract.Deposits(&_Hub.CallOpts, arg0)
}

// Deposits is a free data retrieval call binding the contract method 0x3d4dff7b.
//
// Solidity: function deposits(bytes32 ) constant returns(address recipient, uint256 adaptorPubKey, uint256 value, uint256 blockNumber)
func (_Hub *HubCallerSession) Deposits(arg0 [32]byte) (struct {
	Recipient     common.Address
	AdaptorPubKey *big.Int
	Value         *big.Int
	BlockNumber   *big.Int
}, error) {
	return _Hub.Contract.Deposits(&_Hub.CallOpts, arg0)
}

// Hash is a free data retrieval call binding the contract method 0xb189fd4c.
//
// Solidity: function hash(uint256 id) constant returns(bytes32)
func (_Hub *HubCaller) Hash(opts *bind.CallOpts, id *big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Hub.contract.Call(opts, out, "hash", id)
	return *ret0, err
}

// Hash is a free data retrieval call binding the contract method 0xb189fd4c.
//
// Solidity: function hash(uint256 id) constant returns(bytes32)
func (_Hub *HubSession) Hash(id *big.Int) ([32]byte, error) {
	return _Hub.Contract.Hash(&_Hub.CallOpts, id)
}

// Hash is a free data retrieval call binding the contract method 0xb189fd4c.
//
// Solidity: function hash(uint256 id) constant returns(bytes32)
func (_Hub *HubCallerSession) Hash(id *big.Int) ([32]byte, error) {
	return _Hub.Contract.Hash(&_Hub.CallOpts, id)
}

// Scalarmult is a free data retrieval call binding the contract method 0xe49cf911.
//
// Solidity: function scalarmult(uint256 s) constant returns(uint256, uint256)
func (_Hub *HubCaller) Scalarmult(opts *bind.CallOpts, s *big.Int) (*big.Int, *big.Int, error) {
	var (
		ret0 = new(*big.Int)
		ret1 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := _Hub.contract.Call(opts, out, "scalarmult", s)
	return *ret0, *ret1, err
}

// Scalarmult is a free data retrieval call binding the contract method 0xe49cf911.
//
// Solidity: function scalarmult(uint256 s) constant returns(uint256, uint256)
func (_Hub *HubSession) Scalarmult(s *big.Int) (*big.Int, *big.Int, error) {
	return _Hub.Contract.Scalarmult(&_Hub.CallOpts, s)
}

// Scalarmult is a free data retrieval call binding the contract method 0xe49cf911.
//
// Solidity: function scalarmult(uint256 s) constant returns(uint256, uint256)
func (_Hub *HubCallerSession) Scalarmult(s *big.Int) (*big.Int, *big.Int, error) {
	return _Hub.Contract.Scalarmult(&_Hub.CallOpts, s)
}

// BurnAntiSpamFee is a paid mutator transaction binding the contract method 0xab80cdc2.
//
// Solidity: function burnAntiSpamFee(bytes32 hashedID) returns()
func (_Hub *HubTransactor) BurnAntiSpamFee(opts *bind.TransactOpts, hashedID [32]byte) (*types.Transaction, error) {
	return _Hub.contract.Transact(opts, "burnAntiSpamFee", hashedID)
}

// BurnAntiSpamFee is a paid mutator transaction binding the contract method 0xab80cdc2.
//
// Solidity: function burnAntiSpamFee(bytes32 hashedID) returns()
func (_Hub *HubSession) BurnAntiSpamFee(hashedID [32]byte) (*types.Transaction, error) {
	return _Hub.Contract.BurnAntiSpamFee(&_Hub.TransactOpts, hashedID)
}

// BurnAntiSpamFee is a paid mutator transaction binding the contract method 0xab80cdc2.
//
// Solidity: function burnAntiSpamFee(bytes32 hashedID) returns()
func (_Hub *HubTransactorSession) BurnAntiSpamFee(hashedID [32]byte) (*types.Transaction, error) {
	return _Hub.Contract.BurnAntiSpamFee(&_Hub.TransactOpts, hashedID)
}

// ClaimDeposit is a paid mutator transaction binding the contract method 0xea32a89e.
//
// Solidity: function claimDeposit(uint256 adaptorPrivKey, uint256 antiSpamID) returns()
func (_Hub *HubTransactor) ClaimDeposit(opts *bind.TransactOpts, adaptorPrivKey *big.Int, antiSpamID *big.Int) (*types.Transaction, error) {
	return _Hub.contract.Transact(opts, "claimDeposit", adaptorPrivKey, antiSpamID)
}

// ClaimDeposit is a paid mutator transaction binding the contract method 0xea32a89e.
//
// Solidity: function claimDeposit(uint256 adaptorPrivKey, uint256 antiSpamID) returns()
func (_Hub *HubSession) ClaimDeposit(adaptorPrivKey *big.Int, antiSpamID *big.Int) (*types.Transaction, error) {
	return _Hub.Contract.ClaimDeposit(&_Hub.TransactOpts, adaptorPrivKey, antiSpamID)
}

// ClaimDeposit is a paid mutator transaction binding the contract method 0xea32a89e.
//
// Solidity: function claimDeposit(uint256 adaptorPrivKey, uint256 antiSpamID) returns()
func (_Hub *HubTransactorSession) ClaimDeposit(adaptorPrivKey *big.Int, antiSpamID *big.Int) (*types.Transaction, error) {
	return _Hub.Contract.ClaimDeposit(&_Hub.TransactOpts, adaptorPrivKey, antiSpamID)
}

// DepositEther is a paid mutator transaction binding the contract method 0xb90d104d.
//
// Solidity: function depositEther(address recipient, uint256 adaptorPubKey, bytes32 hashedAntiSpamID) returns()
func (_Hub *HubTransactor) DepositEther(opts *bind.TransactOpts, recipient common.Address, adaptorPubKey *big.Int, hashedAntiSpamID [32]byte) (*types.Transaction, error) {
	return _Hub.contract.Transact(opts, "depositEther", recipient, adaptorPubKey, hashedAntiSpamID)
}

// DepositEther is a paid mutator transaction binding the contract method 0xb90d104d.
//
// Solidity: function depositEther(address recipient, uint256 adaptorPubKey, bytes32 hashedAntiSpamID) returns()
func (_Hub *HubSession) DepositEther(recipient common.Address, adaptorPubKey *big.Int, hashedAntiSpamID [32]byte) (*types.Transaction, error) {
	return _Hub.Contract.DepositEther(&_Hub.TransactOpts, recipient, adaptorPubKey, hashedAntiSpamID)
}

// DepositEther is a paid mutator transaction binding the contract method 0xb90d104d.
//
// Solidity: function depositEther(address recipient, uint256 adaptorPubKey, bytes32 hashedAntiSpamID) returns()
func (_Hub *HubTransactorSession) DepositEther(recipient common.Address, adaptorPubKey *big.Int, hashedAntiSpamID [32]byte) (*types.Transaction, error) {
	return _Hub.Contract.DepositEther(&_Hub.TransactOpts, recipient, adaptorPubKey, hashedAntiSpamID)
}
