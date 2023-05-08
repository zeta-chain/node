package cli_test

//TODO : FIX this test
//func networkWithTSSObjects(t *testing.T) (*network.Network, *types.TSS) {
//	t.Helper()
//	cfg := network.DefaultConfig()
//	state := types.GenesisState{}
//	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))
//
//	state.Tss = &types.TSS{
//		TssPubkey:           "pubkey",
//		TssParticipantList:  []string{"ANY"},
//		OperatorAddressList: []string{"ANY"},
//		FinalizedZetaHeight: 11,
//		KeyGenZetaHeight:    10,
//	}
//	buf, err := cfg.Codec.MarshalJSON(&state)
//	require.NoError(t, err)
//	cfg.GenesisState[types.ModuleName] = buf
//	return network.New(t, cfg), state.Tss
//}
//
//func TestShowTSS(t *testing.T) {
//	net, objs := networkWithTSSObjects(t)
//
//	ctx := net.Validators[0].ClientCtx
//	common := []string{
//		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
//	}
//	for _, tc := range []struct {
//		desc string
//		args []string
//		err  error
//		obj  *types.TSS
//	}{
//		{
//			desc: "get",
//			args: common,
//			obj:  objs,
//			err:  nil,
//		},
//	} {
//		tc := tc
//		t.Run(tc.desc, func(t *testing.T) {
//			var args []string
//			args = append(args, tc.args...)
//			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowTSS(), args)
//			if tc.err != nil {
//				stat, ok := status.FromError(tc.err)
//				require.True(t, ok)
//				require.ErrorIs(t, stat.Err(), tc.err)
//			} else {
//				require.NoError(t, err)
//				var resp types.QueryGetTSSResponse
//				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
//				require.NotNil(t, resp.TSS)
//				require.Equal(t, tc.obj, resp.TSS)
//			}
//		})
//	}
//}
