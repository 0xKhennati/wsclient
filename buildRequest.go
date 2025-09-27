package wsClient

import (
	"encoding/hex"
	"fmt"
	"math/big"

	// "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// BuildEthCall creates an eth_call request with CallMsg and state overrides
func Build_eth_call_request(id int64, callMsg CallMsg, stateOverrides map[common.Address]StateOverride, blockNumber ...any) *Request {
	params := []interface{}{callMsg, buildBlockNumber(blockNumber)}
	if stateOverrides != nil {
		params = []interface{}{callMsg, buildBlockNumber(blockNumber), stateOverrides}
	}
	return NewRequest(id, "eth_call", params)
}

// BuildStateDiff builds StateDiff to update holder balance in token
func BuildStateDiff(tokenContract, holder common.Address, slot int64, newBalance *big.Int) (map[common.Address]StateOverride, error) {
	stateOverrides := make(map[common.Address]StateOverride)
	if newBalance == nil {
		return stateOverrides, fmt.Errorf("newBalance is nil")
	}

	// Define ABI type for address and uint256
	addressType, err := abi.NewType("address", "", nil)
	if err != nil {
		return stateOverrides, fmt.Errorf("failed to create address type: %w", err)
	}
	uintType, err := abi.NewType("uint256", "", nil)
	if err != nil {
		return stateOverrides, fmt.Errorf("failed to create uint256 type: %w", err)
	}

	// Pack the address and slot number according to ABI encoding
	encodedData, err := abi.Arguments{{Type: addressType}, {Type: uintType}}.Pack(holder, big.NewInt(slot))
	if err != nil {
		return stateOverrides, fmt.Errorf("failed to pack data: %w", err)
	}

	// Hash the encoded data using Keccak256
	m_slot := crypto.Keccak256(encodedData)

	// Convert balance to 32-byte padded hex
	newBalanceBytes := common.LeftPadBytes(newBalance.Bytes(), 32)

	stateOverrides[tokenContract] = StateOverride{
		StateDiff: map[string]string{
			"0x" + hex.EncodeToString(m_slot): "0x" + hex.EncodeToString(newBalanceBytes),
		},
	}

	return stateOverrides, nil
}

// castom function api

// GasSoldier to get
func Build_GasSoldier_request(id int64, contractAddress, senderAddress common.Address, nonceStr, maxGasStr, key string) *Request {
	// params := []interface{}{contractAddress, senderAddress, nonceStr, maxGasStr, key}
	return NewRequest(id, "eth_gasSoldier", []interface{}{contractAddress, senderAddress, nonceStr, maxGasStr, key})
}

// GetTargetTx to get target tx for arbt transaction, the response like this:
//
//	type GasTargetTxResult struct {
//		Logs     []*types.Log    `json:"logs"`
//		TargetTx *RPCTransaction `json:"targetTx"`
//	}
func Build_GetTargetTx_request(id int64, arbtHash common.Hash, skipContract []common.Address) *Request {
	// params := []interface{}{arbtHash, skipContract}
	return NewRequest(id, "eth_getTargetTx", []interface{}{arbtHash, skipContract})
}

// GetPengingBlockLog to get logs for pending block
//
// the response like this:
// []*types.Log
func Build_GetPengingBlockLog_request(id int64) *Request {
	return NewRequest(id, "eth_getPengingBlockLog", []interface{}{"pending"})
}

// GetAccountsData to get accounts data
// the response like this:
//
//	type accountData struct {
//		Account      common.Address  `json:"account"`
//		Nonce        *hexutil.Uint64 `json:"nonce"`
//		PendingNonce *hexutil.Uint64 `json:"pendingNonce"`
//		Balance      *hexutil.Big    `json:"balance"`
//	}
func Build_GetAccountsData_request(id int64, addressList []common.Address) *Request {
	return NewRequest(id, "eth_getAccountsData", []interface{}{addressList})
}

// MultiCall to call multiple contracts in one request, and return the result of last callMsg
func Build_MultiCall_request(id int64, args []CallMsg, blockNumber ...any) *Request {
	return NewRequest(id, "eth_multiCall", []interface{}{args, buildBlockNumber(blockNumber)})
}

// GetTransactionLog to get logs for transaction still not mined
// we can apply this transaction in latest or in pendnig state
// if use the pendnig state, only use the cash pending state
// the response like this:
// []*types.Log
func Build_GetTransactionLog_request(id int64, args CallMsg, blockNumber ...any) *Request {
	return NewRequest(id, "eth_getTransactionLog", []interface{}{args, buildBlockNumber(blockNumber)})
}

// SendRawTransactions to send multi raws of transactions direct to the conneted peers.
func Build_SendRawTransactions_request(id int64, args []hexutil.Bytes) *Request {
	return NewRequest(id, "eth_sendRawTransactions", []interface{}{args})
}

// SendRawTransaction to send one raw of transaction and return the hash of the transaction.
func Build_SendRawTransaction_request(id int64, args hexutil.Bytes) *Request {
	return NewRequest(id, "eth_sendRawTransaction", []interface{}{args})
}
