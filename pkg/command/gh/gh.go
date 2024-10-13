package gh

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// Repo struct to represent a repository
type Repo struct {
	ID    string    `json:"id"`
	Name  string    `json:"name"`
	Owner RepoOwner `json:"owner"`
}

// RepoOwner struct to represent the owner of the repository
type RepoOwner struct {
	ID    string `json:"id"`
	Login string `json:"login"`
}

// RepoDetail struct to represent the detail of a repository
type RepoDetail struct {
	AssignableUsers  []RepoDetailAssignableUser `json:"assignableUsers"`
	DefaultBranchRef RepoDetailDefaultBranchRef `json:"defaultBranchRef"`
}

// RepoDetailAssignableUser struct to represent an assignable user
type RepoDetailAssignableUser struct {
	ID    string `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
}

// RepoDetailDefaultBranchRef struct to represent the default branch of a repository
type RepoDetailDefaultBranchRef struct {
	Name string `json:"name"`
}

// PullRequest struct to represent a pull request
type PullRequest struct {
	ID     string            `json:"id"`
	Number int               `json:"number"`
	Title  string            `json:"title"`
	Author PullRequestAuthor `json:"author"`
}

// PullRequestAuthor struct to represent the author of a pull request
type PullRequestAuthor struct {
	ID    string `json:"id"`
	IsBot bool   `json:"is_bot"`
	Login string `json:"login"`
	Name  string `json:"name"`
}

type Commit struct {
	Sha     string `json:"sha"`
	Message string `json:"message"`
	Author  string `json:"author"`
}

// CreatePullRequestParams struct to represent the parameters for creating a pull request
type CreatePullRequestParams struct {
	BaseBranch string
	HeadBranch string
	Title      string
	Body       string
	Reviewers  []string
	IsDraft    bool
}

// GetRepositories function to get the repositories from GitHub
func GetRepositories(org string, limit int) []Repo {
	// Run the GitHub CLI command and capture the output
	cmd := exec.Command(
		"gh",
		"repo",
		"list",
		org,
		"--limit",
		fmt.Sprintf("%d", limit),
		"--json",
		"id,name,owner",
	)
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to execute gh command: %v", err)
	}

	// Parse the JSON output into a slice of repo structs
	var rs []Repo
	err = json.Unmarshal(output, &rs)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	return rs
}

// GetRepositoryDetail function to get the detail of a repository from GitHub
func GetRepositoryDetail(
	repo string,
	ch chan<- RepoDetail,
) RepoDetail {
	// Run the GitHub CLI command and capture the output
	cmd := exec.Command(
		"gh",
		"repo",
		"view",
		repo,
		"--json",
		"assignableUsers,defaultBranchRef",
	)
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to execute gh command: %v", err)
	}

	// Parse the JSON output into a slice of repo detail structs
	var repoDetail RepoDetail
	err = json.Unmarshal(output, &repoDetail)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if ch != nil {
		ch <- repoDetail

		return RepoDetail{}
	}

	return repoDetail
}

// GetPullRequests function to get the pull requests per repository from GitHub
func GetPullRequests(repo string, limit int, ch chan<- []PullRequest) []PullRequest {
	// Run the GitHub CLI command and capture the output
	cmd := exec.Command(
		"gh",
		"pr",
		"list",
		"--repo",
		repo,
		"--limit",
		fmt.Sprintf("%d", limit),
		"--json",
		"title,url",
	)
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to execute gh command: %v", err)
	}

	// Parse the JSON output into a slice of pull request structs
	var prs []PullRequest
	err = json.Unmarshal(output, &prs)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if ch != nil {
		ch <- prs

		return nil
	}

	return prs
}

// SplitCommitSummaryAndDescription function to split the commit message into summary and description
func SplitCommitSummaryAndDescription(commitMessage string) (string, string) {
	parts := strings.SplitN(commitMessage, "\n\n", 2)

	if len(parts) == 2 {
		return parts[0], parts[1]
	}

	return commitMessage, ""
}

// GetBranchCommits function to get the commits between two branches from GitHub
func GetBranchCommits(owner string, repo string, baseBranch string, headBranch string) []Commit {
	// Run the GitHub CLI command and capture the output
	cmd := exec.Command(
		"gh",
		"api",
		fmt.Sprintf("repos/%s/%s/compare/%s...%s", owner, repo, baseBranch, headBranch),
		"--jq",
		"[.commits[] | {sha: .sha, message: .commit.message, author: .commit.author.name}]",
	)
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to execute gh command: %v", err)
	}

	// Parse the JSON output into a slice of strings
	var commits []Commit
	err = json.Unmarshal(output, &commits)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	return commits
}

// CreatePullRequest function to create a pull request on GitHub
func CreatePullRequest(owner string, repo string, options CreatePullRequestParams) bool {
	if options.BaseBranch == "" || options.HeadBranch == "" || options.Title == "" {
		log.Fatalf("BaseBranch, HeadBranch, and Title are required")
	}

	args := []string{
		"pr",
		"create",
		"--repo", fmt.Sprintf("%s/%s", owner, repo),
		"--base", options.BaseBranch,
		"--head", options.HeadBranch,
		"--title", options.Title,
	}

	// Append "--body" and options.body if options.body is available
	if options.Body != "" {
		args = append(args, "--body", options.Body)
	}

	// Append draft flag if options.IsDraft is true
	if options.IsDraft == true {
		args = append(args, "--draft")
	}

	// Append reviewers if options.Reviewers is available
	if len(options.Reviewers) > 0 {
		args = append(args, "--reviewer", strings.Join(options.Reviewers, ","))
	}

	// Run the GitHub CLI command and capture the output
	cmd := exec.Command("gh", args...)
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to execute gh command: %v", err)
	}

	// Parse the JSON output to check if the pull request was created successfully
	var pr map[string]interface{}
	err = json.Unmarshal(output, &pr)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	return pr["number"] != nil
}
