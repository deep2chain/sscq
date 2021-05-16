package faucet

import (
	"encoding/json"

	sdk "github.com/deep2chain/sscq/types"
)

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// MsgAdd defines a Add message ///////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////
type MsgAdd struct {
	SystemIssuer sdk.AccAddress
	Amount       sdk.Coins
}

var _ sdk.Msg = MsgAdd{}

// NewMsgAdd is a constructor function for Msgadd
func NewMsgAdd(addr sdk.AccAddress, amount sdk.Coins) MsgAdd {
	return MsgAdd{
		SystemIssuer: addr,
		Amount:       amount,
	}
}

// Route should return the name of the module
func (msg MsgAdd) Route() string { return "faucet" }

// Type should return the action
func (msg MsgAdd) Type() string { return "add" }

// ValidateBasic runs stateless checks on the message
func (msg MsgAdd) ValidateBasic() sdk.Error {
	if msg.SystemIssuer.Empty() {
		return sdk.ErrInvalidAddress(msg.SystemIssuer.String())
	}
	if !msg.Amount.IsAllPositive() {
		return sdk.ErrInsufficientCoins("Amount must be positive")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgAdd) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgAdd) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.SystemIssuer}
}

// GetStringAddr defines whose fromaddr is required
func (msg MsgAdd) GetSystemIssuerStr() string {
	return sdk.AccAddress.String(msg.SystemIssuer)
}

//
func (msg MsgAdd) GeSystemIssuer() sdk.AccAddress {
	return msg.SystemIssuer
}
