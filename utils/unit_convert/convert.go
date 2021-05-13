package unit_convert

import (
	sdk "github.com/deep2chain/sscq/types"
	"github.com/shopspring/decimal"
)

//
func DefaultCoinToBigCoin(defaultCoin *sdk.Coin, bigCoin *sdk.BigCoin) {
	//convert when default denom
	if defaultCoin.Denom == sdk.DefaultDenom {
		bigCoin.Denom = sdk.BigDenom
		bigCoin.Amount = RightShift(*defaultCoin.Amount.BigInt(), 8).String()
	} else {
		bigCoin.Denom = defaultCoin.Denom
		bigCoin.Amount = defaultCoin.Amount.String()
	}
}

//
func DefaultCoinsToBigCoins(defaultCoins []sdk.Coin) (bigCoins []sdk.BigCoin) {
	for _, coin := range defaultCoins {
		var bigCoin sdk.BigCoin
		DefaultCoinToBigCoin(&coin, &bigCoin)
		bigCoins = append(bigCoins, bigCoin)
	}

	return bigCoins
}

//
func DefaultAmoutToBigAmount(defaultAmount string) (bigAmount string) {
	iDefaultAmount, _ := sdk.NewIntFromString(defaultAmount)
	return RightShift(*iDefaultAmount.BigInt(), 8).String()
}

//
func BigCoinToDefaultCoin(bigCoin *sdk.BigCoin, defaultCoin *sdk.Coin) {
	//convert when bug denom
	if bigCoin.Denom == sdk.BigDenom {
		defaultCoin.Denom = sdk.DefaultDenom

		decAmount, _ := decimal.NewFromString(bigCoin.Amount)
		defaultCoin.Amount, _ = sdk.NewIntFromString(LeftShift(decAmount, 8).String())
	} else {
		defaultCoin.Denom = bigCoin.Denom
		defaultCoin.Amount, _ = sdk.NewIntFromString(bigCoin.Amount)
	}
}

//
func BigCoinsToDefaultCoins(bigCoins []sdk.BigCoin) (defaultCoins []sdk.Coin) {

	for _, bigCoin := range bigCoins {
		var defaultCoin sdk.Coin
		BigCoinToDefaultCoin(&bigCoin, &defaultCoin)
		defaultCoins = append(defaultCoins, defaultCoin)
	}

	return defaultCoins
}

//
func BigAmountToDefaultAmount(bigAmount string) (defaultAmount string) {
	decBigAmount, _ := decimal.NewFromString(bigAmount)
	return LeftShift(decBigAmount, 8).String()
}
