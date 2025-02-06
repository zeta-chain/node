package sui

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"

	"golang.org/x/crypto/blake2b"
)

// test vector:
// see https://github.com/MystenLabs/sui/blob/199f06d25ce85f0270a1a5a0396156bb2b83122c/sdk/typescript/test/unit/cryptography/secp256k1-keypair.test.ts

var VALID_SECP256K1_SECRET_KEY = []byte{
	59, 148, 11, 85, 134, 130, 61, 253, 2, 174, 59, 70, 27, 180, 51, 107, 94, 203,
	174, 253, 102, 39, 170, 146, 46, 252, 4, 143, 236, 12, 136, 28,
}

var VALID_SECP256K1_PUBLIC_KEY = []byte{
	2, 29, 21, 35, 7, 198, 183, 43, 14, 208, 65, 139, 14, 112, 205, 128, 231, 245,
	41, 91, 141, 134, 245, 114, 45, 63, 82, 19, 251, 210, 57, 79, 54,
}

// it('create keypair from secret key', () => {
//    const secret_key = new Uint8Array(VALID_SECP256K1_SECRET_KEY);
//    const pub_key = new Uint8Array(VALID_SECP256K1_PUBLIC_KEY);
//    let pub_key_base64 = toB64(pub_key);
//    const keypair = Secp256k1Keypair.fromSecretKey(secret_key);
//    expect(keypair.getPublicKey().toBytes()).toEqual(new Uint8Array(pub_key));
//    expect(keypair.getPublicKey().toBase64()).toEqual(pub_key_base64);
//  });

func TestSignerSecp256k1FromSecretKey(t *testing.T) {
	// Create signer from secret key
	signer := NewSignerSecp256k1FromSecretKey(VALID_SECP256K1_SECRET_KEY)

	// Get public key bytes
	pubKey := signer.GetPublicKey()

	// Compare with expected public key
	if !bytes.Equal(pubKey, VALID_SECP256K1_PUBLIC_KEY) {
		t.Errorf("Public key mismatch\nexpected: %v\ngot: %v", VALID_SECP256K1_PUBLIC_KEY, pubKey)
	}
}

func Test2(t *testing.T) {
	// test case generated with sui keytool command
	pkbytes, err := base64.StdEncoding.DecodeString("AQLsqLZ8MUNEsBvUXILavCu4dSAtsyg/8Jz0gssssyT63w==")
	if err != nil {
		t.Fatal(err)
	}
	// len(pkbytes) == 34
	// byte 0: 0x01 (flag byte) byte 1-33: 33 compressed public key bytes
	expectedAddress := "0xce97523ae06726043aa5783c7b37964fa92c561fc65515a8d9bf13d3f9c9eae8"

	// Prepare the input for hashing: flag byte + public key bytes
	input := make([]byte, len(pkbytes))
	input[0] = FLAG_SECP256K1
	fmt.Printf("input[0]=%d\n", input[0])
	copy(input[1:], pkbytes[1:])

	// Create BLAKE2b hash
	hash, err := blake2b.New256(nil)
	if err != nil {
		t.Fatal(err)
	}

	// Write input to hash
	hash.Write(pkbytes)

	// Get the final hash
	addressBytes := hash.Sum(nil)

	// Convert to hex string with 0x prefix
	addressHex := "0x" + hex.EncodeToString(addressBytes)

	if addressHex != expectedAddress {
		t.Errorf("Address mismatch\nexpected: %s\ngot: %s", expectedAddress, addressHex)
	}
}
