package sacco

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcutil/bech32"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/crypto/ripemd160"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmBech32 "github.com/tendermint/tendermint/libs/bech32"
)

const ecPubKeyLen = 33

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

// Bech32AminoPubKey returns the Bech32 representation of pkBytes using hrp as its human-readable part.
// pkBytes must contain the Compressed public key (its length must be 33 bytes).
func Bech32AminoPubKey(pkBytes []byte, hrp string) (string, error) {
	var cdc = amino.NewCodec()

	if len(pkBytes) != ecPubKeyLen {
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

	return tmBech32.ConvertAndEncode(hrp+"pub", cdc.MustMarshalBinaryBare(pubkTm))
}

// Bech32Address returns the Bech32 address for a pkBytes public key, using hrp as its human-readable part.
// pkBytes must contain the Compressed public key (its length must be 33 bytes).
func Bech32Address(pkBytes []byte, hrp string) ([]byte, error) {
	if len(pkBytes) != ecPubKeyLen {
		return nil, fmt.Errorf("argument length is %d bytes, must be 33 bytes", len(pkBytes))
	}

	if strings.TrimSpace(hrp) == "" {
		return nil, errors.New("hrp is empty")
	}

	sha := sha256.Sum256(pkBytes)
	s := sha[:]
	r := ripemd160.New()
	_, err := r.Write(s)
	if err != nil {
		return nil, err
	}
	pub := r.Sum(nil)

	converted, err := bech32.ConvertBits(pub, 8, 5, true)
	if err != nil {
		return nil, err
	}

	addr, err := bech32.Encode(hrp, converted)
	if err != nil {
		return nil, err
	}

	return []byte(addr), nil
}
