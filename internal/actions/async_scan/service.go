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
	scans map[string]*Scan
}

type Scan struct {
	repos      []string
	rulePath   string
	progress   chan int
	result     chan (Result)
	cancelChan chan struct{}
}

func NewService() Service {
	return Service{
		scans: make(map[string]*Scan, 0),
	}
}

func newScan(repos []string) *Scan {
	return &Scan{
		repos:      repos,
		progress:   make(chan int, 1),
		result:     make(chan Result),
		cancelChan: make(chan struct{}),
	}
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

	scan := newScan(repos)
	go performScan(context.Background(), scan)

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
		// Construct the command to run semgrep
		fmt.Printf("Scanning %s\n", repo)
		cmd := exec.CommandContext(ctx, "semgrep", "-f", scan.rulePath, "--json", repo)
		// Run semgrep command and capture the output
		output, err := cmd.Output()
		if err != nil {
			fmt.Printf("Error exec: %+v\n", err)
			scan.result <- Result{}
			return
		}
		err = results.addRepoResult(output)
		if err != nil {
			fmt.Printf("Error add result: %+v\n", err)
			scan.result <- Result{}
			return
		}
		if len(scan.progress) == 1 {
			<-scan.progress
		}
		fmt.Printf("progress: %d\n", i)
		scan.progress <- i
	}
	scan.result <- results
	close(scan.progress)
	close(scan.result)
	close(scan.cancelChan)
}

func (s *Service) queryProgress(scanID string) (int, error) {
	scan, found := s.scans[scanID]
	if !found {
		return 0, fmt.Errorf("no scan active for %s", scanID)
	}
	fmt.Printf("found scan %+v\n", scan)
	select {
	case progress, ok := <-scan.progress:
		if !ok {
			delete(s.scans, scanID)
			return 0, fmt.Errorf("channel closed for %s", scanID) // Channel closed, exit the function
		}
		return progress, nil
	case <-scan.cancelChan:
		fmt.Println("Cancelling async operation...")
		delete(s.scans, scanID)
		return 0, fmt.Errorf("scan canceled %s", scanID)
	}
}

func (s *Service) getResult(scanID string) (Result, error) {
	scan, found := s.scans[scanID]
	if !found {
		return Result{}, fmt.Errorf("no scan active for %s", scanID)
	}
	defer func() { delete(s.scans, scanID) }()
	result := <-scan.result
	//return error result not ready if it is not in the queue
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
