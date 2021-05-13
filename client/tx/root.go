package tx

import (
	"github.com/gorilla/mux"

	"github.com/deep2chain/sscq/client/context"
	"github.com/deep2chain/sscq/codec"
)

// register REST routes
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/txs/{hash}", QueryTxRequestHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/txs", QueryTxsByTagsRequestHandlerFn(cliCtx, cdc)).Methods("GET")
	r.HandleFunc("/txs", BroadcastTxRequest(cliCtx, cdc)).Methods("POST")
	r.HandleFunc("/txs/encode", EncodeTxRequestHandlerFn(cdc, cliCtx)).Methods("POST")

	// for mempool query, fix #issue 13 , yqq 2020-12-24
	r.HandleFunc("/mempool/txs/{hash}", QueryMempoolTxRequestHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/mempool/txs", QueryMempoolTxsRequestHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/mempool/txscount", QueryMempoolTxsNumRequestHandlerFn(cdc, cliCtx)).Methods("GET")
}
