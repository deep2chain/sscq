package auth

import (
	codec "github.com/deep2chain/sscq/codec"
	sdk "github.com/deep2chain/sscq/types"
)

var (
	collectedFeesKey = []byte("collectedFees")
	feeAuthKey       = []byte("feeAuth")
)

// FeeCollectionKeeper handles collection of fees in the anteHandler
// and setting of MinFees for different fee tokens
type FeeCollectionKeeper struct {

	// The (unexposed) key used to access the fee store from the Context.
	key sdk.StoreKey

	// The codec codec for binary encoding/decoding of accounts.
	cdc *codec.Codec

	//paramSpace params.Subspace
}

// NewFeeCollectionKeeper returns a new FeeCollectionKeeper
func NewFeeCollectionKeeper(cdc *codec.Codec, key sdk.StoreKey) FeeCollectionKeeper {
	return FeeCollectionKeeper{
		key: key,
		cdc: cdc,
		// paramSpace: paramSpace.WithKeyTable(ParamKeyTable()),
	}
}

// GetCollectedFees - retrieves the collected fee pool
func (fck FeeCollectionKeeper) GetCollectedFees(ctx sdk.Context) sdk.Coins {
	store := ctx.KVStore(fck.key)
	bz := store.Get(collectedFeesKey)
	if bz == nil {
		return sdk.NewCoins()
	}

	emptyFees := sdk.NewCoins()
	feePool := &emptyFees
	fck.cdc.MustUnmarshalBinaryLengthPrefixed(bz, feePool)
	return *feePool
}

func (fck FeeCollectionKeeper) setCollectedFees(ctx sdk.Context, coins sdk.Coins) {
	bz := fck.cdc.MustMarshalBinaryLengthPrefixed(coins)
	store := ctx.KVStore(fck.key)
	store.Set(collectedFeesKey, bz)
}

// AddCollectedFees - add to the fee pool
func (fck FeeCollectionKeeper) AddCollectedFees(ctx sdk.Context, coins sdk.Coins) sdk.Coins {
	newCoins := fck.GetCollectedFees(ctx).Add(coins)
	fck.setCollectedFees(ctx, newCoins)

	return newCoins
}

// ClearCollectedFees - clear the fee pool
func (fck FeeCollectionKeeper) ClearCollectedFees(ctx sdk.Context) {
	fck.setCollectedFees(ctx, sdk.NewCoins())
}

func (fk FeeCollectionKeeper) GetFeeAuth(ctx sdk.Context) (feeAuth FeeAuth) {
	store := ctx.KVStore(fk.key)
	b := store.Get(feeAuthKey)
	if b == nil {
		panic("GetFeeAuth Stored fee pool should not have been nil")
	}
	fk.cdc.MustUnmarshalBinaryLengthPrefixed(b, &feeAuth)
	return
}

// func (fk FeeCollectionKeeper) GetParamSet(ctx sdk.Context) Params {
// 	var feeParams Params
// 	fk.paramSpace.GetParamSet(ctx, &feeParams)
// 	return feeParams
// }

// RefundCollectedFees deducts fees from fee collector
func (fk FeeCollectionKeeper) RefundCollectedFees(ctx sdk.Context, coins sdk.Coins) sdk.Coins {
	newCoins := fk.GetCollectedFees(ctx).Sub(coins)
	if !newCoins.IsAnyNegative() {
		panic("fee collector contains negative coins")
	}
	fk.setCollectedFees(ctx, newCoins)
	return newCoins
}
