module example::connected;

use sui::address::from_bytes;
use sui::coin::Coin;

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
    let receiver = decode_receiver(data);

    // transfer the coins to the provided address
    transfer::public_transfer(in_coins, receiver);

    // increment the called count
    cetus_config.called_count = cetus_config.called_count + 1;
}

fun decode_receiver(data: vector<u8>): address {
    from_bytes(data)
}