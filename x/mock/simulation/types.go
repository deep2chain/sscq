package simulation

import (
	"math/rand"
	"time"

	sdk "github.com/deep2chain/sscq/types"
	bam "github.com/deep2chain/sscq/x/mock/baseapp"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
)

type (
	// Operation runs a state machine transition,
	// and ensures the transition happened as expected.
	// The operation could be running and testing a fuzzed transaction,
	// or doing the same for a message.
	//
	// For ease of debugging,
	// an operation returns a descriptive message "action",
	// which details what this fuzzed state machine transition actually did.
	//
	// Operations can optionally provide a list of "FutureOperations" to run later
	// These will be ran at the beginning of the corresponding block.
	Operation func(r *rand.Rand, app *bam.BaseApp, ctx sdk.Context,
		accounts []Account, event func(string),
	) (action string, futureOperations []FutureOperation, err error)

	// RandSetup performs the random setup the mock module needs.
	RandSetup func(r *rand.Rand, accounts []Account)

	// An Invariant is a function which tests a particular invariant.
	// If the invariant has been broken, it should return an error
	// containing a descriptive message about what happened.
	// The simulator will then halt and print the logs.
	Invariant = sdk.Invariant

	// Account contains a privkey, pubkey, address tuple
	// eventually more useful data can be placed in here.
	// (e.g. number of coins)
	Account struct {
		PrivKey crypto.PrivKey
		PubKey  crypto.PubKey
		Address sdk.AccAddress
	}

	mockValidator struct {
		val           abci.ValidatorUpdate
		livenessState int
	}

	// FutureOperation is an operation which will be ran at the
	// beginning of the provided BlockHeight.
	// If both a BlockHeight and BlockTime are specified, it will use the BlockHeight.
	// In the (likely) event that multiple operations are queued at the same
	// block height, they will execute in a FIFO pattern.
	FutureOperation struct {
		BlockHeight int
		BlockTime   time.Time
		Op          Operation
	}

	// WeightedOperation is an operation with associated weight.
	// This is used to bias the selection operation within the simulator.
	WeightedOperation struct {
		Weight int
		Op     Operation
	}
)

// TODO remove? not being called anywhere
// PeriodicInvariant returns an Invariant function closure that asserts
// a given invariant if the mock application's last block modulo the given
// period is congruent to the given offset.
func PeriodicInvariant(invariant Invariant, period int, offset int) Invariant {
	return func(ctx sdk.Context) error {
		if int(ctx.BlockHeight())%period == offset {
			return invariant(ctx)
		}
		return nil
	}
}

// nolint
func (acc Account) Equals(acc2 Account) bool {
	return acc.Address.Equals(acc2.Address)
}
