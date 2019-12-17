package sacco

import (
	"reflect"
	"testing"

	"github.com/cosmos/go-bip39"
	"github.com/stretchr/testify/assert"
)

func Test_deriveFromMnemonic(t *testing.T) {
	type args struct {
		hrp      string
		mnemonic string
		path     string
	}
	tests := []struct {
		name        string
		args        args
		wantAddress string
		assertion   assert.ErrorAssertionFunc
	}{
		{
			"a well-formed mnemonic",
			args{
				hrp:      "cosmos",
				mnemonic: "final random flame cinnamon grunt hazard easily mutual resist pond solution define knife female tongue crime atom jaguar alert library best forum lesson rigid",
				path:     CosmosDerivationPath,
			},
			"cosmos1huydeevpz37sd9snkgul6070mstupukw00xkw9",
			assert.NoError,
		},
		{
			"non-valid mnemonic",
			args{
				hrp:      "cosmos",
				mnemonic: "no",
				path:     CosmosDerivationPath,
			},
			"",
			assert.Error,
		},
		{
			"empty mnemonic",
			args{
				hrp:      "cosmos",
				mnemonic: "",
				path:     CosmosDerivationPath,
			},
			"",
			assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotAddress, err := deriveFromMnemonic(tt.args.hrp, tt.args.mnemonic, tt.args.path)
			tt.assertion(t, err)
			assert.Equal(t, tt.wantAddress, gotAddress)
		})
	}
}

func Test_derivePath(t *testing.T) {
	ane := reflect.ValueOf(assert.NoError)

	type args struct {
		seed []byte
		path string
	}
	tests := []struct {
		name        string
		args        args
		wantPrivKey string
		assertion   assert.ErrorAssertionFunc
	}{
		{
			"empty seed",
			args{
				seed: []byte{},
				path: CosmosDerivationPath,
			},
			"",
			assert.Error,
		},
		{
			"invalid derivation path",
			args{
				seed: bip39.NewSeed("seed", ""),
				path: CosmosDerivationPath + "wrong",
			},
			"",
			assert.Error,
		},
		{
			"valid seed, derivation ok",
			args{
				seed: bip39.NewSeed("seed", ""),
				path: CosmosDerivationPath,
			},
			"xprvA3bZykkgQzVPPkaBMmXtVHAfj5fKojZSSyS38EuUQP1Ks78Q1gSbuaViTAWH12q7D5YRfG4qpd1DSRN3RuDtxZgzpLYxi36prgdiUVLECzW",
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := derivePath(tt.args.seed, tt.args.path)
			tt.assertion(t, err)

			// if we're expecting assert.NoError, that means
			// a good ExtendedKey has been derived, hence
			// we can test for its Stringer output equality to our
			// test case
			tta := reflect.ValueOf(tt.assertion)
			if tta.Pointer() == ane.Pointer() {
				assert.Equal(t, tt.wantPrivKey, got.String())
			}
		})
	}
}

func Test_stringToComponents(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		want      []derivationComponent
		assertion assert.ErrorAssertionFunc
	}{
		{
			"a well-formed derivation path with spaces between characters",
			"m / 44' / 0' / 0' / 0 / 0 ",
			[]derivationComponent{
				{
					Path:     44,
					Hardened: true,
				},
				{
					Path:     0,
					Hardened: true,
				},
				{
					Path:     0,
					Hardened: true,
				},
				{
					Path:     0,
					Hardened: false,
				},
				{
					Path:     0,
					Hardened: false,
				},
			},
			assert.NoError,
		},
		{
			"a well-formed derivation path without spaces between characters",
			"m/44'/0'/0'/0/0",
			[]derivationComponent{
				{
					Path:     44,
					Hardened: true,
				},
				{
					Path:     0,
					Hardened: true,
				},
				{
					Path:     0,
					Hardened: true,
				},
				{
					Path:     0,
					Hardened: false,
				},
				{
					Path:     0,
					Hardened: false,
				},
			},
			assert.NoError,
		},
		{
			"derivation path which doesn't begin with \"m\"",
			"/44'/0'/0'/0/0",
			[]derivationComponent{},
			assert.Error,
		},
		{
			"random text, invalid derivation path",
			"i am invalid!",
			[]derivationComponent{},
			assert.Error,
		},
		{
			"a partially-valid derivation path",
			"m/44'/0'/0'/k/0",
			[]derivationComponent{},
			assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := stringToComponents(tt.path)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_hardened(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		wantIsHardened bool
		wantDestStr    string
	}{
		{
			"hardened component path",
			"44'",
			true,
			"44",
		},
		{
			"non-hardened component path",
			"44",
			false,
			"44",
		},
		{
			"empty component path",
			"",
			false,
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIsHardened, gotDestStr := hardened(tt.path)
			assert.Equal(t, tt.wantIsHardened, gotIsHardened)
			assert.Equal(t, tt.wantDestStr, gotDestStr)
		})
	}
}
