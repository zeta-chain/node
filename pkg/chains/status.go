package chains

import "errors"

// ReceiveStatusFromString returns a ReceiveStatus from a string using in CLI
// 0 for success, 1 for failed
// TODO: remove "receive" naming ans use outbound
// https://github.com/zeta-chain/node/issues/1797
func ReceiveStatusFromString(str string) (ReceiveStatus, error) {
	switch str {
	case "0":
		return ReceiveStatus_success, nil
	case "1":
		return ReceiveStatus_failed, nil
	default:
		return ReceiveStatus(0), errors.New("wrong status, must be 0 for success or 1 for failed")
	}
}
