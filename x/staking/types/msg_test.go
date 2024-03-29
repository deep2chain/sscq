package types

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/deep2chain/sscq/types"
)

var (
	coinPos  = sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)
	coinZero = sdk.NewInt64Coin(sdk.DefaultBondDenom, 0)
)

// test ValidateBasic for MsgCreateValidator
func TestMsgCreateValidator(t *testing.T) {
	commission1 := NewCommissionMsg(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
	commission2 := NewCommissionMsg(sdk.NewDec(5), sdk.NewDec(5), sdk.NewDec(5))

	tests := []struct {
		name, moniker, identity, website, details string
		commissionMsg                             CommissionMsg
		minSelfDelegation                         sdk.Int
		validatorAddr                             sdk.ValAddress
		pubkey                                    crypto.PubKey
		bond                                      sdk.Coin
		expectPass                                bool
	}{
		{"basic good", "a", "b", "c", "d", commission1, sdk.OneInt(), addr1, pk1, coinPos, true},
		{"partial description", "", "", "c", "", commission1, sdk.OneInt(), addr1, pk1, coinPos, true},
		{"empty description", "", "", "", "", commission2, sdk.OneInt(), addr1, pk1, coinPos, false},
		{"empty address", "a", "b", "c", "d", commission2, sdk.OneInt(), emptyAddr, pk1, coinPos, false},
		{"empty pubkey", "a", "b", "c", "d", commission1, sdk.OneInt(), addr1, emptyPubkey, coinPos, true},
		{"empty bond", "a", "b", "c", "d", commission2, sdk.OneInt(), addr1, pk1, coinZero, false},
		{"zero min self delegation", "a", "b", "c", "d", commission1, sdk.ZeroInt(), addr1, pk1, coinPos, false},
		{"negative min self delegation", "a", "b", "c", "d", commission1, sdk.NewInt(-1), addr1, pk1, coinPos, false},
		{"delegation less than min self delegation", "a", "b", "c", "d", commission1, coinPos.Amount.Add(sdk.OneInt()), addr1, pk1, coinPos, false},
	}

	for _, tc := range tests {
		description := NewDescription(tc.moniker, tc.identity, tc.website, tc.details)
		msg := NewMsgCreateValidator(tc.validatorAddr, tc.pubkey, tc.bond, description, tc.commissionMsg, tc.minSelfDelegation)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}

// test ValidateBasic for MsgEditValidator
func TestMsgEditValidator(t *testing.T) {
	tests := []struct {
		name, moniker, identity, website, details string
		validatorAddr                             sdk.ValAddress
		expectPass                                bool
	}{
		{"basic good", "a", "b", "c", "d", addr1, true},
		{"partial description", "", "", "c", "", addr1, true},
		{"empty description", "", "", "", "", addr1, false},
		{"empty address", "a", "b", "c", "d", emptyAddr, false},
	}

	for _, tc := range tests {
		description := NewDescription(tc.moniker, tc.identity, tc.website, tc.details)
		newRate := sdk.ZeroDec()
		newMinSelfDelegation := sdk.OneInt()

		msg := NewMsgEditValidator(tc.validatorAddr, description, &newRate, &newMinSelfDelegation)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}

// test ValidateBasic for MsgDelegate
func TestMsgDelegate(t *testing.T) {
	tests := []struct {
		name          string
		delegatorAddr sdk.AccAddress
		validatorAddr sdk.ValAddress
		bond          sdk.Coin
		expectPass    bool
	}{
		{"basic good", sdk.AccAddress(addr1), addr2, coinPos, true},
		{"self bond", sdk.AccAddress(addr1), addr1, coinPos, true},
		{"empty delegator", sdk.AccAddress(emptyAddr), addr1, coinPos, false},
		{"empty validator", sdk.AccAddress(addr1), emptyAddr, coinPos, false},
		{"empty bond", sdk.AccAddress(addr1), addr2, coinZero, false},
	}

	for _, tc := range tests {
		msg := NewMsgDelegate(tc.delegatorAddr, tc.validatorAddr, tc.bond)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}

// test ValidateBasic for MsgUnbond
func TestMsgBeginRedelegate(t *testing.T) {
	tests := []struct {
		name             string
		delegatorAddr    sdk.AccAddress
		validatorSrcAddr sdk.ValAddress
		validatorDstAddr sdk.ValAddress
		amount           sdk.Coin
		expectPass       bool
	}{
		{"regular", sdk.AccAddress(addr1), addr2, addr3, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), true},
		{"zero amount", sdk.AccAddress(addr1), addr2, addr3, sdk.NewInt64Coin(sdk.DefaultBondDenom, 0), false},
		{"empty delegator", sdk.AccAddress(emptyAddr), addr1, addr3, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), false},
		{"empty source validator", sdk.AccAddress(addr1), emptyAddr, addr3, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), false},
		{"empty destination validator", sdk.AccAddress(addr1), addr2, emptyAddr, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), false},
	}

	for _, tc := range tests {
		msg := NewMsgBeginRedelegate(tc.delegatorAddr, tc.validatorSrcAddr, tc.validatorDstAddr, tc.amount)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}

// test ValidateBasic for MsgUnbond
func TestMsgUndelegate(t *testing.T) {
	tests := []struct {
		name          string
		delegatorAddr sdk.AccAddress
		validatorAddr sdk.ValAddress
		amount        sdk.Coin
		expectPass    bool
	}{
		{"regular", sdk.AccAddress(addr1), addr2, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), true},
		{"zero amount", sdk.AccAddress(addr1), addr2, sdk.NewInt64Coin(sdk.DefaultBondDenom, 0), false},
		{"empty delegator", sdk.AccAddress(emptyAddr), addr1, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), false},
		{"empty validator", sdk.AccAddress(addr1), emptyAddr, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), false},
	}

	for _, tc := range tests {
		msg := NewMsgUndelegate(tc.delegatorAddr, tc.validatorAddr, tc.amount)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}
