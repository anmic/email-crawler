package config

import (
	"io/ioutil"
	"log"

	"github.com/naoina/toml"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"gopkg.in/pg.v4"
)

type Config struct {
	Gmail *oauth2.Config

	Postgres *pg.Options
}

func Read() (*Config, error) {
	var config Config

	b, err := ioutil.ReadFile("config.toml")
	if err != nil {
		return nil, err
	}

	if err := toml.Unmarshal(b, &config); err != nil {
		return nil, err
	}

	b, err = ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
		return nil, err
	}

	gmailConfig, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
		return nil, err
	}

	config.Gmail = gmailConfig
	return &config, nil
}
