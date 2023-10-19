package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestMsgServer_AddToOutTxTracker(t *testing.T) {
	t.Run("Add proof based tracker with correct proof", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)
		//txHash := "string"
		chainID := int64(5)
		txIndex, block, header, headerRLP, _, err := sample.Proof()
		require.NoError(t, err)
		SetupVerificationParams(zk, ctx, txIndex, chainID, header, headerRLP, block)
		k.SetTSS(ctx, types.TSS{
			TssPubkey: "zetapub1addwnpepq28c57cvcs0a2htsem5zxr6qnlvq9mzhmm76z3jncsnzz32rclangr2g35p",
			TssParticipantList: []string{
				"zetapub1addwnpepqdljmtptxflv3cs3q79v8pvaynu7ep9eglwvyxv7v3a5mpccvmcg7hpywrc",
				"zetapub1addwnpepqf8ngvcs27m4hm8m2e8fupulu3qsvutf7lflds2c0deqmamkdrymvp5u0g6",
				"zetapub1addwnpepq2ayeqvgmdau45rqxapal8j5jd6we55y3j7jlphqner8mvtj3jl454mpuhn",
				"zetapub1addwnpepqf7dw2fnf6ntl6f8d3yw7jly8a535cgum9tuur55k2c064pq0ysl5u2fw5m",
				"zetapub1addwnpepq09ywklml6wts0eua3zhmtcw38rf4kel0zg2d4dhd46zv00mkjz45gn9phl",
				"zetapub1addwnpepqdj8ww4u40wkg2whl22qvmm2zdsdn5rt4yegr6zc6sw9zqdvr8mx54nhke5",
				"zetapub1addwnpepqwsl5cqjj57aw5rn79j44pyp0322nu6cxx28hfrqn5gwfcaau27tgm3p0z2",
				"zetapub1addwnpepqtft2rac0xtu25v3knnhktyjyp0wehfwj0lg8z24g5asfu0lf0dh7xqyfpf",
				"zetapub1addwnpepq02ppdhgdnffaz9mkswx8xu04wlcqjtqc6tgnmsj62fa9ef202v77c9ugp6",
				"zetapub1addwnpepq2vduhrgfmvjjjvmh7h90ewm9d98gtj2xkvvduqcgr7lk0ejxgltsngchvu",
				"zetapub1addwnpepq08m4vw832r4hwt6nwa0fa6ze8c4ue6k9yqjul9ap25hcndv44fsquw5eg2",
				"zetapub1addwnpepqdy8mm6jemgeyv9g7f7ymk5cymal4zms6t2n0ml0tk9ajak0c8n77790s52",
			},
			OperatorAddressList: []string{
				"zeta15ruj2tc76pnj9xtw64utktee7cc7w6vzaes73z",
				"zeta18f7wch6kpfdmk6dk9qqhkszpjwrymev4fpte8p",
				"zeta18pksjzclks34qkqyaahf2rakss80mnusju77cm",
				"zeta1dxyzsket66vt886ap0gnzlnu5pv0y99v086wnz",
				"zeta1g323lusfa9qqvjvupajre2dphuem999fahc086",
				"zeta1ggqzjf5726uu7xc6pfwg00lny79w6t3a3utpw5",
				"zeta1hk05v9len8u0c2xrwxgfknvcskpd4vncm7ehch",
				"zeta1j8g8ch4uqgl3gtet3nntvczaeppmlxajqwh5u6",
				"zeta1mte0r3jzkf2rkd7ex4p3xsd3fxqg7q29q0wxl5",
				"zeta1w5czgpk5kc9etxw2anzhr0uyrr4fqks32qmk6k",
				"zeta1w8qa37h22h884vxedmprvwtd3z2nwakxu9k935",
				"zeta1ymnrwg9e3xr9xkw42ygzjx34dyvwvtc23cnnxz",
			},
			FinalizedZetaHeight: 287833,
			KeyGenZetaHeight:    287830,
		})

	})
}
