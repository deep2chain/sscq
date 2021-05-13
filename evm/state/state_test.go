package state

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"github.com/deep2chain/sscq/store"
	sdk "github.com/deep2chain/sscq/types"
	"github.com/deep2chain/sscq/utils"
	"github.com/deep2chain/sscq/x/auth"
	"github.com/deep2chain/sscq/x/params"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"os"

	newevmtypes "github.com/deep2chain/sscq/evm/types"

	"github.com/deep2chain/sscq/codec"

	tmlog "github.com/tendermint/tendermint/libs/log"

	"math/big"
	"testing"
)

var (
	accKey     = sdk.NewKVStoreKey("acc")
	authCapKey = sdk.NewKVStoreKey("authCapKey")
	fckCapKey  = sdk.NewKVStoreKey("fckCapKey")
	keyParams  = sdk.NewKVStoreKey("params")
	tkeyParams = sdk.NewTransientStoreKey("transient_params")

	storageKey = sdk.NewKVStoreKey("storage")
	codeKey    = sdk.NewKVStoreKey("code")

	testHash    = utils.StringToHash("zhoushx")
	fromAddress = utils.StringToAddress("UserA")
	toAddress   = utils.StringToAddress("UserB")
	amount      = big.NewInt(0)
	nonce       = uint64(0)
	gasLimit    = big.NewInt(100000)
	coinbase    = fromAddress

	logger = tmlog.NewNopLogger()
)

func newTestCodec1() *codec.Codec {
	cdc := codec.New()
	newevmtypes.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func cleanup(dataDir string) {
	fmt.Printf("cleaning up db dir|dataDir=%s\n", dataDir)
	os.RemoveAll(dataDir)
}

func TestStateDB(t *testing.T) {

	//---------------------stateDB test--------------------------------------
	dataPath := "/tmp/sscqStateDB"
	db := dbm.NewDB("state", dbm.LevelDBBackend, dataPath)

	cdc := newTestCodec1()
	cms := store.NewCommitMultiStore(db)

	cms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, nil)
	cms.MountStoreWithDB(codeKey, sdk.StoreTypeIAVL, nil)
	cms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, nil)

	pk := params.NewKeeper(cdc, keyParams, tkeyParams)
	ak := auth.NewAccountKeeper(cdc, accKey, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)

	cms.MountStoreWithDB(accKey, sdk.StoreTypeIAVL, nil)
	cms.MountStoreWithDB(storageKey, sdk.StoreTypeIAVL, nil)

	cms.SetPruning(store.PruneNothing)

	err := cms.LoadLatestVersion()
	require.NoError(t, err)

	ms := cms.CacheMultiStore()
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	stateDB, err := NewCommitStateDB(ctx, &ak, storageKey, codeKey)
	must(err)

	fmt.Printf("addr=%s|testBalance=%v\n", fromAddress.String(), stateDB.GetBalance(fromAddress))
	stateDB.AddBalance(fromAddress, big.NewInt(1e18))
	fmt.Printf("addr=%s|testBalance=%v\n", fromAddress.String(), stateDB.GetBalance(fromAddress))

	assert.Equal(t, stateDB.GetBalance(fromAddress).String() == "1000000000000000000", true)

	//remove DB dir
	cleanup(dataPath)
}
