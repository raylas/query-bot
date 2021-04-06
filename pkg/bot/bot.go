package bot

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/raylas/query-bot/pkg/config"
	"github.com/slack-go/slack"
)

type Slack struct {
	Name  string
	Token string

	Logger *log.Logger

	Client *slack.Client
}

func New() (*Slack, error) {
	token := os.Getenv("SLACK_TOKEN")
	if len(token) == 0 {
		return nil, fmt.Errorf("SLACK_TOKEN not in environment") // Probably don't need to call fmt here
	}

	return &Slack{
		Name:   "",
		Token:  token,
		Logger: &log.Logger{},
		Client: slack.New(token),
	}, nil
}

func (s *Slack) Listen(ctx context.Context, config config.Configuration) error {
	slack.OptionLog(s.Logger)

	rtm := s.Client.NewRTM()
	go rtm.ManageConnection()

	// Loop to process channel messages
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			msg := ev.Msg

			if msg.SubType != "" {
				break // Normal messages only
			}

			user, err := s.Client.GetUserInfo(ev.User)
			if err != nil {
				s.Logger.Printf("[WARN] Could not get user information: %s", ev.User)
				continue
			}

			// s.Logger.Printf("[DEBUG] Heard message from `%s` (%s)\n", user.Profile.RealName, ev.User)
			err = s.Parse(ev, user, config)
			if err != nil {
				s.Logger.Printf("[ERROR] Could not call parser function")
			}

		case *slack.ConnectedEvent:
			s.Logger.Println("[INFO] Bot connected")

		case *slack.InvalidAuthEvent:
			s.Logger.Println("[ERROR] Invalid token")

		case *slack.RTMError:
			s.Logger.Printf("[ERROR] %s\n", ev.Error())
		}
	}

	return nil
}
