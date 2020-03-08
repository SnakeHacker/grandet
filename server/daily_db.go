package server

import (
	"bytes"
	"fmt"
	"time"

	"github.com/SnakeHacker/grandet/common"
	"github.com/golang/glog"
)

const (
	STOCK_DAILY_TSCODE     = "ts_code"
	STOCK_DAILY_TRADE_DATE = "trade_date"
	STOCK_DAILY_OPEN       = "open"
	STOCK_DAILY_HIGH       = "high"
	STOCK_DAILY_LOW        = "low"
	STOCK_DAILY_CLOSE      = "close"
	STOCK_DAILY_PRE_CLOSE  = "pre_close"
	STOCK_DAILY_CHANGE     = "change"
	STOCK_DAILY_PCT_CHG    = "pct_chg"
	STOCK_DAILY_VOL        = "vol"
	STOCK_DAILY_AMOUNT     = "amount"
)

func stockDailyHeaderIdx(headers []string) (
	headerTsCodeIdx,
	headerTradeDateiDx,
	headerOpenIdx,
	headerHighIdx,
	headerLowIdx,
	headerCloseIdx,
	headerPreCloseIdx,
	headerChangeIdx,
	headerPctChgIdx,
	headerVolIdx,
	headerAmountIdx int,
	err error) {

	for i, header := range headers {
		switch header {
		case STOCK_META_TSCODE:
			headerTsCodeIdx = i
		case STOCK_DAILY_TRADE_DATE:
			headerTradeDateiDx = i
		case STOCK_DAILY_OPEN:
			headerOpenIdx = i
		case STOCK_DAILY_HIGH:
			headerHighIdx = i
		case STOCK_DAILY_LOW:
			headerLowIdx = i
		case STOCK_DAILY_CLOSE:
			headerCloseIdx = i
		case STOCK_DAILY_PRE_CLOSE:
			headerPreCloseIdx = i
		case STOCK_DAILY_CHANGE:
			headerChangeIdx = i
		case STOCK_DAILY_PCT_CHG:
			headerPctChgIdx = i
		case STOCK_DAILY_VOL:
			headerVolIdx = i
		case STOCK_DAILY_AMOUNT:
			headerAmountIdx = i
		default:
			err = fmt.Errorf("%v, %v", common.ErrTushareDailyFieldsUnknown, header)
			glog.Error(err)
			return
		}
	}
	return
}

// BatchInsertDaily ...
func (s *Servlet) BatchInsertDaily(headers []string, items [][]interface{}) (err error) {
	var dailys []*StockDaily
	headerTsCodeIdx,
		headerTradeDateiDx,
		headerOpenIdx,
		headerHighIdx,
		headerLowIdx,
		headerCloseIdx,
		headerPreCloseIdx,
		headerChangeIdx,
		headerPctChgIdx,
		headerVolIdx,
		headerAmountIdx,
		err := stockDailyHeaderIdx(headers)
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

		dailys = append(dailys, &StockDaily{
			TsCode:    fields[headerTsCodeIdx].(string),
			TradeDate: fields[headerTradeDateiDx].(string),
			Open:      fields[headerOpenIdx].(float64),
			High:      fields[headerHighIdx].(float64),
			Low:       fields[headerLowIdx].(float64),
			Close:     fields[headerCloseIdx].(float64),
			PreClose:  fields[headerPreCloseIdx].(float64),
			Change:    fields[headerChangeIdx].(float64),
			PctChg:    fields[headerPctChgIdx].(float64),
			Vol:       fields[headerVolIdx].(float64),
			Amount:    fields[headerAmountIdx].(float64),
		})
	}

	if len(dailys) == 0 {
		glog.Warningf("daily data is empty.")
		return
	}

	stt := `INSERT INTO daily (
		ts_code, 
		trade_date, 
		open, 
		high, 
		low, 
		close, 
		pre_close, 
		change, 
		pct_chg, 
		vol, 
		amount) VALUES `
	var buf bytes.Buffer
	buf.WriteString(stt)

	for i, daily := range dailys {
		buf.WriteString(fmt.Sprintf("('%s', '%s', %f, %f, %f, %f, %f, %f, %f, %f, %f)",
			daily.TsCode,
			daily.TradeDate,
			daily.Open,
			daily.High,
			daily.Low,
			daily.Close,
			daily.PreClose,
			daily.Change,
			daily.PctChg,
			daily.Vol,
			daily.Amount))
		if i != len(dailys)-1 {
			buf.WriteString(",")
		} else {
			buf.WriteString(";")
		}
	}

	err = s.DB.Exec(buf.String()).Error
	if err != nil {
		glog.Errorf("%v. SQL: %v", err, buf.String())
		return
	}

	return
}

// FindDaily ...
func (s *Servlet) FindDailys(tsCodes []string, date time.Time) (dailys []StockDaily, err error) {
	err = s.DB.Where(
		"ts_code IN (?) AND trade_date = ?",
		tsCodes, date.Format("20060102")).
		Find(&dailys).Error

	if err != nil {
		glog.Error(err)
		return
	}

	return
}
