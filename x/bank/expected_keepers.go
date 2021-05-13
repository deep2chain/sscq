package bank

import (
	sdk "github.com/deep2chain/sscq/types"
)

// expected crisis keeper
type CrisisKeeper interface {
	RegisterRoute(moduleName, route string, invar sdk.Invariant)
}
