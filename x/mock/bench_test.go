package mock

import (
	"testing"

	sdk "github.com/deep2chain/sscq/types"
	"github.com/deep2chain/sscq/x/auth"
	sscqservice "github.com/deep2chain/sscq/x/core"
	abci "github.com/tendermint/tendermint/abci/types"
)

// getBenchmarkMockApp initializes a mock application for this module, for purposes of benchmarking
// Any long term API support commitments do not apply to this function.
func getBenchmarkMockApp() (*App, error) {
	mapp := NewApp()

	sscqservice.RegisterCodec(mapp.Cdc)
	mapp.Router().AddRoute("sscqservice", []*sdk.KVStoreKey{mapp.KeyAccount}, sscqservice.NewHandler(mapp.AccountKeeper, mapp.FeeKeeper, mapp.KeyStorage, mapp.KeyCode))
	err := mapp.CompleteSetup()
	return mapp, err
}

func BenchmarkOneBankSendTxPerBlock(b *testing.B) {
	benchmarkApp, _ := getBenchmarkMockApp()

	// Add an account at genesis
	acc := &auth.BaseAccount{
		Address: addr1,
		// Some value conceivably higher than the benchmarks would ever go
		Coins: sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 100000000000)},
	}
	accs := []auth.Account{acc}

	// Construct genesis state
	SetGenesis(benchmarkApp, accs)
	// Precompute all txs
	txs := GenSequenceOfTxs([]sdk.Msg{sendMsg1}, []uint64{0}, []uint64{uint64(0)}, b.N, priv1)
	b.ResetTimer()
	// Run this with a profiler, so its easy to distinguish what time comes from
	// Committing, and what time comes from Check/Deliver Tx.
	for i := 0; i < b.N; i++ {
		benchmarkApp.BeginBlock(abci.RequestBeginBlock{})
		x := benchmarkApp.Check(txs[i])
		if !x.IsOK() {
			panic("something is broken in checking transaction")
		}
		benchmarkApp.Deliver(txs[i])
		benchmarkApp.EndBlock(abci.RequestEndBlock{})
		benchmarkApp.Commit()
	}
}
