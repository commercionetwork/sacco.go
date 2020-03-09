package sacco

//go:generate stringer -type=TxMode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

	jdec := json.NewDecoder(resp.Body)

	err = jdec.Decode(&txr)
	if err != nil {
		return "", fmt.Errorf("could not deserialize cosmos txresponse from lcd: %w", err)
	}

	if txr.Code != 0 || len(txr.Logs) > 0 {
		return "", fmt.Errorf(
			"codespace %s: %s, code %d",
			txr.Codespace,
			txr.RawLog,
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
		strconv.FormatInt(accountData.Result.Value.AccountNumber, 10),
		strconv.FormatInt(accountData.Result.Value.Sequence, 10),
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
