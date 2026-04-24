package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/bnb-chain/tss-lib/common"
	"github.com/bnb-chain/tss-lib/crypto"
	"github.com/bnb-chain/tss-lib/crypto/mta"
	"github.com/bnb-chain/tss-lib/crypto/paillier"
	"github.com/bnb-chain/tss-lib/ecdsa/keygen"
	"github.com/bnb-chain/tss-lib/tss"
)

// Simulates the Alpha-Rays / TSSHOCK attack (ZCNode-228).
// A malicious Bob supplies an oversized betaPrm (q^7) to the MtA proof.
// With the patched tss-lib, ProofBobWC.Verify must reject this proof.

func main() {
	fmt.Println("=== ZCNode-228 Verification: Alpha-Rays / TSSHOCK ===")
	fmt.Println()

	ec := tss.EC()
	q := ec.Params().N

	fmt.Println("[1] Generating Paillier keypair (2048-bit)...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	_, pkA, err := paillier.GenerateKeyPair(ctx, 2048)
	if err != nil {
		fmt.Printf("FAIL: Paillier keygen: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("[2] Loading test fixtures for NTilde, h1, h2...")
	NTildei, h1i, h2i, err := keygen.LoadNTildeH1H2FromTestFixture(0)
	if err != nil {
		fmt.Printf("FAIL: load fixture 0: %v\n", err)
		os.Exit(1)
	}
	NTildej, h1j, h2j, err := keygen.LoadNTildeH1H2FromTestFixture(1)
	if err != nil {
		fmt.Printf("FAIL: load fixture 1: %v\n", err)
		os.Exit(1)
	}

	a := common.GetRandomPositiveInt(q)
	b := common.GetRandomPositiveInt(q)

	fmt.Println("[3] Running AliceInit...")
	cA, pfA, err := mta.AliceInit(ec, pkA, a, NTildej, h1j, h2j)
	if err != nil {
		fmt.Printf("FAIL: AliceInit: %v\n", err)
		os.Exit(1)
	}

	// --- Attack simulation: malicious betaPrm = q^7 ---
	fmt.Println("[4] Constructing MALICIOUS proof (betaPrm = q^7)...")
	q7 := new(big.Int).Exp(q, big.NewInt(7), nil)

	cBetaPrmMal, cRandMal, err := pkA.EncryptAndReturnRandomness(q7)
	if err != nil {
		fmt.Printf("FAIL: encrypt malicious betaPrm: %v\n", err)
		os.Exit(1)
	}

	cB, err := pkA.HomoMult(b, cA)
	if err != nil {
		fmt.Printf("FAIL: HomoMult: %v\n", err)
		os.Exit(1)
	}
	cB, err = pkA.HomoAdd(cB, cBetaPrmMal)
	if err != nil {
		fmt.Printf("FAIL: HomoAdd: %v\n", err)
		os.Exit(1)
	}

	gBX, gBY := ec.ScalarBaseMult(b.Bytes())
	B, err := crypto.NewECPoint(ec, gBX, gBY)
	if err != nil {
		fmt.Printf("FAIL: NewECPoint: %v\n", err)
		os.Exit(1)
	}

	pfB, err := mta.ProveBobWC(ec, pkA, NTildei, h1i, h2i, cA, cB, b, q7, cRandMal, B)
	if err != nil {
		fmt.Printf("FAIL: ProveBobWC: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("     malicious T1 bit length: %d (q^7 = %d bits)\n", pfB.T1.BitLen(), q7.BitLen())

	// --- Verification: must reject ---
	fmt.Println("[5] Verifying malicious proof (expecting REJECT)...")
	ok := pfB.Verify(ec, pkA, NTildei, h1i, h2i, cA, cB, B)

	fmt.Println()
	if !ok {
		fmt.Println("PASS: Malicious proof REJECTED — Alpha-Rays attack is mitigated")
	} else {
		fmt.Println("FAIL: Malicious proof ACCEPTED — tss-lib is still vulnerable!")
		os.Exit(1)
	}

	// --- Sanity check: honest proof must still pass ---
	fmt.Println()
	fmt.Println("[6] Verifying honest proof (expecting ACCEPT)...")
	_, cBHonest, _, pfBHonest, err := mta.BobMidWC(ec, pkA, pfA, b, cA, NTildei, h1i, h2i, NTildej, h1j, h2j, B)
	if err != nil {
		fmt.Printf("FAIL: BobMidWC (honest): %v\n", err)
		os.Exit(1)
	}

	okHonest := pfBHonest.Verify(ec, pkA, NTildei, h1i, h2i, cA, cBHonest, B)
	if okHonest {
		fmt.Println("PASS: Honest proof ACCEPTED — no regression")
	} else {
		fmt.Println("FAIL: Honest proof REJECTED — regression in honest path!")
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("=== All checks passed ===")
}
