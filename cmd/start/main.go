package main

import (
	"flag"
	"time"

	"github.com/golang/glog"

	"github.com/SnakeHacker/grandet/server"
)

func main() {
	configPath := flag.String("c", "conf.yaml", "config file")
	startDate := flag.String("start", "", "start date")
	endDate := flag.String("end", "", "end date")

	flag.Parse()
	flag.Set("logtostderr", "true")

	if *startDate == "" || *endDate == "" {
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

	resp, err := s.Tushare.StockBasic()
	if err != nil {
		glog.Error(err)
		return
	}

	glog.Info("Inserting stocks into db ...")
	if err = s.BatchInsertStockMeta(resp.Data.Fields, resp.Data.Items); err != nil {
		glog.Error(err)
		return
	}
	glog.Info("Insert stocks into db successfully!")

	stocks, err := s.FindStocks()
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

	for _, stock := range stocks {
		resp, err = s.Tushare.Daily([]string{stock.TsCode}, start, end)
		if err != nil {
			glog.Error(err)
			return
		}

		glog.Infof("Inserting [%v] daily data into db ...", stock.TsCode)
		if err = s.BatchInsertDaily(resp.Data.Fields, resp.Data.Items); err != nil {
			glog.Error(err)
			return
		}

		resp, err = s.Tushare.ConceptDetail(stock.TsCode)
		if err != nil {
			glog.Error(err)
			return
		}

		glog.Infof("Inserting [%v] concept detail into db ...", stock.TsCode)
		if err = s.InsertConceptDetail(resp.Data.Fields, resp.Data.Items); err != nil {
			glog.Error(err)
			return
		}

		time.Sleep(time.Millisecond * 800)
	}

}
