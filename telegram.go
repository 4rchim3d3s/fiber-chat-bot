package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/bits"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/message"
)

type TelegramResponse struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
	Parameters  struct {
		RetryAfter int `json:"retry_after"`
	} `json:"parameters"`
}

// sendTextToTelegramChat sends a text message to the Telegram chat identified by its chat Id
func sendTextToTelegramChat(botToken string, chatId int, text string) (string, int, error) {

	//log.Printf("Sending to chat_id: %d\n Message:\n %s\n", chatId, text)
	var telegramApi string = "https://api.telegram.org/bot" + botToken + "/sendMessage"
	response, err := http.PostForm(
		telegramApi,
		url.Values{
			"chat_id":    {strconv.Itoa(chatId)},
			"text":       {text},
			"parse_mode": {"HTML"},
		})

	if err != nil {
		log.Printf("error when posting text to the chat: %s", err.Error())
		return "", 1, err
	}
	defer response.Body.Close()

	var bodyBytes, errRead = io.ReadAll(response.Body)
	if errRead != nil {
		log.Printf("error in parsing telegram answer %s", errRead.Error())
		return "", 1, err
	}
	bodyString := string(bodyBytes)
	//log.Printf("Body of Telegram Response:\n %s", bodyString)

	var tgr TelegramResponse
	json.Unmarshal(bodyBytes, &tgr)

	if !tgr.Ok {
		return "", tgr.Parameters.RetryAfter, fmt.Errorf(tgr.Description)
	}

	return bodyString, 1, nil
}

func prettyPrintBlock(block BlockchainBlock, cfg Config) (prettyPrint string) {
	p := message.NewPrinter(NUMBER_FORMAT)

	//TODO: for each transaction --> are there blocks with more than one transaction?

	//// add block number

	blocknumber := p.Sprintf("%d", block.Header.Seq)
	prettyPrint = fmt.Sprintf("<b>Block: %s</b>\n", blocknumber)

	//// add transaction

	txid := block.Body.Txns[0].Txid
	prettyPrint = prettyPrint + fmt.Sprintf("TXN: %s\n", txid[0:4]+"..."+txid[len(txid)-4:])

	////add Coin Whale Emoji

	var fCoinsInput float64
	fCoinsInputOwner := block.Body.Txns[0].Inputs[0].Owner
	for _, input := range block.Body.Txns[0].Inputs {
		_fCoins, err := strconv.ParseFloat(input.Coins, 64)
		if err != nil {
			//Handle error
		}

		fCoinsInput = fCoinsInput + _fCoins
	}

	var fCoinsOutput float64
	for _, output := range block.Body.Txns[0].Outputs {
		if output.Dst == fCoinsInputOwner {
			_fCoins, err := strconv.ParseFloat(output.Coins, 64)
			if err != nil {
				//Handle error
			}

			fCoinsOutput = fCoinsOutput + _fCoins
		}
	}

	fCoins := fCoinsInput - fCoinsOutput

	if cfg.Emojis.PrintCoinEmojis {
		symbols, _ := bits.Div(0, uint(fCoins), uint(cfg.Emojis.CoinEmojiDivisor))
		for i := 0; i < int(symbols); i++ {
			prettyPrint = prettyPrint + fmt.Sprintf("<tg-emoji emoji-id=%s>%s</tg-emoji>", cfg.Emojis.CoinCustomEmojiTelegramID, cfg.Emojis.CoinEmoji)
		}
		if symbols != 0 {
			prettyPrint = prettyPrint + "\n"
		}
	}

	////add Coins
	coins := p.Sprintf("%f", fCoins)
	if strings.HasSuffix(coins, ".000000") {
		prettyPrint = prettyPrint + fmt.Sprintf("%s: %s\n", cfg.FiberNode.FiberCoinTicker, coins[:len(coins)-7])
	} else {
		prettyPrint = prettyPrint + fmt.Sprintf("%s: %s\n", cfg.FiberNode.FiberCoinTicker, coins[:len(coins)-(6-cfg.FiberNode.FiberPrintPrecision)])
	}

	////add SCH

	var hoursOutput int
	hoursInputOwner := block.Body.Txns[0].Inputs[0].Owner
	for _, output := range block.Body.Txns[0].Outputs {
		if output.Dst != hoursInputOwner {
			hoursOutput = hoursOutput + output.Hours
		}
	}

	hours := p.Sprintf("%d", hoursOutput)

	prettyPrint = prettyPrint + fmt.Sprintf("%s: %s %s\n", cfg.FiberNode.FiberHoursTicker, hours, getHoursEmoji(hoursOutput, cfg))

	//// add FEE

	fee := p.Sprintf("%d", block.Header.Fee)
	if cfg.Emojis.PrintFeeEmoji {
		prettyPrint = prettyPrint + fmt.Sprintf("FEE: %s %s %s %s\n", cfg.Emojis.FeeEmoji, fee, cfg.FiberNode.FiberHoursTicker, getHoursEmoji(block.Header.Fee, cfg))

	} else {
		prettyPrint = prettyPrint + fmt.Sprintf("FEE: %s %s\n", fee, cfg.FiberNode.FiberHoursTicker)
	}

	////add SND

	inp := block.Body.Txns[0].Inputs[0].Owner
	vHasAlias, alias, emoji, emoji_id := hasAlias(inp, cfg)
	if vHasAlias {
		if !cfg.Addressbook.PrintAliasEmoji {
			emoji = ""
		} else if emoji_id != "" {
			emoji = fmt.Sprintf("<tg-emoji emoji-id=%s>%s</tg-emoji>", emoji_id, emoji)
		}
		prettyPrint = prettyPrint + fmt.Sprintf("SND: %s #%s \n", emoji, alias)
	} else {
		prettyPrint = prettyPrint + fmt.Sprintf("SND: %s \n", inp[0:4]+"..."+inp[len(inp)-4:])
	}

	////add RCV

	inputOwner := block.Body.Txns[0].Inputs[0].Owner
	var outputs int
	for _, output := range block.Body.Txns[0].Outputs {
		if output.Dst != inputOwner {
			outputs = outputs + 1
		}
	}
	outp := block.Body.Txns[0].Outputs[0].Dst

	if outputs > 1 {
		prettyPrint = prettyPrint + fmt.Sprintf("RCV: %s <i>(+ %d more)</i>", outp[0:4]+"..."+outp[len(outp)-4:], outputs-1)
	} else {
		vHasAlias, alias, emoji, emoji_id := hasAlias(outp, cfg)
		if vHasAlias {
			if !cfg.Addressbook.PrintAliasEmoji {
				emoji = ""
			} else if emoji_id != "" {
				emoji = fmt.Sprintf("<tg-emoji emoji-id=%s>%s</tg-emoji>", emoji_id, emoji)
			}
			prettyPrint = prettyPrint + fmt.Sprintf("RCV: %s #%s ", emoji, alias)
		} else {
			prettyPrint = prettyPrint + fmt.Sprintf("RCV: %s ", outp[0:4]+"..."+outp[len(outp)-4:])
		}
	}
	prettyPrint = prettyPrint + "\n"

	////add Block-Details

	prettyPrint = prettyPrint + "_____________________\n"
	datetime := time.Unix(int64(block.Header.Timestamp), 0).UTC().Format(time.RFC822)
	prettyPrint = prettyPrint + fmt.Sprintf("<i>%s</i>\n", datetime)
	prettyPrint = prettyPrint + fmt.Sprintf("<a href='%s/app/transaction/%s'>Explorer</a>", cfg.Explorer.PublicURL, block.Body.Txns[0].Txid)

	return prettyPrint
}

func getHoursEmoji(hours int, cfg Config) (emoji string) {

	amount := cfg.Emojis.HoursEmojiMultiplicator

	if hours >= (1000000 * amount) {
		emoji = cfg.Emojis.HoursEmoji1M
	} else if hours >= (500000 * amount) {
		emoji = cfg.Emojis.HoursEmoji500K
	} else if hours >= (100000 * amount) {
		emoji = cfg.Emojis.HoursEmoji100K
	} else if hours >= (10000 * amount) {
		emoji = cfg.Emojis.HoursEmoji10K
	} else if hours >= (1000 * amount) {
		emoji = cfg.Emojis.HoursEmoji1K
	} else if hours >= (100 * amount) {
		emoji = cfg.Emojis.HoursEmoji100
	} else if hours >= (10 * amount) {
		emoji = cfg.Emojis.HoursEmoji10
	} else if hours >= (1 * amount) {
		emoji = cfg.Emojis.HoursEmoji1
	}

	return emoji
}
