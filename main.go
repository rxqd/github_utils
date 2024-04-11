package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	APIVersion = "2022-11-28"
	configFile = "config.json"
	linksFile  = "links.json"
	URL        = "https://api.github.com/users/%s/repos"
)

type Config struct {
	URL         string
	AccessToken string `json:"access_token"`
	UserName    string `json:"github_username"`
}

type Repository struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	URL         string `json:"html_url"`
	IsPrivate   bool   `json:"private"`
	IsFork      bool   `json:"fork"`
}

var config Config

func main() {
	err := InitConfig()
	if err != nil {
		printError("Error in reading config", err)
		os.Exit(1)
	}

	fmt.Println("Choose your action")
}

func FetchAndSaveRepositories() {
	repos, err := FetchRepositories()
	if err != nil {
		printError("Error on fetching repos", err)
		os.Exit(1)
	}

	fmt.Println(format(GREEN, "Repositories are fetched successfully"))

	err = SaveRepositories(repos)
	if err != nil {
		printError("Error on saving repos", err)
		os.Exit(1)
	}

	fmt.Println(format(GREEN, "Repositories are saved successfully"))
}

func SaveRepositories(repos []Repository) error {
	file, err := os.Create(linksFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(repos)
	if err != nil {
		return err
	}

	return nil
}

func FetchRepositories() ([]Repository, error) {
	var repositories []Repository
	page := 1

	for {
		repos, err := doRequest(config.URL, page)
		if err != nil {
			return nil, err
		}

		if len(repos) == 0 {
			break
		}

		repositories = append(repositories, repos...)
		page++
	}

	return repositories, nil
}

func doRequest(url string, page int) ([]Repository, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error on creating request: %w", err)
	}

	req.Header.Set("X-GitHub-Api-Version", APIVersion)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", config.UserName)
	req.Header.Set("Authorization", fmt.Sprintf("token %s", config.AccessToken))

	q := req.URL.Query()
	q.Add("page", fmt.Sprintf("%d", page))
	req.URL.RawQuery = q.Encode()

	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error on making request: %w", err)
	}
	defer response.Body.Close()

	var repos []Repository
	err = json.NewDecoder(response.Body).Decode(&repos)

	return repos, nil
}

func InitConfig() error {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("please create \"%s\" from \"%s.tpl\"", configFile, configFile)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	config.URL = fmt.Sprintf(URL, config.UserName)

	return nil
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
