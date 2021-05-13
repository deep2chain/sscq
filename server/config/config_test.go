package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/deep2chain/sscq/types"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	fmt.Printf("minGasPrices=%v\n", cfg.GetMinGasPrices())
	require.True(t, !cfg.GetMinGasPrices().IsZero())
}

func TestSetMinimumFees(t *testing.T) {
	cfg := DefaultConfig()
	cfg.SetMinGasPrices(sdk.Coins{sdk.NewInt64Coin("foo", 5)})
	require.Equal(t, "5foo", cfg.MinGasPrices)
}
