package keystore

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	bip39 "github.com/cosmos/go-bip39"
	"github.com/deep2chain/sscq/accounts"
	"github.com/deep2chain/sscq/crypto/keys/hd"
	"github.com/deep2chain/sscq/crypto/keys/mintkey"
	sdk "github.com/deep2chain/sscq/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

var (
	ErrNoMatch = errors.New("no key for given address or file")
)

type KeyStore struct {
	key *Key
	url accounts.URL
}

func NewKeyStore(path string) *KeyStore {
	ks := &KeyStore{
		url: accounts.URL{
			Path: path,
		},
	}

	return ks
}

func (ks *KeyStore) NewKey(passphrase string) (string, error) {

	key, str, err := newKey(passphrase)
	if err != nil {
		return "", err
	}

	ks.key = key
	return str, ks.storeKey()
}

func (ks *KeyStore) RecoverKey(strPrivKey string, passPhrase string) error {

	key, err := recoverKey(strPrivKey, passPhrase)
	if err != nil {
		return err
	}

	ks.key = key
	return ks.storeKey()
}

func (ks *KeyStore) RecoverKeyByMnemonic(mnemonic string, bip39Passphrase string, passPhrase string, account, index uint32) error {

	// create seed
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, bip39Passphrase)
	if err != nil {
		return  err
	}
	//
	hdPath := hd.NewFundraiserParams(account, index)
	// create master key and derive first key:
	masterPriv, ch := hd.ComputeMastersFromSeed(seed)
	derivedPriv, err := hd.DerivePrivateKeyForPath(masterPriv, ch, hdPath.String())
	if err != nil {
		return err
	}

	privkey := secp256k1.PrivKeySecp256k1(derivedPriv)
	privArmor := mintkey.EncryptArmorPrivKey(privkey, passPhrase)
	ks.key.PrivKey = privArmor

	pkey := privkey.PubKey()
	pubKey, err := sdk.Bech32ifyAccPub(pkey)
	if err != nil {
		fmt.Println("newKey error for Bech32ifyAccPub !", err)
		return  err
	}
	ks.key.PubKey = pubKey

	accAddr := pkey.Address().String()
	ks.key.Address = accAddr
	err = ks.storeKey()
	return err
}

func (ks *KeyStore) Key() *Key {
	return ks.key
}

func (ks *KeyStore) storeKey() error {
	content, err := json.Marshal(ks.key)
	if err != nil {
		return err
	}
	return writeKeyfile(ks.url.Path, ks.key.Address, content)
}

func (ks *KeyStore) update(acc accounts.Account, passphrase, newPassphrase string)error{

	privKey, err := mintkey.UnarmorDecryptPrivKey(ks.key.PrivKey, passphrase)
	if err != nil {
		return err
	}

	privArmor := mintkey.EncryptArmorPrivKey(privKey, newPassphrase)
	ks.key.PrivKey = privArmor

	err = ks.drop(acc)
	if err != nil {
		return err
	}

	return ks.storeKey()
}

func (ks *KeyStore) drop(acc accounts.Account)error{
	err := os.Remove(acc.URL.Path)
	if err != nil {
		return err
	}

	return nil
}

func getKey(addr string, filename string) (*Key, error) {
	keyJSON := new(Key)
	bz, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bz, &keyJSON)
	if keyJSON.Address != addr {
		return nil, fmt.Errorf("key content mismatch: have address %x, want %x", keyJSON.Address, addr)
	}
	
	return keyJSON, nil
}

func writeKeyfile(path string, addr string, content []byte) error {

	tmpName := keyFileName(addr)
	fileName := joinPath(path, tmpName)
	name, err := writeTemporaryKeyFile(fileName, content)
	if err != nil {
		return err
	}

	return os.Rename(name, fileName)
}

func writeTemporaryKeyFile(file string, content []byte) (string, error) {
	// Create the keystore directory with appropriate permissions
	// in case it is not present yet.
	const dirPerm = 0700 //0700-junying-todo-20190422

	if err := os.MkdirAll(filepath.Dir(file), dirPerm); err != nil {
		return "", err
	}

	f, err := ioutil.TempFile(filepath.Dir(file), "."+filepath.Base(file)+".tmp")
	if err != nil {
		return "", err
	}
	if _, err := f.Write(content); err != nil {
		f.Close()
		os.Remove(f.Name())
		return "", err
	}
	f.Close()
	return f.Name(), nil
}

// keyFileName implements the naming convention for keyfiles:
// UTC--<created_at UTC ISO8601>-<address string>
func keyFileName(keyAddr string) string {
	ts := time.Now().UTC()
	return fmt.Sprintf("UTC--%s--%s", toISO8601(ts), keyAddr)
}

func toISO8601(t time.Time) string {
	var tz string
	name, offset := t.Zone()
	if name == "UTC" {
		tz = "Z"
	} else {
		tz = fmt.Sprintf("%03d00", offset/3600)
	}
	return fmt.Sprintf("%04d-%02d-%02dT%02d-%02d-%02d.%09d%s", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), tz)
}

func joinPath(path, filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}
	return filepath.Join(path, filename)
}
