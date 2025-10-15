# P256 Precompile

The P256 precompile implements secp256r1 (also known as P-256 or prime256v1) elliptic curve signature verification as
defined in EIP-7212. This enables smart contracts to verify signatures from devices and systems that use the
secp256r1 curve, such as WebAuthn authenticators, secure enclaves, and many hardware security modules.

## Address

The precompile is available at the fixed address: `0x0000000000000000000000000000000000000100`

## Interface

The P256 precompile doesn't have a Solidity interface as it operates at a lower level.
It accepts raw input data and returns a verification result.

### Input Format

The precompile expects exactly 160 bytes of input data:

| Offset | Length | Description |
|--------|--------|-------------|
| 0 | 32 bytes | Message hash to verify |
| 32 | 32 bytes | Signature r component |
| 64 | 32 bytes | Signature s component |
| 96 | 32 bytes | Public key x coordinate |
| 128 | 32 bytes | Public key y coordinate |

### Output Format

- **Success**: Returns 32 bytes with value `0x0000000000000000000000000000000000000000000000000000000000000001` (1)
- **Failure**: Returns empty data or 32 zero bytes
- **Invalid input length**: Returns empty data

## Gas Cost

Fixed gas cost: **3,450 gas**

This cost is constant regardless of the input values, making gas consumption predictable.

## Implementation Details

### Signature Verification

The precompile performs ECDSA signature verification using the secp256r1 curve parameters:

- Verifies that the signature (r, s) is valid for the given message hash
- Checks that the public key (x, y) corresponds to the signature
- Uses constant-time operations to prevent timing attacks

### Security Features

1. **Input validation**: Strictly enforces 160-byte input length
2. **Cryptographic verification**: Uses proven secp256r1 implementation
3. **No state changes**: Pure function with no side effects
4. **Constant gas cost**: Prevents gas-based attacks

## Usage Example

### Direct Call from Solidity

```solidity
contract P256Verifier {
    // P256 precompile address
    address constant P256_PRECOMPILE = 0x0000000000000000000000000000000000000100;
    
    function verifySignature(
        bytes32 messageHash,
        bytes32 r,
        bytes32 s,
        bytes32 x,
        bytes32 y
    ) public view returns (bool) {
        // Prepare input data
        bytes memory input = abi.encodePacked(messageHash, r, s, x, y);
        
        // Call the precompile
        (bool success, bytes memory result) = P256_PRECOMPILE.staticcall(input);
        
        // Check if call was successful and returned expected result
        if (success && result.length == 32) {
            uint256 verification = abi.decode(result, (uint256));
            return verification == 1;
        }
        
        return false;
    }
}
```

### WebAuthn Integration Example

```solidity
contract WebAuthnWallet {
    address constant P256_PRECOMPILE = 0x0000000000000000000000000000000000000100;
    
    struct PublicKey {
        bytes32 x;
        bytes32 y;
    }
    
    mapping(address => PublicKey) public userKeys;
    
    function verifyWebAuthnSignature(
        address user,
        bytes32 challenge,
        bytes32 r,
        bytes32 s
    ) public view returns (bool) {
        PublicKey memory pubKey = userKeys[user];
        
        bytes memory input = abi.encodePacked(
            challenge,  // WebAuthn challenge hash
            r,          // Signature r
            s,          // Signature s
            pubKey.x,   // Public key x
            pubKey.y    // Public key y
        );
        
        (bool success, bytes memory result) = P256_PRECOMPILE.staticcall(input);
        
        return success && result.length == 32 && uint256(bytes32(result)) == 1;
    }
}
```

## Use Cases

1. **WebAuthn/Passkeys**: Verify signatures from browser-based authenticators
2. **Hardware Security Modules**: Integrate with HSMs that use secp256r1
3. **Mobile Secure Enclaves**: Verify signatures from iOS Secure Enclave or Android Keystore
4. **Cross-chain Bridges**: Verify signatures from chains that use secp256r1
5. **Enterprise Integration**: Many enterprise systems use P-256 for compliance

## Comparison with ecrecover

| Feature | P256 Precompile | ecrecover |
|---------|----------------|-----------|
| Curve | secp256r1 (P-256) | secp256k1 |
| Gas Cost | 3,450 | 3,000 |
| Input Length | 160 bytes | 128 bytes |
| Output | Verification result | Recovered address |
| Use Cases | WebAuthn, HSMs, Enterprise | Ethereum signatures |

## Security Considerations

1. **Message Hash**: Always hash the actual message before verification; never pass raw data
2. **Signature Malleability**: The implementation should handle signature malleability
3. **Public Key Validation**: The precompile validates that the public key is on the curve
4. **Side-channel Resistance**: Implementation uses constant-time operations

## Integration Notes

- The precompile is stateless and can be called multiple times
- No special initialization or setup required
- Compatible with all Solidity versions that support low-level calls
- Can be used in view/pure functions since it doesn't modify state

