package cli

import (
	"encoding/hex"

	"github.com/deep2chain/sscq/client"
	"github.com/deep2chain/sscq/client/context"
	"github.com/deep2chain/sscq/client/utils"
	"github.com/deep2chain/sscq/codec"
	sdk "github.com/deep2chain/sscq/types"
	authtxb "github.com/deep2chain/sscq/x/auth/client/txbuilder"
	"github.com/deep2chain/sscq/x/ibc"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagTo     = "to"
	flagAmount = "amount"
	flagChain  = "chain"
)

// IBCTransferCmd implements the IBC transfer command.
func IBCTransferCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "transfer",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			from := cliCtx.GetFromAddress()
			msg, err := buildMsg(from)
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	cmd.Flags().String(flagTo, "", "Address to send coins")
	cmd.Flags().String(flagAmount, "", "Amount of coins to send")
	cmd.Flags().String(flagChain, "", "Destination chain to send coins")

	return cmd
}

func buildMsg(from sdk.AccAddress) (sdk.Msg, error) {
	amount := viper.GetString(flagAmount)
	coins, err := sdk.ParseCoins(amount)
	if err != nil {
		return nil, err
	}

	dest := viper.GetString(flagTo)
	bz, err := hex.DecodeString(dest)
	if err != nil {
		return nil, err
	}
	to := sdk.AccAddress(bz)

	packet := ibc.NewIBCPacket(from, to, coins, viper.GetString(client.FlagChainID),
		viper.GetString(flagChain))

	msg := ibc.MsgIBCTransfer{
		IBCPacket: packet,
	}

	return msg, nil
}
