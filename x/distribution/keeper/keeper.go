package keeper

import (
	"github.com/deep2chain/sscq/codec"
	sdk "github.com/deep2chain/sscq/types"
	"github.com/deep2chain/sscq/x/distribution/types"
	"github.com/deep2chain/sscq/x/params"
)

// keeper of the staking store
type Keeper struct {
	storeKey            sdk.StoreKey
	cdc                 *codec.Codec
	paramSpace          params.Subspace
	bankKeeper          types.BankKeeper
	stakingKeeper       types.StakingKeeper
	feeCollectionKeeper types.FeeCollectionKeeper

	// codespace
	codespace sdk.CodespaceType
	// metrics
	metrics *Metrics
}

// create a new keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace, ck types.BankKeeper,
	sk types.StakingKeeper, fck types.FeeCollectionKeeper, codespace sdk.CodespaceType, metrics *Metrics) Keeper {
	keeper := Keeper{
		storeKey:            key,
		cdc:                 cdc,
		paramSpace:          paramSpace.WithKeyTable(ParamKeyTable()),
		bankKeeper:          ck,
		stakingKeeper:       sk,
		feeCollectionKeeper: fck,
		codespace:           codespace,
		metrics:             metrics,
	}
	return keeper
}

// set withdraw address
func (k Keeper) SetWithdrawAddr(ctx sdk.Context, delegatorAddr sdk.AccAddress, withdrawAddr sdk.AccAddress) sdk.Error {
	if !k.GetWithdrawAddrEnabled(ctx) {
		return types.ErrSetWithdrawAddrDisabled(k.codespace)
	}

	k.SetDelegatorWithdrawAddr(ctx, delegatorAddr, withdrawAddr)

	return nil
}

// withdraw rewards from a delegation
func (k Keeper) WithdrawDelegationRewards(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdk.Coins, sdk.Error) {
	val := k.stakingKeeper.Validator(ctx, valAddr)
	if val == nil {
		return nil, types.ErrNoValidatorDistInfo(k.codespace)
	}

	del := k.stakingKeeper.Delegation(ctx, delAddr, valAddr)
	if del == nil {
		return nil, types.ErrNoDelegationDistInfo(k.codespace)
	}

	// withdraw rewards
	rewards, err := k.withdrawDelegationRewards(ctx, val, del)
	if err != nil {
		return nil, err
	}

	// reinitialize the delegation
	k.initializeDelegation(ctx, valAddr, delAddr)

	return rewards, nil
}

// withdraw validator commission
func (k Keeper) WithdrawValidatorCommission(ctx sdk.Context, valAddr sdk.ValAddress) (sdk.Coins, sdk.Error) {
	// fetch validator accumulated commission
	commission := k.GetValidatorAccumulatedCommission(ctx, valAddr)
	if commission.IsZero() {
		return nil, types.ErrNoValidatorCommission(k.codespace)
	}

	coins, remainder := commission.TruncateDecimal()
	k.SetValidatorAccumulatedCommission(ctx, valAddr, remainder) // leave remainder to withdraw later

	// update outstanding
	outstanding := k.GetValidatorOutstandingRewards(ctx, valAddr)
	k.SetValidatorOutstandingRewards(ctx, valAddr, outstanding.Sub(sdk.NewDecCoins(coins)))

	if !coins.IsZero() {
		accAddr := sdk.AccAddress(valAddr)
		withdrawAddr := k.GetDelegatorWithdrawAddr(ctx, accAddr)

		if _, _, err := k.bankKeeper.AddCoins(ctx, withdrawAddr, coins); err != nil {
			return nil, err
		}
	}

	return coins, nil
}
