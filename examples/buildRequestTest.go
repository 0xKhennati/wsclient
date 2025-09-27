package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	wsClient "github.com/khennati22/wsClient"
	"github.com/khennati22/wsClient/helper"
)

// CallMsgForRPC is a wrapper around ethereum.CallMsg that properly handles JSON marshaling for RPC
type CallMsgForRPC struct {
	From     *common.Address `json:"from,omitempty"`
	To       *common.Address `json:"to,omitempty"`
	Gas      *hexutil.Uint64 `json:"gas,omitempty"`
	GasPrice *hexutil.Big    `json:"gasPrice,omitempty"`
	Value    *hexutil.Big    `json:"value,omitempty"`
	Data     hexutil.Bytes   `json:"data,omitempty"`
}

// ToRPCCallMsg converts ethereum.CallMsg to CallMsgForRPC
func ToRPCCallMsg(msg ethereum.CallMsg) CallMsgForRPC {
	rpcMsg := CallMsgForRPC{
		To:   msg.To,
		Data: hexutil.Bytes(msg.Data),
	}

	// Only include fields if they're not zero values
	if msg.From != (common.Address{}) {
		rpcMsg.From = &msg.From
	}
	if msg.Gas != 0 {
		gas := hexutil.Uint64(msg.Gas)
		rpcMsg.Gas = &gas
	}
	if msg.GasPrice != nil && msg.GasPrice.Sign() > 0 {
		rpcMsg.GasPrice = (*hexutil.Big)(msg.GasPrice)
	}
	if msg.Value != nil && msg.Value.Sign() > 0 {
		rpcMsg.Value = (*hexutil.Big)(msg.Value)
	}

	return rpcMsg
}

// Test constants
const (
	WPOL_ADDRESS = "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270"
	MY_WALLET    = "0xaaaaE417c54dF9c3986191B8ad044068ee3Faa64"
	WPOL_SLOT    = 3
	wsUrl        = "ws://65.108.192.189:8548"
)

func main() {
	fmt.Println("=== Testing all buildRequest.go functions ===\n")

	client, err := wsClient.NewClient(wsUrl)
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}
	defer client.Close()

	// Parse addresses
	wpolAddress := common.HexToAddress(WPOL_ADDRESS)
	myWallet := common.HexToAddress(MY_WALLET)

	// Test 1: Build_eth_call_request with balanceOf
	// fmt.Println("1. Testing Build_eth_call_request with WPOL balanceOf")
	// testEthCall(client, wpolAddress, myWallet)
	// Test 2: BuildStateDiff
	// fmt.Println("\n2. Testing BuildStateDiff")
	// testStateDiff(wpolAddress, myWallet)

	// // Test 3: Build_eth_call_request with state overrides
	// fmt.Println("\n3. Testing Build_eth_call_request with state overrides")
	// testEthCallWithStateOverrides(client, wpolAddress, myWallet)
	// return
	// Test 4: Build_GasSoldier_request
	// fmt.Println("\n4. Testing Build_GasSoldier_request")
	// testGasSoldierRequest()

	// // Test 5: Build_GetTargetTx_request
	// fmt.Println("\n5. Testing Build_GetTargetTx_request")
	// testGetTargetTxRequest()

	// Test 6: Build_GetPengingBlockLog_request
	// fmt.Println("\n6. Testing Build_GetPengingBlockLog_request")
	// testGetPendingBlockLogRequest(client)

	// Test 7: Build_GetAccountsData_request
	// fmt.Println("\n7. Testing Build_GetAccountsData_request")
	// testGetAccountsDataRequest(client, myWallet)

	// Test 8: Build_MultiCall_request
	// fmt.Println("\n8. Testing Build_MultiCall_request")
	// testMultiCallRequest(client, wpolAddress, myWallet)

	// Test 9: Build_GetTransactionLog_request
	fmt.Println("\n9. Testing Build_GetTransactionLog_request")
	testGetTransactionLogRequest(client, wpolAddress, myWallet)
	return
	// Test 10: Build_SendRawTransactions_request
	fmt.Println("\n10. Testing Build_SendRawTransactions_request")
	testSendRawTransactionsRequest()

	// Test 11: Build_SendRawTransaction_request
	fmt.Println("\n11. Testing Build_SendRawTransaction_request")
	testSendRawTransactionRequest()

	fmt.Println("\n=== All tests completed ===")
}

func testEthCall(client *wsClient.Client, wpolAddress, myWallet common.Address) {
	// Create balanceOf call data
	balanceOfData, err := helper.BuildERC20BalanceOfCallData(myWallet)
	if err != nil {
		log.Printf("Error creating balanceOf call data: %v", err)
		return
	}

	// Create CallMsg - don't set Gas field to avoid the hexutil.Uint64 issue
	callMsg := wsClient.CallMsg{
		To:   &wpolAddress,
		Data: balanceOfData,
	}

	// Convert to RPC-compatible format
	// rpcCallMsg := ToRPCCallMsg(callMsg)

	request := wsClient.Build_eth_call_request(1, callMsg, nil, "latest")
	// Create request manually with proper formatting

	// Print the request for debugging
	requestJSON, _ := json.MarshalIndent(request, "", "  ")
	fmt.Printf("eth_call request (balanceOf WPOL for wallet):\n%s\n", string(requestJSON))

	var resp wsClient.Response
	if err := client.SendAndReceive(request, &resp); err != nil {
		log.Printf("Failed to get balance: %v", err)
	} else if resp.Error != nil {
		log.Printf("Balance call error: %v", resp.Error)
	} else {
		// Parse the balance result (it's returned as hex string)
		fmt.Printf("✅ Raw balance result: %s\n", resp.Result)

		// Try to decode the hex result to get actual balance
		var balanceHex string
		if err := json.Unmarshal(resp.Result, &balanceHex); err == nil {
			if balance, ok := new(big.Int).SetString(balanceHex[2:], 16); ok {
				// Convert from wei to WPOL (18 decimals)
				divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
				wpolBalance := new(big.Float).SetInt(balance)
				wpolBalance.Quo(wpolBalance, new(big.Float).SetInt(divisor))
				fmt.Printf("✅ Balance: %s WPOL\n", wpolBalance.String())
			}
		}
	}
}

func testStateDiff(wpolAddress, myWallet common.Address) {
	// Create a new balance (1000 WPOL with 18 decimals)
	newBalance := new(big.Int)
	newBalance.SetString("1000000000000000000000", 10) // 1000 * 10^18

	// Build StateDiff
	stateOverrides, err := wsClient.BuildStateDiff(wpolAddress, myWallet, WPOL_SLOT, newBalance)
	if err != nil {
		log.Printf("Error building StateDiff: %v", err)
		return
	}

	// Print the state overrides
	stateJSON, _ := json.MarshalIndent(stateOverrides, "", "  ")
	fmt.Printf("StateDiff for WPOL slot %d:\n%s\n", WPOL_SLOT, string(stateJSON))

}

func testEthCallWithStateOverrides(client *wsClient.Client, wpolAddress, myWallet common.Address) {
	// Create balanceOf call data
	balanceOfData, err := helper.BuildERC20BalanceOfCallData(myWallet)
	if err != nil {
		log.Printf("Error creating balanceOf call data: %v", err)
		return
	}

	// Create CallMsg
	callMsg := ethereum.CallMsg{
		To:   &wpolAddress,
		Data: balanceOfData,
	}

	// Create state overrides with a modified balance
	newBalance := new(big.Int)
	newBalance.SetString("5000000000000000000000", 10) // 5000 * 10^18

	stateOverrides, err := wsClient.BuildStateDiff(wpolAddress, myWallet, WPOL_SLOT, newBalance)
	if err != nil {
		log.Printf("Error building StateDiff: %v", err)
		return
	}

	// Convert to RPC-compatible format
	rpcCallMsg := ToRPCCallMsg(callMsg)

	// Create request manually with proper formatting and state overrides
	request := &wsClient.Request{
		JSONRPC: "2.0",
		ID:      2,
		Method:  "eth_call",
		Params:  []interface{}{rpcCallMsg, "latest", stateOverrides},
	}

	// Print the request
	requestJSON, _ := json.MarshalIndent(request, "", "  ")
	fmt.Printf("eth_call request with state overrides:\n%s\n", string(requestJSON))

	// Execute the call with state overrides
	var resp wsClient.Response
	if err := client.SendAndReceive(request, &resp); err != nil {
		log.Printf("Failed to get balance with state overrides: %v", err)
	} else if resp.Error != nil {
		log.Printf("Balance call with state overrides error: %v", resp.Error)
	} else {
		// Parse the balance result
		fmt.Printf("✅ Raw balance result with state overrides: %s\n", resp.Result)

		var balanceHex string
		if err := json.Unmarshal(resp.Result, &balanceHex); err == nil {
			if balance, ok := new(big.Int).SetString(balanceHex[2:], 16); ok {
				// Convert from wei to WPOL (18 decimals)
				divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
				wpolBalance := new(big.Float).SetInt(balance)
				wpolBalance.Quo(wpolBalance, new(big.Float).SetInt(divisor))
				fmt.Printf("✅ Balance with state overrides: %s WPOL (should be 5000)\n", wpolBalance.String())
			}
		}
	}
}

func testGasSoldierRequest() {
	contractAddress := common.HexToAddress("0x1234567890123456789012345678901234567890")
	senderAddress := common.HexToAddress(MY_WALLET)
	nonceStr := "0x1"
	maxGasStr := "0x5208" // 21000 gas
	key := "test_key"

	request := wsClient.Build_GasSoldier_request(3, contractAddress, senderAddress, nonceStr, maxGasStr, key)

	requestJSON, _ := json.MarshalIndent(request, "", "  ")
	fmt.Printf("GasSoldier request:\n%s\n", string(requestJSON))
}

func testGetTargetTxRequest() {
	// Example arbitrage transaction hash
	arbtHash := common.HexToHash("0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890")

	// Skip contracts list
	skipContracts := []common.Address{
		common.HexToAddress("0x1111111111111111111111111111111111111111"),
		common.HexToAddress("0x2222222222222222222222222222222222222222"),
	}

	request := wsClient.Build_GetTargetTx_request(4, arbtHash, skipContracts)

	requestJSON, _ := json.MarshalIndent(request, "", "  ")
	fmt.Printf("GetTargetTx request:\n%s\n", string(requestJSON))
}

func testGetPendingBlockLogRequest(client *wsClient.Client) {
	request := wsClient.Build_GetPengingBlockLog_request(5)

	requestJSON, _ := json.MarshalIndent(request, "", "  ")
	fmt.Printf("GetPendingBlockLog request:\n%s\n", string(requestJSON))

	var resp wsClient.Response
	if err := client.SendAndReceive(request, &resp); err != nil {
		log.Printf("Failed to get pending block log: %v", err)
	} else if resp.Error != nil {
		log.Printf("Pending block log error: %v", resp.Error)
	} else {
		fmt.Printf("✅ Pending block log: %s\n", resp.Result)
	}
}

func testGetAccountsDataRequest(client *wsClient.Client, myWallet common.Address) {
	// List of addresses to get data for
	addressList := []common.Address{
		myWallet,
		common.HexToAddress("0x1234567890123456789012345678901234567890"),
		common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"),
	}

	request := wsClient.Build_GetAccountsData_request(6, addressList)

	requestJSON, _ := json.MarshalIndent(request, "", "  ")
	fmt.Printf("GetAccountsData request:\n%s\n", string(requestJSON))

	var resp wsClient.Response
	if err := client.SendAndReceive(request, &resp); err != nil {
		log.Printf("Failed to get accounts data: %v", err)
	} else if resp.Error != nil {
		log.Printf("Accounts data error: %v", resp.Error)
	} else {
		fmt.Printf("✅ Accounts data: %s\n", resp.Result)
	}
}

func testMultiCallRequest(client *wsClient.Client, wpolAddress, myWallet common.Address) {
	// Create multiple call messages
	var callMsgs []wsClient.CallMsg

	// 1. balanceOf call
	balanceOfData, err := helper.BuildERC20BalanceOfCallData(myWallet)
	if err != nil {
		log.Printf("Error creating balanceOf call data: %v", err)
		return
	}

	callMsgs = append(callMsgs, wsClient.CallMsg{
		To:   &wpolAddress,
		Data: balanceOfData,
	})

	request := wsClient.Build_MultiCall_request(7, callMsgs, "latest")

	requestJSON, _ := json.MarshalIndent(request, "", "  ")
	fmt.Printf("MultiCall request (balanceOf):\n%s\n", string(requestJSON))

	var resp wsClient.Response
	if err := client.SendAndReceive(request, &resp); err != nil {
		log.Printf("Failed to get balanceOf: %v", err)
	} else if resp.Error != nil {
		log.Printf("BalanceOf error: %v", resp.Error)
	} else {
		fmt.Printf("✅ BalanceOf: %s\n", resp.Result)
	}
	// 2. name call
	nameData, err := helper.BuildERC20NameCallData()
	if err != nil {
		log.Printf("Error creating name call data: %v", err)
		return
	}

	callMsgs = append(callMsgs, wsClient.CallMsg{
		To:   &wpolAddress,
		Data: nameData,
	})

	request = wsClient.Build_MultiCall_request(7, callMsgs, "latest")

	requestJSON, _ = json.MarshalIndent(request, "", "  ")
	fmt.Printf("MultiCall request (name):\n%s\n", string(requestJSON))

	if err := client.SendAndReceive(request, &resp); err != nil {
		log.Printf("Failed to get name: %v", err)
	} else if resp.Error != nil {
		log.Printf("Name error: %v", resp.Error)
	} else {
		fmt.Printf("✅ Name: %s\n", resp.Result)
	}
	// 3. symbol call
	symbolData, err := helper.BuildERC20SymbolCallData()
	if err != nil {
		log.Printf("Error creating symbol call data: %v", err)
		return
	}

	callMsgs = append(callMsgs, wsClient.CallMsg{
		To:   &wpolAddress,
		Data: symbolData,
	})

	request = wsClient.Build_MultiCall_request(7, callMsgs, "latest")

	requestJSON, _ = json.MarshalIndent(request, "", "  ")
	fmt.Printf("MultiCall request (symbol):\n%s\n", string(requestJSON))

	if err := client.SendAndReceive(request, &resp); err != nil {
		log.Printf("Failed to get symbol: %v", err)
	} else if resp.Error != nil {
		log.Printf("Symbol error: %v", resp.Error)
	} else {
		fmt.Printf("✅ Symbol: %s\n", resp.Result)
	}
	// 4. decimals call
	decimalsData, err := helper.BuildERC20DecimalsCallData()
	if err != nil {
		log.Printf("Error creating decimals call data: %v", err)
		return
	}

	callMsgs = append(callMsgs, wsClient.CallMsg{
		To:   &wpolAddress,
		Data: decimalsData,
	})

	request = wsClient.Build_MultiCall_request(7, callMsgs, "latest")

	requestJSON, _ = json.MarshalIndent(request, "", "  ")
	fmt.Printf("MultiCall request (balanceOf, name, symbol, decimals):\n%s\n", string(requestJSON))

	if err := client.SendAndReceive(request, &resp); err != nil {
		log.Printf("Failed to get decimals: %v", err)
	} else if resp.Error != nil {
		log.Printf("Decimals error: %v", resp.Error)
	} else {
		fmt.Printf("✅ Decimals: %s\n", resp.Result)
	}
}

func testGetTransactionLogRequest(client *wsClient.Client, wpolAddress, myWallet common.Address) {
	// Create a transfer call for transaction log testing
	transferAmount := new(big.Int)
	transferAmount.SetString("1000000000000000000", 10) // 1 WPOL

	transferData, err := helper.BuildERC20TransferCallData(
		common.HexToAddress("0x1234567890123456789012345678901234567890"),
		transferAmount,
	)
	if err != nil {
		log.Printf("Error creating transfer call data: %v", err)
		return
	}

	callMsg := wsClient.CallMsg{
		From: &myWallet,
		To:   &wpolAddress,
		Data: transferData,
	}

	request := wsClient.Build_GetTransactionLog_request(8, callMsg, "latest")

	requestJSON, _ := json.MarshalIndent(request, "", "  ")
	fmt.Printf("GetTransactionLog request:\n%s\n", string(requestJSON))

	var resp wsClient.Response
	if err := client.SendAndReceive(request, &resp); err != nil {
		log.Printf("Failed to get transaction log: %v", err)
	} else if resp.Error != nil {
		log.Printf("Transaction log error: %v", resp.Error)
	} else {
		fmt.Printf("✅ Transaction log: %s\n", resp.Result)
	}
}

func testSendRawTransactionsRequest() {
	// Example raw transaction data (these are just examples, not real transactions)
	rawTxs := []hexutil.Bytes{
		hexutil.MustDecode("0xf86c808504a817c800825208941234567890123456789012345678901234567890880de0b6b3a76400008025a01234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdefa01234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
		hexutil.MustDecode("0xf86c018504a817c800825208941234567890123456789012345678901234567890880de0b6b3a76400008025a01234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdefa01234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
	}

	request := wsClient.Build_SendRawTransactions_request(9, rawTxs)

	requestJSON, _ := json.MarshalIndent(request, "", "  ")
	fmt.Printf("SendRawTransactions request:\n%s\n", string(requestJSON))
}

func testSendRawTransactionRequest() {
	// Example raw transaction data (this is just an example, not a real transaction)
	rawTx := hexutil.MustDecode("0xf86c808504a817c800825208941234567890123456789012345678901234567890880de0b6b3a76400008025a01234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdefa01234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")

	request := wsClient.Build_SendRawTransaction_request(10, rawTx)

	requestJSON, _ := json.MarshalIndent(request, "", "  ")
	fmt.Printf("SendRawTransaction request:\n%s\n", string(requestJSON))
}
