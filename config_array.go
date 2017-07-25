package hocon

type ConfigArray struct {
	idx     int
	content []*Value
}

func (c *ConfigArray) setString(path string, value string) {
	c.append(MakeStringValue(value))
}

func (c *ConfigArray) setInt(path string, value string) {
	c.append(MakeNumericValue(value))
}

func (c *ConfigArray) setObject(path string, value *ConfigObject) {
	c.append(MakeObjectValue(value))
}

func (c *ConfigArray) setArray(path string, value *ConfigArray) {
	c.append(MakeArrayValue(value))
}

func (c *ConfigArray) setReference(path string, value string) {
}

func NewConfigArray() *ConfigArray {
	co := ConfigArray{
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
	return a.content[idx].RefValue.(string)
}

func (a *ConfigArray) GetInt(idx int) int {
	return a.content[idx].RefValue.(int)
}

func (a *ConfigArray) GetObject(idx int) *ConfigObject {
	return a.content[idx].RefValue.(*ConfigObject)
}

func (a *ConfigArray) GetArray(idx int) *ConfigArray {
	return a.content[idx].RefValue.(*ConfigArray)
}

func (a *ConfigArray) GetSize() int {
	return a.idx
}
