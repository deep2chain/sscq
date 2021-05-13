package rest

import (
	"github.com/gorilla/mux"

	"github.com/deep2chain/sscq/client/context"
	"github.com/deep2chain/sscq/codec"
	"github.com/deep2chain/sscq/crypto/keys"
)

// RegisterRoutes registers staking-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, kb keys.Keybase) {
	registerQueryRoutes(cliCtx, r, cdc)
	registerTxRoutes(cliCtx, r, cdc, kb)
}
