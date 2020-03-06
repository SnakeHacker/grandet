package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/golang/glog"

	"github.com/SnakeHacker/grandet/common/utils/io"
	"github.com/SnakeHacker/grandet/server"
)

func main() {
	configPath := flag.String("c", "conf.yaml", "config file")
	startDate := flag.String("start", "", "start date")
	endDate := flag.String("end", "", "end date")
	tsCode := flag.String("tscode", "", "ts code")

	flag.Parse()
	flag.Set("logtostderr", "true")

	if *startDate == "" || *endDate == "" || *tsCode == "" {
		flag.Usage()
		return
	}

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

	start, err := time.Parse("20060102", *startDate)
	if err != nil {
		glog.Error(err)
		return
	}

	end, err := time.Parse("20060102", *endDate)
	if err != nil {
		glog.Error(err)
		return
	}

	resp, err := s.Tushare.Daily([]string{*tsCode}, start, end)
	if err != nil {
		glog.Error(err)
		return
	}

	if conf.StorageDB {
		glog.Info("Inserting daily into db ...")
		if err = s.InsertDaily(resp.Data.Fields, resp.Data.Items); err != nil {
			glog.Error(err)
			return
		}
		glog.Info("Insert daily into db successfully!")
	}

	if conf.StorageExcel {
		glog.Info("Storage daily data into xlsx")

		filename := fmt.Sprintf("%s_%s_%s_daily", *tsCode, *startDate, *endDate)
		outputDir := ""
		sheet := "daily"
		header := resp.Data.Fields
		data := resp.Data.Items

		err = io.SaveSimpleXlsx(filename, outputDir, sheet, header, data)
		if err != nil {
			glog.Error(err)
			return
		}
	}

}
