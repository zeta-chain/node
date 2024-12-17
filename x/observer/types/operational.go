package types

func (f *OperationalFlags) Validate() error {
	if f.RestartHeight < 0 {
		return ErrOperationalFlagsRestartHeightNegative
	}
	return nil
}
