package helper

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// BuildCalculatorCallData creates properly wrapped call data for any calculatorMoudle function
func BuildCalculatorCallData(functionName string, args ...interface{}) (hexutil.Bytes, error) {
	data, err := calculatorMoudleABI.Pack(functionName, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to pack %s: %w", functionName, err)
	}
	return hexutil.Bytes(data), nil
}

// calculatorMoudle convenience functions

func Build_calculatorModule_BalanceCheck(base, pair common.Address, brand uint8) (hexutil.Bytes, error) {
	return BuildCalculatorCallData("BalanceCheck", base, pair, brand)
}
func Build_calculatorModule_GetBaseBalance(base common.Address, pair PairInfoInterface) (hexutil.Bytes, error) {
	return BuildCalculatorCallData("getBaseBalance", base, newPairInfo(pair))
}

func Build_calculatorModule_GetAmountOut(tokenIn common.Address, amountIn *big.Int, pairs ...PairInfoInterface) (hexutil.Bytes, error) {
	return BuildCalculatorCallData("getAmountOut", tokenIn, amountIn, newPairInfoList(pairs...))
}

func Build_calculatorModule_GetAmountOutLoop(tokenIn common.Address, amountIn *big.Int, pairs ...PairInfoInterface) (hexutil.Bytes, error) {
	return BuildCalculatorCallData("getAmountOutLoop", tokenIn, amountIn, newPairInfoList(pairs...))
}

// Build_calculModule_GetMultiPrice to build the data for the getMultiPrice function
func Build_calculModule_GetMultiPrice(tokenIn common.Address, amountsList map[uint8]*big.Int, pairs ...PairInfoInterface) (hexutil.Bytes, error) {

	length := len(amountsList)
	amountsIn := make([]*big.Int, length)
	for i := 0; i < length; i++ {
		amountsIn[i] = amountsList[uint8(i+1)]
	}

	return BuildCalculatorCallData("getMultiPrice", tokenIn, amountsIn, newPairInfoList(pairs...))
}

// decode functions

func decode_calculatorModule(functionName, rawData string) ([]*big.Int, error) {

	// Decode raw data (remove "0x" prefix if present)
	rawBytes, err := hex.DecodeString(strings.TrimPrefix(rawData, "0x"))
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex data: %w", err)
	}

	// Find the "getAmountOut" method
	method, exists := calculatorMoudleABI.Methods[functionName]
	if !exists {
		return nil, fmt.Errorf("method '%s' not found in ABI", functionName)
	}

	// Unpack the result using method outputs
	results, err := method.Outputs.Unpack(rawBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack %s data: %w", functionName, err)
	}

	// The result should be a uint256[] (slice of big.Int)
	uint256Array, ok := results[0].([]*big.Int)
	if !ok {
		return nil, fmt.Errorf("unexpected result type, expected []*big.Int")
	}

	return uint256Array, nil
}

func Decode_calculatorModule_getAmountOut(rawData string) ([]*big.Int, error) {
	return decode_calculatorModule("getAmountOut", rawData)
}

// Decode_calculModule_getAmountOutLoop to decode the data for the getAmountOutLoop function
// the result is a slice of *big.Int
// result[0] is the bestAmountIn
// result[i>0] is the amountOut
func Decode_calculatorModule_getAmountOutLoop(rawData string) ([]*big.Int, error) {
	return decode_calculatorModule("getAmountOutLoop", rawData)
}

func Decode_calculatorModule_GetMultiPrice(rawData string) ([]*big.Int, error) {
	return decode_calculatorModule("getMultiPrice", rawData)
}
