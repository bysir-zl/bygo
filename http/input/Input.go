package input

import (
    "net/http"
    "github.com/bysir-zl/bygo/util"
    "errors"
)

type Input struct {
    valueMap map[string]string
    Request  *http.Request
}

func (i *Input) All() map[string]string {
    if i.valueMap == nil {
        i.valueMap = make(map[string]string);

        i.Request.ParseForm()
        i.Request.ParseMultipartForm(10240000)

        for key, value := range i.Request.Form {
            i.valueMap[key] = value[0];
        }
    }

    return i.valueMap;
}

func (i *Input) Get(key string) string {
    if i.valueMap == nil {
        i.All()
    }
    return i.valueMap[key]
}

func (i *Input) Gets(keys ...string) map[string]string {
    if i.valueMap == nil {
        i.All()
    }
    mapper := map[string]string{};

    for _, key := range keys {
        value := i.valueMap[key]
        if value != "" {
            mapper[key] = value
        }
    }

    return mapper
}

func (i *Input) Set(key string, value string) {
    if i.valueMap == nil {
        i.All()
    }
    i.valueMap[key] = value
}

func (i *Input) SetAll(input map[string]string) {
    i.valueMap = input
}

func (i *Input) SetToObj(obj interface{}, useTag string) (field []string) {
    if i.valueMap == nil {
        i.All()
    }
    mapper := map[string]interface{}{};
    for k, v := range i.valueMap {
        mapper[k] = v
    }

    field = util.MapToObj(obj, mapper, useTag)
    return
}

// 只填充指定的字段
func (i *Input) SetToObjWithAllowFilter(obj interface{}, useTag string, fields ...string) (field []string) {
    if i.valueMap == nil {
        i.All()
    }
    mapper := map[string]interface{}{};
    for k, v := range i.valueMap {
        if !util.ItemInArray(k, fields) {
            continue
        }
        mapper[k] = v
    }

    field = util.MapToObj(obj, mapper, useTag)
    return
}
// 只拒绝指定的字段
func (i *Input) SetToObjWithDenyFilter(obj interface{}, useTag string, fields ...string) (field []string) {
    if i.valueMap == nil {
        i.All()
    }
    field = []string{}
    mapper := map[string]interface{}{};
    for k, v := range i.valueMap {
        if util.ItemInArray(k, fields) {
            continue
        }
        mapper[k] = v
        field = append(field, k)
    }

    util.MapToObj(obj, mapper, useTag)
    return
}

//是否只有指定字段，通常用于控制传入值的过滤与判断是否合法
func (i *Input) IsOnlyField(fields ...string) (bool) {
    if i.valueMap == nil {
        i.All()
    }
    for k, _ := range i.valueMap {
        if !util.ItemInArray(k, fields) {
            return false
        }
    }

    return true
}

func (i *Input) Validate(rule *ValidateRule) (err error) {
    if i.valueMap == nil {
        i.All()
    }
    if len(rule.filter) != 0 {
        ok, m := util.ArrayInArray(util.GetMapKey(i.valueMap), rule.filter)
        if !ok {
            err = errors.New("field :" + m + " is not allowed");
            return
        }
    }

    for field, r := range rule.rules {
        value := i.valueMap[field]
        ok, m := rule.ValidateValue(value, r)
        if !ok {
            err = errors.New(m);
            return
        }
    }

    return
}


