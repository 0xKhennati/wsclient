package helper

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// BuildRouterCallData creates properly wrapped call data for any routerMoudle function
func BuildRouterCallData(functionName string, args ...interface{}) (hexutil.Bytes, error) {
	data, err := routerMoudleABI.Pack(functionName, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to pack %s: %w", functionName, err)
	}
	return hexutil.Bytes(data), nil
}

// routerMoudle convenience functions

// Build_router_FlashSwapWithLoan_tx builds the call data for the FlashSwapWithLoan_tx function
func Build_router_FlashSwapWithLoan_tx(balanceCheck *BalanceCheck, loanPool *LoanPool, tokenIn common.Address, amountIn *big.Int, pairs ...PairInfoInterface) (hexutil.Bytes, error) {
	return BuildRouterCallData("FlashSwapWithLoan_tx", true, *balanceCheck, *loanPool, tokenIn, amountIn, newPairInfoList(pairs...))
}

// Build_router_FlashSwap_tx builds the call data for the FlashSwap_tx function
func Build_router_FlashSwap_tx(balanceCheck *BalanceCheck, tokenIn common.Address, amountIn *big.Int, pairs ...PairInfoInterface) (hexutil.Bytes, error) {
	return BuildRouterCallData("FlashSwap_tx", true, *balanceCheck, tokenIn, amountIn, newPairInfoList(pairs...))
}

// Build_router_SimulateFlashSwapWithLoan builds the call data for the simulateFlashSwapWithLoan function
func Build_router_SimulateFlashSwapWithLoan(loanPool *LoanPool, tokenIn common.Address, amountIn *big.Int, pairs []PairInfoInterface, maticPairs []PairInfoInterface) (hexutil.Bytes, error) {
	return BuildRouterCallData("simulateFlashSwapWithLoan", true, *loanPool, tokenIn, amountIn, newPairInfoList(pairs...), newPairInfoList(maticPairs...))
}

// Build_router_AtlasFlashSwapWithLoan builds the call data for the atlasFlashSwapWithLoan function
func Build_router_AtlasFlashSwapWithLoan(loanPool *LoanPool, tokenIn common.Address, amountIn *big.Int, bidAmount *big.Int, pairs []PairInfoInterface, maticPairs []PairInfoInterface) (hexutil.Bytes, error) {
	return BuildRouterCallData("atlasFlashSwapWithLoan", false, *loanPool, tokenIn, amountIn, bidAmount, newPairInfoList(pairs...), newPairInfoList(maticPairs...))
}

// Build_router_StartFlashSwapWithLoan builds the call data for the startFlashSwapWithLoan function
func Build_router_StartFlashSwapWithLoan(loanPool *LoanPool, tokenIn common.Address, amountIn *big.Int, pairs ...PairInfoInterface) (hexutil.Bytes, error) {
	return BuildRouterCallData("startFlashSwapWithLoan", true, *loanPool, tokenIn, amountIn, newPairInfoList(pairs...))
}

// Build_router_SimulateFlashSwap builds the call data for the simulateFlashSwap function
func Build_router_SimulateFlashSwap(tokenIn common.Address, amountIn *big.Int, pairs []PairInfoInterface, maticPairs []PairInfoInterface) (hexutil.Bytes, error) {
	return BuildRouterCallData("simulateFlashSwap", true, tokenIn, amountIn, newPairInfoList(pairs...), newPairInfoList(maticPairs...))
}

// Build_router_StartFlashSwap builds the call data for the startFlashSwap function
func Build_router_StartFlashSwap(tokenIn common.Address, amountIn *big.Int, pairs ...PairInfoInterface) (hexutil.Bytes, error) {
	return BuildRouterCallData("startFlashSwap", true, tokenIn, amountIn, newPairInfoList(pairs...))
}

// Build_router_AtlasFlashSwap builds the call data for the atlasFlashSwap function
func Build_router_AtlasFlashSwap(tokenIn common.Address, amountIn *big.Int, bidAmount *big.Int, pairs []PairInfoInterface, maticPairs []PairInfoInterface) (hexutil.Bytes, error) {
	return BuildRouterCallData("atlasFlashSwap", false, tokenIn, amountIn, bidAmount, newPairInfoList(pairs...), newPairInfoList(maticPairs...))
}

// Build_router_AtlasSolverCall builds the call data for the atlasSolverCall function
func Build_router_AtlasSolverCall(solverOpFrom common.Address, executionEnvironment common.Address, bidToken common.Address, bidAmount *big.Int, solverOpData []byte, extraData []byte) (hexutil.Bytes, error) {
	return BuildRouterCallData("atlasSolverCall", solverOpFrom, executionEnvironment, bidToken, bidAmount, solverOpData, extraData)
}

// Build_router_AtlasSolverCallSimulation builds the call data for the atlasSolverCallSimulation function
func Build_router_AtlasSolverCallSimulation(bidAmount *big.Int, solverOpData []byte) (hexutil.Bytes, error) {
	return BuildRouterCallData("atlasSolverCallSimulation", bidAmount, solverOpData)
}

// function LoanCheck(uint8 types, address pool, address tokenIn, uint256 amount ) external  {
func Build_router_LoanCheck(loanPool *LoanPool, amountIn *big.Int) (hexutil.Bytes, error) {
	return BuildRouterCallData("LoanCheck", loanPool.Types, loanPool.Pool, loanPool.Token, amountIn)

}
