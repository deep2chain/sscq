package rest

import (
	"bytes"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/deep2chain/sscq/client/context"
	clientrest "github.com/deep2chain/sscq/client/rest"
	"github.com/deep2chain/sscq/codec"
	"github.com/deep2chain/sscq/crypto/keys"
	sdk "github.com/deep2chain/sscq/types"
	"github.com/deep2chain/sscq/types/rest"
	"github.com/deep2chain/sscq/x/slashing"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, kb keys.Keybase) {
	r.HandleFunc(
		"/slashing/validators/{validatorAddr}/unjail",
		unjailRequestHandlerFn(cdc, kb, cliCtx),
	).Methods("POST")
}

// Unjail TX body
type UnjailReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
}

func unjailRequestHandlerFn(cdc *codec.Codec, kb keys.Keybase, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		bech32validator := vars["validatorAddr"]

		var req UnjailReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		valAddr, err := sdk.ValAddressFromBech32(bech32validator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !bytes.Equal(fromAddr, valAddr) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own validator address")
			return
		}

		msg := slashing.NewMsgUnjail(valAddr)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
