package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"git.luzifer.io/luzifer/birthday-notifier/pkg/config"
	"git.luzifer.io/luzifer/birthday-notifier/pkg/dateutil"
	"git.luzifer.io/luzifer/birthday-notifier/pkg/formatter"
	"git.luzifer.io/luzifer/birthday-notifier/pkg/notifier"
	"github.com/emersion/go-vcard"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"

	"github.com/Luzifer/go_helpers/v2/fieldcollection"
	"github.com/Luzifer/rconfig/v2"
)

var (
	cfg = struct {
		Config         string `flag:"config,c" default:"config.yaml" description:"Configuration file path"`
		LogLevel       string `flag:"log-level" default:"info" description:"Log level (debug, info, warn, error, fatal)"`
		VersionAndExit bool   `flag:"version" default:"false" description:"Prints current version and exits"`
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

	configFile, err := config.LoadFromFile(cfg.Config)
	if err != nil {
		logrus.WithError(err).Fatal("loading configuration file")
	}

	if err = validateNotifierConfigs(configFile); err != nil {
		logrus.WithError(err).Fatal("validating configuration")
	}

	if err = formatter.SetTemplate(configFile.Template); err != nil {
		logrus.WithError(err).Fatal("setting template")
	}

	if birthdays, err = fetchBirthdays(configFile.Webdav); err != nil {
		logrus.WithError(err).Fatal("initially fetching birthdays")
	}

	crontab := cron.New()

	// Periodically update birthdays
	if _, err = crontab.AddFunc(
		fmt.Sprintf("@every %s", configFile.Webdav.FetchInterval),
		cronFetchBirthdays(configFile.Webdav),
	); err != nil {
		logrus.WithError(err).Fatal("adding update-cron")
	}

	// Send notifications at midnight
	if _, err = crontab.AddFunc("@midnight", cronSendNotifications(configFile)); err != nil {
		logrus.WithError(err).Fatal("adding update-cron")
	}

	logrus.WithFields(logrus.Fields{
		"advance": configFile.NotifyDaysInAdvance,
		"version": version,
	}).Info("birthday-notifier started")
	crontab.Start()

	for {
		select {}
	}
}

func cronFetchBirthdays(webdavConfig config.WebdavConfig) func() {
	return func() {
		birthdaysLock.Lock()
		defer birthdaysLock.Unlock()

		var err error
		if birthdays, err = fetchBirthdays(webdavConfig); err != nil {
			logrus.WithError(err).Error("updating birthdays")
		}
	}
}

func cronSendNotifications(configFile config.File) func() {
	return func() {
		birthdaysLock.Lock()
		defer birthdaysLock.Unlock()

		for _, b := range birthdays {
			for _, advanceDays := range append(configFile.NotifyDaysInAdvance, 0) {
				if !dateutil.IsToday(notifyDate(dateutil.ProjectToNextBirthday(b.birthday), advanceDays)) {
					continue
				}

				for i := range configFile.Notifiers {
					notifyInstance := getNotifierByName(configFile.Notifiers[i].Type)

					go func(
						n notifier.Notifier,
						settings *fieldcollection.FieldCollection,
						contact vcard.Card,
						when time.Time,
					) {
						if err := n.SendNotification(settings, contact, when); err != nil {
							logrus.
								WithError(err).
								WithField("name", contact.Get(vcard.FieldFormattedName).Value).
								Error("sending notification")
						}
					}(notifyInstance, configFile.Notifiers[i].Settings, b.contact, b.birthday)
				}
			}
		}
	}
}

func notifyDate(t time.Time, daysInAdvance int) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day()-daysInAdvance, 0, 0, 0, 0, time.Local)
}

func validateNotifierConfigs(configFile config.File) (err error) {
	for i := range configFile.Notifiers {
		notifierCfg := configFile.Notifiers[i]

		n := getNotifierByName(notifierCfg.Type)
		if n == nil {
			return fmt.Errorf("notifier %q does not exist", notifierCfg.Type)
		}

		if err = n.ValidateSettings(notifierCfg.Settings); err != nil {
			return fmt.Errorf("settings for %q are invalid: %w", notifierCfg.Type, err)
		}
	}

	return nil
}
