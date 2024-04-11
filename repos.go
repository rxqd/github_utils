package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type Repository struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	URL         string `json:"html_url"`
	IsPrivate   bool   `json:"private"`
	IsFork      bool   `json:"fork"`
}

const (
	GET    = "GET"
	DELETE = "DELETE"
)

func (v *Repository) String() string {
	return fmt.Sprintf("[%s] - %s", v.FullName, v.Description)
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

	if err := json.NewEncoder(file).Encode(repos); err != nil {
		return err
	}

	return nil
}

func FetchRepositories() ([]Repository, error) {
	var repositories []Repository
	page := 1

	for {
		repos, err := doFetchRequest(fmt.Sprintf(FetchReposURL, config.UserName), page)
		if err != nil {
			return nil, err
		}

		if len(repos) == 0 {
			break
		}

		repositories = append(repositories, onlyForks(repos)...)
		page++
	}

	return repositories, nil
}

func onlyForks(repos []Repository) []Repository {
	var forks []Repository

	for _, repo := range repos {
		if !repo.IsFork {
			continue
		}

		forks = append(forks, repo)
	}

	return forks
}

func newClient() *http.Client {
	return &http.Client{Timeout: 10 * time.Second}
}

func newRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", RootURL, url), nil)
	if err != nil {
		return nil, fmt.Errorf("error on creating request: %w", err)
	}

	req.Header.Set("X-GitHub-Api-Version", APIVersion)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", config.UserName)
	req.Header.Set("Authorization", fmt.Sprintf("token %s", config.AccessToken))

	return req, nil
}

func doFetchRequest(url string, page int) ([]Repository, error) {
	client := newClient()
	req, err := newRequest(GET, url)
	if err != nil {
		return nil, fmt.Errorf("doFetchRequest: %w", err)
	}

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
	if err != nil {
		return nil, err
	}

	return repos, nil
}

func deleteRepositories(repos []Repository) error {
	_ = repos
	for _, repo := range repos {
		url := fmt.Sprintf(DeleteRepoURL, repo.FullName)
		err := doDeleteRequest(url)
		if err != nil {
			fmt.Printf("Remove %s error: %v", repo.FullName, err.Error())
		}
	}

	return nil
}

func listRepositories() ([]Repository, error) {
	file, err := os.Open(linksFile)
	if err != nil {
		return nil, fmt.Errorf("listRepositories: %w", err)
	}

	var result []Repository
	err = json.NewDecoder(file).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	return result, nil
}

func doDeleteRequest(url string) error {
	client := newClient()
	req, err := newRequest(DELETE, url)
	if err != nil {
		return fmt.Errorf("doDeleteRequest: %w", err)
	}

	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error on making request: %w", err)
	}
	defer response.Body.Close()

	return nil
}

func removeCmdAll(repos []Repository) {
	fmt.Println("The following repos will be deleted:")
	for _, repo := range repos {
		fmt.Printf("%s\n", repo.String())
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Are you sure you want to delete the following repos? y/n/q: ")
	input, _ := reader.ReadString('\n')
	if strings.TrimSpace(input) == "y" {
		err := deleteRepositories(repos)
		if err != nil {
			printError("On delete repositories", err)
			os.Exit(1)
		}

		fmt.Printf("Repositories deleted successfully")
	} else {
		fmt.Printf("Cancelled")
	}
}

func removeWithCheck(repos []Repository) {
	reader := bufio.NewReader(os.Stdin)
	var reposForDelete []Repository

outOfLoop:
	for _, repo := range repos {
		fmt.Printf("Delete repository '%s - %s'? (y/n/q/s): ",
			repo.FullName, repo.Description)

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch strings.ToLower(input) {
		case "y":
			reposForDelete = append(reposForDelete, repo)
			fmt.Println("Added repository for deletion. It will happen after you check ALL of them")
		case "n":
			fmt.Println("Skipping...")
		case "s":
			fmt.Println("Skipping all next repositories...")
			break outOfLoop
		case "q":
			return
		}

	}

	err := deleteRepositories(reposForDelete)
	if err != nil {
		printError("On delete repositories", err)
		os.Exit(1)
	}

	fmt.Printf("Repositories deleted successfully")
}
