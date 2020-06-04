package sacco

type SignData struct {
	Tx             TransactionPayload
	ChainID        string
	AccountNumber  string
	SequenceNumber string
}

type CryptoProvider interface {
	Sign(sd SignData) (SignedTransactionPayload, error)
	Derive(options interface{}) (string, error)
	Bech32PublicKey() (string, error)
}
