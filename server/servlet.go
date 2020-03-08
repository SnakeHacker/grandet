package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/SnakeHacker/grandet/tushare"

	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
)

// Servlet ...
type Servlet struct {
	Handler http.Handler
	Conf    Conf
	DB      *gorm.DB
	Tushare *tushare.TuShare
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

	// Init tushare
	s.Tushare, err = tushare.New()
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	// drop tables and create tables
	if conf.DB.Reset {
		if err = s.ResetTables(&StockMeta{}, &StockDaily{}, &ConceptDetail{}); err != nil {
			glog.Error(err)
			return
		}

		glog.Info("Reset db tables successfully!")
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
