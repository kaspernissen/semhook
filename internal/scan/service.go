package scan

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Runs semgrep on all repositories cloned into the repoRoot directory
// with the rule located at rulePath. The output of the run of each repo
// is combined and returned as an *bytes.Reader.
func scan(rulePath, repoRoot string) (string, error) {
	var results []string

	// Get a list of all repositories in the repoRoot directory
	repos, err := getRepositories(repoRoot)
	if err != nil {
		return "", err
	}

	// Iterate over each repository and run semgrep
	for _, repo := range repos {
		// Construct the command to run semgrep
		cmd := exec.Command("semgrep", "-f", rulePath, "--json", repo)

		// Run semgrep command and capture the output
		output, err := cmd.CombinedOutput()
		if err != nil {
			return "", err
		}

		// Append the result to the results slice
		results = append(results, string(output))
	}

	return strings.Join(results, "\n"), nil
}

// Helper function to get a list of repositories in the repoRoot directory
func getRepositories(repoRoot string) ([]string, error) {
	var repos []string

	err := filepath.WalkDir(repoRoot, func(path string, d os.DirEntry, err error) error {
		if d.IsDir() && d.Name() == ".git" {
			// Add the parent directory of the .git directory to the list of repositories
			repos = append(repos, filepath.Dir(path))
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return repos, nil
}
