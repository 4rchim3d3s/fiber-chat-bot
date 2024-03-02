package main

import (
	"fmt"
	"time"
)

func main() {

	//TODO: get lastBlock from persistent memory and save the new one
	lastBlock := 186768

	lastBlock = lastBlock + 1

	fmt.Printf("%s Starting bot from block %d\n", time.Now().Local().Format(time.RFC1123), lastBlock)
	fmt.Printf("%s Query Intervall: %d seconds\n", time.Now().Local().Format(time.RFC1123), LOCAL_FIBER_NODE_QUERY_INTERVALL_SECONDS)

	for {
		lastBlock = function(lastBlock)
		time.Sleep(time.Duration(LOCAL_FIBER_NODE_QUERY_INTERVALL_SECONDS) * time.Second)
	}
}

func function(lastBlock int) (newlastBlock int) {

	metadata, _ := queryBlockchainMetadata()
	blockheight := getBlockHeight(metadata)

	fmt.Printf("%s Blockheight: %d\n", time.Now().Local().Format(time.RFC1123), blockheight)

	for lastBlock <= blockheight {
		fmt.Printf("%s New Block:   %d *\n", time.Now().Local().Format(time.RFC1123), lastBlock)
		block, _ := queryBlockData(lastBlock)
		prettyBlockString := prettyPrintBlock(block)
		if PRINT {
			fmt.Println(prettyBlockString)
		}
		fmt.Printf("%s Send Block:  %d to chat\n", time.Now().Local().Format(time.RFC1123), lastBlock)
		_, retry_after, err := sendTextToTelegramChat(TELEGRAM_CHAT_ID, prettyBlockString)
		if err != nil {
			fmt.Println(err)
			fmt.Printf("%s Waiting %d seconds\n", time.Now().Local().Format(time.RFC1123), retry_after)
			time.Sleep(time.Duration(retry_after) * time.Second)
			_, _, err := sendTextToTelegramChat(TELEGRAM_CHAT_ID, prettyBlockString)
			if err != nil {
				fmt.Println(err)
			}
		}

		lastBlock = lastBlock + 1
		time.Sleep(1 * time.Second)
	}
	return lastBlock
}
