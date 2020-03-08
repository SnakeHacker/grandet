package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"sort"
	"time"

	"github.com/golang/glog"
	"github.com/tealeg/xlsx"

	"github.com/SnakeHacker/grandet/common/utils/io"
	"github.com/SnakeHacker/grandet/server"
)

// TODO(mickey): refactor

func main() {
	configPath := flag.String("c", "conf.yaml", "config file")
	date := flag.String("date", "", "trade date")

	flag.Parse()
	flag.Set("logtostderr", "true")

	if *date == "" {
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

	stocks, err := s.FindStocks()
	if err != nil {
		glog.Error(err)
		return
	}

	tradeDate, err := time.Parse("20060102", *date)
	if err != nil {
		glog.Error(err)
		return
	}

	tsCodes := []string{}

	industryMap := make(map[string][]server.StockDaily)

	for _, stock := range stocks {
		tsCodes = append(tsCodes, stock.TsCode)
	}

	dailys, err := s.FindDailys(tsCodes, tradeDate)
	if err != nil {
		glog.Error(err)
		return
	}

	for _, daily := range dailys {
		if daily.PctChg <= 9.9 || daily.PctChg > 11 {
			continue
		}

		conceptDetails, err := s.FindConceptDetail(daily.TsCode)
		if err != nil {
			glog.Error(err)
			return
		}

		for _, cp := range conceptDetails {
			industryMap[cp.ConceptName] = append(industryMap[cp.ConceptName], daily)
		}
	}

	err = SaveLimitUpXlsx(fmt.Sprintf("limit_up_%v", *date), "", industryMap)
	if err != nil {
		glog.Error(err)
		return
	}

}

type IndustryWrappers []IndustryWrapper

type IndustryWrapper struct {
	Industry string
	Dailys   []server.StockDaily
	Summary  float64
}

func (p IndustryWrappers) Len() int { return len(p) }

func (p IndustryWrappers) Less(i, j int) bool {
	return p[i].Summary > p[j].Summary
}

func (p IndustryWrappers) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func SaveLimitUpXlsx(filename string, dir string, data map[string][]server.StockDaily) (err error) {
	industryMap := make(map[string]IndustryWrapper)

	for industry, dailys := range data {
		industry = io.RemovePunctuation(industry)

		var sumAmount float64

		sortedDailys := server.SortedDailys(dailys)
		sort.Sort(sortedDailys)

		for _, daily := range sortedDailys {
			sumAmount += daily.Amount
		}
		industryMap[industry] = IndustryWrapper{
			Industry: industry,
			Dailys:   sortedDailys,
			Summary:  sumAmount,
		}
	}

	var industryWrappers IndustryWrappers

	for _, wrapper := range industryMap {
		industryWrappers = append(industryWrappers, wrapper)
	}

	sort.Sort(industryWrappers)

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("概念")
	if err != nil {
		glog.Error(err)
		return err
	}

	r := sheet.AddRow()
	r.SetHeightCM(1)
	cell := r.AddCell()
	cell.Value = "概念"
	cell = r.AddCell()
	cell.Value = "总成交额（千万）"

	for _, wrapper := range industryWrappers {
		r := sheet.AddRow()
		r.SetHeightCM(1)
		cell := r.AddCell()
		cell.Value = wrapper.Industry
		cell = r.AddCell()
		cell.Value = fmt.Sprintf("%5.5f", wrapper.Summary/10000)

	}

	sheet, err = file.AddSheet("详情")
	if err != nil {
		glog.Error(err)
		return err
	}
	for _, wrapper := range industryWrappers {
		r := sheet.AddRow()
		r.SetHeightCM(1)
		cell := r.AddCell()
		cell.Value = wrapper.Industry

		{
			r = sheet.AddRow()
			r.SetHeightCM(1)
			cell := r.AddCell()
			cell.Value = "股票代码"
			cell = r.AddCell()
			cell.Value = "昨收价"
			cell = r.AddCell()
			cell.Value = "开盘价"
			cell = r.AddCell()
			cell.Value = "收盘价"
			cell = r.AddCell()
			cell.Value = "涨跌额"
			cell = r.AddCell()
			cell.Value = "涨跌幅"
			cell = r.AddCell()
			cell.Value = "成交量（千手）"
			cell = r.AddCell()
			cell.Value = "成交额（千万）"
		}

		for _, daily := range wrapper.Dailys {
			r = sheet.AddRow()
			r.SetHeightCM(1)
			cell := r.AddCell()
			cell.Value = daily.TsCode
			cell = r.AddCell()
			cell.Value = fmt.Sprintf("%v", daily.PreClose)
			cell = r.AddCell()
			cell.Value = fmt.Sprintf("%v", daily.Open)
			cell = r.AddCell()
			cell.Value = fmt.Sprintf("%v", daily.Close)
			cell = r.AddCell()
			cell.Value = fmt.Sprintf("%v", daily.Change)
			cell = r.AddCell()
			cell.Value = fmt.Sprintf("%v", daily.PctChg)
			cell = r.AddCell()
			cell.Value = fmt.Sprintf("%5.5v", daily.Vol/1000)
			cell = r.AddCell()
			cell.Value = fmt.Sprintf("%5.5f", daily.Amount/10000)
		}

		r = sheet.AddRow()
		cell = r.AddCell()
		cell.Value = fmt.Sprintf("总成交额（千万）: %5.5f", wrapper.Summary/10000)
		r = sheet.AddRow()
		r = sheet.AddRow()
	}

	err = file.Save(filepath.Join(dir, fmt.Sprintf("%s.xlsx", filename)))
	if err != nil {
		glog.Error(err)
		return
	}

	return
}
