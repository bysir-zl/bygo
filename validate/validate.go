package validate

import (
	"github.com/bysir-zl/bygo/util"
	"strconv"
	"strings"
	"errors"
)

type rule struct {
	rules string // "need,len[:5],in(1,2),num[0:100]"
	msg   string // notice
}

type ValidateRule struct {
	filter []string // only field
	rules  map[string]rule
	need   []string // need field
}

func NewValidateRule() (r *ValidateRule) {
	r = &ValidateRule{
		filter: []string{},
		rules:  map[string]rule{},
	}
	return
}

func (p *ValidateRule) Only(fields ...string) (r *ValidateRule) {
	p.filter = fields
	return p
}
func (p *ValidateRule) Need(fields ...string) (r *ValidateRule) {
	p.need = fields
	return p
}

func (p *ValidateRule) Rule(field string, ru string, msg string) (r *ValidateRule) {
	p.rules[field] = rule{ru, msg}
	return p
}

func Need(fields ...string) (r *ValidateRule) {
	return NewValidateRule().Need(fields...)
}
func Only(fields ...string) (r *ValidateRule) {
	return NewValidateRule().Only(fields...)
}
func Rule(field string, ru string, msg string) (r *ValidateRule) {
	p := NewValidateRule()
	p.rules[field] = rule{ru, msg}
	return p
}

func (p *ValidateRule) validateValue(value string, rule rule) (isOk bool, notice string) {
	if rule.rules == "" {
		isOk = false
		notice = rule.msg
		return
	}

	isOk = true
	notice = rule.msg
	for _, r := range strings.Split(rule.rules, ",") {

		if r == "need" {
			isOk = value != ""
			if !isOk {
				return
			}
		} else if strings.Index(r, "len") == 0 {
			r = r[3:]
			ranger := r[1 : len(r) - 1]
			mm := strings.Split(ranger, ":")
			minLenString := mm[0]
			maxLenString := mm[1]

			strLen := len(value)
			if minLenString != "" {
				if r[0] == '(' {
					if min, _ := strconv.Atoi(minLenString); strLen <= min {
						isOk = false
						return
					}
				} else {
					if min, _ := strconv.Atoi(minLenString); strLen < min {
						isOk = false
						return
					}
				}

			}
			if maxLenString != "" {
				if r[len(r) - 1] == ')' {
					if max, _ := strconv.Atoi(maxLenString); strLen >= max {
						isOk = false
						return
					}
				} else {
					if max, _ := strconv.Atoi(maxLenString); strLen > max {
						isOk = false
						return
					}
				}
			}
		} else if strings.Index(r, "num") == 0 {
			num, err := strconv.ParseFloat(value, 64)
			if err != nil {
				isOk = false
				return
			}
			r = r[3:]
			ranger := r[1 : len(r) - 1]
			mm := strings.Split(ranger, ":")
			minLenString := mm[0]
			maxLenString := mm[1]

			if minLenString != "" {
				if r[0] == '(' {
					if min, _ := strconv.ParseFloat(minLenString, 64); num <= min {
						isOk = false
						return
					}
				} else {
					if min, _ := strconv.ParseFloat(minLenString, 64); num < min {
						isOk = false
						return
					}
				}
			}
			if maxLenString != "" {
				if r[len(r) - 1] == ')' {
					if max, _ := strconv.ParseFloat(maxLenString, 64); num >= max {
						isOk = false
						return
					}
				} else {
					if max, _ := strconv.ParseFloat(maxLenString, 64); num > max {
						isOk = false
						return
					}
				}
			}
		} else if strings.Index(r, "in") == 0 {
			r = r[2:]
			ranger := r[1 : len(r) - 1]
			ok := util.ItemInArray(value, strings.Split(ranger, ","))
			if !ok {
				isOk = false
				return
			}
		}

	}
	isOk = true
	return
}

// to validata this map
func (p *ValidateRule) Validate(m map[string]string) (err error) {
	if len(p.filter) != 0 {
		ok, m := util.ArrayInArray(util.GetMapKey(m), p.filter)
		if !ok {
			err = errors.New("field :" + m + " is not allowed")
			return
		}
	}

	if len(p.need) != 0 {
		ok, m := util.ArrayInArray(p.need, util.GetMapKey(m))
		if !ok {
			err = errors.New("field :" + m + " must be seted")
			return
		}
	}

	for field, r := range p.rules {
		value := m[field]
		ok, m := p.validateValue(value, r)
		if !ok {
			err = errors.New(m)
			return
		}
	}

	return
}