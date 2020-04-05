package message

import "github.com/t10471/bitcoin-coding/basetype"

type X int

type (
	Hoge struct {
		int
		Fuga
	}
	Fuga struct{}
)

//go:generate coding -t MsgBlock
type MsgBlock struct {
	Header   BlockHeader
	TxnCount basetype.VarInt
	Txn      []MsgTx `coding-count:"TxnCount"`
}

//go:generate coding -t BlockHeader
type BlockHeader struct {
	Version    basetype.Int32
	PrevBlock  basetype.Hash
	MerkleRoot basetype.Hash
	Timestamp  basetype.Uint32Time
	Bits       basetype.Uint32
	Nonce      basetype.Uint32
}

//go:generate coding -t MsgTx
type MsgTx struct {
	Hash  basetype.Hash
	Index basetype.Uint32
}
