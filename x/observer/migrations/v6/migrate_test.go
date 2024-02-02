package v6_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	v6 "github.com/zeta-chain/zetacore/x/observer/migrations/v6"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMigrateObserverParams(t *testing.T) {
	t.Run("Migrate when keygen is Pending", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		k.SetKeygen(ctx, types.Keygen{
			Status:      types.KeygenStatus_PendingKeygen,
			BlockNumber: math.MaxInt64,
		})
		participantList := []string{
			"zetapub1addwnpepqglunjrgl3qg08duxq9pf28jmvrer3crwnnfzp6m0u0yh9jk9mnn5p76utc",
			"zetapub1addwnpepqwwpjwwnes7cywfkr0afme7ymk8rf5jzhn8pfr6qqvfm9v342486qsrh4f5",
			"zetapub1addwnpepq07xj82w5e6vr85qj3r7htmzh2mp3vkvfraapcv6ynhdwankseayk5yh80t",
			"zetapub1addwnpepq0lxqx92m3fhae3usn8jffqvtx6cuzl06xh9r345c2qcqq8zyfs4cdpqcum",
			"zetapub1addwnpepqvzlntzltvpm22ved5gjtn9nzqfz5fun38el4r64njc979rwanxlgq4u3p8",
			"zetapub1addwnpepqg40psrhwwgy257p4xv50xp0asmtwjup66z8vk829289zxge5lyl7sycga8",
			"zetapub1addwnpepqgpr5ffquqchra93r8l6d35q62cv4nsc9d4k2q7kten4sljxg5rluwx29gh",
			"zetapub1addwnpepqdjf3vt8etgdddkghrvxfmmmeatky6m7hx7wjuv86udfghqpty8h5h4r78w",
			"zetapub1addwnpepqtfcfmsdkzdgv03t8392gsh7kzrstp9g864w2ltz9k0xzz33q60dq6mnkex",
		}
		operatorList := []string{
			"zeta19jr7nl82lrktge35f52x9g5y5prmvchmk40zhg",
			"zeta1cxj07f3ju484ry2cnnhxl5tryyex7gev0yzxtj",
			"zeta1hjct6q7npsspsg3dgvzk3sdf89spmlpf7rqmnw",
			"zeta1k6vh9y7ctn06pu5jngznv5dyy0rltl2qp0j30g",
			"zeta1l07weaxkmn6z69qm55t53v4rfr43eys4cjz54h",
			"zeta1p0uwsq4naus5r4l7l744upy0k8ezzj84mn40nf",
			"zeta1rhj4pkp7eygw8lu9wacpepeh0fnzdxrqr27g6m",
			"zeta1t0uj2z93jd2g3w94zl3jhfrn2ek6dnuk3v93j9",
			"zeta1t5pgk2fucx3drkynzew9zln5z9r7s3wqqyy0pe",
		}
		keygenHeight := int64(1440460)
		finalizedZetaHeight := int64(1440680)
		k.SetTSS(ctx, types.TSS{
			KeyGenZetaHeight:    keygenHeight,
			TssParticipantList:  participantList,
			TssPubkey:           "zetapub1addwnpepqtadxdyt037h86z60nl98t6zk56mw5zpnm79tsmvspln3hgt5phdc79kvfc",
			OperatorAddressList: operatorList,
			FinalizedZetaHeight: finalizedZetaHeight,
		})
		err := v6.MigrateStore(ctx, k)
		assert.NoError(t, err)
		keygen, found := k.GetKeygen(ctx)
		assert.True(t, found)
		assert.Equal(t, types.KeygenStatus_KeyGenSuccess, keygen.Status)
		assert.Equal(t, keygenHeight, keygenHeight)
		assert.Equal(t, participantList, participantList)
	})
	t.Run("Migrate when keygen is not Pending", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		participantList := []string{
			"zetapub1addwnpepqglunjrgl3qg08duxq9pf28jmvrer3crwnnfzp6m0u0yh9jk9mnn5p76utc",
			"zetapub1addwnpepqwwpjwwnes7cywfkr0afme7ymk8rf5jzhn8pfr6qqvfm9v342486qsrh4f5",
			"zetapub1addwnpepq07xj82w5e6vr85qj3r7htmzh2mp3vkvfraapcv6ynhdwankseayk5yh80t",
			"zetapub1addwnpepq0lxqx92m3fhae3usn8jffqvtx6cuzl06xh9r345c2qcqq8zyfs4cdpqcum",
			"zetapub1addwnpepqvzlntzltvpm22ved5gjtn9nzqfz5fun38el4r64njc979rwanxlgq4u3p8",
			"zetapub1addwnpepqg40psrhwwgy257p4xv50xp0asmtwjup66z8vk829289zxge5lyl7sycga8",
			"zetapub1addwnpepqgpr5ffquqchra93r8l6d35q62cv4nsc9d4k2q7kten4sljxg5rluwx29gh",
			"zetapub1addwnpepqdjf3vt8etgdddkghrvxfmmmeatky6m7hx7wjuv86udfghqpty8h5h4r78w",
			"zetapub1addwnpepqtfcfmsdkzdgv03t8392gsh7kzrstp9g864w2ltz9k0xzz33q60dq6mnkex",
		}
		keygenHeight := int64(1440460)
		k.SetKeygen(ctx, types.Keygen{
			Status:         types.KeygenStatus_KeyGenSuccess,
			BlockNumber:    keygenHeight,
			GranteePubkeys: participantList,
		})
		err := v6.MigrateStore(ctx, k)
		assert.NoError(t, err)
		keygen, found := k.GetKeygen(ctx)
		assert.True(t, found)
		assert.Equal(t, types.KeygenStatus_KeyGenSuccess, keygen.Status)
		assert.Equal(t, keygen.BlockNumber, keygenHeight)
		assert.Equal(t, keygen.GranteePubkeys, participantList)
	})

}
