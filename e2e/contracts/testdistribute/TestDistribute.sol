// SPDX-License-Identifier: MIT
pragma solidity 0.8.10;

struct DecCoin {
    string denom;
    uint256 amount;
}

// @dev Interface to interact with distribute.
interface IDistribute {
    function distribute(
        address zrc20,
        uint256 amount
    ) external returns (bool success);

    function claimRewards(
        address delegator,
        string memory validator
    ) external returns (bool success);

    function getDelegatorValidators(
        address delegator
    ) external view returns (string[] calldata validators);

    function getRewards(
        address delegator,
        string memory validator
    ) external view returns (DecCoin[] calldata rewards);
}

// @dev Call IBank contract functions
contract TestDistribute {
    IDistribute distr = IDistribute(0x0000000000000000000000000000000000000066);

    fallback() external payable {}

    receive() external payable {}

    function distributeThroughContract(
        address zrc20,
        uint256 amount
    ) external returns (bool) {
        return distr.distribute(zrc20, amount);
    }

    function claimRewardsThroughContract(
        address delegator,
        string memory validator
    ) external returns (bool) {
        return distr.claimRewards(delegator, validator);
    }

    function getDelegatorValidatorsThroughContract(
        address delegator
    ) external view returns (string[] memory) {
        return distr.getDelegatorValidators(delegator);
    }

    function getRewardsThroughContract(
        address delegator,
        string memory validator
    ) external view returns (DecCoin[] memory) {
        return distr.getRewards(delegator, validator);
    }
}
