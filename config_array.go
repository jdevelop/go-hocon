package hocon

import "fmt"

type ConfigArray struct {
	hocon   *hocon
	idx     int
	content []*Value
}

func (ca *ConfigArray) setString(path string, value string) {
	ca.append(MakeStringValue(value))
}

func (ca *ConfigArray) setCompoundString(path string, value *CompoundString) {
	ca.append(MakeCompoundStringValue(value))
}

func (ca *ConfigArray) setInt(path string, value string) {
	ca.append(MakeNumericValue(value))
}

func (ca *ConfigArray) setObject(path string, value *ConfigObject) {
	ca.append(MakeObjectValue(value))
}

func (ca *ConfigArray) setArray(path string, value *ConfigArray) {
	ca.append(MakeArrayValue(value))
}

func (ca *ConfigArray) setReference(path string, value string) {
}

func (ca *ConfigArray) setValue(path string, t ValueType, value interface{}) {
}

func NewConfigArray(hocon *hocon) *ConfigArray {
	co := ConfigArray{
		hocon:   hocon,
		idx:     0,
		content: make([]*Value, 1),
	}
	return &co
}

func (ca *ConfigArray) append(v *Value) {
	size := len(ca.content)
	if ca.idx == size {
		tmp := make([]*Value, size*2)
		copy(tmp, ca.content)
		ca.content = tmp
	}
	ca.content[ca.idx] = v
	ca.idx++
}

func (a *ConfigArray) GetString(idx int) string {
	ref := a.content[idx]
	switch ref.Type {
	case CompoundStringType:
		var result string = ""
		cs := ref.RefValue.(*CompoundString)
		for _, data := range cs.Value {
			switch data.Type {
			case StringType:
				fallthrough
			case ReferenceType:
				result = data.RefValue.(string) + result
			case NumericType:
				result = fmt.Sprintf("%1d", data.RefValue.(int)) + result
			}
		}
		ref.RefValue = result
		ref.Type = StringType
		return result
	default:
		return ref.RefValue.(string)
	}
}

func (ca *ConfigArray) GetInt(idx int) int {
	return ca.content[idx].RefValue.(int)
}

func (ca *ConfigArray) GetObject(idx int) *ConfigObject {
	return ca.content[idx].RefValue.(*ConfigObject)
}

func (ca *ConfigArray) GetArray(idx int) *ConfigArray {
	return ca.content[idx].RefValue.(*ConfigArray)
}

func (ca *ConfigArray) GetSize() int {
	return ca.idx
}
