package util

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"net/url"
	"strings"
)

// 封装一层 getAll 获取列表 ,方便使用
func GetAll(c beego.Controller, model interface{}, contains interface{}, defaultLimit int64) (err error) {
	f, err := BuildFilter(c, defaultLimit)
	if err != nil {
		return
	}
	err = GetAllByFilter(model, contains, f)
	return
}
func GetAllWithTotal(c beego.Controller, model interface{}, contains interface{}, defaultLimit int64) (total int64, err error) {
	f, err := BuildFilter(c, defaultLimit)
	if err != nil {
		return
	}
	total, err = GetAllByFilterWithTotal(model, contains, f)
	return
}

type Filter struct {
	Fields []string                 // 字段
	Order  []string                 // 排序 : -id  id
	Where  map[string][]interface{} // id:1;user_id:in1,2,3|<=2
	Limit  int64                    //
	Offset int64                    //
}

// 根据输入生成过滤条件
// query=id:1;type:in1,3,5 & order=-ascId,descType & limit=0 & offset=0
func BuildFilter(c beego.Controller, defaultLimit int64) (f *Filter, err error) {
	var fields []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = defaultLimit
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
		if limit > 50 {
			limit = 50
		}
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}

	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v;k:v
	if v := c.GetString("query"); v != "" {
		v, _ = url.QueryUnescape(v)
		for _, cond := range strings.Split(v, ";") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				err = errors.New("invalid query key/value pair")
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	where := map[string][]interface{}{} // condition => args
	// condition : inx,x|<x|>x|>=x|<=x|=x|x|!=x|
	for f, cond := range query {
		// 支持一个字段多个条件
		for _, x := range strings.Split(cond, "|") {
			isNot := ""
			if len(x) > 1 && x[0] == '!' {
				isNot = "!"
				x = x[1:]
			}

			if strings.Index(x, "in") == 0 {
				args := []interface{}{}
				for _, v := range strings.Split(x[2:], ",") {
					args = append(args, v)
				}
				where[isNot+f+"__in"] = args
			} else if strings.Index(x, ">=") == 0 {
				where[isNot+f+"__gte"] = []interface{}{x[2:]}
			} else if strings.Index(x, ">") == 0 {
				where[isNot+f+"__gt"] = []interface{}{x[1:]}
			} else if strings.Index(x, "<=") == 0 {
				where[isNot+f+"__lte"] = []interface{}{x[2:]}
			} else if strings.Index(x, "<") == 0 {
				where[isNot+f+"__lt"] = []interface{}{x[1:]}
			} else if strings.Index(x, "=") == 0 {
				where[isNot+f ] = []interface{}{x[1:]}
			} else if strings.Index(x, "like") == 0 {
				where[isNot+f+"__icontains"] = []interface{}{x[4:]}
			} else {
				where[isNot+f ] = []interface{}{x}
			}
		}
	}

	f = &Filter{
		Fields: fields,
		Limit:  limit,
		Offset: offset,
		Order:  order,
		Where:  where,
	}

	return
}

func GetAllByFilter(model interface{}, contains interface{}, f *Filter) (err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(model)
	// query k=v

	for k, v := range f.Where {
		if len(k) > 1 && k[0] == '!' {
			k = k[1:]
			qs = qs.Exclude(k, v...)
		} else {
			qs = qs.Filter(k, v...)
		}
	}

	// order by:
	qs = qs.OrderBy(f.Order...)

	_, err = qs.Limit(f.Limit, f.Offset).All(contains, f.Fields...)
	return
}

func GetAllByFilterWithTotal(model interface{}, contains interface{}, f *Filter) (total int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(model)
	// query k=v

	for k, v := range f.Where {
		if len(k) > 1 && k[0] == '!' {
			k = k[1:]
			qs = qs.Exclude(k, v...)
		} else {
			qs = qs.Filter(k, v...)
		}
	}
	total, err = qs.Count()
	if err != nil {
		return
	}

	// order by:
	qs = qs.OrderBy(f.Order...)

	_, err = qs.Limit(f.Limit, f.Offset).All(contains, f.Fields...)

	return
}
