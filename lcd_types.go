package sacco

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

// TxBody represents the body of a Cosmos transaction
// signed and ready to be sent over the LCD REST service.
type TxBody struct {
	Tx   SignedTransactionPayload `json:"tx"`
	Mode string                   `json:"mode"`
}

// AccountData holds informations about the account number and
// sequence number of a Cosmos account.
type AccountData struct {
	Result AccountDataResult `json:"result"`
}

// AccountDataResult is a wrapper struct for a call to auth/accounts/{address} LCD
// REST endpoint.
type AccountDataResult struct {
	Value AccountDataValue `json:"value"`
}

// AccountDataValue represents the real data obtained by calling /auth/accounts/{address} LCD
// REST endpoint.
type AccountDataValue struct {
	Address       string `json:"address"`
	AccountNumber string `json:"account_number"`
	Sequence      string `json:"sequence"`
}

// NodeInfo is the LCD REST response to a /node_info request,
// and contains the Network attribute (chain ID).
type NodeInfo struct {
	Info struct {
		Network string `json:"network"`
	} `json:"node_info"`
}

// RawLog is the log format returned by the LCD REST service whenever a
// transaction doesn't meet the quality of life properties specified in the
// transaction type itself.
type RawLog struct {
	Codespace string `json:"codespace,omitempty"`
	Code      int    `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
}

// TxResponse represents whatever data the LCD REST service returns to atomicwallet
// after a transaction gets forwarded to it.
type TxResponse struct {
	Height    string                   `json:"height"`
	TxHash    string                   `json:"txhash"`
	Code      uint32                   `json:"code,omitempty"`
	Data      string                   `json:"data,omitempty"`
	RawLog    string                   `json:"raw_log,omitempty"`
	Logs      sdkTypes.ABCIMessageLogs `json:"logs,omitempty"`
	Info      string                   `json:"info,omitempty"`
	GasWanted string                   `json:"gas_wanted,omitempty"`
	GasUsed   string                   `json:"gas_used,omitempty"`
	Codespace string                   `json:"codespace,omitempty"`
	Tx        sdkTypes.Tx              `json:"tx,omitempty"`
	Timestamp string                   `json:"timestamp,omitempty"`

	// DEPRECATED: Remove in the next next major release in favor of using the
	// ABCIMessageLog.Events field.
	Events sdkTypes.StringEvents `json:"events,omitempty"`
}

// Error represents a JSON encoded error message sent whenever something
// goes wrong during the handler processing.
type Error struct {
	Error string `json:"error,omitempty"`
}
