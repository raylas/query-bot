package bot

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/raylas/query-bot/pkg/config"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type Slack struct {
	Name   string
	Token  string
	Logger *log.Logger
	Client *slack.Client
}

func New() (*Slack, error) {
	appToken := os.Getenv("SLACK_APP_TOKEN")
	if len(appToken) == 0 {
		return nil, fmt.Errorf("SLACK_APP_TOKEN not in environment") // Probably don't need to call fmt here
	}
	if !strings.HasPrefix(appToken, "xapp-") {
		return nil, fmt.Errorf("SLACK_APP_TOKEN must have the prefix \"xapp-\"")
	}

	botToken := os.Getenv("SLACK_BOT_TOKEN")
	if len(botToken) == 0 {
		return nil, fmt.Errorf("SLACK_BOT_TOKEN not in environment") // Probably don't need to call fmt here
	}
	if !strings.HasPrefix(botToken, "xoxb-") {
		return nil, fmt.Errorf("SLACK_BOT_TOKEN must have the prefix \"xoxb-\"")
	}

	return &Slack{
		Name:   "",
		Token:  botToken,
		Logger: &log.Logger{},
		Client: slack.New(
			botToken,
			slack.OptionAppLevelToken(appToken),
		),
	}, nil
}

func (s *Slack) Listen(ctx context.Context, config config.Configuration) error {
	fmt.Println("Blop")

	api := socketmode.New(
		s.Client,
		socketmode.OptionDebug(false),
		socketmode.OptionLog(s.Logger),
	)

	go func() {
		// Loop to process channel messages
		for evt := range api.Events {
			switch evt.Type {
			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
				if !ok {
					s.Logger.Printf("[INFO] Ignored %+v\n", evt)

					continue
				}

				s.Logger.Printf("[DEBUG] Event received: %+v\n", eventsAPIEvent)

				api.Ack(*evt.Request)

				switch eventsAPIEvent.Type {
				case slackevents.CallbackEvent:
					innerEvent := eventsAPIEvent.InnerEvent
					switch ev := innerEvent.Data.(type) {
					case *slackevents.MessageEvent:
						if ev.SubType != "" {
							break // Normal messages only
						}

						user, err := s.Client.GetUserInfo(ev.User)
						if err != nil {
							s.Logger.Printf("[WARN] Could not get user information: %s", ev.User)
							continue
						}

						s.Logger.Printf("[DEBUG] Heard message from `%s` (%s)\n", user.Profile.RealName, ev.User)
						err = s.Parse(ev, user, config)
						if err != nil {
							s.Logger.Printf("[ERROR] Could not call parser function")
						}
					default:
						s.Logger.Printf("[INFO] Unsupported Events API event received")
					}
				}
			case socketmode.EventTypeConnecting:
				s.Logger.Println("[INFO] Connecting to Slack with socket mode")
			case socketmode.EventTypeConnectionError:
				s.Logger.Println("[INFO] Connection failed, retrying later")
			case socketmode.EventTypeConnected:
				s.Logger.Println("[INFO] Connected to Slack with socket mode")
			}
		}
	}()

	api.Run()

	return nil
}
