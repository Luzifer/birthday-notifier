// Package slack provides a notifier to send birthday notifications
// through a Slack(-compatible) WebHook
package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"git.luzifer.io/luzifer/birthday-notifier/pkg/formatter"
	"git.luzifer.io/luzifer/birthday-notifier/pkg/notifier"
	"github.com/emersion/go-vcard"
	"github.com/sirupsen/logrus"
)

const webhookPostTimeout = 2 * time.Second

type (
	// Notifier implements the notifier interface
	Notifier struct{}
)

var _ notifier.Notifier = Notifier{}

// SendNotification implements the Notifier interface
func (Notifier) SendNotification(contact vcard.Card, when time.Time) error {
	if contact.Name() == nil {
		return fmt.Errorf("contact has no name")
	}

	webhookURL := os.Getenv("SLACK_WEBHOOK")

	if webhookURL == "" {
		return fmt.Errorf("missing SLACK_WEBHOOK env variable")
	}

	text, err := formatter.FormatNotificationText(contact, when)
	if err != nil {
		return fmt.Errorf("rendering text: %w", err)
	}

	payload := new(bytes.Buffer)
	if err = json.NewEncoder(payload).Encode(struct {
		Channel   string `json:"channel,omitempty"`
		IconEmoji string `json:"icon_emoji,omitempty"`
		Text      string `json:"text"`
		Username  string `json:"username,omitempty"`
	}{
		Channel:   os.Getenv("SLACK_CHANNEL"),
		IconEmoji: os.Getenv("SLACK_ICON_EMOJI"),
		Text:      text,
		Username:  os.Getenv("SLACK_USERNAME"),
	}); err != nil {
		return fmt.Errorf("encoding hook payload: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), webhookPostTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhookURL, payload)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logrus.WithError(err).Error("closing slack response body (leaked fd)")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	return nil
}
