package ibc

import (
	"github.com/deep2chain/sscq/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgIBCTransfer{}, "htdf/MsgIBCTransfer", nil)
	cdc.RegisterConcrete(MsgIBCReceive{}, "htdf/MsgIBCReceive", nil)
}
