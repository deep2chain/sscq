package utils

import common "github.com/ethereum/go-ethereum/common"

func StringToAddress(s string) common.Address { return common.BytesToAddress([]byte(s)) }

func StringToHash(s string) common.Hash { return common.BytesToHash([]byte(s)) }
