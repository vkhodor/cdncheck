package senders

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type RequestBody struct {
	Channel  string `json:"channel"`
	Username string `json:"username"`
	Text     string `json:"text"`
}

type Slack struct {
	URL      string
	Username string
	Channel  string
}

func NewSlack(url string, username string, channel string) *Slack {
	return &Slack{URL: url, Username: username, Channel: channel}
}

func (s *Slack) Send(msg string) error {
	slackBody, err := json.Marshal(RequestBody{Channel: s.Channel, Username: s.Username, Text: msg})
	req, err := http.NewRequest(http.MethodPost, s.URL, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return fmt.Errorf("error ReadFrom %v", err)
	}

	if buf.String() != "ok" {
		return fmt.Errorf("non-ok response returned from Slack")
	}

	return nil
}
