// SPDX-License-Identifier: MIT
pragma solidity 0.8.7;

import "../interfaces/zContract.sol";
import "../interfaces/IZRC20.sol";

/**
 * @dev Custom errors for SystemContract
 */
interface SystemContractErrors {
    error CallerIsNotFungibleModule();

    error InvalidTarget();

    error CantBeIdenticalAddresses();

    error CantBeZeroAddress();
}

/**
 * @dev The system contract it's called by the protocol to interact with the blockchain.
 * Also includes a lot of tools to make easier to interact with ZetaChain.
 */
contract SystemContract is SystemContractErrors {
    // @dev: Map to know the gas price of each chain given a chain id.
    mapping(uint256 => uint256) public gasPriceByChainId;
    // @dev: Map to know the ZRC20 address of a token given a chain id, ex zETH, zBNB etc.
    mapping(uint256 => address) public gasCoinZRC20ByChainId;
    // @dev: Map to know uniswap pool of ZETA/ZRC20 given a chain id. This refer to the build in uniswap deployed at genesis.
    mapping(uint256 => address) public gasZetaPoolByChainId;

    // @dev: Fungible address is always the same, it's on protocol level
    address public constant FUNGIBLE_MODULE_ADDRESS = 0x735b14BB79463307AAcBED86DAf3322B1e6226aB;
    // @dev: Address of the wrapped ZETA to interact with Uniswap
    address public wZetaContractAddress;
    // @dev: Uniswap addresses
    address public uniswapv2FactoryAddress;
    address public uniswapv2Router02Address;
    address public zetaConnectorZEVMAddress;

    event SystemContractDeployed();
    event SetGasPrice(uint256, uint256);
    event SetGasCoin(uint256, address);
    event SetGasZetaPool(uint256, address);
    event SetWZeta(address);
    event SetConnectorZEVM(address);

    /**
     * @dev Only fungible module can deploy a system contract.
     */
    constructor(
        address wzeta_,
        address uniswapv2Factory_,
        address uniswapv2Router02_
    ) {
        if (msg.sender != FUNGIBLE_MODULE_ADDRESS) revert CallerIsNotFungibleModule();
        wZetaContractAddress = wzeta_;
        uniswapv2FactoryAddress = uniswapv2Factory_;
        uniswapv2Router02Address = uniswapv2Router02_;
        emit SystemContractDeployed();
    }

    // deposit foreign coins into ZRC20 and call user specified contract on zEVM
    function depositAndCall(
        address zrc20,
        uint256 amount,
        address target,
        bytes calldata message
    ) external {
        if (msg.sender != FUNGIBLE_MODULE_ADDRESS) revert CallerIsNotFungibleModule();
        if (target == FUNGIBLE_MODULE_ADDRESS || target == address(this)) revert InvalidTarget();

        IZRC20(zrc20).deposit(target, amount);
        zContract(target).onCrossChainCall(zrc20, amount, message);
    }

    // returns sorted token addresses, used to handle return values from pairs sorted in this order.
    function sortTokens(address tokenA, address tokenB) internal pure returns (address token0, address token1) {
        if (tokenA == tokenB) revert CantBeIdenticalAddresses();
        (token0, token1) = tokenA < tokenB ? (tokenA, tokenB) : (tokenB, tokenA);
        if (token0 == address(0)) revert CantBeZeroAddress();
    }

    // calculates the CREATE2 address for a pair without making any external calls.
    function uniswapv2PairFor(
        address factory,
        address tokenA,
        address tokenB
    ) public pure returns (address pair) {
        (address token0, address token1) = sortTokens(tokenA, tokenB);
        pair = address(
            uint160(
                uint256(
                    keccak256(
                        abi.encodePacked(
                            hex"ff",
                            factory,
                            keccak256(abi.encodePacked(token0, token1)),
                            hex"96e8ac4277198ff8b6f785478aa9a39f403cb768dd02cbee326c3e7da348845f" // init code hash
                        )
                    )
                )
            )
        );
    }

    // fungible module updates the gas price oracle periodically.
    function setGasPrice(uint256 chainID, uint256 price) external {
        if (msg.sender != FUNGIBLE_MODULE_ADDRESS) revert CallerIsNotFungibleModule();
        gasPriceByChainId[chainID] = price;
        emit SetGasPrice(chainID, price);
    }

    // @dev: Setter for gasCoinZRC20ByChainId map.
    function setGasCoinZRC20(uint256 chainID, address zrc20) external {
        if (msg.sender != FUNGIBLE_MODULE_ADDRESS) revert CallerIsNotFungibleModule();
        gasCoinZRC20ByChainId[chainID] = zrc20;
        emit SetGasCoin(chainID, zrc20);
    }

    // set the pool wzeta/erc20 address.
    function setGasZetaPool(uint256 chainID, address erc20) external {
        if (msg.sender != FUNGIBLE_MODULE_ADDRESS) revert CallerIsNotFungibleModule();
        address pool = uniswapv2PairFor(uniswapv2FactoryAddress, wZetaContractAddress, erc20);
        gasZetaPoolByChainId[chainID] = pool;
        emit SetGasZetaPool(chainID, pool);
    }

    // @dev: Setter for wrapped ZETA address.
    function setWZETAContractAddress(address addr) external {
        if (msg.sender != FUNGIBLE_MODULE_ADDRESS) revert CallerIsNotFungibleModule();
        wZetaContractAddress = addr;
        emit SetWZeta(wZetaContractAddress);
    }

    function setConnectorZEVMAddress(address addr) external {
        if (msg.sender != FUNGIBLE_MODULE_ADDRESS) revert CallerIsNotFungibleModule();
        zetaConnectorZEVMAddress = addr;
        emit SetConnectorZEVM(zetaConnectorZEVMAddress);
    }
}
