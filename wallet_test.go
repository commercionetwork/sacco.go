package sacco

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromMnemonic(t *testing.T) {
	type args struct {
		hrp      string
		mnemonic string
		path     string
	}
	tests := []struct {
		name      string
		args      args
		wantAddr  string
		assertion assert.ErrorAssertionFunc
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
			got, err := FromMnemonic(tt.args.hrp, tt.args.mnemonic, tt.args.path)
			tt.assertion(t, err)

			// if Wallet was derived correctly, address must match what we expect
			if got != nil {
				assert.Equal(t, tt.wantAddr, got.Address)
			}
		})
	}
}

func TestWallet_Export(t *testing.T) {
	type args struct {
		hrp      string
		mnemonic string
		path     string
	}
	tests := []struct {
		name      string
		args      args
		wantJSON  string
		assertion assert.ErrorAssertionFunc
	}{
		{
			"a well-formed mnemonic",
			args{
				hrp:      "cosmos",
				mnemonic: "final random flame cinnamon grunt hazard easily mutual resist pond solution define knife female tongue crime atom jaguar alert library best forum lesson rigid",
				path:     CosmosDerivationPath,
			},
			`{"public_key":"xpub6FW9dWDyi8m8todcGW5YDVbzoUx4rgBWZ7nsQ8tDyVyyv4yyc1mo9ca3cRhDHfr2V3xhcHj5GDrBMoHCBZti5LRz1XrsVxSKWrPYbQFssKo","path":"m/44'/118'/0'/0/0","hrp":"cosmos","address":"cosmos1huydeevpz37sd9snkgul6070mstupukw00xkw9"}`,
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
			got, err := FromMnemonic(tt.args.hrp, tt.args.mnemonic, tt.args.path)
			tt.assertion(t, err)

			if got != nil {
				gotExport, err := got.Export()
				tt.assertion(t, err)

				assert.Equal(t, tt.wantJSON, gotExport)
			}
		})
	}
}

func TestWallet_ExportWithPrivateKey(t *testing.T) {
	type args struct {
		hrp      string
		mnemonic string
		path     string
	}
	tests := []struct {
		name      string
		args      args
		wantJSON  string
		assertion assert.ErrorAssertionFunc
	}{
		{
			"a well-formed mnemonic",
			args{
				hrp:      "cosmos",
				mnemonic: "final random flame cinnamon grunt hazard easily mutual resist pond solution define knife female tongue crime atom jaguar alert library best forum lesson rigid",
				path:     CosmosDerivationPath,
			},
			`{"public_key":"xpub6FW9dWDyi8m8todcGW5YDVbzoUx4rgBWZ7nsQ8tDyVyyv4yyc1mo9ca3cRhDHfr2V3xhcHj5GDrBMoHCBZti5LRz1XrsVxSKWrPYbQFssKo","private_key":"xprvA2WoDzh5smCqgKZ9AUYXrMfGFT7aTDTfBtsGbkUcRAT13Geq4UTYbpFZm9BYmxMBtn4fK8LYndQ7HaneCLGwT35iW2VDmPKRdErwJHRkLgX","path":"m/44'/118'/0'/0/0","hrp":"cosmos","address":"cosmos1huydeevpz37sd9snkgul6070mstupukw00xkw9"}`,
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
			got, err := FromMnemonic(tt.args.hrp, tt.args.mnemonic, tt.args.path)
			tt.assertion(t, err)

			if got != nil {
				gotExport, err := got.ExportWithPrivateKey()
				tt.assertion(t, err)

				assert.Equal(t, tt.wantJSON, gotExport)
			}
		})
	}
}

func TestGenerateMnemonic(t *testing.T) {
	tests := []struct {
		name      string
		assertion assert.ErrorAssertionFunc
	}{
		{
			"generate a new random mnemonic",
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// we can only test whether we have an error or not
			_, err := GenerateMnemonic()
			tt.assertion(t, err)
		})
	}
}

func TestWallet_Sign(t *testing.T) {
	type fields struct {
		hrp      string
		mnemonic string
		path     string
	}
	type args struct {
		tx             TransactionPayload
		chainID        string
		accountNumber  string
		sequenceNumber string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      SignedTransactionPayload
		assertion assert.ErrorAssertionFunc
	}{
		{
			"signing a transaction with a valid mnemonic",
			fields{
				hrp:      "com:did:",
				mnemonic: "innocent pony teach letter mask bulk stuff pool more work cute prepare forest simple sunset sphere aisle luggage task drama fire clutch trial search",
				path:     CosmosDerivationPath,
			},
			args{
				tx: TransactionPayload{
					Message: []json.RawMessage{
						json.RawMessage([]byte(`{"type":"cosmos-sdk/MsgSend","value":{"from_address":"did:com:1sfjela2snk9rmmcfh773gm50476w0ur5pmwuak","to_address":"did:com:1kulfxlg33x9lmxa00gmmaq6j3nshtpnrr24tm9","amount":[{"denom":"ucommercio","amount":"10"}]}}`)),
					},
					Fee: Fee{
						Amount: []Coin{},
						Gas:    "200000",
					},
				},
				chainID:        "test-chain-jVvnJ6",
				accountNumber:  "11",
				sequenceNumber: "0",
			},
			SignedTransactionPayload{
				Message: []json.RawMessage{
					json.RawMessage([]byte(`{"type":"cosmos-sdk/MsgSend","value":{"from_address":"did:com:1sfjela2snk9rmmcfh773gm50476w0ur5pmwuak","to_address":"did:com:1kulfxlg33x9lmxa00gmmaq6j3nshtpnrr24tm9","amount":[{"denom":"ucommercio","amount":"10"}]}}`)),
				},
				Fee: Fee{
					Amount: []Coin{},
					Gas:    "200000",
				},
				// this signature has been generated using a local chain with cncli
				Signatures: []Signature{
					{
						SigPubKey: SigPubKey{
							Type:  "tendermint/PubKeySecp256k1",
							Value: "A6WEhS1jR2qwULCuneR7miIMnzg/lFubu3IaPb0K4TVQ",
						},
						Signature: "z/oFsC5M/7ES9MEef3L6Zf6QKlFTUelpj25w3mrPk292WRYQLIKPuYsywLouIaa4cdHHfqfjSh9J8m+ZEwVK3Q==",
					},
				},
			},
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, err := FromMnemonic(tt.fields.hrp, tt.fields.mnemonic, tt.fields.path)
			tt.assertion(t, err)

			got, err := w.Sign(tt.args.tx, tt.args.chainID, tt.args.accountNumber, tt.args.sequenceNumber)
			tt.assertion(t, err)

			assert.Equal(t, tt.want, got)
		})
	}
}
