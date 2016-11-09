package swagger_test


// @INFO title
// @desc:这是一个
// @version:v0.0.1
// @email:zhangliang@kuaifazs.com
// @license:-
// @host:-
// @basePath:/v1

// @BASE BaseInfo
// @debug: 调试模式,boolean , formData,true ,false;
// @debug2: 调试模式,boolean , formData,true ,false

// @API 首页信息
// @router : /index/ , post , tags , operation
//
// @parameters :
// @name: 姓名,string , formData,123123,true ;
// @sex: 姓名,string , formData,0 ,true ;
// @BASE:BaseInfo
//
// @responses :
// @200: 成功,error
// @400: 失败
func index() {
	var name string
	var sex string
	name = name
	sex = sex
}

// @API 生成快发的TOKEN,OPENID
// @router : oauth/request, post, tags, operationId
// @params :
// @original : base64加密后的数据,string , formData, default;
// @guid : guid,string, formData

func ind() {

}

// @DEF error,object
// @code:返回马,string,0为成功|其他code见errorinfo;
// @errorinfo:错误详情



// @DEF e,object
// @x:123

// @DEF gameinfo,object
// @logLevel:名称,类型,默认值;
// @uploadLog:;
// @versionControl:;
// @version:;
// @downloadUrl:;
// @release_note;;
// @ad_status:;
// @reserveNotify:;
// @reserveDialog:,ref,e
