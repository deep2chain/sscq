package client

import (
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"

	"github.com/deep2chain/sscq/client"
	"github.com/deep2chain/sscq/x/slashing"
	slashingcli "github.com/deep2chain/sscq/x/slashing/client/cli"
	hslashingcli "github.com/deep2chain/sscq/x/slashing/client/cli"
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
	// Group slashing queries under a subcommand
	slashingQueryCmd := &cobra.Command{
		Use:   slashing.ModuleName,
		Short: "Querying commands for the slashing module",
	}

	slashingQueryCmd.AddCommand(
		client.GetCommands(
			slashingcli.GetCmdQuerySigningInfo(mc.storeKey, mc.cdc),
			slashingcli.GetCmdQueryParams(mc.cdc),
		)...,
	)

	return slashingQueryCmd

}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	slashingTxCmd := &cobra.Command{
		Use:   slashing.ModuleName,
		Short: "Slashing transactions subcommands",
	}

	slashingTxCmd.AddCommand(client.PostCommands(
		hslashingcli.GetCmdUnjail(mc.cdc),
	)...)

	return slashingTxCmd
}
