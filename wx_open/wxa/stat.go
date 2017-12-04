package wxa

import (
	"encoding/json"
	"github.com/bysir-zl/bygo/wx_open/util"
	"github.com/bysir-zl/bygo/wx_open"
	"github.com/lunny/log"
	"time"
)

// key	标识，0开始，0表示当月，1表示1月后，key取值分别是：0,1
// value	key对应日期的新增用户数/活跃用户数（key=0时）或留存用户数（k>0时）
type VisitUv struct {
	Key   int `json:"key"`
	Value int `json:"value"`
}

type MonthRetainRsp struct {
	wx_open.WxResponse
	RefDate    string    `json:"ref_date"`
	VisitUvNew []VisitUv `json:"visit_uv_new"`
	VisitUv    []VisitUv `json:"visit_uv"`
}

// 月留存
// month: "201701"
// 测试month在201712时微信会报错, 不知道是不是微信的bug,还是不支持查当月, 还是不支持查月初, 因为今天是20171201
func GetMonthRetain(accessToken string, month string) (data MonthRetainRsp, err error) {
	beginDate, err := time.Parse("200601", month)
	if err != nil {
		return
	}
	endDate := beginDate.AddDate(0, 1, -1)

	req, _ := json.Marshal(map[string]interface{}{
		"begin_date": beginDate.Format("20060102"),
		"end_date":   endDate.Format("20060102"),
	})
	rsp, err := util.Post(("https://api.weixin.qq.com/datacube/getweanalysisappidmonthlyretaininfo?access_token=")+accessToken, req)
	if err != nil {
		return
	}

	r := MonthRetainRsp{}
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return
	}
	err = r.Error()
	if err != nil {
		return
	}

	return r, nil
}

type (
	MonthVisit struct {
		VisitUvNew      int `json:"visit_uv_new"`
		StayTimeSession int `json:"stay_time_session"`
		VisitDepth      int `json:"visit_depth"`
		RefDate         string  `json:"ref_date"`
		SessionCnt      int `json:"session_cnt"`
		VisitPv         int `json:"visit_pv"`
		VisitUv         int `json:"visit_uv"`
	}

	MonthVisitRsp struct {
		wx_open.WxResponse
		List []MonthVisit `json:"list"`
	}
)

// 月访问
// month: "201701"
// 测试month在201712时微信会报错, 不知道是不是微信的bug,还是不支持查当月, 还是不支持查月初, 因为今天是20171201
func GetMonthVisit(accessToken string, month string) (data MonthVisitRsp, err error) {
	beginDate, err := time.Parse("200601", month)
	if err != nil {
		return
	}
	endDate := beginDate.AddDate(0, 1, -1)

	req, _ := json.Marshal(map[string]interface{}{
		"begin_date": beginDate.Format("20060102"),
		"end_date":   endDate.Format("20060102"),
	})
	rsp, err := util.Post(("https://api.weixin.qq.com/datacube/getweanalysisappidmonthlyvisittrend?access_token=")+accessToken, req)
	if err != nil {
		return
	}
	log.Info(string(rsp))

	r := MonthVisitRsp{}
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return
	}
	err = r.Error()
	if err != nil {
		return
	}

	return r, nil
}
