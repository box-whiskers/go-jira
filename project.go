package jira

import (
	"fmt"

	"github.com/google/go-querystring/query"
)

// ProjectService handles projects for the JIRA instance / API.
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/latest/#api/2/project
type ProjectService interface {
	GetList() (*ProjectList, *Response, error)
	ListWithOptions(*GetQueryOptions) (*ProjectList, *Response, error)
	Get(string) (*Project, *Response, error)
	GetPermissionScheme(string) (*PermissionScheme, *Response, error)
}

type ProjectServiceImpl struct {
	client *Client
}

// ProjectList represent a list of Projects
type ProjectList []struct {
	Expand          string          `json:"expand" structs:"expand"`
	Self            string          `json:"self" structs:"self"`
	ID              string          `json:"id" structs:"id"`
	Key             string          `json:"key" structs:"key"`
	Name            string          `json:"name" structs:"name"`
	AvatarUrls      AvatarUrls      `json:"avatarUrls" structs:"avatarUrls"`
	ProjectTypeKey  string          `json:"projectTypeKey" structs:"projectTypeKey"`
	ProjectCategory ProjectCategory `json:"projectCategory,omitempty" structs:"projectsCategory,omitempty"`
	IssueTypes      []IssueType     `json:"issueTypes,omitempty" structs:"issueTypes,omitempty"`
}

// ProjectCategory represents a single project category
type ProjectCategory struct {
	Self        string `json:"self" structs:"self,omitempty"`
	ID          string `json:"id" structs:"id,omitempty"`
	Name        string `json:"name" structs:"name,omitempty"`
	Description string `json:"description" structs:"description,omitempty"`
}

// Project represents a JIRA Project.
type Project struct {
	Expand          string             `json:"expand,omitempty" structs:"expand,omitempty"`
	Self            string             `json:"self,omitempty" structs:"self,omitempty"`
	ID              string             `json:"id,omitempty" structs:"id,omitempty"`
	Key             string             `json:"key,omitempty" structs:"key,omitempty"`
	Description     string             `json:"description,omitempty" structs:"description,omitempty"`
	Lead            User               `json:"lead,omitempty" structs:"lead,omitempty"`
	Components      []ProjectComponent `json:"components,omitempty" structs:"components,omitempty"`
	IssueTypes      []IssueType        `json:"issueTypes,omitempty" structs:"issueTypes,omitempty"`
	URL             string             `json:"url,omitempty" structs:"url,omitempty"`
	Email           string             `json:"email,omitempty" structs:"email,omitempty"`
	AssigneeType    string             `json:"assigneeType,omitempty" structs:"assigneeType,omitempty"`
	Versions        []Version          `json:"versions,omitempty" structs:"versions,omitempty"`
	Name            string             `json:"name,omitempty" structs:"name,omitempty"`
	Roles           map[string]string  `json:"roles,omitempty" structs:"roles,omitempty"`
	AvatarUrls      AvatarUrls         `json:"avatarUrls,omitempty" structs:"avatarUrls,omitempty"`
	ProjectCategory ProjectCategory    `json:"projectCategory,omitempty" structs:"projectCategory,omitempty"`
}

// ProjectComponent represents a single component of a project
type ProjectComponent struct {
	Self                string `json:"self" structs:"self,omitempty"`
	ID                  string `json:"id" structs:"id,omitempty"`
	Name                string `json:"name" structs:"name,omitempty"`
	Description         string `json:"description" structs:"description,omitempty"`
	Lead                User   `json:"lead,omitempty" structs:"lead,omitempty"`
	AssigneeType        string `json:"assigneeType" structs:"assigneeType,omitempty"`
	Assignee            User   `json:"assignee" structs:"assignee,omitempty"`
	RealAssigneeType    string `json:"realAssigneeType" structs:"realAssigneeType,omitempty"`
	RealAssignee        User   `json:"realAssignee" structs:"realAssignee,omitempty"`
	IsAssigneeTypeValid bool   `json:"isAssigneeTypeValid" structs:"isAssigneeTypeValid,omitempty"`
	Project             string `json:"project" structs:"project,omitempty"`
	ProjectID           int    `json:"projectId" structs:"projectId,omitempty"`
}

// PermissionScheme represents the permission scheme for the project
type PermissionScheme struct {
	Expand      string       `json:"expand" structs:"expand,omitempty"`
	Self        string       `json:"self" structs:"self,omitempty"`
	ID          int          `json:"id" structs:"id,omitempty"`
	Name        string       `json:"name" structs:"name,omitempty"`
	Description string       `json:"description" structs:"description,omitempty"`
	Permissions []Permission `json:"permissions" structs:"permissions,omitempty"`
}

// GetList gets all projects form JIRA
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/latest/#api/2/project-getAllProjects
func (s *ProjectServiceImpl) GetList() (*ProjectList, *Response, error) {
	return s.ListWithOptions(&GetQueryOptions{})
}

// ListWithOptions gets all projects form JIRA with optional query params, like &GetQueryOptions{Expand: "issueTypes"} to get
// a list of all projects and their supported issuetypes
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/latest/#api/2/project-getAllProjects
func (s *ProjectServiceImpl) ListWithOptions(options *GetQueryOptions) (*ProjectList, *Response, error) {
	apiEndpoint := "rest/api/2/project"
	req, err := s.client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	if options != nil {
		q, err := query.Values(options)
		if err != nil {
			return nil, nil, err
		}
		req.URL.RawQuery = q.Encode()
	}

	projectList := new(ProjectList)
	resp, err := s.client.Do(req, projectList)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}

	return projectList, resp, nil
}

// Get returns a full representation of the project for the given issue key.
// JIRA will attempt to identify the project by the projectIdOrKey path parameter.
// This can be an project id, or an project key.
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/latest/#api/2/project-getProject
func (s *ProjectServiceImpl) Get(projectID string) (*Project, *Response, error) {
	apiEndpoint := fmt.Sprintf("rest/api/2/project/%s", projectID)
	req, err := s.client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	project := new(Project)
	resp, err := s.client.Do(req, project)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}

	return project, resp, nil
}

// GetPermissionScheme returns a full representation of the permission scheme for the project
// JIRA will attempt to identify the project by the projectIdOrKey path parameter.
// This can be an project id, or an project key.
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/latest/#api/2/project-getProject
func (s *ProjectServiceImpl) GetPermissionScheme(projectID string) (*PermissionScheme, *Response, error) {
	apiEndpoint := fmt.Sprintf("/rest/api/2/project/%s/permissionscheme", projectID)
	req, err := s.client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	ps := new(PermissionScheme)
	resp, err := s.client.Do(req, ps)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}

	return ps, resp, nil
}
