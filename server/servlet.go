package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/golang/glog"
)

// Servlet ...
type Servlet struct {
	Handler http.Handler
	Conf    Conf
	DB      *sql.DB
}

// New new servlet
func New(conf Conf) (s *Servlet, err error) {
	s = &Servlet{
		Conf: conf,
	}

	// Init db
	s.DB, err = NewPostgreSQL(conf.DB)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	// clear data and create tables
	if conf.DB.Reset {
		if err = batchExecSQLFileWithSuffix(s.DB, s.Conf.DB.SQLDir, DOWN_TABLE_SUFFIX); err != nil {
			glog.Error(err)
			return
		}

		if err = batchExecSQLFileWithSuffix(s.DB, s.Conf.DB.SQLDir, UP_TABLE_SUFFIX); err != nil {
			glog.Error(err)
			return
		}
	}

	return s, nil
}

// StartOrDie start or die server
func (s *Servlet) StartOrDie() {
	var srv = http.Server{
		Addr:              fmt.Sprintf(":%d", s.Conf.Web.Port),
		WriteTimeout:      time.Second * s.Conf.Web.WriteTimeoutInSec,
		ReadHeaderTimeout: time.Second * s.Conf.Web.ReadHeaderTimeout,
		ReadTimeout:       time.Second * s.Conf.Web.ReadTimeoutInSec,
		IdleTimeout:       time.Second * s.Conf.Web.IdleTimeoutInSec,
		Handler:           s.Handler,
	}

	// auto reconnected db
	go s.autoReConnectDB()

	// Graceful shutdown http service.
	go func() {
		glog.Infof("server listening on :%d", s.Conf.Web.Port)
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				glog.Fatal(err)
			}
		}
	}()
	c := make(chan os.Signal, 1)

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*s.Conf.Web.ShutdownWaitInSec)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)

	glog.Infof("shutting down")
	glog.Flush()
	os.Exit(0)
}
