package main

import (
	"fmt"
	"go-play/git"
	"go-play/gitlab"
	"go-play/gradleparser"
	"go-play/jira"
)

// as inputs, we expect:
//   - the jira ticket for the work
//   - maven group
//   - maven artifact ID
//   - the version to be set
//
// the tool will: 
//   - check the JIRA ticket is a task, story or bug, that it has the expected status (ready_for_development)
//   - find all "jsrvc" repositories that contain a reference to the specified dependency
//   - for each repository that has a dependency, in its own goroutine
//     -- create a subtask under the specified ticket, transition the subtasks to
//        "in progress" status.
//     -- wait for the branch to be created, pull  and check it out
//     -- set the version dependency commit and push
func main() {
	exists, err := jira.NewJiraClient().CheckTicketExistsAndIsStoryTaskOrBug("LS-14523")
	if err != nil {
		println("for an error" + err.Error())
	} else {
		if exists {
			println("ticket status:  exists")
		} else {
			println("ticket status:  exists not")
		}
	}
	/*	gitlabClient := gitlab.NewGitlabApi()
		var repositories, err = gitlabClient.SearchProjects("jsrvc")
		if err != nil {
			println("error ", err)
		} else {
			channel := make(chan string)
			for _, repo := range repositories {
				//fmt.Printf("found repo: %s\n", repo.Name)
				channel = searchDependency(channel, repo)
			}
			var x = ""
			for _, s := range repositories {
				s2 := <-channel
				x += s.Name + "/" + s2 + "-"
			}
			println("\n\n\nDone: " + x)
		}
	*/
}

func searchDependency(channel chan string, repo gitlab.RepositoryDescriptor) chan string {
	go func(s *gitlab.RepositoryDescriptor) {
		//println("Cloning... ", s.Name, " from ", s.HttpUrlToRepo)
		targetDependency := gradleparser.MavenDependency{
			Scope:    "implementation",
			Group:    "io.liquidshare.microservice.staticdata",
			Artefact: "staticdata-service-proto",
			Version:  "3.9",
		}

		uncloned := git.NewRepositoryClone(s.HttpUrlToRepo)
		cloned, _ := uncloned.Clone()
		filesPaths := cloned.FindFilesBySuffix("build.gradle")
		for _, filePath := range filesPaths {
			contains, _ := gradleparser.ContainsDependency(filePath, &targetDependency)
			if contains {
				fmt.Printf("dependency found in repo %s\n", s.Name)
				//fmt.Printf("file %s in repo %s contains targetDependency %s\n", filePath, s.Name, targetDependency.ToString())
				// TODO create subtask in JIRA for the specified story, with the proper repo name, put it "in progress" ,
				// TODO wait for branch to be created, check it out
				// TODO make change on branch, commit and push branch to remote,
				// TODO create MR in gitlab for branch, set to merge automatically
			}
		}
		channel <- s.Name
	}(&repo)
	return channel
}
