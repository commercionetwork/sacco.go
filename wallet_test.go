package sacco_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/sacco.go"
	"github.com/commercionetwork/sacco.go/softwarewallet"
)

func TestWallet_Sign(t *testing.T) {
	type fields struct {
		hrp      string
		mnemonic string
		path     string
	}
	type args struct {
		tx             sacco.TransactionPayload
		chainID        string
		accountNumber  string
		sequenceNumber string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      sacco.SignedTransactionPayload
		assertion require.ErrorAssertionFunc
	}{
		{
			"signing a transaction with a valid mnemonic",
			fields{
				hrp:      "com:did:",
				mnemonic: "innocent pony teach letter mask bulk stuff pool more work cute prepare forest simple sunset sphere aisle luggage task drama fire clutch trial search",
				path:     sacco.CosmosDerivationPath,
			},
			args{
				tx: sacco.TransactionPayload{
					Message: []json.RawMessage{
						json.RawMessage(`{"type":"cosmos-sdk/MsgSend","value":{"from_address":"did:com:1sfjela2snk9rmmcfh773gm50476w0ur5pmwuak","to_address":"did:com:1kulfxlg33x9lmxa00gmmaq6j3nshtpnrr24tm9","amount":[{"denom":"ucommercio","amount":"10"}]}}`),
					},
					Fee: sacco.Fee{
						Amount: []sacco.Coin{},
						Gas:    "200000",
					},
				},
				chainID:        "test-chain-jVvnJ6",
				accountNumber:  "11",
				sequenceNumber: "0",
			},
			sacco.SignedTransactionPayload{
				Message: []json.RawMessage{
					json.RawMessage(`{"type":"cosmos-sdk/MsgSend","value":{"from_address":"did:com:1sfjela2snk9rmmcfh773gm50476w0ur5pmwuak","to_address":"did:com:1kulfxlg33x9lmxa00gmmaq6j3nshtpnrr24tm9","amount":[{"denom":"ucommercio","amount":"10"}]}}`),
				},
				Fee: sacco.Fee{
					Amount: []sacco.Coin{},
					Gas:    "200000",
				},
				// this signature has been generated using a local chain with cncli
				Signatures: []sacco.Signature{
					{
						SigPubKey: sacco.SigPubKey{
							Type:  "tendermint/PubKeySecp256k1",
							Value: "A6WEhS1jR2qwULCuneR7miIMnzg/lFubu3IaPb0K4TVQ",
						},
						Signature: "z/oFsC5M/7ES9MEef3L6Zf6QKlFTUelpj25w3mrPk292WRYQLIKPuYsywLouIaa4cdHHfqfjSh9J8m+ZEwVK3Q==",
					},
				},
			},
			require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sw, err := softwarewallet.Derive(softwarewallet.DeriveOptions{
				Path:     tt.fields.path,
				HRP:      tt.fields.hrp,
				Mnemonic: tt.fields.mnemonic,
			})

			tt.assertion(t, err)

			w, err := sacco.NewWallet(sw)
			require.NoError(t, err)

			got, err := w.Sign(sacco.SignData{
				tt.args.tx,
				tt.args.chainID,
				tt.args.accountNumber,
				tt.args.sequenceNumber,
			})

			tt.assertion(t, err)

			assert.Equal(t, tt.want, got)
		})
	}
}

// errorCProv is a CryptoProvider which always fails, conditionally.
type errorCProv struct {
	addressFails  bool
	pubkeyFails   bool
	signBlobFails bool
	bech32Fails   bool
}

func (e errorCProv) SignBlob(_ []byte) (sacco.ProviderSignature, error) {
	if !e.signBlobFails {
		return sacco.ProviderSignature{
			R: []byte{42},
			S: []byte{42},
		}, nil
	}
	return sacco.ProviderSignature{}, errors.New("error")
}

func (e errorCProv) PublicKey() ([]byte, error) {
	if !e.pubkeyFails {
		return []byte{42}, nil
	}
	return nil, errors.New("error")
}

func (e errorCProv) Address() ([]byte, error) {
	if !e.addressFails {
		return []byte{42}, nil
	}
	return nil, errors.New("error")
}

func (e errorCProv) Bech32PublicKey() (string, error) {
	if !e.bech32Fails {
		return "42", nil
	}
	return "", errors.New("error")
}

func TestNewWallet(t *testing.T) {
	tests := []struct {
		name    string
		p       sacco.CryptoProvider
		want    *sacco.Wallet
		wantErr bool
	}{
		{
			"address fails",
			errorCProv{addressFails: true},
			nil,
			true,
		},
		{
			"bech32 fails",
			errorCProv{bech32Fails: true},
			nil,
			true,
		},
		{
			"all ok",
			errorCProv{},
			&sacco.Wallet{
				Address:   string([]byte{42}),
				PublicKey: "42",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := sacco.NewWallet(tt.p)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, res)
				return
			}

			require.NotNil(t, res)
			require.Equal(t, tt.want.Address, res.Address)
			require.Equal(t, tt.want.PublicKey, res.PublicKey)
			require.NoError(t, err)
		})
	}
}
