package client

import (
	"github.com/deep2chain/sscq/client"
	faucetcmd "github.com/deep2chain/sscq/x/faucet/client/cli"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"
)

// ModuleClient exports all client functionality from this module
type ModuleClient struct {
	storeKey string
	cdc      *amino.Codec
}

func NewModuleClient(storeKey string, cdc *amino.Codec) ModuleClient {
	return ModuleClient{storeKey, cdc}
}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	sscqsvcTxCmd := &cobra.Command{
		Use:   "ss",
		Short: "SscqService transactions subcommands",
	}

	sscqsvcTxCmd.AddCommand(client.PostCommands(
		faucetcmd.GetCmdAdd(mc.cdc),
		// faucetcmd.GetCmdIssue(mc.cdc),
	)...)

	return sscqsvcTxCmd
}
