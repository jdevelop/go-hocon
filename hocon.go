package hocon

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/jdevelop/go-hocon/parser"
	"strconv"
	"strings"
)

type ValueType int

const (
	StringType  ValueType = iota
	NumericType
	ObjectType
	ArrayType
)

type Value struct {
	Type     ValueType
	RefValue interface{}
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

type Content map[string]*Value

type ConfigObject struct {
	content *Content
}

type ConfigArray struct {
	content []*Value
}

type valueProvider interface {
	setString(name string, value string)
	setInt(name string, value string)
	setObject(name string, value *ConfigObject)
}

func splitPath(path string) []string {
	return strings.Split(path, ".")
}

func pathPrefix(path []string) ([]string, string) {
	length := len(path)
	if length == 1 {
		return []string{}, (path)[0]
	} else {
		return path[:len(path)-1], (path)[length-1]
	}
}

func (c *ConfigObject) setString(path string, value string) {
	prefix, key := pathPrefix(splitPath(path))
	resultKey := setObjectKey(prefix, c)
	(*resultKey.content)[key] = MakeStringValue(value)
}

func (c *ConfigObject) setInt(path string, value string) {
	prefix, key := pathPrefix(splitPath(path))
	resultKey := setObjectKey(prefix, c)
	(*resultKey.content)[key] = MakeNumericValue(value)
}

func (c *ConfigObject) setObject(path string, value *ConfigObject) {
	prefix, key := pathPrefix(splitPath(path))
	resultKey := setObjectKey(prefix, c)
	(*resultKey.content)[key] = &Value{
		Type:     ObjectType,
		RefValue: value,
	}
}

type hocon struct {
	*parser.BaseHOCONListener
	stack stack
}

func (r *hocon) ExitObject_data(ctx *parser.Object_dataContext) {
	current, _ := r.stack.Pop()
	parent, _ := r.stack.Peek()
	parent.setObject(ctx.Key().GetText(), current.(*ConfigObject))
}

func (r *hocon) EnterObject_data(ctx *parser.Object_dataContext) {
	r.stack.Push(NewConfigObject())
}

func (r *hocon) ExitString_data(ctx *parser.String_dataContext) {
	sd := ctx.String_value().GetText()
	if sd[0] == '"' || sd[0] == '\'' {
		sd = sd[1: len(sd)-1]
	}
	if v, err := r.stack.Peek(); err == nil {
		v.setString(ctx.Key().GetText(), sd)
	}
}

func (r *hocon) ExitNumber_data(ctx *parser.Number_dataContext) {
	if v, err := r.stack.Peek(); err == nil {
		v.setInt(ctx.Key().GetText(), ctx.NUMBER().GetText())
	}
}

func NewConfigObject() *ConfigObject {
	m := make(Content)
	co := ConfigObject{
		content: &m,
	}
	return &co
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

func newHocon() *hocon {
	h := new(hocon)
	h.stack = *NewStack()
	h.stack.Push(NewConfigObject())
	return h
}

func ParseHocon(stream antlr.CharStream) (o *ConfigObject, err error) {
	h := newHocon()
	ts := parser.NewHOCONLexer(stream)
	p := parser.NewHOCONParser(antlr.NewCommonTokenStream(ts, 0))
	p.AddParseListener(h)
	p.Hocon()
	return o, err
}

func ParseHoconString(data string) (o *ConfigObject, err error) {
	h := newHocon()
	ts := parser.NewHOCONLexer(antlr.NewInputStream(data))
	p := parser.NewHOCONParser(antlr.NewCommonTokenStream(ts, 0))
	p.AddParseListener(h)
	p.Hocon()
	return o, err
}

func ParseHoconFile(filename string) (o *ConfigObject, err error) {
	h := newHocon()
	if fs, err := antlr.NewFileStream(filename); err == nil {
		ts := parser.NewHOCONLexer(fs)
		p := parser.NewHOCONParser(antlr.NewCommonTokenStream(ts, 0))
		p.AddParseListener(h)
		p.Hocon()
		res, _ := h.stack.Pop()
		o = res.(*ConfigObject)
	}
	return o, err
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

func (o *ConfigObject) getString(path string) (res string) {
	if obj, key := traversePath(o, path); obj != nil {
		if v, ok := (*obj.content)[key]; ok {
			res = v.RefValue.(string)
		}
	}
	return res
}

func (o *ConfigObject) getInt(path string) (res int) {
	if obj, key := traversePath(o, path); obj != nil {
		if v, ok := (*obj.content)[key]; ok {
			res = v.RefValue.(int)
		}
	}
	return res
}

func (o *ConfigObject) getObject(path string) (res *ConfigObject) {
	if obj, key := traversePath(o, path); obj != nil {
		if v, ok := (*obj.content)[key]; ok {
			res = v.RefValue.(*ConfigObject)
		}
	}
	return res
}
