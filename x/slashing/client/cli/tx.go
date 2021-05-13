package cli

import (
	"github.com/deep2chain/sscq/client/context"
	"github.com/deep2chain/sscq/client/utils"
	"github.com/deep2chain/sscq/codec"
	sdk "github.com/deep2chain/sscq/types"
	authtxb "github.com/deep2chain/sscq/x/auth/client/txbuilder"
	"github.com/deep2chain/sscq/x/slashing"
	sscorecli "github.com/deep2chain/sscq/x/core/client/cli"
	"github.com/spf13/cobra"
)

// GetCmdUnjail implements the create unjail validator command.
func GetCmdUnjail(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "unjail [keyaddr]",
		Short: "unjail validator previously jailed for downtime",
		Long: `unjail a jailed validator:

$ sscli tx slashing unjail [keyaddr]
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			validatorAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := slashing.NewMsgUnjail(sdk.ValAddress(validatorAddr))
			return sscorecli.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg}, validatorAddr)
		},
	}
}
