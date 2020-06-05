package sacco

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/bech32"
)

// SignBytes transforms a TransactionPayload with its chainID, accountNumber e sequenceNumber
// in a sorted-by-fieldname JSON representation, ready to be signed.
func SignBytes(tx TransactionPayload, chainID, accountNumber, sequenceNumber string) []byte {
	txs := TransactionSignature{
		AccountNumber: accountNumber,
		ChainID:       chainID,
		Fee:           tx.Fee,
		Sequence:      sequenceNumber,
		Msgs:          tx.Message,
		Memo:          tx.Memo,
	}
	txbytes, _ := json.Marshal(txs)
	return sdk.MustSortJSON(txbytes)
}

func Bech32AminoPubKey(pkBytes []byte, hrp string) (string, error) {
	var cdc = amino.NewCodec()

	if len(pkBytes) != 33 {
		return "", fmt.Errorf("argument length is %d bytes, must be 33 bytes", len(pkBytes))
	}

	if strings.TrimSpace(hrp) == "" {
		return "", errors.New("hrp is empty")
	}

	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterConcrete(secp256k1.PubKeySecp256k1{},
		"tendermint/PubKeySecp256k1", nil)

	pubkTm := secp256k1.PubKeySecp256k1{}

	copy(pubkTm[:], pkBytes)

	return bech32.ConvertAndEncode(hrp+"pub", cdc.MustMarshalBinaryBare(pubkTm))
}
