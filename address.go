package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Addresses struct {
	Addresses []Address `json:"addresses"`
}

type Address struct {
	Pk      string `json:"Pk"`
	Alias   string `json:"Alias"`
	Emoji   string `json:"Emoji"`
	EmojiId string `json:"Emoji_Id"`
}

func getAddresses() (adds Addresses, err error) {
	jsonFile, err := os.Open("addresses.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		return adds, err
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := io.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &adds)

	return adds, nil
}

func hasAlias(pk string, knownAdds Addresses) (has bool, alias string, emoji string, emoji_id string) {
	for _, add := range knownAdds.Addresses {
		if add.Pk == pk {
			return true, add.Alias, add.Emoji, add.EmojiId
		}
	}
	return false, alias, emoji, emoji_id
}
