package db

import (
	"errors"
	"github.com/bysir-zl/bygo/util"
	"strconv"
	"strings"
)

//model执行规则
type ModelRules struct {
	whereRules       map[string]string
	whereRulesExtend map[string][]interface{}

	orderRules       map[string]int
	fieldsRules      []fieldsR
	fieldsList       []string
	pageRules        pageRules
}

type pageRules struct {
	WithOutPage  bool
	WithOutTotal bool
}

type fieldsR struct {
	name   string
	rule   string
	errMsg string
}

func (p ModelRules) GetFieldsList() []string {
	return p.fieldsList
}
func (p ModelRules) GetFieldsRules() []fieldsR {
	return p.fieldsRules
}
func (p ModelRules) GetOrderRules() map[string]int {
	return p.orderRules
}
func (p ModelRules) GetWhereRules() map[string]string {
	return p.whereRules
}
func (p ModelRules) GetPageRules() pageRules {
	return p.pageRules
}

func (p ModelRules) FieldsList(f ...string) ModelRules {
	p.fieldsList = f
	return p
}

func (p ModelRules) FieldsRules(field string, rule string, errMsg string) ModelRules {
	if p.fieldsRules == nil {
		p.fieldsRules = []fieldsR{}
	}
	p.fieldsRules = append(p.fieldsRules, fieldsR{name: field, rule: rule, errMsg: errMsg})
	return p
}

// desc 1:both , 2:asc , 3:desc
func (p ModelRules) OrderRules(key string, desc int) ModelRules {
	if p.orderRules == nil {
		p.orderRules = map[string]int{}
	}

	p.orderRules[key] = desc
	return p
}

// key 是前端传的key , rule 是条件字符串(eg:`Id` = ?);
func (p ModelRules) WhereRules(key string, rule string) ModelRules {
	if p.whereRules == nil {
		p.whereRules = map[string]string{}
	}
	p.whereRules[key] = rule
	return p
}

// 扩展条件
// 必要条件, 每次查询都会加上的条件; 这一般用于权限验证
func (p ModelRules) WhereRulesExtend(rule string, values ...interface{}) ModelRules {
	if p.whereRulesExtend == nil {
		p.whereRulesExtend = map[string][]interface{}{}
	}
	p.whereRulesExtend[rule] = values
	return p
}

// withOutPage 是否输出page数据 ; withOutTotal 如果要输出page数据,是否统计Total
func (p ModelRules) PageRules(withOutPage bool, withOutTotal bool) ModelRules {
	p.pageRules = pageRules{WithOutPage: withOutPage, WithOutTotal: withOutTotal}
	return p
}

func NewModelRules() ModelRules {
	return ModelRules{}
}

/** Util **/

// 根据whereString(Id:123|321) 生成 conditionString=>values []string
func (m ModelRules) CheckWhereString(whereString string, need bool) (whereMapper map[string]([]interface{}), err error) {
	whereMapper = map[string]([]interface{}){}
	whereRules := m.GetWhereRules()
	whereRulesExtend := m.whereRulesExtend

	for wk, wv := range whereRulesExtend {
		whereMapper[wk] = wv
	}

	if whereString == "" {
		if need {
			err = errors.New("do you forget Where param ?")
		}
		return
	}
	whereKvs := strings.Split(whereString, ",")

	for _, kvStr := range whereKvs {
		kv := strings.Split(kvStr, ":")
		if len(kv) != 2 || kv[0] == "" || kv[1] == "" {
			err = errors.New("Where format error")
			return
		}
		key := kv[0]

		// whereRules中没指定此字段
		if whereRules[key] == "" {
			err = errors.New("server refuse select by key : " + key)
			return
		}
		conditionString := whereRules[key]
		values := strings.Split("|", kv[1])

		count := strings.Count(conditionString, "?")
		if count != len(values) {
			err = errors.New("server rule '" + key + "' condition params count is " + strconv.Itoa(count))
			return
		}

		valueInterfaces := []interface{}{}
		for v := range values {
			valueInterfaces = append(valueInterfaces, v)
		}

		whereMapper[conditionString] = valueInterfaces
	}

	return
}

func (m ModelRules) CheckOrderString(orderString string) (orderMapper map[string]bool, err error) {
	orderMapper = map[string]bool{}
	if orderString == "" {
		return
	}
	orderRules := m.GetOrderRules()

	// orderString=Id:DESC,Time
	orderKvs := strings.Split(orderString, ",")
	for _, kvStr := range orderKvs {
		kv := strings.Split(kvStr, ":")
		key := kv[0]
		if key == "" {
			err = errors.New("do you forget Order param ?")
			return
		}

		desc := false
		if len(kv) == 2 && kv[1] == "DESC" {
			desc = true
		}

		canOrder := orderRules[key]
		if canOrder == 0 {
			err = errors.New("server refuse order by key : " + key)

			return
		}
		if desc {
			if canOrder == 2 {
				err = errors.New("server refuse order desc by key : " + key)

				return
			}
		} else {
			if canOrder == 3 {
				err = errors.New("server refuse order asc by key : " + key)
				return
			}
		}

		orderMapper[key] = desc
	}
	return
}

// todo 参数验证还没有做
func (m ModelRules) CheckFields(fields []string) (err error) {
	onlyFields := m.GetFieldsList()
	if onlyFields == nil || len(onlyFields) == 0 {
		err = errors.New("server refuse update any field")
		return
	}
	// 如果传入了*,就不进行检查
	if onlyFields[0] != "*" {
		if in, msg := util.ArrayInArray(fields, onlyFields); !in {
			err = errors.New("server refuse update " + msg + " field")
			return
		}
	}
	return
}
