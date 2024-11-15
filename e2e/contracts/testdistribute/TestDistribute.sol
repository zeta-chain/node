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
    event Distributed(
        address indexed zrc20_distributor,
        address indexed zrc20_token,
        uint256 amount
    );

    IDistribute distr = IDistribute(0x0000000000000000000000000000000000000066);

    address immutable owner;

    constructor() {
        owner = msg.sender;
    }

    modifier onlyOwner() {
        require(msg.sender == owner);
        _;
    }

    function distributeThroughContract(
        address zrc20,
        uint256 amount
    ) external onlyOwner returns (bool) {
        return distr.distribute(zrc20, amount);
    }

    fallback() external payable {}

    receive() external payable {}
}
