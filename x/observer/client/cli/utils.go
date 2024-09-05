package cli

import (
	"gitlab.com/thorchain/tss/go-tss/blame"

	"github.com/zeta-chain/node/x/observer/types"
)

func ConvertNodes(n []blame.Node) (nodes []*types.Node) {
	for _, node := range n {
		var entry types.Node
		entry.PubKey = node.Pubkey
		entry.BlameSignature = node.BlameSignature
		entry.BlameData = node.BlameData

		nodes = append(nodes, &entry)
	}
	return
}
