package gov

import (
	"fmt"

	sdk "github.com/deep2chain/sscq/types"
	"github.com/deep2chain/sscq/x/gov/tags"
)

// Called every block, process inflation, update validator set
func EndBlocker(ctx sdk.Context, keeper Keeper) sdk.Tags {
	logger := ctx.Logger().With("module", "x/gov")
	resTags := sdk.NewTags()

	inactiveIterator := keeper.InactiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	defer inactiveIterator.Close()
	for ; inactiveIterator.Valid(); inactiveIterator.Next() {
		var proposalID uint64

		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(inactiveIterator.Value(), &proposalID)
		inactiveProposal, ok := keeper.GetProposal(ctx, proposalID)
		if !ok {
			panic(fmt.Sprintf("proposal %d does not exist", proposalID))
		}

		keeper.DeleteProposal(ctx, proposalID)
		keeper.DeleteDeposits(ctx, proposalID) // delete any associated deposits (burned)

		resTags = resTags.AppendTag(tags.ProposalID, fmt.Sprintf("%d", proposalID))
		resTags = resTags.AppendTag(tags.ProposalResult, tags.ActionProposalDropped)

		logger.Info(
			fmt.Sprintf("proposal %d (%s) didn't meet minimum deposit of %s (had only %s); deleted",
				inactiveProposal.GetProposalID(),
				inactiveProposal.GetTitle(),
				keeper.GetDepositParams(ctx).MinDeposit,
				inactiveProposal.GetTotalDeposit(),
			),
		)
	}

	// fetch active proposals whose voting periods have ended (are passed the block time)
	activeIterator := keeper.ActiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	defer activeIterator.Close()
	for ; activeIterator.Valid(); activeIterator.Next() {
		var proposalID uint64

		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(activeIterator.Value(), &proposalID)
		activeProposal, ok := keeper.GetProposal(ctx, proposalID)
		if !ok {
			panic(fmt.Sprintf("proposal %d does not exist", proposalID))
		}
		result, tallyResults, votingVals := tally(ctx, keeper, activeProposal)

		var tagValue string
		if result == PASS {
			keeper.RefundDeposits(ctx, activeProposal.GetProposalID())
			activeProposal.SetStatus(StatusPassed)
			tagValue = tags.ActionProposalPassed
			Execute(ctx, keeper, activeProposal)
		} else if result == REJECT {
			keeper.RefundDeposits(ctx, activeProposal.GetProposalID())
			activeProposal.SetStatus(StatusRejected)
			tagValue = tags.ActionProposalRejected
		} else if result == REJECTVETO {
			keeper.DeleteDeposits(ctx, activeProposal.GetProposalID())
			activeProposal.SetStatus(StatusRejected)
			tagValue = tags.ActionProposalRejected
		}

		activeProposal.SetTallyResult(tallyResults)
		keeper.SetProposal(ctx, activeProposal)
		keeper.RemoveFromActiveProposalQueue(ctx, activeProposal.GetVotingEndTime(), activeProposal.GetProposalID())

		logger.Info(
			fmt.Sprintf(
				"proposal %d (%s) tallied; passed: %v",
				activeProposal.GetProposalID(), activeProposal.GetTitle(), result,
			),
		)

		resTags = resTags.AppendTag(tags.ProposalID, fmt.Sprintf("%d", proposalID))
		resTags = resTags.AppendTag(tags.ProposalResult, tagValue)

		for _, valAddr := range keeper.GetValidatorSet(ctx, proposalID) {
			if _, ok := votingVals[valAddr.String()]; !ok {
				val := keeper.ds.GetValidatorSet().Validator(ctx, valAddr)
				if val != nil && val.GetStatus() == sdk.Bonded {
					keeper.ds.GetValidatorSet().Slash(ctx,
						val.GetConsAddr(),
						ctx.BlockHeight(),
						val.GetBondedTokens().Int64(),
						keeper.GetTallyParams(ctx).Penalty)
				}
			}
		}

		keeper.DeleteValidatorSet(ctx, activeProposal.GetProposalID())
	}
	return resTags
}
