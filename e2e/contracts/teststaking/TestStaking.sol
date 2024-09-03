// SPDX-License-Identifier: MIT
pragma solidity 0.8.7;

// @dev Interface for IStaking precompile for easier import
interface IStaking {
    function stake(
        address staker,
        string memory validator,
        uint256 amount
    ) external returns (bool success);

    function getShares(address staker, string memory validator) external view returns (uint256 shares);
}

interface IPrototype {
    /// @dev converting a bech32 address to hexadecimal address.
    /// @param bech32 The bech32 address.
    /// @return addr The hexadecimal address.
    function bech32ToHexAddr(
        string memory bech32
    ) external view returns (address addr);
}

// @dev Purpose of this contract is to call staking precompile
// and test permissions when delegator is not calling staking precompile directly
contract TestStaking {
    IStaking staking = IStaking(0x0000000000000000000000000000000000000066);
    address pra = 0x0000000000000000000000000000000000000065;
    IPrototype pr = IPrototype(0x0000000000000000000000000000000000000065);

    function stake(string memory validator, uint256 amount) external {
        bool success = staking.stake(msg.sender, validator, amount);
        require(success == true, "staking failed");
    }

    function getShares(address staker, string memory validator) external view returns(uint256 shares) {
        return staking.getShares(staker, validator);
    }

    function bech32Fn(string memory bech32) external view returns (address addr) {
        return pr.bech32ToHexAddr(bech32);
    }

    function bech32StaticFn(string memory bech32) external view returns (bool) {
        (bool success, ) = pra.staticcall(abi.encodeWithSignature("bech32ToHexAddr(string)", bech32));
        return success;
    }

     function bech32CallFn(string memory bech32) external returns (bool, bytes memory) {
        (bool success, bytes memory data) = pra.call(abi.encodeWithSignature("bech32ToHexAddr(string)", bech32));
        return (success, data);
    }
}