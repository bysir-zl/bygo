## 微信开放平台


- 加解密
- 公众平台发送消息
- 第三方平台授权
- 第三方平台代小程序实现业务


### usage 
先配置app信息, 毕竟微信的接口是加密的
```go
wx_open.Init(token, aeskey, appid, appSecret)
```
