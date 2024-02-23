// This contract is intentionally vulnerable for testing purposes

pragma solidity ^0.8.0;

contract VulnerableContract {
    mapping(address => uint) private balances;

    function deposit() public payable {
        balances[msg.sender] += msg.value;
    }

    function withdraw(uint _amount) public {
        require(balances[msg.sender] >= _amount, "Insufficient balance");
        
        // Vulnerability: Reentrancy
        (bool success, ) = msg.sender.call{value: _amount}("");
        require(success, "Transfer failed");

        balances[msg.sender] -= _amount;
    }

    function getBalance() public view returns (uint) {
        return balances[msg.sender];
    }
}
