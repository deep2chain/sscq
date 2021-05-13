package cli

import (
	"fmt"

	"github.com/deep2chain/sscq/client/context"
	"github.com/deep2chain/sscq/codec"
	sdk "github.com/deep2chain/sscq/types"
	stake "github.com/deep2chain/sscq/x/staking"
	"github.com/deep2chain/sscq/x/staking/types"
	"github.com/deep2chain/sscq/x/upgrade"
	upgcli "github.com/deep2chain/sscq/x/upgrade/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagDetail = "detail"
)

func GetInfoCmd(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "upgrade",
		Short:   "query the information of upgrade module",
		Example: "hscli upgrade info",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res_currentVersion, err := cliCtx.QueryStore(sdk.CurrentVersionKey, sdk.MainStore)
			if(err != nil){
				return err
			}

			var currentVersion uint64
			cdc.MustUnmarshalBinaryLengthPrefixed(res_currentVersion, &currentVersion)

			res_proposalID, err := cliCtx.QueryStore(upgrade.GetSuccessVersionKey(currentVersion), storeName)
			if(err != nil) {
				return err
			}

			var proposalID uint64
			cdc.MustUnmarshalBinaryLengthPrefixed(res_proposalID, &proposalID)

			res_currentVersionInfo, err := cliCtx.QueryStore(upgrade.GetProposalIDKey(proposalID), storeName)
			var currentVersionInfo upgrade.VersionInfo
			cdc.MustUnmarshalBinaryLengthPrefixed(res_currentVersionInfo, &currentVersionInfo)

			res_upgradeInProgress, err := cliCtx.QueryStore(sdk.UpgradeConfigKey, sdk.MainStore)
			if(err != nil){
				return err
			}

			var upgradeInProgress sdk.UpgradeConfig
			if err == nil && len(res_upgradeInProgress) != 0 {
				cdc.MustUnmarshalBinaryLengthPrefixed(res_upgradeInProgress, &upgradeInProgress)
			}

			res_LastFailedVersion, err := cliCtx.QueryStore(sdk.LastFailedVersionKey, sdk.MainStore)
			var lastFailedVersion uint64
			if err == nil && len(res_LastFailedVersion) != 0 {
				cdc.MustUnmarshalBinaryLengthPrefixed(res_LastFailedVersion, &lastFailedVersion)
			} else {
				lastFailedVersion = 0
			}

			upgradeInfoOutput := upgcli.NewUpgradeInfoOutput(currentVersionInfo, lastFailedVersion, upgradeInProgress)

			return cliCtx.PrintOutput(upgradeInfoOutput)
		},
	}
	return cmd
}

func GetCmdQuerySignals(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-signals",
		Short:   "query the information of signals",
		Example: "hscli upgrade status",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res_upgradeConfig, err := cliCtx.QueryStore(sdk.UpgradeConfigKey, sdk.MainStore)
			if err != nil {
				return err
			}
			if len(res_upgradeConfig) == 0 {
				fmt.Println("No Software Upgrade Switch Period is in process.")
				return err
			}

			var upgradeConfig sdk.UpgradeConfig
			if err = cdc.UnmarshalBinaryLengthPrefixed(res_upgradeConfig, &upgradeConfig); err != nil {
				return err
			}

			validatorConsAddrs := make(map[string]bool)
			res, err := cliCtx.QuerySubspace(upgrade.GetSignalPrefixKey(upgradeConfig.Protocol.Version), storeName)
			if err != nil {
				return err
			}

			for _, kv := range res {
				validatorConsAddrs[upgrade.GetAddressFromSignalKey(kv.Key)] = true
			}

			if len(validatorConsAddrs) == 0 {
				fmt.Println("No validator has started the new version.")
				return nil
			}

			key := stake.ValidatorsKey
			resKVs, err := cliCtx.QuerySubspace(key, "stake")
			if err != nil {
				return err
			}

			isDetail := viper.GetBool(flagDetail)
			totalVotingPower := sdk.ZeroDec()
			signalsVotingPower := sdk.ZeroDec()

			for _, kv := range resKVs {
				validator := types.MustUnmarshalValidator(cdc, kv.Value)
				totalVotingPower = totalVotingPower.Add(validator.GetBondedTokens().ToDec())
				if _, ok := validatorConsAddrs[validator.GetConsAddr().String()]; ok {
					signalsVotingPower = signalsVotingPower.Add(validator.GetBondedTokens().ToDec())
					if isDetail {
						fmt.Println(validator.GetOperator().String(), " ", validator.GetBondedTokens().ToDec().String())
					}
				}
			}
			fmt.Println("signalsVotingPower/totalVotingPower = " + signalsVotingPower.Quo(totalVotingPower).String())
			return nil
		},
	}
	cmd.Flags().Bool(flagDetail, false, "details of siganls")
	return cmd
}
