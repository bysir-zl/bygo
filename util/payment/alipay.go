package payment

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/cxuhua/xweb"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

//alipay
type APKeyConfig struct {
	PARTNER_ID          string //商户id
	SELLER_EMAIL        string //商户支付email
	SIGN_TYPE           string //签名类型 RSA
	ALIPAY_KEY          string //阿里支付密钥
	PARTNET_PRIVATE_KEY string //商户私钥
	ALIPAY_PUBLIC_KEY   string //阿里支付公钥
}

var (
	AP_PAY_CONFIG       APKeyConfig     = APKeyConfig{}
	PARTNET_PRIVATE_KEY *rsa.PrivateKey = nil
	ALIPAY_PUBLIC_KEY   *rsa.PublicKey  = nil
)

func InitAPKey(conf APKeyConfig) {
	AP_PAY_CONFIG = conf
	//加载商户私钥
	if block, _ := pem.Decode([]byte(AP_PAY_CONFIG.PARTNET_PRIVATE_KEY)); block != nil {
		if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
			panic(err)
		} else {
			PARTNET_PRIVATE_KEY = key
		}
	} else {
		panic("load PARTNET_PRIVATE_KEY failed")
	}
	//加载支付宝公钥
	if block, _ := pem.Decode([]byte(AP_PAY_CONFIG.ALIPAY_PUBLIC_KEY)); block != nil {
		if pub, err := x509.ParsePKIXPublicKey(block.Bytes); err != nil {
			panic(err)
		} else {
			ALIPAY_PUBLIC_KEY = pub.(*rsa.PublicKey)
		}
	} else {
		panic("load ALIPAY_PUBLIC_KEY failed")
	}
}

func APMD5Sign(v interface{}) string {
	http := APParseSignFields(v)
	str := http.RawEncode() + AP_PAY_CONFIG.ALIPAY_KEY
	return xweb.MD5String(str)
}

func APSHA1Sign(v interface{}) string {
	http := APParseSignFields(v)
	str := http.RawEncode()
	h := crypto.SHA1.New()
	h.Write([]byte(str))
	hashed := h.Sum(nil)
	if s, err := rsa.SignPKCS1v15(rand.Reader, PARTNET_PRIVATE_KEY, crypto.SHA1, hashed); err != nil {
		panic(err)
	} else {
		return base64.StdEncoding.EncodeToString(s)
	}
}

func APParseSignFields(src interface{}) xweb.HTTPValues {
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
		} else if xn = tf.Tag.Get("form"); xn != "" {
			name = strings.Split(xn, ",")[0]
		} else {
			continue
		}
		//for xml
		ns := strings.Split(name, ">")
		if len(ns) > 0 {
			name = ns[len(ns)-1]
		} else {
			name = ns[0]
		}
		values.Add(name, sv)
	}
	return values
}

//校验来自阿里的数据
func APRSAVerify(src string, sign string) error {
	h := crypto.SHA1.New()
	h.Write([]byte(src))
	hashed := h.Sum(nil)
	data, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(ALIPAY_PUBLIC_KEY, crypto.SHA1, hashed, data)
}

const (
	//如果异步订单处理成功返回
	DO_SUCCESS = "success"
)

//TradeStatus
const (
	WAIT_BUYER_PAY  = "WAIT_BUYER_PAY"
	TRADE_CLOSED    = "TRADE_CLOSED"
	TRADE_SUCCESS   = "TRADE_SUCCESS"
	TRADE_FINISHED  = "TRADE_FINISHED"
	REFUND_SUCCESS  = "REFUND_SUCCESS"
	REFUND_CLOSED   = "REFUND_CLOSED"
	TRADE_NOT_EXIST = "TRADE_NOT_EXIST"
)

//https://mapi.alipay.com/gateway.do?
//_input_charset=utf-8&out_trade_no=1201604242253410295&partner=2088121797205248&
//service=single_trade_query&sign=283e436499a7fbfe5e6bccdf72946ec5&sign_type=md5
type APPayQueryOrder struct {
	Service      string `json:"service" sign:"true"`
	Partner      string `json:"partner" sign:"true"`
	InputCharset string `json:"_input_charset" sign:"true"`
	OutTradeNo   string `json:"out_trade_no" sign:"true"`
}

type APPayQueryOrderResponse struct {
	XMLName      struct{} `xml:"alipay" sign:"false"`
	IsSuccess    string   `xml:"is_success" sign:"false"`
	Body         string   `xml:"response>trade>body" sign:"true"`
	BuyerEmail   string   `xml:"response>trade>buyer_email" sign:"true"`
	BuyerId      string   `xml:"response>trade>buyer_id" sign:"true"`
	Discount     string   `xml:"response>trade>discount" sign:"true"`
	Locked       string   `xml:"response>trade>flag_trade_locked" sign:"true"`
	GMTCreate    string   `xml:"response>trade>gmt_create" sign:"true"`
	GMTLast      string   `xml:"response>trade>gmt_last_modified_time" sign:"true"`
	GMTPayment   string   `xml:"response>trade>gmt_payment" sign:"true"`
	IsTotalFee   string   `xml:"response>trade>is_total_fee_adjust" sign:"true"`
	OperatorRole string   `xml:"response>trade>operator_role" sign:"true"`
	OutTradeNo   string   `xml:"response>trade>out_trade_no" sign:"true"`
	PaymentType  string   `xml:"response>trade>payment_type" sign:"true"`
	Price        string   `xml:"response>trade>price" sign:"true"`
	Quantity     string   `xml:"response>trade>quantity" sign:"true"`
	SellerEmail  string   `xml:"response>trade>seller_email" sign:"true"`
	SellerId     string   `xml:"response>trade>seller_id" sign:"true"`
	Subject      string   `xml:"response>trade>subject" sign:"true"`
	Timeout      string   `xml:"response>trade>time_out" sign:"true"`
	TimeoutType  string   `xml:"response>trade>time_out_type" sign:"true"`
	ToBuyerFee   string   `xml:"response>trade>to_buyer_fee" sign:"true"`
	ToSellerFee  string   `xml:"response>trade>to_seller_fee" sign:"true"`
	TotalFee     string   `xml:"response>trade>total_fee" sign:"true"`
	TradeNo      string   `xml:"response>trade>trade_no" sign:"true"`
	TradeStatus  string   `xml:"response>trade>trade_status" sign:"true"`
	UseCoupon    string   `xml:"response>trade>use_coupon" sign:"true"`
	Sign         string   `xml:"sign" sign:"false"`
	SignType     string   `xml:"sign_type" sign:"false"`
	Error        string   `xml:"error" sign:"false"`
}

func (this APPayQueryOrderResponse) IsPaySuccess() bool {
	if this.IsSuccess != "T" {
		return false
	}
	return this.TradeStatus == TRADE_SUCCESS
}

/*
<alipay>
	<is_success>T</is_success>
	<request>
		<param name="_input_charset">utf-8</param>
		<param name="service">single_trade_query</param>
		<param name="partner">2088121797205248</param>
		<param name="out_trade_no">1201604242253410295</param>
	</request>
	<response>
		<trade>
			<body>订单支付</body>
			<buyer_email>cxuhua@gmail.com</buyer_email>
			<buyer_id>2088002003565555</buyer_id>
			<discount>0.00</discount>
			<flag_trade_locked>0</flag_trade_locked>
			<gmt_create>2016-04-24 22:54:05</gmt_create>
			<gmt_last_modified_time>2016-04-24 22:54:06</gmt_last_modified_time>
			<gmt_payment>2016-04-24 22:54:06</gmt_payment>
			<is_total_fee_adjust>F</is_total_fee_adjust>
			<operator_role>B</operator_role>
			<out_trade_no>1201604242253410295</out_trade_no>
			<payment_type>1</payment_type>
			<price>2.12</price>
			<quantity>1</quantity>
			<seller_email>57730141@qq.com</seller_email>
			<seller_id>2088121797205248</seller_id>
			<subject>订单支付</subject>
			<time_out>2016-07-24 22:54:06</time_out>
			<time_out_type>finishFPAction</time_out_type>
			<to_buyer_fee>0.00</to_buyer_fee>
			<to_seller_fee>2.12</to_seller_fee>
			<total_fee>2.12</total_fee>
			<trade_no>2016042421001004550217245009</trade_no>
			<trade_status>TRADE_SUCCESS</trade_status>
			<use_coupon>F</use_coupon>
		</trade>
	</response>
	<sign>ecc83e462668b8a7bc695e24249e2db6</sign>
	<sign_type>MD5</sign_type>
</alipay>
*/
func (this APPayQueryOrder) Get() (APPayQueryOrderResponse, error) {
	ret := APPayQueryOrderResponse{}
	if this.OutTradeNo == "" {
		panic(errors.New("OutTradeNo miss"))
	}
	this.Service = "single_trade_query"
	this.Partner = AP_PAY_CONFIG.PARTNER_ID
	this.InputCharset = "utf-8"
	c := xweb.NewHTTPClient("https://mapi.alipay.com")
	v := xweb.NewHTTPValues()
	v.Set("service", this.Service)
	v.Set("partner", this.Partner)
	v.Set("_input_charset", this.InputCharset)
	v.Set("out_trade_no", this.OutTradeNo)
	v.Set("sign_type", "MD5")
	v.Set("sign", APMD5Sign(this))
	res, err := c.Get("/gateway.do", v)
	if err != nil {
		return ret, err
	}
	if err := res.ToXml(&ret); err != nil {
		return ret, err
	}
	if ret.IsSuccess != "T" {
		return ret, nil
	}
	if ret.Sign != APMD5Sign(ret) {
		return ret, errors.New("sign data error")
	}
	return ret, nil
}

//支付宝服务器异步通知参数说明
type APPayResultNotifyArgs struct {
	NotifyTime   string `form:"notify_time" sign:"true"`
	NotifyType   string `form:"notify_type" sign:"true"`
	NotifyId     string `form:"notify_id" sign:"true"`
	SignType     string `form:"sign_type" sign:"false"` //RSA
	Sign         string `form:"sign" sign:"false"`
	OutTradeNO   string `form:"out_trade_no" sign:"true"`
	Subject      string `form:"subject" sign:"true"`
	PaymentType  string `form:"payment_type" sign:"true"` //1
	TradeNO      string `form:"trade_no" sign:"true"`
	TradeStatus  string `form:"trade_status" sign:"true"`
	SellerId     string `form:"seller_id" sign:"true"`
	SellerEmail  string `form:"seller_email" sign:"true"`
	BuyerId      string `form:"buyer_id" sign:"true"`
	BuyerEmail   string `form:"buyer_email" sign:"true"`
	TotalFee     string `form:"total_fee" sign:"true"`
	Quantity     string `form:"quantity" sign:"true"`
	Price        string `form:"price" sign:"true"`
	Body         string `form:"body" sign:"true"`
	GMTCreate    string `form:"gmt_create" sign:"true"`
	GMTPayment   string `form:"gmt_payment" sign:"true"`
	FeeAdjust    string `form:"is_total_fee_adjust" sign:"true"`
	UseCoupon    string `form:"use_coupon" sign:"true"`
	Discount     string `form:"discount" sign:"true"`
	RefundStatus string `form:"refund_status" sign:"true"`
	GMTRefund    string `form:"gmt_refund" sign:"true"`
	GMTClose     string `form:"gmt_close" sign:"true"`
}

func (this APPayResultNotifyArgs) GetTotalFee() float32 {
	v, err := strconv.ParseFloat(this.TotalFee, 64)
	if err != nil {
		panic(err)
	}
	return float32(v)
}

func (this APPayResultNotifyArgs) IsError() error {
	if !this.IsValid() {
		return errors.New("data sign error")
	}
	if !this.IsFromAlipay() {
		return errors.New("data not form alipay")
	}
	if !this.IsSuccess() {
		return errors.New("pay status not success")
	}
	return nil
}

func (this APPayResultNotifyArgs) IsSuccess() bool {
	return this.TradeStatus == TRADE_SUCCESS || this.TradeStatus == TRADE_FINISHED
}

//是否来自支付宝
func (this APPayResultNotifyArgs) IsFromAlipay() bool {
	q := xweb.NewHTTPValues()
	q.Set("service", "notify_verify")
	q.Set("partner", AP_PAY_CONFIG.PARTNER_ID)
	q.Set("notify_id", this.NotifyId)
	http := xweb.NewHTTPClient("https://mapi.alipay.com")
	res, err := http.Get("/gateway.do", q)
	if err != nil {
		return false
	}
	d, err := res.ToBytes()
	if err != nil {
		return false
	}
	return string(d) == "true"
}

//签名校验
func (this APPayResultNotifyArgs) IsValid() bool {
	v := APParseSignFields(this)
	s := v.RawEncode()
	if err := APRSAVerify(s, this.Sign); err != nil {
		return false
	}
	return true
}

func (this APPayResultNotifyArgs) String() string {
	d, err := json.Marshal(this)
	if err != nil {
		return err.Error()
	}
	return string(d)
}

func NewAPPayResultNotifyArgs(request []byte) (APPayResultNotifyArgs, error) {
	args := APPayResultNotifyArgs{}
	form, err := url.ParseQuery(string(request))
	if err != nil {
		return args, err
	}
	xweb.MapFormBindType(&args, form, nil, nil, nil)
	if err := args.IsError(); err != nil {
		return args, err
	}
	return args, nil
}

//阿里支付请求参数
type APPayReqForApp struct {
	Service      string `json:"service,omitempty" sign:"true"`
	Partner      string `json:"partner,omitempty" sign:"true"`
	InputCharset string `json:"_input_charset,omitempty" sign:"true"`
	SignType     string `json:"sign_type,omitempty" sign:"false"`
	Sign         string `json:"sign,omitempty" sign:"false"`
	NotifyURL    string `json:"notify_url,omitempty" sign:"true"`
	OutTradeNO   string `json:"out_trade_no,omitempty" sign:"true"`
	Subject      string `json:"subject,omitempty" sign:"true"`
	PaymentType  string `json:"payment_type,omitempty" sign:"true"`
	SellerId     string `json:"seller_id,omitempty" sign:"true"`
	TotalFee     string `json:"total_fee,omitempty" sign:"true"`
	Body         string `json:"body,omitempty" sign:"true"`
}

func (this APPayReqForApp) String() string {
	if this.NotifyURL == "" {
		panic(errors.New("NotifyURL miss"))
	}
	if this.OutTradeNO == "" {
		panic(errors.New("OutTradeNO miss"))
	}
	if this.Subject == "" {
		panic(errors.New("Subject miss"))
	}
	if this.TotalFee == "" {
		panic(errors.New("TotalFee miss"))
	}
	if this.Body == "" {
		panic(errors.New("Body miss"))
	}
	this.Sign = url.QueryEscape(APSHA1Sign(this))
	values := xweb.NewHTTPValues()
	t := reflect.TypeOf(this)
	v := reflect.ValueOf(this)
	for i := 0; i < t.NumField(); i++ {
		tf := t.Field(i)
		tv := v.Field(i)
		if !tv.IsValid() {
			continue
		}
		jt := strings.Split(tf.Tag.Get("json"), ",")
		if len(jt) == 0 || jt[0] == "" {
			continue
		}
		sv := fmt.Sprintf(`%v`, tv.Interface())
		if sv == "" {
			continue
		}
		values.Add(jt[0], sv)
	}
	return values.RawEncode()
}

func NewAPPayReqForApp() APPayReqForApp {
	d := APPayReqForApp{}
	d.Service = "mobile.securitypay.pay"
	d.Partner = AP_PAY_CONFIG.PARTNER_ID
	d.InputCharset = "utf-8"
	d.SignType = "RSA"
	d.PaymentType = "1"
	d.SellerId = AP_PAY_CONFIG.SELLER_EMAIL
	return d
}
