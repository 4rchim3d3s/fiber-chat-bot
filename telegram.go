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

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var NUMBER_FORMAT = language.English

type TelegramResponse struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
	Parameters  struct {
		RetryAfter int `json:"retry_after"`
	} `json:"parameters"`
}

// sendTextToTelegramChat sends a text message to the Telegram chat identified by its chat Id
func sendTextToTelegramChat(chatId int, text string) (string, int, error) {

	//log.Printf("Sending to chat_id: %d\n Message:\n %s\n", chatId, text)
	var telegramApi string = "https://api.telegram.org/bot" + TELEGRAM_TOKEN + "/sendMessage"
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

func prettyPrintBlock(block BlockchainBlock) (prettyPrint string) {
	p := message.NewPrinter(NUMBER_FORMAT)

	knownAdds, err := getAddresses()
	if err != nil {
		//handle error
	}

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

	if FIBER_PRINT_WHALE_EMOJI {
		symbols, _ := bits.Div(0, uint(fCoins), FIBER_WHALE_EMOJI_AMOUNT)
		for i := 0; i < int(symbols); i++ {
			prettyPrint = prettyPrint + fmt.Sprintf("<tg-emoji emoji-id=%s>%s</tg-emoji>", FIBER_WHALE_TELEGRAM_CUSTOM_EMOJI_ID, FIBER_WHALE_TELEGRAM_DEFAULT_EMOJI)
		}
		if symbols != 0 {
			prettyPrint = prettyPrint + "\n"
		}
	}

	////add Coins

	coins := p.Sprintf("%f", fCoins)
	if strings.HasSuffix(coins, ".000000") {
		prettyPrint = prettyPrint + fmt.Sprintf("%s: %s\n", FIBER_COIN_SYM, coins[:len(coins)-7])
	} else {
		prettyPrint = prettyPrint + fmt.Sprintf("%s: %s\n", FIBER_COIN_SYM, coins[:len(coins)-(6-(6-FIBER_PRINT_PRECISION))])
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

	prettyPrint = prettyPrint + fmt.Sprintf("%s: %s %s\n", FIBER_HOUR_SYM, hours, getHoursEmoji(hoursOutput))

	//// add FEE

	fee := p.Sprintf("%d", block.Header.Fee)
	if FIBER_PRINT_FEE_EMOJI {
		prettyPrint = prettyPrint + fmt.Sprintf("FEE: %s %s %s %s\n", FIBER_FEE_TELEGRAM_DEFAULT_EMOJI, fee, FIBER_HOUR_SYM, getHoursEmoji(block.Header.Fee))

	} else {
		prettyPrint = prettyPrint + fmt.Sprintf("FEE: %s %s\n", fee, FIBER_HOUR_SYM)
	}

	////add SND

	inp := block.Body.Txns[0].Inputs[0].Owner
	vHasAlias, alias, emoji, emoji_id := hasAlias(inp, knownAdds)
	if vHasAlias {
		if !ADDRESS_ALIAS_PRINT_EMOJI {
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
		vHasAlias, alias, emoji, emoji_id := hasAlias(outp, knownAdds)
		if vHasAlias {
			if !ADDRESS_ALIAS_PRINT_EMOJI {
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
	prettyPrint = prettyPrint + fmt.Sprintf("<a href='%s/app/transaction/%s'>Explorer</a>", EXPLORER_PUBLIC_URL, block.Body.Txns[0].Txid)

	return prettyPrint
}

func getHoursEmoji(hours int) (emoji string) {
	if hours >= (1000000 * FIBER_HOUR_EMOJI_AMOUNT) {
		emoji = FIBER_HOUR_EMOJI_1M
	} else if hours >= (500000 * FIBER_HOUR_EMOJI_AMOUNT) {
		emoji = FIBER_HOUR_EMOJI_500K
	} else if hours >= (100000 * FIBER_HOUR_EMOJI_AMOUNT) {
		emoji = FIBER_HOUR_EMOJI_100K
	} else if hours >= (10000 * FIBER_HOUR_EMOJI_AMOUNT) {
		emoji = FIBER_HOUR_EMOJI_10K
	} else if hours >= (1000 * FIBER_HOUR_EMOJI_AMOUNT) {
		emoji = FIBER_HOUR_EMOJI_1K
	} else if hours >= (100 * FIBER_HOUR_EMOJI_AMOUNT) {
		emoji = FIBER_HOUR_EMOJI_100
	} else if hours >= (10 * FIBER_HOUR_EMOJI_AMOUNT) {
		emoji = FIBER_HOUR_EMOJI_10
	} else if hours >= (1 * FIBER_HOUR_EMOJI_AMOUNT) {
		emoji = FIBER_HOUR_EMOJI_1
	}

	return emoji
}
