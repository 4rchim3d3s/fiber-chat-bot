package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type BlockchainMetadata struct {
	Head struct {
		Seq               int    `json:"seq"`
		BlockHash         string `json:"block_hash"`
		PreviousBlockHash string `json:"previous_block_hash"`
		Timestamp         int    `json:"timestamp"`
		Fee               int    `json:"fee"`
		Version           int    `json:"version"`
		TxBodyHash        string `json:"tx_body_hash"`
	} `json:"head"`
	Unspents    int `json:"unspents"`
	Unconfirmed int `json:"unconfirmed"`
}

type BlockchainBlock struct {
	Header struct {
		Seq               int    `json:"seq"`
		BlockHash         string `json:"block_hash"`
		PreviousBlockHash string `json:"previous_block_hash"`
		Timestamp         int    `json:"timestamp"`
		Fee               int    `json:"fee"`
		Version           int    `json:"version"`
		TxBodyHash        string `json:"tx_body_hash"`
		UxHash            string `json:"ux_hash"`
	} `json:"header"`
	Body struct {
		Txns []struct {
			Length    int      `json:"length"`
			Type      int      `json:"type"`
			Txid      string   `json:"txid"`
			InnerHash string   `json:"inner_hash"`
			Fee       int      `json:"fee"`
			Sigs      []string `json:"sigs"`
			Inputs    []struct {
				Uxid            string `json:"uxid"`
				Owner           string `json:"owner"`
				Coins           string `json:"coins"`
				SrcTxid         string `json:"src_txid"`
				Hours           int    `json:"hours"`
				CalculatedHours int    `json:"calculated_hours"`
			} `json:"inputs"`
			Outputs []struct {
				Uxid  string `json:"uxid"`
				Dst   string `json:"dst"`
				Coins string `json:"coins"`
				Hours int    `json:"hours"`
			} `json:"outputs"`
		} `json:"txns"`
	} `json:"body"`
	Size int `json:"size"`
}

func queryBlockData(block int) (resp BlockchainBlock, err error) {

	url := "http://" + LOCAL_FIBER_NODE_URL + ":" + LOCAL_FIBER_NODE_PORT + "/api/v1/block?seq=" + strconv.Itoa(block) + "&verbose=1"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	if PRINT {
		fmt.Println(string(body))
	}

	var result BlockchainBlock
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	return result, nil
}

func queryBlockchainMetadata() (resp BlockchainMetadata, err error) {

	url := "http://" + LOCAL_FIBER_NODE_URL + ":" + LOCAL_FIBER_NODE_PORT + "/api/v1/blockchain/metadata"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	if PRINT {
		fmt.Println(string(body))
	}

	var result BlockchainMetadata
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	return result, nil
}

func getBlockHeight(metadata BlockchainMetadata) (blockheight int) {
	return metadata.Head.Seq
}
