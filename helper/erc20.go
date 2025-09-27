package helper

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// BuildERC20CallData creates properly wrapped call data for any ERC-20 function
func BuildERC20CallData(functionName string, args ...interface{}) (hexutil.Bytes, error) {
	data, err := erc20ABI.Pack(functionName, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to pack %s: %w", functionName, err)
	}
	return hexutil.Bytes(data), nil
}

// Convenience functions using hexutil.Bytes

// BuildERC20BalanceOfCallData creates wrapped call data for balanceOf
func BuildERC20BalanceOfCallData(owner common.Address) (hexutil.Bytes, error) {
	return BuildERC20CallData("balanceOf", owner)
}

// BuildERC20TransferCallData creates wrapped call data for transfer
func BuildERC20TransferCallData(to common.Address, amount *big.Int) (hexutil.Bytes, error) {
	return BuildERC20CallData("transfer", to, amount)
}

// BuildERC20ApproveCallData creates wrapped call data for approve
func BuildERC20ApproveCallData(spender common.Address, amount *big.Int) (hexutil.Bytes, error) {
	return BuildERC20CallData("approve", spender, amount)
}

// BuildERC20AllowanceCallData creates wrapped call data for allowance
func BuildERC20AllowanceCallData(owner, spender common.Address) (hexutil.Bytes, error) {
	return BuildERC20CallData("allowance", owner, spender)
}

// BuildERC20TransferFromCallData creates wrapped call data for transferFrom
func BuildERC20TransferFromCallData(from, to common.Address, amount *big.Int) (hexutil.Bytes, error) {
	return BuildERC20CallData("transferFrom", from, to, amount)
}

// BuildERC20NameCallData creates wrapped call data for name
func BuildERC20NameCallData() (hexutil.Bytes, error) {
	return BuildERC20CallData("name")
}

// BuildERC20SymbolCallData creates wrapped call data for symbol
func BuildERC20SymbolCallData() (hexutil.Bytes, error) {
	return BuildERC20CallData("symbol")
}

// BuildERC20DecimalsCallData creates wrapped call data for decimals
func BuildERC20DecimalsCallData() (hexutil.Bytes, error) {
	return BuildERC20CallData("decimals")
}

// BuildERC20TotalSupplyCallData creates wrapped call data for totalSupply
func BuildERC20TotalSupplyCallData() (hexutil.Bytes, error) {
	return BuildERC20CallData("totalSupply")
}
