package main

import (
	"flag"

	"github.com/golang/glog"

	"github.com/SnakeHacker/grandet/common"
	"github.com/SnakeHacker/grandet/common/tushare"
)

func main() {
	flag.Parse()
	flag.Set("logtostderr", "true")

	resp, err := tushare.StockBasic()
	if err != nil {
		glog.Error(err)
		return
	}

	err = common.SaveSimpleXlsx("stock_list", "", "stock_list", resp.Data.Fields, resp.Data.Items)
	if err != nil {
		glog.Error(err)
		return
	}
}
