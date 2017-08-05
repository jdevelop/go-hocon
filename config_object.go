package hocon

import (
	"strings"
	"fmt"
)

type Content map[string]*Value

type ConfigObject struct {
	content *Content
}

func (o *ConfigObject) resolveKey(path string) (*ConfigObject, string) {
	prefix, key := pathPrefix(splitPath(path))
	resultKey := setObjectKey(prefix, o)
	return resultKey, key
}

func (o *ConfigObject) setString(path string, value string) {
	obj, key := o.resolveKey(path)
	(*obj.content)[key] = MakeStringValue(value)
}

func (o *ConfigObject) setInt(path string, value string) {
	obj, key := o.resolveKey(path)
	(*obj.content)[key] = MakeNumericValue(value)
}

func (o *ConfigObject) setValue(path string, vType ValueType, ref interface{}) {
	obj, key := o.resolveKey(path)
	(*obj.content)[key] = &Value{
		Type:     vType,
		RefValue: ref,
	}
}

func (o *ConfigObject) setObject(path string, value *ConfigObject) {
	o.setValue(path, ObjectType, value)
}

func (o *ConfigObject) setArray(path string, value *ConfigArray) {
	o.setValue(path, ArrayType, value)
}

func (o *ConfigObject) setCompoundString(path string, value *CompoundString) {
	o.setValue(path, CompoundStringType, value)
}

func (o *ConfigObject) setReference(path string, value string) {
	o.setValue(path, ReferenceType, value)
}

func setObjectKey(keys []string, obj *ConfigObject) *ConfigObject {
	for _, key := range keys {
		if v, exists := (*obj.content)[key]; exists {
			switch v.Type {
			case ObjectType:
				obj = v.RefValue.(*ConfigObject)
			default:
				panic("Wrong path")
			}
			continue
		}

		newObject := NewConfigObject()
		(*obj.content)[key] = &Value{
			Type:     ObjectType,
			RefValue: newObject,
		}
		obj = newObject
	}
	return obj
}

func traversePath(o *ConfigObject, path string) (*ConfigObject, string) {
	obj := o
	paths := strings.Split(path, ".")
	for _, p := range paths[:len(paths)-1] {
		if d := (*obj.content)[p]; d == nil {
			return nil, ""
		} else {
			switch d.Type {
			case ObjectType:
				obj = d.RefValue.(*ConfigObject)
			default:
				return nil, ""
			}
		}
	}
	return obj, paths[len(paths)-1]
}

func (o *ConfigObject) getValue(path string) (res *Value) {
	if obj, key := traversePath(o, path); obj != nil {
		if v, ok := (*obj.content)[key]; ok {
			res = v
		}
	}
	return res
}

func referencePath(path string) string {
	if strings.HasPrefix(path, "${") {
		return path[2:len(path)-1]
	} else {
		return path
	}
}

func (o *ConfigObject) resolveStringReference(path string) string {
	if v := o.getValue(path); v != nil {
		switch v.Type {
		case StringType:
			return v.RefValue.(string)
		case NumericType:
			return fmt.Sprintf("%1d", v.RefValue.(int))
		case CompoundStringType:
			var result string = ""
			cs := v.RefValue.(*CompoundString)
			for _, data := range cs.Value {
				switch data.Type {
				case StringType:
					result = data.RefValue.(string) + result
				case ReferenceType:
					result = o.resolveStringReference(referencePath(data.RefValue.(string))) + result
				}
			}
			v.RefValue = result
			v.Type = StringType
			return result
		default:
			return path
		}
	} else {
		return path
	}
}

// =====================================================================

func NewConfigObject() *ConfigObject {
	m := make(Content)
	co := ConfigObject{
		content: &m,
	}
	return &co
}

func (o *ConfigObject) GetString(path string) (res string) {
	res = o.resolveStringReference(path)
	return res
}

func (o *ConfigObject) GetInt(path string) (res int) {
	if obj, key := traversePath(o, path); obj != nil {
		if v, ok := (*obj.content)[key]; ok {
			res = v.RefValue.(int)
		}
	}
	return res
}

func (o *ConfigObject) GetObject(path string) (res *ConfigObject) {
	if obj, key := traversePath(o, path); obj != nil {
		if v, ok := (*obj.content)[key]; ok {
			res = v.RefValue.(*ConfigObject)
		}
	}
	return res
}

func (o *ConfigObject) GetArray(path string) (res *ConfigArray) {
	if obj, key := traversePath(o, path); obj != nil {
		if v, ok := (*obj.content)[key]; ok {
			res = v.RefValue.(*ConfigArray)
		}
	}
	return res
}

func (o *ConfigObject) GetKeys() []string {
	res := make([]string, len(*o.content))
	i := 0
	for k, _ := range *o.content {
		res[i] = k
		i++
	}
	return res
}
