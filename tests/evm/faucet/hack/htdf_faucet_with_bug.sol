/*
yqq 2020-12-11

a faucet contract with a bug (re-entrancy )

compile contract:
solcjs  --bin --abi hack_faucet.sol  sscq_faucet_with_bug.sol


*/

pragma solidity ^0.4.20;



contract SscqFaucet {
    
    uint256 public onceAmount;
    address public owner ;
    
    event SendSscq(address indexed toAddress, uint256 indexed amount);
    event Deposit(address indexed fromAddress, uint256 indexed amount);
    event SetOnceAmount(address indexed fromAddress, uint256 indexed amount);
    mapping (address => uint256) sendRecords;
    
    function SscqFaucet() public payable{
        onceAmount = 100000000;
        owner = msg.sender;
    }
    
    modifier onlyOwner {
        require(msg.sender == owner);
        _;
    }
    
    function setOnceAmount(uint256 amount) public onlyOwner {
        onceAmount = amount;
        SetOnceAmount(msg.sender, amount);
    }
    
    function getOneSscq() public {
        require( sendRecords[msg.sender] == 0 || 
            (sendRecords[msg.sender] > 0 &&  now - sendRecords[msg.sender] > 1 minutes ));
            
        require(address(this).balance >= onceAmount);
        
        // NOTE: THIS IS UNSAFE
        msg.sender.call.value( onceAmount )("");
        sendRecords[msg.sender] = now; // NOTE: probobaly be re-entrancy attacked
        SendSscq(msg.sender, onceAmount);
    }
    
    function deposit() public payable {
        Deposit(msg.sender, msg.value);
    }
    
    // function() public payable{
        
    // }
    
}