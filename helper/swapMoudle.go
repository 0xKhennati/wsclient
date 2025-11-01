package helper

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// BuildSwapV2CallData creates properly wrapped call data for any swapMoudle function
func BuildSwapV2CallData(functionName string, args ...interface{}) (hexutil.Bytes, error) {
	data, err := swapMoudleABI.Pack(functionName, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to pack %s: %w", functionName, err)
	}
	return hexutil.Bytes(data), nil
}

// swapMoudle convenience functions

// Build_swap_MultiSwap builds the call data for the multiSwap function
func Build_swap_MultiSwap(tokenIn common.Address, amountIn, miniAmountOut, index *big.Int, pairs ...PairInfoInterface) (hexutil.Bytes, error) {
	return BuildSwapV2CallData("multiSwap", tokenIn, amountIn, miniAmountOut, index, newPairInfoList(pairs...))
}

// Build_swap_simulateSwapAllBalance builds the call data for the simulateSwapAllBalance function
func Build_swap_simulateSwapAllBalance(tokenIn common.Address, pairs ...PairInfoInterface) (hexutil.Bytes, error) {
	return BuildSwapV2CallData("simulateSwapAllBalance", tokenIn, newPairInfoList(pairs...))
}
