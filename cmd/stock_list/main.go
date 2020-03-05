package main

import (
	"flag"

	"github.com/golang/glog"

	"github.com/SnakeHacker/grandet/common/utils/io"
	"github.com/SnakeHacker/grandet/server"
	"github.com/SnakeHacker/grandet/tushare"
)

func main() {
	configPath := flag.String("c", "conf.yaml", "config file")
	flag.Parse()
	flag.Set("logtostderr", "true")

	conf, err := server.LoadConf(*configPath)
	if err != nil {
		glog.Fatal(err)
	}

	resp, err := tushare.StockBasic()
	if err != nil {
		glog.Error(err)
		return
	}

	if conf.StorageDB {
		glog.Info("storage stocks into db")
		_, err := server.New(conf)
		if err != nil {
			glog.Error(err)
			return
		}
	}

	if conf.StorageExcel {
		glog.Info("storage stocks into xlsx")

		filename := "stock_list"
		outputDir := ""
		sheet := "stock_list"
		header := resp.Data.Fields
		data := resp.Data.Items

		err = io.SaveSimpleXlsx(filename, outputDir, sheet, header, data)
		if err != nil {
			glog.Error(err)
			return
		}
	}
}
