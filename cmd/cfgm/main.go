package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/l-vitaly/go-kit/log/level"
	"github.com/l-vitaly/imhio-task/pkg/config"
	"github.com/l-vitaly/imhio-task/pkg/configurator"
	"github.com/l-vitaly/imhio-task/pkg/postgres"
)

// build vars
var (
	GitHash   = ""
	BuildDate = ""
)

const logLevelEnvName = "CFGM_LOG_LEVEL"

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() + 1)
}

func main() {
	var logger log.Logger

	level.Key = "log_level"

	logger = log.NewJSONLogger(os.Stdout)
	defer level.Info(logger).Log("msg", "goodbye")

	logger = log.With(logger, "@message", "cfgm")
	logger = log.With(logger, "@timestamp", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.Caller(5))
	logger = level.NewFilter(logger, getLevelOption())

	level.Info(logger).Log("verpon", GitHash, "builddate", BuildDate, "msg", "hello")

	cfg, err := config.Get()
	if err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}

	pgOpts, err := pg.ParseURL(cfg.DB.URL)
	if err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}

	db := pg.Connect(pgOpts)

	catalogs := postgres.NewCatalogRepository(db)
	_, _, err = catalogs.CreateSchemas()
	if err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}
	values := postgres.NewValueRepository(db)
	_, _, err = values.CreateSchemas()
	if err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}
	s := configurator.NewService(catalogs, values)

	httpLogger := log.With(logger, "component", "http")

	h := mux.NewRouter().StrictSlash(true)
	h.PathPrefix("/cfg").Handler(configurator.MakeHandler(s, httpLogger))

	srv := &http.Server{
		Addr:         cfg.HTTPAddr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      h,
	}

	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", cfg.HTTPAddr, "msg", "listening")
		errs <- srv.ListenAndServe()
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)
}

func getLevelOption() level.Option {
	switch os.Getenv(logLevelEnvName) {
	case "error":
		return level.AllowError()
	case "info":
		return level.AllowInfo()
	default:
		return level.AllowDebug()
	}
}
