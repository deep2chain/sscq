package types

import (
	"bytes"
	"fmt"

	"github.com/deep2chain/sscq/codec"
	sdk "github.com/deep2chain/sscq/types"
)

// Pool - tracking bonded and not-bonded token supply of the bond denomination
type Pool struct {
	NotBondedTokens           sdk.Int `json:"not_bonded_tokens"`       // tokens which are not bonded to a validator (unbonded or unbonding)
	BondedTokens              sdk.Int `json:"bonded_tokens"`           // tokens which are currently bonded to a validator
	LastZeroRewardBlockHeight int64   `json:"last_zero_block_height"`  // block height at which block reward is zero
	Amplitude                 int64   `json:"amplitude_sine_function"` //
	CycleAsBlocks             int64   `json:"cycle_sine_function"`     // cycle which is randomly generated when block height is last zero-reward block height.
}

// nolint
// TODO: This is slower than comparing struct fields directly
func (p Pool) Equal(p2 Pool) bool {
	bz1 := MsgCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := MsgCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

// initial pool for testing
func InitialPool() Pool {
	return Pool{
		NotBondedTokens:           sdk.ZeroInt(),
		BondedTokens:              sdk.ZeroInt(),
		LastZeroRewardBlockHeight: 1,
		CycleAsBlocks:             1,
		Amplitude:                 0,
	}
}

// // junying-todo, 2019-12-06
// func (p Pool) GetLastZeroRewardBlockHeight() int64 {
// 	return p.LastZeroRewardBlockHeight
// }

// func (p Pool) SetLastZeroRewardBlockHeight(blockheight int64) {
// 	p.LastZeroRewardBlockHeight = blockheight
// }

// func (p Pool) GetCycleAsBlocks() int64 {
// 	return p.CycleAsBlocks
// }

// func (p Pool) SetCycleAsBlocks(nextcycle int64) {
// 	p.CycleAsBlocks = nextcycle
// }

// func (p Pool) GetAmplitude() float64 {
// 	return p.Amplitude
// }

// func (p Pool) SetAmplitude(amp float64) {
// 	p.Amplitude = amp
// }

// Sum total of all staking tokens in the pool
func (p Pool) TokenSupply() sdk.Int {
	return p.NotBondedTokens.Add(p.BondedTokens)
}

// Get the fraction of the staking token which is currently bonded
func (p Pool) BondedRatio() sdk.Dec {
	supply := p.TokenSupply()
	if supply.IsPositive() {
		return p.BondedTokens.ToDec().QuoInt(supply)
	}
	return sdk.ZeroDec()
}

func (p Pool) notBondedTokensToBonded(bondedTokens sdk.Int) Pool {
	p.BondedTokens = p.BondedTokens.Add(bondedTokens)
	p.NotBondedTokens = p.NotBondedTokens.Sub(bondedTokens)
	if p.NotBondedTokens.IsNegative() {
		panic(fmt.Sprintf("sanity check: not-bonded tokens negative, pool: %v", p))
	}
	return p
}

func (p Pool) bondedTokensToNotBonded(bondedTokens sdk.Int) Pool {
	p.BondedTokens = p.BondedTokens.Sub(bondedTokens)
	p.NotBondedTokens = p.NotBondedTokens.Add(bondedTokens)
	if p.BondedTokens.IsNegative() {
		panic(fmt.Sprintf("sanity check: bonded tokens negative, pool: %v", p))
	}
	return p
}

// String returns a human readable string representation of a pool.
func (p Pool) String() string {
	return fmt.Sprintf(`Pool:
  Loose Tokens:  %s
  Bonded Tokens: %s
  Token Supply:  %s
  Bonded Ratio:  %v`, p.NotBondedTokens,
		p.BondedTokens, p.TokenSupply(),
		p.BondedRatio())
}

// unmarshal the current pool value from store key or panics
func MustUnmarshalPool(cdc *codec.Codec, value []byte) Pool {
	pool, err := UnmarshalPool(cdc, value)
	if err != nil {
		panic(err)
	}
	return pool
}

// unmarshal the current pool value from store key
func UnmarshalPool(cdc *codec.Codec, value []byte) (pool Pool, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &pool)
	if err != nil {
		return
	}
	return
}
