package softwarewallet

import (
	"errors"

	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/cosmos/go-bip39"

	"github.com/commercionetwork/sacco.go"
)

// SoftwareWallet is a CryptoProvider-compliant facility that implements a software-based Cosmos wallet.
type SoftwareWallet struct {
	keyPair   *hdkeychain.ExtendedKey
	publicKey *hdkeychain.ExtendedKey
	path      string
	hrp       string
	address   string
}

// DeriveOptions holds data that will be used by the Derive function to later derive a keypair and a wallet.
type DeriveOptions struct {
	Path     string
	HRP      string
	Mnemonic string
}

// Derive derives a SoftwareWallet instance with with given DeriveOptions.
// If DeriveOptions.Mnemonic is empty, a new mnemonic will be generated and used.
func Derive(opts DeriveOptions) (*SoftwareWallet, error) {
	if opts.Mnemonic == "" {
		var err error
		opts.Mnemonic, err = generateMnemonic()
		if err != nil {
			return nil, sacco.ErrCouldNotDerive(err)
		}
	}

	if opts.HRP == "" {
		return nil, errors.New("hrp cannot be empty")
	}

	var w SoftwareWallet
	k, a, err := deriveFromMnemonic(opts.HRP, opts.Mnemonic, opts.Path)
	if err != nil {
		return nil, err
	}

	w.keyPair = k
	w.path = opts.Path
	w.address = a
	w.hrp = opts.HRP

	pk, err := w.keyPair.Neuter()
	if err != nil {
		return nil, sacco.ErrCouldNotNeuter(err)
	}

	w.publicKey = pk

	return &w, nil
}

// generateMnemonic generates a new random mnemonic sequence.
func generateMnemonic() (string, error) {
	sb, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
	if err != nil {
		return "", err
	}
	mnemonic, err := bip39.NewMnemonic(sb)

	return mnemonic, err
}

// PublicKey implements the CryptoProvider interface.
func (sw SoftwareWallet) PublicKey() ([]byte, error) {
	pkec, err := sw.publicKey.ECPubKey()
	if err != nil {
		return nil, sacco.ErrCouldNotBech32(err)
	}

	return pkec.SerializeCompressed(), nil
}

// PublicKey implements the CryptoProvider interface.
func (sw SoftwareWallet) Address() ([]byte, error) {
	return []byte(sw.address), nil
}

// SignBlob implements CryptoProvider interface.
func (sw SoftwareWallet) SignBlob(b []byte) (sacco.ProviderSignature, error) {
	pk, err := sw.keyPair.ECPrivKey()
	if err != nil {
		return sacco.ProviderSignature{}, err
	}
	signatureRaw, err := pk.Sign(b)
	if err != nil {
		return sacco.ProviderSignature{}, err
	}

	return sacco.ProviderSignature{
		R: signatureRaw.R.Bytes(),
		S: signatureRaw.S.Bytes(),
	}, nil
}

// Bech32PublicKey implements CryptoProvider interface.
func (sw SoftwareWallet) Bech32PublicKey() (string, error) {
	pk, err := sw.PublicKey()
	if err != nil {
		return "", err
	}

	return sacco.Bech32AminoPubKey(pk, sw.hrp)
}
