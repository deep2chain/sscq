package server

import (
	"fmt"
	"path/filepath"

	"github.com/deep2chain/sscq/crypto/keys"

	"github.com/deep2chain/sscq/accounts/keystore"
	clkeys "github.com/deep2chain/sscq/client/keys"
	sdk "github.com/deep2chain/sscq/types"
)

// GenerateCoinKey returns the address of a public key, along with the secret
// phrase to recover the private key.
func GenerateCoinKey() (sdk.AccAddress, string, error) {

	// generate a private key, with recovery phrase
	info, secret, err := clkeys.NewInMemoryKeyBase().CreateMnemonic(
		"name", keys.English, "pass", keys.Secp256k1)
	if err != nil {
		return sdk.AccAddress([]byte{}), "", err
	}
	addr := info.GetPubKey().Address()
	return sdk.AccAddress(addr), secret, nil
}

// GenerateSaveCoinKey returns the address of a public key, along with the secret
// phrase to recover the private key.
func GenerateSaveCoinKey(clientRoot, keyName, keyPass string,
	overwrite bool) (sdk.AccAddress, string, error) {

	// get the keystore from the client
	keybase, err := clkeys.NewKeyBaseFromDir(clientRoot)
	if err != nil {
		return sdk.AccAddress([]byte{}), "", err
	}

	// ensure no overwrite
	if !overwrite {
		_, err := keybase.Get(keyName)
		if err == nil {
			return sdk.AccAddress([]byte{}), "", fmt.Errorf(
				"key already exists, overwrite is disabled (clientRoot: %s)", clientRoot)
		}
	}

	// generate a private key, with recovery phrase
	info, secret, err := keybase.CreateMnemonic(keyName, keys.English, keyPass, keys.Secp256k1)
	if err != nil {
		return sdk.AccAddress([]byte{}), "", err
	}

	return sdk.AccAddress(info.GetPubKey().Address()), secret, nil
}

// junying-todo-20190420
// GenerateSaveCoinKey returns the address of a public key, along with the secret
// phrase to recover the private key.
func GenerateSaveCoinKeyEx(clientRoot, keyPass string) (sdk.AccAddress, string, error) {
	defaultKeyStoreHome := filepath.Join(clientRoot, "keystores")
	// generate a private key, with recovery phrase
	ks := keystore.NewKeyStore(defaultKeyStoreHome)
	secret, err := ks.NewKey(keyPass)
	if err != nil {
		return sdk.AccAddress([]byte{}), "", err
	}
	accaddr, err := sdk.AccAddressFromBech32(ks.Key().Address)
	if err != nil {
		return sdk.AccAddress([]byte{}), "", err
	}
	return accaddr, secret, nil
}
