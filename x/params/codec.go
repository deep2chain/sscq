package params

import (
	"github.com/deep2chain/sscq/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*ParamSet)(nil), nil)
}
