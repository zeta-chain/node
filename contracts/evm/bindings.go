//go:generate abigen --abi Connector.abi --pkg evm --type Connector --out Connector.go

package evm

var _ = Connector{}
