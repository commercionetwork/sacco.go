package sacco

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/bech32"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/cosmos/go-bip39"

	// TODO: search for a modern, better implementation of ripemd160
	// nolint:staticcheck
	"golang.org/x/crypto/ripemd160"
)

// derivationComponent holds informations about a single derivation
// path, used during the derivation process.
type derivationComponent struct {
	Path     uint32
	Hardened bool
}

// derivateFromMnemonic derivates an HD keypair and address from a mnemonic, a path and an
// human-readable part.
func derivateFromMnemonic(hrp, mnemonic, path string) (key *hdkeychain.ExtendedKey, address string, err error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, "", fmt.Errorf("invalid mnemonic")
	}

	seed := bip39.NewSeed(mnemonic, "")
	key, err = derivatePath(seed, path)
	if err != nil {
		return nil, "", err
	}

	epk, err := key.ECPubKey()
	if err != nil {
		return nil, "", err
	}

	addr, err := addressFromPublicKey(epk, hrp)
	if err != nil {
		return nil, "", err
	}

	return key, addr, nil
}

// addressFromPublicKey returns a bech32-encoded address given a public key
// and an human-readable part.
func addressFromPublicKey(pk *btcec.PublicKey, hrp string) (string, error) {
	sha := sha256.Sum256(pk.SerializeCompressed())
	s := sha[:]
	r := ripemd160.New()
	_, err := r.Write(s)
	if err != nil {
		return "", err
	}
	pub := r.Sum(nil)

	converted, err := bech32.ConvertBits(pub, 8, 5, true)
	if err != nil {
		return "", err
	}

	addr, err := bech32.Encode(hrp, converted)
	if err != nil {
		return "", err
	}

	return addr, nil
}

// derivatePath derivates an HD keypair from a seed, and a derivation path.
func derivatePath(seed []byte, path string) (*hdkeychain.ExtendedKey, error) {
	params := chaincfg.MainNetParams
	master, err := hdkeychain.NewMaster(seed, &params)
	if err != nil {
		return nil, ErrKeyGeneration(err)
	}

	components, err := stringToComponents(path)
	if err != nil {
		return nil, err
	}

	var child *hdkeychain.ExtendedKey

	for _, component := range components {
		// If k is nil, this means we're deriving a child key for the first time,
		// so use the master key.
		// This condition is true only on the first components element (the purpose).
		k := child
		if k == nil {
			k = master
		}

		if component.Hardened {
			child, err = k.Child(component.Path + hdkeychain.HardenedKeyStart)
		} else {
			child, err = k.Child(component.Path)
		}

		if err != nil {
			return nil, ErrKeyGeneration(err)
		}
	}

	return child, nil
}

// stringToComponents transforms a derivation path string into a slice
// of DerivationComponents.
func stringToComponents(path string) ([]derivationComponent, error) {
	path = strings.Replace(path, " ", "", -1)

	components := strings.Split(path, "/")
	if len(components) <= 1 {
		return []derivationComponent{}, ErrDerivationPathShort
	}

	if components[0] != "m" {
		return []derivationComponent{}, ErrDerivationPathFirstCharNotM
	}

	// ignore the "m", we don't need that
	components = components[1:]

	// build a DerivationComponent for each element in the path
	dcs := make([]derivationComponent, len(components))

	for index, rawComponent := range components {
		isHardened, rawPathNum := hardened(rawComponent)

		pathNum, convErr := strconv.ParseUint(rawPathNum, 10, 32)

		if convErr != nil || rawPathNum == "" {
			return []derivationComponent{}, ErrComponentNaN(rawPathNum, convErr)
		}

		dcs[index] = derivationComponent{
			Path:     uint32(pathNum),
			Hardened: isHardened,
		}
	}

	return dcs, nil
}

// hardened returns true whether s is an hardened derivation path
// component, false otherwise.
// When hardened returns true, destStr will contain s without the
// hardened indicator (the "'"), otherwise destStr will be equal to
// s.
func hardened(s string) (isHardened bool, destStr string) {
	if len(s) <= 0 {
		return false, ""
	}

	isHardened = s[len(s)-1] == '\''
	destStr = s

	if isHardened {
		destStr = string(s[:len(s)-1])
	}

	return
}
