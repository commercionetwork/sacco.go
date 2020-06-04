package sacco

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/cosmos/go-bip39"

	"github.com/commercionetwork/sacco.go/softwarewallet"
)

// Wallet is a facility used to manipulate private and public keys associated
// to a BIP-32 mnemonic.
type Wallet struct {
	// TODO: just leave publicKey (extendedkey) and Address here,
	// leave all the other fields as a CryptoProvider-dependent implementation.
	keyPair         *hdkeychain.ExtendedKey
	publicKey       *hdkeychain.ExtendedKey
	PublicKey       string `json:"public_key,omitempty"`
	PublicKeyBech32 string `json:"public_key_bech_32,omitempty"`
	PrivateKey      string `json:"private_key,omitempty"`
	Path            string `json:"path,omitempty"`
	HRP             string `json:"hrp,omitempty"`
	Address         string `json:"address,omitempty"`
}

// TODO: save a "cryptoprovider" instance here, defaulting to the software one.
// A crypto provider must provide implementations for
// - FromMnemonic
// - Sign
// - Bech32PublicKey
// - Derive
// Derive and FromMnemonic can be merged - each CryptoProvider can define a set of options.

var DefaultProvider = softwarewallet.SoftwareWallet{}

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

	pkec, err := pk.ECPubKey()
	if err != nil {
		return nil, ErrCouldNotBech32(err)
	}

	pkb32, err := Bech32AminoPubKey(pkec.SerializeCompressed(), w.HRP)
	if err != nil {
		return nil, ErrCouldNotBech32(err)
	}

	w.PublicKeyBech32 = pkb32

	return &w, nil
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

// Sign signs tx with given chainID, accountNumber and sequenceNumber, with w's private key.
// The resulting computation must be enclosed in a Transaction struct to be sent over the wire
// to a Cosmos LCD.
func (w Wallet) Sign(tx TransactionPayload, chainID, accountNumber, sequenceNumber string) (SignedTransactionPayload, error) {
	signBytes := SignBytes(tx, chainID, accountNumber, sequenceNumber)

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
