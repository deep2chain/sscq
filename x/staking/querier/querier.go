package querier

import (
	"fmt"
	"strings"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/deep2chain/sscq/codec"
	sdk "github.com/deep2chain/sscq/types"
	keep "github.com/deep2chain/sscq/x/staking/keeper"
	"github.com/deep2chain/sscq/x/staking/types"
)

// query endpoints supported by the staking Querier
const (
	QueryValidators                    = "validators"
	QueryValidator                     = "validator"
	QueryDelegatorDelegations          = "delegatorDelegations"
	QueryDelegatorDelegationsEx        = "delegatorDelegationsEx"
	QueryDelegatorUnbondingDelegations = "delegatorUnbondingDelegations"
	QueryRedelegations                 = "redelegations"
	QueryValidatorDelegations          = "validatorDelegations"
	QueryValidatorDelegationsEx        = "validatorDelegationsEx"
	QueryValidatorRedelegations        = "validatorRedelegations"
	QueryValidatorUnbondingDelegations = "validatorUnbondingDelegations"
	QueryDelegator                     = "delegator"
	QueryDelegation                    = "delegation"
	QueryDelegationEx                  = "delegationEx"
	QueryUnbondingDelegation           = "unbondingDelegation"
	QueryDelegatorValidators           = "delegatorValidators"
	QueryDelegatorValidator            = "delegatorValidator"
	QueryPool                          = "pool"
	QueryParameters                    = "parameters"
)

// creates a querier for staking REST endpoints
func NewQuerier(k keep.Keeper, cdc *codec.Codec) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryValidators:
			return queryValidators(ctx, cdc, req, k)
		case QueryValidator:
			return queryValidator(ctx, cdc, req, k)
		case QueryValidatorDelegations:
			return queryValidatorDelegations(ctx, cdc, req, k)
		case QueryValidatorDelegationsEx:
			return queryValidatorDelegationsEx(ctx, cdc, req, k)
		case QueryValidatorUnbondingDelegations:
			return queryValidatorUnbondingDelegations(ctx, cdc, req, k)
		case QueryDelegation:
			return queryDelegation(ctx, cdc, req, k)
		case QueryDelegationEx:
			return queryDelegationEx(ctx, cdc, req, k)
		case QueryUnbondingDelegation:
			return queryUnbondingDelegation(ctx, cdc, req, k)
		case QueryDelegatorDelegations:
			return queryDelegatorDelegations(ctx, cdc, req, k)
		case QueryDelegatorDelegationsEx:
			return queryDelegatorDelegationsEx(ctx, cdc, req, k)
		case QueryDelegatorUnbondingDelegations:
			return queryDelegatorUnbondingDelegations(ctx, cdc, req, k)
		case QueryRedelegations:
			return queryRedelegations(ctx, cdc, req, k)
		case QueryDelegatorValidators:
			return queryDelegatorValidators(ctx, cdc, req, k)
		case QueryDelegatorValidator:
			return queryDelegatorValidator(ctx, cdc, req, k)
		case QueryPool:
			return queryPool(ctx, cdc, k)
		case QueryParameters:
			return queryParameters(ctx, cdc, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown staking query endpoint")
		}
	}
}

// defines the params for the following queries:
// - 'custom/staking/delegatorDelegations'
// - 'custom/staking/delegatorUnbondingDelegations'
// - 'custom/staking/delegatorRedelegations'
// - 'custom/staking/delegatorValidators'
type QueryDelegatorParams struct {
	DelegatorAddr sdk.AccAddress
}

func NewQueryDelegatorParams(delegatorAddr sdk.AccAddress) QueryDelegatorParams {
	return QueryDelegatorParams{
		DelegatorAddr: delegatorAddr,
	}
}

// defines the params for the following queries:
// - 'custom/staking/validator'
// - 'custom/staking/validatorDelegations'
// - 'custom/staking/validatorUnbondingDelegations'
// - 'custom/staking/validatorRedelegations'
type QueryValidatorParams struct {
	ValidatorAddr sdk.ValAddress
}

func NewQueryValidatorParams(validatorAddr sdk.ValAddress) QueryValidatorParams {
	return QueryValidatorParams{
		ValidatorAddr: validatorAddr,
	}
}

// defines the params for the following queries:
// - 'custom/staking/delegation'
// - 'custom/staking/unbondingDelegation'
// - 'custom/staking/delegatorValidator'
type QueryBondsParams struct {
	DelegatorAddr sdk.AccAddress
	ValidatorAddr sdk.ValAddress
}

func NewQueryBondsParams(delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress) QueryBondsParams {
	return QueryBondsParams{
		DelegatorAddr: delegatorAddr,
		ValidatorAddr: validatorAddr,
	}
}

// defines the params for the following queries:
// - 'custom/staking/redelegation'
type QueryRedelegationParams struct {
	DelegatorAddr    sdk.AccAddress
	SrcValidatorAddr sdk.ValAddress
	DstValidatorAddr sdk.ValAddress
}

func NewQueryRedelegationParams(delegatorAddr sdk.AccAddress, srcValidatorAddr sdk.ValAddress, dstValidatorAddr sdk.ValAddress) QueryRedelegationParams {
	return QueryRedelegationParams{
		DelegatorAddr:    delegatorAddr,
		SrcValidatorAddr: srcValidatorAddr,
		DstValidatorAddr: dstValidatorAddr,
	}
}

func queryValidators(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, k keep.Keeper) ([]byte, sdk.Error) {
	var params QueryValidatorsParams

	err := cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	stakingParams := k.GetParams(ctx)
	if params.Limit == 0 {
		params.Limit = int(stakingParams.MaxValidators)
	}

	validators := k.GetAllValidators(ctx)
	filteredVals := make([]types.Validator, 0, len(validators))

	for _, val := range validators {
		if strings.ToLower(val.GetStatus().String()) == strings.ToLower(params.Status) {
			filteredVals = append(filteredVals, val)
		}
	}

	// get pagination bounds
	start := (params.Page - 1) * params.Limit
	end := params.Limit + start
	if end >= len(filteredVals) {
		end = len(filteredVals)
	}

	if start >= len(filteredVals) {
		// page is out of bounds
		filteredVals = []types.Validator{}
	} else {
		filteredVals = filteredVals[start:end]
	}

	res, err := codec.MarshalJSONIndent(cdc, filteredVals)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}

	return res, nil
}

func queryValidator(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, k keep.Keeper) (res []byte, err sdk.Error) {
	var params QueryValidatorParams

	errRes := cdc.UnmarshalJSON(req.Data, &params)
	if errRes != nil {
		return []byte{}, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	validator, found := k.GetValidator(ctx, params.ValidatorAddr)
	if !found {
		return []byte{}, types.ErrNoValidatorFound(types.DefaultCodespace)
	}

	res, errRes = codec.MarshalJSONIndent(cdc, validator)
	if errRes != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
	}
	return res, nil
}

func queryValidatorDelegations(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, k keep.Keeper) (res []byte, err sdk.Error) {
	var params QueryValidatorParams

	errRes := cdc.UnmarshalJSON(req.Data, &params)
	if errRes != nil {
		return []byte{}, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	delegations := k.GetValidatorDelegations(ctx, params.ValidatorAddr)

	res, errRes = codec.MarshalJSONIndent(cdc, delegations)
	if errRes != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
	}
	return res, nil
}

// junying-todo, 2020-05-06
func queryValidatorDelegationsEx(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, k keep.Keeper) (res []byte, err sdk.Error) {
	var params QueryValidatorParams

	errRes := cdc.UnmarshalJSON(req.Data, &params)
	if errRes != nil {
		return []byte{}, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	validator, found := k.GetValidator(ctx, params.ValidatorAddr)
	if !found {
		return []byte{}, types.ErrNoValidatorFound(types.DefaultCodespace)
	}
	delegations := k.GetValidatorDelegations(ctx, params.ValidatorAddr)

	var dels types.DelegationsEx
	for _, delegation := range delegations {
		dels = append(dels, types.NewDelegationEx(delegation, validator.TokensFromShares(delegation.Shares).RoundInt64()))
	}

	res, errRes = codec.MarshalJSONIndent(cdc, dels)
	if errRes != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
	}
	return res, nil
}

func queryValidatorUnbondingDelegations(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, k keep.Keeper) (res []byte, err sdk.Error) {
	var params QueryValidatorParams

	errRes := cdc.UnmarshalJSON(req.Data, &params)
	if errRes != nil {
		return []byte{}, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	unbonds := k.GetUnbondingDelegationsFromValidator(ctx, params.ValidatorAddr)

	res, errRes = codec.MarshalJSONIndent(cdc, unbonds)
	if errRes != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
	}
	return res, nil
}

func queryDelegatorDelegations(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, k keep.Keeper) (res []byte, err sdk.Error) {
	var params QueryDelegatorParams

	errRes := cdc.UnmarshalJSON(req.Data, &params)
	if errRes != nil {
		return []byte{}, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	delegations := k.GetAllDelegatorDelegations(ctx, params.DelegatorAddr)

	res, errRes = codec.MarshalJSONIndent(cdc, delegations)
	if errRes != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
	}
	return res, nil
}

// junying-todo, 2020-05-6
func queryDelegatorDelegationsEx(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, k keep.Keeper) (res []byte, err sdk.Error) {
	var params QueryDelegatorParams

	errRes := cdc.UnmarshalJSON(req.Data, &params)
	if errRes != nil {
		return []byte{}, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	delegations := k.GetAllDelegatorDelegations(ctx, params.DelegatorAddr)

	var dels types.DelegationsEx
	for _, delegation := range delegations {
		validator, found := k.GetValidator(ctx, delegation.ValidatorAddress)
		if !found {
			return []byte{}, types.ErrNoValidatorFound(types.DefaultCodespace)
		}
		dels = append(dels, types.NewDelegationEx(delegation, validator.TokensFromShares(delegation.Shares).RoundInt64()))
	}

	res, errRes = codec.MarshalJSONIndent(cdc, dels)
	if errRes != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
	}
	return res, nil
}

func queryDelegatorUnbondingDelegations(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, k keep.Keeper) (res []byte, err sdk.Error) {
	var params QueryDelegatorParams

	errRes := cdc.UnmarshalJSON(req.Data, &params)
	if errRes != nil {
		return []byte{}, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	unbondingDelegations := k.GetAllUnbondingDelegations(ctx, params.DelegatorAddr)

	res, errRes = codec.MarshalJSONIndent(cdc, unbondingDelegations)
	if errRes != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
	}
	return res, nil
}

func queryDelegatorValidators(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, k keep.Keeper) (res []byte, err sdk.Error) {
	var params QueryDelegatorParams

	stakingParams := k.GetParams(ctx)

	errRes := cdc.UnmarshalJSON(req.Data, &params)
	if errRes != nil {
		return []byte{}, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	validators := k.GetDelegatorValidators(ctx, params.DelegatorAddr, stakingParams.MaxValidators)

	res, errRes = codec.MarshalJSONIndent(cdc, validators)
	if errRes != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
	}
	return res, nil
}

func queryDelegatorValidator(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, k keep.Keeper) (res []byte, err sdk.Error) {
	var params QueryBondsParams

	errRes := cdc.UnmarshalJSON(req.Data, &params)
	if errRes != nil {
		return []byte{}, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	validator, err := k.GetDelegatorValidator(ctx, params.DelegatorAddr, params.ValidatorAddr)
	if err != nil {
		return
	}

	res, errRes = codec.MarshalJSONIndent(cdc, validator)
	if errRes != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
	}
	return res, nil
}

func queryDelegation(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, k keep.Keeper) (res []byte, err sdk.Error) {
	var params QueryBondsParams

	errRes := cdc.UnmarshalJSON(req.Data, &params)
	if errRes != nil {
		return []byte{}, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	delegation, found := k.GetDelegation(ctx, params.DelegatorAddr, params.ValidatorAddr)
	if !found {
		return []byte{}, types.ErrNoDelegation(types.DefaultCodespace)
	}

	res, errRes = codec.MarshalJSONIndent(cdc, delegation)
	if errRes != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
	}
	return res, nil
}

// junying-todo, 2020-05-06
func queryDelegationEx(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, k keep.Keeper) (res []byte, err sdk.Error) {
	var params QueryBondsParams

	errRes := cdc.UnmarshalJSON(req.Data, &params)
	if errRes != nil {
		return []byte{}, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	delegation, found := k.GetDelegation(ctx, params.DelegatorAddr, params.ValidatorAddr)
	if !found {
		return []byte{}, types.ErrNoDelegation(types.DefaultCodespace)
	}

	validator, found := k.GetValidator(ctx, delegation.ValidatorAddress)
	if !found {
		return []byte{}, types.ErrNoValidatorFound(types.DefaultCodespace)
	}
	delegationEx := types.NewDelegationEx(delegation, validator.TokensFromShares(delegation.Shares).RoundInt64())

	res, errRes = codec.MarshalJSONIndent(cdc, delegationEx)
	if errRes != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
	}
	return res, nil
}

func queryUnbondingDelegation(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, k keep.Keeper) (res []byte, err sdk.Error) {
	var params QueryBondsParams

	errRes := cdc.UnmarshalJSON(req.Data, &params)
	if errRes != nil {
		return []byte{}, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	unbond, found := k.GetUnbondingDelegation(ctx, params.DelegatorAddr, params.ValidatorAddr)
	if !found {
		return []byte{}, types.ErrNoUnbondingDelegation(types.DefaultCodespace)
	}

	res, errRes = codec.MarshalJSONIndent(cdc, unbond)
	if errRes != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
	}
	return res, nil
}

func queryRedelegations(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, k keep.Keeper) (res []byte, err sdk.Error) {
	var params QueryRedelegationParams

	errRes := cdc.UnmarshalJSON(req.Data, &params)
	if errRes != nil {
		return []byte{}, sdk.ErrUnknownRequest(string(req.Data))
	}

	var redels []types.Redelegation

	if !params.DelegatorAddr.Empty() && !params.SrcValidatorAddr.Empty() && !params.DstValidatorAddr.Empty() {
		redel, found := k.GetRedelegation(ctx, params.DelegatorAddr, params.SrcValidatorAddr, params.DstValidatorAddr)
		if !found {
			return []byte{}, types.ErrNoRedelegation(types.DefaultCodespace)
		}
		redels = []types.Redelegation{redel}
	} else if params.DelegatorAddr.Empty() && !params.SrcValidatorAddr.Empty() && params.DstValidatorAddr.Empty() {
		redels = k.GetRedelegationsFromValidator(ctx, params.SrcValidatorAddr)
	} else {
		redels = k.GetAllRedelegations(ctx, params.DelegatorAddr, params.SrcValidatorAddr, params.DstValidatorAddr)
	}

	res, errRes = codec.MarshalJSONIndent(cdc, redels)
	if errRes != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
	}
	return res, nil
}

func queryPool(ctx sdk.Context, cdc *codec.Codec, k keep.Keeper) (res []byte, err sdk.Error) {
	pool := k.GetPool(ctx)

	res, errRes := codec.MarshalJSONIndent(cdc, pool)
	if errRes != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
	}
	return res, nil
}

func queryParameters(ctx sdk.Context, cdc *codec.Codec, k keep.Keeper) (res []byte, err sdk.Error) {
	params := k.GetParams(ctx)

	res, errRes := codec.MarshalJSONIndent(cdc, params)
	if errRes != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
	}
	return res, nil
}

// QueryValidatorsParams defines the params for the following queries:
// - 'custom/staking/validators'
type QueryValidatorsParams struct {
	Page, Limit int
	Status      string
}

func NewQueryValidatorsParams(page, limit int, status string) QueryValidatorsParams {
	return QueryValidatorsParams{page, limit, status}
}
