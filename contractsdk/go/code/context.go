package code

import (
	"math/big"
)

// Context is the context in which the contract runs
type Context interface {
	Args() map[string]interface{}
	// TxID获取当前合约的交易id, Query调用里面为空
	TxID() []byte
	// Caller获取发起合约调用的用户地址，Query调用里面caller为空
	Caller() string

	PutObject(key []byte, value []byte) error
	GetObject(key []byte) ([]byte, error)
	DeleteObject(key []byte) error
	NewIterator(start, limit []byte) Iterator

	QueryTx(txid []byte) (*TxStatus, error)
	QueryBlock(blockid []byte) (*Block, error)
	Transfer(to string, amount *big.Int) error
	Call(module, method string, args map[string]interface{}) (*Response, error)
}

// Iterator iterates over key/value pairs in key order
type Iterator interface {
	Key() []byte
	Value() []byte
	Next() bool
	Error() error
	// Iterator 必须在使用完毕后关闭
	Close()
}

// PrefixRange returns key range that satisfy the given prefix
func PrefixRange(prefix []byte) ([]byte, []byte) {
	var limit []byte
	for i := len(prefix) - 1; i >= 0; i-- {
		c := prefix[i]
		if c < 0xff {
			limit = make([]byte, i+1)
			copy(limit, prefix)
			limit[i] = c + 1
			break
		}
	}
	return prefix, limit
}
