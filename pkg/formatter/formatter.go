// Package formatter contains a helper to format the date of a birthday
// into a notification text
package formatter

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"

	"git.luzifer.io/luzifer/birthday-notifier/pkg/dateutil"
	"github.com/emersion/go-vcard"
)

const timeDay = 24 * time.Hour

var (
	defaultTemplate = regexp.MustCompile(`\s+`).ReplaceAllString(strings.TrimSpace(strings.ReplaceAll(`
{{ .contact | getName }} has their birthday {{ if .when | isToday -}} today {{- else -}} on {{ (.when | projectToNext).Format "Mon, 02 Jan" }} {{- end }}.
{{ if gt .when.Year 1 -}}They are turning {{ .when | getAge }}.{{- end }}
`, "\n", " ")), " ")

	notifyTpl *template.Template
)

func init() {
	rawTpl := defaultTemplate
	if tpl := os.Getenv("NOTIFICATION_TEMPLATE"); tpl != "" {
		rawTpl = tpl
	}

	var err error
	notifyTpl, err = template.New("notification").Funcs(template.FuncMap{
		"getAge":        getAge,
		"getName":       getContactName,
		"isToday":       dateutil.IsToday,
		"projectToNext": dateutil.ProjectToNextBirthday,
	}).Parse(rawTpl)
	if err != nil {
		panic(fmt.Errorf("parsing notification template: %w", err))
	}
}

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

func getAge(t time.Time) int {
	return dateutil.ProjectToNextBirthday(t).Year() - t.Year()
}

func getContactName(contact vcard.Card) string {
	if contact.Name() != nil && contact.Name().GivenName != "" {
		return contact.Name().GivenName
	}

	return contact.FormattedNames()[0].Value
}
