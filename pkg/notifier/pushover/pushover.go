// Package pushover provides a notifier to send birthday notifications
// using Pushover.net
package pushover

import (
	"fmt"
	"os"
	"time"

	"git.luzifer.io/luzifer/birthday-notifier/pkg/formatter"
	"git.luzifer.io/luzifer/birthday-notifier/pkg/notifier"
	"github.com/emersion/go-vcard"
	"github.com/gregdel/pushover"
)

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

	var (
		apiToken = os.Getenv("PUSHOVER_API_TOKEN")
		userKey  = os.Getenv("PUSHOVER_USER_KEY")
	)

	if apiToken == "" {
		return fmt.Errorf("missing PUSHOVER_API_TOKEN env variable")
	}
	if userKey == "" {
		return fmt.Errorf("missing PUSHOVER_USER_KEY env variable")
	}

	text, err := formatter.FormatNotificationText(contact, when)
	if err != nil {
		return fmt.Errorf("rendering text: %w", err)
	}

	var title string
	for _, fn := range contact.FormattedNames() {
		if fn.Value != "" {
			title = fmt.Sprintf("%s (Birthday)", fn.Value)
		}
	}

	if title == "" {
		title = fmt.Sprintf("%s %s (Birthday)", contact.Name().GivenName, contact.Name().FamilyName)
	}

	message := &pushover.Message{
		Message: text,
		Title:   title,
		Sound:   os.Getenv("PUSHOVER_SOUND"),
	}

	if _, err = pushover.New(apiToken).
		SendMessage(message, pushover.NewRecipient(userKey)); err != nil {
		return fmt.Errorf("sending notification: %w", err)
	}

	return nil
}
