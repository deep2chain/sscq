package simulation

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	baseapp "github.com/deep2chain/sscq/app"
	sdk "github.com/deep2chain/sscq/types"
	"github.com/deep2chain/sscq/x/gov"
	"github.com/deep2chain/sscq/x/simulation"
)

// SimulateSubmittingVotingAndSlashingForProposal simulates creating a msg Submit Proposal
// voting on the proposal, and subsequently slashing the proposal. It is implemented using
// future operations.
// TODO: Vote more intelligently, so we can actually do some checks regarding votes passing or failing
// TODO: Actually check that validator slashings happened
func SimulateSubmittingVotingAndSlashingForProposal(k gov.Keeper) simulation.Operation {
	handler := gov.NewHandler(k)
	// The states are:
	// column 1: All validators vote
	// column 2: 90% vote
	// column 3: 75% vote
	// column 4: 40% vote
	// column 5: 15% vote
	// column 6: noone votes
	// All columns sum to 100 for simplicity, values chosen by @valardragon semi-arbitrarily,
	// feel free to change.
	numVotesTransitionMatrix, _ := simulation.CreateTransitionMatrix([][]int{
		{20, 10, 0, 0, 0, 0},
		{55, 50, 20, 10, 0, 0},
		{25, 25, 30, 25, 30, 15},
		{0, 15, 30, 25, 30, 30},
		{0, 0, 20, 30, 30, 30},
		{0, 0, 0, 10, 10, 25},
	})
	statePercentageArray := []float64{1, .9, .75, .4, .15, 0}
	curNumVotesState := 1
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account) (
		opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		// 1) submit proposal now
		sender := simulation.RandomAcc(r, accs)
		msg, err := simulationCreateMsgSubmitProposal(r, sender)
		if err != nil {
			return simulation.NoOpMsg(), nil, err
		}
		ok := simulateHandleMsgSubmitProposal(msg, handler, ctx)
		opMsg = simulation.NewOperationMsg(msg, ok, "")
		// don't schedule votes if proposal failed
		if !ok {
			return opMsg, nil, nil
		}
		proposalID := k.GetLastProposalID(ctx)
		// 2) Schedule operations for votes
		// 2.1) first pick a number of people to vote.
		curNumVotesState = numVotesTransitionMatrix.NextState(r, curNumVotesState)
		numVotes := int(math.Ceil(float64(len(accs)) * statePercentageArray[curNumVotesState]))
		// 2.2) select who votes and when
		whoVotes := r.Perm(len(accs))
		// didntVote := whoVotes[numVotes:]
		whoVotes = whoVotes[:numVotes]
		votingPeriod := k.GetVotingParams(ctx).VotingPeriod
		fops := make([]simulation.FutureOperation, numVotes+1)
		for i := 0; i < numVotes; i++ {
			whenVote := ctx.BlockHeader().Time.Add(time.Duration(r.Int63n(int64(votingPeriod.Seconds()))) * time.Second)
			fops[i] = simulation.FutureOperation{BlockTime: whenVote, Op: operationSimulateMsgVote(k, accs[whoVotes[i]], proposalID)}
		}
		// 3) Make an operation to ensure slashes were done correctly. (Really should be a future invariant)
		// TODO: Find a way to check if a validator was slashed other than just checking their balance a block
		// before and after.

		return opMsg, fops, nil
	}
}

// SimulateMsgSubmitProposal simulates a msg Submit Proposal
// Note: Currently doesn't ensure that the proposal txt is in JSON form
func SimulateMsgSubmitProposal(k gov.Keeper) simulation.Operation {
	handler := gov.NewHandler(k)
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account) (
		opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		sender := simulation.RandomAcc(r, accs)
		msg, err := simulationCreateMsgSubmitProposal(r, sender)
		if err != nil {
			return simulation.NoOpMsg(), nil, err
		}
		ok := simulateHandleMsgSubmitProposal(msg, handler, ctx)
		opMsg = simulation.NewOperationMsg(msg, ok, "")
		return opMsg, nil, nil
	}
}

func simulateHandleMsgSubmitProposal(msg gov.MsgSubmitProposal, handler sdk.Handler, ctx sdk.Context) (ok bool) {
	ctx, write := ctx.CacheContext()
	ok = handler(ctx, msg).IsOK()
	if ok {
		write()
	}
	return ok
}

func simulationCreateMsgSubmitProposal(r *rand.Rand, sender simulation.Account) (msg gov.MsgSubmitProposal, err error) {
	deposit := randomDeposit(r)
	msg = gov.NewMsgSubmitProposal(
		simulation.RandStringOfLength(r, 5),
		simulation.RandStringOfLength(r, 5),
		gov.ProposalTypeText,
		sender.Address,
		deposit,
	)
	if msg.ValidateBasic() != nil {
		err = fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
	}
	return
}

// SimulateMsgDeposit
func SimulateMsgDeposit(k gov.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account) (
		opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		acc := simulation.RandomAcc(r, accs)
		proposalID, ok := randomProposalID(r, k, ctx)
		if !ok {
			return simulation.NoOpMsg(), nil, nil
		}
		deposit := randomDeposit(r)
		msg := gov.NewMsgDeposit(acc.Address, proposalID, deposit)
		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(), nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}
		ctx, write := ctx.CacheContext()
		ok = gov.NewHandler(k)(ctx, msg).IsOK()
		if ok {
			write()
		}

		opMsg = simulation.NewOperationMsg(msg, ok, "")
		return opMsg, nil, nil
	}
}

// SimulateMsgVote
// nolint: unparam
func SimulateMsgVote(k gov.Keeper) simulation.Operation {
	return operationSimulateMsgVote(k, simulation.Account{}, 0)
}

// nolint: unparam
func operationSimulateMsgVote(k gov.Keeper, acc simulation.Account, proposalID uint64) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account) (
		opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		if acc.Equals(simulation.Account{}) {
			acc = simulation.RandomAcc(r, accs)
		}

		if proposalID < 0 {
			var ok bool
			proposalID, ok = randomProposalID(r, k, ctx)
			if !ok {
				return simulation.NoOpMsg(), nil, nil
			}
		}
		option := randomVotingOption(r)

		msg := gov.NewMsgVote(acc.Address, proposalID, option)
		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(), nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}

		ctx, write := ctx.CacheContext()
		ok := gov.NewHandler(k)(ctx, msg).IsOK()
		if ok {
			write()
		}

		opMsg = simulation.NewOperationMsg(msg, ok, "")
		return opMsg, nil, nil
	}
}

// Pick a random deposit
func randomDeposit(r *rand.Rand) sdk.Coins {
	// TODO Choose based on account balance and min deposit
	amount := int64(r.Intn(20)) + 1
	return sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, amount)}
}

// Pick a random proposal ID
func randomProposalID(r *rand.Rand, k gov.Keeper, ctx sdk.Context) (proposalID uint64, ok bool) {
	lastProposalID := k.GetLastProposalID(ctx)
	if lastProposalID < 1 || lastProposalID == (2<<63-1) {
		return 0, false
	}
	proposalID = uint64(r.Intn(1+int(lastProposalID)) - 1)
	return proposalID, true
}

// Pick a random voting option
func randomVotingOption(r *rand.Rand) gov.VoteOption {
	switch r.Intn(4) {
	case 0:
		return gov.OptionYes
	case 1:
		return gov.OptionAbstain
	case 2:
		return gov.OptionNo
	case 3:
		return gov.OptionNoWithVeto
	}
	panic("should not happen")
}
