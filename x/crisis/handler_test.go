package crisis

import (
	"errors"
	"fmt"
	"testing"

	sdk "github.com/deep2chain/sscq/types"
	"github.com/deep2chain/sscq/x/auth"
	distr "github.com/deep2chain/sscq/x/distribution"
	"github.com/stretchr/testify/require"
)

var (
	testModuleName        = "dummy"
	dummyRouteWhichPasses = NewInvarRoute(testModuleName, "which-passes", func(_ sdk.Context) error { return nil })
	dummyRouteWhichFails  = NewInvarRoute(testModuleName, "which-fails", func(_ sdk.Context) error { return errors.New("whoops") })
	addrs                 = distr.TestAddrs
)

func CreateTestInput(t *testing.T) (sdk.Context, Keeper, auth.AccountKeeper, distr.Keeper) {

	communityTax := sdk.NewDecWithPrec(2, 2)
	ctx, accKeeper, bankKeeper, distrKeeper, _, feeCollectionKeeper, paramsKeeper :=
		distr.CreateTestInputAdvanced(t, false, 10, communityTax)

	paramSpace := paramsKeeper.Subspace(DefaultParamspace)
	crisisKeeper := NewKeeper(paramSpace, distrKeeper, bankKeeper, feeCollectionKeeper)
	constantFee := sdk.NewInt64Coin("satoshi", 10000000)
	crisisKeeper.SetConstantFee(ctx, constantFee)

	crisisKeeper.RegisterRoute(testModuleName, dummyRouteWhichPasses.Route, dummyRouteWhichPasses.Invar)
	crisisKeeper.RegisterRoute(testModuleName, dummyRouteWhichFails.Route, dummyRouteWhichFails.Invar)

	// set the community pool to pay back the constant fee
	feePool := distr.InitialFeePool()
	feePool.CommunityPool = sdk.NewDecCoins(sdk.NewCoins(constantFee))
	distrKeeper.SetFeePool(ctx, feePool)

	return ctx, crisisKeeper, accKeeper, distrKeeper
}

//____________________________________________________________________________

func TestHandleMsgVerifyInvariantWithNotEnoughSenderCoins(t *testing.T) {
	ctx, crisisKeeper, accKeeper, _ := CreateTestInput(t)
	sender := addrs[0]
	coin := accKeeper.GetAccount(ctx, sender).GetCoins()[0]
	excessCoins := sdk.NewCoin(coin.Denom, coin.Amount.AddRaw(1))
	crisisKeeper.SetConstantFee(ctx, excessCoins)

	msg := NewMsgVerifyInvariant(sender, testModuleName, dummyRouteWhichPasses.Route)
	res := handleMsgVerifyInvariant(ctx, msg, crisisKeeper)
	require.False(t, res.IsOK())
}

func TestHandleMsgVerifyInvariantWithBadInvariant(t *testing.T) {
	ctx, crisisKeeper, _, _ := CreateTestInput(t)
	sender := addrs[0]

	msg := NewMsgVerifyInvariant(sender, testModuleName, "route-that-doesnt-exist")
	res := handleMsgVerifyInvariant(ctx, msg, crisisKeeper)
	require.False(t, res.IsOK())
}

func TestHandleMsgVerifyInvariantWithInvariantBroken(t *testing.T) {
	ctx, crisisKeeper, _, _ := CreateTestInput(t)
	sender := addrs[0]

	msg := NewMsgVerifyInvariant(sender, testModuleName, dummyRouteWhichFails.Route)
	var res sdk.Result
	require.Panics(t, func() {
		res = handleMsgVerifyInvariant(ctx, msg, crisisKeeper)
	}, fmt.Sprintf("%v", res))
}

func TestHandleMsgVerifyInvariantWithInvariantBrokenAndNotEnoughPoolCoins(t *testing.T) {
	ctx, crisisKeeper, _, distrKeeper := CreateTestInput(t)
	sender := addrs[0]

	// set the community pool to empty
	feePool := distrKeeper.GetFeePool(ctx)
	feePool.CommunityPool = sdk.DecCoins{}
	distrKeeper.SetFeePool(ctx, feePool)

	msg := NewMsgVerifyInvariant(sender, testModuleName, dummyRouteWhichFails.Route)
	var res sdk.Result
	require.Panics(t, func() {
		res = handleMsgVerifyInvariant(ctx, msg, crisisKeeper)
	}, fmt.Sprintf("%v", res))
}

func TestHandleMsgVerifyInvariantWithInvariantNotBroken(t *testing.T) {
	ctx, crisisKeeper, _, _ := CreateTestInput(t)
	sender := addrs[0]

	msg := NewMsgVerifyInvariant(sender, testModuleName, dummyRouteWhichPasses.Route)
	res := handleMsgVerifyInvariant(ctx, msg, crisisKeeper)
	require.True(t, res.IsOK())
}
