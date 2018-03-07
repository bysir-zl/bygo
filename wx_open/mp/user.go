package mp

import (
	"errors"
	"fmt"
	"github.com/bysir-zl/bygo/wx_open/util"
	"encoding/json"
	"github.com/bysir-zl/bygo/wx_open"
)

type UserInfoRsp struct {
	wx_open.WxResponse
	SubscribeTime int64  `json:"subscribe_time"` //关注时间
	Subscribe     int    `json:"subscribe"`      //是否关注
	OpenId        string `json:"openid"`
	NickName      string `json:"nickname"`
	Language      string `json:"language"`
	Sex           int    `json:"sex"`
	Province      string `json:"province"`
	City          string `json:"city"`
	Remark        string `json:"remark"` //备注
	Country       string `json:"country"`
	HeadImgURL    string `json:"headimgurl"`
	UnionId       string `json:"unionid"`
	GroupId       int    `json:"groupid"`
	TagIdList     []int  `json:"tagid_list"`
}

// 授权之后获取用户信息,包含是否关注公众号, 注意这里的accessToken不是用户的token, 而是公众号的
func GetUserInfo(access_token, openid, lang string) (r *UserInfoRsp, err error) {
	if access_token == "" {
		err = errors.New("AccessToken miss")
		return
	}
	if openid == "" {
		err = errors.New("OpenId miss")
		return
	}

	url := "https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=%s"
	rsp, err := util.Get(fmt.Sprintf(url, access_token, openid, lang))
	if err != nil {
		err = errors.New("RefreshUserAccessToken err:" + err.Error())
		return
	}
	var rs UserInfoRsp
	err = json.Unmarshal(rsp, &rs)
	if err != nil {
		return
	}
	err = rs.HasError()
	if err != nil {
		return
	}

	r = &rs
	return

}
