# discord-log-bot
Discord Log Bot


## Usage
```docker
docker run \
  --name='discord-log-bot' \
  -e 'LOGBOT_TOKEN'='### \
  -e 'LOGBOT_CHANNEL'='###' \
  ghcr.io/bym0/discord-log-bot:latest
```

## Variables

- LOGBOT_TOKEN = Your Discord Bot Token.
- LOGBOT_CHANNEL = The Channel the bot should send messages to.
