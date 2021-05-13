package ibc

import (
	"github.com/deep2chain/sscq/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgIBCTransfer{}, "sscq/MsgIBCTransfer", nil)
	cdc.RegisterConcrete(MsgIBCReceive{}, "sscq/MsgIBCReceive", nil)
}
