package mode

type ClientMode uint8

const (
	StandardMode ClientMode = iota
	DryMode
	ChaosMode
)
