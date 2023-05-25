package app

import (
	"fmt"
	"github/com/hoeg/semhook/internal/scan"
	"github/com/hoeg/semhook/internal/sync"
	"log"
	"net/http"
	"os"

	"github.com/g4s8/go-lifecycle/pkg/adaptors"
	"github.com/g4s8/go-lifecycle/pkg/lifecycle"
	"github.com/gin-gonic/gin"
)

func Start() {
	router := gin.Default()

	repoRoot := repositoryRoot()

	router.POST("/scan", scan.Handler(repoRoot))
	router.GET("/sync", sync.Handler())

	svc := adaptors.NewHTTPService(&http.Server{
		Addr:    fmt.Sprintf(":%s", port()),
		Handler: router,
	})

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
