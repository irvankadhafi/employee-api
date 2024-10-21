package console

import (
	"context"
	"fmt"
	"github.com/irvankadhafi/employee-api/cacher"
	"github.com/irvankadhafi/employee-api/internal/config"
	"github.com/irvankadhafi/employee-api/internal/db"
	httpsvc "github.com/irvankadhafi/employee-api/internal/delivery/http"
	"github.com/irvankadhafi/employee-api/internal/helper"
	"github.com/irvankadhafi/employee-api/internal/repository"
	"github.com/irvankadhafi/employee-api/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var runServerCmd = &cobra.Command{
	Use:   "server",
	Short: "run server",
	Long:  `This subcommand start the server`,
	Run:   runServer,
}

func init() {
	RootCmd.AddCommand(runServerCmd)
}

func runServer(cmd *cobra.Command, args []string) {
	// Initiate all connection like db, redis, etc
	db.InitializePostgresConn()

	pgDB, err := db.PostgreSQL.DB()
	continueOrFatal(err)
	defer helper.WrapCloser(pgDB.Close)

	cacheManager := cacher.NewCacheManager()

	cacheManager.SetDisableCaching(config.DisableCaching())

	if !config.DisableCaching() {
		redisConn, err := db.NewRedigoRedisConnectionPool(config.RedisCacheHost(), redisOpts)
		continueOrFatal(err)
		defer helper.WrapCloser(redisConn.Close)

		redisLockConn, err := db.NewRedigoRedisConnectionPool(config.RedisLockHost(), redisOpts)
		continueOrFatal(err)
		defer helper.WrapCloser(redisLockConn.Close)

		cacheManager.SetConnectionPool(redisConn)
		cacheManager.SetLockConnectionPool(redisLockConn)
		cacheManager.SetDefaultTTL(config.CacheTTL())
	}

	location, err := time.LoadLocation("Asia/Jakarta")
	continueOrFatal(err)

	time.Local = location

	employeeRepository := repository.NewEmployeeRepository(db.PostgreSQL, cacheManager)
	employeeUsecase := usecase.NewEmployeeUsecase(employeeRepository)

	httpServer := echo.New()
	httpServer.Pre(middleware.AddTrailingSlash())
	httpServer.Use(middleware.Logger())
	httpServer.Use(middleware.Recover())
	httpServer.Use(middleware.CORS())

	apiGroup := httpServer.Group("/api")
	httpsvc.RouteService(apiGroup, employeeUsecase)

	sigCh := make(chan os.Signal, 1)
	errCh := make(chan error, 1)
	quitCh := make(chan bool, 1)
	signal.Notify(sigCh, os.Interrupt)

	go func() {
		for {
			select {
			case <-sigCh:
				gracefulShutdown(httpServer)
				quitCh <- true
			case e := <-errCh:
				log.Error(e)
				gracefulShutdown(httpServer)
				quitCh <- true
			}
		}
	}()

	setupLogger()

	go func() {
		// Start HTTP server
		if err := httpServer.Start(fmt.Sprintf(":%s", config.HTTPPort())); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	<-quitCh
	log.Info("exiting")
}

func gracefulShutdown(httpSvr *echo.Echo) {
	db.StopTickerCh <- true

	if httpSvr != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := httpSvr.Shutdown(ctx); err != nil {
			httpSvr.Logger.Fatal(err)
		}
	}
}

func continueOrFatal(err error) {
	if err != nil {
		logrus.Fatal(err)
	}
}
