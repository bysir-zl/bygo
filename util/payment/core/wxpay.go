package core

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/cxuhua/xweb"
	"html/template"
	"io"
	"reflect"
	"strings"
)

/**
 * TODO: 修改这里配置为您自己申请的商户信息
 * 微信公众号信息配置
 *
 * AppId：绑定支付的APPID（必须配置，开户邮件中可查看）
 *
 * MchId：商户号（必须配置，开户邮件中可查看）
 *
 * MchKey：商户支付密钥，参考开户邮件设置（必须配置，登录商户平台自行设置）
 * 设置地址：https://pay.weixin.qq.com/index.php/account/api_cert
 *
 * AppSecret：公众帐号secert（仅JSAPI支付的时候需要配置， 登录公众平台，进入开发者中心可设置），
 * 获取地址：https://mp.weixin.qq.com/advanced/advanced?action=dev&t=advanced/dev&token=2005451881&lang=zh_CN
 * @var string
 */

type WXKeyConfig struct {
	APP_ID     string
	APP_SECRET string
	MCH_ID     string
	MCH_KEY    string
	CRT_PATH   string
	KEY_PATH   string
	CA_PATH    string
	TLSConfig  *tls.Config
}

type WxPay struct {
	config WXKeyConfig
}

func NewWxPay(conf WXKeyConfig) (*WxPay) {
	if conf.CA_PATH != "" && conf.CRT_PATH != "" && conf.KEY_PATH != "" {
		conf.TLSConfig = xweb.MustLoadTLSFileConfig(conf.CA_PATH, conf.CRT_PATH, conf.KEY_PATH)
	}

	return &WxPay{
		config: conf,
	}
}

func (p *WxPay) WXSign(v interface{}) string {
	http := WXParseSignFields(v)
	str := http.RawEncode() + "&key=" + p.config.MCH_KEY
	return strings.ToUpper(xweb.MD5String(str))
}

func WXParseSignFields(src interface{}) xweb.HTTPValues {
	values := xweb.NewHTTPValues()
	t := reflect.TypeOf(src)
	v := reflect.ValueOf(src)
	for i := 0; i < t.NumField(); i++ {
		tf := t.Field(i)
		if tf.Tag.Get("sign") != "true" {
			continue
		}
		tv := v.Field(i)
		if !tv.IsValid() {
			continue
		}
		sv := fmt.Sprintf("%v", tv.Interface())
		if sv == "" {
			continue
		}
		name := ""
		if xn := tf.Tag.Get("xml"); xn != "" {
			name = strings.Split(xn, ",")[0]
		} else if xn = tf.Tag.Get("json"); xn != "" {
			name = strings.Split(xn, ",")[0]
		} else {
			continue
		}
		values.Add(name, sv)
	}
	return values
}

//支付类型
const (
	TRADE_TYPE_JSAPI  = "JSAPI"
	TRADE_TYPE_NATIVE = "NATIVE"
	TRADE_TYPE_APP    = "APP"
)

//转账校验
const (
	NO_CHECK     = "NO_CHECK"     //不校验真实姓名
	FORCE_CHECK  = "FORCE_CHECK"  //强制校验
	OPTION_CHECK = "OPTION_CHECK" //有则校验
)

//返回字符串
const (
	FAIL       = "FAIL"
	NOTPAY     = "NOTPAY"
	SUCCESS    = "SUCCESS"
	USERPAYING = "USERPAYING"
)

//主机地址
const (
	//公众号
	WX_API_HOST = "https://api.weixin.qq.com"
	//支付
	WX_PAY_HOST = "https://api.mch.weixin.qq.com"
)

//应用授权作用域
const (
	WX_BASE_SCOPE = "snsapi_base"     //只能获得openid,不需要用户确认
	WX_INFO_SCOPE = "snsapi_userinfo" //需要用户确认,能够获得用户信息
)

type WXError struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (this WXError) Error() error {
	if this.ErrCode == 0 {
		return nil
	}
	return errors.New(fmt.Sprintf("ERROR:%d,%s", this.ErrCode, this.ErrMsg))
}

//二维码生产
type WXQRCodeCreateRequest struct {
	ActionName string `json:"action_name"`
	ActionInfo struct {
		Scene struct {
			SceneStr string `json:"scene_str"`
		} `json:"scene"`
	} `json:"action_info"`
}

func (this WXQRCodeCreateRequest) ToReader() (io.Reader, error) {
	data, err := json.Marshal(this)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), nil
}

type WXQRCodeCreateResponse struct {
	WXError
	Ticket        string `json:"ticket"`
	ExpireSeconds int    `json:"expire_seconds"`
	URL           string `json:"url"`
}

func (this WXQRCodeCreateRequest) Post(token string, info string) (WXQRCodeCreateResponse, error) {
	ret := WXQRCodeCreateResponse{}
	this.ActionName = "QR_LIMIT_STR_SCENE"
	this.ActionInfo.Scene.SceneStr = info
	http := xweb.NewHTTPClient(WX_API_HOST)
	body, err := this.ToReader()
	if err != nil {
		return ret, err
	}
	res, err := http.Post("/cgi-bin/qrcode/create?access_token="+token, "application/json", body)
	if err != nil {
		return ret, err
	}
	if err := res.ToJson(&ret); err != nil {
		return ret, err
	}
	if ret.ErrCode != 0 {
		return ret, errors.New(ret.ErrMsg)
	}
	return ret, nil
}

//发送消息
const (
	MSG_TEXT  = "text"
	MSG_IMAGE = "image"
	MSG_NEW   = "news"
)

type CustomMsg struct {
	ToUser  string        `json:"touser"`
	MsgType string        `json:"msgtype"`
	Text    CustomMsgText `json:"text"`
}

func (this CustomMsg) ToReader() (io.Reader, error) {
	this.MsgType = MSG_TEXT
	if this.ToUser == "" {
		return nil, errors.New("ToUser empty")
	}
	if this.Text.Content == "" {
		return nil, errors.New("msg content null")
	}
	data, err := json.Marshal(this)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), nil
}

type CustomMsgText struct {
	Content string `json:"content"`
}

type WXGetAccessTokenResponse struct {
	WXError
	AccessToken string `json:"access_token"`
	Expires     int    `json:"expires_in"`
}

//微信转账
type WXTransfersRequest struct {
	AppId          string `xml:"mch_appid" sign:"true"`
	MchId          string `xml:"mchid" sign:"true"`
	NonceStr       string `xml:"nonce_str" sign:"true"`
	Sign           string `xml:"sign" sign:"false"`
	PartnerTradeNo string `xml:"partner_trade_no" sign:"true"`
	OpenId         string `xml:"openid" sign:"true"`
	CheckName      string `xml:"check_name" sign:"true"`
	Amount         int    `xml:"amount" sign:"true"`
	Desc           string `xml:"desc" sign:"true"`
	SpbillCreateIp string `xml:"spbill_create_ip" sign:"true"`
}

func (r WXTransfersRequest) ToXml() string {
	data, err := xml.Marshal(r)
	if err != nil {
		panic(err)
	}
	return string(data)
}

type WXTransfersResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	ResultCode string `xml:"result_code"`
	ErrCodeDes string `xml:"err_code_des"`
}

func (p *WxPay) WxTransfersRequest(q WXTransfersRequest) (WXTransfersResponse, error) {
	ret := WXTransfersResponse{}
	q.PartnerTradeNo = xweb.GenId()
	q.NonceStr = RandStr()
	q.MchId = p.config.MCH_ID
	q.AppId = p.config.APP_ID
	if q.Amount <= 0 {
		return ret, errors.New("Amount error")
	}
	if q.SpbillCreateIp == "" {
		return ret, errors.New("SpbillCreateIp miss")
	}
	if q.Desc == "" {
		return ret, errors.New("Desc miss")
	}
	if q.CheckName == "" {
		return ret, errors.New("CheckName miss")
	}
	vs := WXParseSignFields(q)
	q.Sign = strings.ToUpper(vs.MD5Sign(p.config.MCH_KEY))
	body := strings.NewReader(q.ToXml())
	http := xweb.NewHTTPClient(WX_PAY_HOST, p.config.TLSConfig)
	res, err := http.Post("/mmpaymkttransfers/promotion/transfers", "application/xml", body)
	if err != nil {
		return ret, err
	}
	if err := res.ToXml(&ret); err != nil {
		return ret, err
	}
	if ret.ReturnCode != SUCCESS {
		return ret, errors.New(ret.ReturnMsg)
	}
	if ret.ResultCode != SUCCESS {
		return ret, errors.New(ret.ErrCodeDes)
	}
	return ret, nil
}

//微信红包发送
type WXRedPackageRequest struct {
	MchBillno   string `xml:"mch_billno" sign:"true"`
	NonceStr    string `xml:"nonce_str" sign:"true"`
	MchId       string `xml:"mch_id" sign:"true"`
	AppId       string `xml:"wxappid" sign:"true"`
	SendName    string `xml:"send_name" sign:"true"`
	ReOpenId    string `xml:"re_openid" sign:"true"`
	TotalAmount int    `xml:"total_amount" sign:"true"`
	TotalNum    int    `xml:"total_num" sign:"true"`
	Wishing     string `xml:"wishing" sign:"true"`
	ClientIp    string `xml:"client_ip" sign:"true"`
	ActName     string `xml:"act_name" sign:"true"`
	Remark      string `xml:"remark" sign:"true"`
	Sign        string `xml:"sign" sign:"false"`
}

func (this WXRedPackageRequest) ToXml() string {
	data, err := xml.Marshal(this)
	if err != nil {
		panic(err)
	}
	return string(data)
}

type WXRedPackageResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	ResultCode string `xml:"result_code"`
	ErrCodeDes string `xml:"err_code_des"`
}

func (p *WxPay) RedPackageRequest(q WXRedPackageRequest) (WXRedPackageResponse, error) {
	ret := WXRedPackageResponse{}
	q.MchBillno = xweb.GenId()
	q.NonceStr = RandStr()
	q.MchId = p.config.MCH_ID
	q.AppId = p.config.APP_ID
	if q.SendName == "" {
		return ret, errors.New("SendName miss")
	}
	if q.ReOpenId == "" {
		return ret, errors.New("ReOpenId miss")
	}
	if q.TotalAmount <= 0 {
		return ret, errors.New("TotalAmount error")
	}
	if q.TotalNum <= 0 {
		return ret, errors.New("TotalNum error")
	}
	if q.ClientIp == "" {
		return ret, errors.New("ClientIp miss")
	}
	if q.ActName == "" {
		return ret, errors.New("ActName miss")
	}
	if q.Remark == "" {
		return ret, errors.New("Remark miss")
	}
	if q.Wishing == "" {
		return ret, errors.New("Wishing miss")
	}
	vs := WXParseSignFields(q)
	q.Sign = strings.ToUpper(vs.MD5Sign(p.config.MCH_KEY))
	body := strings.NewReader(q.ToXml())
	http := xweb.NewHTTPClient(WX_PAY_HOST, p.config.TLSConfig)
	res, err := http.Post("/mmpaymkttransfers/sendredpack", "application/xml", body)
	if err != nil {
		return ret, err
	}
	if err := res.ToXml(&ret); err != nil {
		return ret, err
	}
	if ret.ReturnCode != SUCCESS {
		return ret, errors.New("1," + ret.ReturnCode + ":" + ret.ReturnMsg)
	}
	if ret.ResultCode != SUCCESS {
		return ret, errors.New("2," + ret.ResultCode + ":" + ret.ErrCodeDes)
	}
	return ret, nil
}

//{
//"openId":"oUzzq0GSPdp2XEmi7l5g1y8jnGX8",
//"nickName":"芒果玛奇朵",
//"gender":1,
//"language":"zh_CN",
//"city":"Panzhihua",
//"province":"Sichuan",
//"country":"CN",
//"avatarUrl":"http://wx.qlogo.cn/mmopen/vi_3"
//}

type EncryptedInfo struct {
	OpenId    string `json:"openId"`
	UnionId   string `json:"unionId"`
	NickName  string `json:"nickName"`
	Gender    int    `json:"gender"`
	Language  string `json:"language"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Country   string `json:"country"`
	AvatarUrl string `json:"avatarUrl"`
}

// 解密用户信息
func WXAppDecodeEncryptedData(skey string, siv string, sdata string) (EncryptedInfo, error) {
	info := EncryptedInfo{}
	data, err := base64.StdEncoding.DecodeString(sdata)
	if err != nil {
		return info, err
	}
	iv, err := base64.StdEncoding.DecodeString(siv)
	if err != nil {
		return info, err
	}
	key, err := base64.StdEncoding.DecodeString(skey)
	if err != nil {
		return info, err
	}
	aes, err := xweb.NewAESChpher(key)
	if err != nil {
		return info, err
	}
	idata, err := xweb.AesDecryptWithIV(aes, data, iv)
	if err != nil {
		return info, err
	}
	if err := json.Unmarshal(idata, &info); err != nil {
		return info, err
	}
	return info, nil
}

//https://api.mch.weixin.qq.com/secapi/pay/refund
//微信退款发起请求
type WXRefundRequest struct {
	XMLName       struct{} `xml:"xml"`
	AppId         string   `xml:"appid,omitempty" sign:"true"`
	MchId         string   `xml:"mch_id,omitempty" sign:"true"`
	NonceStr      string   `xml:"nonce_str,omitempty" sign:"true"`
	OPUserId      string   `xml:"op_user_id,omitempty" sign:"true"`
	OutRefundNO   string   `xml:"out_refund_no,omitempty" sign:"true"`
	OutTradeNO    string   `xml:"out_trade_no,omitempty" sign:"true"`
	RefundFee     string   `xml:"refund_fee,omitempty" sign:"true"`
	Sign          string   `xml:"sign,omitempty" sign:"false"`
	TotalFee      string   `xml:"total_fee,omitempty" sign:"true"`
	TransactionId string   `xml:"transaction_id,omitempty" sign:"true"`
	p             *WxPay
}

func (p *WxPay) NewRefundRequest() (r *WXRefundRequest, err error) {
	return &WXRefundRequest{
		p: p,
	}, nil
}

func (r WXRefundRequest) ToXML() string {
	data, err := xml.Marshal(r)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func (r WXRefundRequest) Post() (WXRefundResponse, error) {
	ret := WXRefundResponse{p:r.p}
	r.AppId = r.p.config.APP_ID
	r.MchId = r.p.config.MCH_ID
	r.OPUserId = r.p.config.MCH_ID
	r.NonceStr = RandStr()
	if r.TransactionId == "" && r.OutTradeNO == "" {
		panic(errors.New("TransactionId or OutTradeNO miss"))
	}
	if r.OutRefundNO == "" {
		panic(errors.New("OutRefundNO miss"))
	}
	if r.TotalFee == "" {
		panic(errors.New("TotalFee miss"))
	}
	if r.RefundFee == "" {
		panic(errors.New("RefundFee miss"))
	}
	if r.p.config.TLSConfig == nil {
		panic(errors.New("wx pay key config miss"))
	}
	r.Sign = r.p.WXSign(r)
	http := xweb.NewHTTPClient(WX_PAY_HOST, r.p.config.TLSConfig)
	res, err := http.Post("/secapi/pay/refund", "application/xml", strings.NewReader(r.ToXML()))
	if err != nil {
		return ret, NET_ERROR
	}
	if err := res.ToXml(&ret); err != nil {
		return ret, DATA_UNMARSHAL_ERROR
	}
	if !ret.SignValid() {
		return ret, errors.New("sign error")
	}
	if ret.ReturnCode != SUCCESS {
		return ret, errors.New(ret.ReturnMsg)
	}
	if ret.ResultCode != SUCCESS {
		return ret, errors.New(fmt.Sprintf("code:%d,error:%d", ret.ErrCode, ret.ErrCodeDes))
	}
	return ret, nil
}

/*
<xml>
<return_code><![CDATA[SUCCESS]]></return_code>
<return_msg><![CDATA[OK]]></return_msg>
<appid><![CDATA[wx21b3ee9bd6d16364]]></appid>
<mch_id><![CDATA[1230573602]]></mch_id>
<nonce_str><![CDATA[WDrQTCrHR0KuJVyC]]></nonce_str>
<sign><![CDATA[71B70EE065F17DB4BFDF21D40B4346C9]]></sign>
<result_code><![CDATA[SUCCESS]]></result_code>
<transaction_id><![CDATA[4003952001201608030506918893]]></transaction_id>
<out_trade_no><![CDATA[2016080384874864043]]></out_trade_no>
<out_refund_no><![CDATA[2016080384902087406]]></out_refund_no>
<refund_id><![CDATA[2003952001201608030360901594]]></refund_id>
<refund_channel><![CDATA[]]></refund_channel>
<refund_fee>3900</refund_fee>
<coupon_refund_fee>0</coupon_refund_fee>
<total_fee>3900</total_fee>
<cash_fee>3900</cash_fee>
<coupon_refund_count>0</coupon_refund_count>
<cash_refund_fee>3900</cash_refund_fee>
</xml>
*/

type WXRefundResponse struct {
	XMLName            struct{} `xml:"xml"`
	AppId              string   `xml:"appid,omitempty" sign:"true"`
	CashFee            string   `xml:"cash_fee,omitempty" sign:"true"`
	CashRefundFee      string   `xml:"cash_refund_fee,omitempty" sign:"true"`
	DeviceInfo         string   `xml:"device_info,omitempty" sign:"true"`
	ErrCode            string   `xml:"err_code,omitempty" sign:"true"`
	ErrCodeDes         string   `xml:"err_code_des,omitempty" sign:"true"`
	FeeType            string   `xml:"fee_type,omitempty" sign:"true"`
	MchId              string   `xml:"mch_id,omitempty" sign:"true"`
	NonceStr           string   `xml:"nonce_str,omitempty" sign:"true"`
	OutRefundNO        string   `xml:"out_refund_no,omitempty" sign:"true"`
	OutTradeNO         string   `xml:"out_trade_no,omitempty" sign:"true"`
	RefundChannel      string   `xml:"refund_channel,omitempty" sign:"true"`
	RefundFee          string   `xml:"refund_fee,omitempty" sign:"true"`
	RefundId           string   `xml:"refund_id,omitempty" sign:"true"`
	ResultCode         string   `xml:"result_code,omitempty" sign:"true"`
	ReturnCode         string   `xml:"return_code,omitempty" sign:"true"`
	ReturnMsg          string   `xml:"return_msg,omitempty" sign:"true"`
	SettlementTotalFee string   `xml:"settlement_total_fee,omitempty" sign:"true"`
	Sign               string   `xml:"sign" sign:"false"`
	TotalFee           string   `xml:"total_fee,omitempty" sign:"true"`
	TransactionId      string   `xml:"transaction_id,omitempty" sign:"true"`
	CouponRefundFee    string   `xml:"coupon_refund_fee,omitempty" sign:"true"`
	CouponRefundCount  string   `xml:"coupon_refund_count,omitempty" sign:"true"`
	p                  *WxPay
}

func (r WXRefundResponse) SignValid() bool {
	sign := r.p.WXSign(r)
	return sign == r.Sign
}


//刷新网页授权凭证
//https://api.weixin.qq.com/sns/oauth2/refresh_token
type WXOAuth2RefreshTokenRequest struct {
	AppId        string `json:"appid,omitempty" sign:"true"`
	RefreshToken string `json:"refresh_token,omitempty" sign:"true"`
	GrantType    string `json:"grant_type,omitempty" sign:"true"`
}

type WXOAuth2RefreshTokenResponse struct {
	WXError
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
	Scope        string `json:"scope"`
}

func (p *WxPay) OAuth2RefreshTokenRequest(r WXOAuth2RefreshTokenRequest) (WXOAuth2RefreshTokenResponse, error) {
	ret := WXOAuth2RefreshTokenResponse{}
	r.AppId = p.config.APP_ID
	if r.RefreshToken == "" {
		panic(errors.New("RefreshToken miss"))
	}
	r.GrantType = "refresh_token"
	v := WXParseSignFields(r)
	http := xweb.NewHTTPClient(WX_API_HOST)
	res, err := http.Get("/sns/oauth2/refresh_token", v)
	if err != nil {
		return ret, err
	}
	if err := res.ToJson(&ret); err != nil {
		return ret, err
	}
	if err := ret.Error(); err != nil {
		return ret, err
	}
	return ret, nil
}

//拉取用户信息 AccessToken并非网页授权token
//https://api.weixin.qq.com/cgi-bin/user/info
type WXUserInfoRequest struct {
	AccessToken string `json:"access_token" sign:"true"`
	OpenId      string `json:"openid" sign:"true"`
	Lang        string `json:"lang" sign:"true"`
}

type WXUserInfoResponse struct {
	WXError
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

func (this WXUserInfoRequest) Get() (WXUserInfoResponse, error) {
	ret := WXUserInfoResponse{}
	if this.AccessToken == "" {
		panic(errors.New("AccessToken miss"))
	}
	if this.OpenId == "" {
		panic(errors.New("OpenId miss"))
	}
	if this.Lang == "" {
		this.Lang = "zh_CN"
	}
	v := WXParseSignFields(this)
	http := xweb.NewHTTPClient(WX_API_HOST)
	res, err := http.Get("/cgi-bin/user/info", v)
	if err != nil {
		return ret, err
	}
	if err := res.ToJson(&ret); err != nil {
		return ret, err
	}
	if err := ret.Error(); err != nil {
		return ret, err
	}
	return ret, nil
}

//检验授权凭证（access_token,openid）是否有效
//GET https://api.weixin.qq.com/sns/auth?access_token=ACCESS_TOKEN&openid=OPENID
func AuthGet(token, openid string) WXError {
	ret := WXError{}
	if token == "" {
		panic(errors.New("token error"))
	}
	if openid == "" {
		panic(errors.New("openid error"))
	}
	http := xweb.NewHTTPClient(WX_API_HOST)
	v := xweb.NewHTTPValues()
	v.Set("access_token", token)
	v.Set("openid", openid)
	res, err := http.Get("/sns/auth", v)
	if err != nil {
		ret.ErrCode = 1000000
		ret.ErrMsg = err.Error()
		return ret
	}
	if err := res.ToJson(&ret); err != nil {
		ret.ErrCode = 1000001
		ret.ErrMsg = err.Error()
		return ret
	}
	return ret
}

type WXPayQueryOrderResponse struct {
	XMLName        struct{} `xml:"xml"`
	AppId          string   `xml:"appid" sign:"true"`
	Attach         string   `xml:"attach" sign:"true"`
	BankType       string   `xml:"bank_type" sign:"true"`
	CashFee        string   `xml:"cash_fee" sign:"true"`
	ErrCode        string   `xml:"err_code" sign:"true"`
	ErrCodeDes     string   `xml:"err_code_des" sign:"true"`
	FeeType        string   `xml:"fee_type" sign:"true"`
	IsSubScribe    string   `xml:"is_subscribe" sign:"true"`
	MchId          string   `xml:"mch_id" sign:"true"`
	NonceStr       string   `xml:"nonce_str" sign:"true"`
	OpenId         string   `xml:"openid" sign:"true"`
	OutTradeNo     string   `xml:"out_trade_no" sign:"true"`
	ResultCode     string   `xml:"result_code" sign:"true"`
	ReturnCode     string   `xml:"return_code" sign:"true"`
	ReturnMsg      string   `xml:"return_msg" sign:"true"`
	Sign           string   `xml:"sign" sign:"false"`
	TimeEnd        string   `xml:"time_end" sign:"true"`
	TotalFee       string   `xml:"total_fee" sign:"true"`
	TradeState     string   `xml:"trade_state" sign:"true"`
	TradeStateDesc string   `xml:"trade_state_desc" sign:"true"`
	TradeType      string   `xml:"trade_type" sign:"true"`
	TransactionId  string   `xml:"transaction_id" sign:"true"`
}

func (p *WxPay) SignValidPayQueryOrderResponse(r WXPayQueryOrderResponse) bool {
	sign := p.WXSign(r)
	return sign == r.Sign
}

//正在支付
func (this WXPayQueryOrderResponse) IsPaying() bool {
	if this.ReturnCode != SUCCESS {
		return false
	}
	if this.ResultCode != SUCCESS {
		return false
	}
	return this.TradeState == USERPAYING
}

//支付成功
func (this WXPayQueryOrderResponse) IsPaySuccess() bool {
	if this.ReturnCode != SUCCESS {
		return false
	}
	if this.ResultCode != SUCCESS {
		return false
	}
	return this.TradeState == SUCCESS
}

//2201604122135130001
//https://api.mch.weixin.qq.com/pay/orderquery
type WXPayQueryOrder struct {
	XMLName    struct{} `xml:"xml"`
	AppId      string   `xml:"appid" sign:"true"`
	MchId      string   `xml:"mch_id" sign:"true"`
	OutTradeNo string   `xml:"out_trade_no" sign:"true"`
	NonceStr   string   `xml:"nonce_str" sign:"true"`
	Sign       string   `xml:"sign" sign:"false"` //sign=false表示不参与签名
}

func (this WXPayQueryOrder) ToXML() string {
	data, err := xml.Marshal(this)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func (p *WxPay) PostPayQueryOrder(r WXPayQueryOrder) (WXPayQueryOrderResponse, error) {
	ret := WXPayQueryOrderResponse{}
	r.NonceStr = RandStr()
	r.AppId = p.config.APP_ID
	r.MchId = p.config.MCH_ID
	if r.AppId == "" {
		panic(errors.New("AppId miss"))
	}
	if r.MchId == "" {
		panic(errors.New("MchId miss"))
	}
	r.Sign = p.WXSign(r)
	http := xweb.NewHTTPClient(WX_PAY_HOST)
	res, err := http.Post("/pay/orderquery", "application/xml", strings.NewReader(r.ToXML()))
	if err != nil {
		return ret, NET_ERROR
	}
	if err := res.ToXml(&ret); err != nil {
		return ret, DATA_UNMARSHAL_ERROR
	}
	if !p.SignValidPayQueryOrderResponse(ret) {
		return ret, errors.New("sign error")
	}
	return ret, nil
}

//支付结果通用通知
//微信服务器将会根据统一下单的NotifyURL POST以下数据到商机服务器处理
type PayNotify struct {
	xweb.XMLArgs           `xml:"-"`
	XMLName       struct{} `xml:"xml"` //root node name
	AppId         string   `xml:"appid" sign:"true"`
	Attach        string   `xml:"attach" sign:"true"`
	BankType      string   `xml:"bank_type" sign:"true"`
	CashFee       string   `xml:"cash_fee" sign:"true"`
	CashFeeType   string   `xml:"cash_fee_type" sign:"true"`
	CouponCount   string   `xml:"coupon_count" sign:"true"`
	CouponFee     string   `xml:"coupon_fee" sign:"true"`
	DeviceInfo    string   `xml:"device_info" sign:"true"`
	ErrCode       string   `xml:"err_code" sign:"true"`
	ErrCodeDes    string   `xml:"err_code_des" sign:"true"`
	FeeType       string   `xml:"fee_type" sign:"true"`
	IsSubScribe   string   `xml:"is_subscribe" sign:"true"` //Y or N
	MchId         string   `xml:"mch_id" sign:"true"`
	NonceStr      string   `xml:"nonce_str" sign:"true"`
	OpenId        string   `xml:"openid" sign:"true"`
	OutTradeNo    string   `xml:"out_trade_no" sign:"true"`
	ResultCode    string   `xml:"result_code" sign:"true"` //SUCCESS or FAIL
	ReturnCode    string   `xml:"return_code" sign:"true"` //SUCCESS or FAIL
	ReturnMsg     string   `xml:"return_msg" sign:"true"`  //返回信息，如非空，为错误原因
	Sign          string   `xml:"sign" sign:"false"`       //sign=false表示不参与签名
	TimeEnd       string   `xml:"time_end" sign:"true"`
	TotalFee      string   `xml:"total_fee" sign:"true"`
	TradeType     string   `xml:"trade_type" sign:"true"` //JSAPI、NATIVE、APP
	TransactionId string   `xml:"transaction_id" sign:"true"`
	p             *WxPay
}

func (p *WxPay) NewWXPayNotify(body []byte) (n *PayNotify, err error) {
	n = &PayNotify{
		p: p,
	}
	err = xml.Unmarshal(body, n)
	return
}

//签名校验
func (r PayNotify) SignValid() bool {
	sign := r.p.WXSign(r)
	return sign == r.Sign
}

//nil表示没有错误
func (r PayNotify) IsError() error {
	if r.ReturnCode != SUCCESS {
		return errors.New(r.ReturnMsg)
	}
	if r.ResultCode != SUCCESS {
		return errors.New(fmt.Sprintf("ERROR:%d,%s", r.ErrCode, r.ErrCodeDes))
	}
	if !r.SignValid() {
		return errors.New("sign valid error")
	}
	return nil
}

//商户处理后返回格式
type WXPayResultResponse struct {
	xweb.XMLModel       `xml:"-"`
	XMLName    struct{} `xml:"xml"`                   //root node name
	ReturnCode string   `xml:"return_code,omitempty"` //SUCCESS or FAIL
	ReturnMsg  string   `xml:"return_msg,omitempty"`  //OK
}

func (r WXPayResultResponse) ToXML() string {
	data, err := xml.Marshal(r)
	if err != nil {
		panic(err)
	}
	return string(data)
}

//恶心的微信签名用noncestr,脚本里用nonceStr
type WXConfigForJS struct {
	Debug     bool     `json:"debug" sign:"false"`
	AppId     string   `json:"appId" sign:"false"`
	Timestamp string   `json:"timestamp" sign:"true"`
	NonceStr  string   `json:"nonceStr" sign:"true"`
	Signature string   `json:"signature" sign:"false"`
	JSApiList []string `json:"jsApiList" sign:"false"`
	p         *WxPay
}

func (p *WxPay) NewConfigForJS() (r *WXConfigForJS, err error) {
	return &WXConfigForJS{
		p: p,
	}, nil
}

func (r WXConfigForJS) ToScript(jsticket string, url string) (template.JS, error) {
	r.AppId = r.p.config.APP_ID
	r.Timestamp = TimeNowString()
	r.NonceStr = RandStr()
	if r.JSApiList == nil {
		r.JSApiList = []string{}
	}
	v := xweb.NewHTTPValues()
	v.Set("timestamp", r.Timestamp)
	v.Set("noncestr", r.NonceStr)
	v.Set("jsapi_ticket", jsticket)
	v.Set("url", url)
	r.Signature = xweb.SHA1String(v.RawEncode())
	data, err := json.Marshal(r)
	if err != nil {
		return template.JS(""), err
	}
	return template.JS(data), nil
}

//为jsapi支付返回给客户端用于客户端发起支付
type WXPayReqForJS struct {
	AppId     string `json:"appId,omitempty" sign:"true"`
	Timestamp int64  `json:"timeStamp,omitempty" sign:"true"`
	Package   string `json:"package,omitempty" sign:"true"`
	NonceStr  string `json:"nonceStr,omitempty" sign:"true"`
	SignType  string `json:"signType,omitempty" sign:"true"`
	PaySign   string `json:"paySign,omitempty" sign:"false"`
}

type WXPayReqScript struct {
	Timestamp int64  `json:"timestamp,omitempty"`
	Package   string `json:"package,omitempty"`
	NonceStr  string `json:"nonceStr,omitempty"`
	SignType  string `json:"signType,omitempty"`
	PaySign   string `json:"paySign,omitempty"`
}

func (this WXPayReqForJS) ToScript() (template.JS, error) {
	s := WXPayReqScript{}
	s.NonceStr = this.NonceStr
	s.Package = this.Package
	s.PaySign = this.PaySign
	s.SignType = this.SignType
	s.Timestamp = this.Timestamp
	data, err := json.Marshal(s)
	if err != nil {
		return template.JS(""), err
	}
	return template.JS(data), nil
}

func (p *WxPay) NewWXPayReqScript(prepayid string) WXPayReqScript {
	d := WXPayReqForJS{}
	d.AppId = p.config.APP_ID
	d.Package = "prepay_id=" + prepayid
	d.NonceStr = RandStr()
	d.Timestamp = TimeNow()
	d.SignType = "MD5"
	d.PaySign = p.WXSign(d)
	s := WXPayReqScript{}
	s.NonceStr = d.NonceStr
	s.Package = d.Package
	s.PaySign = d.PaySign
	s.SignType = d.SignType
	s.Timestamp = d.Timestamp
	return s
}

//新建jsapi支付返回
func (p *WxPay) NewWXPayReqForJS(prepayid string) WXPayReqForJS {
	d := WXPayReqForJS{}
	d.AppId = p.config.APP_ID
	d.Package = "prepay_id=" + prepayid
	d.NonceStr = RandStr()
	d.Timestamp = TimeNow()
	d.SignType = "MD5"
	d.PaySign = p.WXSign(d)
	return d
}

//为app支付返回给客户端用于客户端发起支付
type WXPayReqForApp struct {
	AppId     string `json:"appid,omitempty" sign:"true"`
	NonceStr  string `json:"noncestr,omitempty" sign:"true"`
	Package   string `json:"package,omitempty" sign:"true"` //APP支付固定(Sign=WXPay)
	PartnerId string `json:"partnerid,omitempty" sign:"true"`
	PrepayId  string `json:"prepayid,omitempty" sign:"true"` //统一下单返回
	Sign      string `json:"sign,omitempty" sign:"false"`
	Timestamp int64  `json:"timestamp,omitempty" sign:"true"`
}

func (this WXPayReqForApp) String() string {
	data, err := json.Marshal(this)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

//新建APP支付返回
func (p *WxPay) NewWXPayReqForApp(prepayid string) WXPayReqForApp {
	d := WXPayReqForApp{}
	d.AppId = p.config.APP_ID
	d.PrepayId = prepayid
	d.PartnerId = p.config.MCH_ID
	d.Package = "Sign=WXPay"
	d.NonceStr = RandStr()
	d.Timestamp = TimeNow()
	d.Sign = p.WXSign(d)
	return d
}

//微信支付:统一下单
//https://api.mch.weixin.qq.com/pay/unifiedorder
type WXUnifiedOrderRequest struct {
	XMLName        struct{} `xml:"xml"` //root node name
	AppId          string   `xml:"appid,omitempty" sign:"true"`
	Attach         string   `xml:"attach,omitempty" sign:"true"`
	Body           string   `xml:"body,omitempty" sign:"true"`
	Detail         string   `xml:"detail,omitempty" sign:"true"`
	DeviceInfo     string   `xml:"device_info,omitempty" sign:"true"`
	FeeType        string   `xml:"fee_type,omitempty" sign:"true"`
	GoodsTag       string   `xml:"goods_tag,omitempty" sign:"true"`
	LimitPay       string   `xml:"limit_pay,omitempty" sign:"true"`
	MchId          string   `xml:"mch_id,omitempty" sign:"true"`
	NonceStr       string   `xml:"nonce_str,omitempty" sign:"true"`
	NotifyURL      string   `xml:"notify_url,omitempty" sign:"true"`
	OpenId         string   `xml:"openid,omitempty" sign:"true"` //TradeType=TRADE_TYPE_JSAPI 必须
	OutTradeNo     string   `xml:"out_trade_no,omitempty" sign:"true"`
	ProductId      string   `xml:"product_id,omitempty" sign:"true"` //TradeType=TRADE_TYPE_NATIVE 必须
	Sign           string   `xml:"sign,omitempty"  sign:"false"`     //sign=false表示不参与签名
	SpBillCreateIp string   `xml:"spbill_create_ip,omitempty" sign:"true"`
	TimeExpire     string   `xml:"time_expire,omitempty" sign:"true"`
	TimeStart      string   `xml:"time_start,omitempty" sign:"true"`
	TotalFee       string   `xml:"total_fee,omitempty" sign:"true"`
	TradeType      string   `xml:"trade_type,omitempty" sign:"true"`
	p              *WxPay
}

func (p *WxPay) NewUnifiedOrderRequest() (r *WXUnifiedOrderRequest) {
	return &WXUnifiedOrderRequest{
		p: p,
	}
}

//微信支付:统一下单返回数据
type WXUnifiedOrderResponse struct {
	XMLName    struct{} `xml:"xml"` //root node name
	AppId      string   `xml:"appid,omitempty" sign:"true"`
	CodeURL    string   `xml:"code_url,omitempty" sign:"true"` //trade_type=NATIVE返回code url
	DeviceInfo string   `xml:"device_info,omitempty" sign:"true"`
	ErrCode    string   `xml:"err_code,omitempty" sign:"true"`
	ErrCodeDes string   `xml:"err_code_des,omitempty" sign:"true"`
	MchId      string   `xml:"mch_id,omitempty" sign:"true"`
	NonceStr   string   `xml:"nonce_str,omitempty" sign:"true"`
	PrePayId   string   `xml:"prepay_id,omitempty" sign:"true"`
	ResultCode string   `xml:"result_code,omitempty" sign:"true"` //SUCCESS or FAIL
	ReturnCode string   `xml:"return_code,omitempty" sign:"true"` //SUCCESS or FAIL
	ReturnMsg  string   `xml:"return_msg,omitempty" sign:"true"`  //返回信息，如非空，为错误原因
	Sign       string   `xml:"sign,omitempty"  sign:"false"`      //sign=false表示不参与签名
	TradeType  string   `xml:"trade_type,omitempty" sign:"true"`
}

func (r WXUnifiedOrderResponse) Error() error {
	if r.ReturnCode != SUCCESS {
		return errors.New("ERROR:" + r.ReturnMsg)
	}
	if r.ResultCode != SUCCESS {
		return errors.New("ERROR:" + r.ErrCode + "," + r.ErrCodeDes)
	}
	return nil
}

func (r WXUnifiedOrderRequest) Post() (WXUnifiedOrderResponse, error) {
	ret := WXUnifiedOrderResponse{}
	if r.TotalFee == "" {
		panic(errors.New("TotalFee must > 0 "))
	}
	r.NonceStr = RandStr()
	if r.NotifyURL == "" {
		panic(errors.New("NotifyURL miss"))
	}
	r.AppId = r.p.config.APP_ID
	r.MchId = r.p.config.MCH_ID
	if r.AppId == "" {
		panic(errors.New("AppId miss"))
	}
	if r.MchId == "" {
		panic(errors.New("MchId miss"))
	}
	if r.NotifyURL == "" {
		panic(errors.New("NotifyURL miss"))
	}
	if r.TradeType == "" {
		panic(errors.New("TradeType must set"))
	}
	if r.TradeType == TRADE_TYPE_JSAPI && r.OpenId == "" {
		panic(errors.New(TRADE_TYPE_JSAPI + " openid empty"))
	}
	if r.TradeType == TRADE_TYPE_NATIVE && r.ProductId == "" {
		panic(errors.New(TRADE_TYPE_NATIVE + " product_id empty"))
	}
	r.Sign = r.p.WXSign(r)
	http := xweb.NewHTTPClient(WX_PAY_HOST)
	res, err := http.Post("/pay/unifiedorder", "application/xml", strings.NewReader(r.ToXML()))
	if err != nil {
		return ret, err
	}
	if err := res.ToXml(&ret); err != nil {
		return ret, err
	}
	if err := ret.Error(); err != nil {
		return ret, err
	}
	return ret, nil
}

func (r WXUnifiedOrderRequest) ToXML() string {
	data, err := xml.Marshal(r)
	if err != nil {
		panic(err)
	}
	return string(data)
}
