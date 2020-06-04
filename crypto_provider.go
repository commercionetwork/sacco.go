package sacco

type SignData struct {
	Tx             TransactionPayload
	ChainID        string
	AccountNumber  string
	SequenceNumber string
}

type ProviderSignature struct {
	R []byte
	S []byte
}

type CryptoProvider interface {
	SignBlob([]byte) (ProviderSignature, error)
	PublicKey() ([]byte, error)
	Address() ([]byte, error)
	Bech32PublicKey() (string, error)
}
