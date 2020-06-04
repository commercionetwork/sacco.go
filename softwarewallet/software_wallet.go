package softwarewallet

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/btcsuite/btcutil/hdkeychain"

	"github.com/commercionetwork/sacco.go"
)

type SoftwareWallet struct {
	keyPair         *hdkeychain.ExtendedKey
	publicKey       *hdkeychain.ExtendedKey
	PublicKey       string `json:"public_key,omitempty"`
	PublicKeyBech32 string `json:"public_key_bech_32,omitempty"`
	PrivateKey      string `json:"private_key,omitempty"`
	Path            string `json:"path,omitempty"`
	HRP             string `json:"hrp,omitempty"`
	Address         string `json:"address,omitempty"`
}

func (sw SoftwareWallet) Sign(sd sacco.SignData) (sacco.SignedTransactionPayload, error) {
	signBytes := sacco.SignBytes(sd.Tx, sd.ChainID, sd.AccountNumber, sd.SequenceNumber)

	pk, err := sw.keyPair.ECPrivKey()
	if err != nil {
		return sacco.SignedTransactionPayload{}, err
	}

	hashSb := sha256.Sum256(signBytes)
	signatureRaw, err := pk.Sign(hashSb[:])
	if err != nil {
		return sacco.SignedTransactionPayload{}, err
	}

	signatureRaw.Serialize()
	rBytes := signatureRaw.R.Bytes()
	sBytes := signatureRaw.S.Bytes()

	pubKey, err := sw.publicKey.ECPubKey()
	if err != nil {
		return sacco.SignedTransactionPayload{}, err
	}

	r := []byte{}
	r = append(r, rBytes...)
	r = append(r, sBytes...)
	signature := base64.StdEncoding.EncodeToString(r)
	compressedPubKey := base64.StdEncoding.EncodeToString(pubKey.SerializeCompressed())

	sd.Tx.Signatures = []sacco.Signature{
		{
			Signature: signature,
			SigPubKey: sacco.SigPubKey{
				Type:  "tendermint/PubKeySecp256k1",
				Value: compressedPubKey,
			},
		},
	}

	return sacco.SignedTransactionPayload(sd.Tx), nil
}

func (SoftwareWallet) Derive(options interface{}) (string, error) {
	panic("implement me")
}

func (sw SoftwareWallet) Bech32PublicKey() (string, error) {
	pkec, err := sw.publicKey.ECPubKey()
	if err != nil {
		return "", sacco.ErrCouldNotBech32(err)
	}

	return sacco.Bech32AminoPubKey(pkec.SerializeCompressed(), sw.HRP)
}
