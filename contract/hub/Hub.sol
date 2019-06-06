pragma solidity ^0.5.9;

import "./Ed25519.sol";

contract Hub is Ed25519 {
    address payable constant BLACK_HOLE = 0x0000000000000000000000000000000000000000;

    struct AntiSpamFee {
        uint fee;
        uint blockNumber;
    }

    struct Deposit {
        address recipient;
        uint adaptorPubKey;
        uint value;
        uint blockNumber;
    }

    mapping(bytes32 => AntiSpamFee) public antiSpamFees;
    mapping(bytes32 => Deposit) public deposits;
    mapping(uint => uint) public adaptorPrivKeys;

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
            return block.number - antiSpamFees[hashedID].blockNumber + 1;
        }
    }

    function depositEther(address recipient, uint adaptorPubKey, bytes32 hashedAntiSpamID) external payable {
        require(deposits[hashedAntiSpamID].blockNumber == 0);

        deposits[hashedAntiSpamID].recipient = recipient;
        deposits[hashedAntiSpamID].adaptorPubKey = adaptorPubKey;
        deposits[hashedAntiSpamID].value = msg.value;
        deposits[hashedAntiSpamID].blockNumber = block.number;
    }

    function checkDepositConfirmations(address recipient, uint adaptorPubKey,
                                       uint value, bytes32 hashedAntiSpamID) external view returns (uint) {
        if (deposits[hashedAntiSpamID].recipient != recipient ||
            deposits[hashedAntiSpamID].adaptorPubKey != adaptorPubKey ||
            deposits[hashedAntiSpamID].value < value) {
            return 0;
        } else {
            return block.number - deposits[hashedAntiSpamID].blockNumber + 1;
        }
    }

    function claimDeposit(uint adaptorPrivKey, uint antiSpamID) external {
        bytes32 hashedAntiSpamID = hash(antiSpamID);
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

    function hash(uint id) public pure returns (bytes32) {
        return sha256(abi.encode(id));
    }
}
