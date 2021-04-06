# query-bot

query-bot is an unintelligent Slack bot for writing HTTP query responses to a given channel when user messages match pre-defined commands.

It's basically a Slack bot for curl.

Todo:
- [ ] Slack integration updated to use [Slack apps](https://api.slack.com/start)
- [ ] Some basic tests
- [ ] Config boolean for help/queries menu
- [ ] Logic for handling methods other than `GET`

## Prerequisites
- A legacy Slack bot integraton
- Invite the bot into a channel

## Configure
`config.yaml`
```yaml
queries:
- command: !ifconfig
  url: https://ifconfig.co
  file: false
```

## Usage
### Deployment
Requirements:
- Configuration present, baked-in, or volume mounted at `/config.yaml` or `./config.yaml`
- Environment variable `SLACK_TOKEN` is set

### Development
```shell
SLACK_TOKEN=<token> docker-compose run develop
```

### Build
```shell
docker-compose run build
```

### App
```shell
SLACK_TOKEN=<token> docker-compose run app
```
