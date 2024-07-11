package main

import (
	"context"
	_ "embed"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
)

const (
	pythProgramDevnet = "gSbePebfvPy7tRqimPoVecS2UsBvYv46ynrzWocc92s" // this program has many many txs
)

//go:embed gateway.json
var GatewayIDLJSON []byte

func main() {
	// devnet RPC
	client := rpc.New("https://solana-devnet.g.allthatnode.com/archive/json_rpc/842c667c947e42e2a9995ac2ec75026d")

	limit := 10
	out, err := client.GetSignaturesForAddressWithOpts(
		context.TODO(),
		solana.MustPublicKeyFromBase58(pythProgramDevnet),
		&rpc.GetSignaturesForAddressOpts{
			Limit: &limit,
			Before: solana.MustSignatureFromBase58(
				"5pLBywq74Nc6jYrWUqn9KjnYXHbQEY2UPkhWefZF5u4NYaUvEwz1Cirqaym9wDeHNAjiQwuLBfrdhXo8uFQA45jL",
			),
			Until: solana.MustSignatureFromBase58(
				"2coX9CckSmJWeHVqJNANeD7m4J7pctpSomxMon3h36droxCVB3JDbLyWQKMjnf85ntuFGxMLySykEMaRd5MDw35e",
			),
		},
	)

	if err != nil {
		panic(err)
	}
	fmt.Printf("len(out) = %d\n", len(out))
	//spew.Dump(out)
	for _, sig := range out {
		fmt.Printf("%s %d %v\n", sig.Signature, sig.Slot, sig.Err == nil)
	}

	{
		bn, _ := client.GetFirstAvailableBlock(context.TODO())
		fmt.Printf("first available bn = %d\n", bn)
		cutoffTimestamp, _ := client.GetBlockTime(context.TODO(), bn)
		fmt.Printf("cutoffTimestamp = %s\n", cutoffTimestamp.Time())
		block, _ := client.GetBlock(context.TODO(), bn)
		//spew.Dump(block)
		fmt.Printf("block time %s, block height %d\n", block.BlockTime.Time(), *block.BlockHeight)
		fmt.Printf("block #%d\n", len(block.Transactions))
		//first_tx := block.Signatures[0]
		//spew.Dump(first_tx)
	}

	{
		// Parsing a Deposit Instruction
		// devnet tx: deposit with memo
		// https://solana.fm/tx/51746triQeve21zP1bcVEPvvsoXt94B57TU5exBvoy938bhGCfzBtsvKJbLpS1zRc2dmb3S3HBHnhTfbtKCBpmqg
		const depositTx = "51746triQeve21zP1bcVEPvvsoXt94B57TU5exBvoy938bhGCfzBtsvKJbLpS1zRc2dmb3S3HBHnhTfbtKCBpmqg"

		tx, err := client.GetTransaction(
			context.TODO(),
			solana.MustSignatureFromBase58(depositTx),
			&rpc.GetTransactionOpts{})
		if err != nil {
			log.Fatalf("Error getting transaction: %v", err)
		}
		fmt.Printf("tx status: %v", tx.Meta.Err == nil)
		//spew.Dump(tx)
		type DepositInstructionParams struct {
			Discriminator [8]byte
			Amount        uint64
			Memo          []byte
		}
		//hexString := "f223c68952e1f2b6390500000000000014000000dead000000000000000042069420694206942069"
		// Decode hex string to byte slice
		//data, _ := hex.DecodeString(hexString)
		transaction, _ := tx.Transaction.GetTransaction()
		instruction := transaction.Message.Instructions[0]
		data := instruction.Data
		pk, _ := transaction.Message.Program(instruction.ProgramIDIndex)
		fmt.Printf("Program ID: %s\n", pk)
		var inst DepositInstructionParams
		err = borsh.Deserialize(&inst, data)
		if err != nil {
			log.Fatalf("Error deserializing: %v", err)
		}
		fmt.Printf("Discriminator: %016x\n", inst.Discriminator)
		fmt.Printf("U64 Parameter: %d\n", inst.Amount)
		fmt.Printf("Vec<u8> (%d): %x\n", len(inst.Memo), inst.Memo)
	}

	{
		var idl IDL
		err := json.Unmarshal(GatewayIDLJSON, &idl)
		if err != nil {
			panic(err)
		}
		//spew.Dump(idl)
	}

	{
		// explore failed transaction
		//https://explorer.solana.com/tx/2LbBdmCkuVyQhHAvsZhZ1HLdH12jQbHY7brwH6xUBsZKKPuV8fomyz1Qh9CaCZSqo8FNefaR8ir7ngo7H3H2VfWv
		txSig := solana.MustSignatureFromBase58(
			"2LbBdmCkuVyQhHAvsZhZ1HLdH12jQbHY7brwH6xUBsZKKPuV8fomyz1Qh9CaCZSqo8FNefaR8ir7ngo7H3H2VfWv",
		)
		client2 := rpc.New("https://solana-mainnet.g.allthatnode.com/archive/json_rpc/842c667c947e42e2a9995ac2ec75026d")
		tx, err := client2.GetTransaction(
			context.TODO(),
			txSig,
			&rpc.GetTransactionOpts{})
		if err != nil {
			log.Fatalf("Error getting transaction: %v", err)
		}
		fmt.Printf("tx successful?: %v\n", tx.Meta.Err == nil)
		spew.Dump(tx)
	}

	pk := os.Getenv("SOLANA_WALLET_PK")
	if pk == "" {
		log.Fatal("SOLANA_WALLET_PK must be set (base58 encoded private key)")
	}

	privkey, err := solana.PrivateKeyFromBase58(pk)
	if err != nil {
		log.Fatalf("Error getting private key: %v", err)
	}
	fmt.Println("account public key:", privkey.PublicKey())

	ethPk := os.Getenv("ETH_WALLET_PK")
	if ethPk == "" {
		log.Fatal("ETH_WALLET_PK must be set (hex encoded private key)")
	}
	privkeyBytes, err := hex.DecodeString(ethPk)
	if err != nil {
		log.Fatalf("Error decoding hex private key: %v", err)
	}
	ethPrivkey, err := crypto.ToECDSA(privkeyBytes)
	if err != nil {
		log.Fatalf("Error converting to ECDSA: %v", err)
	}

	{
		// build & bcast a Depsosit tx
		bal, err := client.GetBalance(context.TODO(), privkey.PublicKey(), rpc.CommitmentFinalized)
		if err != nil {
			log.Fatalf("Error getting balance: %v", err)
		}
		fmt.Println("account balance in SOL ", float64(bal.Value)/1e9)

		// building the transaction
		recent, err := client.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
		if err != nil {
			panic(err)
		}
		fmt.Println("recent blockhash:", recent.Value.Blockhash)

		programID := solana.MustPublicKeyFromBase58("4Nt8tsYWQj3qC1TbunmmmDbzRXE4UQuzcGcqqgwy9bvX")
		seed := []byte("meta")
		pdaComputed, bump, err := solana.FindProgramAddress([][]byte{seed}, programID)
		if err != nil {
			panic(err)
		}
		fmt.Printf("computed pda: %s, bump %d\n", pdaComputed, bump)

		//pdaAccount := solana.MustPublicKeyFromBase58("4hA43LCh2Utef8EwCyWwYmWBoSeNq6RS2HdoLkWGm5z5")
		var inst solana.GenericInstruction
		accountSlice := []*solana.AccountMeta{}
		accountSlice = append(accountSlice, solana.Meta(privkey.PublicKey()).WRITE().SIGNER())
		accountSlice = append(accountSlice, solana.Meta(pdaComputed).WRITE())
		accountSlice = append(accountSlice, solana.Meta(solana.SystemProgramID))
		accountSlice = append(accountSlice, solana.Meta(programID))
		inst.ProgID = programID
		inst.AccountValues = accountSlice

		type DepositInstructionParams struct {
			Discriminator [8]byte
			Amount        uint64
			Memo          []byte
		}

		inst.DataBytes, err = borsh.Serialize(DepositInstructionParams{
			Discriminator: [8]byte{0xf2, 0x23, 0xc6, 0x89, 0x52, 0xe1, 0xf2, 0xb6},
			Amount:        1338,
			Memo:          []byte("hello this is a good memo for you to enjoy"),
		})
		//inst.DataBytes, err = hex.DecodeString("f223c68952e1f2b6390500000000000014000000dead000000000000000042069420694206942069")
		if err != nil {
			panic(err)
		}

		tx, err := solana.NewTransaction(
			[]solana.Instruction{&inst},
			recent.Value.Blockhash,
			solana.TransactionPayer(privkey.PublicKey()),
		)
		if err != nil {
			panic(err)
		}
		_, err = tx.Sign(
			func(key solana.PublicKey) *solana.PrivateKey {
				if privkey.PublicKey().Equals(key) {
					return &privkey
				}
				return nil
			},
		)
		if err != nil {
			panic(fmt.Errorf("unable to sign transaction: %w", err))
		}

		spew.Dump(tx)
		//wsClient, err := ws.Connect(context.Background(), rpc.DevNet_WS)
		//if err != nil {
		//	panic(err)
		//}
		//sig, err := confirm.SendAndConfirmTransaction(
		//	context.TODO(),
		//	client,
		//	wsClient,
		//	tx,
		//)
		// tx: 33cVywTwufSy5NsNSnJS87wmkPwVAr9iiJqxAhhny9pazxWpiH6L24c6ruVnSjctcGasyt2ngnrtx3TqK6KU6x6j

		//sig, err := client.SendTransactionWithOpts(
		//	context.TODO(),
		//	tx,
		//	rpc.TransactionOpts{},
		//)
		// broadcast success! see
		// https://solana.fm/tx/43hXUywVouKeG5V98mjPysPWG9eKyKo6XDVHuoQs5YP1gJfa5z2UtU6hjJGgscrWzmYbhbqNW2hykvV6HYfBXATD

		//if err != nil {
		//	panic(err)
		//}
		//spew.Dump(sig)
	}

	{
		fmt.Printf("Build and broadcast a withdraw tx\n")
		type WithdrawInstructionParams struct {
			Discriminator [8]byte
			Amount        uint64
			Signature     [64]byte
			RecoveryID    uint8
			MessageHash   [32]byte
			Nonce         uint64
		}
		// fetch PDA account
		programID := solana.MustPublicKeyFromBase58("4Nt8tsYWQj3qC1TbunmmmDbzRXE4UQuzcGcqqgwy9bvX")
		seed := []byte("meta")
		pdaComputed, bump, err := solana.FindProgramAddress([][]byte{seed}, programID)
		if err != nil {
			panic(err)
		}
		fmt.Printf("computed pda: %s, bump %d\n", pdaComputed, bump)
		type PdaInfo struct {
			Discriminator [8]byte
			Nonce         uint64
			TssAddress    [20]byte
			Authority     [32]byte
		}
		pdaInfo, err := client.GetAccountInfo(context.TODO(), pdaComputed)
		if err != nil {
			panic(err)
		}

		// deserialize PDA account
		var pda PdaInfo
		err = borsh.Deserialize(&pda, pdaInfo.Bytes())
		if err != nil {
			panic(err)
		}

		//spew.Dump(pda)
		// building the transaction
		recent, err := client.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
		if err != nil {
			panic(err)
		}
		fmt.Println("recent blockhash:", recent.Value.Blockhash)
		var inst solana.GenericInstruction

		pdaBalance, err := client.GetBalance(context.TODO(), pdaComputed, rpc.CommitmentFinalized)
		if err != nil {
			panic(err)
		}
		fmt.Printf("PDA balance in SOL %f\n", float64(pdaBalance.Value)/1e9)
		var message []byte

		amount := uint64(2_337_000)
		to := privkey.PublicKey()
		bytes := make([]byte, 8)
		nonce := pda.Nonce
		binary.BigEndian.PutUint64(bytes, nonce)
		message = append(message, bytes...)
		binary.BigEndian.PutUint64(bytes, amount)
		message = append(message, bytes...)
		message = append(message, to.Bytes()...)
		messageHash := crypto.Keccak256Hash(message)
		// this sig will be 65 bytes; R || S || V, where V is 0 or 1
		signature, err := crypto.Sign(messageHash.Bytes(), ethPrivkey)
		if err != nil {
			panic(err)
		}
		var sig [64]byte
		copy(sig[:], signature[:64])
		inst.DataBytes, err = borsh.Serialize(WithdrawInstructionParams{
			Discriminator: [8]byte{183, 18, 70, 156, 148, 109, 161, 34},
			Amount:        amount,
			Signature:     sig,
			RecoveryID:    signature[64],
			MessageHash:   messageHash,
			Nonce:         nonce,
		})
		if err != nil {
			panic(err)
		}

		var accountSlice []*solana.AccountMeta
		accountSlice = append(accountSlice, solana.Meta(privkey.PublicKey()).WRITE().SIGNER())
		accountSlice = append(accountSlice, solana.Meta(pdaComputed).WRITE())
		accountSlice = append(accountSlice, solana.Meta(to).WRITE())
		accountSlice = append(accountSlice, solana.Meta(programID))
		inst.ProgID = programID
		inst.AccountValues = accountSlice
		tx, err := solana.NewTransaction(
			[]solana.Instruction{&inst},
			recent.Value.Blockhash,
			solana.TransactionPayer(privkey.PublicKey()),
		)
		if err != nil {
			panic(err)
		}
		_, err = tx.Sign(
			func(key solana.PublicKey) *solana.PrivateKey {
				if privkey.PublicKey().Equals(key) {
					return &privkey
				}
				return nil
			},
		)
		if err != nil {
			panic(fmt.Errorf("unable to sign transaction: %w", err))
		}

		spew.Dump(tx)
		txsig, err := client.SendTransactionWithOpts(
			context.TODO(),
			tx,
			rpc.TransactionOpts{},
		)
		//broadcast success! see
		if err != nil {
			panic(err)
		}
		spew.Dump(txsig)
	}
}
