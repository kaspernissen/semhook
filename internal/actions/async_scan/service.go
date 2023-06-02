package asyncscan

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/Microsoft/go-winio/pkg/guid"
)

type Service struct {
	scans map[string]*Scan
	mu    sync.Mutex
}

func NewService() Service {
	return Service{
		scans: make(map[string]*Scan, 0),
	}
}

type Scan struct {
	repos        []string
	rulePath     string
	progress     int64
	result       chan Result
	cancelChan   chan struct{}
	progressLock sync.Mutex
}

func newScan(repos []string) *Scan {
	return &Scan{
		repos:      repos,
		result:     make(chan Result),
		cancelChan: make(chan struct{}),
	}
}

func (s *Service) scan(ctx context.Context, rulePath, repoRoot string) (string, error) {
	// Get a list of all repositories in the repoRoot directory
	repos, err := getRepositories(repoRoot)
	if err != nil {
		return "", err
	}

	scan := newScan(repos)

	scanID, err := guid.NewV4()
	if err != nil {
		return "", err
	}

	s.mu.Lock()
	s.scans[scanID.String()] = scan
	s.mu.Unlock()

	var eg sync.WaitGroup
	eg.Add(1)

	go func() {
		defer eg.Done()
		performScan(ctx, scan)
	}()

	return scanID.String(), nil
}

func performScan(ctx context.Context, scan *Scan) {
	results := NewResult(scan.rulePath)
	var scanErr error

	atomic.StoreInt64(&scan.progress, 0)

	var eg sync.WaitGroup
	eg.Add(len(scan.repos))

	for _, repo := range scan.repos {
		go func(repo string) {
			defer eg.Done()

			// Construct the command to run semgrep
			fmt.Printf("Scanning %s\n", repo)
			cmd := exec.CommandContext(context.Background(), "semgrep", "-f", scan.rulePath, "--json", repo)

			// Run semgrep command and capture the output
			output, err := cmd.Output()
			if err != nil {
				fmt.Printf("Error exec: %+v\n", err)
				scanErr = err
				return
			}

			atomic.AddInt64(&scan.progress, 1) // Increment progress counter

			err = results.addRepoResult(output)
			if err != nil {
				fmt.Printf("Error add result: %+v\n", err)
				scanErr = err
				return
			}
		}(repo)
	}

	go func() {
		eg.Wait()

		// Close channels and handle errors
		close(scan.cancelChan)

		if scanErr != nil {
			scan.result <- Result{} // Send empty result to indicate an error
		} else {
			scan.result <- results
		}
	}()
}

func (s *Service) queryProgress(scanID string) (int64, error) {
	s.mu.Lock()
	scan, found := s.scans[scanID]
	s.mu.Unlock()

	if !found {
		return 0, fmt.Errorf("no scan active for %s", scanID)
	}

	scan.progressLock.Lock()
	defer scan.progressLock.Unlock()

	progress := atomic.LoadInt64(&scan.progress)
	return progress, nil
}

func (s *Service) getResult(scanID string) (Result, error) {
	scan, found := s.scans[scanID]
	if !found {
		return Result{}, fmt.Errorf("no scan active for %s", scanID)
	}
	defer func() {
		delete(s.scans, scanID)
	}()
	result := <-scan.result
	close(scan.result)
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
