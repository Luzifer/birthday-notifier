package main

import (
	"git.luzifer.io/luzifer/birthday-notifier/pkg/notifier"
	"git.luzifer.io/luzifer/birthday-notifier/pkg/notifier/log"
	"git.luzifer.io/luzifer/birthday-notifier/pkg/notifier/pushover"
)

func getNotifierByName(name string) notifier.Notifier {
	switch name {
	case "log":
		return log.Notifier{}

	case "pushover":
		return pushover.Notifier{}

	default:
		return nil
	}
}
