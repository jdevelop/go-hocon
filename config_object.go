package hocon

import (
	"strings"
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
	if obj := co.GetValue(path); obj != nil && obj.Type == ObjectType {
		obj.RefValue.(*ConfigObject).Merge(value)
	} else {
		co.setValue(path, ObjectType, value)
	}
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

func referencePath(path string) string {
	if strings.HasPrefix(path, "${") {
		return path[2:len(path)-1]
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

func (co *ConfigObject) GetValue(path string) (res *Value) {
	if obj, key := traversePath(co, path); obj != nil {
		if v, ok := (*obj.content)[key]; ok {
			res = v
		}
	}
	return res
}

func (co *ConfigObject) GetString(path string) (res string) {
	if obj, key := traversePath(co, path); obj != nil {
		if v, ok := (*obj.content)[key]; ok {
			switch v.Type {
			case StringType:
				fallthrough
			case ReferenceType:
				res = v.RefValue.(string)
			case CompoundStringType:
				res = ""
				for _, v := range v.RefValue.(*CompoundString).Value {
					res = v.RefValue.(string) + res
				}
			}
		}
	}
	return
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
			switch v.Type {
			case ObjectType:
				res = v.RefValue.(*ConfigObject)
			case CompoundStringType:
				res = v.RefValue.(*CompoundString).Value[0].RefValue.(*ConfigObject)
			}
		}
	}
	return res
}

func (co *ConfigObject) GetArray(path string) (res *ConfigArray) {
	if obj, key := traversePath(co, path); obj != nil {
		if v, ok := (*obj.content)[key]; ok {
			switch v.Type {
			case ArrayType:
				res = v.RefValue.(*ConfigArray)
			case CompoundStringType:
				res = v.RefValue.(*CompoundString).Value[0].RefValue.(*ConfigArray)
			}
		}
	}
	return
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

func (co *ConfigObject) Merge(ref ConfigInterface) {
	var right *ConfigObject
	switch ref.(type) {
	case *hocon:
		right = ref.(*hocon).root
	case *ConfigObject:
		right = ref.(*ConfigObject)
	}
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

func (co *ConfigObject) ResolveReferences() {

}
