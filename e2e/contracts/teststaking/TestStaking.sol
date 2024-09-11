// SPDX-License-Identifier: MIT
pragma solidity 0.8.10;

// @dev Interface for IStaking precompile for easier import
enum BondStatus {
    Unspecified,
    Unbonded,
    Unbonding,
    Bonded
}

struct Validator {
    string operatorAddress;
    string consensusPubKey;
    bool jailed;
    BondStatus bondStatus;
}

interface IStaking {
    function stake(
        address staker,
        string memory validator,
        uint256 amount
    ) external returns (bool success);

    function unstake(
        address staker,
        string memory validator,
        uint256 amount
    ) external returns (int64 completionTime);

    function moveStake(
        address staker,
        string memory validatorSrc,
        string memory validatorDst,
        uint256 amount
    ) external returns (int64 completionTime);

    function getAllValidators() external view returns (Validator[] calldata validators);

    function getShares(address staker, string memory validator) external view returns (uint256 shares);
}

interface WZETA {
    function deposit() external payable;
    function withdraw(uint256 wad) external;
}

// @dev Purpose of this contract is to call staking precompile
contract TestStaking {
    event Stake(
        address indexed staker,
        address indexed validator,
        uint256 amount
    );

    event Unstake(
        address indexed staker,
        address indexed validator,
        uint256 amount
    );

    event MoveStake(
        address indexed staker,
        address indexed validatorSrc,
        address indexed validatorDst,
        uint256 amount
    );

    IStaking staking = IStaking(0x0000000000000000000000000000000000000066);
    WZETA wzeta;
    address owner;

    // @dev used to test state change in smart contract
    uint256 public counter = 0;

    constructor(address _wzeta) {
        wzeta = WZETA(_wzeta);
        owner = msg.sender;
    }

    // simple protection to not be able to call contract by anyone
    // not relevant for e2e tests
    modifier onlyOwner() {
        require(msg.sender == owner);
        _;
    }

    function depositWZETA() external payable onlyOwner {
        wzeta.deposit{value: msg.value}();
    }

    function withdrawWZETA(uint256 wad) external onlyOwner {
        wzeta.withdraw(wad);
    }

    function stake(address staker, string memory validator, uint256 amount) external onlyOwner returns (bool)  {
        return staking.stake(staker, validator, amount);
    }

    function stakeWithStateUpdate(address staker, string memory validator, uint256 amount) external onlyOwner returns (bool)  {
        counter = counter + 1;
        bool success = staking.stake(staker, validator, amount);
        counter = counter + 1;
        return success;
    }

    function stakeAndRevert(address staker, string memory validator, uint256 amount) external onlyOwner returns (bool)  {
        counter = counter + 1;
        staking.stake(staker, validator, amount);
        counter = counter + 1;
        revert("testrevert");
    }

    function unstake(
        address staker,
        string memory validator,
        uint256 amount
    ) external onlyOwner returns (int64 completionTime) {
        return staking.unstake(staker, validator, amount);
    }

    function moveStake(
        address staker,
        string memory validatorSrc,
        string memory validatorDst,
        uint256 amount
    ) external onlyOwner returns (int64 completionTime) {
        return staking.moveStake(staker, validatorSrc, validatorDst, amount);
    }

    function getShares(address staker, string memory validator) external view returns(uint256 shares) {
        return staking.getShares(staker, validator);
    }

    function getAllValidators() external view returns (Validator[] memory validators) {
        return staking.getAllValidators();
    }

    fallback() external payable {}

    receive() external payable {}
}