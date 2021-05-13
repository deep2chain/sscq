package client

import (
	"github.com/deep2chain/sscq/client"
	sscqservicecmd "github.com/deep2chain/sscq/x/core/client/cli"
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

// GetQueryCmd returns the cli query commands for this module
func (mc ModuleClient) GetQueryCmd() *cobra.Command {
	// Group sscqservice queries under a subcommand
	htdfsvcQueryCmd := &cobra.Command{
		Use:   "hs",
		Short: "Querying commands for the sscqservice module",
	}

	htdfsvcQueryCmd.AddCommand(client.GetCommands()...)

	return htdfsvcQueryCmd
}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	htdfsvcTxCmd := &cobra.Command{
		Use:   "hs",
		Short: "HtdfService transactions subcommands",
	}

	htdfsvcTxCmd.AddCommand(client.PostCommands(
		//sscqservicecmd.GetCmdAdd(mc.cdc),
		//sscqservicecmd.GetCmdIssue(mc.cdc),
		sscqservicecmd.GetCmdSend(mc.cdc),
		sscqservicecmd.GetCmdCreate(mc.cdc),
		sscqservicecmd.GetCmdSign(mc.cdc),
		sscqservicecmd.GetCmdBroadCast(mc.cdc),
	)...)

	return htdfsvcTxCmd
}
