// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

/**
 * @dev Custom errors for ZRC20
 */
interface ZRC20Errors {
    error CallerIsNotFungibleModule();
    error InvalidSender();
    error GasFeeTransferFailed();
    error ZeroGasCoin();
    error ZeroGasPrice();
    error LowAllowance();
    error LowBalance();
    error ZeroAddress();
}

/**
 * @dev Interfaces of SystemContract and ZRC20 to make easier to import.
 */
interface ISystem {
    function FUNGIBLE_MODULE_ADDRESS() external view returns (address);
    function wZetaContractAddress() external view returns (address);
    function uniswapv2FactoryAddress() external view returns (address);
    function gasPriceByChainId(uint256 chainID) external view returns (uint256);
    function gasCoinZRC20ByChainId(uint256 chainID) external view returns (address);
    function gasZetaPoolByChainId(uint256 chainID) external view returns (address);
}

interface IZRC20 {
    function totalSupply() external view returns (uint256);
    function balanceOf(address account) external view returns (uint256);
    function transfer(address recipient, uint256 amount) external returns (bool);
    function allowance(address owner, address spender) external view returns (uint256);
    function approve(address spender, uint256 amount) external returns (bool);
    function transferFrom(address sender, address recipient, uint256 amount) external returns (bool);
    function deposit(address to, uint256 amount) external returns (bool);
    function withdraw(bytes memory to, uint256 amount) external returns (bool);
    function withdrawGasFee() external view returns (address, uint256);

    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);
    event Deposit(bytes from, address indexed to, uint256 value);
    event Withdrawal(address indexed from, bytes to, uint256 value, uint256 gasfee, uint256 protocolFlatFee);
    event UpdatedSystemContract(address systemContract);
    event UpdatedGasLimit(uint256 gasLimit);
    event UpdatedProtocolFlatFee(uint256 protocolFlatFee);
}

interface IZRC20Metadata is IZRC20 {
    function name() external view returns (string memory);
    function symbol() external view returns (string memory);
    function decimals() external view returns (uint8);
}

/// @dev Coin types for ZRC20. Zeta value should not be used.
enum CoinType {
    Zeta,
    Gas,
    ERC20
}

/**
 * @dev TestZRC20 is a test implementation of ZRC20 that extends the contract with new fields to test contract
 bytecode upgrade
 */
contract TestZRC20 is IZRC20, IZRC20Metadata, ZRC20Errors {
    /// @notice Fungible address is always the same, maintained at the protocol level
    address public constant FUNGIBLE_MODULE_ADDRESS = 0x735b14BB79463307AAcBED86DAf3322B1e6226aB;
    /// @notice Chain id.abi
    uint256 public immutable CHAIN_ID;
    /// @notice Coin type, checkout Interfaces.sol.
    CoinType public immutable COIN_TYPE;
    /// @notice System contract address.
    address public SYSTEM_CONTRACT_ADDRESS;
    /// @notice Gas limit.
    uint256 public GAS_LIMIT;
    /// @notice Protocol flat fee.
    uint256 public PROTOCOL_FLAT_FEE;

    mapping(address => uint256) private _balances;
    mapping(address => mapping(address => uint256)) private _allowances;
    uint256 private _totalSupply;
    string private _name;
    string private _symbol;
    uint8 private _decimals;

    /// @notice extend the contract with new fields to test contract bytecode upgrade
    uint256 public newField;
    string public newPublicField;

    function _msgSender() internal view virtual returns (address) {
        return msg.sender;
    }

    function _msgData() internal view virtual returns (bytes calldata) {
        return msg.data;
    }

    /**
     * @dev Only fungible module modifier.
     */
    modifier onlyFungible() {
        if (msg.sender != FUNGIBLE_MODULE_ADDRESS) revert CallerIsNotFungibleModule();
        _;
    }

    /**
     * @dev Constructor
     */
    constructor(
        uint256 chainid_,
        CoinType coinType_
    ) {
        CHAIN_ID = chainid_;
        COIN_TYPE = coinType_;
    }

    /**
     * @dev ZRC20 name
     * @return name as string
     */
    function name() public view virtual override returns (string memory) {
        return _name;
    }

    /**
     * @dev ZRC20 symbol.
     * @return symbol as string.
     */
    function symbol() public view virtual override returns (string memory) {
        return _symbol;
    }

    /**
     * @dev ZRC20 decimals.
     * @return returns uint8 decimals.
     */
    function decimals() public view virtual override returns (uint8) {
        return _decimals;
    }

    /**
     * @dev ZRC20 total supply.
     * @return returns uint256 total supply.
     */
    function totalSupply() public view virtual override returns (uint256) {
        return _totalSupply;
    }

    /**
     * @dev Returns ZRC20 balance of an account.
     * @param account, account address for which balance is requested.
     * @return uint256 account balance.
     */
    function balanceOf(address account) public view virtual override returns (uint256) {
        return _balances[account];
    }

    /**
     * @dev Returns ZRC20 balance of an account.
     * @param recipient, recipiuent address to which transfer is done.
     * @return true/false if transfer succeeded/failed.
     */
    function transfer(address recipient, uint256 amount) public virtual override returns (bool) {
        _transfer(_msgSender(), recipient, amount);
        return true;
    }

    /**
     * @dev Returns token allowance from owner to spender.
     * @param owner, owner address.
     * @return uint256 allowance.
     */
    function allowance(address owner, address spender) public view virtual override returns (uint256) {
        return _allowances[owner][spender];
    }

    /**
     * @dev Approves amount transferFrom for spender.
     * @param spender, spender address.
     * @param amount, amount to approve.
     * @return true/false if succeeded/failed.
     */
    function approve(address spender, uint256 amount) public virtual override returns (bool) {
        _approve(_msgSender(), spender, amount);
        return true;
    }

    /**
     * @dev Increases allowance by amount for spender.
     * @param spender, spender address.
     * @param amount, amount by which to increase allownace.
     * @return true/false if succeeded/failed.
     */
    function increaseAllowance(address spender, uint256 amount) external virtual returns (bool) {
        _allowances[spender][_msgSender()] += amount;
        return true;
    }

    /**
     * @dev Decreases allowance by amount for spender.
     * @param spender, spender address.
     * @param amount, amount by which to decrease allownace.
     * @return true/false if succeeded/failed.
     */
    function decreaseAllowance(address spender, uint256 amount) external virtual returns (bool) {
        if (_allowances[spender][_msgSender()] < amount) revert LowAllowance();
        _allowances[spender][_msgSender()] -= amount;
        return true;
    }

    /**
     * @dev Transfers tokens from sender to recipient.
     * @param sender, sender address.
     * @param recipient, recipient address.
     * @param amount, amount to transfer.
     * @return true/false if succeeded/failed.
     */
    function transferFrom(address sender, address recipient, uint256 amount) public virtual override returns (bool) {
        _transfer(sender, recipient, amount);

        uint256 currentAllowance = _allowances[sender][_msgSender()];
        if (currentAllowance < amount) revert LowAllowance();

        _approve(sender, _msgSender(), currentAllowance - amount);

        return true;
    }

    /**
     * @dev Burns an amount of tokens.
     * @param amount, amount to burn.
     * @return true/false if succeeded/failed.
     */
    function burn(uint256 amount) external returns (bool) {
        _burn(msg.sender, amount);
        return true;
    }

    function _transfer(address sender, address recipient, uint256 amount) internal virtual {
        if (sender == address(0)) revert ZeroAddress();
        if (recipient == address(0)) revert ZeroAddress();

        uint256 senderBalance = _balances[sender];
        if (senderBalance < amount) revert LowBalance();

        _balances[sender] = senderBalance - amount;
        _balances[recipient] += amount;

        emit Transfer(sender, recipient, amount);
    }

    function _mint(address account, uint256 amount) internal virtual {
        if (account == address(0)) revert ZeroAddress();

        _totalSupply += amount;
        _balances[account] += amount;
        emit Transfer(address(0), account, amount);
    }

    function _burn(address account, uint256 amount) internal virtual {
        if (account == address(0)) revert ZeroAddress();

        uint256 accountBalance = _balances[account];
        if (accountBalance < amount) revert LowBalance();

        _balances[account] = accountBalance - amount;
        _totalSupply -= amount;

        emit Transfer(account, address(0), amount);
    }

    function _approve(address owner, address spender, uint256 amount) internal virtual {
        if (owner == address(0)) revert ZeroAddress();
        if (spender == address(0)) revert ZeroAddress();

        _allowances[owner][spender] = amount;
        emit Approval(owner, spender, amount);
    }

    /**
     * @dev Deposits corresponding tokens from external chain, only callable by Fungible module.
     * @param to, recipient address.
     * @param amount, amount to deposit.
     * @return true/false if succeeded/failed.
     */
    function deposit(address to, uint256 amount) external override returns (bool) {
        if (msg.sender != FUNGIBLE_MODULE_ADDRESS && msg.sender != SYSTEM_CONTRACT_ADDRESS) revert InvalidSender();
        _mint(to, amount);
        emit Deposit(abi.encodePacked(FUNGIBLE_MODULE_ADDRESS), to, amount);
        return true;
    }

    /**
     * @dev Withdraws gas fees.
     * @return returns the ZRC20 address for gas on the same chain of this ZRC20, and calculates the gas fee for withdraw()
     */
    function withdrawGasFee() public view override returns (address, uint256) {
        address gasZRC20 = ISystem(SYSTEM_CONTRACT_ADDRESS).gasCoinZRC20ByChainId(CHAIN_ID);
        if (gasZRC20 == address(0)) {
            revert ZeroGasCoin();
        }
        uint256 gasPrice = ISystem(SYSTEM_CONTRACT_ADDRESS).gasPriceByChainId(CHAIN_ID);
        if (gasPrice == 0) {
            revert ZeroGasPrice();
        }
        uint256 gasFee = gasPrice * GAS_LIMIT + PROTOCOL_FLAT_FEE;
        return (gasZRC20, gasFee);
    }

    /**
     * @dev Withraws ZRC20 tokens to external chains, this function causes cctx module to send out outbound tx to the outbound chain
     * this contract should be given enough allowance of the gas ZRC20 to pay for outbound tx gas fee.
     * @param to, recipient address.
     * @param amount, amount to deposit.
     * @return true/false if succeeded/failed.
     */
    function withdraw(bytes memory to, uint256 amount) external override returns (bool) {
        (address gasZRC20, uint256 gasFee) = withdrawGasFee();
        if (!IZRC20(gasZRC20).transferFrom(msg.sender, FUNGIBLE_MODULE_ADDRESS, gasFee)) {
            revert GasFeeTransferFailed();
        }
        _burn(msg.sender, amount);
        emit Withdrawal(msg.sender, to, amount, gasFee, PROTOCOL_FLAT_FEE);
        return true;
    }

    /**
     * @dev Updates system contract address. Can only be updated by the fungible module.
     * @param addr, new system contract address.
     */
    function updateSystemContractAddress(address addr) external onlyFungible {
        SYSTEM_CONTRACT_ADDRESS = addr;
        emit UpdatedSystemContract(addr);
    }

    /**
     * @dev Updates gas limit. Can only be updated by the fungible module.
     * @param gasLimit, new gas limit.
     */
    function updateGasLimit(uint256 gasLimit) external onlyFungible {
        GAS_LIMIT = gasLimit;
        emit UpdatedGasLimit(gasLimit);
    }

    /**
     * @dev Updates protocol flat fee. Can only be updated by the fungible module.
     * @param protocolFlatFee, new protocol flat fee.
     */
    function updateProtocolFlatFee(uint256 protocolFlatFee) external onlyFungible {
        PROTOCOL_FLAT_FEE = protocolFlatFee;
        emit UpdatedProtocolFlatFee(protocolFlatFee);
    }

    /**
     * @dev Updates newField. Can only be updated by the fungible module.
     * @param newField_, new newField.
     */
    function updateNewField(uint256 newField_) external {
        newField = newField_;
    }
}
