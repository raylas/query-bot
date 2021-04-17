# query-bot

query-bot is an unintelligent Slack bot for writing HTTP query responses to a given channel when user messages match pre-defined commands.

It's basically a Slack bot for curl.

Todo:
- [ ] Some basic tests
- [ ] Config boolean for help/queries menu
- [ ] Logic for handling methods other than `GET`

## Slack App
1. Create [a new Slack App](https://api.slack.com/apps) 
2. Enable **Socket Mode** and copy App Level Token
3. Under **OAuth & Permissions** set the following Bot Token Scopes, and copy Bot Token:
- `channels:history`
- `chat:write`
- `files:write`
- `users:read`
4. Toggle **Enable Events** under **Event Subscriptions**
5. Install the Slack App to your workspace
6. Invite the bot into channel(s)

## Configuration
`config.yaml`
```yaml
queries:
- command: !ifconfig
  url: https://ifconfig.co
  file: false
```

## Usage
### Deployment requirements
- Environment variables `SLACK_APP_TOKEN` and `SLACK_BOT_TOKEN` are set
- Configuration present or volume mounted at `/config.yaml` or `./config.yaml`

### Development
```shell
SLACK_APP_TOKEN=<app_token> \
SLACK_BOT_TOKEN=<bot_token> \
docker-compose run develop
```

### Build
```shell
docker-compose run build
```

### App
```shell
SLACK_APP_TOKEN=<app_token> \
SLACK_BOT_TOKEN=<bot_token> \
docker-compose run app
```
