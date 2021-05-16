package faucet

import (
	"github.com/deep2chain/sscq/codec"
)

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgAdd{}, "sscq/add", nil)
}
