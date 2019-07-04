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
const HubABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"deposits\",\"outputs\":[{\"name\":\"sender\",\"type\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\"},{\"name\":\"adaptorPubKey\",\"type\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\"},{\"name\":\"adaptorPubKey\",\"type\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"hashedAntiSpamID\",\"type\":\"bytes32\"}],\"name\":\"checkDepositConfirmations\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"antiSpamFees\",\"outputs\":[{\"name\":\"fee\",\"type\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"},{\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"checkAntiSpamConfirmations\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"hashedID\",\"type\":\"bytes32\"}],\"name\":\"burnAntiSpamFee\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"hash\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\"},{\"name\":\"adaptorPubKey\",\"type\":\"uint256\"},{\"name\":\"hashedAntiSpamID\",\"type\":\"bytes32\"}],\"name\":\"depositEther\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"s\",\"type\":\"uint256\"}],\"name\":\"scalarmult\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"adaptorPrivKeys\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"adaptorPrivKey\",\"type\":\"uint256\"},{\"name\":\"antiSpamID\",\"type\":\"uint256\"}],\"name\":\"claimDeposit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"hashedAntiSpamID\",\"type\":\"bytes32\"}],\"name\":\"reclaimDeposit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// HubBin is the compiled bytecode used for deploying new contracts.
const HubBin = `608060405234801561001057600080fd5b50610c87806100206000396000f3fe60806040526004361061009c5760003560e01c8063b189fd4c11610064578063b189fd4c146101f3578063b90d104d1461021d578063e49cf9111461024f578063e74db5a914610279578063ea32a89e146102a3578063fa79c259146102d35761009c565b80633d4dff7b146100a157806357888e921461010a5780635a161ba51461016157806366db09c6146101a4578063ab80cdc2146101d4575b600080fd5b3480156100ad57600080fd5b506100cb600480360360208110156100c457600080fd5b50356102fd565b604080516001600160a01b039788168152959096166020860152848601939093526060840191909152608083015260a082015290519081900360c00190f35b34801561011657600080fd5b5061014f6004803603608081101561012d57600080fd5b506001600160a01b038135169060208101359060408101359060600135610341565b60408051918252519081900360200190f35b34801561016d57600080fd5b5061018b6004803603602081101561018457600080fd5b50356103e6565b6040805192835260208301919091528051918290030190f35b3480156101b057600080fd5b5061014f600480360360408110156101c757600080fd5b50803590602001356103ff565b6101f1600480360360208110156101ea57600080fd5b503561044b565b005b3480156101ff57600080fd5b5061014f6004803603602081101561021657600080fd5b5035610498565b6101f16004803603606081101561023357600080fd5b506001600160a01b038135169060208101359060400135610547565b34801561025b57600080fd5b5061018b6004803603602081101561027257600080fd5b50356105be565b34801561028557600080fd5b5061014f6004803603602081101561029c57600080fd5b50356106b1565b3480156102af57600080fd5b506101f1600480360360408110156102c657600080fd5b50803590602001356106c3565b3480156102df57600080fd5b506101f1600480360360208110156102f657600080fd5b50356107e6565b60016020819052600091825260409091208054918101546002820154600383015460048401546005909401546001600160a01b039586169590931693919290919086565b6000818152600160208190526040822001546001600160a01b03868116911614158061037e57506000828152600160205260409020600201548414155b80610399575060008281526001602052604090206003015483115b806103bb57506000828152600160205260409020600501544261070719909101105b156103c8575060006103de565b5060008181526001602052604090206004015443035b949350505050565b6000602081905290815260409020805460019091015482565b60008061040b84610498565b60008181526020819052604090205490915083111561042e576000915050610445565b600090815260208190526040902060010154430390505b92915050565b600081815260208190526040808220805434908101825543600190920191909155905181156108fc02919083818181858288f19350505050158015610494573d6000803e3d6000fd5b5050565b6000600282604051602001808281526020019150506040516020818303038152906040526040518082805190602001908083835b602083106104eb5780518252601f1990920191602091820191016104cc565b51815160209384036101000a60001901801990921691161790526040519190930194509192505080830381855afa15801561052a573d6000803e3d6000fd5b5050506040513d602081101561053f57600080fd5b505192915050565b6000818152600160205260409020600401541561056357600080fd5b60009081526001602081905260409091208054336001600160a01b031991821617825591810180549092166001600160a01b03949094169390931790556002820155346003820155436004820155611c204201600590910155565b6000806105c9610bec565b6105d1610bec565b7f216936d3cd6e53fec0a4e231fdd6dc5c692cc7609525a7b2c9562d608f25d51a82527f66666666666666666666666666666666666666666666666666666666666666586020808401919091526001604080850182905260008452918301819052908201525b841561066d5784600116600114156106565761065381836108b3565b90505b600185901c945061066682610a49565b9150610637565b600061067c8260400151610b8c565b90506013600160ff1b03825182900982526013600160ff1b038183602001510960208301819052915194509092505050915091565b60026020526000908152604090205481565b60006106ce82610498565b6000818152600160205260409020600501549091504211156106ef57600080fd5b600081815260016020819052604090912001546001600160a01b0316331461071657600080fd5b8261072057600080fd5b600061072b846105be565b6000848152600160205260409020600201549092508214905061074d57600080fd5b6000818152600260208181526040808420889055858452600180835281852060038101805482546001600160a01b03199081168455838501805490911690559582018790558690556004810186905560050185905591849052808420848155909101839055519091339183156108fc0291849190818181858888f193505050501580156107de573d6000803e3d6000fd5b505050505050565b600081815260016020526040902060050154421161080357600080fd5b6000818152600160205260409020546001600160a01b0316331461082657600080fd5b600081815260016020818152604080842060038101805482546001600160a01b031990811684558387018054909116905560028301879055908690556004820186905560059091018590559184905280842084815590920183905590519091339183156108fc0291849190818181858888f193505050501580156108ae573d6000803e3d6000fd5b505050565b6108bb610bec565b6108c3610c0d565b6013600160ff1b03836040015185604001510981526013600160ff1b038151800960208201526013600160ff1b03835185510960408201526013600160ff1b03836020015185602001510960608201526013600160ff1b038082606001518360400151097f52036cee2b6ffe738cc740797779e89800700a4d4141d8ab75eb4dca135978a30960808201526013600160ff1b0381608001516013600160ff1b030382602001510860a08201526013600160ff1b03816080015182602001510860c08201526013600160ff1b038082606001516013600160ff1b03036013600160ff1b03806109ad57fe5b84604001516013600160ff1b03036013600160ff1b03806109ca57fe5b6013600160ff1b0360208a01518a51086013600160ff1b0360208c01518c51080908086013600160ff1b0360a08401518451090982526013600160ff1b038082604001518360600151086013600160ff1b0360c08401518451090960208301526013600160ff1b038160c001518260a001510960408301525092915050565b610a51610bec565b610a59610c0d565b6013600160ff1b03602084015184510881526013600160ff1b038151800960208201526013600160ff1b038351800960408201526013600160ff1b03602084015180096060820181905260408201516013600160ff1b03908103608084018190529091900860a08201526013600160ff1b036040840151800960e08201526013600160ff1b03808260e001516002096013600160ff1b03038260a001510860c08201526013600160ff1b0360c08201516013600160ff1b0383606001516013600160ff1b03036013600160ff1b0380610b2e57fe5b85604001516013600160ff1b0303866020015108080982526013600160ff1b038082606001516013600160ff1b03038360800151088260a001510960208301526013600160ff1b038160c001518260a0015109604083015250919050565b60008060026013600160ff1b0303905060006013600160ff1b03905060405160208152602080820152602060408201528460608201528260808201528160a082015260208160c0836005600019fa610be357600080fd5b51949350505050565b60405180606001604052806000815260200160008152602001600081525090565b6040518061010001604052806000815260200160008152602001600081526020016000815260200160008152602001600081526020016000815260200160008152509056fea265627a7a72305820efbcbb6c2af9aa94474b41081a8eb70920710a020d23cbd53e7ce082362c569364736f6c63430005090032`

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
// Solidity: function deposits(bytes32 ) constant returns(address sender, address recipient, uint256 adaptorPubKey, uint256 value, uint256 blockNumber, uint256 deadline)
func (_Hub *HubCaller) Deposits(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Sender        common.Address
	Recipient     common.Address
	AdaptorPubKey *big.Int
	Value         *big.Int
	BlockNumber   *big.Int
	Deadline      *big.Int
}, error) {
	ret := new(struct {
		Sender        common.Address
		Recipient     common.Address
		AdaptorPubKey *big.Int
		Value         *big.Int
		BlockNumber   *big.Int
		Deadline      *big.Int
	})
	out := ret
	err := _Hub.contract.Call(opts, out, "deposits", arg0)
	return *ret, err
}

// Deposits is a free data retrieval call binding the contract method 0x3d4dff7b.
//
// Solidity: function deposits(bytes32 ) constant returns(address sender, address recipient, uint256 adaptorPubKey, uint256 value, uint256 blockNumber, uint256 deadline)
func (_Hub *HubSession) Deposits(arg0 [32]byte) (struct {
	Sender        common.Address
	Recipient     common.Address
	AdaptorPubKey *big.Int
	Value         *big.Int
	BlockNumber   *big.Int
	Deadline      *big.Int
}, error) {
	return _Hub.Contract.Deposits(&_Hub.CallOpts, arg0)
}

// Deposits is a free data retrieval call binding the contract method 0x3d4dff7b.
//
// Solidity: function deposits(bytes32 ) constant returns(address sender, address recipient, uint256 adaptorPubKey, uint256 value, uint256 blockNumber, uint256 deadline)
func (_Hub *HubCallerSession) Deposits(arg0 [32]byte) (struct {
	Sender        common.Address
	Recipient     common.Address
	AdaptorPubKey *big.Int
	Value         *big.Int
	BlockNumber   *big.Int
	Deadline      *big.Int
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

// ReclaimDeposit is a paid mutator transaction binding the contract method 0xfa79c259.
//
// Solidity: function reclaimDeposit(bytes32 hashedAntiSpamID) returns()
func (_Hub *HubTransactor) ReclaimDeposit(opts *bind.TransactOpts, hashedAntiSpamID [32]byte) (*types.Transaction, error) {
	return _Hub.contract.Transact(opts, "reclaimDeposit", hashedAntiSpamID)
}

// ReclaimDeposit is a paid mutator transaction binding the contract method 0xfa79c259.
//
// Solidity: function reclaimDeposit(bytes32 hashedAntiSpamID) returns()
func (_Hub *HubSession) ReclaimDeposit(hashedAntiSpamID [32]byte) (*types.Transaction, error) {
	return _Hub.Contract.ReclaimDeposit(&_Hub.TransactOpts, hashedAntiSpamID)
}

// ReclaimDeposit is a paid mutator transaction binding the contract method 0xfa79c259.
//
// Solidity: function reclaimDeposit(bytes32 hashedAntiSpamID) returns()
func (_Hub *HubTransactorSession) ReclaimDeposit(hashedAntiSpamID [32]byte) (*types.Transaction, error) {
	return _Hub.Contract.ReclaimDeposit(&_Hub.TransactOpts, hashedAntiSpamID)
}
