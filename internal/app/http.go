package app

import (
	"fmt"
	"github/com/hoeg/semhook/internal/actions/repo"
	"github/com/hoeg/semhook/internal/actions/scan"
	"github/com/hoeg/semhook/internal/actions/sync"
	"net/http"
	"time"

	"github.com/g4s8/go-lifecycle/pkg/adaptors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func WireHTTP(repoRoot string) *adaptors.HTTPService {
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

	scanHandler := scan.NewScanHandler()

	router.GET("/sync", sync.Handler())
	router.GET("/repo", repo.Handler())

	router.POST("/scan", scanHandler.HandlerStart(repoRoot))
	router.GET("/scan/progress", scanHandler.HandlerStatus())
	router.GET("/scan", scanHandler.HandlerGetResult())

	return adaptors.NewHTTPService(&http.Server{
		Addr:    fmt.Sprintf(":%s", port()),
		Handler: router,
	})
}
