package bank

import (
	"github.com/deep2chain/sscq/x/params"
)

const (
	// DefaultParamspace for params keeper
	DefaultParamspace = "bank"
	// DefaultSendEnabled enabled
	DefaultSendEnabled = true
)

// ParamStoreKeySendEnabled is store's key for SendEnabled
var ParamStoreKeySendEnabled = []byte("sendenabled")

// ParamKeyTable type declaration for parameters
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable(
		ParamStoreKeySendEnabled, false,
	)
}
