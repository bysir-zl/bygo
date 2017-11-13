package core

import (
	"encoding/json"
	"errors"
	"github.com/bysir-zl/bjson"
	"github.com/bysir-zl/bygo/util"
	"github.com/bysir-zl/bygo/util/encoder"
	"github.com/bysir-zl/bygo/util/http_util"
	"net/url"
	"strconv"
	"strings"
)

const (
	apiHost       = "https://payh5.bbnpay.com"
	apiPlaceOrder = apiHost + "/cpapi/place_order.php"

	BbnResponse_Success = "SUCCESS"
)

// H5网页收银台 第三方支付
type BbnPay struct {
	BbnPayConfig
}

type BbnPayConfig struct {
	AppId, Key string
}

type BbnPayPlaceOrder struct {
	AppId     string `json:"appid"`
	GoodsId   int    `json:"goodsid"`   // 应用中的商品编号
	PcorderId string `json:"pcorderid"` // 商户生成的订单号，需要保证系统唯一
	Money     int    `json:"money"`     // 支付金额
	Currency  string `json:"currency"`  // 货币类型以及单位： CHY – 人民币（单位：分）
	PcuserId  string `json:"pcuserid"`  // 用户在商户应用的唯一标识，建议为用户帐号。
	NotifyUrl string `json:"notifyurl"` // 商户服务端接收支付结果通知的地址
	GoodsName string `json:"goodsname"`
}

type BbnPayNotify struct {
	Appid      string `json:"appid,omitempty"`
	Goodsid    string `json:"goodsid,omitempty"`
	TransType  int8   `json:"transtype,omitempty"`    // 交易类型：0–支付交易；
	Cporderid  string `json:"cporderid,omitempty"`    // 商户订单号
	Transid    string `json:"transid,omitempty"`      // 交易流水号
	Pcuserid   string `json:"pcuserid,omitempty"`     // 用户在商户应用 的唯一标识
	Feetype    int8   `json:"feetype,omitempty"`      // 计费方式
	Money      int    `json:"money,omitempty"`        // 本次交易的金额（请务必严格校验商品金额与交易的金额是否一致）
	FactMoney  int    `json:"fact_money,omitempty"`   // 实际付款金额
	Result     int    `json:"result,omitempty"`       // 交易结果： 1–支付成功 2–支付失败
	Paytype    string `json:"paytype,omitempty"`      // 支付方式
	Currency   string `json:"currency,omitempty"`     // 货币类型
	Transtime  string `json:"transtime,omitempty"`    // 交易完成时间
	PcPrivInfo string `json:"pc_priv_info,omitempty"` // 商户私有信息
	err        error
}

func (p BbnPayNotify) IsError() error {
	return p.err
}

// 下单
func (p *BbnPay) PlaceOrder(i *BbnPayPlaceOrder) (transId string, err error) {
	i.AppId = p.AppId
	if i.Currency == "" {
		i.Currency = "CHY"
	}

	signKv := util.OrderKV{}
	signKv.Add("appid", i.AppId)
	signKv.Add("currency", i.Currency)
	signKv.Add("pcuserid", i.PcuserId)
	signKv.Add("notifyurl", i.NotifyUrl)
	signKv.Add("pcorderid", i.PcorderId)
	signKv.Add("goodsid", strconv.Itoa(i.GoodsId))
	signKv.Add("money", strconv.Itoa(i.Money))
	signKv.Add("goodsname", i.GoodsName)
	signKv.Sort()
	// 签名
	signStr := signKv.EncodeStringWithoutEscape() + "&key=" + p.Key
	sign := encoder.Md5String(signStr)

	ps := util.OrderKV{}
	transdataByte, _ := json.Marshal(i)
	ps.Add("transdata", string(transdataByte))
	ps.Add("sign", sign)
	ps.Add("signtype", "MD5")

	_, rsp, err := http_util.Post(apiPlaceOrder, ps, nil)
	if err != nil {
		return
	}
	rsp, _ = url.QueryUnescape(rsp)
	transDataRsp := strings.Split(rsp, "=")[1]
	transDataRsp = strings.Split(transDataRsp, "&")[0]
	bj, err := bjson.New([]byte(transDataRsp))
	if err != nil {
		return
	}
	if bj.Pos("code").Int() != 200 {
		err = errors.New(bj.Pos("errmsg").String())
		return
	}

	transId = bj.Pos("transid").String()

	return
}

func (p *BbnPay) Notify(data, sign string) (resp BbnPayNotify, err error) {
	err = json.Unmarshal([]byte(data), &resp)
	if err != nil {
		return
	}

	signKv := util.OrderKV{}
	signKv.Add("transtype", strconv.Itoa(int(resp.TransType)))
	signKv.Add("cporderid", resp.Cporderid)
	signKv.Add("transid", resp.Transid)
	signKv.Add("pcuserid", resp.Pcuserid)
	signKv.Add("appid", resp.Appid)
	signKv.Add("goodsid", resp.Goodsid)
	signKv.Add("feetype", strconv.Itoa(int(resp.Feetype)))
	signKv.Add("money", strconv.Itoa(resp.Money))
	signKv.Add("fact_money", strconv.Itoa(resp.FactMoney))
	signKv.Add("currency", resp.Currency)
	signKv.Add("result", strconv.Itoa(resp.Result))
	signKv.Add("transtime", resp.Transtime)
	signKv.Add("pc_priv_info", resp.PcPrivInfo)
	signKv.Add("paytype", resp.Paytype)
	signKv.Sort()
	// 签名
	signStr := signKv.EncodeStringWithoutEscape() + "&key=" + p.Key
	signT := encoder.Md5String(signStr)

	if sign != signT {
		resp.err = errors.New("sign error")
		err = errors.New("sign error")
		return
	}

	return
}

func NewBbnPay(config BbnPayConfig) *BbnPay {
	b := &BbnPay{
		BbnPayConfig: config,
	}

	return b
}
