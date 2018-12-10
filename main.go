package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	addonURL, ok := os.LookupEnv("addon_url")
	if !ok {
		fmt.Println("No addon URL is specified")
		os.Exit(1)
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/ascii-art/%s", addonURL, os.Getenv("BITRISE_APP_SLUG")), nil)
	if err != nil {
		fmt.Printf("Failed to create HTTP request: %s", err)
		os.Exit(1)
	}
	accessKey, ok := os.LookupEnv("access_token")
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

	cmdLog, err := exec.Command("bitrise", "envman", "add", "--key", "BITRISE_ASCII_ART", "--value", asciiArt, "-n").CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to expose output with envman, error: %#v | output: %s", err, cmdLog)
		os.Exit(1)
	}

	os.Exit(0)
}
