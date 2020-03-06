package tushare

import (
	"fmt"

	"github.com/golang/glog"
)

// StockBasic ...
func (t *TuShare) StockBasic() (resp TushareHttpResponse, err error) {
	req := TushareHttpRequest{
		APIName: API_STOCK_BASIC,
		Token:   t.Token,
	}

	if err = t.API(req, &resp); err != nil {
		glog.Error(err)
		return
	}

	if resp.Code != 0 {
		err = fmt.Errorf("[%v] %v", resp.Code, resp.Msg)
		glog.Error(err)
		return
	}

	return
}
