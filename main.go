package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	addonURL := "https://bitrise-sample-addon.herokuapp.com/"
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/ascii-art/%s", addonURL, os.Getenv("BITRISE_APP_SLUG")), nil)
	if err != nil {
		fmt.Printf("Failed to create HTTP request: %s", err)
		os.Exit(1)
	}
	accessKey, ok := os.LookupEnv("BITRISE_SAMPLE_ADDON_ACCESS_KEY")
	if !ok {
		fmt.Println("Failed to look up environment variable")
		os.Exit(1)
	}
	req.Header.Add("Authentication", accessKey)

	var client http.Client
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to send HTTP request: %s", err)
		os.Exit(1)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to fetch ASCII art from add-on: %s", err)
		os.Exit(1)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %s", err)
		os.Exit(1)
	}
	asciiArt := string(bodyBytes)

	//
	// --- Step Outputs: Export Environment Variables for other Steps:
	// You can export Environment Variables for other Steps with
	//  envman, which is automatically installed by `bitrise setup`.
	// A very simple example:
	cmdLog, err := exec.Command("bitrise", "envman", "add", "--key", "BITRISE_ASCII_ART", "--value", asciiArt, "-n").CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to expose output with envman, error: %#v | output: %s", err, cmdLog)
		os.Exit(1)
	}
	// You can find more usage examples on envman's GitHub page
	//  at: https://github.com/bitrise-io/envman

	//
	// --- Exit codes:
	// The exit code of your Step is very important. If you return
	//  with a 0 exit code `bitrise` will register your Step as "successful".
	// Any non zero exit code will be registered as "failed" by `bitrise`.
	os.Exit(0)
}
