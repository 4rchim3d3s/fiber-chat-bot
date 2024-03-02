# Fiber Chat Bot (Telegram)

Fiber Chat Bot is a bot that queries a local fiber node and when detecting a new block broadcasts it via chat (currently only telegram).

The bot currently supports a local fiber node that is accessible via HTTP. \
HTTPS might be implemented in a later version

![image](/images/message_example1.png "Skycoin Blockhain Bot Message")


# Table of Contents

<!-- MarkdownTOC levels="1,2,3,4,5" autolink="true" bracket="round" -->

- [Installation](#installation)
    - [Software requirements](#software-requirements)
    - [Initial configuration steps](#initial-configuration-steps)


## Installation

### Software Requirements

The bot depends on a local instance of a fiber node.\
See information about a fiber node on https://github.com/skycoin/skycoin.git.

### Installing the Bot

TODO

### Initial Configuration Steps

#### Fiber Configuration

To configure the bot for your own fiber coin you have to edit the file `config_variables.go`
If you have known addresses and their alias you can add them to `addresses.json`

#### Telegram Configuration

##### New Bot

1. Open Telegram application then search for [@BotFather](https://t.me/BotFather)
2. Click Start
3. Click Menu -> `/newbot` or type `/newbot` and hit Send
4. Copy your token and paste it in `config_variables.go`
    ```go
    const TELEGRAM_TOKEN = "1xxxxx2:AxxxxxxxxxxxxZ"
    ```
Attention: Never share this token. Everyone having this token can fully control your bot

##### Make a new channel

1. Make a new Telegram channel
2. Add your bot as admin to the channel
3. Add the bot [@IDBot](https://t.me/myidbot) as admin to the channel
4. type `/getgroupid` and hit Send
5. copy your channel/groupid and paste it in `config_variables.go`
    ```go
    const TELEGRAM_CHAT_ID = -100xxxxxxx9
    ```

### Running the bot

TODO