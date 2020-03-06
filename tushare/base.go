package tushare

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
)

const (
	API_STOCK_BASIC = "stock_basic"
	API_DAILY       = "daily"
	TS_CODE         = "ts_code"
	START_DATE      = "start_date"
	END_DATE        = "end_date"
)

// TushareHttpRequest ...
// Reference: https://tushare.pro/document/1?doc_id=130
type TushareHttpRequest struct {
	APIName string            `json:"api_name"`
	Token   string            `json:"token"`
	Params  map[string]string `json:"params,omitempty"`
	Fields  string            `json:"fields,omitempty"`
}

// StockBasicResponse ...
// Reference: https://tushare.pro/document/2?doc_id=25
type TushareHttpResponse struct {
	RequestID string      `json:"request_id"`
	Code      int         `json:"code"`
	Msg       string      `json:"mgs"`
	Data      TushareData `json:"data"`
}

// TushareData ...
type TushareData struct {
	Fields  []string        `json:"fields"`
	Items   [][]interface{} `json:"items"`
	HasMore bool            `json:"has_more"`
}

// API ...
func (t *TuShare) API(req interface{}, resp interface{}) (err error) {
	reqBodyJSON, err := json.Marshal(req)
	if err != nil {
		glog.Error(err)
		return
	}

	request, err := http.NewRequest("POST", t.URL, bytes.NewReader(reqBodyJSON))
	if err != nil {
		glog.Error(err)
		return
	}

	request.Header.Set("Content-Type", "application/json;charset=UTF-8")

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		glog.Error(err)
		return
	}
	defer response.Body.Close()

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		glog.Error(err)
		return
	}

	err = json.Unmarshal(responseBytes, resp)
	if err != nil {
		glog.Error(err)
		return
	}

	return
}
