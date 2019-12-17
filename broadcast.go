package sacco

//go:generate stringer -type=TxMode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// TxMode identifies when an LCD should replies to a client
// after a transaction broadcast.
type TxMode string

// String implements the stringer interface for TxMode.
func (txm TxMode) String() string {
	return string(txm)
}

const (
	// ModeAsync waits for the tx to pass/fail CheckTx
	ModeAsync TxMode = "async"

	// ModeSync doesn't wait for pass/fail CheckTx and send and return tx immediately
	ModeSync TxMode = "sync"

	// ModeBlock waits for the tx to pass/fail CheckTx, DeliverTx, and be committed in a block (not recommended, slow)
	ModeBlock TxMode = "block"
)

// broadcastTx broadcasts a tx to the Cosmos LCD identified by lcdEndpoint.
func broadcastTx(tx SignedTransactionPayload, lcdEndpoint string, txMode TxMode) (string, error) {
	endpoint := fmt.Sprintf("%s/txs", lcdEndpoint)

	// assemble a tx transaction
	txBody := TxBody{
		Tx:   tx,
		Mode: txMode.String(),
	}
	requestBody, err := json.Marshal(txBody)
	if err != nil {
		return "", err
	}

	// send tx to lcdEndpoint
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	// we had a serious problem with the LCD request,
	// decode the json error and return its value to caller
	if resp.StatusCode != http.StatusOK {
		var jerr Error
		jd := json.NewDecoder(resp.Body)
		err := jd.Decode(&jerr)
		if err != nil {
			return "", fmt.Errorf("could not process error json decoding: %w", err)
		}

		return "", fmt.Errorf("error while processing tx send request: %s", jerr.Error)
	}

	// deserialize LCD response into a cosmos TxResponse
	var txr TxResponse

	txrJdec := json.NewDecoder(resp.Body)
	err = txrJdec.Decode(&txr)
	if err != nil {
		return "", fmt.Errorf("could not deserialize cosmos txresponse from lcd: %w", err)
	}

	// something went wrong in tx validation on the LCD side,
	// deserialize raw_log and return the error inside
	allOk := true
	for _, l := range txr.Logs {
		allOk = allOk && l.Success
	}

	if !allOk || len(txr.Logs) <= 0 {
		// txr.RawLog can either be a RawLog
		// or it might be a lot of good stuff,
		// we better cycle txr.Logs and build a nice error message
		var jerr RawLog
		jd := json.NewDecoder(strings.NewReader(txr.RawLog))
		jsonErr := jd.Decode(&jerr)

		if jsonErr == nil {
			return "", fmt.Errorf(
				"codespace %s: %s, code %d",
				jerr.Codespace,
				jerr.Message,
				jerr.Code,
			)
		}

		// parse logs!
		message := ""
		for i, log := range txr.Logs {
			var jerr RawLog
			jd := json.NewDecoder(strings.NewReader(log.Log))
			jsonErr := jd.Decode(&jerr)

			if jsonErr != nil {
				// bail out
				message = txr.RawLog
				break
			}

			if i == 0 {
				message = message + jerr.Message
			} else {
				message = message + ", " + jerr.Message
			}
		}

		return "", fmt.Errorf(
			"codespace %s: %s, code %d",
			txr.Codespace,
			message,
			txr.Code,
		)
	}

	return txr.TxHash, nil
}

// SignAndBroadcast signs tx and broadcast it to the LCD specified by lcdEndpoint.
func (w *Wallet) SignAndBroadcast(tx TransactionPayload, lcdEndpoint string, txMode TxMode) (string, error) {
	// get network (chain) name
	nodeInfo, err := getNodeInfo(lcdEndpoint)
	if err != nil {
		return "", fmt.Errorf("could not get LCD node informations: %w", err)
	}

	// get account sequence and account number
	accountData, err := getAccountData(lcdEndpoint, w.Address)
	if err != nil {
		return "", fmt.Errorf("could not get Account informations for address %s: %w", w.Address, err)
	}

	// sign transaction
	signedTx, err := w.Sign(
		tx,
		nodeInfo.Info.Network,
		accountData.Result.Value.AccountNumber,
		accountData.Result.Value.Sequence,
	)
	if err != nil {
		return "", fmt.Errorf("could not sign transaction: %w", err)
	}

	// broadcast transaction to the LCD
	txHash, err := broadcastTx(signedTx, lcdEndpoint, txMode)
	if err != nil {
		return "", fmt.Errorf("could not broadcast transaction to the Cosmos network: %w", err)
	}

	// return transaction hash!
	return txHash, nil
}
