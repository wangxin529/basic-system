package meta

import "reflect"

func (l *ListOption) AddToCondition(key string, val interface{}) {
	if val == nil || (reflect.ValueOf(val).Kind() == reflect.Ptr && reflect.ValueOf(val).IsNil()) {
		return
	}
	if l.Condition == nil {
		l.Condition = make(map[string]interface{})
	}
	l.Condition[key] = val
}

func isNull(val interface{}) (ok bool) {
	switch val.(type) {
	case string:
		if val.(string) == "" {
			return true
		}
	case int64:
		if val.(int64) == 0 {
			return true
		}
	}
	return false
}
