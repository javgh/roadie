pragma solidity ^0.5.9;

import "./Ed25519.sol";

contract Hub is Ed25519 {
    address payable constant BLACK_HOLE = 0x0000000000000000000000000000000000000000;
    uint constant DEPOSIT_DURATION = 2 hours;
    uint constant DEPOSIT_DURATION_MARGIN = 30 minutes;

    struct AntiSpamFee {
        uint fee;
        uint blockNumber;
    }

    struct Deposit {
        address sender;
        address recipient;
        uint adaptorPubKey;
        uint value;
        uint blockNumber;
        uint deadline;
    }

    struct Server {
        string target;
        bytes cert;
        uint timestamp;
    }

    mapping(bytes32 => AntiSpamFee) public antiSpamFees;
    mapping(bytes32 => Deposit) public deposits;
    mapping(uint => uint) public adaptorPrivKeys;

    mapping(uint => Server) public servers;
    uint public nextServerID = 0;

    string public version = "0.1.0";
    bool public deprecated = false;
    address public admin;

    modifier onlyAdmin {
        require(msg.sender == admin);
        _;
    }

    constructor() public {
        admin = msg.sender;
    }

    function burnAntiSpamFee(bytes32 hashedID) external payable {
        antiSpamFees[hashedID].fee += msg.value;
        antiSpamFees[hashedID].blockNumber = block.number;
        BLACK_HOLE.transfer(msg.value);
    }

    function checkAntiSpamConfirmations(uint id, uint fee) external view returns (uint) {
        bytes32 hashedID = hash(id);

        if (antiSpamFees[hashedID].fee < fee) {
            return 0;
        } else {
            return block.number - antiSpamFees[hashedID].blockNumber;
        }
    }

    function depositEther(address recipient, uint adaptorPubKey, bytes32 hashedAntiSpamID) external payable {
        require(deposits[hashedAntiSpamID].blockNumber == 0);

        deposits[hashedAntiSpamID].sender = msg.sender;
        deposits[hashedAntiSpamID].recipient = recipient;
        deposits[hashedAntiSpamID].adaptorPubKey = adaptorPubKey;
        deposits[hashedAntiSpamID].value = msg.value;
        deposits[hashedAntiSpamID].blockNumber = block.number;
        deposits[hashedAntiSpamID].deadline = now + DEPOSIT_DURATION;
    }

    function checkDepositConfirmations(address recipient, uint adaptorPubKey,
                                       uint value, bytes32 hashedAntiSpamID) external view returns (uint) {
        if (deposits[hashedAntiSpamID].recipient != recipient ||
            deposits[hashedAntiSpamID].adaptorPubKey != adaptorPubKey ||
            deposits[hashedAntiSpamID].value < value ||
            deposits[hashedAntiSpamID].deadline - DEPOSIT_DURATION_MARGIN < now) {
            return 0;
        } else {
            return block.number - deposits[hashedAntiSpamID].blockNumber;
        }
    }

    function claimDeposit(uint adaptorPrivKey, uint antiSpamID) external {
        bytes32 hashedAntiSpamID = hash(antiSpamID);
        require(deposits[hashedAntiSpamID].deadline >= now);
        require(deposits[hashedAntiSpamID].recipient == msg.sender);
        require(adaptorPrivKey != 0);

        (, uint adaptorPubKey) = scalarmult(adaptorPrivKey);    // check via Ed25519.sol
        require(deposits[hashedAntiSpamID].adaptorPubKey == adaptorPubKey);
        adaptorPrivKeys[adaptorPubKey] = adaptorPrivKey;

        uint value = deposits[hashedAntiSpamID].value;
        delete deposits[hashedAntiSpamID];
        delete antiSpamFees[hashedAntiSpamID];
        msg.sender.transfer(value);
    }

    function reclaimDeposit(bytes32 hashedAntiSpamID) external {
        require(deposits[hashedAntiSpamID].deadline < now);
        require(deposits[hashedAntiSpamID].sender == msg.sender);

        uint value = deposits[hashedAntiSpamID].value;
        delete deposits[hashedAntiSpamID];
        delete antiSpamFees[hashedAntiSpamID];
        msg.sender.transfer(value);
    }

    function registerServer(string calldata target, bytes calldata cert) external {
        servers[nextServerID].target = target;
        servers[nextServerID].cert = cert;
        servers[nextServerID].timestamp = now;
        nextServerID += 1;
    }

    function fetchServer(uint maxAge,
                         uint offset) external view
                         returns (bool, string memory, bytes memory) {
        if (offset >= nextServerID) {
            return (false, "", "");
        }

        uint id = nextServerID - offset - 1;
        if (servers[id].timestamp + maxAge < now) {
            return (false, "", "");
        }

        return (true, servers[id].target, servers[id].cert);
    }

    function hash(uint id) public pure returns (bytes32) {
        return sha256(abi.encode(id));
    }

    function setVersion(string calldata _version) external onlyAdmin {
        version = _version;
    }

    function setDeprecated(bool _deprecated) external onlyAdmin {
        deprecated = _deprecated;
    }
}
