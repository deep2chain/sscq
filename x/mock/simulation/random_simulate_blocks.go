package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"runtime/debug"
	"sort"
	"strings"
	"syscall"
	"testing"
	"time"

	sdk "github.com/deep2chain/sscq/types"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	tmtypes "github.com/tendermint/tendermint/types"

	bam "github.com/deep2chain/sscq/x/mock/baseapp"
)

// Simulate tests application by sending random messages.
func Simulate(t *testing.T, app *bam.BaseApp,
	appStateFn func(r *rand.Rand, accs []Account) json.RawMessage,
	ops []WeightedOperation, setups []RandSetup,
	invariants []Invariant, numBlocks int, blockSize int, commit bool) error {

	time := time.Now().UnixNano()
	return SimulateFromSeed(t, app, appStateFn, time, ops, setups, invariants, numBlocks, blockSize, commit)
}

func initChain(r *rand.Rand, params Params, accounts []Account, setups []RandSetup, app *bam.BaseApp,
	appStateFn func(r *rand.Rand, accounts []Account) json.RawMessage) (validators map[string]mockValidator) {
	res := app.InitChain(abci.RequestInitChain{AppStateBytes: appStateFn(r, accounts)})
	validators = make(map[string]mockValidator)
	for _, validator := range res.Validators {
		str := fmt.Sprintf("%v", validator.PubKey)
		validators[str] = mockValidator{
			val:           validator,
			livenessState: GetMemberOfInitialState(r, params.InitialLivenessWeightings),
		}
	}

	for i := 0; i < len(setups); i++ {
		setups[i](r, accounts)
	}

	return
}

func randTimestamp(r *rand.Rand) time.Time {
	// json.Marshal breaks for timestamps greater with year greater than 9999
	unixTime := r.Int63n(253373529600)
	return time.Unix(unixTime, 0)
}

// SimulateFromSeed tests an application by running the provided
// operations, testing the provided invariants, but using the provided seed.
func SimulateFromSeed(tb testing.TB, app *bam.BaseApp,
	appStateFn func(r *rand.Rand, accs []Account) json.RawMessage,
	seed int64, ops []WeightedOperation, setups []RandSetup, invariants []Invariant,
	numBlocks int, blockSize int, commit bool) (simError error) {

	// in case we have to end early, don't os.Exit so that we can run cleanup code.
	stopEarly := false
	testingMode, t, b := getTestingMode(tb)
	fmt.Printf("Starting SimulateFromSeed with randomness created with seed %d\n", int(seed))
	r := rand.New(rand.NewSource(seed))
	params := RandomParams(r) // := DefaultParams()
	fmt.Printf("Randomized simulation params: %+v\n", params)
	timestamp := randTimestamp(r)
	fmt.Printf("Starting the simulation from time %v, unixtime %v\n", timestamp.UTC().Format(time.UnixDate), timestamp.Unix())
	timeDiff := maxTimePerBlock - minTimePerBlock

	accs := RandomAccounts(r, params.NumKeys)

	// Setup event stats
	events := make(map[string]uint)
	event := func(what string) {
		events[what]++
	}

	validators := initChain(r, params, accs, setups, app, appStateFn)
	// Second variable to keep pending validator set (delayed one block since TM 0.24)
	// Initially this is the same as the initial validator set
	nextValidators := validators

	header := abci.Header{Height: 1, Time: timestamp, ProposerAddress: randomProposer(r, validators)}
	opCount := 0

	// Setup code to catch SIGTERM's
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		receivedSignal := <-c
		fmt.Printf("\nExiting early due to %s, on block %d, operation %d\n", receivedSignal, header.Height, opCount)
		simError = fmt.Errorf("Exited due to %s", receivedSignal)
		stopEarly = true
	}()

	var pastTimes []time.Time
	var pastVoteInfos [][]abci.VoteInfo

	request := RandomRequestBeginBlock(r, params, validators, pastTimes, pastVoteInfos, event, header)
	// These are operations which have been queued by previous operations
	operationQueue := make(map[int][]Operation)
	timeOperationQueue := []FutureOperation{}
	var blockLogBuilders []*strings.Builder

	if testingMode {
		blockLogBuilders = make([]*strings.Builder, numBlocks)
	}
	displayLogs := logPrinter(testingMode, blockLogBuilders)
	blockSimulator := createBlockSimulator(testingMode, tb, t, params, event, invariants, ops, operationQueue, timeOperationQueue, numBlocks, blockSize, displayLogs)
	if !testingMode {
		b.ResetTimer()
	} else {
		// Recover logs in case of panic
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Panic with err\n", r)
				stackTrace := string(debug.Stack())
				fmt.Println(stackTrace)
				displayLogs()
				simError = fmt.Errorf("Simulation halted due to panic on block %d", header.Height)
			}
		}()
	}

	for i := 0; i < numBlocks && !stopEarly; i++ {
		// Log the header time for future lookup
		pastTimes = append(pastTimes, header.Time)
		pastVoteInfos = append(pastVoteInfos, request.LastCommitInfo.Votes)

		// Construct log writer
		logWriter := addLogMessage(testingMode, blockLogBuilders, i)

		// Run the BeginBlock handler
		logWriter("BeginBlock")
		app.BeginBlock(request)

		if testingMode {
			// Make sure invariants hold at beginning of block
			assertAllInvariants(t, app, invariants, "BeginBlock", displayLogs)
		}

		ctx := app.NewContext(false, header)

		// Run queued operations. Ignores blocksize if blocksize is too small
		logWriter("Queued operations")
		numQueuedOpsRan := runQueuedOperations(operationQueue, int(header.Height), tb, r, app, ctx, accs, logWriter, displayLogs, event)
		numQueuedTimeOpsRan := runQueuedTimeOperations(timeOperationQueue, header.Time, tb, r, app, ctx, accs, logWriter, displayLogs, event)
		if testingMode && onOperation {
			// Make sure invariants hold at end of queued operations
			assertAllInvariants(t, app, invariants, "QueuedOperations", displayLogs)
		}

		logWriter("Standard operations")
		operations := blockSimulator(r, app, ctx, accs, header, logWriter)
		opCount += operations + numQueuedOpsRan + numQueuedTimeOpsRan
		if testingMode {
			// Make sure invariants hold at end of block
			assertAllInvariants(t, app, invariants, "StandardOperations", displayLogs)
		}

		res := app.EndBlock(abci.RequestEndBlock{})
		header.Height++
		header.Time = header.Time.Add(time.Duration(minTimePerBlock) * time.Second).Add(time.Duration(int64(r.Intn(int(timeDiff)))) * time.Second)
		header.ProposerAddress = randomProposer(r, validators)
		logWriter("EndBlock")

		if testingMode {
			// Make sure invariants hold at end of block
			assertAllInvariants(t, app, invariants, "EndBlock", displayLogs)
		}
		if commit {
			app.Commit()
		}

		if header.ProposerAddress == nil {
			fmt.Printf("\nSimulation stopped early as all validators have been unbonded, there is nobody left propose a block!\n")
			stopEarly = true
			break
		}

		// Generate a random RequestBeginBlock with the current validator set for the next block
		request = RandomRequestBeginBlock(r, params, validators, pastTimes, pastVoteInfos, event, header)

		// Update the validator set, which will be reflected in the application on the next block
		validators = nextValidators
		nextValidators = updateValidators(tb, r, params, validators, res.ValidatorUpdates, event)
	}
	if stopEarly {
		DisplayEvents(events)
		return
	}
	fmt.Printf("\nSimulation complete. Final height (blocks): %d, final time (seconds), : %v, operations ran %d\n", header.Height, header.Time, opCount)
	DisplayEvents(events)
	return nil
}

type blockSimFn func(
	r *rand.Rand, app *bam.BaseApp, ctx sdk.Context,
	accounts []Account, header abci.Header, logWriter func(string),
) (opCount int)

// Returns a function to simulate blocks. Written like this to avoid constant parameters being passed everytime, to minimize
// memory overhead
func createBlockSimulator(testingMode bool, tb testing.TB, t *testing.T, params Params,
	event func(string), invariants []Invariant,
	ops []WeightedOperation, operationQueue map[int][]Operation, timeOperationQueue []FutureOperation,
	totalNumBlocks int, avgBlockSize int, displayLogs func()) blockSimFn {

	var (
		lastBlocksizeState = 0 // state for [4 * uniform distribution]
		totalOpWeight      = 0
		blocksize          int
	)

	for i := 0; i < len(ops); i++ {
		totalOpWeight += ops[i].Weight
	}
	selectOp := func(r *rand.Rand) Operation {
		x := r.Intn(totalOpWeight)
		for i := 0; i < len(ops); i++ {
			if x <= ops[i].Weight {
				return ops[i].Op
			}
			x -= ops[i].Weight
		}
		// shouldn't happen
		return ops[0].Op
	}

	return func(r *rand.Rand, app *bam.BaseApp, ctx sdk.Context,
		accounts []Account, header abci.Header, logWriter func(string)) (opCount int) {
		fmt.Printf("\rSimulating... block %d/%d, operation %d/%d. ", header.Height, totalNumBlocks, opCount, blocksize)
		lastBlocksizeState, blocksize = getBlockSize(r, params, lastBlocksizeState, avgBlockSize)
		for j := 0; j < blocksize; j++ {
			logUpdate, futureOps, err := selectOp(r)(r, app, ctx, accounts, event)
			logWriter(logUpdate)
			if err != nil {
				displayLogs()
				tb.Fatalf("error on operation %d within block %d, %v", header.Height, opCount, err)
			}

			queueOperations(operationQueue, timeOperationQueue, futureOps)
			if testingMode {
				if onOperation {
					assertAllInvariants(t, app, invariants, fmt.Sprintf("operation: %v", logUpdate), displayLogs)
				}
				if opCount%50 == 0 {
					fmt.Printf("\rSimulating... block %d/%d, operation %d/%d. ", header.Height, totalNumBlocks, opCount, blocksize)
				}
			}
			opCount++
		}
		return opCount
	}
}

func getTestingMode(tb testing.TB) (testingMode bool, t *testing.T, b *testing.B) {
	testingMode = false
	if _t, ok := tb.(*testing.T); ok {
		t = _t
		testingMode = true
	} else {
		b = tb.(*testing.B)
	}
	return
}

// getBlockSize returns a block size as determined from the transition matrix.
// It targets making average block size the provided parameter. The three
// states it moves between are:
// "over stuffed" blocks with average size of 2 * avgblocksize,
// normal sized blocks, hitting avgBlocksize on average,
// and empty blocks, with no txs / only txs scheduled from the past.
func getBlockSize(r *rand.Rand, params Params, lastBlockSizeState, avgBlockSize int) (state, blocksize int) {
	// TODO: Make default blocksize transition matrix actually make the average
	// blocksize equal to avgBlockSize.
	state = params.BlockSizeTransitionMatrix.NextState(r, lastBlockSizeState)
	if state == 0 {
		blocksize = r.Intn(avgBlockSize * 4)
	} else if state == 1 {
		blocksize = r.Intn(avgBlockSize * 2)
	} else {
		blocksize = 0
	}
	return
}

// adds all future operations into the operation queue.
func queueOperations(queuedOperations map[int][]Operation, queuedTimeOperations []FutureOperation, futureOperations []FutureOperation) {
	if futureOperations == nil {
		return
	}
	for _, futureOp := range futureOperations {
		if futureOp.BlockHeight != 0 {
			if val, ok := queuedOperations[futureOp.BlockHeight]; ok {
				queuedOperations[futureOp.BlockHeight] = append(val, futureOp.Op)
			} else {
				queuedOperations[futureOp.BlockHeight] = []Operation{futureOp.Op}
			}
		} else {
			// TODO: Replace with proper sorted data structure, so don't have the copy entire slice
			index := sort.Search(len(queuedTimeOperations), func(i int) bool { return queuedTimeOperations[i].BlockTime.After(futureOp.BlockTime) })
			queuedTimeOperations = append(queuedTimeOperations, FutureOperation{})
			copy(queuedTimeOperations[index+1:], queuedTimeOperations[index:])
			queuedTimeOperations[index] = futureOp
		}
	}
}

// nolint: errcheck
func runQueuedOperations(queueOperations map[int][]Operation, height int, tb testing.TB, r *rand.Rand, app *bam.BaseApp, ctx sdk.Context,
	accounts []Account, logWriter func(string), displayLogs func(), event func(string)) (numOpsRan int) {
	if queuedOps, ok := queueOperations[height]; ok {
		numOps := len(queuedOps)
		for i := 0; i < numOps; i++ {
			// For now, queued operations cannot queue more operations.
			// If a need arises for us to support queued messages to queue more messages, this can
			// be changed.
			logUpdate, _, err := queuedOps[i](r, app, ctx, accounts, event)
			logWriter(logUpdate)
			if err != nil {
				displayLogs()
				tb.FailNow()
			}
		}
		delete(queueOperations, height)
		return numOps
	}
	return 0
}

func runQueuedTimeOperations(queueOperations []FutureOperation, currentTime time.Time, tb testing.TB, r *rand.Rand, app *bam.BaseApp, ctx sdk.Context,
	accounts []Account, logWriter func(string), displayLogs func(), event func(string)) (numOpsRan int) {

	numOpsRan = 0
	for len(queueOperations) > 0 && currentTime.After(queueOperations[0].BlockTime) {
		// For now, queued operations cannot queue more operations.
		// If a need arises for us to support queued messages to queue more messages, this can
		// be changed.
		logUpdate, _, err := queueOperations[0].Op(r, app, ctx, accounts, event)
		logWriter(logUpdate)
		if err != nil {
			displayLogs()
			tb.FailNow()
		}
		queueOperations = queueOperations[1:]
		numOpsRan++
	}
	return numOpsRan
}

func getKeys(validators map[string]mockValidator) []string {
	keys := make([]string, len(validators))
	i := 0
	for key := range validators {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	return keys
}

// randomProposer picks a random proposer from the current validator set
func randomProposer(r *rand.Rand, validators map[string]mockValidator) cmn.HexBytes {
	keys := getKeys(validators)
	if len(keys) == 0 {
		return nil
	}
	key := keys[r.Intn(len(keys))]
	proposer := validators[key].val
	pk, err := tmtypes.PB2TM.PubKey(proposer.PubKey)
	if err != nil {
		panic(err)
	}
	return pk.Address()
}

// RandomRequestBeginBlock generates a list of signing validators according to the provided list of validators, signing fraction, and evidence fraction
// nolint: unparam
func RandomRequestBeginBlock(r *rand.Rand, params Params, validators map[string]mockValidator,
	pastTimes []time.Time, pastVoteInfos [][]abci.VoteInfo, event func(string), header abci.Header) abci.RequestBeginBlock {
	if len(validators) == 0 {
		fmt.Printf("validator is nil\n")
		return abci.RequestBeginBlock{Header: header}
	}
	voteInfos := make([]abci.VoteInfo, len(validators))
	i := 0
	for _, key := range getKeys(validators) {
		mVal := validators[key]
		mVal.livenessState = params.LivenessTransitionMatrix.NextState(r, mVal.livenessState)
		signed := true

		if mVal.livenessState == 1 {
			// spotty connection, 50% probability of success
			// See https://github.com/golang/go/issues/23804#issuecomment-365370418
			// for reasoning behind computing like this
			signed = r.Int63()%2 == 0
		} else if mVal.livenessState == 2 {
			// offline
			signed = false
		}
		if signed {
			event("beginblock/signing/signed")
		} else {
			event("beginblock/signing/missed")
		}
		pubkey, err := tmtypes.PB2TM.PubKey(mVal.val.PubKey)
		if err != nil {
			panic(err)
		}
		voteInfos[i] = abci.VoteInfo{
			Validator: abci.Validator{
				Address: pubkey.Address(),
				Power:   mVal.val.Power,
			},
			SignedLastBlock: signed,
		}
		i++
	}
	// TODO: Determine capacity before allocation
	evidence := make([]abci.Evidence, 0)
	// Anything but the first block
	if len(pastTimes) > 0 {
		for r.Float64() < params.EvidenceFraction {
			height := header.Height
			time := header.Time
			vals := voteInfos
			if r.Float64() < params.PastEvidenceFraction {
				height = int64(r.Intn(int(header.Height) - 1))
				time = pastTimes[height]
				vals = pastVoteInfos[height]
			}
			validator := vals[r.Intn(len(vals))].Validator
			var totalVotingPower int64
			for _, val := range vals {
				totalVotingPower += val.Validator.Power
			}
			evidence = append(evidence, abci.Evidence{
				Type:             tmtypes.ABCIEvidenceTypeDuplicateVote,
				Validator:        validator,
				Height:           height,
				Time:             time,
				TotalVotingPower: totalVotingPower,
			})
			event("beginblock/evidence")
		}
	}
	return abci.RequestBeginBlock{
		Header: header,
		LastCommitInfo: abci.LastCommitInfo{
			Votes: voteInfos,
		},
		ByzantineValidators: evidence,
	}
}

// updateValidators mimicks Tendermint's update logic
// nolint: unparam
func updateValidators(tb testing.TB, r *rand.Rand, params Params, current map[string]mockValidator, updates []abci.ValidatorUpdate, event func(string)) map[string]mockValidator {

	for _, update := range updates {
		str := fmt.Sprintf("%v", update.PubKey)
		switch {
		case update.Power == 0:
			if _, ok := current[str]; !ok {
				tb.Fatalf("tried to delete a nonexistent validator")
			}

			event("endblock/validatorupdates/kicked")
			delete(current, str)
		default:
			// Does validator already exist?
			if mVal, ok := current[str]; ok {
				mVal.val = update
				event("endblock/validatorupdates/updated")
			} else {
				// Set this new validator
				current[str] = mockValidator{update, GetMemberOfInitialState(r, params.InitialLivenessWeightings)}
				event("endblock/validatorupdates/added")
			}
		}
	}

	return current
}
