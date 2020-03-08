package tushare

import (
	"fmt"

	"github.com/golang/glog"
)

// ConceptDetail ...
func (t *TuShare) ConceptDetail(tsCode string) (resp TushareHttpResponse, err error) {
	req := &TushareHttpRequest{
		APIName: API_CONCEPT_DETAIL,
		Token:   t.Token,
		Params:  map[string]string{TS_CODE: tsCode},
	}

	if err = t.API(req, &resp); err != nil {
		glog.Error(err)
		return
	}

	if resp.Code != 0 {
		err = fmt.Errorf("[%v] errCode:%v msg: %v.", tsCode, resp.Code, resp.Msg)
		glog.Error(err)
		return
	}

	return
}
