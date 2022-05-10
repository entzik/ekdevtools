package gitlab

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const apiBase = "/api/v4"
const searchApiBase = apiBase + "/search"

type GitlabClient interface {
	SearchProjects(term string) ([]RepositoryDescriptor, error)
}

// gitlabClient a structure that contains information required to invoke gitlab API, such the API URL and
// the authentication token
type gitlabClient struct {
	// the URL of the gitlab API
	Url string
	// the authentication token to be used with the gitlab server
	Token string
}

// NewGitlabApi creates a new gitlab api object from environment variables
func NewGitlabApi() GitlabClient {
	return &gitlabClient{
		Url:   os.Getenv("GITLAB_API_BASE_URL"),
		Token: os.Getenv("GITLAB_API_TOKEN"),
	}
}

// SearchProjects search gitlab for project that match the specified term
func (gitlabClient *gitlabClient) SearchProjects(term string) ([]RepositoryDescriptor, error) {
	var pageSize = 100
	var parallelism = 4
	var page = 1

	var repositories []RepositoryDescriptor
	var errors []error
	var stay = true

	// create a channel to collect result pages
	resultsChannel := make(chan []RepositoryDescriptor)
	// create a channel to collect errors
	errorsChannel := make(chan error)
	for stay {
		for crt := 0; crt < parallelism; crt++ {
			// fetch a page concurrently
			go gitlabClient.searchProjectsPage(resultsChannel, errorsChannel, term, page+crt, pageSize)
		}
		// collect results
		for crt := 0; crt < parallelism; crt++ {
			// wait for responses to be delivered, concat items into the repositories slice
			repositories = append(repositories, <-resultsChannel...)
		}
		// collect errors
		for crt := 0; crt < parallelism; crt++ {
			// wait for errors to be delivered, concat items into the repositories slice
			err := <-errorsChannel
			errors = append(errors, err)
		}

		page += parallelism
		// if not all pages are full, we're done here
		stay = len(repositories) == pageSize*parallelism
	}

	if len(errors) > 0 {
		// if errors, return the first one that occurred
		return nil, errors[0]
	} else {
		// return results
		return repositories, nil
	}
}

func (gitlabClient *gitlabClient) searchProjectsPage(resultsChannel chan []RepositoryDescriptor, errorChannel chan error, term string, page int, pageSize int) {
	var repositories []RepositoryDescriptor

	var urlString = fmt.Sprintf("%s%s?per_page=%d&page=%d&scope=projects&search=%s", gitlabClient.Url, searchApiBase, pageSize, page, term)
	var client = &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", urlString, nil)
	req.Header.Add("PRIVATE-TOKEN", gitlabClient.Token)
	resp, err := client.Do(req)
	if err != nil {
		println("error -> ", err)
		resultsChannel <- repositories
		errorChannel <- err
	} else {
		defer resp.Body.Close()
		err = json.NewDecoder(resp.Body).Decode(&repositories)
		resultsChannel <- repositories
		errorChannel <- nil
	}
}
