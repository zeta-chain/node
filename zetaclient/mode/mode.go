package mode

import "fmt"

type ClientMode uint8

const (
	StandardMode ClientMode = iota
	DryMode
	ChaosMode
)

func (mode ClientMode) String() string {
	switch mode {
	case StandardMode:
		return "standard-mode"
	case DryMode:
		return "dry-mode"
	case ChaosMode:
		return "chaos-mode"
	default:
		return fmt.Sprintf("invalid mode: %d", mode)
	}
}
