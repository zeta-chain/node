package types

import "github.com/zeta-chain/go-tss/blame"

func ConvertNodes(n []blame.Node) (nodes []*Node) {
	for _, node := range n {
		var entry Node
		entry.PubKey = node.Pubkey
		entry.BlameSignature = node.BlameSignature
		entry.BlameData = node.BlameData

		nodes = append(nodes, &entry)
	}
	return
}
