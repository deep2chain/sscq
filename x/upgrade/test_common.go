package upgrade

import (
	"encoding/hex"
	"os"
	"testing"

	"github.com/deep2chain/sscq/codec"
	"github.com/deep2chain/sscq/store"
	sdk "github.com/deep2chain/sscq/types"
	"github.com/deep2chain/sscq/x/auth"
	"github.com/deep2chain/sscq/x/bank"
	"github.com/deep2chain/sscq/x/params"
	stake "github.com/deep2chain/sscq/x/staking"
	stakekeeper "github.com/deep2chain/sscq/x/staking/keeper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/stretchr/testify/require"
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
	initCoins sdk.Int = sdk.NewInt(200)
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
	// bank.RegisterCodec(cdc)
	stake.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

func createTestInput(t *testing.T) (sdk.Context, Keeper, params.Keeper) {
	keyMain := sdk.NewKVStoreKey("main")
	keyAcc := sdk.NewKVStoreKey("acc")
	keyStake := sdk.NewKVStoreKey("stake")
	keyUpgrade := sdk.NewKVStoreKey("upgrade")
	keyParams := sdk.NewKVStoreKey("params")
	tkeyStake := sdk.NewTransientStoreKey("transient_stake")
	tkeyParams := sdk.NewTransientStoreKey("transient_params")

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyMain, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyStake, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyUpgrade, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyStake, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)

	cdc := createTestCodec()
	paramsKeeper := params.NewKeeper(
		cdc,
		keyParams, tkeyParams,
	)

	err := ms.LoadLatestVersion()
	require.Nil(t, err)
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewTMLogger(os.Stdout))

	AccountKeeper := auth.NewAccountKeeper(cdc, keyAcc, paramsKeeper.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	ck := bank.NewBaseKeeper(AccountKeeper, paramsKeeper.Subspace(bank.DefaultParamspace), bank.DefaultCodespace)

	sk := stake.NewKeeper(
		cdc,
		keyStake, tkeyStake,
		ck, paramsKeeper.Subspace(stake.DefaultParamspace),
		stake.DefaultCodespace,
		stakekeeper.NopMetrics(),
	)
	keeper := NewKeeper(cdc, keyUpgrade, sdk.NewProtocolKeeper(sdk.NewKVStoreKey("main")), sk, NopMetrics())

	return ctx, keeper, paramsKeeper
}
