package hocon

import (
	"strings"
	"fmt"
)

type Content map[string]*Value

type ConfigObject struct {
	content *Content
}

func (co *ConfigObject) resolveKey(path string) (*ConfigObject, string) {
	prefix, key := pathPrefix(splitPath(path))
	resultKey := setObjectKey(prefix, co)
	return resultKey, key
}

func (co *ConfigObject) setString(path string, value string) {
	obj, key := co.resolveKey(path)
	(*obj.content)[key] = MakeStringValue(value)
}

func (co *ConfigObject) setInt(path string, value string) {
	obj, key := co.resolveKey(path)
	(*obj.content)[key] = MakeNumericValue(value)
}

func (co *ConfigObject) setValue(path string, vType ValueType, ref interface{}) {
	obj, key := co.resolveKey(path)
	(*obj.content)[key] = &Value{
		Type:     vType,
		RefValue: ref,
	}
}

func (co *ConfigObject) setObject(path string, value *ConfigObject) {
	co.setValue(path, ObjectType, value)
}

func (co *ConfigObject) setArray(path string, value *ConfigArray) {
	co.setValue(path, ArrayType, value)
}

func (co *ConfigObject) setCompoundString(path string, value *CompoundString) {
	co.setValue(path, CompoundStringType, value)
}

func (co *ConfigObject) setReference(path string, value string) {
	co.setValue(path, ReferenceType, value)
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

func (co *ConfigObject) getValue(path string) (res *Value) {
	if obj, key := traversePath(co, path); obj != nil {
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

func (co *ConfigObject) resolveStringReference(path string) string {
	if v := co.getValue(path); v != nil {
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
					result = co.resolveStringReference(referencePath(data.RefValue.(string))) + result
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

func (co *ConfigObject) GetString(path string) (res string) {
	res = co.resolveStringReference(path)
	return res
}

func (co *ConfigObject) GetInt(path string) (res int) {
	if obj, key := traversePath(co, path); obj != nil {
		if v, ok := (*obj.content)[key]; ok {
			res = v.RefValue.(int)
		}
	}
	return res
}

func (co *ConfigObject) GetObject(path string) (res *ConfigObject) {
	if obj, key := traversePath(co, path); obj != nil {
		if v, ok := (*obj.content)[key]; ok {
			res = v.RefValue.(*ConfigObject)
		}
	}
	return res
}

func (co *ConfigObject) GetArray(path string) (res *ConfigArray) {
	if obj, key := traversePath(co, path); obj != nil {
		if v, ok := (*obj.content)[key]; ok {
			res = v.RefValue.(*ConfigArray)
		}
	}
	return res
}

func (co *ConfigObject) GetKeys() []string {
	res := make([]string, len(*co.content))
	i := 0
	for k := range *co.content {
		res[i] = k
		i++
	}
	return res
}

func (co *ConfigObject) Merge(right *ConfigObject) {
	for k, v := range *right.content {
		if obj, ok := (*co.content)[k]; ok {
			if obj.Type == v.Type {
				switch obj.Type {
				case ObjectType:
					obj.RefValue.(*ConfigObject).Merge(v.RefValue.(*ConfigObject))
				case ArrayType:
					array := obj.RefValue.(*ConfigArray)
					array.content = append(array.content, v.RefValue.(*ConfigArray).content...)
					array.idx = array.idx + v.RefValue.(*ConfigArray).idx
				case StringType:
					fallthrough
				case CompoundStringType:
					fallthrough
				case NumericType:
					fallthrough
				case ReferenceType:
					obj.cloneFrom(v)
				}
			}
		} else {
			(*co.content)[k] = v
		}
	}
}
