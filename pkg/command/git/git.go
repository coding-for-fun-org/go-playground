package git

import (
	"encoding/json"
	"log"
	"os/exec"
	"strings"
)

type Branch struct {
	Ref    string `json:"ref"`
	Commit string `json:"commit"`
	Date   string `json:"date"`
}

func GetLatestBranches() []Branch {
	cmd := exec.Command(
		"git",
		"--git-dir=/Users/jiyeollee/code/rockrabbit/web/.git",
		"--work-tree=/Users/jiyeollee/code/rockrabbit/web",
		"for-each-ref",
		"refs/heads/",
		"--sort=-committerdate",
		"--format={\"ref\": \"%(refname:short)\", \"commit\": \"%(objectname)\", \"date\": \"%(authordate:iso8601)\"}",
	)
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to execute gh command: %v", err)
	}

	// Convert output to string and split into lines
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	// Join lines with "," and wrap with "[" and "]"
	jsonArray := "[" + strings.Join(lines, ",") + "]"

	// Parse the JSON output into a slice of Branch structs
	var bs []Branch
	err = json.Unmarshal([]byte(jsonArray), &bs)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	return bs
}
