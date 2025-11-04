package errors_test

import (
	"testing"

	"cosmossdk.io/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	cctxerror "github.com/zeta-chain/node/pkg/errors"
)

var (
	sampleErr  = errors.New("test", 9991, "error_cause")
	sampleErr2 = errors.New("test", 9992, "error_cause2")
)

func Test_NewCCTXErrorMessage(t *testing.T) {
	t.Run("TestNewCCTXErrorMessage", func(t *testing.T) {
		m := cctxerror.NewCCTXErrorMessage("message")
		jsonString, err := m.ToJSON()
		require.NoError(t, err)
		require.Equal(
			t,
			`{"type":"internal_error","message":"message"}`,
			jsonString,
		)
	})

}

func Test_ZEvmErrorMessage(t *testing.T) {
	t.Run("TestEvmErrorMessage", func(t *testing.T) {
		contractAddress := "0xE97Ac2CA30D30de65a6FE0Ab20EDC39a623c18df"

		msg := cctxerror.NewZEVMErrorMessage(
			"method",
			common.HexToAddress(contractAddress),
			"args",
			"message",
			sampleErr,
		)
		msg.AddRevertReason("revert_reason")

		s, err := msg.ToJSON()
		require.NoError(t, err)

		require.Equal(
			t,
			`{"type":"contract_call_error","message":"message","error":"error_cause","method":"method","contract":"0xE97Ac2CA30D30de65a6FE0Ab20EDC39a623c18df","args":"args","revert_reason":"revert_reason"}`,
			s,
		)
	})
}

func Test_ParseCCTXErrorMessage(t *testing.T) {
	t.Run("TestParseCCTXErrorMessage", func(t *testing.T) {
		m := `{"message":"contract call failed when calling EVM with data","method":"depositAndCall0","contract":"0x733aB8b06DDDEf27Eaa72294B0d7c9cEF7f12db9","args":"[{[]0xdFb74337c53141bf912101b0Ee770FA8e2DCB921 1337} 0x13A0c5930C028511Dc02665E7285134B6d11A5f410000000000000000 0xD28D6A0b8189305551a0A8bd247a6ECa9CE781Ca [114 101 118 101114 116]]","error":"execution reverted: ret 0x: evm transaction execution failed","revert_reason":""}`
		parsedMsg, err := cctxerror.ParseCCTXErrorMessage(m)
		require.NoError(t, err)

		require.Equal(t, "contract call failed when calling EVM with data", parsedMsg.Message)
		require.Equal(t, "depositAndCall0", parsedMsg.Method)
		require.Equal(t, "0x733aB8b06DDDEf27Eaa72294B0d7c9cEF7f12db9", parsedMsg.Contract)
		require.Equal(
			t,
			"[{[]0xdFb74337c53141bf912101b0Ee770FA8e2DCB921 1337} 0x13A0c5930C028511Dc02665E7285134B6d11A5f410000000000000000 0xD28D6A0b8189305551a0A8bd247a6ECa9CE781Ca [114 101 118 101114 116]]",
			parsedMsg.Args,
		)
		require.Equal(t, "execution reverted: ret 0x: evm transaction execution failed", parsedMsg.Error)
		require.Equal(t, "", parsedMsg.RevertReason)
	})
}

func Test_NewCCTXErrorJsonMessage(t *testing.T) {
	t.Run("Pack error into CCTXErrorMessage", func(t *testing.T) {
		m := cctxerror.NewCCTXErrorJSONMessage("message", sampleErr)
		require.Equal(
			t,
			`{"type":"internal_error","message":"message","error":"error_cause"}`,
			m,
		)
	})

	t.Run("do not repack into CCTXErrorMessage if older error has been formatted already", func(t *testing.T) {
		m := cctxerror.NewCCTXErrorJSONMessage("message", sampleErr)
		require.Equal(
			t,
			`{"type":"internal_error","message":"message","error":"error_cause"}`,
			m,
		)
		errPacked := errors.New("test", 9993, m)
		m2 := cctxerror.NewCCTXErrorJSONMessage("", errPacked)
		require.Equal(
			t,
			`{"type":"internal_error","message":"message","error":"error_cause"}`,
			m2,
		)
	})

	t.Run("unpack json CCTXErrorMessage and wrap other errors into the error field", func(t *testing.T) {
		m := cctxerror.NewCCTXErrorJSONMessage("message", sampleErr)
		require.Equal(
			t,
			`{"type":"internal_error","message":"message","error":"error_cause"}`,
			m,
		)
		errPacked := errors.Wrap(sampleErr2, m)
		errPacked2 := errors.Wrap(sampleErr, errPacked.Error())
		m2 := cctxerror.NewCCTXErrorJSONMessage("", errPacked2)
		require.Equal(
			t,
			`{"type":"internal_error","message":"message","error":"error_cause:error_cause2:error_cause"}`,
			m2,
		)
	})

}

func Test_SplitErrorMessage(t *testing.T) {
	t.Run("SplitErrorMessage", func(t *testing.T) {
		m := `{"message":"message","error":"error_cause","method":"method","contract":"contract","args":"args","revert_reason":"revert_reason"}`
		errorLists := cctxerror.SplitErrorMessage(m)

		require.Len(t, errorLists, 1)
		require.Equal(t, m, errorLists[0])
	})

	t.Run("SplitErrorMessage with prefix and suffix errors", func(t *testing.T) {
		m := `error_cause1:{"message":"message","error":"error_cause","method":"","contract":"","args":"","revert_reason":""}: error_cause2`
		errorLists := cctxerror.SplitErrorMessage(m)

		require.Len(t, errorLists, 3)
		require.Equal(t, "error_cause1", errorLists[0])
		require.Equal(
			t,
			`{"message":"message","error":"error_cause","method":"","contract":"","args":"","revert_reason":""}`,
			errorLists[1],
		)
		require.Equal(t, "error_cause2", errorLists[2])
	})

	t.Run("SplitErrorMessage with wrapped errors", func(t *testing.T) {
		m := cctxerror.NewCCTXErrorJSONMessage("message", sampleErr)
		require.Equal(
			t,
			`{"type":"internal_error","message":"message","error":"error_cause"}`,
			m,
		)
		errPacked := errors.Wrap(sampleErr2, m)
		errPacked2 := errors.Wrap(sampleErr, errPacked.Error())

		errorLists := cctxerror.SplitErrorMessage(errPacked2.Error())
		require.Len(t, errorLists, 3)
		require.Equal(
			t,
			`{"type":"internal_error","message":"message","error":"error_cause"}`,
			errorLists[0],
		)
		require.Equal(t, "error_cause2", errorLists[1])
		require.Equal(t, "error_cause", errorLists[2])
	})
}

func TestCCTXErrorMessage_WrapError(t *testing.T) {
	t.Run("WrapError when error field exists", func(t *testing.T) {
		contractAddress := "0xE97Ac2CA30D30de65a6FE0Ab20EDC39a623c18df"

		msg := cctxerror.NewZEVMErrorMessage(
			"method",
			common.HexToAddress(contractAddress),
			"args",
			"message",
			sampleErr,
		)
		msg.AddRevertReason("revert_reason")
		msg.WrapError("test_error")

		require.Equal(t, "error_cause:test_error", msg.Error)
	})

	t.Run("WrapError when error field does not exist", func(t *testing.T) {
		contractAddress := "0xE97Ac2CA30D30de65a6FE0Ab20EDC39a623c18df"

		msg := cctxerror.NewZEVMErrorMessage(
			"method",
			common.HexToAddress(contractAddress),
			"args",
			"message",
			nil,
		)
		msg.AddRevertReason("revert_reason")
		msg.WrapError("test_error")

		require.Equal(t, "test_error", msg.Error)
	})
}
