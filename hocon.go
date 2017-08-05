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

type hocon struct {
	*parser.BaseHOCONListener
	stack       stack
	root        *ConfigObject
	compoundRef *CompoundString
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
	h.stack.Push(co)
	return h
}

func commonParse(p *parser.HOCONParser, h *hocon) (o *ConfigObject, err error) {
	p.AddParseListener(h)
	p.Hocon()
	res, _ := h.stack.Pop()
	o = res.(*ConfigObject)
	return o, err

}

func ParseHocon(stream antlr.CharStream) (o *ConfigObject, err error) {
	h := newHocon()
	ts := parser.NewHOCONLexer(stream)
	p := parser.NewHOCONParser(antlr.NewCommonTokenStream(ts, 0))
	return commonParse(p, h)
}

func ParseHoconString(data *string) (o *ConfigObject, err error) {
	h := newHocon()
	ts := parser.NewHOCONLexer(antlr.NewInputStream(*data))
	p := parser.NewHOCONParser(antlr.NewCommonTokenStream(ts, 0))
	return commonParse(p, h)
}

func ParseHoconFile(filename string) (o *ConfigObject, err error) {
	h := newHocon()
	if fs, err := antlr.NewFileStream(filename); err == nil {
		ts := parser.NewHOCONLexer(fs)
		p := parser.NewHOCONParser(antlr.NewCommonTokenStream(ts, 0))
		return commonParse(p, h)
	}
	return o, err
}