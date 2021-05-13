package cli

import (
	"github.com/deep2chain/sscq/client"
	"github.com/deep2chain/sscq/client/context"
	"github.com/deep2chain/sscq/codec"
	sscqservice "github.com/deep2chain/sscq/x/core"
	"github.com/spf13/cobra"
)

// junying-todo-20190327
// GetCmdBroadCast is the CLI command for broadcasting a signed transaction
/*
	inspired by
	sscli tx broadcast signed.json
*/
func GetCmdBroadCast(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "broadcast [rawdata]",
		Short: "broadcast signed transaction",
		Long:  "sscli tx broadcast 72032..13123",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// load sign tx from string
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			stdTx, err := sscqservice.ReadStdTxFromRawData(cliCtx.Codec, args[0])
			if err != nil {
				return err
			}
			// convert tx to bytes
			txBytes, err := cliCtx.Codec.MarshalBinaryLengthPrefixed(stdTx)
			if err != nil {
				return err
			}
			// broadcast
			res, err := cliCtx.BroadcastTx(txBytes)
			cliCtx.PrintOutput(res)
			return err

		},
	}
	return client.PostCommands(cmd)[0]
}
