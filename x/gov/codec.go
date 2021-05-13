package gov

import (
	"github.com/deep2chain/sscq/codec"
)

var msgCdc = codec.New()

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSubmitSoftwareUpgradeProposal{}, "htdf/gov/MsgSubmitSoftwareUpgradeProposal", nil)
	cdc.RegisterConcrete(MsgSubmitProposal{}, "htdf/gov/MsgSubmitProposal", nil)
	cdc.RegisterConcrete(MsgDeposit{}, "htdf/gov/MsgDeposit", nil)
	cdc.RegisterConcrete(MsgVote{}, "htdf/gov/MsgVote", nil)

	cdc.RegisterInterface((*ProposalContent)(nil), nil)
	cdc.RegisterConcrete(&Proposal{}, "htdf/gov/Proposal", nil)
	cdc.RegisterConcrete(&SoftwareUpgradeProposal{}, "htdf/gov/SoftwareUpgradeProposal", nil)
}

func init() {
	RegisterCodec(msgCdc)
}
