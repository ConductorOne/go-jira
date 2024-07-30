package onpremise

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// StatusService handles staties for the Jira instance / API.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v2/#api-group-Workflow-statuses
type StatusService service

// Status represents the current status of a Jira issue.
// Typical status are "Open", "In Progress", "Closed", ...
// Status can be user defined in every Jira instance.
type Status struct {
	Self           string         `json:"self" structs:"self"`
	Description    string         `json:"description" structs:"description"`
	IconURL        string         `json:"iconUrl" structs:"iconUrl"`
	Name           string         `json:"name" structs:"name"`
	ID             string         `json:"id" structs:"id"`
	StatusCategory StatusCategory `json:"statusCategory" structs:"statusCategory"`
}

type StatusSearchOptions struct {
	// StartAt: The starting index of the returned projects. Base index: 0.
	StartAt int `url:"startAt,omitempty"`
	// MaxResults: The maximum number of projects to return per page. Default: 50.
	MaxResults   int      `url:"maxResults,omitempty"`
	Query        string   `url:"query,omitempty"`
	ProjectIds   []string `url:"projectIds,omitempty"`
	IssueTypeIds []string `url:"issueTypeIds,omitempty"`
}

// GetAllStatuses returns a list of all statuses associated with workflows.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v2/#api-rest-api-2-status-get
//
// TODO Double check this method if this works as expected, is using the latest API and the response is complete
// This double check effort is done for v2 - Remove this two lines if this is completed.
func (s *StatusService) GetAllStatuses(ctx context.Context) ([]Status, *Response, error) {
	apiEndpoint := "rest/api/2/status"
	req, err := s.client.NewRequest(ctx, http.MethodGet, apiEndpoint, nil)

	if err != nil {
		return nil, nil, err
	}

	statusList := []Status{}
	resp, err := s.client.Do(req, &statusList)
	if err != nil {
		return nil, resp, NewJiraError(resp, err)
	}

	return statusList, resp, nil
}

// https://docs.atlassian.com/software/jira/docs/api/REST/9.17.0/#api/2/status-getPaginatedStatuses
func (s *StatusService) GetStatusesPaginated(ctx context.Context, options *StatusSearchOptions) ([]Status, *Response, error) {
	u := url.URL{
		Path: "rest/api/2/status/page",
	}
	uv := url.Values{}
	if options != nil {
		if options.StartAt != 0 {
			uv.Add("startAt", strconv.Itoa(options.StartAt))
		}
		if options.MaxResults != 0 {
			uv.Add("maxResults", strconv.Itoa(options.MaxResults))
		}
		if options.Query != "" {
			uv.Add("query", options.Query)
		}
		if len(options.ProjectIds) > 0 {
			uv.Add("projectIds", strings.Join(options.ProjectIds, ","))
		}
		if len(options.IssueTypeIds) > 0 {
			uv.Add("issueTypeIds", strings.Join(options.IssueTypeIds, ","))
		}
	}

	u.RawQuery = uv.Encode()

	req, err := s.client.NewRequest(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	statusList := []Status{}
	resp, err := s.client.Do(req, &statusList)
	if err != nil {
		return nil, resp, NewJiraError(resp, err)
	}

	return statusList, resp, nil
}
