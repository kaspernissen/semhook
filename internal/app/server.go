package app

import (
	"fmt"
	"github/com/hoeg/semhook/internal/actions/repo"
	"github/com/hoeg/semhook/internal/actions/scan"
	"github/com/hoeg/semhook/internal/actions/sync"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/g4s8/go-lifecycle/pkg/adaptors"
	"github.com/g4s8/go-lifecycle/pkg/lifecycle"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Start() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "*"
		},
		MaxAge: 12 * time.Hour,
	}))

	repoRoot := repositoryRoot()
	asyncHandler := scan.NewScanHandler()

	router.GET("/sync", sync.Handler())
	router.GET("/repo", repo.Handler())

	router.POST("/scan", asyncHandler.HandlerStart(repoRoot))
	router.GET("/scan/progress", asyncHandler.HandlerStatus())
	router.GET("/scan", asyncHandler.HandlerGetResult())

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
