package wasm

type ZetaCoreMsg struct {
	AddToWatchList *AddToWatchList `json:"add_to_watch_list"`
}

type AddToWatchList struct {
	Chain  string `protobuf:"bytes,2,opt,name=chain,proto3" json:"chain,omitempty"`
	Nonce  uint64 `protobuf:"varint,3,opt,name=nonce,proto3" json:"nonce,omitempty"`
	TxHash string `protobuf:"bytes,4,opt,name=txHash,proto3" json:"tx_hash,omitempty"`
}
