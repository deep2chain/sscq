package sscqservice

import (
	"encoding/hex"
	"strings"

	"github.com/deep2chain/sscq/x/auth"
	"github.com/tendermint/go-amino"
)

//
func Decode_Hex(str string) ([]byte, error) {
	b, err := hex.DecodeString(strings.Replace(str, " ", "", -1))
	if err != nil {
		//panic(fmt.Sprintf("invalid hex string: %q", str))
		return nil, err
	}
	return b, nil
}

//
func Encode_Hex(str []byte) string {
	return hex.EncodeToString(str)
}

// Read and decode a StdTx from rawdata
func ReadStdTxFromRawData(cdc *amino.Codec, str string) (stdTx auth.StdTx, err error) {
	bytes, err := Decode_Hex(str)
	if err = cdc.UnmarshalJSON(bytes, &stdTx); err != nil {
		return stdTx, err
	}
	return stdTx, err
}

// Read and decode a StdTx from rawdata
func ReadStdTxFromString(cdc *amino.Codec, str string) (stdTx auth.StdTx, err error) {
	bytes := []byte(str)
	if err = cdc.UnmarshalJSON(bytes, &stdTx); err != nil {
		return stdTx, err
	}
	return stdTx, err
}
