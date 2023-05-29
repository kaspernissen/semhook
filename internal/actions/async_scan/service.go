package asyncscan

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Microsoft/go-winio/pkg/guid"
)

type Service struct {
	scans map[string]Scan
}

type Scan struct {
	repos      []string
	rulePath   string
	progress   chan int
	result     chan (Result)
	cancelChan chan struct{}
}

// Runs semgrep on all repositories cloned into the repoRoot directory
// with the rule located at rulePath. The output of the run of each repo
// is combined and returned as an *bytes.Reader.
func (s *Service) scan(ctx context.Context, rulePath, repoRoot string) (string, error) {
	// Get a list of all repositories in the repoRoot directory
	repos, err := getRepositories(repoRoot)
	if err != nil {
		return "", err
	}
	scan := Scan{
		repos: repos,
	}
	go performScan(ctx, &scan)
	scanID, err := guid.NewV4()
	if err != nil {
		return "", err
	}
	s.scans[scanID.String()] = scan
	return scanID.String(), nil
}

func performScan(ctx context.Context, scan *Scan) {
	// Iterate over each repository and run semgrep
	results := NewResult(scan.rulePath)
	for i, repo := range scan.repos {
		select {
		case <-scan.cancelChan:
			fmt.Println("Async operation cancelled")
			return
		default:
			// Construct the command to run semgrep
			cmd := exec.CommandContext(ctx, "semgrep", "-f", scan.rulePath, "--json", repo)

			// Run semgrep command and capture the output
			output, err := cmd.Output()
			if err != nil {
				scan.result <- Result{}
			}
			err = results.addRepoResult(output)
			if err != nil {
				scan.result <- Result{}
			}
			scan.progress <- i
		}
	}
	scan.result <- results
	close(scan.progress)
	close(scan.result)
}

func (s *Service) queryProgress(scanID string) (int, error) {
	scan, found := s.scans[scanID]
	if !found {
		return 0, fmt.Errorf("no scan active for %s", scanID)
	}
	select {
	case progress, ok := <-scan.progress:
		if !ok {
			return 0, fmt.Errorf("channel closed for %s", scanID) // Channel closed, exit the function
		}
		return progress, nil
	case <-scan.cancelChan:
		fmt.Println("Cancelling async operation...")
		return 0, fmt.Errorf("scan canceled %s", scanID)
	}
}

func (s *Service) getResult(scanID string) (Result, error) {
	scan, found := s.scans[scanID]
	if !found {
		return Result{}, fmt.Errorf("no scan active for %s", scanID)
	}
	result := <-scan.result
	return result, nil
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
