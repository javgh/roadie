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
const HubABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"deprecated\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"deposits\",\"outputs\":[{\"name\":\"sender\",\"type\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\"},{\"name\":\"adaptorPubKey\",\"type\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\"},{\"name\":\"adaptorPubKey\",\"type\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"hashedAntiSpamID\",\"type\":\"bytes32\"}],\"name\":\"checkDepositConfirmations\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"antiSpamFees\",\"outputs\":[{\"name\":\"fee\",\"type\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"servers\",\"outputs\":[{\"name\":\"target\",\"type\":\"string\"},{\"name\":\"cert\",\"type\":\"bytes\"},{\"name\":\"timestamp\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"},{\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"checkAntiSpamConfirmations\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_version\",\"type\":\"string\"}],\"name\":\"setVersion\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"nextServerID\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"target\",\"type\":\"string\"},{\"name\":\"cert\",\"type\":\"bytes\"}],\"name\":\"registerServer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"hashedID\",\"type\":\"bytes32\"}],\"name\":\"burnAntiSpamFee\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"hash\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\"},{\"name\":\"adaptorPubKey\",\"type\":\"uint256\"},{\"name\":\"hashedAntiSpamID\",\"type\":\"bytes32\"}],\"name\":\"depositEther\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"s\",\"type\":\"uint256\"}],\"name\":\"scalarMultBase\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_deprecated\",\"type\":\"bool\"}],\"name\":\"setDeprecated\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"adaptorPrivKeys\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"maxAge\",\"type\":\"uint256\"},{\"name\":\"offset\",\"type\":\"uint256\"}],\"name\":\"fetchServer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"},{\"name\":\"\",\"type\":\"string\"},{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"adaptorPrivKey\",\"type\":\"uint256\"},{\"name\":\"antiSpamID\",\"type\":\"uint256\"}],\"name\":\"claimDeposit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"hashedAntiSpamID\",\"type\":\"bytes32\"}],\"name\":\"reclaimDeposit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// HubBin is the compiled bytecode used for deploying new contracts.
var HubBin = "0x600060045560c0604052600560808190527f302e312e3000000000000000000000000000000000000000000000000000000060a09081526200004391908162000079565b506006805460ff191690553480156200005b57600080fd5b5060068054610100600160a81b03191633610100021790556200011e565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10620000bc57805160ff1916838001178555620000ec565b82800160010185558215620000ec579182015b82811115620000ec578251825591602001919060010190620000cf565b50620000fa929150620000fe565b5090565b6200011b91905b80821115620000fa576000815560010162000105565b90565b61166a806200012e6000396000f3fe60806040526004361061011f5760003560e01c8063ab80cdc2116100a0578063e74db5a911610064578063e74db5a91461064b578063e86ef23b14610675578063ea32a89e1461074d578063f851a4401461077d578063fa79c259146107ae5761011f565b8063ab80cdc21461057c578063b189fd4c14610599578063b90d104d146105c3578063c4f4912b146105f5578063d848dee71461061f5761011f565b80635cf0f357116100e75780635cf0f357146102da57806366db09c6146103e9578063788bc78c1461041957806395fcfa0c146104985780639f64195d146104ad5761011f565b80630e136b19146101245780633d4dff7b1461014d57806354fd4d50146101b657806357888e92146102405780635a161ba514610297575b600080fd5b34801561013057600080fd5b506101396107d8565b604080519115158252519081900360200190f35b34801561015957600080fd5b506101776004803603602081101561017057600080fd5b50356107e1565b604080516001600160a01b039788168152959096166020860152848601939093526060840191909152608083015260a082015290519081900360c00190f35b3480156101c257600080fd5b506101cb610825565b6040805160208082528351818301528351919283929083019185019080838360005b838110156102055781810151838201526020016101ed565b50505050905090810190601f1680156102325780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561024c57600080fd5b506102856004803603608081101561026357600080fd5b506001600160a01b0381351690602081013590604081013590606001356108b3565b60408051918252519081900360200190f35b3480156102a357600080fd5b506102c1600480360360208110156102ba57600080fd5b5035610958565b6040805192835260208301919091528051918290030190f35b3480156102e657600080fd5b50610304600480360360208110156102fd57600080fd5b5035610971565b604051808060200180602001848152602001838103835286818151815260200191508051906020019080838360005b8381101561034b578181015183820152602001610333565b50505050905090810190601f1680156103785780820380516001836020036101000a031916815260200191505b50838103825285518152855160209182019187019080838360005b838110156103ab578181015183820152602001610393565b50505050905090810190601f1680156103d85780820380516001836020036101000a031916815260200191505b509550505050505060405180910390f35b3480156103f557600080fd5b506102856004803603604081101561040c57600080fd5b5080359060200135610ab6565b34801561042557600080fd5b506104966004803603602081101561043c57600080fd5b81019060208101813564010000000081111561045757600080fd5b82018360208201111561046957600080fd5b8035906020019184600183028401116401000000008311171561048b57600080fd5b509092509050610b02565b005b3480156104a457600080fd5b50610285610b2f565b3480156104b957600080fd5b50610496600480360360408110156104d057600080fd5b8101906020810181356401000000008111156104eb57600080fd5b8201836020820111156104fd57600080fd5b8035906020019184600183028401116401000000008311171561051f57600080fd5b91939092909160208101903564010000000081111561053d57600080fd5b82018360208201111561054f57600080fd5b8035906020019184600183028401116401000000008311171561057157600080fd5b509092509050610b35565b6104966004803603602081101561059257600080fd5b5035610b98565b3480156105a557600080fd5b50610285600480360360208110156105bc57600080fd5b5035610be5565b610496600480360360608110156105d957600080fd5b506001600160a01b038135169060208101359060400135610c94565b34801561060157600080fd5b506102c16004803603602081101561061857600080fd5b5035610d0b565b34801561062b57600080fd5b506104966004803603602081101561064257600080fd5b50351515610dfe565b34801561065757600080fd5b506102856004803603602081101561066e57600080fd5b5035610e2d565b34801561068157600080fd5b506106a56004803603604081101561069857600080fd5b5080359060200135610e3f565b60405180841515151581526020018060200180602001838103835285818151815260200191508051906020019080838360005b838110156106f05781810151838201526020016106d8565b50505050905090810190601f16801561071d5780820380516001836020036101000a031916815260200191505b508381038252845181528451602091820191860190808383600083156103ab578181015183820152602001610393565b34801561075957600080fd5b506104966004803603604081101561077057600080fd5b5080359060200135610ffc565b34801561078957600080fd5b5061079261111f565b604080516001600160a01b039092168252519081900360200190f35b3480156107ba57600080fd5b50610496600480360360208110156107d157600080fd5b5035611133565b60065460ff1681565b60016020819052600091825260409091208054918101546002820154600383015460048401546005909401546001600160a01b039586169590931693919290919086565b6005805460408051602060026001851615610100026000190190941693909304601f810184900484028201840190925281815292918301828280156108ab5780601f10610880576101008083540402835291602001916108ab565b820191906000526020600020905b81548152906001019060200180831161088e57829003601f168201915b505050505081565b6000818152600160208190526040822001546001600160a01b0386811691161415806108f057506000828152600160205260409020600201548414155b8061090b575060008281526001602052604090206003015483115b8061092d57506000828152600160205260409020600501544261070719909101105b1561093a57506000610950565b5060008181526001602052604090206004015443035b949350505050565b6000602081905290815260409020805460019091015482565b60036020908152600091825260409182902080548351601f60026000196101006001861615020190931692909204918201849004840281018401909452808452909291839190830182828015610a085780601f106109dd57610100808354040283529160200191610a08565b820191906000526020600020905b8154815290600101906020018083116109eb57829003601f168201915b505050505090806001018054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610aa65780601f10610a7b57610100808354040283529160200191610aa6565b820191906000526020600020905b815481529060010190602001808311610a8957829003601f168201915b5050505050908060020154905083565b600080610ac284610be5565b600081815260208190526040902054909150831115610ae5576000915050610afc565b600090815260208190526040902060010154430390505b92915050565b60065461010090046001600160a01b03163314610b1e57600080fd5b610b2a60058383611534565b505050565b60045481565b6004546000908152600360205260409020610b51908585611534565b506004546000908152600360205260409020610b71906001018383611534565b50506004805460009081526003602052604090204260029091015580546001019055505050565b600081815260208190526040808220805434908101825543600190920191909155905181156108fc02919083818181858288f19350505050158015610be1573d6000803e3d6000fd5b5050565b6000600282604051602001808281526020019150506040516020818303038152906040526040518082805190602001908083835b60208310610c385780518252601f199092019160209182019101610c19565b51815160209384036101000a60001901801990921691161790526040519190930194509192505080830381855afa158015610c77573d6000803e3d6000fd5b5050506040513d6020811015610c8c57600080fd5b505192915050565b60008181526001602052604090206004015415610cb057600080fd5b60009081526001602081905260409091208054336001600160a01b031991821617825591810180549092166001600160a01b03949094169390931790556002820155346003820155436004820155611c204201600590910155565b600080610d166115b2565b610d1e6115b2565b7f216936d3cd6e53fec0a4e231fdd6dc5c692cc7609525a7b2c9562d608f25d51a82527f66666666666666666666666666666666666666666666666666666666666666586020808401919091526001604080850182905260008452918301819052908201525b8415610dba578460011660011415610da357610da081836111fb565b90505b600185901c9450610db382611391565b9150610d84565b6000610dc982604001516114d4565b90506013600160ff1b03825182900982526013600160ff1b038183602001510960208301819052915194509092505050915091565b60065461010090046001600160a01b03163314610e1a57600080fd5b6006805460ff1916911515919091179055565b60026020526000908152604090205481565b60006060806004548410610e7357505060408051602080820183526000808352835191820190935282815291925090610ff5565b60045484900360001901600081815260036020526040902060020154429087011015610ec05750506040805160208082018352600080835283519182019093528281529193509150610ff5565b600081815260036020908152604091829020805483516002600180841615610100026000190190931604601f810185900485028201850190955284815290939192848401928491830182828015610f585780601f10610f2d57610100808354040283529160200191610f58565b820191906000526020600020905b815481529060010190602001808311610f3b57829003601f168201915b5050845460408051602060026001851615610100026000190190941693909304601f810184900484028201840190925281815295975086945092508401905082828015610fe65780601f10610fbb57610100808354040283529160200191610fe6565b820191906000526020600020905b815481529060010190602001808311610fc957829003601f168201915b50505050509050935093509350505b9250925092565b600061100782610be5565b60008181526001602052604090206005015490915042111561102857600080fd5b600081815260016020819052604090912001546001600160a01b0316331461104f57600080fd5b8261105957600080fd5b600061106484610d0b565b6000848152600160205260409020600201549092508214905061108657600080fd5b6000818152600260208181526040808420889055858452600180835281852060038101805482546001600160a01b03199081168455838501805490911690559582018790558690556004810186905560050185905591849052808420848155909101839055519091339183156108fc0291849190818181858888f19350505050158015611117573d6000803e3d6000fd5b505050505050565b60065461010090046001600160a01b031681565b600081815260016020526040902060050154421161115057600080fd5b6000818152600160205260409020546001600160a01b0316331461117357600080fd5b600081815260016020818152604080842060038101805482546001600160a01b031990811684558387018054909116905560028301879055908690556004820186905560059091018590559184905280842084815590920183905590519091339183156108fc0291849190818181858888f19350505050158015610b2a573d6000803e3d6000fd5b6112036115b2565b61120b6115d3565b6013600160ff1b03836040015185604001510981526013600160ff1b038151800960208201526013600160ff1b03835185510960408201526013600160ff1b03836020015185602001510960608201526013600160ff1b038082606001518360400151097f52036cee2b6ffe738cc740797779e89800700a4d4141d8ab75eb4dca135978a30960808201526013600160ff1b0381608001516013600160ff1b030382602001510860a08201526013600160ff1b03816080015182602001510860c08201526013600160ff1b038082606001516013600160ff1b03036013600160ff1b03806112f557fe5b84604001516013600160ff1b03036013600160ff1b038061131257fe5b6013600160ff1b0360208a01518a51086013600160ff1b0360208c01518c51080908086013600160ff1b0360a08401518451090982526013600160ff1b038082604001518360600151086013600160ff1b0360c08401518451090960208301526013600160ff1b038160c001518260a001510960408301525092915050565b6113996115b2565b6113a16115d3565b6013600160ff1b03602084015184510881526013600160ff1b038151800960208201526013600160ff1b038351800960408201526013600160ff1b03602084015180096060820181905260408201516013600160ff1b03908103608084018190529091900860a08201526013600160ff1b036040840151800960e08201526013600160ff1b03808260e001516002096013600160ff1b03038260a001510860c08201526013600160ff1b0360c08201516013600160ff1b0383606001516013600160ff1b03036013600160ff1b038061147657fe5b85604001516013600160ff1b0303866020015108080982526013600160ff1b038082606001516013600160ff1b03038360800151088260a001510960208301526013600160ff1b038160c001518260a0015109604083015250919050565b60008060026013600160ff1b0303905060006013600160ff1b03905060405160208152602080820152602060408201528460608201528260808201528160a082015260208160c0836005600019fa61152b57600080fd5b51949350505050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106115755782800160ff198235161785556115a2565b828001600101855582156115a2579182015b828111156115a2578235825591602001919060010190611587565b506115ae929150611618565b5090565b60405180606001604052806000815260200160008152602001600081525090565b60405180610100016040528060008152602001600081526020016000815260200160008152602001600081526020016000815260200160008152602001600081525090565b61163291905b808211156115ae576000815560010161161e565b9056fea265627a7a723058204f33fc8e9418b12b01999f27372001917f3907dd3ce8141bbbbe0773260b588e64736f6c634300050a0032"

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

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() constant returns(address)
func (_Hub *HubCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Hub.contract.Call(opts, out, "admin")
	return *ret0, err
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() constant returns(address)
func (_Hub *HubSession) Admin() (common.Address, error) {
	return _Hub.Contract.Admin(&_Hub.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() constant returns(address)
func (_Hub *HubCallerSession) Admin() (common.Address, error) {
	return _Hub.Contract.Admin(&_Hub.CallOpts)
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

// Deprecated is a free data retrieval call binding the contract method 0x0e136b19.
//
// Solidity: function deprecated() constant returns(bool)
func (_Hub *HubCaller) Deprecated(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Hub.contract.Call(opts, out, "deprecated")
	return *ret0, err
}

// Deprecated is a free data retrieval call binding the contract method 0x0e136b19.
//
// Solidity: function deprecated() constant returns(bool)
func (_Hub *HubSession) Deprecated() (bool, error) {
	return _Hub.Contract.Deprecated(&_Hub.CallOpts)
}

// Deprecated is a free data retrieval call binding the contract method 0x0e136b19.
//
// Solidity: function deprecated() constant returns(bool)
func (_Hub *HubCallerSession) Deprecated() (bool, error) {
	return _Hub.Contract.Deprecated(&_Hub.CallOpts)
}

// FetchServer is a free data retrieval call binding the contract method 0xe86ef23b.
//
// Solidity: function fetchServer(uint256 maxAge, uint256 offset) constant returns(bool, string, bytes)
func (_Hub *HubCaller) FetchServer(opts *bind.CallOpts, maxAge *big.Int, offset *big.Int) (bool, string, []byte, error) {
	var (
		ret0 = new(bool)
		ret1 = new(string)
		ret2 = new([]byte)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
	}
	err := _Hub.contract.Call(opts, out, "fetchServer", maxAge, offset)
	return *ret0, *ret1, *ret2, err
}

// FetchServer is a free data retrieval call binding the contract method 0xe86ef23b.
//
// Solidity: function fetchServer(uint256 maxAge, uint256 offset) constant returns(bool, string, bytes)
func (_Hub *HubSession) FetchServer(maxAge *big.Int, offset *big.Int) (bool, string, []byte, error) {
	return _Hub.Contract.FetchServer(&_Hub.CallOpts, maxAge, offset)
}

// FetchServer is a free data retrieval call binding the contract method 0xe86ef23b.
//
// Solidity: function fetchServer(uint256 maxAge, uint256 offset) constant returns(bool, string, bytes)
func (_Hub *HubCallerSession) FetchServer(maxAge *big.Int, offset *big.Int) (bool, string, []byte, error) {
	return _Hub.Contract.FetchServer(&_Hub.CallOpts, maxAge, offset)
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

// NextServerID is a free data retrieval call binding the contract method 0x95fcfa0c.
//
// Solidity: function nextServerID() constant returns(uint256)
func (_Hub *HubCaller) NextServerID(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Hub.contract.Call(opts, out, "nextServerID")
	return *ret0, err
}

// NextServerID is a free data retrieval call binding the contract method 0x95fcfa0c.
//
// Solidity: function nextServerID() constant returns(uint256)
func (_Hub *HubSession) NextServerID() (*big.Int, error) {
	return _Hub.Contract.NextServerID(&_Hub.CallOpts)
}

// NextServerID is a free data retrieval call binding the contract method 0x95fcfa0c.
//
// Solidity: function nextServerID() constant returns(uint256)
func (_Hub *HubCallerSession) NextServerID() (*big.Int, error) {
	return _Hub.Contract.NextServerID(&_Hub.CallOpts)
}

// ScalarMultBase is a free data retrieval call binding the contract method 0xc4f4912b.
//
// Solidity: function scalarMultBase(uint256 s) constant returns(uint256, uint256)
func (_Hub *HubCaller) ScalarMultBase(opts *bind.CallOpts, s *big.Int) (*big.Int, *big.Int, error) {
	var (
		ret0 = new(*big.Int)
		ret1 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := _Hub.contract.Call(opts, out, "scalarMultBase", s)
	return *ret0, *ret1, err
}

// ScalarMultBase is a free data retrieval call binding the contract method 0xc4f4912b.
//
// Solidity: function scalarMultBase(uint256 s) constant returns(uint256, uint256)
func (_Hub *HubSession) ScalarMultBase(s *big.Int) (*big.Int, *big.Int, error) {
	return _Hub.Contract.ScalarMultBase(&_Hub.CallOpts, s)
}

// ScalarMultBase is a free data retrieval call binding the contract method 0xc4f4912b.
//
// Solidity: function scalarMultBase(uint256 s) constant returns(uint256, uint256)
func (_Hub *HubCallerSession) ScalarMultBase(s *big.Int) (*big.Int, *big.Int, error) {
	return _Hub.Contract.ScalarMultBase(&_Hub.CallOpts, s)
}

// Servers is a free data retrieval call binding the contract method 0x5cf0f357.
//
// Solidity: function servers(uint256 ) constant returns(string target, bytes cert, uint256 timestamp)
func (_Hub *HubCaller) Servers(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Target    string
	Cert      []byte
	Timestamp *big.Int
}, error) {
	ret := new(struct {
		Target    string
		Cert      []byte
		Timestamp *big.Int
	})
	out := ret
	err := _Hub.contract.Call(opts, out, "servers", arg0)
	return *ret, err
}

// Servers is a free data retrieval call binding the contract method 0x5cf0f357.
//
// Solidity: function servers(uint256 ) constant returns(string target, bytes cert, uint256 timestamp)
func (_Hub *HubSession) Servers(arg0 *big.Int) (struct {
	Target    string
	Cert      []byte
	Timestamp *big.Int
}, error) {
	return _Hub.Contract.Servers(&_Hub.CallOpts, arg0)
}

// Servers is a free data retrieval call binding the contract method 0x5cf0f357.
//
// Solidity: function servers(uint256 ) constant returns(string target, bytes cert, uint256 timestamp)
func (_Hub *HubCallerSession) Servers(arg0 *big.Int) (struct {
	Target    string
	Cert      []byte
	Timestamp *big.Int
}, error) {
	return _Hub.Contract.Servers(&_Hub.CallOpts, arg0)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(string)
func (_Hub *HubCaller) Version(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Hub.contract.Call(opts, out, "version")
	return *ret0, err
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(string)
func (_Hub *HubSession) Version() (string, error) {
	return _Hub.Contract.Version(&_Hub.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(string)
func (_Hub *HubCallerSession) Version() (string, error) {
	return _Hub.Contract.Version(&_Hub.CallOpts)
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

// RegisterServer is a paid mutator transaction binding the contract method 0x9f64195d.
//
// Solidity: function registerServer(string target, bytes cert) returns()
func (_Hub *HubTransactor) RegisterServer(opts *bind.TransactOpts, target string, cert []byte) (*types.Transaction, error) {
	return _Hub.contract.Transact(opts, "registerServer", target, cert)
}

// RegisterServer is a paid mutator transaction binding the contract method 0x9f64195d.
//
// Solidity: function registerServer(string target, bytes cert) returns()
func (_Hub *HubSession) RegisterServer(target string, cert []byte) (*types.Transaction, error) {
	return _Hub.Contract.RegisterServer(&_Hub.TransactOpts, target, cert)
}

// RegisterServer is a paid mutator transaction binding the contract method 0x9f64195d.
//
// Solidity: function registerServer(string target, bytes cert) returns()
func (_Hub *HubTransactorSession) RegisterServer(target string, cert []byte) (*types.Transaction, error) {
	return _Hub.Contract.RegisterServer(&_Hub.TransactOpts, target, cert)
}

// SetDeprecated is a paid mutator transaction binding the contract method 0xd848dee7.
//
// Solidity: function setDeprecated(bool _deprecated) returns()
func (_Hub *HubTransactor) SetDeprecated(opts *bind.TransactOpts, _deprecated bool) (*types.Transaction, error) {
	return _Hub.contract.Transact(opts, "setDeprecated", _deprecated)
}

// SetDeprecated is a paid mutator transaction binding the contract method 0xd848dee7.
//
// Solidity: function setDeprecated(bool _deprecated) returns()
func (_Hub *HubSession) SetDeprecated(_deprecated bool) (*types.Transaction, error) {
	return _Hub.Contract.SetDeprecated(&_Hub.TransactOpts, _deprecated)
}

// SetDeprecated is a paid mutator transaction binding the contract method 0xd848dee7.
//
// Solidity: function setDeprecated(bool _deprecated) returns()
func (_Hub *HubTransactorSession) SetDeprecated(_deprecated bool) (*types.Transaction, error) {
	return _Hub.Contract.SetDeprecated(&_Hub.TransactOpts, _deprecated)
}

// SetVersion is a paid mutator transaction binding the contract method 0x788bc78c.
//
// Solidity: function setVersion(string _version) returns()
func (_Hub *HubTransactor) SetVersion(opts *bind.TransactOpts, _version string) (*types.Transaction, error) {
	return _Hub.contract.Transact(opts, "setVersion", _version)
}

// SetVersion is a paid mutator transaction binding the contract method 0x788bc78c.
//
// Solidity: function setVersion(string _version) returns()
func (_Hub *HubSession) SetVersion(_version string) (*types.Transaction, error) {
	return _Hub.Contract.SetVersion(&_Hub.TransactOpts, _version)
}

// SetVersion is a paid mutator transaction binding the contract method 0x788bc78c.
//
// Solidity: function setVersion(string _version) returns()
func (_Hub *HubTransactorSession) SetVersion(_version string) (*types.Transaction, error) {
	return _Hub.Contract.SetVersion(&_Hub.TransactOpts, _version)
}
