package init

// DONTCOVER

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/deep2chain/sscq/client"
	"github.com/deep2chain/sscq/client/context"
	"github.com/deep2chain/sscq/client/utils"
	"github.com/deep2chain/sscq/codec"
	sdk "github.com/deep2chain/sscq/types"
	"github.com/deep2chain/sscq/x/auth"
	authtxb "github.com/deep2chain/sscq/x/auth/client/txbuilder"
	"github.com/deep2chain/sscq/x/staking/client/cli"
	"github.com/deep2chain/sscq/app"
	"github.com/deep2chain/sscq/app/v0"
	"github.com/deep2chain/sscq/server"
	"github.com/deep2chain/sscq/client/keys"
	hscorecli "github.com/deep2chain/sscq/x/core/client/cli"
	hstakingcli "github.com/deep2chain/sscq/x/staking/client/cli"
	"github.com/deep2chain/sscq/accounts/keystore"
)

var (
	defaultTokens                  = sdk.TokensFromTendermintPower(100)
	defaultAmount                  = defaultTokens.String() + sdk.DefaultBondDenom
	defaultCommissionRate          = "0.1"
	defaultCommissionMaxRate       = "0.2"
	defaultCommissionMaxChangeRate = "0.01"
	defaultMinSelfDelegation       = "1"
)

// GenTxCmd builds the hsd gentx command.
// nolint: errcheck
func GenTxCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gentx [validatorAddr]",
		Short: "Generate a genesis tx carrying a self delegation",
		Long: fmt.Sprintf(`This command is an alias of the 'hsd tx create-validator' command'.

It creates a genesis piece carrying a self delegation with the
following delegation and commission default parameters:

	delegation amount:           %s
	commission rate:             %s
	commission max rate:         %s
	commission max change rate:  %s
	minimum self delegation:     %s
`, defaultAmount, defaultCommissionRate, defaultCommissionMaxRate, defaultCommissionMaxChangeRate, defaultMinSelfDelegation),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			config := ctx.Config
			config.SetRoot(viper.GetString(tmcli.HomeFlag))
			nodeID, valPubKey, err := InitializeNodeValidatorFiles(ctx.Config)
			if err != nil {
				return err
			}

			// Read --nodeID, if empty take it from priv_validator.json
			if nodeIDString := viper.GetString(cli.FlagNodeID); nodeIDString != "" {
				nodeID = nodeIDString
			}

			ip := viper.GetString(cli.FlagIP)
			if ip == "" {
				fmt.Fprintf(os.Stderr, "couldn't retrieve an external IP; "+
					"the tx's memo field will be unset")
			}

			genDoc, err := LoadGenesisDoc(cdc, config.GenesisFile())
			if err != nil {
				return err
			}

			genesisState := v0.GenesisState{}
			if err = cdc.UnmarshalJSON(genDoc.AppState, &genesisState); err != nil {
				return err
			}

			validatorAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			// Read --pubkey, if empty take it from priv_validator.json
			if valPubKeyString := viper.GetString(cli.FlagPubKey); valPubKeyString != "" {
				valPubKey, err = sdk.GetConsPubKeyBech32(valPubKeyString)
				if err != nil {
					return err
				}
			}
			// Set flags for creating gentx
			prepareFlagsForTxCreateValidator(config, nodeID, ip, genDoc.ChainID, valPubKey)

			// Fetch the amount of coins staked
			amount := viper.GetString(cli.FlagAmount)
			coins, err := sdk.ParseCoins(amount)
			if err != nil {
				return err
			}
			err = accountInGenesis(genesisState, validatorAddr, coins)
			if err != nil {
				return err
			}

			// Run hsd tx create-validator
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr, msg, err := hstakingcli.BuildCreateValidatorMsg(cliCtx, txBldr, validatorAddr)
			if err != nil {
				return err
			}

			// write the unsigned transaction to the buffer
			w := bytes.NewBuffer([]byte{})
			cliCtx = cliCtx.WithOutput(w)
			if err = hscorecli.PrintUnsignedStdTx(txBldr, cliCtx, []sdk.Msg{msg}, false); err != nil {
				return err
			}

			// read the transaction
			stdTx, err := readUnsignedGenTxFile(cdc, w)
			if err != nil {
				return err
			}

			passphrase, err := keys.ReadShortPassphraseFromStdin(sdk.AccAddress.String(stdTx.GetSigners()[0]))
			if err != nil {
				return err
			}

			fromaddr := stdTx.GetSigners()[0]
			ksw := keystore.NewKeyStoreWallet(keystore.DefaultKeyStoreHome())
			signedTx, err := ksw.SignStdTx(txBldr,stdTx,sdk.AccAddress.String(fromaddr), passphrase)
			if err != nil {
				return err
			}

			// Fetch output file name
			outputDocument := viper.GetString(client.FlagOutputDocument)
			if outputDocument == "" {
				outputDocument, err = makeOutputFilepath(config.RootDir, nodeID)
				if err != nil {
					return err
				}
			}

			if err := writeSignedGenTx(cdc, outputDocument, signedTx); err != nil {
				return err
			}

			fmt.Fprintf(os.Stderr, "Genesis transaction written to %q\n", outputDocument)
			return nil
		},
	}

	ip, _ := server.ExternalIP()

	cmd.Flags().String(tmcli.HomeFlag, app.DefaultNodeHome, "node's home directory")
	cmd.Flags().String(flagClientHome, app.DefaultCLIHome, "client's home directory")
	cmd.Flags().String(client.FlagName, "", "name of private key with which to sign the gentx")
	cmd.Flags().String(client.FlagOutputDocument, "",
		"write the genesis transaction JSON document to the given file instead of the default location")
	cmd.Flags().String(cli.FlagIP, ip, "The node's public IP")
	cmd.Flags().String(cli.FlagNodeID, "", "The node's NodeID")
	cmd.Flags().AddFlagSet(cli.FsCommissionCreate)
	cmd.Flags().AddFlagSet(cli.FsMinSelfDelegation)
	cmd.Flags().AddFlagSet(cli.FsAmount)
	cmd.Flags().AddFlagSet(cli.FsPk)
	return cmd
}

func accountInGenesis(genesisState v0.GenesisState, key sdk.AccAddress, coins sdk.Coins) error {
	accountIsInGenesis := false
	bondDenom := genesisState.StakeData.Params.BondDenom

	// Check if the account is in genesis
	for _, acc := range genesisState.Accounts {
		// Ensure that account is in genesis
		if acc.Address.Equals(key) {

			// Ensure account contains enough funds of default bond denom
			if coins.AmountOf(bondDenom).GT(acc.Coins.AmountOf(bondDenom)) {
				return fmt.Errorf(
					"account %v is in genesis, but it only has %v%v available to stake, not %v%v",
					key.String(), acc.Coins.AmountOf(bondDenom), bondDenom, coins.AmountOf(bondDenom), bondDenom,
				)
			}
			accountIsInGenesis = true
			break
		}
	}

	if accountIsInGenesis {
		return nil
	}

	return fmt.Errorf("account %s in not in the app_state.accounts array of genesis.json", key)
}

func prepareFlagsForTxCreateValidator(config *cfg.Config, nodeID, ip, chainID string,
	valPubKey crypto.PubKey) {
	viper.Set(tmcli.HomeFlag, viper.GetString(flagClientHome)) // --home
	viper.Set(client.FlagChainID, chainID)
	viper.Set(client.FlagFrom, viper.GetString(client.FlagName))   // --from
	viper.Set(cli.FlagNodeID, nodeID)                              // --node-id
	viper.Set(cli.FlagIP, ip)                                      // --ip
	viper.Set(cli.FlagPubKey, sdk.MustBech32ifyConsPub(valPubKey)) // --pubkey
	viper.Set(client.FlagGenerateOnly, true)                       // --genesis-format
	viper.Set(cli.FlagMoniker, config.Moniker)                     // --moniker
	if config.Moniker == "" {
		viper.Set(cli.FlagMoniker, viper.GetString(client.FlagName))
	}
	if viper.GetString(cli.FlagAmount) == "" {
		viper.Set(cli.FlagAmount, defaultAmount)
	}
	if viper.GetString(cli.FlagCommissionRate) == "" {
		viper.Set(cli.FlagCommissionRate, defaultCommissionRate)
	}
	if viper.GetString(cli.FlagCommissionMaxRate) == "" {
		viper.Set(cli.FlagCommissionMaxRate, defaultCommissionMaxRate)
	}
	if viper.GetString(cli.FlagCommissionMaxChangeRate) == "" {
		viper.Set(cli.FlagCommissionMaxChangeRate, defaultCommissionMaxChangeRate)
	}
	if viper.GetString(cli.FlagMinSelfDelegation) == "" {
		viper.Set(cli.FlagMinSelfDelegation, defaultMinSelfDelegation)
	}
}

func makeOutputFilepath(rootDir, nodeID string) (string, error) {
	writePath := filepath.Join(rootDir, "config", "gentx")
	if err := common.EnsureDir(writePath, 0700); err != nil {
		return "", err
	}
	return filepath.Join(writePath, fmt.Sprintf("gentx-%v.json", nodeID)), nil
}

func readUnsignedGenTxFile(cdc *codec.Codec, r io.Reader) (auth.StdTx, error) {
	var stdTx auth.StdTx
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return stdTx, err
	}
	err = cdc.UnmarshalJSON(bytes, &stdTx)
	return stdTx, err
}

// nolint: errcheck
func writeSignedGenTx(cdc *codec.Codec, outputDocument string, tx auth.StdTx) error {
	outputFile, err := os.OpenFile(outputDocument, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	json, err := cdc.MarshalJSON(tx)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(outputFile, "%s\n", json)
	return err
}
