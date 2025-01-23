package errors_test

import (
	"fmt"
	"testing"

	"cosmossdk.io/errors"
	cctxerror "github.com/zeta-chain/node/pkg/errors"
)

//func TestEvmErrorMessage(t *testing.T) {
//	t.Run("TestEvmErrorMessage", func(t *testing.T) {
//		contractAddress := "0xE97Ac2CA30D30de65a6FE0Ab20EDC39a623c18df"
//		msg := cctxerror.NewZEVMErrorMessage("method", common.HexToAddress(contractAddress), "args", "message", errors.New("error_cause"))
//		msg.AddRevertReason("revert_reason")
//
//		s, err := msg.ToJSON()
//		require.NoError(t, err)
//
//		require.Equal(
//			t,
//			`{"message":"message","error":"error_cause","method":"method","contract":"0xE97Ac2CA30D30de65a6FE0Ab20EDC39a623c18df","args":"args","revert_reason":"revert_reason"}`,
//			s,
//		)
//	})
//}
//
//func TestParseEvmErrorMessage(t *testing.T) {
//	t.Run("TestParseEvmErrorMessage", func(t *testing.T) {
//		m := `{"message":"contract call failed when calling EVM with data","method":"depositAndCall0","contract":"0x733aB8b06DDDEf27Eaa72294B0d7c9cEF7f12db9","args":"[{[]0xdFb74337c53141bf912101b0Ee770FA8e2DCB921 1337} 0x13A0c5930C028511Dc02665E7285134B6d11A5f410000000000000000 0xD28D6A0b8189305551a0A8bd247a6ECa9CE781Ca [114 101 118 101114 116]]","error":"execution reverted: ret 0x: evm transaction execution failed","revert_reason":""}`
//		parsedMsg, err := cctxerror.ParseCCTXErrorMessage(m)
//		require.NoError(t, err)
//
//		require.Equal(t, "contract call failed when calling EVM with data", parsedMsg.Message)
//		require.Equal(t, "depositAndCall0", parsedMsg.Method)
//		require.Equal(t, "0x733aB8b06DDDEf27Eaa72294B0d7c9cEF7f12db9", parsedMsg.Contract)
//		require.Equal(
//			t,
//			"[{[]0xdFb74337c53141bf912101b0Ee770FA8e2DCB921 1337} 0x13A0c5930C028511Dc02665E7285134B6d11A5f410000000000000000 0xD28D6A0b8189305551a0A8bd247a6ECa9CE781Ca [114 101 118 101114 116]]",
//			parsedMsg.Args,
//		)
//		require.Equal(t, "execution reverted: ret 0x: evm transaction execution failed", parsedMsg.Error)
//		require.Equal(t, "", parsedMsg.RevertReason)
//
//	})
//}

func Test_NewCCTXErrorJsonMessage(t *testing.T) {
	t.Run("Test_NewCCTXErrorJsonMessage", func(t *testing.T) {
		m := cctxerror.NewCCTXErrorJsonMessage("message", errors.New("test", 999110, "error_message"))
		fmt.Println(m)
	})

	t.Run("Unwrap error from NewCCTXErrorJsonMessage", func(t *testing.T) {
		m := cctxerror.NewCCTXErrorJsonMessage("message", errors.New("test", 999111, "error_message"))
		fmt.Println("m1", m)
		errm := errors.New("test", 2, m)
		m2 := cctxerror.NewCCTXErrorJsonMessage("", errm)
		fmt.Println("m2", m2)
	})

}

/*

can''t call a non-contract addres
*/

/*
'{"message":"outbound tx failed to be executed on connected chain","error":"","method":"","contract":"","args":"","revert_reason":""}'

'{"message":"contract call failed when calling EVM with data","error":"execution
      reverted: ret 0x: evm transaction execution failed","method":"depositAndCall0","contract":"0x733aB8b06DDDEf27Eaa72294B0d7c9cEF7f12db9","args":"[{[]
      0xdFb74337c53141bf912101b0Ee770FA8e2DCB921 1337} 0x13A0c5930C028511Dc02665E7285134B6d11A5f4
      10000000000000000 0xD28D6A0b8189305551a0A8bd247a6ECa9CE781Ca [114 101 118 101
      114 116]]","revert_reason":""}
*/
