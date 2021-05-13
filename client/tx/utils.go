package tx

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/deep2chain/sscq/client/context"
	"github.com/deep2chain/sscq/codec"
	sdk "github.com/deep2chain/sscq/types"
	"github.com/deep2chain/sscq/x/auth"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// SearchTxs performs a search for transactions for a given set of tags via
// Tendermint RPC. It returns a slice of Info object containing txs and metadata.
// An error is returned if the query fails.
func SearchTxs(cliCtx context.CLIContext, cdc *codec.Codec, tags []string, page, limit int) ([]sdk.TxResponse, error) {
	if len(tags) == 0 {
		return nil, errors.New("must declare at least one tag to search")
	}

	if page <= 0 {
		return nil, errors.New("page must greater than 0")
	}

	if limit <= 0 {
		return nil, errors.New("limit must greater than 0")
	}

	// XXX: implement ANY
	query := strings.Join(tags, " AND ")

	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	prove := !cliCtx.TrustNode

	resTxs, err := node.TxSearch(query, prove, page, limit)
	if err != nil {
		return nil, err
	}

	if prove {
		for _, tx := range resTxs.Txs {
			err := ValidateTxResult(cliCtx, tx)
			if err != nil {
				return nil, err
			}
		}
	}

	resBlocks, err := getBlocksForTxResults(cliCtx, resTxs.Txs)
	if err != nil {
		return nil, err
	}

	txs, err := formatTxResults(cdc, resTxs.Txs, resBlocks)
	if err != nil {
		return nil, err
	}

	return txs, nil
}

// formatTxResults parses the indexed txs into a slice of TxResponse objects.
func formatTxResults(cdc *codec.Codec, resTxs []*ctypes.ResultTx, resBlocks map[int64]*ctypes.ResultBlock) ([]sdk.TxResponse, error) {
	var err error
	out := make([]sdk.TxResponse, len(resTxs))
	for i := range resTxs {
		out[i], err = formatTxResult(cdc, resTxs[i], resBlocks[resTxs[i].Height])
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

// ValidateTxResult performs transaction verification.
func ValidateTxResult(cliCtx context.CLIContext, resTx *ctypes.ResultTx) error {
	if !cliCtx.TrustNode {
		check, err := cliCtx.Verify(resTx.Height)
		if err != nil {
			return err
		}

		err = resTx.Proof.Validate(check.Header.DataHash)
		if err != nil {
			return err
		}
	}

	return nil
}

func getBlocksForTxResults(cliCtx context.CLIContext, resTxs []*ctypes.ResultTx) (map[int64]*ctypes.ResultBlock, error) {
	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	resBlocks := make(map[int64]*ctypes.ResultBlock)

	for _, resTx := range resTxs {
		if _, ok := resBlocks[resTx.Height]; !ok {
			resBlock, err := node.Block(&resTx.Height)
			if err != nil {
				return nil, err
			}

			resBlocks[resTx.Height] = resBlock
		}
	}

	return resBlocks, nil
}

func formatTxResult(cdc *codec.Codec, resTx *ctypes.ResultTx, resBlock *ctypes.ResultBlock) (sdk.TxResponse, error) {
	tx, err := parseTx(cdc, resTx.Tx)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	return sdk.NewResponseResultTx(resTx, tx, resBlock.Block.Time.Format(time.RFC3339)), nil
}

func parseTx(cdc *codec.Codec, txBytes []byte) (sdk.Tx, error) {
	var tx auth.StdTx

	err := cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func queryTx(cdc *codec.Codec, cliCtx context.CLIContext, hashHexStr string) (sdk.TxResponse, error) {
	hash, err := hex.DecodeString(hashHexStr)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	node, err := cliCtx.GetNode()
	if err != nil {
		return sdk.TxResponse{}, err
	}

	resTx, err := node.Tx(hash, !cliCtx.TrustNode)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	if !cliCtx.TrustNode {
		if err = ValidateTxResult(cliCtx, resTx); err != nil {
			return sdk.TxResponse{}, err
		}
	}

	resBlocks, err := getBlocksForTxResults(cliCtx, []*ctypes.ResultTx{resTx})
	if err != nil {
		return sdk.TxResponse{}, err
	}

	out, err := formatTxResult(cdc, resTx, resBlocks[resTx.Height])
	if err != nil {
		return out, err
	}

	return out, nil
}

func queryTxInMempool(cdc *codec.Codec, cliCtx context.CLIContext, hashHexStr string) (sdk.TxResponse, error) {
	hash, err := hex.DecodeString(hashHexStr)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	httpcli := cliCtx.GetNewHttpClient()
	if httpcli == nil {
		return sdk.TxResponse{}, fmt.Errorf("node is empty")
	}

	// default limit 30, max limit 100
	// tendermint/mempool ReapMaxTxs  `max` accept negative number as no limit
	// but a mutex in ReapMaxTxs , if no limit , maybe lock mempool too long time.
	rst, err := httpcli.UnconfirmedTxs(100)
	if err != nil || rst == nil {
		return sdk.TxResponse{}, err
	}
	for _, txBytes := range rst.Txs {
		if 0 == bytes.Compare(hash, txBytes.Hash()) {
			sdktx, err := parseTx(cdc, txBytes)
			if err != nil {
				return sdk.TxResponse{}, err
			}
			return sdk.TxResponse{TxHash: hashHexStr, Height: 0, Tx: sdktx}, nil
		}
	}

	return sdk.TxResponse{}, fmt.Errorf("not found tx %s in mempool", hashHexStr)
}

func queryTxsInMempool(cdc *codec.Codec, cliCtx context.CLIContext) (rettxs MempoolTxsResponse, err error) {
	httpcli := cliCtx.GetNewHttpClient()
	if httpcli == nil {
		return MempoolTxsResponse{}, fmt.Errorf("node is empty")
	}

	// default limit 30, max limit 100
	// tendermint/mempool ReapMaxTxs  `max` accept negative number as no limit
	// but a mutex in ReapMaxTxs , if no limit , maybe lock mempool too long time.
	rst, err := httpcli.UnconfirmedTxs(100)
	if err != nil || rst == nil {
		return MempoolTxsResponse{}, err
	}
	if err = rettxs.ParseUnconfirmedTxs(cdc, rst); err != nil {
		return MempoolTxsResponse{}, err
	}

	return
}

func queryTxsNumInMempool(cdc *codec.Codec, cliCtx context.CLIContext) (rsp MempoolTxNumResponse, err error) {
	httpcli := cliCtx.GetNewHttpClient()
	if httpcli == nil {
		err = fmt.Errorf("node is empty")
		return
	}

	// default limit 30, max limit 100
	// tendermint/mempool ReapMaxTxs  `max` accept negative number as no limit
	// but a mutex in ReapMaxTxs , if no limit , maybe lock mempool too long time.
	rst, err := httpcli.NumUnconfirmedTxs()
	if err == nil {
		if rst != nil {
			rsp = NewMempoolTxNumResponse(rst.Total, rst.TotalBytes)
		} else {
			rsp = NewMempoolTxNumResponse(0, 0)
		}
		return
	}
	return NewMempoolTxNumResponse(0, 0), err
}
