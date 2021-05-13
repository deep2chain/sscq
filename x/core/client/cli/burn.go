package cli

import (
	"fmt"

	"github.com/deep2chain/sscq/client"
	"github.com/deep2chain/sscq/client/context"
	"github.com/deep2chain/sscq/client/utils"
	"github.com/deep2chain/sscq/codec"
	sdk "github.com/deep2chain/sscq/types"
	authtxb "github.com/deep2chain/sscq/x/auth/client/txbuilder"
	sscqservice "github.com/deep2chain/sscq/x/core"
	"github.com/spf13/cobra"
)

const (
	// hscli bech32 h2b 0000000000000000000000000000000000000000
	blackholeAddr = "htdf1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq0d4n7t"
)

// junying-todo-20200525
// GetCmdBurn burn owner's coin
func GetCmdBurn(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn [burn] [amount]",
		Short: "burn own balance",
		Long: `hscli tx burn htdf1qn38r8re3lwlf5t6zgrdycrerd5w0 \
							 5satoshi \
							 --gas=30000 \
							 --gas-price=100`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			fmt.Println("GetCmdBurn:txBldr.GasWanted()", txBldr.GasWanted())

			burnaddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			blackholeaddr, err := sdk.AccAddressFromBech32(blackholeAddr)
			if err != nil {
				return err
			}

			coins, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}

			if txBldr.GasPrice() == 0 {
				return sdk.ErrTxDecode("no gasprice")
			}

			gas := txBldr.GasWanted()
			fmt.Println("GetCmdBurn:txBldr.GasPrices():", txBldr.GasPrice())
			msg := sscqservice.NewMsgSend(burnaddr, blackholeaddr, coins, txBldr.GasPrice(), gas)

			cliCtx.PrintResponse = true

			return CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg}, burnaddr) //not completed yet, need account name
		},
	}
	return client.PostCommands(cmd)[0]
}
