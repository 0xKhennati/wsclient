# ABI JSON Files

This directory contains the ABI (Application Binary Interface) definitions for various smart contracts in JSON format. These files are embedded into the Go binary at compile time for efficient loading and parsing.

## Structure

```
helper/abi/
├── README.md          # This file
├── erc20.json         # Standard ERC-20 token ABI
├── uniswapv2.json     # Uniswap V2 pair contract ABI
├── uniswapv3.json     # Uniswap V3 pool contract ABI
└── aave.json          # Aave V3 pool contract ABI
```

## ABI Files

### erc20.json
Contains the standard ERC-20 token interface functions:
- `balanceOf(address)` - Get token balance
- `transfer(address,uint256)` - Transfer tokens
- `approve(address,uint256)` - Approve spending
- `allowance(address,address)` - Check allowance
- `transferFrom(address,address,uint256)` - Transfer from approved amount
- `name()` - Get token name
- `symbol()` - Get token symbol
- `decimals()` - Get token decimals
- `totalSupply()` - Get total token supply

### uniswapv2.json
Contains Uniswap V2 pair contract functions:
- `getReserves()` - Get pair reserves and timestamp
- `token0()` - Get first token address
- `token1()` - Get second token address
- `swap(uint256,uint256,address,bytes)` - Execute swap
- `price0CumulativeLast()` - Get price cumulative for token0
- `price1CumulativeLast()` - Get price cumulative for token1

### uniswapv3.json
Contains Uniswap V3 pool contract functions:
- `slot0()` - Get pool state (price, tick, etc.)
- `token0()` - Get first token address
- `token1()` - Get second token address
- `fee()` - Get pool fee tier
- `liquidity()` - Get current liquidity
- `swap(address,bool,int256,uint160,bytes)` - Execute swap

### aave.json
Contains Aave V3 pool contract functions:
- `supply(address,uint256,address,uint16)` - Supply assets
- `withdraw(address,uint256,address)` - Withdraw assets
- `borrow(address,uint256,uint256,uint16,address)` - Borrow assets
- `repay(address,uint256,uint256,address)` - Repay debt
- `getUserAccountData(address)` - Get user account information
- `getReserveData(address)` - Get reserve information

## Usage

These JSON files are automatically loaded by the helper package using Go's embed feature. You don't need to manually read or parse these files - the helper functions handle everything for you.

```go
import "github.com/khennati22/wsClient/helper"

// The ABI is automatically loaded and ready to use
data, err := helper.BuildERC20BalanceOfCallData(walletAddress)
```

## Adding New ABIs

To add a new contract ABI:

1. Create a new JSON file in this directory (e.g., `compound.json`)
2. Add the ABI functions in standard JSON format
3. Update `abi_loader.go` to load the new ABI
4. Create corresponding helper functions in a new Go file
5. Add exports to `helper.go` for convenience

### JSON Format

Each ABI file should contain an array of function definitions:

```json
[
  {
    "constant": true,
    "inputs": [{"name": "param1", "type": "address"}],
    "name": "functionName",
    "outputs": [{"name": "result", "type": "uint256"}],
    "type": "function"
  }
]
```

### Properties:
- `constant`: `true` for view/pure functions, `false` for state-changing functions
- `inputs`: Array of input parameters with name and type
- `name`: Function name as it appears in the contract
- `outputs`: Array of return values with name and type
- `type`: Always "function" for function definitions

## Benefits of JSON Separation

1. **Clean Code**: Removes large ABI strings from Go source files
2. **Easy Editing**: JSON files can be easily edited and validated
3. **Version Control**: Better diff tracking for ABI changes
4. **Reusability**: JSON files can be used by other tools
5. **Validation**: JSON format is validated at compile time
6. **Embedded**: Files are embedded in the binary, no runtime file access needed

## File Validation

All JSON files are validated during compilation. If a file contains invalid JSON or ABI format, the build will fail with a clear error message.

## Performance

- ABIs are loaded once during package initialization
- Files are embedded in the binary (no runtime file I/O)
- Parsing is done once and cached for the lifetime of the application
- Function call data generation is very fast after initial loading
