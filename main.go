package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/stepconf"
)

type config struct {
	AddonURL       string `env:"addon_url"`
	AccessToken    string `env:"access_token"`
	BitriseAppSlug string `env:"BITRISE_APP_SLUG"`
}

func failf(msg string, args ...interface{}) {
	log.Errorf(msg, args...)
	os.Exit(1)
}

func main() {
	var cfg config
	if err := stepconf.Parse(&cfg); err != nil {
		failf("Issue with input: %s", err)
	}
	stepconf.Print(cfg)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/ascii-art/%s", cfg.AddonURL, cfg.BitriseAppSlug), nil)
	if err != nil {
		failf("Failed to create HTTP request: %s", err)
	}

	req.Header.Add("Authentication", cfg.AccessToken)

	var client http.Client
	resp, err := client.Do(req)
	if err != nil {
		failf("Failed to send HTTP request: %s", err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Warnf("Failed to close response body: %s", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		failf("Failed to fetch ASCII art from add-on: %s", err)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		failf("Failed to read response body: %s", err)
	}
	asciiArt := string(bodyBytes)

	cmdLog := command.New("bitrise", "envman", "add", "--key", "BITRISE_ASCII_ART", "--value", asciiArt, "-n")

	fmt.Println()
	log.Donef("$ %s", cmdLog.PrintableCommandArgs())
	fmt.Println()

	if err := cmdLog.Run(); err != nil {
		failf("Failed to expose output with envman, error: %#v | output: %s", err, cmdLog)
	}
}
