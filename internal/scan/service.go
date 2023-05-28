package scan

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Runs semgrep on all repositories cloned into the repoRoot directory
// with the rule located at rulePath. The output of the run of each repo
// is combined and returned as an *bytes.Reader.
func scan(rulePath, repoRoot string) (Result, error) {
	results := NewResult("test")

	// Get a list of all repositories in the repoRoot directory
	repos, err := getRepositories(repoRoot)
	if err != nil {
		return Result{}, err
	}

	// Iterate over each repository and run semgrep
	for _, repo := range repos {
		// Construct the command to run semgrep
		cmd := exec.Command("semgrep", "-f", rulePath, "--json", repo)

		// Run semgrep command and capture the output
		output, err := cmd.Output()
		if err != nil {
			return Result{}, err
		}
		fmt.Printf("repo: %s - result: %s\n", repo, output)
		err = results.addRepoResult(output)
		if err != nil {
			return Result{}, err
		}
	}
	fmt.Printf("All: %+v", results)
	return results, nil
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
