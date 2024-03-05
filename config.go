package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	Explorer struct {
		PublicURL string `json:"public_url"`
	} `json:"explorer"`
	FiberNode struct {
		URL                   string `json:"url"`
		Port                  string `json:"port"`
		QueryIntervallSeconds int    `json:"query_intervall_seconds"`
		FiberCoinTicker       string `json:"fiber_coin_ticker"`
		FiberHoursTicker      string `json:"fiber_hours_ticker"`
		FiberPrintPrecision   int    `json:"fiber_print_precision"`
	} `json:"fiber_node"`
	Telegram struct {
		BotToken string `json:"bot_token"`
		ChatID   int    `json:"chat_id"`
	} `json:"telegram"`
	Emojis struct {
		PrintCoinEmojis           bool   `json:"print_coin_emojis"`
		CoinEmojiDivisor          int    `json:"coin_emoji_divisor"`
		CoinEmoji                 string `json:"coin_emoji"`
		CoinCustomEmojiTelegramID string `json:"coin_custom_emoji_telegram_id"`
		PrintHoursEmoji           bool   `json:"print_hours_emoji"`
		HoursEmojiMultiplicator   int    `json:"hours_emoji_multiplicator"`
		HoursEmoji1               string `json:"hours_emoji_1"`
		HoursEmoji10              string `json:"hours_emoji_10"`
		HoursEmoji100             string `json:"hours_emoji_100"`
		HoursEmoji1K              string `json:"hours_emoji_1K"`
		HoursEmoji10K             string `json:"hours_emoji_10K"`
		HoursEmoji100K            string `json:"hours_emoji_100K"`
		HoursEmoji500K            string `json:"hours_emoji_500K"`
		HoursEmoji1M              string `json:"hours_emoji_1M"`
		PrintFeeEmoji             bool   `json:"print_fee_emoji"`
		FeeEmoji                  string `json:"fee_emoji"`
	} `json:"emojis"`
	Addressbook struct {
		PrintAliasEmoji bool `json:"print_alias_emoji"`
		Addresses       []struct {
			Pk                    string `json:"pk"`
			Alias                 string `json:"alias"`
			Emoji                 string `json:"emoji"`
			CustomEmojiTelegramID string `json:"custom_emoji_telegram_id"`
		} `json:"addresses"`
	} `json:"addressbook"`
}

func readConfig(filedir string) (cfg Config, err error) {
	jsonFile, err := os.Open(filedir)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		return cfg, err
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := io.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &cfg)

	return cfg, nil
}

func hasAlias(pk string, cfg Config) (has bool, alias string, emoji string, emoji_id string) {
	for _, add := range cfg.Addressbook.Addresses {
		if add.Pk == pk {
			return true, add.Alias, add.Emoji, add.CustomEmojiTelegramID
		}
	}
	return false, alias, emoji, emoji_id
}
