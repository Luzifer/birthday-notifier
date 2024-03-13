package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"git.luzifer.io/luzifer/birthday-notifier/pkg/dateutil"
	"git.luzifer.io/luzifer/birthday-notifier/pkg/notifier"
	"github.com/emersion/go-vcard"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"

	"github.com/Luzifer/rconfig/v2"
)

var (
	cfg = struct {
		FetchInterval       time.Duration `flag:"fetch-interval" default:"1h" description:"How often to fetch birthdays from CardDAV"`
		LogLevel            string        `flag:"log-level" default:"info" description:"Log level (debug, info, warn, error, fatal)"`
		NotifyDaysInAdvance []int         `flag:"notify-days-in-advance" default:"1" description:"Send notification X days before birthday"`
		NotifyVia           []string      `flag:"notify-via" default:"log" description:"How to send the notification (one of: log, pushover)"`
		WebdavBaseURL       string        `flag:"webdav-base-url" default:"" description:"Webdav server to connect to"`
		WebdavPass          string        `flag:"webdav-pass" default:"" description:"Password for the Webdav user"`
		WebdavPrincipal     string        `flag:"webdav-principal" default:"principals/users/%s" description:"Principal format to fetch the addressbooks for (%s will be replaced with the webdav-user)"`
		WebdavUser          string        `flag:"webdav-user" default:"" description:"Username for Webdav login"`
		VersionAndExit      bool          `flag:"version" default:"false" description:"Prints current version and exits"`
	}{}

	birthdays     []birthdayEntry
	birthdaysLock sync.Mutex

	version = "dev"
)

func initApp() error {
	rconfig.AutoEnv(true)
	if err := rconfig.ParseAndValidate(&cfg); err != nil {
		return errors.Wrap(err, "parsing cli options")
	}

	l, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		return errors.Wrap(err, "parsing log-level")
	}
	logrus.SetLevel(l)

	return nil
}

func main() {
	var err error
	if err = initApp(); err != nil {
		logrus.WithError(err).Fatal("initializing app")
	}

	if cfg.VersionAndExit {
		logrus.WithField("version", version).Info("birthday-notifier")
		os.Exit(0)
	}

	var notifiers []notifier.Notifier
	for _, nv := range cfg.NotifyVia {
		notify := getNotifierByName(nv)
		if notify == nil {
			logrus.Fatal("unknown notifier specified")
		}
		notifiers = append(notifiers, notify)
	}

	if birthdays, err = fetchBirthdays(); err != nil {
		logrus.WithError(err).Fatal("initially fetching birthdays")
	}

	crontab := cron.New()

	// Periodically update birthdays
	if _, err = crontab.AddFunc(fmt.Sprintf("@every %s", cfg.FetchInterval), cronFetchBirthdays); err != nil {
		logrus.WithError(err).Fatal("adding update-cron")
	}

	// Send notifications at midnight
	if _, err = crontab.AddFunc("@midnight", cronSendNotifications(notifiers)); err != nil {
		logrus.WithError(err).Fatal("adding update-cron")
	}

	logrus.WithFields(logrus.Fields{
		"advance": cfg.NotifyDaysInAdvance,
		"version": version,
	}).Info("birthday-notifier started")
	crontab.Start()

	for {
		select {}
	}
}

func cronFetchBirthdays() {
	birthdaysLock.Lock()
	defer birthdaysLock.Unlock()

	var err error
	if birthdays, err = fetchBirthdays(); err != nil {
		logrus.WithError(err).Error("updating birthdays")
	}
}

func cronSendNotifications(notifiers []notifier.Notifier) func() {
	return func() {
		birthdaysLock.Lock()
		defer birthdaysLock.Unlock()

		for _, b := range birthdays {
			for _, advanceDays := range append(cfg.NotifyDaysInAdvance, 0) {
				if !dateutil.IsToday(notifyDate(dateutil.ProjectToNextBirthday(b.birthday), advanceDays)) {
					continue
				}

				for i := range notifiers {
					go func(n notifier.Notifier, contact vcard.Card, when time.Time) {
						if err := n.SendNotification(contact, when); err != nil {
							logrus.
								WithError(err).
								WithField("name", contact.Get(vcard.FieldFormattedName).Value).
								Error("sending notification")
						}
					}(notifiers[i], b.contact, b.birthday)
				}
			}
		}
	}
}

func notifyDate(t time.Time, daysInAdvance int) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day()-daysInAdvance, 0, 0, 0, 0, time.Local)
}
