package subspace

import (
	"fmt"
	"strings"
	"testing"

	"github.com/deep2chain/sscq/codec"
	sdk "github.com/deep2chain/sscq/types"
	"github.com/stretchr/testify/require"
)

var (
	TestDefaultParamSpace = "test"
)

type testparams struct {
	i int64
	b bool
}

func (tp *testparams) ParamSetPairs() ParamSetPairs {
	return ParamSetPairs{
		{[]byte("i"), &tp.i},
		{[]byte("b"), &tp.b},
	}
}
func (tp *testparams) GetParamSpace() string {
	return TestDefaultParamSpace
}

func (tp *testparams) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("i: %d\n", tp.i))
	return sb.String()
}

func (tp *testparams) StringFromBytes(cdc *codec.Codec, key string, bytes []byte) (string, error) {
	return "", nil
}

func (tp *testparams) Validate(key string, value string) (interface{}, sdk.Error) {
	return nil, nil
}

func TestKeyTable(t *testing.T) {
	table := NewKeyTable()

	require.Panics(t, func() { table.RegisterType([]byte(""), nil) })
	require.Panics(t, func() { table.RegisterType([]byte("!@#$%"), nil) })
	require.Panics(t, func() { table.RegisterType([]byte("hello,"), nil) })
	require.Panics(t, func() { table.RegisterType([]byte("hello"), nil) })

	require.NotPanics(t, func() { table.RegisterType([]byte("hello"), bool(false)) })
	require.NotPanics(t, func() { table.RegisterType([]byte("world"), int64(0)) })
	require.Panics(t, func() { table.RegisterType([]byte("hello"), bool(false)) })

	require.NotPanics(t, func() { table.RegisterParamSet(&testparams{}) })
	require.Panics(t, func() { table.RegisterParamSet(&testparams{}) })
}
