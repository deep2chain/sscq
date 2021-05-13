package client

import (
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"

	"github.com/deep2chain/sscq/client"
	sstakingcli "github.com/deep2chain/sscq/x/staking/client/cli"
	stakingcli "github.com/deep2chain/sscq/x/staking/client/cli"
	"github.com/deep2chain/sscq/x/staking/types"
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
	stakingQueryCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Querying commands for the staking module",
	}

	stakingQueryCmd.AddCommand(client.GetCommands(
		stakingcli.GetCmdQueryDelegation(mc.storeKey, mc.cdc),
		stakingcli.GetCmdQueryDelegations(mc.storeKey, mc.cdc),
		stakingcli.GetCmdQueryUnbondingDelegation(mc.storeKey, mc.cdc),
		stakingcli.GetCmdQueryUnbondingDelegations(mc.storeKey, mc.cdc),
		stakingcli.GetCmdQueryRedelegation(mc.storeKey, mc.cdc),
		stakingcli.GetCmdQueryRedelegations(mc.storeKey, mc.cdc),
		stakingcli.GetCmdQueryValidator(mc.storeKey, mc.cdc),
		stakingcli.GetCmdQueryValidators(mc.storeKey, mc.cdc),
		stakingcli.GetCmdQueryValidatorDelegations(mc.storeKey, mc.cdc),
		stakingcli.GetCmdQueryValidatorUnbondingDelegations(mc.storeKey, mc.cdc),
		stakingcli.GetCmdQueryValidatorRedelegations(mc.storeKey, mc.cdc),
		stakingcli.GetCmdQueryParams(mc.storeKey, mc.cdc),
		stakingcli.GetCmdQueryPool(mc.storeKey, mc.cdc))...)

	return stakingQueryCmd

}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	stakingTxCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Staking transaction subcommands",
	}

	stakingTxCmd.AddCommand(client.PostCommands(
		sstakingcli.GetCmdCreateValidator(mc.cdc),
		sstakingcli.GetCmdEditValidator(mc.cdc),
		sstakingcli.GetCmdDelegate(mc.cdc),
		sstakingcli.GetCmdRedelegate(mc.storeKey, mc.cdc),
		sstakingcli.GetCmdUnbond(mc.storeKey, mc.cdc),
		sstakingcli.GetCmdUpgradeDelStatus(mc.storeKey, mc.cdc),
	)...)

	return stakingTxCmd
}
