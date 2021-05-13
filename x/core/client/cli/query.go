package cli

import (
	"fmt"

	"github.com/deep2chain/sscq/client"
	"github.com/deep2chain/sscq/client/context"
	"github.com/deep2chain/sscq/codec"
	sdk "github.com/deep2chain/sscq/types"
	sscqservice "github.com/deep2chain/sscq/x/core"
	"github.com/spf13/cobra"
)

// junying-todo, 2020-04-01
// GetCmdCall is the CLI command for call contract.
func GetCmdCall(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contract [contract-address] [callcode]",
		Short: "query contract data",
		Long:  "sscli query contract sscq...  7839124400000000...",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// load sign tx from string
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			contractaddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			callcode := args[1]
			//
			bz, err := cliCtx.Codec.MarshalJSON(sscqservice.NewQueryContractParams(contractaddr, callcode))
			if err != nil {
				return err
			}
			route := fmt.Sprintf("custom/%s/%s", sscqservice.QuerierRoute, sscqservice.QueryContract)
			res, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}
			var answer string
			if err := cliCtx.Codec.UnmarshalJSON(res, &answer); err != nil {
				return err
			}
			//
			// cliCtx.PrintOutput(res)
			fmt.Println(answer)
			return nil

		},
	}
	return client.PostCommands(cmd)[0]
}
