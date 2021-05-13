package guardian

import (
	"github.com/deep2chain/sscq/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgAddProfiler{}, "sscq/x/guardian/MsgAddProfiler", nil)
	cdc.RegisterConcrete(MsgAddTrustee{}, "sscq/x/guardian/MsgAddTrustee", nil)
	cdc.RegisterConcrete(MsgDeleteProfiler{}, "sscq/x/guardian/MsgDeleteProfiler", nil)
	cdc.RegisterConcrete(MsgDeleteTrustee{}, "sscq/x/guardian/MsgDeleteTrustee", nil)

	cdc.RegisterConcrete(Guardian{}, "sscq/x/guardian/Guardian", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
