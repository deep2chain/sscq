package transient

import (
	"testing"

	"github.com/deep2chain/sscq/store/types"
	"github.com/stretchr/testify/require"
)

var k, v = []byte("hello"), []byte("world")

func TestTransientStore(t *testing.T) {
	tstore := NewStore()

	require.Nil(t, tstore.Get(k))

	tstore.Set(k, v)

	require.Equal(t, v, tstore.Get(k))

	tstore.Commit([]*types.KVStoreKey{})

	require.Nil(t, tstore.Get(k))
}
