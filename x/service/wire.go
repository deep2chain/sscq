package service

import (
	"github.com/deep2chain/sscq/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSvcDef{}, "htdf/service/MsgSvcDef", nil)
	cdc.RegisterConcrete(MsgSvcBind{}, "htdf/service/MsgSvcBinding", nil)
	cdc.RegisterConcrete(MsgSvcBindingUpdate{}, "htdf/service/MsgSvcBindingUpdate", nil)
	cdc.RegisterConcrete(MsgSvcDisable{}, "htdf/service/MsgSvcDisable", nil)
	cdc.RegisterConcrete(MsgSvcEnable{}, "htdf/service/MsgSvcEnable", nil)
	cdc.RegisterConcrete(MsgSvcRefundDeposit{}, "htdf/service/MsgSvcRefundDeposit", nil)
	cdc.RegisterConcrete(MsgSvcRequest{}, "htdf/service/MsgSvcRequest", nil)
	cdc.RegisterConcrete(MsgSvcResponse{}, "htdf/service/MsgSvcResponse", nil)
	cdc.RegisterConcrete(MsgSvcRefundFees{}, "htdf/service/MsgSvcRefundFees", nil)
	cdc.RegisterConcrete(MsgSvcWithdrawFees{}, "htdf/service/MsgSvcWithdrawFees", nil)
	cdc.RegisterConcrete(MsgSvcWithdrawTax{}, "htdf/service/MsgSvcWithdrawTax", nil)

	cdc.RegisterConcrete(SvcDef{}, "htdf/service/SvcDef", nil)
	cdc.RegisterConcrete(MethodProperty{}, "htdf/service/MethodProperty", nil)
	cdc.RegisterConcrete(SvcBinding{}, "htdf/service/SvcBinding", nil)
	cdc.RegisterConcrete(SvcRequest{}, "htdf/service/SvcRequest", nil)
	cdc.RegisterConcrete(SvcResponse{}, "htdf/service/SvcResponse", nil)
	cdc.RegisterConcrete(IncomingFee{}, "htdf/service/IncomingFee", nil)
	cdc.RegisterConcrete(ReturnedFee{}, "htdf/service/ReturnedFee", nil)

	cdc.RegisterConcrete(&Params{}, "htdf/service/Params", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
