package cli

import (
	"github.com/spf13/cobra"

	"github.com/deep2chain/sscq/client"
	"github.com/deep2chain/sscq/client/context"
	"github.com/deep2chain/sscq/codec"
	sdk "github.com/deep2chain/sscq/types"
)

// GetAccountCmd returns a query account that will display the state of the
// account at a given address.
// nolint: unparam
func GetAccountCmd(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account [address]",
		Short: "Query account balance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).WithAccountDecoder(cdc)

			key, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			if err = cliCtx.EnsureAccountExistsFromAddr(key); err != nil {
				return err
			}

			acc, err := cliCtx.GetAccount(key)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(acc)
		},
	}
	return client.GetCommands(cmd)[0]
}
