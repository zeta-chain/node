package runner

import (
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/constant"
	solanacontracts "github.com/zeta-chain/node/pkg/contracts/solana"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// SetupSolanaAccount imports the deployer's private key
func (r *E2ERunner) SetupSolanaAccount() {
	privateKey, err := solana.PrivateKeyFromBase58(r.Account.SolanaPrivateKey.String())
	require.NoError(r, err)
	r.SolanaDeployerAddress = privateKey.PublicKey()

	r.Logger.Info("SolanaDeployerAddress: %s", r.SolanaDeployerAddress)
}

// SetupSolana sets Solana contracts and params
func (r *E2ERunner) SetupSolana(gatewayID, deployerPrivateKey, splAccountPrivateKey string) {
	r.Logger.Print("⚙️ initializing gateway program on Solana")

	// set Solana contracts
	r.GatewayProgram = solana.MustPublicKeyFromBase58(gatewayID)

	// get deployer account balance
	privkey, err := solana.PrivateKeyFromBase58(deployerPrivateKey)
	require.NoError(r, err)
	bal, err := r.SolanaClient.GetBalance(r.Ctx, privkey.PublicKey(), rpc.CommitmentConfirmed)
	require.NoError(r, err)
	r.Logger.Info("deployer address: %s, balance: %f SOL", privkey.PublicKey().String(), float64(bal.Value)/1e9)

	// compute the gateway PDA address
	pdaComputed := r.ComputePdaAddress()

	// create 'initialize' instruction
	var inst solana.GenericInstruction
	accountSlice := []*solana.AccountMeta{}
	accountSlice = append(accountSlice, solana.Meta(privkey.PublicKey()).WRITE().SIGNER())
	accountSlice = append(accountSlice, solana.Meta(pdaComputed).WRITE())
	accountSlice = append(accountSlice, solana.Meta(solana.SystemProgramID))
	inst.ProgID = r.GatewayProgram
	inst.AccountValues = accountSlice

	inst.DataBytes, err = borsh.Serialize(solanacontracts.InitializeParams{
		Discriminator: solanacontracts.DiscriminatorInitialize,
		TssAddress:    r.TSSAddress,
		// #nosec G115 chain id always positive
		ChainID: uint64(chains.SolanaLocalnet.ChainId),
	})
	require.NoError(r, err)

	// create and sign the transaction
	signedTx := r.CreateSignedTransaction([]solana.Instruction{&inst}, privkey, []solana.PrivateKey{})

	// broadcast the transaction and wait for finalization
	_, out := r.BroadcastTxSync(signedTx)
	r.Logger.Info("initialize gateway logs: %v", out.Meta.LogMessages)

	// initialize connected program, use hardcoded program ID if not set
	if (r.ConnectedProgram == solana.PublicKey{}) {
		r.ConnectedProgram = ConnectedProgramID
	}
	connectedPda, err := solanacontracts.ComputeConnectedPdaAddress(r.ConnectedProgram)
	require.NoError(r, err)

	var instConnected solana.GenericInstruction
	accountSliceConnected := []*solana.AccountMeta{}
	accountSliceConnected = append(accountSliceConnected, solana.Meta(privkey.PublicKey()).WRITE().SIGNER())
	accountSliceConnected = append(accountSliceConnected, solana.Meta(connectedPda).WRITE())
	accountSliceConnected = append(accountSliceConnected, solana.Meta(solana.SystemProgramID))
	instConnected.ProgID = r.ConnectedProgram
	instConnected.AccountValues = accountSliceConnected

	type InitializeConnected struct {
		Discriminator [8]byte
	}
	instConnected.DataBytes, err = borsh.Serialize(InitializeConnected{
		Discriminator: solanacontracts.DiscriminatorInitialize,
	})
	require.NoError(r, err)

	// create and sign the transaction
	signedTx = r.CreateSignedTransaction([]solana.Instruction{&instConnected}, privkey, []solana.PrivateKey{})

	// broadcast the transaction and wait for finalization
	_, out = r.BroadcastTxSync(signedTx)
	r.Logger.Info("initialize connected logs: %v", out.Meta.LogMessages)

	// initialize connected_spl program, use hardcoded program ID if not set
	if (r.ConnectedSPLProgram == solana.PublicKey{}) {
		r.ConnectedSPLProgram = ConnectedSPLProgramID
	}
	connectedSPLPda, err := solanacontracts.ComputeConnectedPdaAddress(r.ConnectedSPLProgram)
	require.NoError(r, err)

	var instConnectedSPL solana.GenericInstruction
	accountSliceConnectedSPL := []*solana.AccountMeta{}
	accountSliceConnectedSPL = append(accountSliceConnectedSPL, solana.Meta(privkey.PublicKey()).WRITE().SIGNER())
	accountSliceConnectedSPL = append(accountSliceConnectedSPL, solana.Meta(connectedSPLPda).WRITE())
	accountSliceConnectedSPL = append(accountSliceConnectedSPL, solana.Meta(solana.SystemProgramID))
	instConnectedSPL.ProgID = r.ConnectedSPLProgram
	instConnectedSPL.AccountValues = accountSliceConnectedSPL

	type InitializeConnectedSPL struct {
		Discriminator [8]byte
	}
	instConnectedSPL.DataBytes, err = borsh.Serialize(InitializeConnectedSPL{
		Discriminator: solanacontracts.DiscriminatorInitialize,
	})
	require.NoError(r, err)

	// create and sign the transaction
	signedTx = r.CreateSignedTransaction([]solana.Instruction{&instConnectedSPL}, privkey, []solana.PrivateKey{})

	// broadcast the transaction and wait for finalization
	_, out = r.BroadcastTxSync(signedTx)
	r.Logger.Info("initialize connected SPL logs: %v", out.Meta.LogMessages)

	// retrieve the PDA account info
	pdaInfo, err := r.SolanaClient.GetAccountInfoWithOpts(r.Ctx, pdaComputed, &rpc.GetAccountInfoOpts{
		Commitment: rpc.CommitmentConfirmed,
	})
	require.NoError(r, err)

	// deserialize the PDA info
	pda := solanacontracts.PdaInfo{}
	err = borsh.Deserialize(&pda, pdaInfo.Bytes())
	require.NoError(r, err)
	tssAddress := ethcommon.BytesToAddress(pda.TssAddress[:])

	// check the TSS address
	require.Equal(r, r.TSSAddress, tssAddress, "TSS address mismatch")

	// show the PDA balance
	balance, err := r.SolanaClient.GetBalance(r.Ctx, pdaComputed, rpc.CommitmentConfirmed)
	require.NoError(r, err)
	r.Logger.Info("initial PDA balance: %d lamports", balance.Value)

	err = r.ensureSolanaChainParams()
	require.NoError(r, err)

	// deploy test spl
	mintAccount := r.DeploySPL(&privkey, true)
	r.SPLAddr = mintAccount.PublicKey()

	// get spl account private key
	splPrivkey, err := solana.PrivateKeyFromBase58(splAccountPrivateKey)
	require.NoError(r, err)

	// minting some tokens to spl account for testing
	ata := r.ResolveSolanaATA(privkey, splPrivkey.PublicKey(), mintAccount.PublicKey())

	mintToInstruction := token.NewMintToInstruction(uint64(100_000_000_000_000), mintAccount.PublicKey(), ata, privkey.PublicKey(), []solana.PublicKey{}).
		Build()
	signedTx = r.CreateSignedTransaction(
		[]solana.Instruction{mintToInstruction},
		privkey,
		[]solana.PrivateKey{},
	)

	// broadcast the transaction and wait for finalization
	_, out = r.BroadcastTxSync(signedTx)
	r.Logger.Info("mint spl logs: %v", out.Meta.LogMessages)
}

func (r *E2ERunner) ensureSolanaChainParams() error {
	if r.ZetaTxServer == nil {
		return errors.New("ZetaTxServer is not initialized")
	}

	creator := r.ZetaTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName)

	chainID := chains.SolanaLocalnet.ChainId

	chainParams := &observertypes.ChainParams{
		ChainId:                     chainID,
		ConfirmationCount:           32,
		ZetaTokenContractAddress:    constant.EVMZeroAddress,
		ConnectorContractAddress:    constant.EVMZeroAddress,
		Erc20CustodyContractAddress: constant.EVMZeroAddress,
		GasPriceTicker:              5,
		WatchUtxoTicker:             0,
		InboundTicker:               2,
		OutboundTicker:              2,
		OutboundScheduleInterval:    2,
		OutboundScheduleLookahead:   20,
		BallotThreshold:             observertypes.DefaultBallotThreshold,
		MinObserverDelegation:       observertypes.DefaultMinObserverDelegation,
		IsSupported:                 true,
		GatewayAddress:              r.GatewayProgram.String(),
		ConfirmationParams: &observertypes.ConfirmationParams{
			SafeInboundCount:  32,
			SafeOutboundCount: 32,
		},
		GasPriceMultiplier: observertypes.DefaultGasPriceMultiplier,
	}

	if err := r.ZetaTxServer.UpdateChainParams(chainParams); err != nil {
		return errors.Wrap(err, "unable to broadcast solana chain params tx")
	}

	resetMsg := observertypes.NewMsgResetChainNonces(creator, chainID, 0, 0)
	if _, err := r.ZetaTxServer.BroadcastTx(utils.OperationalPolicyName, resetMsg); err != nil {
		return errors.Wrap(err, "unable to broadcast solana chain nonce reset tx")
	}

	r.Logger.Print("⚙️ voted for adding solana chain params (localnet). Waiting for confirmation")

	query := &observertypes.QueryGetChainParamsForChainRequest{ChainId: chainID}

	const duration = 2 * time.Second

	for i := 0; i < 10; i++ {
		_, err := r.ObserverClient.GetChainParamsForChain(r.Ctx, query)
		if err == nil {
			r.Logger.Print("⚙️ solana chain params are set")
			return nil
		}

		time.Sleep(duration)
	}

	return errors.New("unable to set Solana chain params")
}

// UpdateTSSAddressSolana updates the TSS address on the Solana gateway program
func (r *E2ERunner) UpdateTSSAddressSolana(gatewayID, deployerPrivateKey string) {
	r.Logger.Print("⚙️ updating tss on the gateway program on Solana")

	// set Solana contracts
	r.GatewayProgram = solana.MustPublicKeyFromBase58(gatewayID)

	// get deployer account balance
	privkey, err := solana.PrivateKeyFromBase58(deployerPrivateKey)
	require.NoError(r, err)
	pdaComputed := r.ComputePdaAddress()

	// create 'initialize' instruction
	var inst solana.GenericInstruction
	accountSlice := []*solana.AccountMeta{}
	accountSlice = append(accountSlice, solana.Meta(privkey.PublicKey()).WRITE().SIGNER())
	accountSlice = append(accountSlice, solana.Meta(pdaComputed).WRITE())
	accountSlice = append(accountSlice, solana.Meta(solana.SystemProgramID))
	inst.ProgID = r.GatewayProgram
	inst.AccountValues = accountSlice

	inst.DataBytes, err = borsh.Serialize(solanacontracts.UpdateTssParams{
		Discriminator: solanacontracts.DiscriminatorUpdateTss,
		TssAddress:    r.TSSAddress,
	})
	require.NoError(r, err)

	// create and sign the transaction
	signedTx := r.CreateSignedTransaction([]solana.Instruction{&inst}, privkey, []solana.PrivateKey{})

	// broadcast the transaction and wait for finalization
	_, out := r.BroadcastTxSync(signedTx)
	r.Logger.Info("update TSS gateway logs: %v", out.Meta.LogMessages)

	// retrieve the PDA account info
	pdaInfo, err := r.SolanaClient.GetAccountInfoWithOpts(r.Ctx, pdaComputed, &rpc.GetAccountInfoOpts{
		Commitment: rpc.CommitmentConfirmed,
	})
	require.NoError(r, err)

	pda := solanacontracts.PdaInfo{}
	err = borsh.Deserialize(&pda, pdaInfo.Bytes())
	require.NoError(r, err)
	tssAddress := ethcommon.BytesToAddress(pda.TssAddress[:])

	// verify updated TSS address
	require.Equal(r, r.TSSAddress, tssAddress, "TSS address mismatch")
}
