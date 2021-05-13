package cli

import (
	"fmt"

	"github.com/deep2chain/sscq/app/protocol"
	"github.com/deep2chain/sscq/client/context"
	"github.com/deep2chain/sscq/codec"
	"github.com/deep2chain/sscq/x/guardian"
	"github.com/spf13/cobra"
)

func GetCmdQueryProfilers(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "profilers",
		Short:   "Query for all profilers",
		Example: "sscli guardian profilers",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", protocol.GuardianRoute, guardian.QueryProfilers), nil)

			if err != nil {
				return err
			}

			var profilers guardian.Profilers
			err = cdc.UnmarshalJSON(res, &profilers)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(profilers)
		},
	}
	return cmd
}

func GetCmdQueryTrustees(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "trustees",
		Short:   "Query for all trustees",
		Example: "sscli guardian trustees",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", protocol.GuardianRoute, guardian.QueryTrustees), nil)

			if err != nil {
				return err
			}

			var trustees guardian.Trustees
			err = cdc.UnmarshalJSON(res, &trustees)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(trustees)
		},
	}
	return cmd
}
