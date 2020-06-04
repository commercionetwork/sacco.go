package sacco

import "encoding/json"

// TransactionPayload is the body of a Cosmos transaction.
type TransactionPayload struct {
	Message    []json.RawMessage `json:"msg"`
	Fee        Fee               `json:"fee"`
	Signatures []Signature       `json:"signatures"`
	Memo       string            `json:"memo"`
}

// SignedTransactionPayload is a TransactionPayload which has been signed
// by wallet.SignBlob().
type SignedTransactionPayload TransactionPayload

// TransactionSignature is a Transaction with AccountNumber, ChainID (Network) and
// sequence number, ready to be signed with a wallet's private key.
type TransactionSignature struct {
	AccountNumber string            `json:"account_number" yaml:"account_number"`
	ChainID       string            `json:"chain_id" yaml:"chain_id"`
	Fee           Fee               `json:"fee" yaml:"fee"`
	Sequence      string            `json:"sequence" yaml:"sequence"`
	Memo          string            `json:"memo" yaml:"memo"`
	Msgs          []json.RawMessage `json:"msgs" yaml:"msgs"`
}

// Fee represents the fee needed for a Cosmos request,
// with both an amount in Coins and a Gas quantity.
type Fee struct {
	Amount []Coin `json:"amount"`
	Gas    string `json:"gas"`
}

// Coin is an entity describing a token,
// with a specific denomination and amount.
type Coin struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

// Coins is a set of Coin.
type Coins []Coin

// Signature is an object holding a cryptographic signature
// and the public key associated with it.
type Signature struct {
	SigPubKey SigPubKey `json:"pub_key"`
	Signature string    `json:"signature"`
}

// SigPubKey represents the public key used to create a Signature.
type SigPubKey struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
