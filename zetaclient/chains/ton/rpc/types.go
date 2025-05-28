package rpc

import "github.com/tidwall/gjson"

type MasterchainInfo struct {
	Seqno    uint64
	RootHash string
	FileHash string
}

/*
	{
	    "@type": "blocks.masterchainInfo",
	    "last": {
	      "@type": "ton.blockIdExt",
	      "workchain": -1,
	      "shard": "-9223372036854775808",
	      "seqno": 48264168,
	      "root_hash": "KMDqzkt/GnyPWviwg6FgDPImNRg95Fmd9WBwREkBG+M=",
	      "file_hash": "eE+Ij4upTsZfzlEh7+n323gNzsl+3fsHtxJysIweU0k="
	    },
	    "state_root_hash": "jDOC8hLoeS+Kx0s4yI9IEqHxlWyImnBBvpP7ljMXF4o=",
	    "init": {
	      "@type": "ton.blockIdExt",
	      "workchain": -1,
	      "shard": "0",
	      "seqno": 0,
	      "root_hash": "F6OpKZKqvqeFp6CQmFomXNMfMj2EnaUSOXN+Mh+wVWk=",
	      "file_hash": "XplPz01CXAps5qeSWUtxcyBfdAo5zVb1N979KLSKD24="
	    },
	    "@extra": "1748432129.6833205:6:0.23348377421037714"
	}
*/
func (m *MasterchainInfo) UnmarshalJSON(data []byte) error {
	result := gjson.GetManyBytes(
		data,
		"last.seqno",
		"last.root_hash",
		"last.file_hash",
	)

	m.Seqno = result[0].Uint()
	m.RootHash = result[1].String()
	m.FileHash = result[2].String()

	return nil
}
