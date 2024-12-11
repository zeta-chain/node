package types

import (
	"errors"
)

var (
	ErrOperationalFlagsRestartHeightNegative = errors.New("restart height cannot be negative")
)

func (f *OperationalFlags) Validate() error {
	if f.RestartHeight < 0 {
		return ErrOperationalFlagsRestartHeightNegative
	}
	return nil
}
