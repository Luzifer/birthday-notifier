// Package formatter contains a helper to format the date of a birthday
// into a notification text
package formatter

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"
	"time"

	"git.luzifer.io/luzifer/birthday-notifier/pkg/dateutil"
	"github.com/emersion/go-vcard"
)

const timeDay = 24 * time.Hour

var (
	// DefaultTemplate contains the template used in testing and as a
	// default in the config package
	DefaultTemplate = regexp.MustCompile(`\s+`).ReplaceAllString(strings.TrimSpace(strings.ReplaceAll(`
{{ .contact | getName }} has their birthday {{ if .when | isToday -}} today {{- else -}} on {{ (.when | projectToNext).Format "Mon, 02 Jan" }} {{- end }}.
{{ if gt .when.Year 1 -}}They are turning {{ .when | getAge }}.{{- end }}
`, "\n", " ")), " ")

	notifyTpl *template.Template
)

// FormatNotificationText takes the notification template and renders
// the contact / birthday date into a text to submit in the notification
func FormatNotificationText(contact vcard.Card, when time.Time) (text string, err error) {
	buf := new(bytes.Buffer)

	if err = notifyTpl.Execute(buf, map[string]any{
		"contact": contact,
		"when":    when,
	}); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return buf.String(), nil
}

// FormatNotificationTitle provides a title from the contacts formatted
// name or from given and family name
func FormatNotificationTitle(contact vcard.Card) (title string) {
	for _, fn := range contact.FormattedNames() {
		if fn.Value != "" {
			title = fmt.Sprintf("%s (Birthday)", fn.Value)
		}
	}

	if title == "" {
		title = fmt.Sprintf("%s %s (Birthday)", contact.Name().GivenName, contact.Name().FamilyName)
	}

	return title
}

// SetTemplate initializes the template to use in the
// FormatNotificationText function. This MUST be called before first
// use of the FormatNotificationText function.
func SetTemplate(rawTpl string) error {
	var err error
	notifyTpl, err = template.New("notification").Funcs(template.FuncMap{
		"getAge":        getAge,
		"getName":       getContactName,
		"isToday":       dateutil.IsToday,
		"projectToNext": dateutil.ProjectToNextBirthday,
	}).Parse(rawTpl)
	if err != nil {
		return fmt.Errorf("parsing notification template: %w", err)
	}

	return nil
}

func getAge(t time.Time) int {
	return dateutil.ProjectToNextBirthday(t).Year() - t.Year()
}

func getContactName(contact vcard.Card) string {
	if contact.Name() != nil && contact.Name().GivenName != "" {
		return contact.Name().GivenName
	}

	return contact.FormattedNames()[0].Value
}
