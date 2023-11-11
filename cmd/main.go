package main

import (
	"Darkyfun/UrlShortener/internal/config"
	"Darkyfun/UrlShortener/internal/logging"
	"Darkyfun/UrlShortener/internal/logging/path"
	"Darkyfun/UrlShortener/internal/server/middleware"
	"Darkyfun/UrlShortener/internal/storage/cache"
	"Darkyfun/UrlShortener/internal/storage/persistent"
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("Starting service")

	configPath := flag.String("config", "", "path for config file")
	flag.Parse()

	conf, err := config.GetConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	// initializing logger and logging path
	fmt.Println("Reading config")
	logPaths := path.DestinationLog("./logs")
	defer logPaths.IncomeLog.Close()
	defer logPaths.ErrorLog.Close()

	fmt.Println("Setting destination for logs")

	baseLogger := logging.NewLogger(conf.GetString("OutputType"), logPaths.ErrorLog)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	cacheOpts := cache.Opts{
		Addr:       conf.GetString("CacheAddr"),
		User:       conf.GetString("CacheUser"),
		Password:   conf.GetString("CachePass"),
		MaxRetries: conf.GetInt("MaxRetries"),
		PoolSize:   conf.GetInt("PoolSize"),
	}

	rdb := cache.NewCacheDb(cacheOpts, baseLogger)
	defer rdb.Close()
	fmt.Println("Connected to cache database")

	db := persistent.NewDb(ctx, baseLogger, conf.GetString("SqlConnString"))
	defer db.Close()
	fmt.Println("Connected to persistence database")

	// cache healthcheck
	go cache.PingCache(rdb, baseLogger)

	// storage healthcheck
	go persistent.PingStorage(db, baseLogger)

	// initializing Gin framework
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	middleLogger := middleware.NewLogHandler(logPaths.IncomeLog)
	router.Use(middleLogger.Logger())
	router.Use(gin.Recovery())

	router.GET("/redirect/:alias", middleware.Redirect(rdb, db, baseLogger))
	router.POST("/receive", middleware.Validate(), middleware.Saver(rdb, db, conf.GetString("ServerAddr")))

	server := &http.Server{
		Addr:         conf.GetString("ServerAddr"),
		Handler:      router,
		ReadTimeout:  conf.GetDuration("ReadTimeout") * time.Second,
		WriteTimeout: conf.GetDuration("WriteTimeout") * time.Second,
		IdleTimeout:  conf.GetDuration("IdleTimeout") * time.Second,
	}

	// Starting server
	fmt.Println("Starting server")
	go func() { baseLogger.Log("info", server.ListenAndServe().Error()) }()

	fmt.Println("Server has been started")

	quit := make(chan os.Signal)
	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	<-quit

	fmt.Println("Shutting down the server")
	serveCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err = server.Shutdown(serveCtx); err != nil {
		log.Fatalf("Shutting down: %v\n", err)
	}

	select {
	case <-serveCtx.Done():
		baseLogger.Log("warn", "Server shutdown by timeout")
	default:
		fmt.Println("Done")
	}
}
