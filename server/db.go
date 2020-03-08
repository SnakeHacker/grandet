package server

import (
	"fmt"
	"time"

	"github.com/golang/glog"

	"github.com/SnakeHacker/grandet/common"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//DBConf ...
type DBConf struct {
	Username     string        `yaml:"username"`
	Password     string        `yaml:"password"`
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	Database     string        `yaml:"database"`
	Reset        bool          `yaml:"reset"`
	ReconnectSec time.Duration `yaml:"reconnect_sec"`
}

func (c DBConf) validate() error {
	if c.Username == "" {
		return common.ErrDBEmptyUsername
	}
	if c.Host == "" {
		return common.ErrDBEmptyHost
	}
	if c.Port == 0 {
		return common.ErrDBEmptyPort
	}
	if c.Database == "" {
		return common.ErrDBEmptyDatabase
	}

	return nil
}

// NewPostgreSQL ...
func NewPostgreSQL(conf DBConf) (db *gorm.DB, err error) {
	db, err = gorm.Open("postgres",
		fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s",
			conf.Host,
			conf.Port,
			conf.Username,
			conf.Database,
			conf.Password))

	if err != nil {
		glog.Error(err)
		return
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	return
}

func (s *Servlet) autoReConnectDB() (err error) {
	for {
		if s.DB.DB().Ping() != nil {
			s.DB, err = NewPostgreSQL(s.Conf.DB)
			if err != nil {
				glog.Error(err)
			}
		}
		time.Sleep(3000)
	}
}

// CreateTables create all tables
func (s *Servlet) CreateTables(values ...interface{}) (err error) {
	errs := s.DB.CreateTable(values...).GetErrors()
	if len(errs) != 0 {
		for _, err := range errs {
			glog.Error(err)
		}
		return
	}

	return
}

// DropTables drop all tables
func (s *Servlet) DropTables(values ...interface{}) (err error) {
	errs := s.DB.DropTableIfExists(values...).GetErrors()
	if len(errs) != 0 {
		for _, err := range errs {
			glog.Error(err)
		}
		return
	}

	return
}

// ResetTables drop and create tables
func (s *Servlet) ResetTables(values ...interface{}) (err error) {
	if err = s.DropTables(values...); err != nil {
		glog.Error(err)
		return
	}

	if err = s.CreateTables(values...); err != nil {
		glog.Error(err)
		return
	}

	return
}
