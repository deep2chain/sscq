package service

import (
	"github.com/deep2chain/sscq/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSvcDef{}, "sscq/service/MsgSvcDef", nil)
	cdc.RegisterConcrete(MsgSvcBind{}, "sscq/service/MsgSvcBinding", nil)
	cdc.RegisterConcrete(MsgSvcBindingUpdate{}, "sscq/service/MsgSvcBindingUpdate", nil)
	cdc.RegisterConcrete(MsgSvcDisable{}, "sscq/service/MsgSvcDisable", nil)
	cdc.RegisterConcrete(MsgSvcEnable{}, "sscq/service/MsgSvcEnable", nil)
	cdc.RegisterConcrete(MsgSvcRefundDeposit{}, "sscq/service/MsgSvcRefundDeposit", nil)
	cdc.RegisterConcrete(MsgSvcRequest{}, "sscq/service/MsgSvcRequest", nil)
	cdc.RegisterConcrete(MsgSvcResponse{}, "sscq/service/MsgSvcResponse", nil)
	cdc.RegisterConcrete(MsgSvcRefundFees{}, "sscq/service/MsgSvcRefundFees", nil)
	cdc.RegisterConcrete(MsgSvcWithdrawFees{}, "sscq/service/MsgSvcWithdrawFees", nil)
	cdc.RegisterConcrete(MsgSvcWithdrawTax{}, "sscq/service/MsgSvcWithdrawTax", nil)

	cdc.RegisterConcrete(SvcDef{}, "sscq/service/SvcDef", nil)
	cdc.RegisterConcrete(MethodProperty{}, "sscq/service/MethodProperty", nil)
	cdc.RegisterConcrete(SvcBinding{}, "sscq/service/SvcBinding", nil)
	cdc.RegisterConcrete(SvcRequest{}, "sscq/service/SvcRequest", nil)
	cdc.RegisterConcrete(SvcResponse{}, "sscq/service/SvcResponse", nil)
	cdc.RegisterConcrete(IncomingFee{}, "sscq/service/IncomingFee", nil)
	cdc.RegisterConcrete(ReturnedFee{}, "sscq/service/ReturnedFee", nil)

	cdc.RegisterConcrete(&Params{}, "sscq/service/Params", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
