package helper

import (
	"bytes"
	"embed"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

//go:embed abi/*.json
var abiFiles embed.FS

// loadABI loads an ABI from the embedded JSON files
func loadABI(filename string) (abi.ABI, error) {
	data, err := abiFiles.ReadFile(fmt.Sprintf("abi/%s.json", filename))
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to read ABI file %s: %w", filename, err)
	}

	parsedABI, err := abi.JSON(bytes.NewReader(data))
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to parse ABI for %s: %w", filename, err)
	}

	return parsedABI, nil
}

// Initialize all ABIs
var (
	erc20ABI            abi.ABI
	routerMoudleABI     abi.ABI
	swapMoudleABI       abi.ABI
	calculatorMoudleABI abi.ABI
)

func init() {
	var err error

	// Load ERC20 ABI
	erc20ABI, err = loadABI("erc20")
	if err != nil {
		panic(fmt.Sprintf("Failed to load ERC20 ABI: %v", err))
	}

	// Load calculatorMoudle ABI
	calculatorMoudleABI, err = loadABI("calculatorMoudle")
	if err != nil {
		panic(fmt.Sprintf("Failed to load calculatorMoudle ABI: %v", err))
	}

	// Load routerMoudle ABI
	routerMoudleABI, err = loadABI("routerMoudle")
	if err != nil {
		panic(fmt.Sprintf("Failed to load routerMoudle ABI: %v", err))
	}

	// Load swapMoudle ABI
	swapMoudleABI, err = loadABI("swapMoudle")
	if err != nil {
		panic(fmt.Sprintf("Failed to load swapMoudle ABI: %v", err))
	}
}
