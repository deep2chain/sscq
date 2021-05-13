package cli

import (
	"fmt"

	"github.com/deep2chain/sscq/accounts/keystore"
	"github.com/deep2chain/sscq/client"
	"github.com/deep2chain/sscq/client/context"
	"github.com/deep2chain/sscq/codec"
	sdk "github.com/deep2chain/sscq/types"
	"github.com/deep2chain/sscq/x/auth"
	authtxb "github.com/deep2chain/sscq/x/auth/client/txbuilder"
	sscqservice "github.com/deep2chain/sscq/x/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/deep2chain/sscq/client/keys"
)

// junying-todo-20190327
// GetCmdSign is the CLI command for signing unsigned transaction
/*
	inspired by
	hscli tx sign unsigned.json --name junying >> signed.json
	hscli tx sign --validate-signatures signed.json
	hscli tx sign --signature-only  test.json --name junying
*/
func GetCmdSign(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sign [unsignedtransaction]",
		Short: "sign a transaction",
		Long:  "hscli tx sign 7b0a202...23 --sequence 1 --account-number 0 --offline=true --encode=false",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			// load sign tx from string
			stdTx, err := sscqservice.ReadStdTxFromRawData(cliCtx.Codec, args[0])
			if err != nil {
				return err
			}

			// if no signers
			if len(stdTx.GetSigners()) == 0 {
				return err //err.
			}
			
			passphrase, err := keys.ReadShortPassphraseFromStdin(sdk.AccAddress.String(stdTx.GetSigners()[0]))
			if err != nil {
				return err
			}

			offlineflag := viper.GetBool(sscqservice.FlagOffline)

			// sign
			res, err := SignTransaction(authtxb.NewTxBuilderFromCLI(), cliCtx, stdTx, passphrase, offlineflag)
			if err != nil {
				return err
			}

			// print
			encodeflag := viper.GetBool(sscqservice.FlagEncode)
			if !encodeflag {
				fmt.Printf("%s\n", res)
			} else {
				fmt.Printf("%s\n", sscqservice.Encode_Hex(res))
			}
			return nil
		},
	}
	cmd.Flags().Bool(sscqservice.FlagEncode, true, "encode enabled")
	cmd.Flags().Bool(sscqservice.FlagOffline, false, "offline disabled")
	return client.PostCommands(cmd)[0]
}

func populateAccountFromState(txBldr authtxb.TxBuilder, cliCtx context.CLIContext,
	addr sdk.AccAddress) (authtxb.TxBuilder, error) {
	if txBldr.AccountNumber() == 0 {
		accNum, err := cliCtx.GetAccountNumber(addr)
		if err != nil {
			return txBldr, err
		}
		txBldr = txBldr.WithAccountNumber(accNum)
	}

	if txBldr.Sequence() == 0 {
		accSeq, err := cliCtx.GetAccountSequence(addr)
		if err != nil {
			return txBldr, err
		}
		txBldr = txBldr.WithSequence(accSeq)
	}
	return txBldr, nil
}

//
func SignStdTx(txBldr authtxb.TxBuilder, cliCtx context.CLIContext, stdTx auth.StdTx, passphrase string, offline bool) (signedTx auth.StdTx, err error) {
	// from address
	if len(stdTx.GetSigners()) == 0 {
		return signedTx, nil
	}
	fromaddr := stdTx.GetSigners()[0]
	// accountnumber, accountsequence check
	if !offline {
		txBldr, err = populateAccountFromState(txBldr, cliCtx, fromaddr)
		if err != nil {
			return signedTx, err
		}
	}

	ksw := keystore.NewKeyStoreWallet(keystore.DefaultKeyStoreHome())

	// signature
	return ksw.SignStdTx(txBldr,stdTx,sdk.AccAddress.String(fromaddr), passphrase)
}

//
func SignTransaction(txBldr authtxb.TxBuilder, cliCtx context.CLIContext, stdTx auth.StdTx, passphrase string, offline bool) (res []byte, err error) {
	// signature
	signedTx, err := SignStdTx(txBldr, cliCtx, stdTx, passphrase, offline)
	if err != nil {
		return []byte("signing failed"), err
	}

	switch cliCtx.Indent {
	case true:
		res, err = cliCtx.Codec.MarshalJSONIndent(signedTx, "", "  ")
	default:
		res, err = cliCtx.Codec.MarshalJSON(signedTx)
	}

	if err != nil {
		return []byte("json creating failed"), err
	}
	return res, err
}
