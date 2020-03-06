package tushare

import (
	"bytes"
	"fmt"
	"time"

	"github.com/SnakeHacker/grandet/common"
	"github.com/golang/glog"
)

// Daily
// 数据说明：交易日每天15点～16点之间。本接口是未复权行情，停牌期间不提供数据。
// 调取说明：基础积分每分钟内最多调取500次，每次5000条数据，相当于23年历史，用户获得超过5000积分正常调取无频次限制。
func (t *TuShare) Daily(tsCodes []string, startDate time.Time, endDate time.Time) (resp TushareHttpResponse, err error) {
	params, err := dailyParams(tsCodes, startDate, endDate)
	if err != nil {
		glog.Error(err)
		return
	}
	req := &TushareHttpRequest{
		APIName: API_DAILY,
		Token:   t.Token,
		Params:  params,
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

func dailyParams(tsCodes []string, startDate time.Time, endDate time.Time) (params map[string]string, err error) {
	params = make(map[string]string)

	if len(tsCodes) == 0 {
		err = common.ErrTushareDailyParamsEmpty
		glog.Error(err)
		return
	}

	var tsCodesParam bytes.Buffer

	for i, tsCode := range tsCodes {
		tsCodesParam.WriteString(tsCode)
		if i != len(tsCodes) {
			tsCodesParam.WriteString(",")
		}
	}

	params["ts_code"] = tsCodesParam.String()

	start := startDate.Format("20060102")
	params[START_DATE] = start

	end := endDate.Format("20060102")
	params[END_DATE] = end

	return
}
