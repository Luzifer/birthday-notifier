// Package notifier includes the interface to implement in a notifier
package notifier

import (
	"time"

	"github.com/emersion/go-vcard"
)

type (
	// Notifier specifies what a Notifier can do
	Notifier interface {
		// SendNotification will be called with the contact and the
		// time when the birthday actually is. The method is therefore
		// also called when a notification in advance is configured and
		// needs to properly format the notification for that.
		SendNotification(contact vcard.Card, when time.Time) error
	}
)
