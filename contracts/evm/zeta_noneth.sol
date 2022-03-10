// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";

contract Zeta is ERC20Burnable {
    // TSSAddress is the TSS address collectively possessed by zeta blockchain validators. 
    // It's the only address that can mint. Threshold Signature Scheme (TSS) [GG20] is a multi-sig ECDSA/EdDSA protocol. 
    address public TSSAddress; 
    // The TSSAddressUpdater can change the TSSAddress, in case Zeta blockchain node churn changes the TSS key pairs. 
    // At launch, TSSAddressUpdater is controlled by multi-sig wallet; Eventually, TSSAddressUpdater will be the same as TSSAddress, via renounceTSSAddressUpdater()
    address public TSSAddressUpdater;
    // Message Passing Interface address; can mint.
    address public MPIAddress; 

    // nonces for permit_* functions; prevent replay attack. 
    mapping(address => uint) public nonces;

    event MMinted(address indexed mintee, uint amount, bytes32 indexed sendHash);
    event MBurnt(address indexed burnee, uint amount); 

    constructor(uint256 initialSupply, string memory name, string memory symbol,  address _TSSAddress, address _TSSAddressUpdater) ERC20(name, symbol) {
        _mint(msg.sender, initialSupply * (10 ** uint256(decimals())));
        TSSAddress = _TSSAddress; 
        TSSAddressUpdater = _TSSAddressUpdater; 
    }

    // update the TSSAddress in case of Zeta blockchain validator nodes churn
    function updateTSSAndMPIAddresses(address _tss, address _mpi) external {
        require(msg.sender == TSSAddressUpdater, "updateTSSAddress: need TSSAddressUpdater permission");
        TSSAddress = _tss;
        MPIAddress = _mpi; 
    }

    // Change the ownership of TSSAddressUpdater to the Zeta blockchain TSS nodes. 
    // Effectively, only Zeta blockchain validators collectively can update TSSAddress afterwards. 
    function renounceTSSAddressUpdater() external {
        require(msg.sender == TSSAddressUpdater, "renounceTSSAddressUpdater: need TSSAddressUpdater permission");
        require(TSSAddress != address(0));
        TSSAddressUpdater = TSSAddress;
    }

    // TSS can mint to anyone with commensurate burns on other chains. Only TSS/MPI address can mint. 
    function mint(address mintee, uint value, bytes32 sendHash) external {
        require(msg.sender == TSSAddress || msg.sender == MPIAddress, "Only TSSAddress can mint"); 
        _mint(mintee, value); 
        emit MMinted(mintee, value, sendHash);
    }

}