package app

import (
	"log"
	"os"

	"github.com/g4s8/go-lifecycle/pkg/lifecycle"
)

func Start() {
	repoRoot := repositoryRoot()
	// read github token from file
	// call init on starhook

	svc := WireHTTP(repoRoot)

	lf := lifecycle.New(lifecycle.DefaultConfig)
	svc.RegisterLifecycle("web", lf)
	lf.Start()
	sig := lifecycle.NewSignalHandler(lf, nil)
	sig.Start(lifecycle.DefaultShutdownConfig)
	if err := sig.Wait(); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
	log.Print("shutdown complete")
}

func port() string {
	port := os.Getenv("SEMHOOK_PORT")
	if port == "" {
		return "8080"
	}
	return port
}

func repositoryRoot() string {
	path := os.Getenv("SEMHOOK_REPO_ROOT")
	if path == "" {
		panic("SEMHOOK_REPO_ROOT must be set in env")
	}
	return path
}
