package gov

import (
	"fmt"

	sdk "github.com/deep2chain/sscq/types"
)

func Execute(ctx sdk.Context, gk Keeper, p ProposalContent) (err error) {
	switch p.GetProposalType() {
	case ProposalTypeSoftwareUpgrade:
		return SoftwareUpgradeProposalExecute(ctx, gk, p.(*SoftwareUpgradeProposal))
	}
	return nil
}

func SoftwareUpgradeProposalExecute(ctx sdk.Context, gk Keeper, sp *SoftwareUpgradeProposal) error {

	if _, ok := gk.protocolKeeper.GetUpgradeConfig(ctx); ok {
		ctx.Logger().Info("Execute SoftwareProposal Failure", "info",
			fmt.Sprintf("Software Upgrade Switch Period is in process."))
		return nil
	}
	if !gk.protocolKeeper.IsValidVersion(ctx, sp.ProtocolDefinition.Version) {
		ctx.Logger().Info("Execute SoftwareProposal Failure", "info",
			fmt.Sprintf("version [%v] in SoftwareUpgradeProposal isn't valid ", sp.ProposalID))
		return nil
	}
	if uint64(ctx.BlockHeight())+1 >= sp.ProtocolDefinition.Height {
		ctx.Logger().Info("Execute SoftwareProposal Failure", "info",
			fmt.Sprintf("switch height must be more than blockHeight + 1"))
		return nil
	}

	gk.protocolKeeper.SetUpgradeConfig(ctx, sdk.NewUpgradeConfig(sp.ProposalID, sp.ProtocolDefinition))
	fmt.Println("Execute SoftwareProposal Success")
	ctx.Logger().Info("Execute SoftwareProposal Success")

	return nil
}
