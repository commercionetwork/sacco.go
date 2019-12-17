package sacco

import "fmt"

// ErrDerivationPathShort represents an error that happens when the derivation path
// is too short.
var ErrDerivationPathShort = fmt.Errorf("derivation path string too short")

// ErrDerivationPathFirstCharNotM happens whenever the derivation path string doesn't
// begin with "m".
var ErrDerivationPathFirstCharNotM = fmt.Errorf("derivation path invalid, first character isn't 'm'")

// ErrComponentNaN happens when a component of a derivation path
// isn't a number.
var ErrComponentNaN = func(component string, err error) error {
	return fmt.Errorf("derivation component \"%s\" not a number: %w", component, err)
}

// ErrKeyGeneration represents an error thrown when there was some kind of
// error while generating a key.
var ErrKeyGeneration = func(err error) error {
	return fmt.Errorf("cannot derive key: %w", err)
}

// ErrCouldNotNeuter happens when neutering a key is not possible.
var ErrCouldNotNeuter = func(err error) error {
	return fmt.Errorf("could not derive neutered public key: %w", err)
}
