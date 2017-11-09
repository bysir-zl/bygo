package wx_third_party

const (
    URL_ComponentToken = "https://api.weixin.qq.com/cgi-bin/component/api_component_token"

    URL_PreAuthcode = "https://api.weixin.qq.com/cgi-bin/component/api_create_preauthcode?component_access_token="
    // 公众号或小程序的接口调用凭据
    URL_OtherAuthToken = "https://api.weixin.qq.com/cgi-bin/component/api_query_auth?component_access_token="
    // 获取（刷新）授权公众号或小程序的接口调用凭据（令牌）
    URL_RefreshOtherAuthToken = "https://api.weixin.qq.com/cgi-bin/component/api_authorizer_token?component_access_token="
)
