package git

import (
	"github.com/go-git/go-git/v5"
	http2 "github.com/go-git/go-git/v5/plumbing/transport/http"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type GitRepositoryClone interface {
	Clone() (GitRepositoryClone, error)
	FindFilesBySuffix(suffix string) []string
	WaitForRemoteBranchAndCheckOut(jira string) error
}

type gitRepositoryClone struct {
	// url the URL
	url string
	// clonePath the path to the clone
	clonePath string
	// clone the gir repository structure
	clone *git.Repository
}

func NewRepositoryClone(url string) GitRepositoryClone {
	return &gitRepositoryClone{
		url:       url,
		clonePath: "",
	}
}

func (gitRepoClone gitRepositoryClone) Clone() (GitRepositoryClone, error) {
	dir, dirCreationErr := ioutil.TempDir("", "clone-example")
	if dirCreationErr != nil {
		log.Fatal(dirCreationErr)
	}
	if dirCreationErr != nil {
		println("got an error: " + dirCreationErr.Error())
		return nil, dirCreationErr
	}
	// defer os.RemoveAll(dir)

	clone, dirCreationErr := git.PlainClone(dir, false, &git.CloneOptions{
		URL: gitRepoClone.url,
		Auth: &http2.BasicAuth{
			Username: "something",
			Password: os.Getenv("GITLAB_API_TOKEN"),
		},
	})
	if dirCreationErr != nil {
		println("got an error: " + dirCreationErr.Error())
		return nil, dirCreationErr
	} else {
		return &gitRepositoryClone{
			clonePath: dir,
			url:       gitRepoClone.url,
			clone:     clone,
		}, nil
	}
}

func (gitRepoClone gitRepositoryClone) WaitForRemoteBranchAndCheckOut(jira string) error {
	// find branch for created for ticket
	remote, err := gitRepoClone.clone.Remote("origin")
	if err != nil {
		return err
	}
	var rounds = 15
	for rounds > 0 {
		refList, err := remote.List(&git.ListOptions{})
		if err != nil {
			return err
		}
		refPrefix := "refs/heads/"
		for _, ref := range refList {
			refName := ref.Name().String()
			if !strings.HasPrefix(refName, refPrefix) {
				continue
			}
			branchName := refName[len(refPrefix):]
			if strings.Contains(branchName, jira) { // we contain the
				rounds = 0

			}
		}
	}

	// code to check out a branch
	/*	err := r.Fetch(&git.FetchOptions{
			RefSpecs: []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
		})
		if err != nil {
			fmt.Println(err)
		}

		err = w.Checkout(&git.CheckoutOptions{
			Branch: fmt.Sprintf("refs/heads/%s", branch),
			Force:  true,
		})
		if err != nil {
			fmt.Println(err)
		}
	*/
	return nil
}

func (gitRepositoryClone gitRepositoryClone) FindFilesBySuffix(suffix string) []string {
	var files []string

	filepath.Walk(gitRepositoryClone.clonePath, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), suffix) {
			files = append(files, path)
		}
		return nil
	})

	return files
}
