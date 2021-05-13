package mint

import (
	"math/big"
	"os"

	sdk "github.com/deep2chain/sscq/types"
	log "github.com/sirupsen/logrus"
)

func init() {
	// This decodes a valid hex string into a sepc256k1Pubkey for use in transaction simulation
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

// junying-todo, 2019-07-17
//	6,000,000	25
//	6,000,000	12.5
// 	6,000,000	6.25
//	...
// ex:
//	BlksPerRound = 100
//	rewards+commission+community-pool
//	sscli query distr rewards sscq1zulqmaqlsgrgmagenaqf02p8kfgsuqkdwgwj80
//	121793749706.0satoshi
//	* 4 = 487174998824
//	not true becasue proper get more rewards,that's, different rewards on every node.
//  sscli query distr commission cosmosvaloper1lwjmdnks33xwnmfayc64ycprww49n33mtm92ne
// 	sscli query distr community-pool
const (
	// Block Reward of First Round
	InitialReward = 25 * 100000000 //25sscq = 2500000000satoshi
	// Block Count Per Round
	BlksPerRound = 6000000 //10 //6,000,000
	// Last Round Index with Block Rewards
	LastRoundIndex = 31
)

// junying-todo, 2019-07-15
// single node: 88.2 for delegators, 11.8 for validator(commission)
// commission is validating fee
// commission rate changes?
func calcParams(ctx sdk.Context, k Keeper) (sdk.Dec, sdk.Dec, sdk.Dec) {
	// fetch params
	totalSupply := k.sk.TotalTokens(ctx)
	log.Infoln("totalSupply", totalSupply)
	// block index
	curBlkHeight := ctx.BlockHeight()
	////fmt.Printf("current Blk Height: %d\n", curBlkHeight)
	// roundIndex = curBlkHeight / BlkCountPerRound
	roundIndex := new(big.Int).Div(big.NewInt(curBlkHeight), big.NewInt(BlksPerRound))
	log.Infoln("curBlkHeight:", curBlkHeight)
	log.Infoln("roundIndex:", roundIndex)
	// BlockProvision = 25 / 2**roundIndex
	BlockProvision := sdk.ZeroDec()
	if roundIndex.Cmp(big.NewInt(LastRoundIndex)) == -1 {
		division := new(big.Int).Exp(big.NewInt(2), roundIndex, nil)
		divisionDec, err := sdk.NewDecFromStr(division.String())
		if err == nil {
			BlockProvision = sdk.NewDec(int64(InitialReward)).Quo(divisionDec)
		}
	}
	// AnnualProvisions = fRatio * FirstRoundBlkReward
	AnnualProvisions := BlockProvision.Mul(sdk.NewDec(int64(BlksPerRound)))
	// junying-todo, 2020-02-04
	k.SetReward(ctx, curBlkHeight, BlockProvision.TruncateInt64())
	Inflation := AnnualProvisions.Quo(sdk.NewDecFromInt(totalSupply))
	return AnnualProvisions, Inflation, BlockProvision
}

// Inflate every block, update inflation parameters once per hour
func BeginBlocker(ctx sdk.Context, k Keeper) {

	// fetch stored minter & params
	minter := k.GetMinter(ctx)
	//params := k.GetParams(ctx)

	// recalculate inflation rate
	var provisionAmt sdk.Dec
	minter.AnnualProvisions, minter.Inflation, provisionAmt = calcParams(ctx, k)

	k.SetMinter(ctx, minter)

	// mint coins, add to collected fees, update supply
	//fmt.Printf("AnnualProvisions: %s, Inflation: %s, provisionAmt: %s\n", minter.AnnualProvisions.String(), minter.Inflation.String(), provisionAmt.TruncateInt().String())
	mintedCoin := sdk.NewCoin(sdk.DefaultDenom, provisionAmt.TruncateInt())
	k.fck.AddCollectedFees(ctx, sdk.Coins{mintedCoin})
	k.sk.InflateSupply(ctx, mintedCoin.Amount)

}
