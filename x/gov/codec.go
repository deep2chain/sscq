package gov

import (
	"github.com/deep2chain/sscq/codec"
)

var msgCdc = codec.New()

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSubmitSoftwareUpgradeProposal{}, "sscq/gov/MsgSubmitSoftwareUpgradeProposal", nil)
	cdc.RegisterConcrete(MsgSubmitProposal{}, "sscq/gov/MsgSubmitProposal", nil)
	cdc.RegisterConcrete(MsgDeposit{}, "sscq/gov/MsgDeposit", nil)
	cdc.RegisterConcrete(MsgVote{}, "sscq/gov/MsgVote", nil)

	cdc.RegisterInterface((*ProposalContent)(nil), nil)
	cdc.RegisterConcrete(&Proposal{}, "sscq/gov/Proposal", nil)
	cdc.RegisterConcrete(&SoftwareUpgradeProposal{}, "sscq/gov/SoftwareUpgradeProposal", nil)
}

func init() {
	RegisterCodec(msgCdc)
}
