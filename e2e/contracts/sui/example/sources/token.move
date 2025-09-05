module example::token;

use sui::coin::{Self, TreasuryCap};

public struct TOKEN has drop {}

fun init(witness: TOKEN, ctx: &mut TxContext) {
    let (treasury, metadata) = coin::create_currency(
        witness,
        6,
        b"TOKEN",
        b"",
        b"",
        option::none(),
        ctx,
    );
    transfer::public_freeze_object(metadata);
    transfer::public_transfer(treasury, ctx.sender())
}

public entry fun mint(
    treasury_cap: &mut TreasuryCap<TOKEN>,
    amount: u64,
    recipient: address,
    ctx: &mut TxContext,
) {
    let coin = coin::mint(treasury_cap, amount, ctx);
    transfer::public_transfer(coin, recipient)
}