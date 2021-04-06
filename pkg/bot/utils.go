package bot

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/raylas/query-bot/pkg/config"
	"github.com/slack-go/slack"
)

func (s *Slack) Parse(ev *slack.MessageEvent, user *slack.User, config config.Configuration) error {
	msg := ev.Msg

	// Iterate over queries in configuration
	for _, q := range config.Queries {
		// Look for a query command that matches channel message
		if msg.Text == q.Command {
			resp, err := Query(q.URL, q.File)
			if err != nil {
				s.Logger.Printf("[ERROR] %s\n", err)
			}

			if q.File {
				attachment := slack.FileUploadParameters{
					File:     resp,
					Channels: []string{ev.Channel},
				}

				_, err = s.Client.UploadFile(attachment)
				if err != nil {
					s.Logger.Printf("[ERROR] %s\n", err)
					return err
				}
			} else {
				attachment := slack.Attachment{
					Pretext: "```\n" + resp + "```\n",
				}

				_, _, err = s.Client.PostMessage(
					ev.Channel,
					slack.MsgOptionAttachments(attachment),
					slack.MsgOptionAsUser(true),
				)
				if err != nil {
					s.Logger.Printf("[ERROR] %s\n", err)
					return err
				}
			}
		}
	}

	return nil
}

func Query(url string, file bool) (string, error) {
	c := &http.Client{
		Timeout: 15 * time.Second,
	}

	// Run query
	resp, err := c.Get(url)
	if err != nil {
		return "", err
	}

	if file {
		// Define local file path (scratch dir)
		path := "/tmp/" + path.Base(resp.Request.URL.Path)

		// Create file
		out, err := os.Create(path)
		if err != nil {
			return path, err
		}
		defer out.Close()

		defer resp.Body.Close()

		// Write data to file
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return path, err
		}

		// Return file path
		return path, err
	} else {
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return string(bytes), err
		}

		// Return stringified response bytes
		return string(bytes), err
	}
}
