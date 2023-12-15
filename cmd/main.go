package main

import (
	"Darkyfun/UrlShortener/internal/config"
	"Darkyfun/UrlShortener/internal/logging"
	"Darkyfun/UrlShortener/internal/logging/logpath"
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

	configPath := flag.String("config", "", "logpath for config file")
	flag.Parse()

	// парсим файл конфигурации.
	conf, err := config.GetConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	// инициализируем логеры.
	fmt.Println("Reading config")
	logPaths := logpath.DestinationLog("./logs")

	fmt.Println("Setting destination for logs")

	defer func() {
		if err := logPaths.CloseFiles(); err != nil {
			fmt.Println(err)
		}
	}()

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

	// подключаемся к кэшу.
	rdb := cache.NewCacheDb(cacheOpts, baseLogger)
	defer func() {
		err = rdb.Close()
		if err != nil {
			baseLogger.Log("error", "can not close the connection to Cache Db: "+err.Error())
		}
	}()
	fmt.Println("Connected to cache database")

	// подключаемся в SQL-базе данных.
	db := persistent.NewDb(ctx, baseLogger, conf.GetString("SqlConnString"))
	defer db.Close()
	fmt.Println("Connected to persistence database")

	// бесконечный healthcheck к кэшу.
	go cache.PingCache(rdb, baseLogger)

	// бесконечный healthcheck к SQL-базе данных.
	go persistent.PingStorage(&db, baseLogger)

	// инициализируем gin.
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	middleLogger := middleware.NewLogHandler(logPaths.IncomeLog)
	router.Use(middleLogger.Logger())
	router.Use(gin.Recovery())

	router.GET("/redirect/:alias", middleware.Redirect(rdb, &db, baseLogger))
	router.POST("/receive", middleware.Validate(), middleware.Saver(rdb, &db, conf.GetString("ServerAddr")))

	server := &http.Server{
		Addr:         conf.GetString("ServerAddr"),
		Handler:      router,
		ReadTimeout:  conf.GetDuration("ReadTimeout") * time.Second,
		WriteTimeout: conf.GetDuration("WriteTimeout") * time.Second,
		IdleTimeout:  conf.GetDuration("IdleTimeout") * time.Second,
	}

	// запускаем сервер.
	fmt.Println("Starting server")
	go func() { baseLogger.Log("info", server.ListenAndServe().Error()) }()

	fmt.Println("Server has been started")

	// graceful shutdown.
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
