package guardian

import (
	"github.com/deep2chain/sscq/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgAddProfiler{}, "htdf/x/guardian/MsgAddProfiler", nil)
	cdc.RegisterConcrete(MsgAddTrustee{}, "htdf/x/guardian/MsgAddTrustee", nil)
	cdc.RegisterConcrete(MsgDeleteProfiler{}, "htdf/x/guardian/MsgDeleteProfiler", nil)
	cdc.RegisterConcrete(MsgDeleteTrustee{}, "htdf/x/guardian/MsgDeleteTrustee", nil)

	cdc.RegisterConcrete(Guardian{}, "htdf/x/guardian/Guardian", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
