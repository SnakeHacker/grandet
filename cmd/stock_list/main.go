package main

import (
	"flag"

	"github.com/golang/glog"

	"github.com/SnakeHacker/grandet/common"
	"github.com/SnakeHacker/grandet/tushare"
)

func main() {
	configPath := flag.String("c", "conf.yaml", "config file")
	flag.Parse()
	flag.Set("logtostderr", "true")

	conf, err := common.LoadConf(*configPath)
	if err != nil {
		glog.Fatal(err)
	}

	resp, err := tushare.StockBasic()
	if err != nil {
		glog.Error(err)
		return
	}

	if conf.StorageDB {

	}

	if conf.StorageExcel {
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
}
