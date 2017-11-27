## 第三方平台
- 全网发布
- 小程序管理
- 认证

### 全网发布
在接收到微信回调以后:
```go
// body 就是微信回调的请求body
msg, err := common_party.DecodeMessage(msgSignature, timestamp, nonce, body)
if err != nil {
    log.Errorf("Decrypt err:%v", err)
    ctx.String(500, "Decrypt err:%v", err)
    return
}
// 如果这是测试请求就会被拦截
stop, rsp, err := ready.FilterReady(appId, msg)
if err != nil {
    log.Errorf("FilterReady err:%v", err)
    ctx.String(500, "FilterReady err:%v", err)
    return
}
if stop {
    ctx.Writer.Write(rsp)
    return
}
```
