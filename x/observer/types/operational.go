package types

import (
	"time"

	cosmoserrors "cosmossdk.io/errors"
	"golang.org/x/mod/semver"
)

const (
	signerBlockTimeOffsetLimit = time.Second * 10
)

func (f *OperationalFlags) Validate() error {
	if f.RestartHeight < 0 {
		return ErrOperationalFlagsRestartHeightNegative
	}
	if f.SignerBlockTimeOffset != nil {
		signerBlockTimeOffset := *f.SignerBlockTimeOffset
		if signerBlockTimeOffset < 0 {
			return ErrOperationalFlagsSignerBlockTimeOffsetNegative
		}
		if signerBlockTimeOffset > signerBlockTimeOffsetLimit {
			return cosmoserrors.Wrapf(ErrOperationalFlagsSignerBlockTimeOffsetLimit, "(%s)", signerBlockTimeOffset)
		}
	}
	if f.MinimumVersion != "" && !semver.IsValid(f.MinimumVersion) {
		return ErrOperationalFlagsInvalidMinimumVersion
	}
	return nil
}
