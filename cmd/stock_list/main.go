package main

import (
	"flag"

	"github.com/golang/glog"

	"github.com/SnakeHacker/grandet/common/utils/io"
	"github.com/SnakeHacker/grandet/server"
)

func main() {
	configPath := flag.String("c", "conf.yaml", "config file")
	flag.Parse()
	flag.Set("logtostderr", "true")

	conf, err := server.LoadConf(*configPath)
	if err != nil {
		glog.Fatal(err)
	}

	glog.Info("Start new server...")
	s, err := server.New(conf)
	if err != nil {
		glog.Error(err)
		return
	}

	resp, err := s.Tushare.StockBasic()
	if err != nil {
		glog.Error(err)
		return
	}

	if conf.StorageDB {
		glog.Info("Inserting stocks into db ...")
		if err = s.BatchInsertStockMeta(resp.Data.Fields, resp.Data.Items); err != nil {
			glog.Error(err)
			return
		}
		glog.Info("Insert stocks into db successfully!")
	}

	if conf.StorageExcel {
		glog.Info("Storage stocks into xlsx")

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
