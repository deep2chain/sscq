package types

const (
	BigDenom     = "htdf"
	DefaultDenom = "satoshi"
)

// BigCoin
//	use BigDenom
type BigCoin struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}
