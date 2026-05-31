// Package log contains a log-notifier for debugging
package log //revive:disable-line:package-naming // it's a package logging the birthdays

import (
	"fmt"
	"time"

	"github.com/Luzifer/go_helpers/fieldcollection"
	"github.com/emersion/go-vcard"
	"github.com/sirupsen/logrus"

	"git.luzifer.io/luzifer/birthday-notifier/pkg/formatter"
	"git.luzifer.io/luzifer/birthday-notifier/pkg/notifier"
)

type (
	// Notifier implements the notifier interface
	Notifier struct{}
)

var _ notifier.Notifier = Notifier{}

// SendNotification implements the Notifier interface
func (Notifier) SendNotification(_ *fieldcollection.FieldCollection, contact vcard.Card, when time.Time) error {
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

// ValidateSettings implements the Notifier interface
func (Notifier) ValidateSettings(*fieldcollection.FieldCollection) error {
	// We don't take settings so everything is fine
	return nil
}
