package params

import (
	"testing"

	"github.com/deep2chain/sscq/types"
	sdk "github.com/deep2chain/sscq/types"

	"github.com/deep2chain/sscq/x/params/subspace"
)

// re-export types from subspace
type (
	Subspace         = subspace.Subspace
	ReadOnlySubspace = subspace.ReadOnlySubspace
	ParamSet         = subspace.ParamSet
	ParamSetPairs    = subspace.ParamSetPairs
	KeyTable         = subspace.KeyTable
)

// nolint - re-export functions from subspace
func NewKeyTable(keytypes ...interface{}) KeyTable {
	return subspace.NewKeyTable(keytypes...)
}
func DefaultTestComponents(t *testing.T) (sdk.Context, Subspace, func([]*types.KVStoreKey) sdk.CommitID) {
	return subspace.DefaultTestComponents(t)
}
