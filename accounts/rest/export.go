package rest

import (
	"encoding/json"
	"fmt"
	"github.com/deep2chain/sscq/accounts/cli"
	"github.com/deep2chain/sscq/accounts/keystore"
	"github.com/deep2chain/sscq/client/context"
	"github.com/deep2chain/sscq/codec"
	"github.com/deep2chain/sscq/types/rest"

	"net/http"
)

type ExportAccountBody struct {
	Address  string `json:"address"`
	Password string `json:"password"`
}

type ResultExportAccount struct {
	PrivateKey string `json:"private_key"`
}

func ExportAccountRequestHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ExportAccountBody

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		ksw := keystore.NewKeyStoreWallet(keystore.DefaultKeyStoreHome())

		priv, err := cli.GetPrivateKey(ksw, req.Address, req.Password)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var result ResultExportAccount
		result.PrivateKey = priv

		fmt.Printf("result.PrivateKey=%s\n", result.PrivateKey)

		rest.PostProcessResponse(w, cdc, &result, cliCtx.Indent)

		return
	}

}
