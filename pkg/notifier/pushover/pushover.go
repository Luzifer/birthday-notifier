// Package pushover provides a notifier to send birthday notifications
// using Pushover.net
package pushover

import (
	"fmt"
	"time"

	"git.luzifer.io/luzifer/birthday-notifier/pkg/formatter"
	"git.luzifer.io/luzifer/birthday-notifier/pkg/notifier"
	"github.com/Luzifer/go_helpers/v2/fieldcollection"
	"github.com/emersion/go-vcard"
	"github.com/gregdel/pushover"
)

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

	message := &pushover.Message{
		Message: text,
		Title:   formatter.FormatNotificationTitle(contact),
		Sound:   settings.MustString("sound", ptrStrEmpty),
	}

	if _, err = pushover.New(settings.MustString("apiToken", nil)).
		SendMessage(message, pushover.NewRecipient(settings.MustString("userKey", nil))); err != nil {
		return fmt.Errorf("sending notification: %w", err)
	}

	return nil
}

// ValidateSettings implements the Notifier interface
func (Notifier) ValidateSettings(settings *fieldcollection.FieldCollection) (err error) {
	if v, err := settings.String("apiToken"); err != nil || v == "" {
		return fmt.Errorf("apiToken is expected to be non-empty string")
	}

	if v, err := settings.String("userKey"); err != nil || v == "" {
		return fmt.Errorf("userKey is expected to be non-empty string")
	}

	return nil
}
