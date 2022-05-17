package client

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/BESTSELLER/terraform-provider-harbor/models"
)

// GetProjectRepositories returns a list of repositories for the given project.
func (c *Client) GetProjectRepositories(projectName string) ([]models.RepositoryBody, error) {
	var allRepos []models.RepositoryBody

	page := 1

	// Page through the repository list until the API returns an empty result
	// or an error.
	for {
		reposPath := fmt.Sprintf("/projects/%s/repositories?page=%d&page_size=100", projectName, page)

		resp, _, _, err := c.SendRequest("GET", reposPath, nil, 200)
		if err != nil {
			return nil, err
		}

		var repos []models.RepositoryBody

		if err := json.Unmarshal([]byte(resp), &repos); err != nil {
			return nil, err
		}

		if len(repos) == 0 {
			return allRepos, nil
		}

		allRepos = append(allRepos, repos...)
		page++
	}
}

// DeleteProjectRepositories deletes all repositories of a given project.
func (c *Client) DeleteProjectRepositories(projectName string) error {
	repos, err := c.GetProjectRepositories(projectName)
	if err != nil {
		return err
	}

	// Repository names returned by the API have the form
	// <project_name>/<repository_name>.
	projectNamePrefix := fmt.Sprintf("%s/", projectName)

	for _, repo := range repos {
		repoName := strings.TrimPrefix(repo.Name, projectNamePrefix)

		// Encode slashes in the repository name as mandated by the API.
		repoName = strings.ReplaceAll(repoName, "/", "%252F")

		repoPath := fmt.Sprintf("/projects/%s/repositories/%s", projectName, repoName)

		_, _, _, err := c.SendRequest("DELETE", repoPath, nil, 200)
		if err != nil {
			return err
		}
	}

	return nil
}
