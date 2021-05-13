package unit_convert

import (
	"fmt"
	"github.com/shopspring/decimal"
	"math/big"
)

//unit shift big.Int to decimal
func RightShift(v big.Int, decimalShift uint8) decimal.Decimal {
	decImput := decimal.NewFromBigInt(&v, 0)
	return decImput.Mul(decimal.New(1, -1*int32(decimalShift)))
}

//unit shift decimal to big.Int
func LeftShift(v decimal.Decimal, decimalShift uint8) *big.Int {
	biRet := big.NewInt(0)

	//decimal转换成big.Int
	_, bRet := biRet.SetString(v.Mul(decimal.New(1, int32(decimalShift))).String(), 10)
	if bRet == false {
		fmt.Printf("SetString error|v=%s", v.String())
		return biRet
	}
	//fmt.Printf("biRet=%s\n", biRet.String())

	return biRet
}
