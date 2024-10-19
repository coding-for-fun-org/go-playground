package main

import (
	"github.com/coding-for-fun-org/go-playground/pkg/tui/github"
)

func main() {
	// branches := git.GetLatestBranches()
	//
	// for _, branch := range branches {
	// 	log.Printf(
	// 		"Branch: %s, Commit: %s, Date: %s\n",
	// 		branch.RemoteBranch,
	// 		branch.Commit,
	// 		branch.Date,
	// 	)
	// }

	pr := &github.CreatePullRequest{}
	pr.Run()

	// repos := gh.GetRepositories("coding-for-fun-org", 5)
	// c := make(chan gh.RepoDetail)
	// // Print the result
	// for _, repo := range repos {
	// 	fmt.Printf("Repository Name: %s\n", fmt.Sprintf("%s/%s", repo.Owner.Login, repo.Name))
	//
	// 	go gh.GetRepositoryDetail(fmt.Sprintf("%s/%s", repo.Owner.Login, repo.Name), c)
	//
	// }
	//
	// for range repos {
	// 	repoDetail := <-c
	//
	// 	for _, user := range repoDetail.AssignableUsers {
	// 		fmt.Printf("User ID: %s, Login: %s, Name: %s\n", user.ID, user.Login, user.Name)
	// 	}
	// }
	//
	// headBranch := "feat/KPC-3130/hello-greeting"
	// branchCommits := gh.GetBranchCommits(
	// 	"coding-for-fun-org",
	// 	"frontend",
	// 	"main",
	// 	headBranch,
	// )
	//
	// if len(branchCommits) == 1 {
	// 	commit := branchCommits[0]
	//
	// 	title, body := gh.SplitCommitSummaryAndDescription(commit.Message)
	//
	// 	rawCommitSummary, _ := json.Marshal(title)
	// 	log.Println(string(rawCommitSummary))
	//
	// 	rawCommitDescription, _ := json.Marshal(body)
	// 	log.Println(string(rawCommitDescription))
	// } else {
	// 	for _, commit := range branchCommits {
	// 		title, body := gh.SplitCommitSummaryAndDescription(commit.Message)
	//
	// 		rawCommitSummary, _ := json.Marshal(title)
	// 		log.Println(string(rawCommitSummary))
	//
	// 		rawCommitDescription, _ := json.Marshal(body)
	// 		log.Println(string(rawCommitDescription))
	// 	}
	// }

	// github.Run()
}
