package client

import (
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"

	sdkclient "github.com/deep2chain/sscq/client"
	"github.com/deep2chain/sscq/x/mint"
	"github.com/deep2chain/sscq/x/mint/client/cli"
)

// ModuleClient exports all CLI client functionality from the minting module.
type ModuleClient struct {
	storeKey string
	cdc      *amino.Codec
}

func NewModuleClient(storeKey string, cdc *amino.Codec) ModuleClient {
	return ModuleClient{storeKey, cdc}
}

// GetQueryCmd returns the cli query commands for the minting module.
func (mc ModuleClient) GetQueryCmd() *cobra.Command {
	mintingQueryCmd := &cobra.Command{
		Use:   mint.ModuleName,
		Short: "Querying commands for the minting module",
	}

	mintingQueryCmd.AddCommand(
		sdkclient.GetCommands(
			cli.GetCmdQueryParams(mc.cdc),
			cli.GetCmdQueryInflation(mc.cdc),
			cli.GetCmdQueryAnnualProvisions(mc.cdc),
			cli.GetCmdQueryTotalProvisions(mc.cdc),
			cli.GetCmdQueryBlockRewards(mc.cdc),
		)...,
	)

	return mintingQueryCmd
}

// GetTxCmd returns the transaction commands for the minting module.
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	mintTxCmd := &cobra.Command{
		Use:   mint.ModuleName,
		Short: "Minting transaction subcommands",
	}

	return mintTxCmd
}
