package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/khennati22/wsClient/helper"
)

func main() {
	fmt.Println("=== Testing ABI Loading and Call Data Generation ===\n")

	// Test ERC20 functions
	fmt.Println("1. Testing ERC20 ABI:")
	testERC20ABI()
	return
	// Test UniswapV2 functions
	fmt.Println("\n2. Testing UniswapV2 ABI:")
	testUniswapV2ABI()

	// Test UniswapV3 functions
	fmt.Println("\n3. Testing UniswapV3 ABI:")
	testUniswapV3ABI()

	// Test Aave functions
	fmt.Println("\n4. Testing Aave V3 ABI:")
	testAaveABI()

	// Test utility functions
	fmt.Println("\n5. Testing Utility Functions:")
	testUtilityFunctions()

	fmt.Println("\n=== All ABI tests passed! ===")
}

func testERC20ABI() {
	wallet := common.HexToAddress("0xaaaaE417c54dF9c3986191B8ad044068ee3Faa64")

	// Test balanceOf
	data, err := helper.BuildERC20BalanceOfCallData(wallet)
	if err != nil {
		log.Fatalf("Failed to build balanceOf call data: %v", err)
	}
	fmt.Println("✅ balanceOf call data:", data)

	// Test name
	data, err = helper.BuildERC20NameCallData()
	if err != nil {
		log.Fatalf("Failed to build name call data: %v", err)
	}
	fmt.Println("✅ name call data:", data)

	// Test transfer
	amount := big.NewInt(1000000000000000000) // 1 token
	to := common.HexToAddress("0x1234567890123456789012345678901234567890")
	data, err = helper.BuildERC20TransferCallData(to, amount)
	if err != nil {
		log.Fatalf("Failed to build transfer call data: %v", err)
	}
	fmt.Println("✅ transfer call data:", data)
}

func testUniswapV2ABI() {
	// Test getReserves
	data, err := helper.BuildUniswapV2GetReservesCallData()
	if err != nil {
		log.Fatalf("Failed to build getReserves call data: %v", err)
	}
	fmt.Printf("✅ getReserves call data: 0x%x\n", data)

	// Test token0
	data, err = helper.BuildUniswapV2Token0CallData()
	if err != nil {
		log.Fatalf("Failed to build token0 call data: %v", err)
	}
	fmt.Printf("✅ token0 call data: 0x%x\n", data)

	// Test swap
	amount0Out := big.NewInt(0)
	amount1Out := big.NewInt(1000000)
	to := common.HexToAddress("0x1234567890123456789012345678901234567890")
	swapData := []byte{}
	data, err = helper.BuildUniswapV2SwapCallData(amount0Out, amount1Out, to, swapData)
	if err != nil {
		log.Fatalf("Failed to build swap call data: %v", err)
	}
	fmt.Printf("✅ swap call data: 0x%x\n", data)
}

func testUniswapV3ABI() {
	// Test slot0
	data, err := helper.BuildUniswapV3Slot0CallData()
	if err != nil {
		log.Fatalf("Failed to build slot0 call data: %v", err)
	}
	fmt.Printf("✅ slot0 call data: 0x%x\n", data)

	// Test liquidity
	data, err = helper.BuildUniswapV3LiquidityCallData()
	if err != nil {
		log.Fatalf("Failed to build liquidity call data: %v", err)
	}
	fmt.Printf("✅ liquidity call data: 0x%x\n", data)

	// Test fee
	data, err = helper.BuildUniswapV3FeeCallData()
	if err != nil {
		log.Fatalf("Failed to build fee call data: %v", err)
	}
	fmt.Printf("✅ fee call data: 0x%x\n", data)
}

func testAaveABI() {
	asset := common.HexToAddress("0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270")
	user := common.HexToAddress("0xaaaaE417c54dF9c3986191B8ad044068ee3Faa64")
	amount := big.NewInt(1000000000000000000)

	// Test getUserAccountData
	data, err := helper.BuildAaveV3GetUserAccountDataCallData(user)
	if err != nil {
		log.Fatalf("Failed to build getUserAccountData call data: %v", err)
	}
	fmt.Printf("✅ getUserAccountData call data: 0x%x\n", data)

	// Test getReserveData
	data, err = helper.BuildAaveV3GetReserveDataCallData(asset)
	if err != nil {
		log.Fatalf("Failed to build getReserveData call data: %v", err)
	}
	fmt.Printf("✅ getReserveData call data: 0x%x\n", data)

	// Test supply
	data, err = helper.BuildAaveV3SupplyCallData(asset, amount, user, 0)
	if err != nil {
		log.Fatalf("Failed to build supply call data: %v", err)
	}
	fmt.Printf("✅ supply call data: 0x%x\n", data)
}

func testUtilityFunctions() {
	// Test MaxUint256
	maxUint := helper.MaxUint256()
	fmt.Printf("✅ MaxUint256: %s\n", maxUint.String())

	// Test WeiToEther conversion
	weiAmount := big.NewInt(1000000000000000000) // 1 ETH in wei
	etherAmount := helper.WeiToEther(weiAmount)
	fmt.Printf("✅ 1 ETH in wei (%s) = %s ETH\n", weiAmount.String(), etherAmount.String())

	// Test token amount conversions
	tokenAmount := big.NewFloat(100.5) // 100.5 tokens
	weiWithDecimals := helper.TokenAmountToWei(tokenAmount, 18)
	fmt.Printf("✅ 100.5 tokens (18 decimals) = %s wei\n", weiWithDecimals.String())

	// Test contract helper
	wpolHelper := helper.NewContractHelper(common.HexToAddress("0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270"), helper.ERC20Token)
	fmt.Printf("✅ Contract helper created: %s (type: %s)\n", wpolHelper.Address.Hex(), wpolHelper.ContractType)
}
