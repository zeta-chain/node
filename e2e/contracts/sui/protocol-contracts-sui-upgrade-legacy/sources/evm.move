module gateway::evm;

use std::ascii::{String, into_bytes};

/// Check if a given string is a valid Ethereum address.
public fun is_valid_evm_address(addr: String): bool {
    if (addr.length() != 42) {
        return false
    };

    let mut addrBytes = addr.into_bytes();

    // check prefix 0x, 0=48, x=120
    if (addrBytes[0] != 48 || addrBytes[1] != 120) {
        return false
    };

    // remove 0x prefix
    addrBytes.remove(0);
    addrBytes.remove(0);

    // check if remaining characters are hex (0-9, a-f, A-F)
    is_hex_vec(addrBytes)
}

/// Check that vector contains only hex chars (0-9, a-f, A-F).
fun is_hex_vec(input: vector<u8>): bool {
    let mut i = 0;
    let len = input.length();

    while (i < len) {
        let c = input[i];

        let is_hex = (c >= 48 && c <= 57) ||  // '0' to '9'
                     (c >= 97 && c <= 102) || // 'a' to 'f'
                     (c >= 65 && c <= 70);    // 'A' to 'F'

        if (!is_hex) {
            return false
        };

        i = i + 1;
    };

    true
}
