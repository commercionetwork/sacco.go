package sacco

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func Test_getAccountData(t *testing.T) {
	mockHTTPEndpoint := "http://127.0.0.1:3333/"
	tests := []struct {
		name       string
		address    string
		jsonResp   string
		statusResp int
		want       AccountData
		assertion  assert.ErrorAssertionFunc
	}{
		{
			"successful request on a live account",
			"did:com:1sfjela2snk9rmmcfh773gm50476w0ur5pmwuak",
			`{"height":"1590","result":{"type":"cosmos-sdk/Account","value":{"address":"did:com:1sfjela2snk9rmmcfh773gm50476w0ur5pmwuak","coins":[{"denom":"ucommercio","amount":"10"}],"public_key":null,"account_number":"11","sequence":"0"}}}`,
			http.StatusOK,
			AccountData{
				AccountDataResult{
					AccountDataValue{
						Sequence:      "0",
						AccountNumber: "11",
						Address:       "did:com:1sfjela2snk9rmmcfh773gm50476w0ur5pmwuak",
					},
				},
			},
			assert.NoError,
		},
		{
			"successful request on a non-live account",
			"did:com:13lsdhm9gmxhmm0lksvv042ufx2ykwfqj2julet",
			`{"height":"1809","result":{"type":"cosmos-sdk/Account","value":{"address":"","coins":[],"public_key":null,"account_number":"0","sequence":"0"}}}`,
			http.StatusOK,
			AccountData{
				AccountDataResult{
					AccountDataValue{
						Sequence:      "",
						AccountNumber: "",
						Address:       "",
					},
				},
			},
			assert.Error,
		},
		{
			"unsuccessful request with a JSON error",
			"fakeaccount",
			`{"error":"decoding bech32 failed: invalid index of 1"}`,
			http.StatusInternalServerError,
			AccountData{},
			assert.Error,
		},
		{
			"unsuccessful request with a malformed error",
			"fakeaccount",
			`malformed error`,
			http.StatusInternalServerError,
			AccountData{},
			assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			httpmock.RegisterResponder("GET", mockHTTPEndpoint+"/auth/accounts/"+tt.address,
				httpmock.NewStringResponder(tt.statusResp, tt.jsonResp))

			got, err := getAccountData(mockHTTPEndpoint, tt.address)

			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_getNodeInfo(t *testing.T) {
	mockHTTPEndpoint := "http://127.0.0.1:3333/"

	tests := []struct {
		name       string
		jsonResp   string
		statusResp int
		want       NodeInfo
		assertion  assert.ErrorAssertionFunc
	}{
		{
			"successful call",
			`{"node_info":{"protocol_version":{"p2p":"7","block":"10","app":"0"},"id":"4bc6d5af186f705316620bffd5e2cefaea11bd59","listen_addr":"tcp://0.0.0.0:26656","network":"test-chain-jVvnJ6","version":"0.32.7","channels":"4020212223303800","moniker":"testchain","other":{"tx_index":"on","rpc_address":"tcp://127.0.0.1:26657"}},"application_version":{"name":"commercionetwork","server_name":"cnd","client_name":"cndcli","version":"1.3.3-9-gef69043","commit":"ef69043933adaefed4803c5032f27c8ab6280bbd","build_tags":"netgo","go":"go version go1.13.4 darwin/amd64"}}`,
			http.StatusOK,
			NodeInfo{
				struct{Network string "json:\"network\""}{
					"test-chain-jVvnJ6",
				},
			},
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			httpmock.RegisterResponder("GET", mockHTTPEndpoint+"/node_info",
				httpmock.NewStringResponder(tt.statusResp, tt.jsonResp))

			got, err := getNodeInfo(mockHTTPEndpoint)

			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
