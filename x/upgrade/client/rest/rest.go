package rest

import (
	"fmt"

	"github.com/deep2chain/sscq/client/context"
	"github.com/deep2chain/sscq/codec"

	"github.com/gorilla/mux"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc(fmt.Sprintf("/upgrade_info"), QueryUpgradeInfoRequestHandlerFn(cliCtx, cdc)).Methods("GET")
}
