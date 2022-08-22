package main

import (
	"context"
	"log"
	"net/url"
	"os"
	"path"
	"runtime"

	"github.com/IPFS-NEXIVIL/orbit-db-gateway/cache"
	"github.com/IPFS-NEXIVIL/orbit-db-gateway/config"
	"github.com/IPFS-NEXIVIL/orbit-db-gateway/database"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewLogger(filename string) (*zap.Logger, error) {
	if runtime.GOOS == "windows" {
		zap.RegisterSink("winfile", func(u *url.URL) (zap.Sink, error) {
			return os.OpenFile(u.Path[1:], os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		})
	}

	cfg := zap.NewDevelopmentConfig()
	if runtime.GOOS == "windows" {
		cfg.OutputPaths = []string{
			"stdout",
			"winfile:///" + filename,
		}
	} else {
		cfg.OutputPaths = []string{
			filename,
		}
	}

	return cfg.Build()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("loading configuration ...")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Panicln(err)
	}
	if cfg.WasSetup() == false {
		cfg.Setup()
	}

	log.Println("initializing logger ...")
	logger, err := NewLogger(cfg.Logfile)
	if err != nil {
		log.Panicln(err)
	}

	log.Println("initializing cache ...")
	cch, err := cache.NewCache(cfg.ProgramCachePath)
	if err != nil {
		log.Panicln(err)
	}
	defer cch.Close()

	log.Println("initializing database ...")
	db, err := database.NewDatabase(ctx, cfg.ConnectionString, cfg.DatabaseCachePath, cch, logger)
	if err != nil {
		log.Panicln(err)
	}
	defer db.Disconnect()

	// Create content storage directories
	if err := os.MkdirAll(path.Join("files", "text"), 0755); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(path.Join("files", "images"), 0755); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(path.Join("files", "other"), 0755); err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	ipfs := router.Group("/ipfs")
	{
		ipfs.POST("/paste", paste)
		ipfs.POST("/upload", upload)
	}

	router.Run("localhost:8001")
}
