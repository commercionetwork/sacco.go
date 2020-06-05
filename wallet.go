package sacco

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
)

// Wallet is a facility used to manipulate private and public keys, send transaction to LCD nodes.
type Wallet struct {
	cp        CryptoProvider
	Address   string
	PublicKey string
}

func NewWallet(p CryptoProvider) (*Wallet, error) {
	address, err := p.Address()
	if err != nil {
		return nil, err
	}

	pubKey, err := p.Bech32PublicKey()
	if err != nil {
		return nil, err
	}
	return &Wallet{
		cp:        p,
		Address:   string(address),
		PublicKey: pubKey,
	}, nil
}

// SignAndBroadcast signs tx and broadcast it to the LCD specified by lcdEndpoint.
func (w Wallet) SignAndBroadcast(tx TransactionPayload, lcdEndpoint string, txMode TxMode) (string, error) {
	// get network (chain) name
	nodeInfo, err := getNodeInfo(lcdEndpoint)
	if err != nil {
		return "", fmt.Errorf("could not get LCD node informations: %w", err)
	}

	addressBytes, err := w.cp.Address()
	if err != nil {
		return "", err
	}

	address := string(addressBytes)

	// get account sequence and account number
	accountData, err := getAccountData(lcdEndpoint, address)
	if err != nil {
		return "", fmt.Errorf("could not get Account informations for address %s: %w", address, err)
	}

	// sign transaction
	signedTx, err := w.Sign(
		SignData{
			tx,
			nodeInfo.Info.Network,
			strconv.FormatInt(accountData.Result.Value.AccountNumber, 10),
			strconv.FormatInt(accountData.Result.Value.Sequence, 10),
		},
	)
	if err != nil {
		return "", fmt.Errorf("could not sign transaction: %w", err)
	}

	// broadcast transaction to the LCD
	txHash, err := broadcastTx(signedTx, lcdEndpoint, txMode)
	if err != nil {
		return "", fmt.Errorf("could not broadcast transaction to the Cosmos network: %w", err)
	}

	// return transaction hash!
	return txHash, nil
}

func (w Wallet) Sign(sd SignData) (SignedTransactionPayload, error) {
	signBytes := SignBytes(sd.Tx, sd.ChainID, sd.AccountNumber, sd.SequenceNumber)

	hashSb := sha256.Sum256(signBytes)
	signatureRaw, err := w.cp.SignBlob(hashSb[:])
	if err != nil {
		return SignedTransactionPayload{}, err
	}

	r := []byte{}
	r = append(r, signatureRaw.R...)
	r = append(r, signatureRaw.S...)
	signature := base64.StdEncoding.EncodeToString(r)

	pubKey, err := w.cp.PublicKey()
	if err != nil {
		return SignedTransactionPayload{}, err
	}
	compressedPubKey := base64.StdEncoding.EncodeToString(pubKey)

	sd.Tx.Signatures = []Signature{
		{
			Signature: signature,
			SigPubKey: SigPubKey{
				Type:  "tendermint/PubKeySecp256k1",
				Value: compressedPubKey,
			},
		},
	}

	return SignedTransactionPayload(sd.Tx), nil
}
