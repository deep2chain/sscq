package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/deep2chain/sscq/accounts/keystore"
	v0 "github.com/deep2chain/sscq/app/v0"
	"github.com/deep2chain/sscq/client"
	"github.com/deep2chain/sscq/client/context"
	"github.com/deep2chain/sscq/client/keys"
	"github.com/deep2chain/sscq/client/utils"
	"github.com/deep2chain/sscq/codec"
	sdk "github.com/deep2chain/sscq/types"
	authtxb "github.com/deep2chain/sscq/x/auth/client/txbuilder"
	faucet "github.com/deep2chain/sscq/x/faucet"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// junying-todo-20190409
func GetCmdAdd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [amount]",
		Short: "publish new coin or add existing coin to system issuer except stake",
		Long:  "hscli tx add 5satoshi --fees=1satoshi --genfile /home/xxx/.hsd/config/genesis.json",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			GenFilePath := viper.GetString(v0.FlagGenFilePath)
			if GenFilePath == "" {
				fmt.Print("--genfile required. please indicate the path of genesis.json.\n")
				return nil
			}

			systemissuer, err := faucet.GetSystemIssuerFromFile(GenFilePath)
			if err != nil {
				return err
			}

			if strings.Contains(args[0], "stake") {
				fmt.Print("stake can't be added. Or, system will panic. \n")
				return nil
			}

			coins, err := sdk.ParseCoins(args[0])
			if err != nil {
				return err
			}

			msg := faucet.NewMsgAdd(systemissuer, coins)
			cliCtx.PrintResponse = true

			return CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg}, systemissuer) //not completed yet, need account name
		},
	}

	cmd.Flags().String(v0.FlagGenFilePath, "", "genesis.json path")
	cmd.MarkFlagRequired(v0.FlagGenFilePath)
	return client.PostCommands(cmd)[0]
}

func PrepareTxBuilder(txBldr authtxb.TxBuilder, cliCtx context.CLIContext, fromaddr sdk.AccAddress) (authtxb.TxBuilder, error) {

	// TODO: (ref #1903) Allow for user supplied account number without
	// automatically doing a manual lookup.
	if txBldr.AccountNumber() == 0 {
		accNum, err := cliCtx.GetAccountNumber(fromaddr)
		if err != nil {
			return txBldr, err
		}
		txBldr = txBldr.WithAccountNumber(accNum)
	}

	// TODO: (ref #1903) Allow for user supplied account sequence without
	// automatically doing a manual lookup.
	if txBldr.Sequence() == 0 {
		accSeq, err := cliCtx.GetAccountSequence(fromaddr)
		if err != nil {
			return txBldr, err
		}
		txBldr = txBldr.WithSequence(accSeq)
	}
	return txBldr, nil
}

func CompleteAndBroadcastTxCLI(txBldr authtxb.TxBuilder, cliCtx context.CLIContext, msgs []sdk.Msg, fromaddr sdk.AccAddress) error {
	//
	txBldr, err := PrepareTxBuilder(txBldr, cliCtx, fromaddr)
	if err != nil {
		return err
	}

	if txBldr.SimulateAndExecute() || cliCtx.Simulate {
		txBldr, err := utils.EnrichWithGas(txBldr, cliCtx, msgs)
		if err != nil {
			return err
		}

		gasEst := utils.GasEstimateResponse{GasEstimate: txBldr.GasWanted()}
		fmt.Fprintf(os.Stderr, "%s\n", gasEst.String())
	}

	passphrase, err := keys.ReadShortPassphraseFromStdin(sdk.AccAddress.String(fromaddr))
	if err != nil {
		return err
	}
	addr := sdk.AccAddress.String(fromaddr)
	ksw := keystore.NewKeyStoreWallet(keystore.DefaultKeyStoreHome())
	txBytes, err := ksw.BuildAndSign(txBldr, addr, passphrase, msgs)
	if err != nil {
		return err
	}
	// broadcast to a Tendermint node
	res, err := cliCtx.BroadcastTx(txBytes)
	cliCtx.PrintOutput(res)
	return err
}
