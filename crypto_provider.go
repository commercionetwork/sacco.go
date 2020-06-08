package sacco

// SignData represents an unit of work (an unsigned transaction) which will be processed and signed into a full-blown Tendermint
// transaction.
type SignData struct {
	// Tx is the transaction to be signed, containing messages and fees.
	Tx TransactionPayload

	// ChainID is the Cosmos chain ID where the transaction will be sent.
	ChainID string

	// AccountNumber is the account number assigned to the key which signs the transaction.
	AccountNumber string

	// SequenceNumber is the number associated to the next sign-able transaction (the one being processed).
	SequenceNumber string
}

// ProviderSignature is a Secp256k1 R || s signature obtained by calling a CryptoProvider SignBlob() method.
type ProviderSignature struct {
	R []byte
	S []byte
}

// CryptoProvider represents a logical unit which implements methods used by sacco.go to sign transactions and handle
// Secp256k1 keypairs.
//
// While any implementer must adhere to the semantics of this interface, the way they gets initializer are
// implementation-dependant.
//
// See the `softwarewallet' package for an implementation example.
//
// CryptoProvider is meant to be used along with the `NewWallet()' sacco.go function, and its semantics are designed to
// be used as a throw-away object: create a new CryptoProvider and sacco.Wallet instance every time you need one
// (multiple private keys need multiple CryptoProvider instances).
type CryptoProvider interface {
	// SignBlob signs b with a private key and returns a ProviderSignature with the R || s signature.
	SignBlob(b []byte) (ProviderSignature, error)

	// PublicKey returns a byte slice containing the public key associated with the CryptoProvider.
	PublicKey() ([]byte, error)

	// Address returns a byte slice containing the Bech32-encoded Cosmos address, derived by the CryptoProvider's
	// public key.
	Address() ([]byte, error)

	// Bech32PublicKey returns a string containing the CryptoProvider public key as a Bech32-encoded string.
	Bech32PublicKey() (string, error)
}
