package model

type JSONRpcRequest struct {
	JsonRPC string
	ID      string
	Method  string
	Params  string
}
