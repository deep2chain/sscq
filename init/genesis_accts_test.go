package init

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	app "github.com/deep2chain/sscq/app/v0"
	"github.com/deep2chain/sscq/codec"
	sdk "github.com/deep2chain/sscq/types"
)

func TestAddGenesisAccount(t *testing.T) {
	cdc := codec.New()
	addr1 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	type args struct {
		appState     app.GenesisState
		addr         sdk.AccAddress
		coins        sdk.Coins
		vestingAmt   sdk.Coins
		vestingStart int64
		vestingEnd   int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"valid account",
			args{
				app.GenesisState{},
				addr1,
				sdk.Coins{},
				sdk.Coins{},
				0,
				0,
			},
			false,
		},
		{
			"dup account",
			args{
				app.GenesisState{Accounts: []app.GenesisAccount{{Address: addr1}}},
				addr1,
				sdk.Coins{},
				sdk.Coins{},
				0,
				0,
			},
			true,
		},
		{
			"invalid vesting amount",
			args{
				app.GenesisState{},
				addr1,
				sdk.Coins{sdk.NewInt64Coin("stake", 50)},
				sdk.Coins{sdk.NewInt64Coin("stake", 100)},
				0,
				0,
			},
			true,
		},
		{
			"invalid vesting times",
			args{
				app.GenesisState{},
				addr1,
				sdk.Coins{sdk.NewInt64Coin("stake", 50)},
				sdk.Coins{sdk.NewInt64Coin("stake", 50)},
				1654668078,
				1554668078,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := addGenesisAccount(
				cdc, tt.args.appState, tt.args.addr, tt.args.coins,
				tt.args.vestingAmt, tt.args.vestingStart, tt.args.vestingEnd,
			)
			require.Equal(t, tt.wantErr, (err != nil))
		})
	}
}
