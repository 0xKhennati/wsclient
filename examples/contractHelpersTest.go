package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	wsClient "github.com/khennati22/wsClient"
	"github.com/khennati22/wsClient/helper"
)

// Test constants - you can replace these with real contract addresses
const (
	// ERC20 Token (WPOL)
	WPOL_ADDRESS = "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270"
	MY_WALLET    = "0xaaaaE417c54dF9c3986191B8ad044068ee3Faa64"

	// UniswapV2 Pair (example: WETH-USDC pair)
	UNISWAP_V2_PAIR = "0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc"

	// UniswapV3 Pool (example: WETH-USDC 0.3% pool)
	UNISWAP_V3_POOL = "0x8ad599c3A0ff1De082011EFDDc58f1908eb6e6D8"

	// Aave V3 Pool (Polygon)
	AAVE_V3_POOL = "0x794a61358D6845594F94dc1DB02A252b5b4814aD"

	// WebSocket URL
	wsUrl = "ws://65.108.192.189:8548"
)

func main() {
	fmt.Println("=== Testing All Contract Helpers ===\n")

	client, err := wsClient.NewClient(wsUrl)
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}
	defer client.Close()

	// Parse addresses
	wpolAddress := common.HexToAddress(WPOL_ADDRESS)
	myWallet := common.HexToAddress(MY_WALLET)
	uniV2Pair := common.HexToAddress(UNISWAP_V2_PAIR)
	uniV3Pool := common.HexToAddress(UNISWAP_V3_POOL)
	aaveV3Pool := common.HexToAddress(AAVE_V3_POOL)

	// Test 1: ERC20 Functions
	fmt.Println("1. Testing ERC20 Helper Functions")
	testERC20Functions(client, wpolAddress, myWallet)

	// Test 2: UniswapV2 Functions
	fmt.Println("\n2. Testing UniswapV2 Helper Functions")
	testUniswapV2Functions(client, uniV2Pair)

	// Test 3: UniswapV3 Functions
	fmt.Println("\n3. Testing UniswapV3 Helper Functions")
	testUniswapV3Functions(client, uniV3Pool)

	// Test 4: Aave V3 Functions
	fmt.Println("\n4. Testing Aave V3 Helper Functions")
	testAaveV3Functions(client, aaveV3Pool, myWallet)

	// Test 5: Utility Functions
	fmt.Println("\n5. Testing Utility Functions")
	testUtilityFunctions()

	fmt.Println("\n=== All contract helper tests completed ===")
}

func testERC20Functions(client *wsClient.Client, tokenAddress, wallet common.Address) {
	fmt.Printf("Testing ERC20 functions for token: %s\n", tokenAddress.Hex())

	// Test balanceOf
	balanceOfData, err := helper.BuildERC20BalanceOfCallData(wallet)
	if err != nil {
		log.Printf("Error creating balanceOf call data: %v", err)
		return
	}

	callMsg := wsClient.CallMsg{
		To:   &tokenAddress,
		Data: balanceOfData,
	}

	request := wsClient.Build_eth_call_request(1, callMsg, nil, "latest")
	var resp wsClient.Response
	if err := client.SendAndReceive(request, &resp); err != nil {
		log.Printf("Failed to get balance: %v", err)
	} else if resp.Error != nil {
		log.Printf("Balance call error: %v", resp.Error)
	} else {
		fmt.Printf("✅ Balance: %s\n", resp.Result)
	}

	// Test name
	nameData, err := helper.BuildERC20NameCallData()
	if err == nil {
		callMsg.Data = nameData
		request = wsClient.Build_eth_call_request(2, callMsg, nil, "latest")
		if err := client.SendAndReceive(request, &resp); err == nil && resp.Error == nil {
			fmt.Printf("✅ Name: %s\n", resp.Result)
		}
	}

	// Test symbol
	symbolData, err := helper.BuildERC20SymbolCallData()
	if err == nil {
		callMsg.Data = symbolData
		request = wsClient.Build_eth_call_request(3, callMsg, nil, "latest")
		if err := client.SendAndReceive(request, &resp); err == nil && resp.Error == nil {
			fmt.Printf("✅ Symbol: %s\n", resp.Result)
		}
	}

	// Test decimals
	decimalsData, err := helper.BuildERC20DecimalsCallData()
	if err == nil {
		callMsg.Data = decimalsData
		request = wsClient.Build_eth_call_request(4, callMsg, nil, "latest")
		if err := client.SendAndReceive(request, &resp); err == nil && resp.Error == nil {
			fmt.Printf("✅ Decimals: %s\n", resp.Result)
		}
	}
}

func testUniswapV2Functions(client *wsClient.Client, pairAddress common.Address) {
	fmt.Printf("Testing UniswapV2 functions for pair: %s\n", pairAddress.Hex())

	// Test getReserves
	reservesData, err := helper.BuildUniswapV2GetReservesCallData()
	if err != nil {
		log.Printf("Error creating getReserves call data: %v", err)
		return
	}

	callMsg := wsClient.CallMsg{
		To:   &pairAddress,
		Data: reservesData,
	}

	request := wsClient.Build_eth_call_request(10, callMsg, nil, "latest")
	var resp wsClient.Response
	if err := client.SendAndReceive(request, &resp); err != nil {
		log.Printf("Failed to get reserves: %v", err)
	} else if resp.Error != nil {
		log.Printf("Reserves call error: %v", resp.Error)
	} else {
		fmt.Printf("✅ Reserves: %s\n", resp.Result)
	}

	// Test token0
	token0Data, err := helper.BuildUniswapV2Token0CallData()
	if err == nil {
		callMsg.Data = token0Data
		request = wsClient.Build_eth_call_request(11, callMsg, nil, "latest")
		if err := client.SendAndReceive(request, &resp); err == nil && resp.Error == nil {
			fmt.Printf("✅ Token0: %s\n", resp.Result)
		}
	}

	// Test token1
	token1Data, err := helper.BuildUniswapV2Token1CallData()
	if err == nil {
		callMsg.Data = token1Data
		request = wsClient.Build_eth_call_request(12, callMsg, nil, "latest")
		if err := client.SendAndReceive(request, &resp); err == nil && resp.Error == nil {
			fmt.Printf("✅ Token1: %s\n", resp.Result)
		}
	}
}

func testUniswapV3Functions(client *wsClient.Client, poolAddress common.Address) {
	fmt.Printf("Testing UniswapV3 functions for pool: %s\n", poolAddress.Hex())

	// Test slot0
	slot0Data, err := helper.BuildUniswapV3Slot0CallData()
	if err != nil {
		log.Printf("Error creating slot0 call data: %v", err)
		return
	}

	callMsg := wsClient.CallMsg{
		To:   &poolAddress,
		Data: slot0Data,
	}

	request := wsClient.Build_eth_call_request(20, callMsg, nil, "latest")
	var resp wsClient.Response
	if err := client.SendAndReceive(request, &resp); err != nil {
		log.Printf("Failed to get slot0: %v", err)
	} else if resp.Error != nil {
		log.Printf("Slot0 call error: %v", resp.Error)
	} else {
		fmt.Printf("✅ Slot0: %s\n", resp.Result)
	}

	// Test liquidity
	liquidityData, err := helper.BuildUniswapV3LiquidityCallData()
	if err == nil {
		callMsg.Data = liquidityData
		request = wsClient.Build_eth_call_request(21, callMsg, nil, "latest")
		if err := client.SendAndReceive(request, &resp); err == nil && resp.Error == nil {
			fmt.Printf("✅ Liquidity: %s\n", resp.Result)
		}
	}

	// Test fee
	feeData, err := helper.BuildUniswapV3FeeCallData()
	if err == nil {
		callMsg.Data = feeData
		request = wsClient.Build_eth_call_request(22, callMsg, nil, "latest")
		if err := client.SendAndReceive(request, &resp); err == nil && resp.Error == nil {
			fmt.Printf("✅ Fee: %s\n", resp.Result)
		}
	}
}

func testAaveV3Functions(client *wsClient.Client, poolAddress, userAddress common.Address) {
	fmt.Printf("Testing Aave V3 functions for pool: %s\n", poolAddress.Hex())

	// Test getUserAccountData
	userDataCallData, err := helper.BuildAaveV3GetUserAccountDataCallData(userAddress)
	if err != nil {
		log.Printf("Error creating getUserAccountData call data: %v", err)
		return
	}

	callMsg := wsClient.CallMsg{
		To:   &poolAddress,
		Data: userDataCallData,
	}

	request := wsClient.Build_eth_call_request(30, callMsg, nil, "latest")
	var resp wsClient.Response
	if err := client.SendAndReceive(request, &resp); err != nil {
		log.Printf("Failed to get user account data: %v", err)
	} else if resp.Error != nil {
		log.Printf("User account data call error: %v", resp.Error)
	} else {
		fmt.Printf("✅ User Account Data: %s\n", resp.Result)
	}

	// Test getReserveData for WPOL
	wpolAddress := common.HexToAddress(WPOL_ADDRESS)
	reserveDataCallData, err := helper.BuildAaveV3GetReserveDataCallData(wpolAddress)
	if err == nil {
		callMsg.Data = reserveDataCallData
		request = wsClient.Build_eth_call_request(31, callMsg, nil, "latest")
		if err := client.SendAndReceive(request, &resp); err == nil && resp.Error == nil {
			fmt.Printf("✅ Reserve Data (WPOL): %s\n", resp.Result)
		}
	}
}

func testUtilityFunctions() {
	fmt.Println("Testing utility functions:")

	// Test MaxUint256
	maxUint := helper.MaxUint256()
	fmt.Printf("✅ MaxUint256: %s\n", maxUint.String())

	// Test WeiToEther conversion
	weiAmount := big.NewInt(1000000000000000000) // 1 ETH in wei
	etherAmount := helper.WeiToEther(weiAmount)
	fmt.Printf("✅ 1 ETH in wei (%s) = %s ETH\n", weiAmount.String(), etherAmount.String())

	// Test EtherToWei conversion
	etherFloat := big.NewFloat(1.5) // 1.5 ETH
	weiResult := helper.EtherToWei(etherFloat)
	fmt.Printf("✅ 1.5 ETH = %s wei\n", weiResult.String())

	// Test token amount conversions
	tokenAmount := big.NewFloat(100.5) // 100.5 tokens
	weiWithDecimals := helper.TokenAmountToWei(tokenAmount, 18)
	fmt.Printf("✅ 100.5 tokens (18 decimals) = %s wei\n", weiWithDecimals.String())

	backToTokens := helper.WeiToTokenAmount(weiWithDecimals, 18)
	fmt.Printf("✅ %s wei = %s tokens (18 decimals)\n", weiWithDecimals.String(), backToTokens.String())

	// Test with different decimals (USDC has 6 decimals)
	usdcAmount := big.NewFloat(1000.50) // 1000.50 USDC
	usdcWei := helper.TokenAmountToWei(usdcAmount, 6)
	fmt.Printf("✅ 1000.50 USDC (6 decimals) = %s smallest units\n", usdcWei.String())

	backToUSDC := helper.WeiToTokenAmount(usdcWei, 6)
	fmt.Printf("✅ %s smallest units = %s USDC (6 decimals)\n", usdcWei.String(), backToUSDC.String())

	// Demonstrate contract helper creation
	wpolHelper := helper.NewContractHelper(common.HexToAddress(WPOL_ADDRESS), helper.ERC20Token)
	fmt.Printf("✅ Contract helper created for WPOL: %s (type: %s)\n", wpolHelper.Address.Hex(), wpolHelper.ContractType)

	uniV2Helper := helper.NewContractHelper(common.HexToAddress(UNISWAP_V2_PAIR), helper.UniswapV2Pair)
	fmt.Printf("✅ Contract helper created for UniV2 Pair: %s (type: %s)\n", uniV2Helper.Address.Hex(), uniV2Helper.ContractType)
}
