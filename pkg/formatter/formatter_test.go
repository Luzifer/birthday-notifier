package formatter

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"git.luzifer.io/luzifer/birthday-notifier/pkg/dateutil"
	"github.com/emersion/go-vcard"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestVCard(t *testing.T, content string) vcard.Card {
	c, err := vcard.NewDecoder(strings.NewReader(content)).Decode()
	require.NoError(t, err)

	return c
}

func TestFormatNotificationText(t *testing.T) {
	card := getTestVCard(t, `BEGIN:VCARD
VERSION:4.0
N:Bloggs;Joe;;;
FN:Joe Bloggs
EMAIL;TYPE=home;PREF=1:me@joebloggs.com
TEL;TYPE="cell,home";PREF=1:tel:+44 20 1234 5678
ADR;TYPE=home;PREF=1:;;1 Trafalgar Square;London;;WC2N;United Kingdom
URL;TYPE=home;PREF=1:http://joebloggs.com
IMPP;TYPE=home;PREF=1:skype:joe.bloggs
X-SOCIALPROFILE;TYPE=home;PREF=1:twitter:https://twitter.com/joebloggs
END:VCARD`)

	bday := time.Date(time.Now().Year()-30, time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)

	txt, err := FormatNotificationText(card, bday)
	require.NoError(t, err)
	assert.Equal(t, "Joe has their birthday today. They are turning 30.", txt)

	bday = bday.Add(timeDay)
	txt, err = FormatNotificationText(card, bday)
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf(
		"Joe has their birthday on %s. They are turning 30.",
		time.Now().Add(timeDay).Format("Mon, 02 Jan"),
	), txt)

	bday = bday.Add(-2 * timeDay)
	txt, err = FormatNotificationText(card, bday)
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf(
		"Joe has their birthday on %s. They are turning 31.",
		dateutil.ProjectToNextBirthday(time.Now().Add(-timeDay)).Format("Mon, 02 Jan"),
	), txt)
}
