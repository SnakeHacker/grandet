package server

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/SnakeHacker/grandet/common"
	"github.com/SnakeHacker/grandet/common/utils/io"
	"github.com/golang/glog"
	_ "github.com/lib/pq"
)

const (
	CREATE_TABLE_SUFFIX = ".create.sql"
	UP_TABLE_SUFFIX     = ".up.sql"
	DOWN_TABLE_SUFFIX   = ".down.sql"
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
	SQLDir       string        `yaml:"sql_dir"`
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
func NewPostgreSQL(conf DBConf) (db *sql.DB, err error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Database,
	)
	return sql.Open("postgres", connStr)
}

func (s *Servlet) autoReConnectDB() (err error) {
	for {
		if s.DB.Ping() != nil {
			s.DB, err = NewPostgreSQL(s.Conf.DB)
			if err != nil {
				glog.Error(err)
			}
		}
		time.Sleep(s.Conf.DB.ReconnectSec * time.Second)
	}
}

func execSQLFile(db *sql.DB, sqlPath string) (err error) {
	content, err := ioutil.ReadFile(sqlPath)
	if err != nil {
		glog.Error(err)
		return
	}
	sqls := strings.Split(string(content), ";")
	for _, s := range sqls {
		if s != "" {
			if _, err = db.Exec(s); err != nil {
				glog.Errorf("%s, Err: %v", s, err)
				return
			}
		}
	}
	return
}

func batchExecSQLFileWithSuffix(db *sql.DB, sqlDir, suffix string) (err error) {
	files, err := io.GetFilesWithSuffix(sqlDir, suffix)
	if err != nil {
		glog.Error(err)
		return
	}

	for _, f := range files {
		if err = execSQLFile(db, f); err != nil {
			glog.Error(err)
			return
		}
	}
	return
}

// // InsertStocks ...
// func (s *Servlet) InsertStocks(stock tushare.StockBasicResponse) (err error) {
// 	s.DB.
// }
