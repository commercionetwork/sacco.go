package sacco

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBech32AminoPubKey(t *testing.T) {
	tests := []struct {
		name    string
		hrp     string
		data    []byte
		want    string
		wantErr bool
	}{
		{
			"argument length is not 33",
			"",
			[]byte{1, 2},
			"",
			true,
		},
		{
			"hrp empty",
			"",
			[]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
			"",
			true,
		},
		{
			"all ok",
			"did:com:",
			[]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
			"did:com:pub1addwnpepqqqsyqcyq5rqwzqfpg9scrgwpugpzysnzs23v9ccrydpk8qarc0jq8mqqvq",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Bech32AminoPubKey(tt.data, tt.hrp)

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

func TestSignBytes(t *testing.T) {
	type args struct {
		tx             TransactionPayload
		chainID        string
		accountNumber  string
		sequenceNumber string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			"SignBytes creates the payload as instructed",
			args{
				tx:             TransactionPayload{},
				chainID:        "1",
				accountNumber:  "2",
				sequenceNumber: "3",
			},
			[]byte{0x7b, 0x22, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x22, 0x3a, 0x22, 0x32, 0x22, 0x2c, 0x22, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x5f, 0x69, 0x64, 0x22, 0x3a, 0x22, 0x31, 0x22, 0x2c, 0x22, 0x66, 0x65, 0x65, 0x22, 0x3a, 0x7b, 0x22, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x3a, 0x6e, 0x75, 0x6c, 0x6c, 0x2c, 0x22, 0x67, 0x61, 0x73, 0x22, 0x3a, 0x22, 0x22, 0x7d, 0x2c, 0x22, 0x6d, 0x65, 0x6d, 0x6f, 0x22, 0x3a, 0x22, 0x22, 0x2c, 0x22, 0x6d, 0x73, 0x67, 0x73, 0x22, 0x3a, 0x6e, 0x75, 0x6c, 0x6c, 0x2c, 0x22, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x22, 0x3a, 0x22, 0x33, 0x22, 0x7d},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, SignBytes(
				tt.args.tx,
				tt.args.chainID,
				tt.args.accountNumber,
				tt.args.sequenceNumber,
			))
		})
	}
}
