package wsClient

import (
	"encoding/json"
	"fmt"
	"sync/atomic"

	// "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Global ID counter for all requests
var globalID int64 = 0

// Request represents a JSON-RPC request
type Request struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int64       `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// Use the official ethereum.CallMsg struct from go-ethereum
// This is the standard struct used for eth_call parameters
// type CallMsg = ethereum.CallMsg
type CallMsg struct {
	From      *common.Address `json:"from,omitempty"`
	To        *common.Address `json:"to,omitempty"`
	Gas       *hexutil.Uint64 `json:"gas,omitempty"`
	GasPrice  *hexutil.Big    `json:"gasPrice,omitempty"`
	Value     *hexutil.Big    `json:"value,omitempty"`
	GasFeeCap *hexutil.Big    `json:"gasFeeCap,omitempty"`
	GasTipCap *hexutil.Big    `json:"gasTipCap,omitempty"`
	Data      hexutil.Bytes   `json:"data,omitempty"`
}

// String returns a string representation of the request
func (r *Request) String() string {
	return fmt.Sprintf("Request{ID: %d, Method: %s, Params: %v}", r.ID, r.Method, r.Params)
}

// GetID returns the ID of the request
func (r *Request) GetID() int64 {
	return r.ID
}

// NewRequest creates a new JSON-RPC request with auto-incrementing global ID
// if id is 0, it will be auto-incremented
func NewRequest(id int64, method string, params interface{}) *Request {
	if id == 0 {
		id = atomic.AddInt64(&globalID, 1)
	} else {
		atomic.AddInt64(&globalID, 1)
	}

	return &Request{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      id,
	}
}

// Response represents a JSON-RPC response
type Response struct {
	ID     int64           `json:"id"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *RPCError       `json:"error,omitempty"`
}

// String returns a string representation of the response
func (r *Response) String() string {
	return fmt.Sprintf("Response{ID: %d, Result: %s, Error: %v}", r.ID, string(r.Result), r.Error)
}

// GetID returns the ID of the response
func (r *Response) GetID() int64 {
	return r.ID
}

// RPCError represents a JSON-RPC error
type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (e *RPCError) Error() string {
	return e.Message
}

// StateOverride represents a state override for a specific address
type StateOverride struct {
	Balance   *string           `json:"balance,omitempty"`   // Override balance
	Nonce     *uint64           `json:"nonce,omitempty"`     // Override nonce
	Code      *string           `json:"code,omitempty"`      // Override code
	State     map[string]string `json:"state,omitempty"`     // Override individual storage slots
	StateDiff map[string]string `json:"stateDiff,omitempty"` // Override storage slots as diff
}
