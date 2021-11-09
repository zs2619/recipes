package recipesgo

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type JSONOP int8

const (
	NoneOP  JSONOP = 0
	SetOP   JSONOP = 1
	UnSetOP JSONOP = 2
	PushOP  JSONOP = 3
	PullOP  JSONOP = 4
	IncOP   JSONOP = 5
)

func ConvertJsonOp(op string) (JSONOP, error) {
	switch op {
	case "s":
		return SetOP, nil
	case "us":
		return UnSetOP, nil
	case "i":
		return IncOP, nil
	case "pl":
		return PullOP, nil
	case "ps":
		return PushOP, nil
	default:
		return NoneOP, fmt.Errorf("convertJsonOp " + op)
	}

}

type DiffSymbol struct {
	Path     string
	OP       JSONOP
	OldValue interface{}
	NewValue interface{}
}

func DiffJsonValue(src, dest reflect.Value, prefix string) *map[string]DiffSymbol {
	ret := make(map[string]DiffSymbol)

	switch src.Kind() {
	case reflect.String:
		if src.String() == dest.String() {
			return &ret
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
		if src.Float() == dest.Float() {
			return &ret
		}
	case reflect.Bool:
		if src.Bool() == dest.Bool() {
			return &ret
		}
	case reflect.Slice:
		if dest.Len() == 0 && src.Len() == 0 {
			return &ret
		}
		if dest.Len() == src.Len() {
			for i := 0; i < dest.Len(); i++ {
				v1 := src.Index(i).Elem()
				v2 := dest.Index(i).Elem()
				if v1.Kind() == v2.Kind() {
					var m *map[string]DiffSymbol
					if v1.Kind() == reflect.Map {
						m = diffJSObj(v1.Interface().(*map[string]interface{}), v2.Interface().(*map[string]interface{}), prefix+"."+strconv.Itoa(i))
					} else {
						m = DiffJsonValue(v1, v2, prefix+"."+strconv.Itoa(i))
					}
					if m != nil {
						for k, v := range *m {
							ret[k] = v
						}
					}
					return &ret
				}
			}
		}
	default:
		panic(fmt.Sprintf("%v %v type error", src, dest))
	}

	ret[prefix] = DiffSymbol{Path: prefix, OP: SetOP, OldValue: src.Interface(), NewValue: dest.Interface()}

	return &ret
}

func diffJSObj(src, dest *map[string]interface{}, position string) *map[string]DiffSymbol {
	if src == nil {
		src = new(map[string]interface{})
	}
	if dest == nil {
		dest = new(map[string]interface{})
	}
	if len(*src) == 0 && len(*dest) == 0 {
		return nil
	}
	ret := make(map[string]DiffSymbol)
	if len(*src) == 0 && len(*dest) > 0 {
		for k, v := range *dest {
			ret[k] = DiffSymbol{Path: k, OP: SetOP, OldValue: nil, NewValue: v}
		}
		return &ret
	}
	if len(*src) > 0 && len(*dest) == 0 {
		for k, v := range *src {
			ret[k] = DiffSymbol{Path: k, OP: UnSetOP, OldValue: v, NewValue: nil}
		}
		return &ret
	}

	for k, dv := range *dest {
		sv, ok := (*src)[k]
		path := position + "." + k
		if ok {
			svv := reflect.ValueOf(sv)
			dvv := reflect.ValueOf(dv)
			if svv.Kind() == dvv.Kind() {
				var ms *map[string]DiffSymbol
				if svv.Kind() == reflect.Map {
					tsv := sv.(map[string]interface{})
					tdv := dv.(map[string]interface{})
					ms = diffJSObj(&tsv, &tdv, path)
				} else {
					ms = DiffJsonValue(svv, dvv, path)
				}
				if ms != nil {
					for k, v := range *ms {
						ret[k] = v
					}
				}
				continue
			}
		}
		ret[path] = DiffSymbol{Path: path, OP: SetOP, OldValue: sv, NewValue: dv}
	}

	for k, sv := range *src {
		_, ok := (*dest)[k]
		if !ok {
			path := position + "." + k
			ret[path] = DiffSymbol{Path: path, OP: UnSetOP, OldValue: sv, NewValue: nil}
		}
	}

	return &ret
}

/**深度遍历 json序列化的对象 找出差异*/
func DiffJSON(src, dest *map[string]interface{}) *map[string]DiffSymbol {
	return diffJSObj(src, dest, "$")
}

/** json path 格式$.field.field*/
func JsonPathGet(jsonObj *map[string]interface{}, path string) (interface{}, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("path error")
	}
	fields := strings.Split(path, ".")
	if fields[0] != "$" {
		return nil, fmt.Errorf("path error")
	}

	var retv interface{}
	tmpObj := *jsonObj
	for i := 1; i < len(fields); i++ {
		field := fields[i]
		tmp, ok := (tmpObj)[field]
		if !ok {
			return nil, fmt.Errorf("json object no (%s) field", path)
		}
		if i == len(fields)-1 {
			retv = tmp
			return retv, nil
		} else {
			tmpObj, ok = (tmp).(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("json object no (%s) field", path)
			}
		}
	}
	return nil, fmt.Errorf("path error")
}

func formatValue(param interface{}) interface{} {
	var p interface{}
	switch param.(type) {
	case float64:
		p = float64(param.(float64))
	case float32:
		p = float64(param.(float32))
	case int:
		p = float64(param.(int))
	case int8:
		p = float64(param.(int8))
	case int16:
		p = float64(param.(int16))
	case int32:
		p = float64(param.(int32))
	case int64:
		p = float64(param.(int64))
	case uint8:
		p = float64(param.(uint8))
	case uint16:
		p = float64(param.(uint16))
	case uint32:
		p = float64(param.(uint32))
	case uint64:
		p = float64(param.(uint64))
	default:
		p = param
	}
	return p
}

type jsNodeOP interface {
	op(*map[string]interface{}, string) error
}
type jsNodeSetOP struct {
	param interface{}
}

func (this *jsNodeSetOP) op(pnode *map[string]interface{}, field string) error {
	(*pnode)[field] = formatValue(this.param)
	return nil
}

type jsNodeUnSetOP struct {
}

func (this *jsNodeUnSetOP) op(pnode *map[string]interface{}, field string) error {
	delete(*pnode, field)
	return nil
}

type jsNodePushOP struct {
	param interface{}
}

func (this *jsNodePushOP) op(pnode *map[string]interface{}, field string) error {
	v := formatValue(this.param)
	s, ok := ((*pnode)[field]).([]interface{})
	if !ok {
		(*pnode)[field] = []interface{}{v}
	} else {
		(*pnode)[field] = append(s, v)
	}
	return nil
}

type jsNodePullOP struct {
	param interface{}
}

func (this *jsNodePullOP) op(pnode *map[string]interface{}, field string) error {
	p := formatValue(this.param)
	s, ok := ((*pnode)[field]).([]interface{})
	if ok {
		for k, v := range s {
			if reflect.DeepEqual(v, p) {
				s = append(s[:k], s[k+1:]...)
			}
		}
		(*pnode)[field] = s
	}
	return nil
}

type jsNodeIncOP struct {
	param float64
}

func (this *jsNodeIncOP) op(pnode *map[string]interface{}, field string) error {
	v, ok := (*pnode)[field]
	if !ok {
		(*pnode)[field] = reflect.Zero(reflect.ValueOf(this.param).Type()).Interface()
		v = (*pnode)[field]
	}
	switch v.(type) {
	case float64:
		(*pnode)[field] = v.(float64) + this.param
	case float32:
		(*pnode)[field] = float64(v.(float32)) + this.param
	case int:
		(*pnode)[field] = float64(v.(int)) + this.param
	case int8:
		(*pnode)[field] = float64(v.(int8)) + this.param
	case int16:
		(*pnode)[field] = float64(v.(int16)) + this.param
	case int32:
		(*pnode)[field] = float64(v.(int32)) + this.param
	case int64:
		(*pnode)[field] = float64(v.(int64)) + this.param
	case uint8:
		(*pnode)[field] = float64(v.(uint8)) + this.param
	case uint16:
		(*pnode)[field] = float64(v.(uint16)) + this.param
	case uint32:
		(*pnode)[field] = float64(v.(uint32)) + this.param
	case uint64:
		(*pnode)[field] = float64(v.(uint64)) + this.param
	default:
		return fmt.Errorf("JSNodeIncOP  %s type error", field)
	}
	return nil
}

func jsNodeOPFactroy(op JSONOP, param interface{}) (jsNodeOP, error) {
	switch op {
	case SetOP:
		return &jsNodeSetOP{param: param}, nil
	case UnSetOP:
		return &jsNodeUnSetOP{}, nil
	case PushOP:
		return &jsNodePushOP{param: param}, nil
	case PullOP:
		return &jsNodePullOP{param: param}, nil
	case IncOP:
		switch param.(type) {
		case float64:
			return &jsNodeIncOP{param: float64(param.(float64))}, nil
		case float32:
			return &jsNodeIncOP{param: float64(param.(float32))}, nil
		case int:
			return &jsNodeIncOP{param: float64(param.(int))}, nil
		case int8:
			return &jsNodeIncOP{param: float64(param.(int8))}, nil
		case int16:
			return &jsNodeIncOP{param: float64(param.(int16))}, nil
		case int32:
			return &jsNodeIncOP{param: float64(param.(int32))}, nil
		case int64:
			return &jsNodeIncOP{param: float64(param.(int64))}, nil
		case uint8:
			return &jsNodeIncOP{param: float64(param.(uint8))}, nil
		case uint16:
			return &jsNodeIncOP{param: float64(param.(uint16))}, nil
		case uint32:
			return &jsNodeIncOP{param: float64(param.(int32))}, nil
		case uint64:
			return &jsNodeIncOP{param: float64(param.(uint64))}, nil
		default:
			return nil, fmt.Errorf("jsNodeOPFactroy  param %s", param)
		}

	default:
		panic("op error")
	}
}

/** json 设置属性 path 格式$.field.field*/
func JsonPathOP(jsonObj map[string]interface{}, path string, op JSONOP, value interface{}) (*map[string]interface{}, error) {
	if len(path) <= 1 {
		return nil, fmt.Errorf("path error")
	}
	fields := strings.Split(path, ".")
	if fields[0] != "$" {
		return nil, fmt.Errorf("path error")
	}
	f, err := jsNodeOPFactroy(op, value)
	if err != nil {
		return nil, err
	}

	tmpObj := &jsonObj
	for i := 1; i < len(fields); i++ {
		field := fields[i]
		if i == len(fields)-1 {
			f.op(tmpObj, field)
			return &jsonObj, nil
		}
		tmp, ok := (*tmpObj)[field]
		if ok {
			t, ok := tmp.(map[string]interface{})
			if ok {
				tmpObj = &t
			} else {
				v := make(map[string]interface{})
				(*tmpObj)[field] = v
				tmpObj = &v
			}

		} else {
			v := make(map[string]interface{})
			(*tmpObj)[field] = v
			tmpObj = &v
		}
	}
	return nil, fmt.Errorf("JsonPathOP error")
}
func CheckJsonOpJs(opjs *string) error {
	var op map[string]map[string]interface{}
	err := json.Unmarshal([]byte(*opjs), &op)
	if err != nil {
		return fmt.Errorf("CheckJsonOpJs")
	}
	set := make(map[string]bool)
	for k, v := range op {
		_, err = ConvertJsonOp(k)
		if err != nil {
			return fmt.Errorf("CheckJsonOpJs (%s) op error", k)
		}
		for k1, _ := range v {
			_, ok := set[k1]
			if ok {
				return fmt.Errorf("CheckJsonOpJs (%s)key %s  error", k1, *opjs)
			}
			set[k1] = true
		}
	}
	return nil
}
