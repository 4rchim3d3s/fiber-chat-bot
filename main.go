package main

import (
	"flag"
	"fmt"
	"time"

	"golang.org/x/text/language"
)

const PRINT = false

var NUMBER_FORMAT = language.English

func main() {
	var startBlock int
	var config_filename string
	flag.IntVar(&startBlock, "block", 0, "the first block number to be querried from the fiber node")
	flag.StringVar(&config_filename, "cfg", "default_config.json", "the config to run the bot")
	flag.Parse()

	if config_filename == "" || config_filename == "default_config.json" {
		fmt.Println("No config file defined. Rename 'default_config.json', edit it correctly and define it with -cfg")
		return
	}

	if startBlock == 0 {
		fmt.Println("No start block defined. Specify block number >1 with -block")
	}

	cfg, err := readConfig(config_filename)
	if err != nil {
		fmt.Println(err)
	}

	lastBlock := startBlock

	fmt.Printf("%s Starting bot from block %d\n", time.Now().Local().Format(time.RFC1123), lastBlock)
	fmt.Printf("%s Query Intervall: %d seconds\n", time.Now().Local().Format(time.RFC1123), cfg.FiberNode.QueryIntervallSeconds)

	for {
		lastBlock = function(lastBlock, cfg)
		time.Sleep(time.Duration(cfg.FiberNode.QueryIntervallSeconds) * time.Second)
	}
}

func function(lastBlock int, cfg Config) (newlastBlock int) {

	metadata, _ := queryBlockchainMetadata(cfg)
	blockheight := getBlockHeight(metadata)

	fmt.Printf("%s Blockheight: %d\n", time.Now().Local().Format(time.RFC1123), blockheight)

	for lastBlock <= blockheight {
		fmt.Printf("%s New Block:   %d *\n", time.Now().Local().Format(time.RFC1123), lastBlock)
		block, _ := queryBlockData(lastBlock, cfg)
		prettyBlockString := prettyPrintBlock(block, cfg)
		if PRINT {
			fmt.Println(prettyBlockString)
		}
		fmt.Printf("%s Send Block:  %d to chat\n", time.Now().Local().Format(time.RFC1123), lastBlock)
		_, retry_after, err := sendTextToTelegramChat(cfg.Telegram.BotToken, cfg.Telegram.ChatID, prettyBlockString)
		if err != nil {
			fmt.Println(err)
			fmt.Printf("%s Waiting %d seconds\n", time.Now().Local().Format(time.RFC1123), retry_after)
			time.Sleep(time.Duration(retry_after) * time.Second)
			_, _, err := sendTextToTelegramChat(cfg.Telegram.BotToken, cfg.Telegram.ChatID, prettyBlockString)
			if err != nil {
				fmt.Println(err)
			}
		}

		lastBlock = lastBlock + 1
		time.Sleep(1 * time.Second)
	}
	return lastBlock
}
