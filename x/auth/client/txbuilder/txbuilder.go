package context

import (
	"fmt"
	"strings"

	crkeys "github.com/deep2chain/sscq/crypto/keys"

	"github.com/deep2chain/sscq/client"
	"github.com/deep2chain/sscq/client/keys"
	sdk "github.com/deep2chain/sscq/types"
	"github.com/deep2chain/sscq/x/auth"

	"github.com/spf13/viper"
)

// TxBuilder implements a transaction context created in SDK modules.
type TxBuilder struct {
	txEncoder          sdk.TxEncoder
	keybase            crkeys.Keybase
	accountNumber      uint64
	sequence           uint64
	gasWanted          uint64
	gasPrice           uint64
	gasAdjustment      float64
	simulateAndExecute bool
	chainID            string
	memo               string
}

// NewTxBuilder returns a new initialized TxBuilder.
func NewTxBuilder(
	txEncoder sdk.TxEncoder, accNumber, seq, gasWanted uint64, gasAdj float64,
	simulateAndExecute bool, chainID, memo string, gasPrice uint64,
) TxBuilder {

	return TxBuilder{
		txEncoder:          txEncoder,
		keybase:            nil,
		accountNumber:      accNumber,
		sequence:           seq,
		gasWanted:          gasWanted,
		gasPrice:           gasPrice,
		gasAdjustment:      gasAdj,
		simulateAndExecute: simulateAndExecute,
		chainID:            chainID,
		memo:               memo,
	}
}

// NewTxBuilderFromCLI returns a new initialized TxBuilder with parameters from
// the command line using Viper.
func NewTxBuilderFromCLI() TxBuilder {
	kb, err := keys.NewKeyBaseFromHomeFlag()
	if err != nil {
		panic(err)
	}
	txbldr := TxBuilder{
		keybase:       kb,
		accountNumber: uint64(viper.GetInt64(client.FlagAccountNumber)),
		sequence:      uint64(viper.GetInt64(client.FlagSequence)),
		// gas:                client.GasFlagVar.Gas, // commented by junying, 2019-10-21, gas_wanted set to 200000 as default here.
		gasWanted:          uint64(viper.GetInt64(client.FlagGasWanted)), // added by junying, 2019-11-07
		gasAdjustment:      viper.GetFloat64(client.FlagGasAdjustment),
		simulateAndExecute: client.GasFlagVar.Simulate,
		chainID:            viper.GetString(client.FlagChainID),
		memo:               viper.GetString(client.FlagMemo),
	}

	txbldr = txbldr.WithGasPrices(uint64(viper.GetInt64(client.FlagGasPrices)))
	txbldr = txbldr.WithGasWanted(uint64(viper.GetInt64(client.FlagGasWanted))) // added by junying, 2019-11-07

	return txbldr
}

// TxEncoder returns the transaction encoder
func (bldr TxBuilder) TxEncoder() sdk.TxEncoder { return bldr.txEncoder }

// AccountNumber returns the account number
func (bldr TxBuilder) AccountNumber() uint64 { return bldr.accountNumber }

// Sequence returns the transaction sequence
func (bldr TxBuilder) Sequence() uint64 { return bldr.sequence }

// Gas returns the gas for the transaction
func (bldr TxBuilder) GasWanted() uint64 { return bldr.gasWanted }

// GasAdjustment returns the gas adjustment
func (bldr TxBuilder) GasAdjustment() float64 { return bldr.gasAdjustment }

// Keybase returns the keybase
func (bldr TxBuilder) Keybase() crkeys.Keybase { return bldr.keybase }

// SimulateAndExecute returns the option to simulate and then execute the transaction
// using the gas from the simulation results
func (bldr TxBuilder) SimulateAndExecute() bool { return bldr.simulateAndExecute }

// ChainID returns the chain id
func (bldr TxBuilder) ChainID() string { return bldr.chainID }

// Memo returns the memo message
func (bldr TxBuilder) Memo() string { return bldr.memo }

// Fees returns the fees for the transaction
// func (bldr TxBuilder) Fees() sdk.Coins { return bldr.fees }?

// GasPrices returns the gas prices set for the transaction, if any.
func (bldr TxBuilder) GasPrice() uint64 { return bldr.gasPrice }

// WithTxEncoder returns a copy of the context with an updated codec.
func (bldr TxBuilder) WithTxEncoder(txEncoder sdk.TxEncoder) TxBuilder {
	bldr.txEncoder = txEncoder
	return bldr
}

// WithChainID returns a copy of the context with an updated chainID.
func (bldr TxBuilder) WithChainID(chainID string) TxBuilder {
	bldr.chainID = chainID
	return bldr
}

// WithGas returns a copy of the context with an updated gas.
func (bldr TxBuilder) WithGasWanted(gasWanted uint64) TxBuilder {
	bldr.gasWanted = gasWanted
	return bldr
}

// // WithFees returns a copy of the context with an updated fee.
// func (bldr TxBuilder) WithFees(fees string) TxBuilder {
// 	parsedFees, err := sdk.ParseCoins(fees)
// 	if err != nil {
// 		panic(err)
// 	}

// 	bldr.fees = parsedFees
// 	return bldr
// }

// WithGasPrices returns a copy of the context with updated gas prices.
func (bldr TxBuilder) WithGasPrices(gasPrice uint64) TxBuilder {
	bldr.gasPrice = gasPrice
	return bldr
}

// WithKeybase returns a copy of the context with updated keybase.
func (bldr TxBuilder) WithKeybase(keybase crkeys.Keybase) TxBuilder {
	bldr.keybase = keybase
	return bldr
}

// WithSequence returns a copy of the context with an updated sequence number.
func (bldr TxBuilder) WithSequence(sequence uint64) TxBuilder {
	bldr.sequence = sequence
	return bldr
}

// WithMemo returns a copy of the context with an updated memo.
func (bldr TxBuilder) WithMemo(memo string) TxBuilder {
	bldr.memo = strings.TrimSpace(memo)
	return bldr
}

// WithAccountNumber returns a copy of the context with an account number.
func (bldr TxBuilder) WithAccountNumber(accnum uint64) TxBuilder {
	bldr.accountNumber = accnum
	return bldr
}

// BuildSignMsg builds a single message to be signed from a TxBuilder given a
// set of messages. It returns an error if a fee is supplied but cannot be
// parsed.
func (bldr TxBuilder) BuildSignMsg(msgs []sdk.Msg) (StdSignMsg, error) {
	chainID := bldr.chainID
	if chainID == "" {
		return StdSignMsg{}, fmt.Errorf("chain ID required but not specified")
	}
	// junying-todo, 2019-11-08
	// converted from fee based to gas*gasprice based
	// if bldr.gasPrices.IsZero() {
	// 	return StdSignMsg{}, errors.New("gasprices can't not be zero")
	// }
	// if bldr.gasWanted <= 0 {
	// 	return StdSignMsg{}, errors.New("gasWanted must be provided")
	// }
	return StdSignMsg{
		ChainID:       bldr.chainID,
		AccountNumber: bldr.accountNumber,
		Sequence:      bldr.sequence,
		Memo:          bldr.memo,
		Msgs:          msgs,
		Fee:           auth.NewStdFee(bldr.gasWanted, bldr.gasPrice), // junying-todo, 2019-11-07
	}, nil
}

// Sign signs a transaction given a name, passphrase, and a single message to
// signed. An error is returned if signing fails.
func (bldr TxBuilder) Sign(name, passphrase string, msg StdSignMsg) ([]byte, error) {
	sig, err := MakeSignature(bldr.keybase, name, passphrase, msg)
	if err != nil {
		return nil, err
	}

	return bldr.txEncoder(auth.NewStdTx(msg.Msgs, msg.Fee, []auth.StdSignature{sig}, msg.Memo))
}

// BuildAndSign builds a single message to be signed, and signs a transaction
// with the built message given a name, passphrase, and a set of messages.
func (bldr TxBuilder) BuildAndSign(name, passphrase string, msgs []sdk.Msg) ([]byte, error) {
	msg, err := bldr.BuildSignMsg(msgs)
	if err != nil {
		return nil, err
	}

	return bldr.Sign(name, passphrase, msg)
}

// BuildTxForSim creates a StdSignMsg and encodes a transaction with the
// StdSignMsg with a single empty StdSignature for tx simulation.
func (bldr TxBuilder) BuildTxForSim(msgs []sdk.Msg) ([]byte, error) {
	signMsg, err := bldr.BuildSignMsg(msgs)
	if err != nil {
		return nil, err
	}

	// the ante handler will populate with a sentinel pubkey
	sigs := []auth.StdSignature{{}}
	return bldr.txEncoder(auth.NewStdTx(signMsg.Msgs, signMsg.Fee, sigs, signMsg.Memo))
}

// SignStdTx appends a signature to a StdTx and returns a copy of it. If append
// is false, it replaces the signatures already attached with the new signature.
func (bldr TxBuilder) SignStdTx(name, passphrase string, stdTx auth.StdTx, appendSig bool) (signedStdTx auth.StdTx, err error) {
	if bldr.chainID == "" {
		return auth.StdTx{}, fmt.Errorf("chain ID required but not specified")
	}

	stdSignature, err := MakeSignature(bldr.keybase, name, passphrase, StdSignMsg{
		ChainID:       bldr.chainID,
		AccountNumber: bldr.accountNumber,
		Sequence:      bldr.sequence,
		Fee:           stdTx.Fee,
		Msgs:          stdTx.GetMsgs(),
		Memo:          stdTx.GetMemo(),
	})
	if err != nil {
		return
	}

	sigs := stdTx.GetSignatures()
	if len(sigs) == 0 || !appendSig {
		sigs = []auth.StdSignature{stdSignature}
	} else {
		sigs = append(sigs, stdSignature)
	}
	signedStdTx = auth.NewStdTx(stdTx.GetMsgs(), stdTx.Fee, sigs, stdTx.GetMemo())
	return
}

// MakeSignature builds a StdSignature given keybase, key name, passphrase, and a StdSignMsg.
func MakeSignature(keybase crkeys.Keybase, name, passphrase string,
	msg StdSignMsg) (sig auth.StdSignature, err error) {
	if keybase == nil {
		keybase, err = keys.NewKeyBaseFromHomeFlag()
		if err != nil {
			return
		}
	}

	sigBytes, pubkey, err := keybase.Sign(name, passphrase, msg.Bytes())
	if err != nil {
		return
	}
	return auth.StdSignature{
		PubKey:    pubkey,
		Signature: sigBytes,
	}, nil
}
