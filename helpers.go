package wsClient

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// Helper functions to build common Ethereum JSON-RPC requests

// buildBlockNumber builds the block number for the request
// If blockNumber slice is empty or nil, defaults to "latest"
func buildBlockNumber(blockNumber []any) string {
	// Handle empty slice case - default to "latest"
	if len(blockNumber) == 0 {
		return "latest"
	}

	// Handle nil case - default to "latest"
	if blockNumber[0] == nil {
		return "latest"
	}

	switch block := blockNumber[0].(type) {
	case string:
		return block
	case int64, uint64, int32, uint32, int16, uint16, int8, uint8, int, uint:
		return fmt.Sprintf("0x%x", block)
	}
	panic(fmt.Sprintf("invalid block number type: %T", blockNumber[0]))
}

// BuildGetTransactionCount creates a request to get transaction count for an address
// If blockNumber is not provided, defaults to "latest"
func BuildRequestGetTransactionCount(address common.Address, blockNumber ...any) *Request {
	params := []interface{}{address.Hex(), buildBlockNumber(blockNumber)}
	return NewRequest(0, "eth_getTransactionCount", params)
}

// BuildGetBalance creates a request to get balance for an address
// If blockNumber is not provided, defaults to "latest"
func BuildRequestGetBalance(address common.Address, blockNumber ...any) *Request {
	params := []interface{}{address.Hex(), buildBlockNumber(blockNumber)}
	return NewRequest(0, "eth_getBalance", params)
}

// BuildGetBlockNumber creates a request to get the latest block number
func BuildRequestGetBlockNumber() *Request {
	return NewRequest(0, "eth_blockNumber", nil)
}

// BuildGetBlockByNumber creates a request to get block by number
func BuildRequestGetBlockByNumber(number int64, fullTx bool) *Request {
	blockNum := fmt.Sprintf("0x%x", number)
	params := []interface{}{blockNum, fullTx} // true for full transaction objects
	return NewRequest(0, "eth_getBlockByNumber", params)
}

// BuildGetTransactionByHash creates a request to get transaction by hash
func BuildRequestGetTransactionByHash(hash common.Hash) *Request {
	params := []interface{}{hash.Hex()}
	return NewRequest(0, "eth_getTransactionByHash", params)
}

// BuildGetTransactionReceipt creates a request to get transaction receipt
func BuildRequestGetTransactionReceipt(hash common.Hash) *Request {
	params := []interface{}{hash.Hex()}
	return NewRequest(0, "eth_getTransactionReceipt", params)
}

// // BuildStateOverride creates a new StateOverride struct
// func BuildStateOverride() *StateOverride {
// 	return &StateOverride{
// 		StateDiff: make(map[string]string),
// 	}
// }
