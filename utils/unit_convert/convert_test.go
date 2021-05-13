package unit_convert

import (
	"fmt"
	"testing"

	"github.com/magiconair/properties/assert"
	sdk "github.com/deep2chain/sscq/types"
)

func TestConvert1(t *testing.T) {

	defaultCoin := sdk.NewCoin("satoshi", sdk.NewInt(1000000))
	var bigCoin sdk.BigCoin

	DefaultCoinToBigCoin(&defaultCoin, &bigCoin)
	assert.Equal(t, bigCoin.Denom, sdk.BigDenom)
	assert.Equal(t, bigCoin.Amount, "0.01")

	fmt.Printf("defaultCoin=%v\n", defaultCoin)
	fmt.Printf("bigCoin=%v\n", bigCoin)

	defaultCoin = sdk.NewCoin("stake", sdk.NewInt(1000000))
	DefaultCoinToBigCoin(&defaultCoin, &bigCoin)
	assert.Equal(t, bigCoin.Denom, "stake")
	assert.Equal(t, bigCoin.Amount, "1000000")

	fmt.Printf("defaultCoin=%v\n", defaultCoin)
	fmt.Printf("bigCoin=%v\n", bigCoin)

	var defaultCoins []sdk.Coin
	defaultCoins = append(defaultCoins, sdk.NewCoin("satoshi", sdk.NewInt(1000000)))
	defaultCoins = append(defaultCoins, sdk.NewCoin("satoshi", sdk.NewInt(2000000)))
	defaultCoins = append(defaultCoins, sdk.NewCoin("satoshi", sdk.NewInt(3000000)))

	fmt.Printf("defaultCoins=%v\n", defaultCoins)

	var bigCoins []sdk.BigCoin
	bigCoins = DefaultCoinsToBigCoins(defaultCoins)
	fmt.Printf("bigCoins=%v\n", bigCoins)

	//clear slice
	defaultCoins = defaultCoins[:0]
	defaultCoins = append(defaultCoins, sdk.NewCoin("stake", sdk.NewInt(1000000)))
	defaultCoins = append(defaultCoins, sdk.NewCoin("stake", sdk.NewInt(2000000)))
	defaultCoins = append(defaultCoins, sdk.NewCoin("stake", sdk.NewInt(3000000)))
	fmt.Printf("defaultCoins=%v\n", defaultCoins)

	bigCoins = bigCoins[:0]
	bigCoins = DefaultCoinsToBigCoins(defaultCoins)
	fmt.Printf("bigCoins=%v\n", bigCoins)

	defaultAmount := "123456789"
	bigAmount := DefaultAmoutToBigAmount(defaultAmount)
	assert.Equal(t, bigAmount, "1.23456789")

	fmt.Printf("bigAmount=%v\n", bigAmount)

}

func TestConvert2(t *testing.T) {

	var bigCoin sdk.BigCoin
	bigCoin.Denom = "htdf"
	bigCoin.Amount = "1.2345678"

	var defaultCoin sdk.Coin
	BigCoinToDefaultCoin(&bigCoin, &defaultCoin)

	assert.Equal(t, defaultCoin.String(), "123456780satoshi")

	fmt.Printf("bigCoin=%v\n", bigCoin)
	fmt.Printf("defaultCoin=%v\n", defaultCoin)

	bigCoin.Denom = "stake"
	bigCoin.Amount = "12345678"

	BigCoinToDefaultCoin(&bigCoin, &defaultCoin)

	assert.Equal(t, defaultCoin.String(), "12345678stake")

	fmt.Printf("bigCoin=%v\n", bigCoin)
	fmt.Printf("defaultCoin=%v\n", defaultCoin)

	//htdf to satoshi
	var bigCoins []sdk.BigCoin
	bigCoins = append(bigCoins, sdk.BigCoin{Denom: "htdf", Amount: "1.2345"})
	bigCoins = append(bigCoins, sdk.BigCoin{Denom: "htdf", Amount: "2.3456"})
	bigCoins = append(bigCoins, sdk.BigCoin{Denom: "htdf", Amount: "3.4567"})
	fmt.Printf("bigCoins=%v\n", bigCoins)

	var defaultCoins []sdk.Coin
	defaultCoins = BigCoinsToDefaultCoins(bigCoins)

	fmt.Printf("defaultCoins=%v\n", defaultCoins)

	//stake
	bigCoins = bigCoins[:0]
	bigCoins = append(bigCoins, sdk.BigCoin{Denom: "stake", Amount: "12345"})
	bigCoins = append(bigCoins, sdk.BigCoin{Denom: "stake", Amount: "23456"})
	bigCoins = append(bigCoins, sdk.BigCoin{Denom: "stake", Amount: "34567"})
	fmt.Printf("bigCoins=%v\n", bigCoins)

	defaultCoins = defaultCoins[:0]
	defaultCoins = BigCoinsToDefaultCoins(bigCoins)
	fmt.Printf("defaultCoins=%v\n", defaultCoins)

	bigAmount := "1.2345678"
	defaultAmount := BigAmountToDefaultAmount(bigAmount)
	assert.Equal(t, defaultAmount, "123456780")

	fmt.Printf("defaultAmount=%v\n", defaultAmount)

}
