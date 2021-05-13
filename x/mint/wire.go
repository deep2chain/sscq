package mint

import (
	"github.com/deep2chain/sscq/codec"
	sdk "github.com/deep2chain/sscq/types"
)

// Register concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	// Not Register mint codec in app, deprecated now
	//cdc.RegisterConcrete(Minter{}, "sscq/mint/Minter", nil)
	cdc.RegisterConcrete(&Params{}, "sscq/mint/Params", nil)
	cdc.RegisterConcrete(&sdk.Dec{}, "sscq/mint/rewards", nil)
	// cdc.RegisterConcrete(&Params{}, "mint/Params", nil)
	// cdc.RegisterConcrete(&BlockReward{}, "sscq/mint/BlockReward", nil)
	// cdc.RegisterConcrete(&sdk.Dec{}, "types/Dec", nil)
	// cdc.RegisterConcrete(&sdk.Int{}, "types/Int", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
