package client

import (
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"

	"github.com/deep2chain/sscq/client"
	distcli "github.com/deep2chain/sscq/x/distribution/client/cli"
	ssdistcli "github.com/deep2chain/sscq/x/distribution/client/cli"
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
	distQueryCmd := &cobra.Command{
		Use:   "distr",
		Short: "Querying commands for the distribution module",
	}

	distQueryCmd.AddCommand(client.GetCommands(
		distcli.GetCmdQueryParams(mc.storeKey, mc.cdc),
		distcli.GetCmdQueryValidatorOutstandingRewards(mc.storeKey, mc.cdc),
		distcli.GetCmdQueryValidatorCommission(mc.storeKey, mc.cdc),
		distcli.GetCmdQueryValidatorSlashes(mc.storeKey, mc.cdc),
		distcli.GetCmdQueryDelegatorRewards(mc.storeKey, mc.cdc),
		distcli.GetCmdQueryCommunityPool(mc.storeKey, mc.cdc),
	)...)

	return distQueryCmd
}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	distTxCmd := &cobra.Command{
		Use:   "distr",
		Short: "Distribution transactions subcommands",
	}

	distTxCmd.AddCommand(client.PostCommands(
		ssdistcli.GetCmdWithdrawRewards(mc.cdc),
		ssdistcli.GetCmdSetWithdrawAddr(mc.cdc),
		ssdistcli.GetCmdWithdrawAllRewards(mc.cdc, mc.storeKey),
	)...)

	return distTxCmd
}
