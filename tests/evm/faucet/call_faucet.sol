/*
yqq  2020-12-11
test contract call contract
*/

pragma solidity ^0.4.20;

import "./sscq_faucet.sol";

contract CallFaucet {

    SscqFaucet public faucet;
    uint256 public stackDepth = 0;
    address public addr;
    address public owner;

    function CallFaucet() public payable {
        // hard coding the faucet contract address
        addr = address(0x18EDA861679664967c067bA3068414339E5B49e9);
        faucet = SscqFaucet(addr);
        owner = msg.sender;
    }

    function  getOneSscq() public {
        stackDepth = 0;
        faucet.getOneSscq(); // it's ok
    }

    // fallback function
    function () external payable {
       
    }

}