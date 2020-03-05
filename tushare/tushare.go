package tushare

import (
	"os"

	"github.com/SnakeHacker/grandet/common"
	"github.com/golang/glog"
)

// API get tushare API URL and Token
func API() (url, token string, err error) {
	url, has := os.LookupEnv("TUSHARE_API")
	if !has {
		err = common.ErrTushareURL
		glog.Error(err)
		return
	}

	token, has = os.LookupEnv("TUSHARE_TOKEN")
	if !has {
		err = common.ErrTushareToken
		glog.Error(err)
		return
	}

	return
}
