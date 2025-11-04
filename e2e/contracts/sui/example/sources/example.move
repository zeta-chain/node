module example::connected;

use std::ascii;
use std::ascii::String;
use sui::address::from_bytes;
use sui::coin::Coin;
use gateway::gateway::{MessageContext, message_context_sender, message_context_target};

// === Errors ===

const EInvalidPayload: u64 = 1;

const EUnauthorizedSender: u64 = 2;

// ENonceMismatch is a fabricated nonce mismatch error code emitted from the on_call function
// zetaclient should be able to differentiate this error from real withdraw_impl nonce mismatch
const ENonceMismatch: u64 = 3;

const EInactiveMessageContext: u64 = 4;

const EPackageMismatch: u64 = 5;

// stub for shared objects
public struct GlobalConfig has key {
    id: UID,
    called_count: u64,
}

public struct Partner has key {
    id: UID,
}

public struct Clock has key {
    id: UID,
}

public struct Pool<phantom CoinA, phantom CoinB> has key {
    id: UID,
}

// share objects
fun init(ctx: &mut TxContext) {
    let global_config = GlobalConfig {
        id: object::new(ctx),
        called_count: 0,
    };
    let pool = Pool<sui::sui::SUI, example::token::TOKEN> {
        id: object::new(ctx),
    };
    let partner = Partner {
        id: object::new(ctx),
    };
    let clock = Clock {
        id: object::new(ctx),
    };

    transfer::share_object(global_config);
    transfer::share_object(pool);
    transfer::share_object(partner);
    transfer::share_object(clock);
}

public entry fun on_call<SOURCE_COIN>(
    message_context: &MessageContext,
    in_coins: Coin<SOURCE_COIN>,
    cetus_config: &mut GlobalConfig,
    // Note: this pool type is hardcoded as <SUI, TOKEN> and therefore causes type mismatch error in the
    // fungible token withdrawAndCall test, where the SOURCE_COIN type is FAKE_USDC instead of TOKEN.
    // Disabling the pool object for now is the easiest solution to allow the E2E tests to go through.
    // _pool: &mut Pool<SOURCE_COIN, TARGET_COIN>,
    _cetus_partner: &mut Partner,
    _clock: &Clock,
    data: vector<u8>,
    _ctx: &mut TxContext,
) {
    // check if the message is "revert" and revert with faked ENonceMismatch if so
    if (data == b"revert") {
        assert!(false, ENonceMismatch);
    };

    // decode the sender, target package, and receiver from the payload
    let (authenticated_sender, target_package, receiver) = decode_sender_target_and_receiver(data);

    // check if the sender is the authorized sender
    let actual_sender = message_context_sender(message_context);
    assert!(authenticated_sender == actual_sender, EUnauthorizedSender);

    // check if the target package is my own package
    // this prevents other package routing TSS calls to my package
    let actual_target = message_context_target(message_context);
    assert!(actual_target == target_package, EPackageMismatch);

    // transfer the coins to the provided address
    transfer::public_transfer(in_coins, receiver);

    // increment the called count
    cetus_config.called_count = cetus_config.called_count + 1;
}

// decode the sender, target package, and receiver from the payload data
fun decode_sender_target_and_receiver(data: vector<u8>): (String, address, address) {
    // [42-byte ZEVM sender] + [32-byte Sui target package] + [32-byte Sui receiver] = 106 bytes
    assert!(vector::length(&data) >= 106, EInvalidPayload);

    // extract ZEVM sender address (first 42 bytes)
    // this allows E2E test to feed a custom authenticated address
    let sender_bytes = extract_bytes(&data, 0, 42);
    let sender_str = ascii::string(sender_bytes);

    // extract target package address (bytes 42-74)
    let target_bytes = extract_bytes(&data, 42, 74);
    let target_package = from_bytes(target_bytes);

    // extract receiver address (bytes 74-106)
    let receiver_bytes = extract_bytes(&data, 74, 106);
    let receiver = from_bytes(receiver_bytes);
    
    (sender_str, target_package, receiver)
}

// helper function to extract a subslice of bytes from a vector
fun extract_bytes(data: &vector<u8>, start: u64, end: u64): vector<u8> {
    let mut result = vector::empty<u8>();
    let mut i = start;
    
    while (i < end) {
        vector::push_back(&mut result, *vector::borrow(data, i));
        i = i + 1;
    };
    
    result
}
