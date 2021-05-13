// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"math/big"
	"time"

	"github.com/deep2chain/sscq/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/deep2chain/sscq/evm/vm"
)

// ChainContext supports retrieving headers and consensus parameters from the
// current blockchain to be used during transaction processing.
type ChainContext interface {
	// Engine retrieves the chain's consensus engine.
	//Engine() consensus.Engine

	// GetHeader returns the hash corresponding to their hash.
	GetHeader(common.Hash, uint64) *types.Header
}

type IMessage interface {
	FromAddress() common.Address
}

type FooChainContext struct {
	blockTime time.Time
}

func (self FooChainContext) GetHeader(hash common.Hash, number uint64) *types.Header {

	return &types.Header{
		Difficulty: big.NewInt(1), // NOTE: DO NOT use difficulty to generate random number in contract  !! 2020-12-09 yqq
		Number:     big.NewInt(int64(number)),
		GasLimit:   0,
		GasUsed:    0,
		// Time:       big.NewInt(time.Now().Unix()).Uint64(), // fix issue #15 yqq 2020-12-09
		Time:  uint64(self.blockTime.Unix()),
		Extra: nil,
	}
}

// NewEVMContext creates a new context for use in the EVM.
func NewEVMContext(msg IMessage, author *common.Address, height uint64, blockTime time.Time) vm.Context {
	// If we don't have an explicit author (i.e. not mining), extract from the header
	var beneficiary common.Address
	// if author == nil {
	// 	beneficiary, _ = chain.Engine().Author(header) // Ignore error, we're past header validation
	// } else {
	// 	beneficiary = *author
	// }
	beneficiary = *author

	fooChainContext := FooChainContext{blockTime: blockTime}
	fooHash := utils.StringToHash("xxx")
	fooHeader := fooChainContext.GetHeader(fooHash, height)

	return vm.Context{
		CanTransfer: CanTransfer,
		Transfer:    Transfer,
		GetHash:     GetHashFn(fooHeader, fooChainContext),
		Origin:      msg.FromAddress(),
		Coinbase:    beneficiary,
		BlockNumber: new(big.Int).Set(fooHeader.Number),
		Time:        new(big.Int).Set(big.NewInt(int64(fooHeader.Time))),
		Difficulty:  new(big.Int).Set(fooHeader.Difficulty),
		GasLimit:    fooHeader.GasLimit,
		GasPrice:    big.NewInt(0),
	}
}

// GetHashFn returns a GetHashFunc which retrieves header hashes by number
func GetHashFn(ref *types.Header, chain ChainContext) func(n uint64) common.Hash {
	return func(n uint64) common.Hash {
		for header := chain.GetHeader(ref.ParentHash, ref.Number.Uint64()-1); header != nil; header = chain.GetHeader(header.ParentHash, header.Number.Uint64()-1) {
			if header.Number.Uint64() == n {
				return header.Hash()
			}
		}

		return common.Hash{}
	}
}

// CanTransfer checks wether there are enough funds in the address' account to make a transfer.
// This does not take the necessary gas in to account to make the transfer valid.
func CanTransfer(db vm.StateDB, addr common.Address, amount *big.Int) bool {
	return db.GetBalance(addr).Cmp(amount) >= 0
}

// Transfer subtracts amount from sender and adds amount to recipient using the given Db
func Transfer(db vm.StateDB, sender, recipient common.Address, amount *big.Int) {
	db.SubBalance(sender, amount)
	db.AddBalance(recipient, amount)
}
