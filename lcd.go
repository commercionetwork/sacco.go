package sacco

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Retrieve the account data related to the given wallet address, like
// account number and sequence number.
func getAccountData(lcdEndpoint, address string) (AccountData, error) {
	endpoint := fmt.Sprintf("%s/auth/accounts/%s", lcdEndpoint, address)

	// nolint:gosec
	resp, err := http.Get(endpoint)
	if err != nil {
		return AccountData{}, err
	}

	jdec := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		// we had an error, deserialize it and return
		var jsonError Error
		err := jdec.Decode(&jsonError)
		if err != nil {
			return AccountData{}, fmt.Errorf("error deserializing account data JSON error: %w", err)
		}

		return AccountData{}, fmt.Errorf("error during get account data: %s", jsonError.Error)
	}

	var accountData AccountData

	errCdc := jdec.Decode(&accountData)
	if errCdc != nil {
		return AccountData{}, fmt.Errorf("could not unmarshal node response: %w", errCdc)
	}

	if accountData.Result.Value.Address == "" {
		return AccountData{}, fmt.Errorf("account with address %s is not online", address)
	}

	return accountData, nil
}

// Return useful information of the full node, like the Network
// (chain) name.
func getNodeInfo(lcdEndpoint string) (NodeInfo, error) {
	endpoint := fmt.Sprintf("%s/node_info", lcdEndpoint)

	// nolint:gosec
	resp, err := http.Get(endpoint)
	if err != nil {
		return NodeInfo{}, err
	}

	var nodeInfo NodeInfo
	jdec := json.NewDecoder(resp.Body)
	err = jdec.Decode(&nodeInfo)
	if err != nil {
		return NodeInfo{}, err
	}

	return nodeInfo, nil
}
