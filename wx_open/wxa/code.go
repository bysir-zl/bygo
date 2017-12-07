package wxa

import (
	"github.com/bysir-zl/bygo/wx_open/util"
	"encoding/json"
	"github.com/bysir-zl/bygo/wx_open"
)

// 获取体验小程序的体验二维码
func GetQrcode(accessToken string) (image []byte, err error) {
	rsp, err := util.Get(("https://api.weixin.qq.com/wxa/get_qrcode?access_token=") + accessToken)
	if err != nil {
		return
	}
	if len(rsp) > 0 && rsp[0] == '{' {
		r := wx_open.WxResponse{}
		err = json.Unmarshal(rsp, &r)
		if err != nil {
			return
		}
		err = r.Error()
		if err != nil {
			return
		}
	}

	image = rsp
	return
}

type Color struct {
	R uint `json:"r"`
	G uint `json:"g"`
	B uint `json:"b"`
}

// 获取小程序码
// 接口B：适用于需要的码数量极多，或仅临时使用的业务场景
// 注意：通过该接口生成的小程序码，永久有效，数量暂无限制。用户扫描该码进入小程序后，开发者需在对应页面获取的码中 scene 字段的值，再做处理逻辑。
// 使用如下代码可以获取到二维码中的 scene 字段的值。调试阶段可以使用开发工具的条件编译自定义参数 scene=xxxx 进行模拟，开发工具模拟时的 scene 的参数值需要进行 urlencode
func GetAppCode(accessToken string, scene, page string, width int, autoColor bool, lineColor Color) (image []byte, err error) {
	req, _ := json.Marshal(map[string]interface{}{
		"scene":      scene,
		"page":       page,
		"width":      width,
		"auto_color": autoColor,
		"line_color": lineColor,
	})
	rsp, err := util.Post(("https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=")+accessToken, req)
	if err != nil {
		return
	}
	if len(rsp) > 0 && rsp[0] == '{' {
		r := wx_open.WxResponse{}
		err = json.Unmarshal(rsp, &r)
		if err != nil {
			return
		}
		err = r.Error()
		if err != nil {
			return
		}
	}

	image = rsp
	return
}

type Category struct {
	FirstClass  string `json:"first_class"`
	SecondClass string `json:"second_class"`
	FirstId     int    `json:"first_id"`
	SecondId    int    `json:"second_id"`
}
type CategoryRsp struct {
	wx_open.WxResponse
	CategoryList []Category `json:"category_list"`
}

// 获取授权小程序帐号的可选类目
func GetCategory(accessToken string) (categoryList []Category, err error) {
	rsp, err := util.Get(("https://api.weixin.qq.com/wxa/get_category?access_token=") + accessToken)
	if err != nil {
		return
	}
	r := CategoryRsp{}
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return
	}
	err = r.Error()
	if err != nil {
		return
	}

	categoryList = r.CategoryList
	return
}
