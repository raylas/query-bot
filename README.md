# query-bot

query-bot is an unintelligent Slack bot for writing HTTP query responses to a given channel when user messages match pre-defined commands.

It's basically a Slack bot for curl.

Todo:
- [ ] Some basic tests
- [ ] Config boolean for help/queries menu
- [ ] Logic for handling methods other than `GET`

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
ENV=<slack_token> docker-compose run develop
```

### Build
```shell
docker-compose run build
```

### App
```shell
ENV=<slack_token> docker-compose run app
```
