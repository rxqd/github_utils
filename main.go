package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	URL         string `json:"api_url"`
	AccessToken string `json:"access_token"`
	UserName    string `json:"github_username"`
}

func main() {
	config, err := NewConfig()
	if err != nil {
		printError("Error in reading config", err)
		os.Exit(1)
	}

	fmt.Println(format(GREEN, config))
}

func NewConfig() (*Config, error) {
	data, err := os.ReadFile("config.json")
	if err != nil {
		return nil, fmt.Errorf("please create \"config.json\" from \"config.json.tpl\"")
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func printError(text string, err error) {
	fmt.Fprintf(os.Stderr, format(RED, "%s: %s\n"), text, err)
}

const (
	NONE = iota
	RED
	GREEN
	YELLOW
	BLUE
	PURPLE
)

func format(c int, text any) string {
	const escape = "\x1b"

	color := func(c int) string {
		var term string
		if c != NONE {
			term = "3"
		}

		return fmt.Sprintf("%s[%s%dm", escape, term, c)
	}

	return fmt.Sprintf("%s%s%s", color(c), text, color(NONE))
}
