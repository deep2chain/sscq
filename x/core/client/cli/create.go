package cli

import (
	"fmt"
	"os"

	"github.com/deep2chain/sscq/client"
	"github.com/deep2chain/sscq/client/context"
	"github.com/deep2chain/sscq/client/utils"
	"github.com/deep2chain/sscq/codec"
	sdk "github.com/deep2chain/sscq/types"
	"github.com/deep2chain/sscq/x/auth"
	authtxb "github.com/deep2chain/sscq/x/auth/client/txbuilder"
	sscqservice "github.com/deep2chain/sscq/x/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// junying-todo-20190326
// GetCmdCreate is the CLI command for creating unsigned transaction
/*
	inspired by
	hscli send --generate-only cosmos1yqgv2rhxcgrf5jqrxlg80at5szzlarlcy254re 5sscqtoken --from junying > unsigned.json
	utils.PrintUnsignedStdTx(txBldr, cliCtx, []sdk.Msg{msg},false)
	Tips:
	check functions in utils
*/
func GetCmdCreate(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [fromaddr] [toaddr] [amount]",
		Short: "create unsigned transaction",
		Long:  "hscli tx create cosmos1tq7zajghkxct4al0yf44ua9rjwnw06vdusflk4 cosmos1yqgv2rhxcgrf5jqrxlg80at5szzlarlcy254re 5satoshi --fees=1satoshi",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			fromaddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			toaddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			coins, err := sdk.ParseCoins(args[2])
			if err != nil {
				return err
			}

			encodeflag := viper.GetBool(sscqservice.FlagEncode)

			msg := sscqservice.NewMsgSendDefault(fromaddr, toaddr, coins)

			return PrintUnsignedStdTx(txBldr, cliCtx, []sdk.Msg{msg}, encodeflag)
		},
	}
	cmd.Flags().Bool(sscqservice.FlagEncode, true, "encode enabled")
	return client.PostCommands(cmd)[0]
}

// PrintUnsignedStdTx builds an unsigned StdTx and prints it to os.Stdout.
// Don't perform online validation or lookups if offline is true.
func PrintUnsignedStdTx(txBldr authtxb.TxBuilder, cliCtx context.CLIContext, msgs []sdk.Msg, encodeflag bool) (err error) {
	if txBldr.SimulateAndExecute() {
		txBldr, err = utils.EnrichWithGas(txBldr, cliCtx, msgs)
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "estimated gas = %v\n", txBldr.GasWanted())
	}
	stdSignMsg, err := txBldr.BuildSignMsg(msgs)
	if err != nil {
		return
	}
	//var stdTx auth.StdTx
	stdTx := auth.NewStdTx(stdSignMsg.Msgs, stdSignMsg.Fee, nil, stdSignMsg.Memo)

	if err != nil {
		return
	}
	json, err := cliCtx.Codec.MarshalJSON(stdTx)
	if err != nil {
		return
	}
	if !encodeflag {
		fmt.Fprintf(cliCtx.Output, "%s\n", json)
	} else {
		encoded := sscqservice.Encode_Hex(json)
		fmt.Fprintf(cliCtx.Output, "%s\n", encoded)
	}

	return
}
