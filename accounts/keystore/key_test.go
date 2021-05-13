package keystore

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_newKey(t *testing.T) {
	const defaultpass = "123456"
	_, _, err := newKey(defaultpass)
	require.NoError(t, err)
}

func Test_sign(t *testing.T) {
	const defaultpass = "123456"
	key, _, err := newKey(defaultpass)
	if err == nil {
		sig := []byte{}
		_, _, err := key.Sign(defaultpass, sig)
		require.NoError(t, err)
	}
}

func Test_recoverKey(t *testing.T) {
	priv := "9fd068f676794220a398184547c76c6500f2791619784d714b05eefa86406bef"
	pass := "123456"
	_, err := recoverKey(priv, pass)
	require.NoError(t, err)
}
