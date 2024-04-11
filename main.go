package main

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	APIVersion    = "2022-11-28"
	configFile    = "config.json"
	linksFile     = "links.json"
	RootURL       = "https://api.github.com"
	FetchReposURL = "/users/%s/repos"
	DeleteRepoURL = "/repos/%s"
)

type Config struct {
	AccessToken string `json:"access_token"`
	UserName    string `json:"github_username"`
}

var config Config

const usageMessage = `
Usage: github_utils <subcommand> [options]

Available Subcommands:
  fetch    Fetches repositories from GitHub and saves to file
  list     Lists repositories from file
  remove   Removes repositories (interactive with confirmation)
`

const removeCmdUsageMessage = `
Usage: github_utils remove <subcommand>

Available Subcommands:
  all      Removes all repositories with confirmation
  					[You can manually remove single repo from json file]
  check    Asks confirmation for each single repository
  					y 	remove
  					n 	skip
  					q 	quit WITHOUT ANY REMOVE
  					s 	Skip ALL NEXT
`

func main() {
	err := initConfig()
	if err != nil {
		printError("Error on reading config", err)
		os.Exit(1)
	}

	if len(os.Args) < 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Print(usageMessage)
		os.Exit(0)
	}

	switch os.Args[1] {
	case "fetch":
		fmt.Println("Fetching repositories from github...")
	case "list":
		repos, err := listRepositories()
		if err != nil {
			printError("Can't list repositories", err)
		}

		for _, repo := range repos {
			fmt.Printf("%s\n", repo.String())
		}
	case "remove":
		if len(os.Args) < 3 {
			fmt.Print(removeCmdUsageMessage)
			os.Exit(0)
		}

		repos, err := listRepositories()
		if err != nil {
			printError("Can't list repositories", err)
		}

		switch os.Args[2] {
		case "all":
			removeCmdAll(repos)
		case "check":
			removeWithCheck(repos)
		default:
			fmt.Print(removeCmdUsageMessage)
			os.Exit(0)
		}

	default:
		fmt.Print(usageMessage)
	}
}

func initConfig() error {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("please create \"%s\" from \"%s.tpl\"", configFile, configFile)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	return nil
}
