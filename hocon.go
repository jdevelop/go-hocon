package hocon

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/jdevelop/go-hocon/parser"
	"strconv"
)

type ValueType int

const (
	StringType         ValueType = iota
	CompoundStringType
	ReferenceType
	NumericType
	ObjectType
	ArrayType
)

type Value struct {
	Type     ValueType
	RefValue interface{}
}

func (v *Value) cloneFrom(that *Value) {
	v.RefValue = that.RefValue
	v.Type = that.Type
}

type hocon struct {
	*parser.BaseHOCONListener
	stack       stack
	root        *ConfigObject
	compoundRef *CompoundString
	refs        map[string]*Value
}

type ConfigInterface interface {
	GetValue(path string) *Value
	GetString(path string) string
	GetInt(path string) int
	GetObject(path string) *ConfigObject
	GetArray(path string) *ConfigArray
	GetKeys() []string
	Merge(ref ConfigInterface)
	ResolveReferences()
}

type CompoundString struct {
	Value []*Value
}

func (cv *CompoundString) addString(src string) {
	cv.Value = append(cv.Value, MakeStringValue(src))
}

func (cv *CompoundString) addReference(src string) {
	cv.Value = append(cv.Value, MakeReferenceValue(src))
}

type valueSetter interface {
	setString(name string, value string)
	setCompoundString(name string, value *CompoundString)
	setReference(name string, value string)
	setInt(name string, value string)
	setObject(name string, value *ConfigObject)
	setArray(name string, value *ConfigArray)
}

func MakeReferenceValue(src string) *Value {
	return &Value{
		Type:     ReferenceType,
		RefValue: src,
	}
}

func MakeCompoundStringValue(src *CompoundString) *Value {
	return &Value{
		Type:     CompoundStringType,
		RefValue: src,
	}
}

func MakeStringValue(src string) *Value {
	return &Value{
		Type:     StringType,
		RefValue: src,
	}
}

func MakeNumericValue(src string) *Value {
	val, _ := strconv.Atoi(src)
	return &Value{
		Type:     NumericType,
		RefValue: val,
	}
}

func MakeObjectValue(src *ConfigObject) *Value {
	return &Value{
		Type:     ObjectType,
		RefValue: src,
	}
}

func MakeArrayValue(src *ConfigArray) *Value {
	return &Value{
		Type:     ArrayType,
		RefValue: src,
	}
}

func newHocon() *hocon {
	h := new(hocon)
	h.stack = *NewStack()
	co := NewConfigObject()
	h.root = co
	h.refs = make(map[string]*Value)
	h.stack.Push(co)
	return h
}

func commonParse(p *parser.HOCONParser, h *hocon) (ConfigInterface, error) {
	p.AddParseListener(h)
	p.Hocon()
	return h, nil
}

func ParseHocon(stream antlr.CharStream) (o ConfigInterface, err error) {
	h := newHocon()
	ts := parser.NewHOCONLexer(stream)
	p := parser.NewHOCONParser(antlr.NewCommonTokenStream(ts, 0))
	return commonParse(p, h)
}

func ParseHoconString(data *string) (o ConfigInterface, err error) {
	h := newHocon()
	ts := parser.NewHOCONLexer(antlr.NewInputStream(*data))
	p := parser.NewHOCONParser(antlr.NewCommonTokenStream(ts, 0))
	return commonParse(p, h)
}

func ParseHoconFile(filename string) (o ConfigInterface, err error) {
	h := newHocon()
	if fs, err := antlr.NewFileStream(filename); err == nil {
		ts := parser.NewHOCONLexer(fs)
		p := parser.NewHOCONParser(antlr.NewCommonTokenStream(ts, 0))
		return commonParse(p, h)
	}
	return o, err
}

func (hc *hocon) GetValue(path string) *Value {
	return hc.root.GetValue(path)
}

func (hc *hocon) GetString(path string) string {
	return hc.root.GetString(path)
}

func (hc *hocon) GetInt(path string) int {
	return hc.root.GetInt(path)
}

func (hc *hocon) GetObject(path string) *ConfigObject {
	return hc.root.GetObject(path)
}

func (hc *hocon) GetArray(path string) *ConfigArray {
	return hc.root.GetArray(path)
}

func (hc *hocon) GetKeys() []string {
	return hc.root.GetKeys()
}

func (hc *hocon) Merge(ref ConfigInterface) {
	hc.root.Merge(ref)
	for k, v := range ref.(*hocon).refs {
		hc.refs[k] = v
	}
}

func (hc *hocon) ResolveReferences() {

	for ; len(hc.refs) > 0; {
		processed := false
		for k, v := range hc.refs {
			val := hc.GetValue(k)
			cleanup := func() {
				v.cloneFrom(val)
				delete(hc.refs, k)
				processed = true
			}
			if val == nil {
				continue
			}
			switch val.Type {
			case ReferenceType:
				continue
			case CompoundStringType:
				hasRef := false
				for _, v := range val.RefValue.(*CompoundString).Value {
					if v.Type == ReferenceType {
						hasRef = true
						break
					}
				}
				if !hasRef {
					cleanup()
				}
			default:
				cleanup()
			}
		}
		if !processed {
			break
		}
	}
}
