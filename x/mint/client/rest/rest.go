package rest

import (
	"github.com/gorilla/mux"

	"github.com/deep2chain/sscq/client/context"
	"github.com/deep2chain/sscq/codec"
)

// RegisterRoutes registers minting module REST handlers on the provided router.
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	registerQueryRoutes(cliCtx, r, cdc)
}
