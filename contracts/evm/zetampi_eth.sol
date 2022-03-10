// SPDX-License-Identifier: MIT
// v2.1, 2022-02-23
pragma solidity ^0.8.4;
// This ERC20 interface comes from OpenZeppelin
// https://github.com/OpenZeppelin/openzeppelin-contracts/blob/master/contracts/token/ERC20/IERC20.sol
interface IERC20 {
    function transfer(address to, uint256 amount) external returns (bool);
    function allowance(address owner, address spender) external view returns (uint256);
    function approve(address spender, uint256 amount) external returns (bool);
    function transferFrom( address from, address to, uint256 amount) external returns (bool);
    function balanceOf(address account) external view returns (uint256);
}

interface ZetaMPIReceiver {
	function uponZetaMessage(bytes calldata sender, uint16 srcChainID, address destContract, uint zetaAmount, bytes calldata message) external; 
}

contract ZetaMPI {
    address public ZETA_TOKEN; // the Zeta token contract
    bool public paused;

    // TSSAddress is the TSS address collectively possessed by Zeta blockchain validators. 
    // Threshold Signature Scheme (TSS) [GG20] is a multi-sig ECDSA/EdDSA protocol. 
    address public TSSAddress; 
    address public TSSAddressUpdater;
    
    event ZetaMessageSendEvent(address indexed sender, uint16 destChainID, bytes destContract, uint zetaAmount, uint gasLimit, bytes message, bytes zetaParams); 
    event ZetaMessageReceiveEvent(bytes sender, uint16 indexed srcChainID, address indexed destContract, uint zetaAmount, bytes message, bytes32 indexed sendHash); 
    event Paused(address sender);
    event Unpaused(address sender);

    constructor(address zetaAddress,  address _TSSAddress, address _TSSAddressUpdater) {       
        ZETA_TOKEN = zetaAddress;
        TSSAddress = _TSSAddress; 
        TSSAddressUpdater = _TSSAddressUpdater; 
        paused = false; 
    }

    // update the TSSAddress in case of Zeta blockchain validator nodes churn
    function updateTSSAddress(address _address) external {
        require(msg.sender == TSSAddressUpdater, "updateTSSAddress: need TSSAddressUpdater permission");
        require(_address != address(0)); 
        TSSAddress = _address;
    }

    // Change the ownership of TSSAddressUpdater to the Zeta blockchain TSS nodes. 
    // Effectively, only Zeta blockchain validators collectively can update TSSAddress afterwards. 
    function renounceTSSAddressUpdater() external {
        require(msg.sender == TSSAddressUpdater, "renounceTSSAddressUpdater: need TSSAddressUpdater permission");
        require(TSSAddress != address(0)); 
        TSSAddressUpdater = TSSAddress;
    }

    function getLockedAmount() public view returns (uint) {
        return IERC20(ZETA_TOKEN).balanceOf(address(this));
    }

    function pause() external {
        require(paused == false, "already paused");
        require(msg.sender == TSSAddressUpdater); 
        paused = true;
        emit Paused(msg.sender);
    }
    function unpause() external {
        require(paused == true, "already unpaused");
        require(msg.sender == TSSAddressUpdater); 
        paused = false;
        emit Unpaused(msg.sender);
    }

    function zetaMessageSend(uint16 destChainID, bytes calldata  destContract, uint zetaAmount, uint gasLimit, bytes calldata message, bytes calldata zetaParams) external {
        require(paused == false, "paused"); 
        bool success = IERC20(ZETA_TOKEN).transferFrom(msg.sender, address(this), zetaAmount); 
        require(success == true, "zetaMessageSend: transfer fails"); 
        emit ZetaMessageSendEvent(msg.sender, destChainID, destContract, zetaAmount, gasLimit, message, zetaParams); 
    }

    function zetaMessageReceive(bytes calldata srcContract, uint16 srcChainID, address destContract, uint zetaAmount, bytes calldata message, bytes32 sendHash) external {
        require(paused == false, "paused"); 
        require(msg.sender == TSSAddress, "zetaMessageReceive: permission error"); 
        bool success = IERC20(ZETA_TOKEN).transfer(destContract, zetaAmount);
        require(success == true, "zetaMessageReceive: transfer failed"); 
        if (message.length > 0) 
            ZetaMPIReceiver(destContract).uponZetaMessage(srcContract, srcChainID, destContract, zetaAmount, message);
        emit ZetaMessageReceiveEvent(srcContract, srcChainID, destContract, zetaAmount, message, sendHash);
    }
}