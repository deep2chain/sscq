package auth

import (
	"fmt"
	"os"

	"github.com/tendermint/tendermint/crypto"

	codec "github.com/deep2chain/sscq/codec"
	sdk "github.com/deep2chain/sscq/types"
	"github.com/deep2chain/sscq/x/params"
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
	// StoreKey is string representation of the store key for auth
	StoreKey = "acc"

	// FeeStoreKey is a string representation of the store key for fees
	FeeStoreKey = "fee"

	// QuerierRoute is the querier route for acc
	QuerierRoute = StoreKey
)

var (
	// AddressStoreKeyPrefix prefix for account-by-address store
	AddressStoreKeyPrefix = []byte{0x01}

	globalAccountNumberKey = []byte("globalAccountNumber")

	TotalLoosenTokenKey = []byte("totalLoosenToken")

	BurnedTokenKey = []byte("burnedToken")
)

// AccountKeeper encodes/decodes accounts using the go-amino (binary)
// encoding/decoding library.
type AccountKeeper struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey

	// The prototypical Account constructor.
	proto func() Account

	// The codec codec for binary encoding/decoding of accounts.
	cdc *codec.Codec

	paramSubspace params.Subspace
}

// NewAccountKeeper returns a new sdk.AccountKeeper that uses go-amino to
// (binary) encode and decode concrete sdk.Accounts.
// nolint
func NewAccountKeeper(
	cdc *codec.Codec, key sdk.StoreKey, paramstore params.Subspace, proto func() Account,
) AccountKeeper {

	return AccountKeeper{
		key:           key,
		proto:         proto,
		cdc:           cdc,
		paramSubspace: paramstore.WithKeyTable(ParamKeyTable()),
	}
}

// NewAccountWithAddress implements sdk.AccountKeeper.
func (ak AccountKeeper) NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) Account {
	acc := ak.proto()
	err := acc.SetAddress(addr)
	if err != nil {
		// Handle w/ #870
		panic(err)
	}
	err = acc.SetAccountNumber(ak.GetNextAccountNumber(ctx))
	if err != nil {
		// Handle w/ #870
		panic(err)
	}
	return acc
}

// NewAccount creates a new account
func (ak AccountKeeper) NewAccount(ctx sdk.Context, acc Account) Account {
	if err := acc.SetAccountNumber(ak.GetNextAccountNumber(ctx)); err != nil {
		panic(err)
	}
	return acc
}

// AddressStoreKey turn an address to key used to get it from the account store
func AddressStoreKey(addr sdk.AccAddress) []byte {
	return append(AddressStoreKeyPrefix, addr.Bytes()...)
}

// GetAccount implements sdk.AccountKeeper.
func (ak AccountKeeper) GetAccount(ctx sdk.Context, addr sdk.AccAddress) Account {
	store := ctx.KVStore(ak.key)
	bz := store.Get(AddressStoreKey(addr))
	if bz == nil {
		return nil
	}
	acc := ak.decodeAccount(bz)
	return acc
}

// GetAllAccounts returns all accounts in the accountKeeper.
func (ak AccountKeeper) GetAllAccounts(ctx sdk.Context) []Account {
	accounts := []Account{}
	appendAccount := func(acc Account) (stop bool) {
		accounts = append(accounts, acc)
		return false
	}
	ak.IterateAccounts(ctx, appendAccount)
	return accounts
}

// Implements sdk.AccountKeeper.
func (ak AccountKeeper) SetGenesisAccount(ctx sdk.Context, acc Account) {
	ak.IncreaseTotalLoosenToken(ctx, acc.GetCoins())
	ak.SetAccount(ctx, acc)
}

// SetAccount implements sdk.AccountKeeper.
func (ak AccountKeeper) SetAccount(ctx sdk.Context, acc Account) {
	addr := acc.GetAddress()
	store := ctx.KVStore(ak.key)
	bz, err := ak.cdc.MarshalBinaryBare(acc)
	if err != nil {
		panic(err)
	}
	store.Set(AddressStoreKey(addr), bz)
}

// RemoveAccount removes an account for the account mapper store.
// NOTE: this will cause supply invariant violation if called
func (ak AccountKeeper) RemoveAccount(ctx sdk.Context, acc Account) {
	addr := acc.GetAddress()
	store := ctx.KVStore(ak.key)
	store.Delete(AddressStoreKey(addr))
}

// IterateAccounts implements sdk.AccountKeeper.
func (ak AccountKeeper) IterateAccounts(ctx sdk.Context, process func(Account) (stop bool)) {
	store := ctx.KVStore(ak.key)
	iter := sdk.KVStorePrefixIterator(store, AddressStoreKeyPrefix)
	defer iter.Close()
	for {
		if !iter.Valid() {
			return
		}
		val := iter.Value()
		acc := ak.decodeAccount(val)
		if process(acc) {
			return
		}
		iter.Next()
	}
}

// GetPubKey Returns the PubKey of the account at address
func (ak AccountKeeper) GetPubKey(ctx sdk.Context, addr sdk.AccAddress) (crypto.PubKey, sdk.Error) {
	acc := ak.GetAccount(ctx, addr)
	if acc == nil {
		return nil, sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", addr))
	}
	return acc.GetPubKey(), nil
}

// GetSequence Returns the Sequence of the account at address
func (ak AccountKeeper) GetSequence(ctx sdk.Context, addr sdk.AccAddress) (uint64, sdk.Error) {
	acc := ak.GetAccount(ctx, addr)
	if acc == nil {
		return 0, sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", addr))
	}
	return acc.GetSequence(), nil
}

func (ak AccountKeeper) setSequence(ctx sdk.Context, addr sdk.AccAddress, newSequence uint64) sdk.Error {
	acc := ak.GetAccount(ctx, addr)
	if acc == nil {
		return sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", addr))
	}

	if err := acc.SetSequence(newSequence); err != nil {
		panic(err)
	}

	ak.SetAccount(ctx, acc)
	return nil
}

// GetNextAccountNumber Returns and increments the global account number counter
func (ak AccountKeeper) GetNextAccountNumber(ctx sdk.Context) uint64 {
	var accNumber uint64
	store := ctx.KVStore(ak.key)
	bz := store.Get(globalAccountNumberKey)
	if bz == nil {
		accNumber = 0
	} else {
		err := ak.cdc.UnmarshalBinaryLengthPrefixed(bz, &accNumber)
		if err != nil {
			panic(err)
		}
	}

	bz = ak.cdc.MustMarshalBinaryLengthPrefixed(accNumber + 1)
	store.Set(globalAccountNumberKey, bz)

	return accNumber
}

// -----------------------------------------------------------------------------
// Params

// SetParams sets the auth module's parameters.
func (ak AccountKeeper) SetParams(ctx sdk.Context, params Params) {
	ak.paramSubspace.SetParamSet(ctx, &params)
}

// GetParams gets the auth module's parameters.
func (ak AccountKeeper) GetParams(ctx sdk.Context) (params Params) {
	ak.paramSubspace.GetParamSet(ctx, &params)
	return
}

func (am AccountKeeper) GetBurnedToken(ctx sdk.Context) sdk.Coins {
	// read from db
	var burnToken sdk.Coins
	store := ctx.KVStore(am.key)
	bz := store.Get(BurnedTokenKey)
	if bz == nil {
		burnToken = nil
	} else {
		am.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &burnToken)
	}
	return burnToken
}

func (am AccountKeeper) IncreaseBurnedToken(ctx sdk.Context, coins sdk.Coins) {
	// parameter checking
	if coins == nil || !coins.IsValid() {
		return
	}
	burnToken := am.GetBurnedToken(ctx)
	// increase burn token amount
	burnToken = burnToken.Add(coins)
	if !burnToken.IsAllPositive() {
		panic(fmt.Errorf("burn token is negative"))
	}
	// write back to db
	bzNew := am.cdc.MustMarshalBinaryLengthPrefixed(burnToken)
	store := ctx.KVStore(am.key)
	store.Set(BurnedTokenKey, bzNew)
}

func (am AccountKeeper) GetTotalLoosenToken(ctx sdk.Context) sdk.Coins {
	// read from db
	var totalLoosenToken sdk.Coins
	store := ctx.KVStore(am.key)
	bz := store.Get(TotalLoosenTokenKey)
	if bz == nil {
		totalLoosenToken = nil
	} else {
		am.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &totalLoosenToken)
	}
	return totalLoosenToken
}

func (am AccountKeeper) IncreaseTotalLoosenToken(ctx sdk.Context, coins sdk.Coins) {
	// parameter checking
	if coins == nil || !coins.IsValid() {
		return
	}
	// read from db
	totalLoosenToken := am.GetTotalLoosenToken(ctx)
	// increase totalLoosenToken
	totalLoosenToken = totalLoosenToken.Add(coins)
	if !totalLoosenToken.IsAllPositive() {
		panic(fmt.Errorf("total loosen token is overflow"))
	}
	// write back to db
	bzNew := am.cdc.MustMarshalBinaryLengthPrefixed(totalLoosenToken)
	store := ctx.KVStore(am.key)
	store.Set(TotalLoosenTokenKey, bzNew)

	log.Infoln("Execute IncreaseTotalLoosenToken Successed",
		"increaseCoins", coins.String(), "totalLoosenToken", totalLoosenToken.String())
}

func (am AccountKeeper) DecreaseTotalLoosenToken(ctx sdk.Context, coins sdk.Coins) {
	// parameter checking
	if coins == nil || !coins.IsValid() {
		return
	}
	// read from db
	totalLoosenToken := am.GetTotalLoosenToken(ctx)
	// decrease totalLoosenToken
	totalLoosenToken, negative := totalLoosenToken.SafeSub(coins)
	if negative {
		panic(fmt.Errorf("total loosen token is negative"))
	}
	// write back to db
	bzNew := am.cdc.MustMarshalBinaryLengthPrefixed(totalLoosenToken)
	store := ctx.KVStore(am.key)
	store.Set(TotalLoosenTokenKey, bzNew)

	log.Infoln("Execute DecreaseTotalLoosenToken Successed",
		"decreaseCoins", coins.String(), "totalLoosenToken", totalLoosenToken.String())
}

// -----------------------------------------------------------------------------
// Misc.

func (ak AccountKeeper) decodeAccount(bz []byte) (acc Account) {
	err := ak.cdc.UnmarshalBinaryBare(bz, &acc)
	if err != nil {
		panic(err)
	}
	return
}
