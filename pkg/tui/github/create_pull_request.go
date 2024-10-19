package github

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/coding-for-fun-org/go-playground/pkg/command/gh"
	"github.com/coding-for-fun-org/go-playground/pkg/command/git"
)

type CreatePullRequest struct {
	repoOwner       string
	repoName        string
	assignableUsers []gh.RepoDetailAssignableUser
	defaultBranch   string
	latestBranches  []git.Branch
	// prType         string
	// issueId        []string
	title      string
	body       string
	baseBranch string
	headBranch string
	reviewers  []string
	// isDraft         bool
}

// extractPatterns function extracts all occurrences of the pattern [A-Z]{1,}-\d{1,} from the input string.
func extractPatterns(input string) ([]string, error) {
	// Define the regular expression
	pattern := `[A-Z]+-\d+`

	// Compile the regular expression
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	// Find all matches
	matches := re.FindAllString(input, -1)
	return matches, nil
}

// concatenateAndRemoveDuplicates function to concatenate two slices and remove duplicates
func concatenateAndRemoveDuplicates(slice1, slice2 []string) []string {
	// Use a map to track unique items
	uniqueMap := make(map[string]struct{})
	result := make([]string, 0)

	// Helper function to add items to the result while checking uniqueness
	addUnique := func(item string) {
		if _, exists := uniqueMap[item]; !exists {
			uniqueMap[item] = struct{}{}
			result = append(result, item)
		}
	}

	// Add elements from the first slice
	for _, item := range slice1 {
		addUnique(item)
	}

	// Add elements from the second slice
	for _, item := range slice2 {
		addUnique(item)
	}

	return result
}

// splitCommitSummaryAndDescription function to split the commit message into summary and description
func splitCommitSummaryAndDescription(commitMessage string) (string, string) {
	parts := strings.SplitN(commitMessage, "\n\n", 2)

	if len(parts) == 2 {
		return parts[0], parts[1]
	}

	return commitMessage, ""
}

// branchForm method to create a form for selecting the base and head branches
func (p *CreatePullRequest) branchForm() *huh.Form {
	branchForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select the head branch").
				Options((func() []huh.Option[string] {
					// filter out the default branch
					branches := make([]huh.Option[string], 0)
					for _, branch := range p.latestBranches {
						if branch.Ref != p.defaultBranch {
							branches = append(
								branches,
								huh.NewOption(branch.Ref, branch.Ref),
							)
						}
					}

					return branches
				})()...).
				Value(&p.headBranch),

			// I'd like to put this in another group but
			// because of the current bug, I can not do it.
			// https://github.com/charmbracelet/huh/issues/419
			huh.NewSelect[string]().
				Title("Select the base branch").
				OptionsFunc(func() []huh.Option[string] {
					// filter out the default branch
					branches := make([]huh.Option[string], 0)
					branches = append(
						branches,
						huh.NewOption(p.defaultBranch, p.defaultBranch),
					)
					for _, branch := range p.latestBranches {
						if branch.Ref != p.defaultBranch && branch.Ref != p.headBranch {
							branches = append(
								branches,
								huh.NewOption(branch.Ref, branch.Ref),
							)
						}
					}

					return branches
				}, &p.headBranch).
				Value(&p.baseBranch),
		),
	)

	return branchForm
}

// initializeBaseInfo method to initialize the base information for creating a pull request
func (p *CreatePullRequest) initializeBaseInfo() {
	repo := gh.GetRepositoryDetail("RockRabbit-ai/rockrabbit-web", nil)
	p.repoOwner = repo.Owner.Login
	p.repoName = repo.Name
	p.assignableUsers = repo.AssignableUsers
	p.defaultBranch = repo.DefaultBranchRef.Name
	p.latestBranches = git.GetLatestBranches()
}

func (p *CreatePullRequest) getPrePopulatedTitleAndBody(commits []gh.Commit) (string, string) {
	if len(commits) == 1 {
		commitTitle, commitBody := splitCommitSummaryAndDescription(commits[0].Message)
		issueNumbersFromTitle, _ := extractPatterns(commitTitle)
		issueNumbersFromBody, _ := extractPatterns(commitBody)
		issueNumbers := concatenateAndRemoveDuplicates(
			issueNumbersFromTitle,
			issueNumbersFromBody,
		)
		commitBody = commitBody + "\n\n### Jira Link\n\n"
		for _, issueNumber := range issueNumbers {
			commitBody = commitBody + fmt.Sprintf(
				"[%s](https://keends.atlassian.net/browse/%s)\n",
				issueNumber,
				issueNumber,
			)
		}
		return commitTitle, commitBody
	}

	commitFullBody := "### Jira Link\n\n"
	for _, commit := range commits {
		commitTitle, commitBody := splitCommitSummaryAndDescription(commit.Message)
		issueNumbersFromTitle, _ := extractPatterns(commitTitle)
		issueNumbersFromBody, _ := extractPatterns(commitBody)
		issueNumbers := concatenateAndRemoveDuplicates(
			issueNumbersFromTitle,
			issueNumbersFromBody,
		)
		for _, issueNumber := range issueNumbers {
			commitFullBody = commitFullBody + fmt.Sprintf(
				"[%s](https://keends.atlassian.net/browse/%s)\n",
				issueNumber,
				issueNumber,
			)
		}
		return commitTitle, commitFullBody
	}

	return "", ""
}

func (p *CreatePullRequest) initializePullRequestTitleAndBody() {
	commits := gh.GetBranchCommits(p.repoOwner, p.repoName, p.baseBranch, p.headBranch)

	p.getPrePopulatedTitleAndBody(commits)
}

func (p *CreatePullRequest) Run() {
	initializeBaseInfo := p.initializeBaseInfo
	spinner.New().
		Title("Loading base information to create a pull request...").
		Action(initializeBaseInfo).
		Run()

	branchForm := p.branchForm()
	err := branchForm.Run()
	if err != nil {
		log.Fatal(err)
	}

	// If the user stops the program, we don't want to go to the next form
	if branchForm.State == huh.StateAborted {
		fmt.Println("Aborted")
		return
	}

	initializePullRequestTitleAndBody := p.initializePullRequestTitleAndBody

	spinner.New().
		Title("Loading").
		Action(initializePullRequestTitleAndBody).
		Run()

		// Ask title - pre-populate with the first line of the commit message
		// Ask body - pre-populate with the rest of the commit message
		// Ask isDraft
		// Ask reviewers

	huh.NewMultiSelect[string]().
		Title("Select reviewers").
		Options((func() []huh.Option[string] {
			myUserLogin := gh.GetMyUserLogin()
			users := make([]huh.Option[string], 0)
			for _, user := range p.assignableUsers {
				if user.Login == myUserLogin {
					continue
				}
				format := "%s"
				if user.Name != "" {
					format += " (%s)"
				}
				spread := []interface{}{user.Login}
				if user.Name != "" {
					spread = append(spread, user.Name)
				}

				users = append(
					users,
					huh.NewOption(fmt.Sprintf(format, spread...), user.Login),
				)
			}

			return users
		})()...).
		Value(&p.reviewers)
}
