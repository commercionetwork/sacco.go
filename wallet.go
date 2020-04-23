package sacco

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/awnumar/memguard"
	"github.com/btcsuite/btcutil/hdkeychain"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/go-bip39"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/bech32"
)

// Wallet is a facility used to manipulate private and public keys associated
// to a BIP-32 mnemonic.
type Wallet struct {
	keyPair         *hdkeychain.ExtendedKey
	publicKey       *hdkeychain.ExtendedKey
	PublicKey       string `json:"public_key,omitempty"`
	PublicKeyBech32 string `json:"public_key_bech_32,omitempty"`
	PrivateKey      string `json:"private_key,omitempty"`
	Path            string `json:"path,omitempty"`
	HRP             string `json:"hrp,omitempty"`
	Address         string `json:"address,omitempty"`
}

// FromMnemonic returns a new Wallet instance given a human-readable part,
// mnemonic and path.
func FromMnemonic(hrp, mnemonic, path string) (*Wallet, error) {
	var w Wallet
	k, a, err := deriveFromMnemonic(hrp, mnemonic, path)
	if err != nil {
		return nil, err
	}

	w.keyPair = k
	w.Path = path
	w.Address = a
	w.HRP = hrp

	pk, err := w.keyPair.Neuter()
	if err != nil {
		return nil, ErrCouldNotNeuter(err)
	}

	w.publicKey = pk
	w.PublicKey = w.publicKey.String()

	pkb32, err := w.bech32AminoPubKey()
	if err != nil {
		return nil, ErrCouldNotBech32(err)
	}

	w.PublicKeyBech32 = pkb32

	return &w, nil
}

func (w Wallet) bech32AminoPubKey() (string, error) {
	pkec, _ := w.publicKey.ECPubKey()

	var cdc = amino.NewCodec()

	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterConcrete(secp256k1.PubKeySecp256k1{},
		"tendermint/PubKeySecp256k1", nil)

	pubkTm := secp256k1.PubKeySecp256k1{}

	for i, l := range pkec.SerializeCompressed() {
		pubkTm[i] = l
	}

	return bech32.ConvertAndEncode(w.HRP+"pub", cdc.MustMarshalBinaryBare(pubkTm))
}

// Export creates a JSON representation of w.
// Export does not include the private key in the JSON representation.
func (w Wallet) Export() (string, error) {
	w.PrivateKey = ""
	data, err := json.Marshal(w)

	return string(data), err
}

// ExportWithPrivateKey creates a JSON representation of w.
// ExportWithPrivateKey includes the private key in the JSON representation.
func (w Wallet) ExportWithPrivateKey() (string, error) {
	w.PrivateKey = w.keyPair.String()

	s := memguard.NewStream()

	enc := json.NewEncoder(s)
	err := enc.Encode(w)
	if err != nil {
		return "", err
	}

	enclave := s.Front().Value.(*memguard.Enclave)
	data, err := enclave.Open()
	if err != nil {
		return "", err
	}

	data.Melt()

	defer data.Destroy()

	return strings.TrimSpace(string(data.Bytes())), err
}

// GenerateMnemonic generates a new random mnemonic sequence.
func GenerateMnemonic() (string, error) {
	sb, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
	if err != nil {
		return "", err
	}
	mnemonic, err := bip39.NewMnemonic(sb)

	return mnemonic, err
}

// signBytes transforms a TransactionPayload with its chainID, accountNumber e sequenceNumber
// in a sorted-by-fieldname JSON representation, ready to be signed.
func signBytes(tx TransactionPayload, chainID, accountNumber, sequenceNumber string) []byte {
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

// Sign signs tx with given chainID, accountNumber and sequenceNumber, with w's private key.
// The resulting computation must be enclosed in a Transaction struct to be sent over the wire
// to a Cosmos LCD.
func (w Wallet) Sign(tx TransactionPayload, chainID, accountNumber, sequenceNumber string) (SignedTransactionPayload, error) {
	signBytes := signBytes(tx, chainID, accountNumber, sequenceNumber)

	pk, err := w.keyPair.ECPrivKey()
	if err != nil {
		return SignedTransactionPayload{}, err
	}

	hashSb := sha256.Sum256(signBytes)
	signatureRaw, err := pk.Sign(hashSb[:])
	if err != nil {
		return SignedTransactionPayload{}, err
	}

	signatureRaw.Serialize()
	rBytes := signatureRaw.R.Bytes()
	sBytes := signatureRaw.S.Bytes()

	pubKey, err := w.publicKey.ECPubKey()
	if err != nil {
		return SignedTransactionPayload{}, err
	}

	r := []byte{}
	r = append(r, rBytes...)
	r = append(r, sBytes...)
	signature := base64.StdEncoding.EncodeToString(r)
	compressedPubKey := base64.StdEncoding.EncodeToString(pubKey.SerializeCompressed())

	tx.Signatures = []Signature{
		{
			Signature: signature,
			SigPubKey: SigPubKey{
				Type:  "tendermint/PubKeySecp256k1",
				Value: compressedPubKey,
			},
		},
	}

	return SignedTransactionPayload(tx), nil
}
