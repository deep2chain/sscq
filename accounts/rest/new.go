package rest

import (
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/deep2chain/sscq/accounts/keystore"
	"github.com/deep2chain/sscq/types/rest"
)

type newaccountBody struct {
	Password string `json:"password"`
}

func NewAccountRequestHandlerFn(w http.ResponseWriter, r *http.Request) {
	var req newaccountBody

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	ks := keystore.NewKeyStore(keystore.DefaultKeyStoreHome())
	_, err = ks.NewKey(req.Password)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("{\"address\": \"%s\"}", ks.Key().Address)))

	return
}
