package jira

import (
	_ "github.com/joho/godotenv"
	jira "gopkg.in/andygrunwald/go-jira.v1"
	"log"
	"os"
)

type JiraService interface {
	CheckTicketExistsAndIsStoryTaskOrBug(ticket string) (bool, error)
	CreateSubtaskInTaskOrStory(parent string, subtask JiraSubtask)
	TransitionTicket(ticket string, newStatus string)
}

type jiraService struct {
	jiraClient *jira.Client
}

func NewJiraClient() JiraService {
	tp := jira.BasicAuthTransport{
		Username: os.Getenv("JIRA_API_USER"),
		Password: os.Getenv("JIRA_API_TOKEN"),
	}
	// Create a new Jira Client
	client, err := jira.NewClient(tp.Client(), os.Getenv("JIRA_API_BASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	return &jiraService{
		jiraClient: client,
	}
}

func (jira jiraService) CheckTicketExistsAndIsStoryTaskOrBug(ticket string) (bool, error) {
	acceptedIssueTypes := map[string]struct{}{"Bug": {}, "Story": {}, "Task": {}}

	t, r, err := jira.jiraClient.Issue.Get(ticket, nil)
	if err == nil && r.StatusCode == 200 {
		_, ok := acceptedIssueTypes[t.Fields.Type.Name]
		if ok {
			return true, nil
		} else {
			return false, nil
		}
	} else {
		return false, err
	}
}

func (jira jiraService) CreateSubtaskInTaskOrStory(parent string, subtask JiraSubtask) {
	//TODO implement me
	panic("implement me")
}

func (jira jiraService) TransitionTicket(ticket string, newStatus string) {
	//TODO implement me
	panic("implement me")
}
