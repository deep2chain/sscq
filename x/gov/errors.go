//nolint
package gov

import (
	"fmt"

	sdk "github.com/deep2chain/sscq/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeUnknownProposal         sdk.CodeType = 1
	CodeInactiveProposal        sdk.CodeType = 2
	CodeAlreadyActiveProposal   sdk.CodeType = 3
	CodeAlreadyFinishedProposal sdk.CodeType = 4
	CodeAddressNotStaked        sdk.CodeType = 5
	CodeInvalidTitle            sdk.CodeType = 6
	CodeInvalidDescription      sdk.CodeType = 7
	CodeInvalidProposalType     sdk.CodeType = 8
	CodeInvalidVote             sdk.CodeType = 9
	CodeInvalidGenesis          sdk.CodeType = 10
	CodeInvalidProposalStatus   sdk.CodeType = 11
	CodeSwitchPeriodInProcess   sdk.CodeType = 14
	CodeInvalidInput            sdk.CodeType = 17
	CodeInvalidVersion          sdk.CodeType = 18
	CodeInvalidProposal         sdk.CodeType = 19
	CodeNotEnoughInitialDeposit sdk.CodeType = 20
	CodeAlreadyVote             sdk.CodeType = 25
	CodeOnlyValidatorVote       sdk.CodeType = 26
	CodeInvalidUpgradeParams    sdk.CodeType = 28
)

// Error constructors

func ErrUnknownProposal(codespace sdk.CodespaceType, proposalID uint64) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownProposal, fmt.Sprintf("Unknown proposal with id %d", proposalID))
}

func ErrInactiveProposal(codespace sdk.CodespaceType, proposalID uint64) sdk.Error {
	return sdk.NewError(codespace, CodeInactiveProposal, fmt.Sprintf("Inactive proposal with id %d", proposalID))
}

func ErrAlreadyActiveProposal(codespace sdk.CodespaceType, proposalID uint64) sdk.Error {
	return sdk.NewError(codespace, CodeAlreadyActiveProposal, fmt.Sprintf("Proposal %d has been already active", proposalID))
}

func ErrAlreadyFinishedProposal(codespace sdk.CodespaceType, proposalID uint64) sdk.Error {
	return sdk.NewError(codespace, CodeAlreadyFinishedProposal, fmt.Sprintf("Proposal %d has already passed its voting period", proposalID))
}

func ErrAddressNotStaked(codespace sdk.CodespaceType, address sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeAddressNotStaked, fmt.Sprintf("Address %s is not staked and is thus ineligible to vote", address))
}

func ErrInvalidTitle(codespace sdk.CodespaceType, errorMsg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTitle, errorMsg)
}

func ErrInvalidDescription(codespace sdk.CodespaceType, errorMsg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDescription, errorMsg)
}

func ErrInvalidProposalType(codespace sdk.CodespaceType, proposalType ProposalKind) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidProposalType, fmt.Sprintf("Proposal Type '%s' is not valid", proposalType))
}

func ErrInvalidVote(codespace sdk.CodespaceType, voteOption VoteOption) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidVote, fmt.Sprintf("'%v' is not a valid voting option", voteOption))
}

func ErrInvalidGenesis(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidVote, msg)
}

func ErrNotTrustee(codespace sdk.CodespaceType, trustee sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, fmt.Sprintf("[%s] is not a trustee address", trustee))
}

func ErrNotProfiler(codespace sdk.CodespaceType, profiler sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, fmt.Sprintf("[%s] is not a profiler address", profiler))
}

func ErrCodeInvalidVersion(codespace sdk.CodespaceType, version uint64) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidVersion, fmt.Sprintf("Version [%v] in SoftwareUpgradeProposal isn't valid", version))
}
func ErrCodeInvalidSwitchHeight(codespace sdk.CodespaceType, blockHeight uint64, switchHeight uint64) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidVersion, fmt.Sprintf("Protocol switchHeight [%v] in SoftwareUpgradeProposal isn't large than current block height [%v]", switchHeight, blockHeight))
}

func ErrSwitchPeriodInProcess(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSwitchPeriodInProcess, fmt.Sprintf("Software Upgrade Switch Period is in process."))
}

func ErrInvalidUpgradeThreshold(codespace sdk.CodespaceType, Threshold sdk.Dec) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidUpgradeParams, fmt.Sprintf("Invalid Upgrade Threshold( "+Threshold.String()+" ) should be [0.85, 1)"))
}

func ErrAlreadyVote(codespace sdk.CodespaceType, address sdk.AccAddress, proposalID uint64) sdk.Error {
	return sdk.NewError(codespace, CodeAlreadyVote, fmt.Sprintf("Address %s has voted for the proposal [%d]", address, proposalID))
}

func ErrOnlyValidatorVote(codespace sdk.CodespaceType, address sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeOnlyValidatorVote, fmt.Sprintf("Address %s isn't a validator, so can't vote.", address))
}

func ErrNotEnoughInitialDeposit(codespace sdk.CodespaceType, initialDeposit sdk.Coins, minDeposit sdk.Coins) sdk.Error {
	return sdk.NewError(codespace, CodeNotEnoughInitialDeposit, fmt.Sprintf("Initial Deposit [%s] is less than minInitialDeposit [%s]", initialDeposit.String(), minDeposit.String()))
}
