package keystore

import (
	"encoding/hex"
	"fmt"

	"github.com/cosmos/go-bip39"
	"github.com/deep2chain/sscq/crypto/keys"
	"github.com/deep2chain/sscq/crypto/keys/hd"
	"github.com/deep2chain/sscq/crypto/keys/mintkey"
	sdk "github.com/deep2chain/sscq/types"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	flagAccount         string = "account"
	flagIndex           string = "index"
	mnemonicEntropySize int    = 256
)

type Key struct {
	Address string `json:"address"`
	PubKey  string `json:"pubkey"`
	PrivKey string `json:"privkey"`
}

func newKey(passphrase string) (*Key, string, error) {

	s, err := newSeed()
	if err != nil {
		return nil, "", err
	}

	hdPath := hd.NewFundraiserParams(s.fAccount, s.index)
	seed, err := bip39.NewSeedWithErrorChecking(s.mnemonic, s.bip39Passphrase)
	if err != nil {
		return nil, "", err
	}

	// create master key and derive first key:
	masterPriv, ch := hd.ComputeMastersFromSeed(seed)
	derivedPriv, err := hd.DerivePrivateKeyForPath(masterPriv, ch, hdPath.String())
	if err != nil {
		return nil, "", err
	}

	privkey := secp256k1.PrivKeySecp256k1(derivedPriv)

	privArmor := mintkey.EncryptArmorPrivKey(privkey, passphrase)

	pkey := privkey.PubKey()

	accAddr := sdk.AccAddress(pkey.Address().Bytes())
	pubKey, err := sdk.Bech32ifyAccPub(pkey)
	if err != nil {
		fmt.Println("newKey error for Bech32ifyAccPub !", err)
		return nil, "", err
	}

	key := &Key{
		Address: accAddr.String(),
		PubKey:  pubKey,
		PrivKey: privArmor,
	}

	return key, s.mnemonic, err
}

func (k Key) Sign(auth string, msg []byte) (sig []byte, pub crypto.PubKey, err error) {

	privKey, err := mintkey.UnarmorDecryptPrivKey(k.PrivKey, auth)
	if err != nil {
		return nil, nil, err
	}

	sig, err = privKey.Sign(msg)
	if err != nil {
		return nil, nil, err
	}

	pub = privKey.PubKey()

	return sig, pub, nil
}

type Seed struct {
	mnemonic        string
	bip39Passphrase string
	fAccount        uint32
	index           uint32
}

func newSeed() (*Seed, error) {
	entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
	mnemc, err := bip39.NewMnemonic(entropySeed[:])
	if err != nil {
		return nil, err
	}
	seed := &Seed{
		mnemonic:        mnemc,
		bip39Passphrase: keys.DefaultBIP39Passphrase,
		fAccount:        uint32(viper.GetInt(flagAccount)),
		index:           uint32(viper.GetInt(flagIndex)),
	}

	return seed, nil
}

//revocerKey return a Key by hex string private key and passphrase
func recoverKey(strPrivateKey string, passPhrase string) (*Key, error) {
	//validate and conversion
	privKey, vil, err := conversionPrivKey(strPrivateKey)
	if !vil || err != nil {
		return nil, err
	}

	pub := privKey.PubKey()
	pubKey, err := sdk.Bech32ifyAccPub(pub)
	if err != nil {
		fmt.Println("newKey error for Bech32ifyAccPub !", err)
		return nil, err
	}

	address := sdk.AccAddress(pub.Address().Bytes())

	privArmor := mintkey.EncryptArmorPrivKey(privKey, passPhrase)

	key := &Key{
		Address: address.String(),
		PubKey:  pubKey,
		PrivKey: privArmor,
	}

	return key, nil
}

//conversionPrivKey hex string conversion to byte 32
func conversionPrivKey(strPrivateKey string) (crypto.PrivKey, bool, error) {

	tmpPrivKey, err := hex.DecodeString(strPrivateKey)
	if err != nil {
		fmt.Printf("decodeString error|err=%s\n", err)
		return nil, false, err
	}

	// if len(tmpPrivKey) != 32 {
	// 	return nil, false, err
	// }

	var privateKey32 [32]byte
	setBytes(&privateKey32, tmpPrivKey)

	privkey := secp256k1.PrivKeySecp256k1(privateKey32)

	return privkey, true, nil
}

func setBytes(a *[32]byte, b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-32:]
	}
	copy(a[32-len(b):], b)
}