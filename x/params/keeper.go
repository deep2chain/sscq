package params

import (
	"os"

	"github.com/deep2chain/sscq/codec"
	sdk "github.com/deep2chain/sscq/types"

	"github.com/deep2chain/sscq/x/params/subspace"
	log "github.com/sirupsen/logrus"
)

func init() {
	// junying-todo,2020-01-17
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	// LOG_LEVEL not set, let's default to debug
	if !ok {
		lvl = "info" //trace/debug/info/warn/error/parse/fatal/panic
	}
	// parse string, this is built-in feature of logrus
	ll, err := log.ParseLevel(lvl)
	if err != nil {
		ll = log.FatalLevel //TraceLevel/DebugLevel/InfoLevel/WarnLevel/ErrorLevel/ParseLevel/FatalLevel/PanicLevel
	}
	// set global log level
	log.SetLevel(ll)
	log.SetFormatter(&log.TextFormatter{}) //&log.JSONFormatter{})
}

const (
	// StoreKey is the string key for the params store
	StoreKey = subspace.StoreKey

	// TStoreKey is the string key for the params transient store
	TStoreKey = subspace.TStoreKey
)

// Keeper of the global paramstore
type Keeper struct {
	cdc  *codec.Codec
	key  sdk.StoreKey
	tkey sdk.StoreKey

	spaces    map[string]*Subspace
	paramSets map[string]ParamSet
}

// NewKeeper constructs a params keeper
func NewKeeper(cdc *codec.Codec, key *sdk.KVStoreKey, tkey *sdk.TransientStoreKey) (k Keeper) {
	k = Keeper{
		cdc:  cdc,
		key:  key,
		tkey: tkey,

		spaces:    make(map[string]*Subspace),
		paramSets: make(map[string]ParamSet),
	}

	return k
}

// Allocate subspace used for keepers
func (k Keeper) Subspace(spacename string) Subspace {
	_, ok := k.spaces[spacename]
	if ok {
		panic("subspace already occupied")
	}

	if spacename == "" {
		panic("cannot use empty string for subspace")
	}

	space := subspace.NewSubspace(k.cdc, k.key, k.tkey, spacename)
	log.Traceln("=======k.key=======", k.key, "	k.tkey	", k.tkey)
	k.spaces[spacename] = &space

	return space
}

// Get existing substore from keeper
func (k Keeper) GetSubspace(storename string) (Subspace, bool) {
	space, ok := k.spaces[storename]
	if !ok {
		return Subspace{}, false
	}
	return *space, ok
}

func (k Keeper) RegisterParamSet(ps ...ParamSet) {
	RegisterParamSet(k.paramSets, ps...)
}

// Get existing substore from keeper
func (k Keeper) GetParams(paramSpace string) (ParamSet, bool) {
	paramSet, ok := k.paramSets[paramSpace]
	if !ok {
		return nil, false
	}
	return paramSet, ok
}
