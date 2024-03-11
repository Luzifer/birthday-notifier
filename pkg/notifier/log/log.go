// Package log contains a log-notifier for debugging
package log

import (
	"fmt"
	"time"

	"git.luzifer.io/luzifer/birthday-notifier/pkg/formatter"
	"git.luzifer.io/luzifer/birthday-notifier/pkg/notifier"
	"github.com/emersion/go-vcard"
	"github.com/sirupsen/logrus"
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

	text, err := formatter.FormatNotificationText(contact, when)
	if err != nil {
		return fmt.Errorf("rendering text: %w", err)
	}

	logrus.WithField("name", contact.Name().GivenName).Info(text)
	return nil
}
