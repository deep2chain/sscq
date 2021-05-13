package upgrade

import (
	sdk "github.com/deep2chain/sscq/types"
)

func tally(ctx sdk.Context, versionProtocol uint64, k Keeper, threshold sdk.Dec) (passes bool) {

	totalVotingPower := sdk.ZeroInt()
	signalsVotingPower := sdk.ZeroInt()
	k.sk.IterateBondedValidatorsByPower(ctx, func(index int64, validator sdk.Validator) (stop bool) {
		totalVotingPower = totalVotingPower.AddRaw(validator.GetTendermintPower())
		valAcc := validator.GetConsAddr().String()
		if ok := k.GetSignal(ctx, versionProtocol, valAcc); ok {
			signalsVotingPower = signalsVotingPower.AddRaw(validator.GetTendermintPower())
		}
		return true
	})

	ctx.Logger().Info("Tally Start", "SiganlsVotingPower", signalsVotingPower.String(),
		"TotalVotingPower", totalVotingPower.String(),
		"SiganlsVotingPower/TotalVotingPower", signalsVotingPower.Quo(totalVotingPower).String(),
		"Threshold", threshold.String())
	// If more than 95% of validator update , do switch
	if signalsVotingPower.Quo(totalVotingPower).GT(threshold.RoundInt()) {
		return true
	}
	return true
}
