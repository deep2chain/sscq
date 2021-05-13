package init

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/deep2chain/sscq/client/keys"
	"github.com/deep2chain/sscq/codec"
	sdk "github.com/deep2chain/sscq/types"
	"github.com/deep2chain/sscq/x/auth"
	"github.com/deep2chain/sscq/app"
	"github.com/deep2chain/sscq/app/v0"
	"github.com/deep2chain/sscq/server"
)

// AddGenesisAccountCmd returns add-genesis-account cobra Command.
func AddGenesisAccountCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-account [address_or_key_name] [coin][,[coin]]",
		Short: "Add genesis account to genesis.json",
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				kb, err := keys.NewKeyBaseFromDir(viper.GetString(flagClientHome))
				if err != nil {
					return err
				}

				info, err := kb.Get(args[0])
				if err != nil {
					return err
				}

				addr = info.GetAddress()
			}

			coins, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}

			vestingStart := viper.GetInt64(flagVestingStart)
			vestingEnd := viper.GetInt64(flagVestingEnd)
			vestingAmt, err := sdk.ParseCoins(viper.GetString(flagVestingAmt))
			if err != nil {
				return err
			}

			genFile := config.GenesisFile()
			if !common.FileExists(genFile) {
				return fmt.Errorf("%s does not exist, run `hsd init` first", genFile)
			}

			genDoc, err := LoadGenesisDoc(cdc, genFile)
			if err != nil {
				return err
			}

			var appState v0.GenesisState
			if err = cdc.UnmarshalJSON(genDoc.AppState, &appState); err != nil {
				return err
			}

			appState, err = addGenesisAccount(cdc, appState, addr, coins, vestingAmt, vestingStart, vestingEnd)
			if err != nil {
				return err
			}

			appStateJSON, err := cdc.MarshalJSON(appState)
			if err != nil {
				return err
			}

			return ExportGenesisFile(genFile, genDoc.ChainID, nil, appStateJSON)
		},
	}

	cmd.Flags().String(cli.HomeFlag, app.DefaultNodeHome, "node's home directory")
	cmd.Flags().String(flagClientHome, app.DefaultCLIHome, "client's home directory")
	cmd.Flags().String(flagVestingAmt, "", "amount of coins for vesting accounts")
	cmd.Flags().Uint64(flagVestingStart, 0, "schedule start time (unix epoch) for vesting accounts")
	cmd.Flags().Uint64(flagVestingEnd, 0, "schedule end time (unix epoch) for vesting accounts")

	return cmd
}

func addGenesisAccount(
	cdc *codec.Codec, appState v0.GenesisState, addr sdk.AccAddress,
	coins, vestingAmt sdk.Coins, vestingStart, vestingEnd int64,
) (v0.GenesisState, error) {

	for _, stateAcc := range appState.Accounts {
		if stateAcc.Address.Equals(addr) {
			return appState, fmt.Errorf("the application state already contains account %v", addr)
		}
	}

	acc := auth.NewBaseAccountWithAddress(addr)
	acc.Coins = coins

	if !vestingAmt.IsZero() {
		var vacc auth.VestingAccount

		bvacc := &auth.BaseVestingAccount{
			BaseAccount:     &acc,
			OriginalVesting: vestingAmt,
			EndTime:         vestingEnd,
		}

		if bvacc.OriginalVesting.IsAllGT(acc.Coins) {
			return appState, fmt.Errorf("vesting amount cannot be greater than total amount")
		}
		if vestingStart >= vestingEnd {
			return appState, fmt.Errorf("vesting start time must before end time")
		}

		if vestingStart != 0 {
			vacc = &auth.ContinuousVestingAccount{
				BaseVestingAccount: bvacc,
				StartTime:          vestingStart,
			}
		} else {
			vacc = &auth.DelayedVestingAccount{
				BaseVestingAccount: bvacc,
			}
		}

		appState.Accounts = append(appState.Accounts, v0.NewGenesisAccountI(vacc))
	} else {
		appState.Accounts = append(appState.Accounts, v0.NewGenesisAccount(&acc))
	}

	return appState, nil
}
