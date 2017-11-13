# payment
微信支付，支付宝支付

## 微信支付
下单
```go
wxPay := payment.NewWxPayClient(WxAppId, MchId, MchKey)
// 下单并返回小程序调起支付需要的字段
info,err := wxPay.CreateWxPayPayInfo(tranId, order.Remark, strconv.Itoa(order.PricePay), clientIp, userOpenId, wxNotifyUrl)
```
处理回调
```go
wxPay := payment.NewWxPayClient(paySetting.WxAppId, paySetting.MchId, paySetting.MchKey)
// 检查回调是否有效, 并返回必要的交易信息
trade, err := wxPay.CheckWxPayNotify(request)
```

## 支付宝支付
todo