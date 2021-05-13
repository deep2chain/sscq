package guardian

import (
	"testing"
	"os"

	stake "github.com/deep2chain/sscq/x/staking"
	"github.com/stretchr/testify/require"
	sdk "github.com/deep2chain/sscq/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/deep2chain/sscq/codec"
	"github.com/deep2chain/sscq/x/auth"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/deep2chain/sscq/store"
	"encoding/hex"
)

var (
	pks = []crypto.PubKey{
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB50"),
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB51"),
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB52"),
	}
	addrs = []sdk.AccAddress{
		sdk.AccAddress(pks[0].Address()),
		sdk.AccAddress(pks[1].Address()),
		sdk.AccAddress(pks[2].Address()),
	}
)

func newPubKey(pk string) (res crypto.PubKey) {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}
	var pkEd ed25519.PubKeyEd25519
	copy(pkEd[:], pkBytes[:])
	return pkEd
}

func createTestCodec() *codec.Codec {
	cdc := codec.New()
	sdk.RegisterCodec(cdc)
	RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	stake.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

func createTestInput(t *testing.T) (sdk.Context, Keeper) {
	keyProf := sdk.NewKVStoreKey("guardian")

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyProf, sdk.StoreTypeIAVL, db)

	err := ms.LoadLatestVersion()
	require.Nil(t, err)
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewTMLogger(os.Stdout))
	cdc := createTestCodec()

	keeper := NewKeeper(cdc, keyProf, DefaultCodespace)

	return ctx, keeper
}
