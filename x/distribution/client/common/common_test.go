package common

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/deep2chain/sscq/client/context"
	"github.com/deep2chain/sscq/codec"
)

func TestQueryDelegationRewardsAddrValidation(t *testing.T) {
	cdc := codec.New()
	ctx := context.NewCLIContext().WithCodec(cdc)
	type args struct {
		delAddr string
		valAddr string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"invalid delegator address", args{"invalid", ""}, nil, true},
		{"empty delegator address", args{"", ""}, nil, true},
		{"invalid validator address", args{"htdf1keyvaa4u5rcjwq3gncvct4hrmq553fpkremp5v", "invalid"}, nil, true},
		{"empty validator address", args{"htdf1keyvaa4u5rcjwq3gncvct4hrmq553fpkremp5v", ""}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := QueryDelegationRewards(ctx, cdc, "", tt.args.delAddr, tt.args.valAddr)
			require.True(t, err != nil, tt.wantErr)
		})
	}
}
