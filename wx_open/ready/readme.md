## 准备
- 全网发布

### usage
```
msg, err := common_party.DecodeMessage(msgSignature, timestamp, nonce, body)
if err != nil {
    log.Errorf("Decrypt err:%v", err)
    ctx.String(500, "Decrypt err:%v", err)
    return
}

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