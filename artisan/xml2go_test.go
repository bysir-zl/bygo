package artisan

import (
	"testing"
	"regexp"
)

func TestXml(t *testing.T) {
	res := `
<xml>

<return_code><![CDATA[SUCCESS]]></return_code>

<return_msg><![CDATA[]]></return_msg>

<mch_appid><![CDATA[wxec38b8ff840bd989]]></mch_appid>

<mchid><![CDATA[10013274]]></mchid>

<device_info><![CDATA[]]></device_info>

<nonce_str><![CDATA[lxuDzMnRjpcXzxLx0q]]></nonce_str>

<result_code><![CDATA[SUCCESS]]></result_code>

<partner_trade_no><![CDATA[10013574201505191526582441]]></partner_trade_no>

<payment_no><![CDATA[1000018301201505190181489473]]></payment_no>

<payment_time><![CDATA[2015-05-19 15：26：59]]></payment_time>

</xml>
`

	{
		r, _ := regexp.Compile(`<(.*?)>`)
		rsp := r.FindStringSubmatch(res)
		t.Log(rsp)
	}
	{
		r, _ := regexp.Compile(`<(.*?)><!\[CDATA\[(.*?)\]\]></.*>`)
		rsp := r.FindAllStringSubmatch(res, -1)
		t.Log(rsp)
	}
}
