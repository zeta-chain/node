module gateway::gateway;

use gateway::evm;
use std::ascii::String;
use std::type_name::{get, into_string};
use sui::bag::{Self, Bag};
use sui::balance::{Self, Balance};
use sui::coin::{Self, Coin};
use sui::event;
use sui::sui::SUI;

// === Errors ===

const EAlreadyWhitelisted: u64 = 0;
const EInvalidReceiverAddress: u64 = 1;
const ENotWhitelisted: u64 = 2;
const ENonceMismatch: u64 = 3;
const EPayloadTooLong: u64 = 4;
const EInactiveWithdrawCap: u64 = 5;
const EInactiveWhitelistCap: u64 = 6;
const EDepositPaused: u64 = 7;

const PayloadMaxLength: u64 = 1024;

// === Structs ===

// Vault stores the balance of a specific coin type
public struct Vault<phantom T> has store {
    balance: Balance<T>,
    whitelisted: bool,
}

// Gateway stores the vaults and the nonce for withdrawals
public struct Gateway has key {
    id: UID,
    vaults: Bag,
    nonce: u64,
    active_withdraw_cap: ID,
    active_whitelist_cap: ID,
    deposit_paused: bool,
}

// WithdrawCap is a capability object that allows the caller to withdraw tokens from the gateway
public struct WithdrawCap has key, store {
    id: UID,
}

// WhitelistCap is a capability object that allows the caller to whitelist a new vault
public struct WhitelistCap has key, store {
    id: UID,
}

// AdminCap is a capability object that allows to issue new capabilities
public struct AdminCap has key, store {
    id: UID,
}

// === Events ===

// DepositEvent is emitted when a user deposits tokens into the gateway
public struct DepositEvent has copy, drop {
    coin_type: String,
    amount: u64,
    sender: address,
    receiver: String, // 0x hex address
}

// DepositAndCallEvent is emitted when a user deposits tokens into the gateway with a call
public struct DepositAndCallEvent has copy, drop {
    coin_type: String,
    amount: u64,
    sender: address,
    receiver: String, // 0x hex address
    payload: vector<u8>,
}

public struct WithdrawEvent has copy, drop {
    coin_type: String,
    amount: u64,
    sender: address,
    receiver: address,
    nonce: u64,
}

public struct NonceIncreaseEvent has copy, drop {
    sender: address,
    nonce: u64,
}

// === Initialization ===

fun init(ctx: &mut TxContext) {
    // to withdraw tokens from the gateway, the caller must have the WithdrawCap
    let withdraw_cap = WithdrawCap {
        id: object::new(ctx),
    };

    // to whitelist a new vault, the caller must have the WhitelistCap
    let whitelist_cap = WhitelistCap {
        id: object::new(ctx),
    };

    // to whitelist a new vault, the caller must have the AdminCap
    let admin_cap = AdminCap {
        id: object::new(ctx),
    };

    // create and share the gateway object
    let mut gateway = Gateway {
        id: object::new(ctx),
        vaults: bag::new(ctx),
        nonce: 0,
        active_withdraw_cap: object::id(&withdraw_cap),
        active_whitelist_cap: object::id(&whitelist_cap),
        deposit_paused: false,
    };

    // whitelist SUI by default
    whitelist_impl<SUI>(&mut gateway, &whitelist_cap);

    transfer::transfer(withdraw_cap, tx_context::sender(ctx));
    transfer::transfer(whitelist_cap, tx_context::sender(ctx));
    transfer::transfer(admin_cap, tx_context::sender(ctx));
    transfer::share_object(gateway);
}

// === Entrypoints ===

// increase_nonce increases the nonce of the gateway
// it is used when a failed outbound needs to be reported to ZetaChain
// it is sent by the tss and therefore requires the withdraw cap
entry fun increase_nonce(gateway: &mut Gateway, nonce: u64, cap: &WithdrawCap, ctx: &TxContext) {
    assert!(gateway.active_withdraw_cap == object::id(cap), EInactiveWithdrawCap);
    assert!(nonce == gateway.nonce, ENonceMismatch);
    gateway.nonce = nonce + 1;

    // Emit event
    event::emit(NonceIncreaseEvent {
        sender: tx_context::sender(ctx),
        nonce: gateway.nonce,
    });
}

// withdraw allows the TSS to withdraw tokens from the gateway
entry fun withdraw<T>(
    gateway: &mut Gateway,
    amount: u64,
    nonce: u64,
    receiver: address,
    gas_budget: u64,
    cap: &WithdrawCap,
    ctx: &mut TxContext,
) {
    let (coins, coins_gas_budget) = withdraw_impl<T>(gateway, amount, nonce, gas_budget, cap, ctx);

    transfer::public_transfer(coins, receiver);
    transfer::public_transfer(coins_gas_budget, tx_context::sender(ctx));

    // Emit event
    event::emit(WithdrawEvent {
        coin_type: coin_name<T>(),
        amount: amount,
        sender: tx_context::sender(ctx),
        receiver: receiver,
        nonce: nonce,
    });
}

// whitelist whitelists a new coin by creating a new vault for the coin type
entry fun whitelist<T>(gateway: &mut Gateway, cap: &WhitelistCap) {
    whitelist_impl<T>(gateway, cap)
}

// unwhitelist unwhitelists a coin by setting the whitelisted flag to false
entry fun unwhitelist<T>(gateway: &mut Gateway, cap: &AdminCap) {
    unwhitelist_impl<T>(gateway, cap)
}

// issue_withdraw_and_whitelist_cap issues a new WithdrawCap and WhitelistCap and revokes the old ones
entry fun issue_withdraw_and_whitelist_cap(
    gateway: &mut Gateway,
    _cap: &AdminCap,
    ctx: &mut TxContext,
) {
    let (withdraw_cap, whitelist_cap) = issue_withdraw_and_whitelist_cap_impl(gateway, _cap, ctx);
    transfer::transfer(withdraw_cap, tx_context::sender(ctx));
    transfer::transfer(whitelist_cap, tx_context::sender(ctx));
}

// pause pauses the deposit functionality
entry fun pause(gateway: &mut Gateway, cap: &AdminCap) {
    pause_impl(gateway, cap)
}

// unpause unpauses the deposit functionality
entry fun unpause(gateway: &mut Gateway, cap: &AdminCap) {
    unpause_impl(gateway, cap)
}

// reset_nonce resets the nonce to a new value in case it gets corrupted
entry fun reset_nonce(gateway: &mut Gateway, nonce: u64, _cap: &AdminCap) {
    gateway.nonce = nonce;
}

// === Deposit Functions ===

// deposit allows the user to deposit tokens into the gateway
public entry fun deposit<T>(
    gateway: &mut Gateway,
    coins: Coin<T>,
    receiver: String,
    ctx: &mut TxContext,
) {
    let amount = coins.value();
    let coin_name = coin_name<T>();

    check_receiver_and_deposit_to_vault(gateway, coins, receiver);

    // Emit deposit event
    event::emit(DepositEvent {
        coin_type: coin_name,
        amount: amount,
        sender: tx_context::sender(ctx),
        receiver: receiver,
    });
}

// deposit_and_call allows the user to deposit tokens into the gateway and call a contract
public entry fun deposit_and_call<T>(
    gateway: &mut Gateway,
    coins: Coin<T>,
    receiver: String,
    payload: vector<u8>,
    ctx: &mut TxContext,
) {
    assert!(payload.length() <= PayloadMaxLength, EPayloadTooLong);

    let amount = coins.value();
    let coin_name = coin_name<T>();

    check_receiver_and_deposit_to_vault(gateway, coins, receiver);

    // Emit deposit event
    event::emit(DepositAndCallEvent {
        coin_type: coin_name,
        amount: amount,
        sender: tx_context::sender(ctx),
        receiver: receiver,
        payload: payload,
    });
}

// check_receiver_and_deposit_to_vault is a helper function that checks the receiver address and deposits the coin
fun check_receiver_and_deposit_to_vault<T>(
    gateway: &mut Gateway,
    coins: Coin<T>,
    receiver: String,
) {
    assert!(evm::is_valid_evm_address(receiver), EInvalidReceiverAddress);
    assert!(is_whitelisted<T>(gateway), ENotWhitelisted);
    assert!(!gateway.deposit_paused, EDepositPaused);

    // Deposit the coin into the vault
    let coin_name = coin_name<T>();
    let vault = bag::borrow_mut<String, Vault<T>>(&mut gateway.vaults, coin_name);
    balance::join(&mut vault.balance, coins.into_balance());
}

// === Withdraw Functions ===

public fun withdraw_impl<T>(
    gateway: &mut Gateway,
    amount: u64,
    nonce: u64,
    gas_budget: u64,
    cap: &WithdrawCap,
    ctx: &mut TxContext,
): (Coin<T>, Coin<sui::sui::SUI>) {
    assert!(gateway.active_withdraw_cap == object::id(cap), EInactiveWithdrawCap);
    assert!(is_whitelisted<T>(gateway), ENotWhitelisted);
    assert!(nonce == gateway.nonce, ENonceMismatch); // prevent replay
    gateway.nonce = nonce + 1;

    // Withdraw the coin from the vault
    let coin_name = coin_name<T>();
    let vault = bag::borrow_mut<String, Vault<T>>(&mut gateway.vaults, coin_name);
    let coins_out = coin::take(&mut vault.balance, amount, ctx);

    // Withdraw SUI to cover the gas budget
    let sui_vault = bag::borrow_mut<String, Vault<sui::sui::SUI>>(
        &mut gateway.vaults,
        coin_name<sui::sui::SUI>(),
    );
    let coins_gas_budget = coin::take(&mut sui_vault.balance, gas_budget, ctx);

    (coins_out, coins_gas_budget)
}

// === Admin Functions ===

public fun whitelist_impl<T>(gateway: &mut Gateway, cap: &WhitelistCap) {
    assert!(gateway.active_whitelist_cap == object::id(cap), EInactiveWhitelistCap);
    assert!(is_whitelisted<T>(gateway) == false, EAlreadyWhitelisted);

    // if the vault already exists, set it to whitelisted, otherwise create a new vault
    if (bag::contains_with_type<String, Vault<T>>(&gateway.vaults, coin_name<T>())) {
        let vault = bag::borrow_mut<String, Vault<T>>(&mut gateway.vaults, coin_name<T>());
        vault.whitelisted = true;
    } else {
        let vault_name = coin_name<T>();
        let vault = Vault<T> {
            balance: balance::zero<T>(),
            whitelisted: true,
        };
        bag::add(&mut gateway.vaults, vault_name, vault);
    }
}

public fun unwhitelist_impl<T>(gateway: &mut Gateway, _cap: &AdminCap) {
    assert!(is_whitelisted<T>(gateway), ENotWhitelisted);
    let vault = bag::borrow_mut<String, Vault<T>>(&mut gateway.vaults, coin_name<T>());
    vault.whitelisted = false;
}

public fun issue_withdraw_and_whitelist_cap_impl(
    gateway: &mut Gateway,
    _cap: &AdminCap,
    ctx: &mut TxContext,
): (WithdrawCap, WhitelistCap) {
    let withdraw_cap = WithdrawCap {
        id: object::new(ctx),
    };
    let whitelist_cap = WhitelistCap {
        id: object::new(ctx),
    };
    gateway.active_withdraw_cap = object::id(&withdraw_cap);
    gateway.active_whitelist_cap = object::id(&whitelist_cap);
    (withdraw_cap, whitelist_cap)
}

public fun pause_impl(gateway: &mut Gateway, _cap: &AdminCap) {
    gateway.deposit_paused = true;
}

public fun unpause_impl(gateway: &mut Gateway, _cap: &AdminCap) {
    gateway.deposit_paused = false;
}

// === View Functions ===

public fun nonce(gateway: &Gateway): u64 {
    gateway.nonce
}

public fun active_withdraw_cap(gateway: &Gateway): ID {
    gateway.active_withdraw_cap
}

public fun active_whitelist_cap(gateway: &Gateway): ID {
    gateway.active_whitelist_cap
}

public fun vault_balance<T>(gateway: &Gateway): u64 {
    if (!is_whitelisted<T>(gateway)) {
        return 0
    };
    let coin_name = coin_name<T>();
    let vault = bag::borrow<String, Vault<T>>(&gateway.vaults, coin_name);
    balance::value(&vault.balance)
}

public fun is_paused(gateway: &Gateway): bool {
    gateway.deposit_paused
}

// is_whitelisted returns true if a given coin type is whitelisted
public fun is_whitelisted<T>(gateway: &Gateway): bool {
    let vault_name = coin_name<T>();
    if (!bag::contains_with_type<String, Vault<T>>(&gateway.vaults, vault_name)) {
        return false
    };
    let vault = bag::borrow<String, Vault<T>>(&gateway.vaults, vault_name);
    vault.whitelisted
}

// === Helpers ===

// coin_name returns the name of the coin type to index the vault
fun coin_name<T>(): String {
    into_string(get<T>())
}

// === Test Helpers ===

#[test_only]
public fun init_for_testing(ctx: &mut TxContext) {
    init(ctx)
}

#[test_only]
public fun create_test_withdraw_cap(ctx: &mut TxContext): WithdrawCap {
    WithdrawCap {
        id: object::new(ctx),
    }
}

#[test_only]
public fun create_test_whitelist_cap(ctx: &mut TxContext): WhitelistCap {
    WhitelistCap {
        id: object::new(ctx),
    }
}
