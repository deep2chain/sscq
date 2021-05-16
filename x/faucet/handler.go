package faucet

import (
	"encoding/json"
	"math/big"
	"os"

	"github.com/deep2chain/sscq/evm/state"
	sdk "github.com/deep2chain/sscq/types"
	"github.com/deep2chain/sscq/x/auth"
	"github.com/deep2chain/sscq/x/bank"
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

//
type SendTxResp struct {
	ErrCode         sdk.CodeType `json:"code"`
	ErrMsg          string       `json:"message"`
	ContractAddress string       `json:"contract_address"`
	EvmOutput       string       `json:"evm_output"`
}

//
func (rsp SendTxResp) String() string {
	rsp.ErrMsg = sdk.GetErrMsg(rsp.ErrCode)
	data, _ := json.Marshal(&rsp)
	return string(data)
}

// New HTDF Message Handler
// connected to handler.go
// HandleMsgSend, HandleMsgAdd upgraded to EVM version
// commented by junying, 2019-08-21
func NewHandler(bankKeeper bank.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {

		switch msg := msg.(type) {
		case MsgAdd:
			return HandleMsgAdd(ctx, bankKeeper, msg)
		default:
			return HandleUnknownMsg(msg)
		}
	}

}

// junying-todo, 2019-08-26
func HandleUnknownMsg(msg sdk.Msg) sdk.Result {
	var sendTxResp SendTxResp
	log.Debugf("msgType error|mstType=%v\n", msg.Type())
	sendTxResp.ErrCode = sdk.ErrCode_Param
	return sdk.Result{Code: sendTxResp.ErrCode, Log: sendTxResp.String()}
}

// Handle a message to add
func HandleMsgAdd(ctx sdk.Context, keeper bank.Keeper, msg MsgAdd) sdk.Result {
	CurSystemIssuer, err := GetSystemIssuerFromRoot()
	if err != nil {
		return sdk.NewError("htdfservice", 101, "system_issuer failed to be found or genesis.json doesn't exists").Result()
	}

	if !msg.SystemIssuer.Equals(CurSystemIssuer) {
		return sdk.NewError("htdfservice", 101, "requester is not the system_issuer").Result()
	}

	_, tags, err := keeper.AddCoins(ctx, msg.SystemIssuer, msg.Amount)
	if err != nil {
		return sdk.NewError("htdfservice", 101, "keeper failed to add requested amount").Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

//
func FeeCollecting(ctx sdk.Context,
	feeCollectionKeeper auth.FeeCollectionKeeper,
	stateDB *state.CommitStateDB,
	gasused uint64,
	gasprice *big.Int) {
	gasUsed := new(big.Int).Mul(new(big.Int).SetUint64(gasused), gasprice)
	log.Debugf("FeeCollecting:gasUsed=%s\n", gasUsed.String())
	feeCollectionKeeper.AddCollectedFees(ctx, sdk.Coins{sdk.NewCoin(sdk.DefaultDenom, sdk.NewIntFromBigInt(gasUsed))})
	stateDB.Commit(false)
	log.Debugln("FeeCollecting:stateDB commited!")
}
