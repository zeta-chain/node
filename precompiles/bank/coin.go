package bank

import "github.com/ethereum/go-ethereum/common"

func ZRC20ToCosmosDenom(ZRC20Address common.Address) string {
	return ZEVMDenom + ZRC20Address.String()
}
