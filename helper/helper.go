package helper

// This file re-exports all helper functions for easier access
// You can import just "helper" and access all contract functions

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// ContractType represents different contract types for easier identification
type ContractType string

const (
	ERC20Token       ContractType = "erc20"
	RouterModule     ContractType = "router_module"
	SwapModule       ContractType = "swap_module"
	CalculatorModule ContractType = "calculator_module"
)

// ContractHelper provides a unified interface for contract interactions
type ContractHelper struct {
	Address      common.Address
	ContractType ContractType
}

// NewContractHelper creates a new contract helper
func NewContractHelper(address common.Address, contractType ContractType) *ContractHelper {
	return &ContractHelper{
		Address:      address,
		ContractType: contractType,
	}
}

// pairInterface defines the interface for pair-related operations
// It provides methods to access pair information in a standardized way
type PairInfoInterface interface {
	GetPairHex() common.Address
	GetT0Hex() common.Address
	GetT1Hex() common.Address
	GetFee() uint32
	GetBrand() uint8
	GetData() string
	GetCurrency0() *common.Address
	GetCurrency1() *common.Address
}

// PairInfo represents the pairinfo struct
type PairInfo struct {
	PairAddr common.Address
	Token0   common.Address
	Token1   common.Address
	Fee      uint32
	Brand    uint8
	Data     []byte
}

// newPairInfo creates a new PairInfo instance from any type implementing pairInterface
func newPairInfo(pair PairInfoInterface) PairInfo {
	padData := []byte{}
	if pair.GetData() != "" {
		padData, _ = hex.DecodeString(pair.GetData())
	}

	if pair.GetBrand() == 22 {
		currency0 := pair.GetCurrency0()
		currency1 := pair.GetCurrency1()

		return PairInfo{
			PairAddr: pair.GetPairHex(),
			Token0:   *currency0,
			Token1:   *currency1,
			Fee:      uint32(pair.GetFee()),
			Brand:    pair.GetBrand(),
			Data:     padData,
		}
	}

	return PairInfo{
		PairAddr: pair.GetPairHex(),
		Token0:   pair.GetT0Hex(),
		Token1:   pair.GetT1Hex(),
		Fee:      uint32(pair.GetFee()),
		Brand:    pair.GetBrand(),
		Data:     padData,
	}
}

// newPairInfoList creates a list of PairInfo instances from a list of pairInterface implementations
func newPairInfoList(pairs ...PairInfoInterface) []PairInfo {
	pairInfos := make([]PairInfo, len(pairs))
	for i, pair := range pairs {
		pairInfos[i] = newPairInfo(pair)
	}
	return pairInfos
}

func GetNullPairInfo() PairInfo {
	return PairInfo{
		PairAddr: common.Address{},
		Token0:   common.Address{},
		Token1:   common.Address{},
		Fee:      0,
		Brand:    0,
		Data:     []byte{},
	}
}

// LoanPool represents the loanPool struct
// type 0 dodo
// type 1 balancer
// type 2 uniV3
type LoanPool struct {
	Pool          common.Address
	Token         common.Address
	Types         uint8
	IsWrappedAave bool
}

type BalanceCheck struct {
	Latest *big.Int
	Target *big.Int
	End    *big.Int
	Pair   common.Address
	Token  common.Address
	Brand  uint8
}

func NewBalanceCheck(pair, token common.Address, brand uint8) *BalanceCheck {
	return &BalanceCheck{
		Latest: new(big.Int),
		Target: new(big.Int),
		End:    new(big.Int),
		Pair:   pair,
		Token:  token,
		Brand:  brand,
	}
}
func (bc *BalanceCheck) SetLatestBalance(latest string) {
	_, success := bc.Latest.SetString(latest[2:], 16) // Remove "0x" prefix
	if !success {
		fmt.Println("SetLatestBalance: Failed to parse hexadecimal string", latest)
		return
	}
}
func (bc *BalanceCheck) SetTargetBalance(target string) {
	_, success := bc.Target.SetString(target[2:], 16) // Remove "0x" prefix
	if !success {
		fmt.Println("BalanceCheck: Failed to parse hexadecimal string", target)
		return
	}
}
func (bc *BalanceCheck) SetEndBalance(end string) {
	_, success := bc.End.SetString(end[2:], 16) // Remove "0x" prefix
	if !success {
		fmt.Println("BalanceCheck: Failed to parse hexadecimal string", end)
		return
	}
}

func NewLoanPool(pool, token common.Address, types uint8, isWrappedAave bool) *LoanPool {
	return &LoanPool{
		Pool:          pool,
		Token:         token,
		Types:         types,
		IsWrappedAave: isWrappedAave,
	}
}

type Amounts struct {
	Amount1  *big.Int
	Amount2  *big.Int
	Amount3  *big.Int
	Amount4  *big.Int
	Amount5  *big.Int
	Amount6  *big.Int
	Amount7  *big.Int
	Amount8  *big.Int
	Amount9  *big.Int
	Amount10 *big.Int
}

// Function exports for easier access
var (
	// ERC20 Functions
	ERC20BalanceOf    = BuildERC20BalanceOfCallData
	ERC20Transfer     = BuildERC20TransferCallData
	ERC20Approve      = BuildERC20ApproveCallData
	ERC20Allowance    = BuildERC20AllowanceCallData
	ERC20TransferFrom = BuildERC20TransferFromCallData
	ERC20Name         = BuildERC20NameCallData
	ERC20Symbol       = BuildERC20SymbolCallData
	ERC20Decimals     = BuildERC20DecimalsCallData
	ERC20TotalSupply  = BuildERC20TotalSupplyCallData

	// Router Module Functions
	RouterFlashSwapWithLoan         = Build_router_FlashSwapWithLoan_tx
	RouterFlashSwap                 = Build_router_FlashSwap_tx
	RouterSimulateFlashSwapWithLoan = Build_router_SimulateFlashSwapWithLoan
	RouterAtlasFlashSwapWithLoan    = Build_router_AtlasFlashSwapWithLoan
	RouterStartFlashSwapWithLoan    = Build_router_StartFlashSwapWithLoan
	RouterSimulateFlashSwap         = Build_router_SimulateFlashSwap
	RouterStartFlashSwap            = Build_router_StartFlashSwap
	RouterAtlasFlashSwap            = Build_router_AtlasFlashSwap
	RouterAtlasSolverCall           = Build_router_AtlasSolverCall
	RouterAtlasSolverCallSimulation = Build_router_AtlasSolverCallSimulation
	RouterLoanCheck                 = Build_router_LoanCheck

	// Swap Module Functions
	SwapMultiSwap              = Build_swap_MultiSwap
	SwapSimulateSwapAllBalance = Build_swap_simulateSwapAllBalance

	// Calculator Module Functions
	CalculatorBalanceCheck     = Build_calculatorModule_BalanceCheck
	CalculatorGetBaseBalance   = Build_calculatorModule_GetBaseBalance
	CalculatorGetAmountOut     = Build_calculatorModule_GetAmountOut
	CalculatorGetAmountOutLoop = Build_calculatorModule_GetAmountOutLoop
	CalculatorGetMultiPrice    = Build_calculModule_GetMultiPrice

	// Calculator Module Decode Functions
	DecodeCalculatorGetAmountOut     = Decode_calculatorModule_getAmountOut
	DecodeCalculatorGetAmountOutLoop = Decode_calculatorModule_getAmountOutLoop
	DecodeCalculatorGetMultiPrice    = Decode_calculatorModule_GetMultiPrice
)

// Common utility functions

// MaxUint256 returns the maximum uint256 value (useful for approvals)
func MaxUint256() *big.Int {
	max := new(big.Int)
	max.SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)
	return max
}

// WeiToEther converts wei to ether (divides by 10^18)
func WeiToEther(wei *big.Int) *big.Float {
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	ether := new(big.Float).SetInt(wei)
	ether.Quo(ether, new(big.Float).SetInt(divisor))
	return ether
}

// EtherToWei converts ether to wei (multiplies by 10^18)
func EtherToWei(ether *big.Float) *big.Int {
	multiplier := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	wei := new(big.Float).Mul(ether, new(big.Float).SetInt(multiplier))
	result, _ := wei.Int(nil)
	return result
}

// TokenAmountToWei converts token amount to wei based on decimals
func TokenAmountToWei(amount *big.Float, decimals uint8) *big.Int {
	multiplier := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	wei := new(big.Float).Mul(amount, new(big.Float).SetInt(multiplier))
	result, _ := wei.Int(nil)
	return result
}

// WeiToTokenAmount converts wei to token amount based on decimals
func WeiToTokenAmount(wei *big.Int, decimals uint8) *big.Float {
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	amount := new(big.Float).SetInt(wei)
	amount.Quo(amount, new(big.Float).SetInt(divisor))
	return amount
}
