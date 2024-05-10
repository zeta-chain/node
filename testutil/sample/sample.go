package sample

import (
	"encoding/json"
	"errors"
	"hash/fnv"
	"math/rand"
	"strconv"
	"testing"

	"github.com/zeta-chain/zetacore/pkg/chains"

	sdkmath "cosmossdk.io/math"
	"github.com/cometbft/cometbft/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
)

var ErrSample = errors.New("sample error")

func newRandFromSeed(s int64) *rand.Rand {
	// #nosec G404 test purpose - weak randomness is not an issue here
	return rand.New(rand.NewSource(s))
}

func newRandFromStringSeed(t *testing.T, s string) *rand.Rand {
	h := fnv.New64a()
	_, err := h.Write([]byte(s))
	require.NoError(t, err)
	return newRandFromSeed(int64(h.Sum64()))
}

// Rand returns a new random number generator
func Rand() *rand.Rand {
	return newRandFromSeed(42)
}

// Validator returns a sample staking validator
func Validator(t testing.TB, r *rand.Rand) stakingtypes.Validator {
	seed := []byte(strconv.Itoa(r.Int()))
	val, err := stakingtypes.NewValidator(
		ValAddress(r),
		ed25519.GenPrivKeyFromSecret(seed).PubKey(),
		stakingtypes.Description{})
	require.NoError(t, err)
	return val
}

func ZetaIndex(t *testing.T) string {
	msg := CrossChainTx(t, "foo")
	hash := ethcrypto.Keccak256Hash([]byte(msg.String()))
	return hash.Hex()
}

// Bytes returns a sample byte array
func Bytes() []byte {
	return []byte("sample")
}

// String returns a sample string
func String() string {
	return "sample"
}

// StringRandom returns a sample string with random alphanumeric characters
func StringRandom(r *rand.Rand, length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}

// Coins returns a sample sdk.Coins
func Coins() sdk.Coins {
	return sdk.NewCoins(sdk.NewCoin(config.BaseDenom, sdk.NewInt(42)))
}

// Uint64InRange returns a sample uint64 in the given ranges
func Uint64InRange(low, high uint64) uint64 {
	r := newRandFromSeed(int64(low))
	return r.Uint64()%(high-low) + low
}

// Int64InRange returns a sample int64 in the given ranges
func Int64InRange(low, high int64) int64 {
	r := newRandFromSeed(low)
	return r.Int63()%(high-low) + low
}

func UintInRange(low, high uint64) sdkmath.Uint {
	u := Uint64InRange(low, high)
	return sdkmath.NewUint(u)
}

func IntInRange(low, high int64) sdkmath.Int {
	i := Int64InRange(low, high)
	return sdkmath.NewInt(i)
}
func AppState(t *testing.T) map[string]json.RawMessage {
	appState, err := genutiltypes.GenesisStateFromGenDoc(*GenDoc(t))
	require.NoError(t, err)
	return appState
}

func GenDoc(t *testing.T) *types.GenesisDoc {
	jsonBlob := []byte("{\n    \"genesis_time\": \"2024-04-12T05:07:56.004517Z\",\n    \"chain_id\": \"localnet_101-1\",\n    \"initial_height\": \"1\",\n    \"consensus_params\": {\n      \"block\": {\n        \"max_bytes\": \"22020096\",\n        \"max_gas\": \"10000000\",\n        \"time_iota_ms\": \"1000\"\n      },\n      \"evidence\": {\n        \"max_age_num_blocks\": \"100000\",\n        \"max_age_duration\": \"172800000000000\",\n        \"max_bytes\": \"1048576\"\n      },\n      \"validator\": {\n        \"pub_key_types\": [\n          \"ed25519\"\n        ]\n      },\n      \"version\": {}\n    },\n    \"app_hash\": \"\",\n    \"app_state\": {\n      \"auth\": {\n        \"params\": {\n          \"max_memo_characters\": \"256\",\n          \"tx_sig_limit\": \"7\",\n          \"tx_size_cost_per_byte\": \"10\",\n          \"sig_verify_cost_ed25519\": \"590\",\n          \"sig_verify_cost_secp256k1\": \"1000\"\n        },\n        \"accounts\": [\n          {\n            \"@type\": \"/ethermint.types.v1.EthAccount\",\n            \"base_account\": {\n              \"address\": \"zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax\",\n              \"pub_key\": null,\n              \"account_number\": \"0\",\n              \"sequence\": \"0\"\n            },\n            \"code_hash\": \"0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470\"\n          },\n          {\n            \"@type\": \"/ethermint.types.v1.EthAccount\",\n            \"base_account\": {\n              \"address\": \"zeta10up34mvwjhjd9xkq56fwsf0k75vtg287uav69n\",\n              \"pub_key\": null,\n              \"account_number\": \"0\",\n              \"sequence\": \"0\"\n            },\n            \"code_hash\": \"0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470\"\n          },\n          {\n            \"@type\": \"/ethermint.types.v1.EthAccount\",\n            \"base_account\": {\n              \"address\": \"zeta1f203dypqg5jh9hqfx0gfkmmnkdfuat3jr45ep2\",\n              \"pub_key\": null,\n              \"account_number\": \"0\",\n              \"sequence\": \"0\"\n            },\n            \"code_hash\": \"0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470\"\n          },\n          {\n            \"@type\": \"/ethermint.types.v1.EthAccount\",\n            \"base_account\": {\n              \"address\": \"zeta1unzpyll3tmutf0r8sqpxpnj46vtdr59mw8qepx\",\n              \"pub_key\": null,\n              \"account_number\": \"0\",\n              \"sequence\": \"0\"\n            },\n            \"code_hash\": \"0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470\"\n          }\n        ]\n      },\n      \"authority\": {\n        \"policies\": {\n          \"items\": [\n            {\n              \"policy_type\": \"groupEmergency\",\n              \"address\": \"zeta1afk9zr2hn2jsac63h4hm60vl9z3e5u69gndzf7c99cqge3vzwjzsxn0x73\"\n            },\n            {\n              \"policy_type\": \"groupOperational\",\n              \"address\": \"zeta1afk9zr2hn2jsac63h4hm60vl9z3e5u69gndzf7c99cqge3vzwjzsxn0x73\"\n            },\n            {\n              \"policy_type\": \"groupAdmin\",\n              \"address\": \"zeta1afk9zr2hn2jsac63h4hm60vl9z3e5u69gndzf7c99cqge3vzwjzsxn0x73\"\n            }\n          ]\n        }\n      },\n      \"authz\": {\n        \"authorization\": [\n          {\n            \"granter\": \"zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax\",\n            \"grantee\": \"zeta10up34mvwjhjd9xkq56fwsf0k75vtg287uav69n\",\n            \"authorization\": {\n              \"@type\": \"/cosmos.authz.v1beta1.GenericAuthorization\",\n              \"msg\": \"/zetachain.zetacore.crosschain.MsgGasPriceVoter\"\n            },\n            \"expiration\": null\n          },\n          {\n            \"granter\": \"zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax\",\n            \"grantee\": \"zeta10up34mvwjhjd9xkq56fwsf0k75vtg287uav69n\",\n            \"authorization\": {\n              \"@type\": \"/cosmos.authz.v1beta1.GenericAuthorization\",\n              \"msg\": \"/zetachain.zetacore.crosschain.MsgVoteInbound\"\n            },\n            \"expiration\": null\n          },\n          {\n            \"granter\": \"zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax\",\n            \"grantee\": \"zeta10up34mvwjhjd9xkq56fwsf0k75vtg287uav69n\",\n            \"authorization\": {\n              \"@type\": \"/cosmos.authz.v1beta1.GenericAuthorization\",\n              \"msg\": \"/zetachain.zetacore.crosschain.MsgVoteOutbound\"\n            },\n            \"expiration\": null\n          },\n          {\n            \"granter\": \"zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax\",\n            \"grantee\": \"zeta10up34mvwjhjd9xkq56fwsf0k75vtg287uav69n\",\n            \"authorization\": {\n              \"@type\": \"/cosmos.authz.v1beta1.GenericAuthorization\",\n              \"msg\": \"/zetachain.zetacore.crosschain.MsgCreateTSSVoter\"\n            },\n            \"expiration\": null\n          },\n          {\n            \"granter\": \"zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax\",\n            \"grantee\": \"zeta10up34mvwjhjd9xkq56fwsf0k75vtg287uav69n\",\n            \"authorization\": {\n              \"@type\": \"/cosmos.authz.v1beta1.GenericAuthorization\",\n              \"msg\": \"/zetachain.zetacore.crosschain.MsgAddOutboundTracker\"\n            },\n            \"expiration\": null\n          },\n          {\n            \"granter\": \"zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax\",\n            \"grantee\": \"zeta10up34mvwjhjd9xkq56fwsf0k75vtg287uav69n\",\n            \"authorization\": {\n              \"@type\": \"/cosmos.authz.v1beta1.GenericAuthorization\",\n              \"msg\": \"/zetachain.zetacore.observer.MsgAddBlameVote\"\n            },\n            \"expiration\": null\n          },\n          {\n            \"granter\": \"zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax\",\n            \"grantee\": \"zeta10up34mvwjhjd9xkq56fwsf0k75vtg287uav69n\",\n            \"authorization\": {\n              \"@type\": \"/cosmos.authz.v1beta1.GenericAuthorization\",\n              \"msg\": \"/zetachain.zetacore.observer.MsgAddBlockHeader\"\n            },\n            \"expiration\": null\n          },\n          {\n            \"granter\": \"zeta1f203dypqg5jh9hqfx0gfkmmnkdfuat3jr45ep2\",\n            \"grantee\": \"zeta1unzpyll3tmutf0r8sqpxpnj46vtdr59mw8qepx\",\n            \"authorization\": {\n              \"@type\": \"/cosmos.authz.v1beta1.GenericAuthorization\",\n              \"msg\": \"/zetachain.zetacore.crosschain.MsgGasPriceVoter\"\n            },\n            \"expiration\": null\n          },\n          {\n            \"granter\": \"zeta1f203dypqg5jh9hqfx0gfkmmnkdfuat3jr45ep2\",\n            \"grantee\": \"zeta1unzpyll3tmutf0r8sqpxpnj46vtdr59mw8qepx\",\n            \"authorization\": {\n              \"@type\": \"/cosmos.authz.v1beta1.GenericAuthorization\",\n              \"msg\": \"/zetachain.zetacore.crosschain.MsgVoteInbound\"\n            },\n            \"expiration\": null\n          },\n          {\n            \"granter\": \"zeta1f203dypqg5jh9hqfx0gfkmmnkdfuat3jr45ep2\",\n            \"grantee\": \"zeta1unzpyll3tmutf0r8sqpxpnj46vtdr59mw8qepx\",\n            \"authorization\": {\n              \"@type\": \"/cosmos.authz.v1beta1.GenericAuthorization\",\n              \"msg\": \"/zetachain.zetacore.crosschain.MsgVoteOutbound\"\n            },\n            \"expiration\": null\n          },\n          {\n            \"granter\": \"zeta1f203dypqg5jh9hqfx0gfkmmnkdfuat3jr45ep2\",\n            \"grantee\": \"zeta1unzpyll3tmutf0r8sqpxpnj46vtdr59mw8qepx\",\n            \"authorization\": {\n              \"@type\": \"/cosmos.authz.v1beta1.GenericAuthorization\",\n              \"msg\": \"/zetachain.zetacore.crosschain.MsgCreateTSSVoter\"\n            },\n            \"expiration\": null\n          },\n          {\n            \"granter\": \"zeta1f203dypqg5jh9hqfx0gfkmmnkdfuat3jr45ep2\",\n            \"grantee\": \"zeta1unzpyll3tmutf0r8sqpxpnj46vtdr59mw8qepx\",\n            \"authorization\": {\n              \"@type\": \"/cosmos.authz.v1beta1.GenericAuthorization\",\n              \"msg\": \"/zetachain.zetacore.crosschain.MsgAddOutboundTracker\"\n            },\n            \"expiration\": null\n          },\n          {\n            \"granter\": \"zeta1f203dypqg5jh9hqfx0gfkmmnkdfuat3jr45ep2\",\n            \"grantee\": \"zeta1unzpyll3tmutf0r8sqpxpnj46vtdr59mw8qepx\",\n            \"authorization\": {\n              \"@type\": \"/cosmos.authz.v1beta1.GenericAuthorization\",\n              \"msg\": \"/zetachain.zetacore.observer.MsgAddBlameVote\"\n            },\n            \"expiration\": null\n          },\n          {\n            \"granter\": \"zeta1f203dypqg5jh9hqfx0gfkmmnkdfuat3jr45ep2\",\n            \"grantee\": \"zeta1unzpyll3tmutf0r8sqpxpnj46vtdr59mw8qepx\",\n            \"authorization\": {\n              \"@type\": \"/cosmos.authz.v1beta1.GenericAuthorization\",\n              \"msg\": \"/zetachain.zetacore.observer.MsgAddBlockHeader\"\n            },\n            \"expiration\": null\n          }\n        ]\n      },\n      \"bank\": {\n        \"params\": {\n          \"send_enabled\": [],\n          \"default_send_enabled\": true\n        },\n        \"balances\": [\n          {\n            \"address\": \"zeta1f203dypqg5jh9hqfx0gfkmmnkdfuat3jr45ep2\",\n            \"coins\": [\n              {\n                \"denom\": \"azeta\",\n                \"amount\": \"4200000000000000000000000\"\n              }\n            ]\n          },\n          {\n            \"address\": \"zeta10up34mvwjhjd9xkq56fwsf0k75vtg287uav69n\",\n            \"coins\": [\n              {\n                \"denom\": \"azeta\",\n                \"amount\": \"1000000000000000000000\"\n              }\n            ]\n          },\n          {\n            \"address\": \"zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax\",\n            \"coins\": [\n              {\n                \"denom\": \"azeta\",\n                \"amount\": \"4200000000000000000000000\"\n              }\n            ]\n          },\n          {\n            \"address\": \"zeta1unzpyll3tmutf0r8sqpxpnj46vtdr59mw8qepx\",\n            \"coins\": [\n              {\n                \"denom\": \"azeta\",\n                \"amount\": \"1000000000000000000000\"\n              }\n            ]\n          }\n        ],\n        \"supply\": [\n          {\n            \"denom\": \"azeta\",\n            \"amount\": \"8402000000000000000000000\"\n          }\n        ],\n        \"denom_metadata\": []\n      },\n      \"crisis\": {\n        \"constant_fee\": {\n          \"denom\": \"azeta\",\n          \"amount\": \"1000\"\n        }\n      },\n      \"crosschain\": {\n        \"outboundTrackerList\": [],\n        \"inboundHashToCctxList\": [],\n        \"inbound_tracker_list\": [],\n        \"zeta_accounting\": {\n          \"aborted_zeta_amount\": \"0\"\n        }\n      },\n      \"distribution\": {\n        \"params\": {\n          \"community_tax\": \"0.020000000000000000\",\n          \"base_proposer_reward\": \"0.010000000000000000\",\n          \"bonus_proposer_reward\": \"0.040000000000000000\",\n          \"withdraw_addr_enabled\": true\n        },\n        \"fee_pool\": {\n          \"community_pool\": []\n        },\n        \"delegator_withdraw_infos\": [],\n        \"previous_proposer\": \"\",\n        \"outstanding_rewards\": [],\n        \"validator_accumulated_commissions\": [],\n        \"validator_historical_rewards\": [],\n        \"validator_current_rewards\": [],\n        \"delegator_starting_infos\": [],\n        \"validator_slash_events\": []\n      },\n      \"emissions\": {\n        \"params\": {\n          \"max_bond_factor\": \"1.25\",\n          \"min_bond_factor\": \"0.75\",\n          \"avg_block_time\": \"6.00\",\n          \"target_bond_ratio\": \"00.67\",\n          \"validator_emission_percentage\": \"00.50\",\n          \"observer_emission_percentage\": \"00.25\",\n          \"tss_signer_emission_percentage\": \"00.25\",\n          \"duration_factor_constant\": \"0.001877876953694702\",\n          \"observer_slash_amount\": \"0\"\n        },\n        \"withdrawableEmissions\": []\n      },\n      \"evidence\": {\n        \"evidence\": []\n      },\n      \"evm\": {\n        \"accounts\": [],\n        \"params\": {\n          \"evm_denom\": \"azeta\",\n          \"enable_create\": true,\n          \"enable_call\": true,\n          \"extra_eips\": [],\n          \"chain_config\": {\n            \"homestead_block\": \"0\",\n            \"dao_fork_block\": \"0\",\n            \"dao_fork_support\": true,\n            \"eip150_block\": \"0\",\n            \"eip150_hash\": \"0x0000000000000000000000000000000000000000000000000000000000000000\",\n            \"eip155_block\": \"0\",\n            \"eip158_block\": \"0\",\n            \"byzantium_block\": \"0\",\n            \"constantinople_block\": \"0\",\n            \"petersburg_block\": \"0\",\n            \"istanbul_block\": \"0\",\n            \"muir_glacier_block\": \"0\",\n            \"berlin_block\": \"0\",\n            \"london_block\": \"0\",\n            \"arrow_glacier_block\": \"0\",\n            \"gray_glacier_block\": \"0\",\n            \"merge_netsplit_block\": \"0\",\n            \"shanghai_block\": \"0\",\n            \"cancun_block\": \"0\"\n          },\n          \"allow_unprotected_txs\": false\n        }\n      },\n      \"feemarket\": {\n        \"params\": {\n          \"no_base_fee\": false,\n          \"base_fee_change_denominator\": 8,\n          \"elasticity_multiplier\": 2,\n          \"enable_height\": \"0\",\n          \"base_fee\": \"1000000000\",\n          \"min_gas_price\": \"0.000000000000000000\",\n          \"min_gas_multiplier\": \"0.500000000000000000\"\n        },\n        \"block_gas\": \"0\"\n      },\n      \"fungible\": {\n        \"params\": {},\n        \"foreignCoinsList\": [],\n        \"systemContract\": null\n      },\n      \"genutil\": {\n        \"gen_txs\": [\n          {\n            \"body\": {\n              \"messages\": [\n                {\n                  \"@type\": \"/cosmos.staking.v1beta1.MsgCreateValidator\",\n                  \"description\": {\n                    \"moniker\": \"Zetanode-Localnet\",\n                    \"identity\": \"\",\n                    \"website\": \"\",\n                    \"security_contact\": \"\",\n                    \"details\": \"\"\n                  },\n                  \"commission\": {\n                    \"rate\": \"0.100000000000000000\",\n                    \"max_rate\": \"0.200000000000000000\",\n                    \"max_change_rate\": \"0.010000000000000000\"\n                  },\n                  \"min_self_delegation\": \"1\",\n                  \"delegator_address\": \"zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax\",\n                  \"validator_address\": \"zetavaloper13c7p3xrhd6q2rx3h235jpt8pjdwvacyw7tkass\",\n                  \"pubkey\": {\n                    \"@type\": \"/cosmos.crypto.ed25519.PubKey\",\n                    \"key\": \"sBSs5r1vQn1idTp4uRTbdUK0jjmEscI3pn88LUXI4CQ=\"\n                  },\n                  \"value\": {\n                    \"denom\": \"azeta\",\n                    \"amount\": \"1000000000000000000000\"\n                  }\n                }\n              ],\n              \"memo\": \"1db4f4185e68c1c17d508294de2592616dad37a5@192.168.2.12:26656\",\n              \"timeout_height\": \"0\",\n              \"extension_options\": [],\n              \"non_critical_extension_options\": []\n            },\n            \"auth_info\": {\n              \"signer_infos\": [\n                {\n                  \"public_key\": {\n                    \"@type\": \"/cosmos.crypto.secp256k1.PubKey\",\n                    \"key\": \"A05F6QuFVpb/5KrIPvlHr209ZsD22gW0omhLSXWAtQrh\"\n                  },\n                  \"mode_info\": {\n                    \"single\": {\n                      \"mode\": \"SIGN_MODE_DIRECT\"\n                    }\n                  },\n                  \"sequence\": \"0\"\n                }\n              ],\n              \"fee\": {\n                \"amount\": [],\n                \"gas_limit\": \"200000\",\n                \"payer\": \"\",\n                \"granter\": \"\"\n              },\n              \"tip\": null\n            },\n            \"signatures\": [\n              \"y5YROwZmV0jcgv5BgRJCDE+Kq5OsX8+88or1ogekPLBw3ecPt8GsCeEbPQ24JONLzNwQEIUDNYTeSQnXnCfzyg==\"\n            ]\n          }\n        ]\n      },\n      \"gov\": {\n        \"starting_proposal_id\": \"1\",\n        \"deposits\": [],\n        \"votes\": [],\n        \"proposals\": [],\n        \"deposit_params\": {\n          \"min_deposit\": [\n            {\n              \"denom\": \"azeta\",\n              \"amount\": \"10000000\"\n            }\n          ],\n          \"max_deposit_period\": \"172800s\"\n        },\n        \"voting_params\": {\n          \"voting_period\": \"10s\"\n        },\n        \"tally_params\": {\n          \"quorum\": \"0.334000000000000000\",\n          \"threshold\": \"0.500000000000000000\",\n          \"veto_threshold\": \"0.334000000000000000\"\n        }\n      },\n      \"group\": {\n        \"group_seq\": \"0\",\n        \"groups\": [],\n        \"group_members\": [],\n        \"group_policy_seq\": \"0\",\n        \"group_policies\": [],\n        \"proposal_seq\": \"0\",\n        \"proposals\": [],\n        \"votes\": []\n      },\n      \"mint\": {\n        \"params\": {\n          \"mint_denom\": \"azeta\"\n        }\n      },\n      \"observer\": {\n        \"observers\": {\n          \"observer_list\": [\n            \"zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax\",\n            \"zeta1f203dypqg5jh9hqfx0gfkmmnkdfuat3jr45ep2\"\n          ]\n        },\n        \"nodeAccountList\": [\n          {\n            \"operator\": \"zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax\",\n            \"granteeAddress\": \"zeta10up34mvwjhjd9xkq56fwsf0k75vtg287uav69n\",\n            \"granteePubkey\": {\n              \"secp256k1\": \"zetapub1addwnpepqtlu7fykuh875xjckz4mn4x0mzc25rrqk5qne7mrwxqmatgllv3nx6lrkdp\"\n            },\n            \"nodeStatus\": 4\n          },\n          {\n            \"operator\": \"zeta1f203dypqg5jh9hqfx0gfkmmnkdfuat3jr45ep2\",\n            \"granteeAddress\": \"zeta1unzpyll3tmutf0r8sqpxpnj46vtdr59mw8qepx\",\n            \"granteePubkey\": {\n              \"secp256k1\": \"zetapub1addwnpepqwy5pmg39regpq0gkggxehmfm8hwmxxw94sch7qzh4smava0szs07kk5045\"\n            },\n            \"nodeStatus\": 4\n          }\n        ],\n        \"crosschain_flags\": {\n          \"isInboundEnabled\": true,\n          \"isOutboundEnabled\": true\n        },\n        \"keygen\": {\n          \"status\": 1,\n          \"granteePubkeys\": [\n            \"zetapub1addwnpepqtlu7fykuh875xjckz4mn4x0mzc25rrqk5qne7mrwxqmatgllv3nx6lrkdp\",\n            \"zetapub1addwnpepqwy5pmg39regpq0gkggxehmfm8hwmxxw94sch7qzh4smava0szs07kk5045\"\n          ]\n        },\n        \"chain_params_list\": {},\n        \"tss\": {\n          \"tss_pubkey\": \"zetapub1addwnpepq28c57cvcs0a2htsem5zxr6qnlvq9mzhmm76z3jncsnzz32rclangr2g35p\",\n          \"tss_participant_list\": [\n            \"zetapub1addwnpepqtlu7fykuh875xjckz4mn4x0mzc25rrqk5qne7mrwxqmatgllv3nx6lrkdp\",\n            \"zetapub1addwnpepqwy5pmg39regpq0gkggxehmfm8hwmxxw94sch7qzh4smava0szs07kk5045\"\n          ],\n          \"operator_address_list\": [\n            \"zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax\",\n            \"zeta1f203dypqg5jh9hqfx0gfkmmnkdfuat3jr45ep2\"\n          ]\n        },\n        \"tss_history\": [],\n        \"tss_fund_migrators\": [],\n        \"blame_list\": [],\n        \"pending_nonces\": [],\n        \"chain_nonces\": [],\n        \"nonce_to_cctx\": []\n      },\n      \"params\": null,\n      \"slashing\": {\n        \"params\": {\n          \"signed_blocks_window\": \"100\",\n          \"min_signed_per_window\": \"0.500000000000000000\",\n          \"downtime_jail_duration\": \"600s\",\n          \"slash_fraction_double_sign\": \"0.050000000000000000\",\n          \"slash_fraction_downtime\": \"0.010000000000000000\"\n        },\n        \"signing_infos\": [],\n        \"missed_blocks\": []\n      },\n      \"staking\": {\n        \"params\": {\n          \"unbonding_time\": \"1814400s\",\n          \"max_validators\": 100,\n          \"max_entries\": 7,\n          \"historical_entries\": 10000,\n          \"bond_denom\": \"azeta\",\n          \"min_commission_rate\": \"0.000000000000000000\"\n        },\n        \"last_total_power\": \"0\",\n        \"last_validator_powers\": [],\n        \"validators\": [],\n        \"delegations\": [],\n        \"unbonding_delegations\": [],\n        \"redelegations\": [],\n        \"exported\": false\n      },\n      \"upgrade\": {},\n      \"vesting\": {}\n    }\n  }")
	genDoc, err := types.GenesisDocFromJSON(jsonBlob)
	require.NoError(t, err)
	return genDoc
}

func Chain(chainID int64) *chains.Chain {
	r := newRandFromSeed(chainID)

	return &chains.Chain{
		ChainName: chains.ChainName(r.Intn(4)),
		ChainId:   chainID,
	}
}

func EventIndex() uint64 {
	r := newRandFromSeed(1)
	return r.Uint64()
}
