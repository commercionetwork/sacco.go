package softwarewallet_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/sacco.go"
	"github.com/commercionetwork/sacco.go/softwarewallet"
)

func workingSW(t *testing.T) *softwarewallet.SoftwareWallet {
	opts := softwarewallet.DeriveOptions{
		Path:     sacco.CosmosDerivationPath,
		HRP:      "did:com:",
		Mnemonic: "final random flame cinnamon grunt hazard easily mutual resist pond solution define knife female tongue crime atom jaguar alert library best forum lesson rigid",
	}

	sw, err := softwarewallet.Derive(opts)
	require.NoError(t, err)

	return sw
}

func TestDerive(t *testing.T) {

	tests := []struct {
		name    string
		opts    softwarewallet.DeriveOptions
		wantErr bool
	}{
		{
			"opts is empty",
			softwarewallet.DeriveOptions{},
			true,
		},
		{
			"opts has bad mnemonic",
			softwarewallet.DeriveOptions{
				Mnemonic: "bad",
			},
			true,
		},
		{
			"opts has bad path",
			softwarewallet.DeriveOptions{
				Path:     "",
				HRP:      "did:com:",
				Mnemonic: "final random flame cinnamon grunt hazard easily mutual resist pond solution define knife female tongue crime atom jaguar alert library best forum lesson rigid",
			},
			true,
		},
		{
			"opts has empty hrp",
			softwarewallet.DeriveOptions{
				Path:     sacco.CosmosDerivationPath,
				HRP:      "",
				Mnemonic: "final random flame cinnamon grunt hazard easily mutual resist pond solution define knife female tongue crime atom jaguar alert library best forum lesson rigid",
			},
			true,
		},
		{
			"all ok",
			softwarewallet.DeriveOptions{
				Path:     sacco.CosmosDerivationPath,
				HRP:      "did:com:",
				Mnemonic: "final random flame cinnamon grunt hazard easily mutual resist pond solution define knife female tongue crime atom jaguar alert library best forum lesson rigid",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := softwarewallet.Derive(tt.opts)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, res)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, res)
		})
	}
}

func TestSoftwareWallet_PublicKey(t *testing.T) {
	sw := workingSW(t)

	tests := []struct {
		name    string
		wantPk  []byte
		wantErr bool
	}{
		{
			"public key is correctly derived",
			[]byte{3, 107, 56, 31, 200, 141, 91, 113, 119, 176, 150, 120, 71, 113, 74, 34, 136, 141, 224, 143, 159, 173, 253, 82, 69, 208, 40, 129, 37, 102, 208, 178, 149},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := sw.PublicKey()

			if tt.wantErr {
				require.Error(t, err)
				require.Empty(t, res)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantPk, res)
		})
	}
}

func TestSoftwareWallet_Address(t *testing.T) {
	sw := workingSW(t)

	tests := []struct {
		name    string
		want    []byte
		wantErr bool
	}{
		{
			"address is derived correctly",
			[]byte{0x64, 0x69, 0x64, 0x3a, 0x63, 0x6f, 0x6d, 0x3a, 0x31, 0x68, 0x75, 0x79, 0x64, 0x65, 0x65, 0x76, 0x70, 0x7a, 0x33, 0x37, 0x73, 0x64, 0x39, 0x73, 0x6e, 0x6b, 0x67, 0x75, 0x6c, 0x36, 0x30, 0x37, 0x30, 0x6d, 0x73, 0x74, 0x75, 0x70, 0x75, 0x6b, 0x77, 0x63, 0x79, 0x61, 0x61, 0x77, 0x6b},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := sw.Address()

			if tt.wantErr {
				require.Error(t, err)
				require.Empty(t, res)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, res)
		})
	}
}

func TestSoftwareWallet_Bech32PublicKey(t *testing.T) {
	sw := workingSW(t)

	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			"bech32 public key is correctly derived",
			"did:com:pub1addwnpepqd4ns87g34dhzaasjeuywu22y2ygmcy0n7kl65j96q5gzftx6zef2l3c0zk",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := sw.Bech32PublicKey()

			if tt.wantErr {
				require.Error(t, err)
				require.Empty(t, res)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, res)
		})
	}
}

func TestSoftwareWallet_SignBlob(t *testing.T) {
	sw := workingSW(t)

	tests := []struct {
		name    string
		blob    []byte
		want    sacco.ProviderSignature
		wantErr bool
	}{
		{
			"signature works on arbitrary blob",
			[]byte("blob"),
			sacco.ProviderSignature{
				R: []byte{0x9f, 0x23, 0xf7, 0x81, 0xf9, 0xb5, 0x52, 0x7, 0x21, 0x1, 0x7d, 0x20, 0xd7, 0xba, 0xd, 0x38, 0xb0, 0xc1, 0xb1, 0x16, 0x64, 0xb1, 0x81, 0x45, 0x69, 0xd1, 0xbd, 0xe6, 0x1e, 0x73, 0x7e, 0xfb},
				S: []byte{0x12, 0x25, 0x88, 0x87, 0xab, 0x13, 0xc3, 0x1a, 0x97, 0x31, 0x92, 0x8d, 0x99, 0x13, 0x9a, 0x45, 0x76, 0x42, 0x3a, 0x77, 0x2f, 0x47, 0x6e, 0xb8, 0xbd, 0x50, 0xc, 0xa, 0x43, 0x3c, 0xe8, 0xe0},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := sw.SignBlob(tt.blob)

			if tt.wantErr {
				require.Error(t, err)
				require.Empty(t, res)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, res)
		})
	}
}
