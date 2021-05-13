package types

import (
	"github.com/deep2chain/sscq/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateValidator{}, "sscq/MsgCreateValidator", nil)
	cdc.RegisterConcrete(MsgEditValidator{}, "sscq/MsgEditValidator", nil)
	cdc.RegisterConcrete(MsgDelegate{}, "sscq/MsgDelegate", nil)
	cdc.RegisterConcrete(MsgUndelegate{}, "sscq/MsgUndelegate", nil)
	cdc.RegisterConcrete(MsgBeginRedelegate{}, "sscq/MsgBeginRedelegate", nil)
	cdc.RegisterConcrete(MsgSetUndelegateStatus{}, "sscq/MsgSetUndelegateStatus", nil)
}

// generic sealed codec to be used throughout sdk
var MsgCdc *codec.Codec

func init() {
	cdc := codec.New()
	RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	MsgCdc = cdc.Seal()
}
