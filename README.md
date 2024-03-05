# Fiber Chat Bot (Telegram)

Fiber Chat Bot is a bot that queries a local fiber node and when detecting a new block broadcasts it via chat (currently only telegram).

The bot currently supports a local fiber node that is accessible via HTTP. \
HTTPS might be implemented in a later version

![image](/images/message_example1.png "Skycoin Blockhain Bot Message")


# Table of Contents

<!-- MarkdownTOC levels="1,2,3,4,5" autolink="true" bracket="round" -->

- [Installation](#installation)
    - [Software Requirements](#software-requirements)
    - [Cloning the Bot](#cloning-the-bot)
    - [Configuring the Bot](#configuring-the-bot)
        - [Telegram Configuration](#telegram-configuration)
            - [New Bot](#new-bot)
            - [New Channel](#new-channel)
    - [Installing the Bot](#installing-the-bot)
    - [Running the Bot](#running-the-bot)


## Installation

### Software Requirements

The bot depends on a local instance of a fiber node.\
See information about a fiber node on https://github.com/skycoin/skycoin.git.

!! You have to install a local node and then let it synchronize to the latest block to use the bot !!

### Cloning the Bot

1. Clone the bot in your go directory

    ```sh
    git clone https://github.com/4rchim3d3s/fiber-chat-bot.git $GOPATH/src/github.com/4rchim3d3s/fiber-chat-bot
    ```

### Configuring the Bot

1. Copy `default_config.json` to a `*_config.json` of your choice (e.g. fibercoin_config.json)
2. Edit the config depending to your preferences.


#### Telegram Configuration

##### New Bot

1. Open Telegram application then search for [@BotFather](https://t.me/BotFather)
2. Click Start
3. Click Menu -> `/newbot` or type `/newbot` and hit Send
4. Copy your token and paste it in `your_config.json`
    ```json
    "telegram" : {
        "bot_token" : "your-bot-token",
        "chat_id" : YourChatID
    },
    ```
!!! Attention: Never share this token. Everyone having this token can fully control your bot

##### New Channel

1. Make a new Telegram channel
2. Add your bot as admin to the channel
3. Add the bot [@IDBot](https://t.me/myidbot) as admin to the channel
4. type `/getgroupid` and hit Send
5. copy your channel/groupid and paste it in `your_config.json`
    ```json
    "telegram" : {
        "bot_token" : "your-bot-token",
        "chat_id" : YourChatID
    },
    ```

### Installing the Bot

1. Change directory to the cloned repo:\
    `cd $GOPATH/src/github.com/4rchim3d3s/fiber-chat-bot`

2. Install it with:\
    `go install`

### Running the Bot

1. Run the bot with:
    ```sh
    go run fiber-chat-bot -cfg your_config.json -block theBlockYouWantToStartToQueryAndSendFirst
    ```

    e.g. `go run fiber-chat-bot -cfg skycoin_config.json -block 186798`