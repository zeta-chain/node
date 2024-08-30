package bitcoin

import "errors"

// ErrBitcoinNotEnabled is the error returned when bitcoin is not enabled
var ErrBitcoinNotEnabled = errors.New("bitcoin is not enabled")
