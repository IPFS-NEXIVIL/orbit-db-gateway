package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"runtime"
	"time"

	"berty.tech/go-orbit-db/iface"
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
	if !cfg.WasSetup() {
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

	log.Println("connecting database ...")
	err = db.Connect(func(address string) {
	})
	if err != nil {
		log.Panicln(err)
	}

	go func() {
		for {
			_, err := db.IPFSCoreAPI.Swarm().Peers(context.Background())
			if err != nil {
				log.Panicln(err)
			}
			time.Sleep(time.Second * 5)
		}
	}()

	go func() {
		var input string
		for {
			fmt.Scanln(&input)

			switch input {
			case "q":
				return
			case "g":
				fmt.Scanln(&input)
				docs, err := db.Store.Get(ctx, input, &iface.DocumentStoreGetOptions{CaseInsensitive: false})
				if err != nil {

					log.Println(err)
				} else {
					log.Println(docs)
				}
			case "l":
				docs, err := db.Store.Query(ctx, func(e interface{}) (bool, error) {
					return true, nil
				})
				if err != nil {
					log.Println(err)
				} else {
					log.Println(docs)
				}
			}

		}
	}()

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
		dbInfo := DBInfo{
			DB: db,
		}
		// save and get data to orbit db
		ipfs.POST("/paste", dbInfo.paste)
	}

	router.Run("localhost:8001")
}
