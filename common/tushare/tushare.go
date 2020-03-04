package tushare

import (
	"fmt"
	"os"

	"github.com/golang/glog"
)

// API get tushare API URL and Token
func API() (url, token string, err error) {
	url, has := os.LookupEnv("TUSHARE_API")
	if !has {
		err = fmt.Errorf("TUSHARE_API is empty.")
		glog.Error(err)
		return
	}

	token, has = os.LookupEnv("TUSHARE_TOKEN")
	if !has {
		err = fmt.Errorf("TUSHARE_TOKEN is empty.")
		glog.Error(err)
		return
	}

	return
}
