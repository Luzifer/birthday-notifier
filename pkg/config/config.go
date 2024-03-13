// Package config contains parser and structure for the configuration
// of the tool
package config

import (
	"fmt"
	"io"
	"os"
	"time"

	"git.luzifer.io/luzifer/birthday-notifier/pkg/formatter"
	"github.com/Luzifer/go_helpers/v2/fieldcollection"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// WebdavPrincipalNextcloud is the principal default used for config
// files on parse and represents the principal format used by Nextcloud
const WebdavPrincipalNextcloud = "principals/users/%s"

type (
	// File contains the structure of the YAML configuration file
	File struct {
		NotifyDaysInAdvance []int            `yaml:"notifyDaysInAdvance"`
		Notifiers           []NotifierConfig `yaml:"notifiers"`

		Template string `yaml:"template"`

		Webdav WebdavConfig `yaml:"webdav"`
	}

	// NotifierConfig contains the type of the notifier and the settings
	// for it required to execute
	NotifierConfig struct {
		Type     string                           `yaml:"type"`
		Settings *fieldcollection.FieldCollection `yaml:"settings"`
	}

	// WebdavConfig defines how to interact with the Webdav server
	WebdavConfig struct {
		BaseURL       string        `yaml:"baseURL"`
		FetchInterval time.Duration `yaml:"fetchInterval"`
		Pass          string        `yaml:"pass"`
		Principal     string        `yaml:"principal"`
		User          string        `yaml:"user"`
	}
)

// Load parses the given reader over a default configuration replacing
// the fields specified in the reader
func Load(r io.Reader) (f File, err error) {
	f = defaultConfig()
	dec := yaml.NewDecoder(r)

	dec.KnownFields(true)
	if err = dec.Decode(&f); err != nil {
		return f, fmt.Errorf("decoding yaml: %w", err)
	}

	return f, nil
}

// LoadFromFile is a convenience wrapper around Load to read the config
// from file system
func LoadFromFile(filePath string) (f File, err error) {
	inFile, err := os.Open(filePath) //#nosec G304 -- Intended to load a given path
	if err != nil {
		return f, fmt.Errorf("opening file: %w", err)
	}
	defer func() {
		if err := inFile.Close(); err != nil {
			logrus.WithError(err).Error("closing config file (leaked fd)")
		}
	}()

	return Load(inFile)
}

func defaultConfig() File {
	return File{
		NotifyDaysInAdvance: nil,

		Template: formatter.DefaultTemplate,

		Webdav: WebdavConfig{
			FetchInterval: time.Hour,
			Principal:     WebdavPrincipalNextcloud,
		},
	}
}
