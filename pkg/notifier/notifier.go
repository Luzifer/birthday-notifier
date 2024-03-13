// Package notifier includes the interface to implement in a notifier
package notifier

import (
	"time"

	"github.com/Luzifer/go_helpers/v2/fieldcollection"
	"github.com/emersion/go-vcard"
)

type (
	// Notifier specifies what a Notifier can do
	Notifier interface {
		// SendNotification will be called with the contact and the
		// time when the birthday actually is. The method is therefore
		// also called when a notification in advance is configured and
		// needs to properly format the notification for that. The settings
		// passed through this call MUST NOT be stored.
		SendNotification(settings *fieldcollection.FieldCollection, contact vcard.Card, when time.Time) error

		// ValidateSettings is called after configuration load to validate
		// the settings are suitable for the notifier and do not yield
		// surprising errors when doing the real notifications
		ValidateSettings(settings *fieldcollection.FieldCollection) error
	}
)
