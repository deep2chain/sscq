package faucet

import (
	"os"

	"github.com/deep2chain/sscq/types"
	sdk "github.com/deep2chain/sscq/types"
	"github.com/ethereum/go-ethereum/common"
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

// Query endpoints supported by the core querier
const (
	QueryContract = "contract"
	//
	ZeroAddress = "0000000000000000000000000000000000000000"
	//
	TxGasLimit = 100000
)

//
type MsgTest struct {
	From sdk.AccAddress
}

func NewMsgTest(addr sdk.AccAddress) MsgTest {
	return MsgTest{
		From: addr,
	}
}
func (msg MsgTest) FromAddress() common.Address {
	return types.ToEthAddress(msg.From)
}

func isZeroByte(data []byte) bool {
	for index := 0; index < len(data); index++ {
		if data[index] != 0 {
			return false
		}
	}
	return true
}
