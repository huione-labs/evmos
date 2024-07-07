// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)
package types

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"

	errorsmod "cosmossdk.io/errors"
)

var (
	EvmChainID_Testnet int64 = 9000
	EvmChainID_Mainnet int64 = 9001
)

var (
	regexChainID         = `[a-z]{1,}`
	regexEIP155Separator = `_{1}`
	regexEIP155          = `[1-9][0-9]*`
	regexEpochSeparator  = `-{1}`
	regexEpoch           = `[1-9][0-9]*`
	evmosChainID         = regexp.MustCompile(fmt.Sprintf(`^(%s)%s(%s)%s(%s)$`,
		regexChainID,
		regexEIP155Separator,
		regexEIP155,
		regexEpochSeparator,
		regexEpoch))
)

// IsValidChainID returns false if the given chain identifier is incorrectly formatted.
func IsValidChainID(chainID string) bool {
	if len(chainID) > 48 {
		return false
	}

	return evmosChainID.MatchString(chainID)
}

// ParseChainID parses a string chain identifier's epoch to an Ethereum-compatible
// chain-id in *big.Int format. The function returns an error if the chain-id has an invalid format
func ParseChainID(chainID string) (*big.Int, error) {
	chainID = strings.TrimSpace(chainID)
	if len(chainID) > 48 {
		return nil, errorsmod.Wrapf(ErrInvalidChainID, "chain-id '%s' cannot exceed 48 chars", chainID)
	}

	matches := evmosChainID.FindStringSubmatch(chainID)
	if matches == nil || len(matches) != 4 || matches[1] == "" {
		return nil, errorsmod.Wrapf(ErrInvalidChainID, "%s: %v", chainID, matches)
	}

	// verify that the chain-id entered is a base 10 integer
	chainIDInt, ok := new(big.Int).SetString(matches[2], 10)
	if !ok {
		return nil, errorsmod.Wrapf(ErrInvalidChainID, "epoch %s must be base-10 integer format", matches[2])
	}

	return chainIDInt, nil
}

func CheckEvmChainID(chainID *big.Int) {
	if !(chainID.Cmp(big.NewInt(EvmChainID_Testnet)) == 0 || chainID.Cmp(big.NewInt(EvmChainID_Mainnet)) == 0) {
		panic(fmt.Sprintf("EVM only supports chain identifiers (%v or %v)", EvmChainID_Testnet, EvmChainID_Mainnet))
	}
}

func SetEvmChainIDs(testnet, mainnet int64) {
	EvmChainID_Testnet = testnet
	EvmChainID_Mainnet = mainnet
}

func IsMainnet(chainID string) bool {
	cid, err := ParseChainID(chainID)
	if err != nil {
		panic(err.Error())
	}
	if cid.Int64() == EvmChainID_Mainnet {
		return true
	}
	return false
}
