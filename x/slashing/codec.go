package slashing

import (
	"github.com/deep2chain/sscq/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgUnjail{}, "sscq/MsgUnjail", nil)
}

var cdcEmpty = codec.New()
