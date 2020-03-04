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

	filename := "stock_list"
	outputDir := ""
	sheet := "stock_list"
	header := resp.Data.Fields
	data := resp.Data.Items

	err = common.SaveSimpleXlsx(filename, outputDir, sheet, header, data)
	if err != nil {
		glog.Error(err)
		return
	}
}
