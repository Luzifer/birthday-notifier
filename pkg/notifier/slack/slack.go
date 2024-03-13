// Package slack provides a notifier to send birthday notifications
// through a Slack(-compatible) WebHook
package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"git.luzifer.io/luzifer/birthday-notifier/pkg/formatter"
	"git.luzifer.io/luzifer/birthday-notifier/pkg/notifier"
	"github.com/Luzifer/go_helpers/v2/fieldcollection"
	"github.com/emersion/go-vcard"
	"github.com/sirupsen/logrus"
)

const webhookPostTimeout = 2 * time.Second

type (
	// Notifier implements the notifier interface
	Notifier struct{}
)

var (
	ptrStrEmpty = func(v string) *string { return &v }("")

	_ notifier.Notifier = Notifier{}
)

// SendNotification implements the Notifier interface
func (Notifier) SendNotification(settings *fieldcollection.FieldCollection, contact vcard.Card, when time.Time) error {
	if contact.Name() == nil {
		return fmt.Errorf("contact has no name")
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
		Channel:   settings.MustString("channel", ptrStrEmpty),
		IconEmoji: settings.MustString("iconEmoji", ptrStrEmpty),
		Text:      text,
		Username:  settings.MustString("username", ptrStrEmpty),
	}); err != nil {
		return fmt.Errorf("encoding hook payload: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), webhookPostTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, settings.MustString("webhook", nil), payload)
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

// ValidateSettings implements the Notifier interface
func (Notifier) ValidateSettings(settings *fieldcollection.FieldCollection) (err error) {
	if v, err := settings.String("webhook"); err != nil || v == "" {
		return fmt.Errorf("webhook is expected to be non-empty string")
	}

	return nil
}
