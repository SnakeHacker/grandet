package server

import (
	"bytes"
	"fmt"

	"github.com/SnakeHacker/grandet/common"
	"github.com/golang/glog"
)

const (
	STOCK_META_TSCODE   = "ts_code"
	STOCK_META_Symbol   = "symbol"
	STOCK_META_Name     = "name"
	STOCK_META_Area     = "area"
	STOCK_META_Industry = "industry"
	STOCK_META_Market   = "market"
	STOCK_META_ListDate = "list_date"
)

func stockMetaHeaderIdx(headers []string) (
	headerTsCodeIdx,
	headerSymbolIdx,
	headerNameIdx,
	headerAreaIdx,
	headerIndustryIdx,
	headerMarketIdx,
	headerListDateIdx int,
	err error) {

	for i, header := range headers {
		switch header {
		case STOCK_META_TSCODE:
			headerTsCodeIdx = i
		case STOCK_META_Symbol:
			headerSymbolIdx = i
		case STOCK_META_Name:
			headerNameIdx = i
		case STOCK_META_Area:
			headerAreaIdx = i
		case STOCK_META_Industry:
			headerIndustryIdx = i
		case STOCK_META_Market:
			headerMarketIdx = i
		case STOCK_META_ListDate:
			headerListDateIdx = i
		default:
			err = common.ErrTushareStockFieldsUnknown
			glog.Error(err)
			return
		}
	}
	return
}

// BatchInsertStockMeta ...
func (s *Servlet) BatchInsertStockMeta(headers []string, items [][]interface{}) (err error) {
	var stocks []*StockMeta
	headerTsCodeIdx,
		headerSymbolIdx,
		headerNameIdx,
		headerAreaIdx,
		headerIndustryIdx,
		headerMarketIdx,
		headerListDateIdx,
		err := stockMetaHeaderIdx(headers)
	if err != nil {
		glog.Error(err)
		return
	}

	for _, fields := range items {
		if len(fields) != len(headers) {
			err = common.ErrTushareStockFieldsLen
			glog.Error(err)
			return
		}

		stocks = append(stocks, &StockMeta{
			TsCode:   fields[headerTsCodeIdx].(string),
			Symbol:   fields[headerSymbolIdx].(string),
			Name:     fields[headerNameIdx].(string),
			Area:     fields[headerAreaIdx].(string),
			Industry: fields[headerIndustryIdx].(string),
			Market:   fields[headerMarketIdx].(string),
			ListDate: fields[headerListDateIdx].(string),
		})
	}

	stt := `INSERT INTO stocks (ts_code, symbol, name, area, industry, market, list_date) VALUES `
	var buf bytes.Buffer
	buf.WriteString(stt)
	for i, stockMeta := range stocks {
		buf.WriteString(fmt.Sprintf("('%s', '%s', '%s', '%s', '%s', '%s', '%s')",
			stockMeta.TsCode,
			stockMeta.Symbol,
			stockMeta.Name,
			stockMeta.Area,
			stockMeta.Industry,
			stockMeta.Market,
			stockMeta.ListDate))
		if i != len(stocks)-1 {
			buf.WriteString(",")
		} else {
			buf.WriteString(";")
		}
	}

	err = s.DB.Exec(buf.String()).Error
	if err != nil {
		glog.Error(err)
		return
	}

	return
}

// FindStocks ...
func (s *Servlet) FindStocks() (stocks []StockMeta, err error) {
	err = s.DB.Find(&stocks).Error
	if err != nil {
		glog.Error(err)
		return
	}

	return
}
