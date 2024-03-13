package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"git.luzifer.io/luzifer/birthday-notifier/pkg/config"
	"git.luzifer.io/luzifer/birthday-notifier/pkg/dateutil"
	"github.com/emersion/go-vcard"
	"github.com/emersion/go-webdav"
	"github.com/emersion/go-webdav/carddav"
	"github.com/sirupsen/logrus"
)

type (
	birthdayEntry struct {
		contact  vcard.Card
		birthday time.Time
	}
)

func fetchBirthdays(webdavConfig config.WebdavConfig) (birthdays []birthdayEntry, err error) {
	client, err := carddav.NewClient(
		webdav.HTTPClientWithBasicAuth(http.DefaultClient, webdavConfig.User, webdavConfig.Pass),
		webdavConfig.BaseURL,
	)
	if err != nil {
		return nil, fmt.Errorf("creating carddav client: %w", err)
	}

	homeSet, err := client.FindAddressBookHomeSet(
		context.Background(),
		fmt.Sprintf(webdavConfig.Principal, webdavConfig.User),
	)
	if err != nil {
		return nil, fmt.Errorf("getting addressbook-home-set: %w", err)
	}

	books, err := client.FindAddressBooks(context.Background(), homeSet)
	if err != nil {
		return nil, fmt.Errorf("getting addressbooks: %w", err)
	}

	for _, book := range books {
		contacts, err := client.QueryAddressBook(context.Background(), book.Path, &carddav.AddressBookQuery{})
		if err != nil {
			return nil, fmt.Errorf("getting contacts: %w", err)
		}

		for _, address := range contacts {
			bday := address.Card.Get(vcard.FieldBirthday)
			if bday == nil {
				continue
			}

			bdayDate, err := dateutil.Parse(bday)
			if err != nil {
				logrus.WithField("date", bday).WithError(err).Error("parsing birthday")
				continue
			}

			birthdays = append(birthdays, birthdayEntry{
				contact:  address.Card,
				birthday: bdayDate,
			})
		}
	}

	logrus.Infof("fetched %d birthdays from contacts", len(birthdays))
	return birthdays, nil
}
