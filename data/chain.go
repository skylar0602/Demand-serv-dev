package data

import (
	"fmt"
	"strings"
)

var ChainIDMap = make(map[string]int)

func GetChainIdByName(chainName string) (int, error) {
	chainId, ok := ChainIDMap[strings.ToLower(chainName)]
	if !ok {
		return 0, fmt.Errorf("chain %s not found", chainName)
	}
	return chainId, nil
}
