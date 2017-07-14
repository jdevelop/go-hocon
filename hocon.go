package hocon

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/jdevelop/go-hocon/parser"
	"strconv"
	"strings"
)

type ValueType int

const (
	StringType ValueType = iota
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
	parent  *ConfigObject
	content *Content
}

type updateKey func(string)

func noop(string) {}

type hocon struct {
	*parser.BaseHOCONListener
	root    *ConfigObject
	updater updateKey
}

func (r *hocon) ExitObject_data(ctx *parser.Object_dataContext) {
	r.root = r.root.parent
	r.updater = noop
}

func (r *hocon) EnterObject_data(ctx *parser.Object_dataContext) {
	r.updater = func(path string) {
		r.root = setObjectKey(path, r.root)
		r.updater = noop
	}
}

func (r *hocon) ExitKey(ctx *parser.KeyContext) {
	r.updater(ctx.GetText())
}

func (r *hocon) ExitString_data(ctx *parser.String_dataContext) {
	sd := ctx.String_value().GetText()
	if sd[0] == '"' || sd[0] == '\'' {
		sd = sd[1 : len(sd)-1]
	}
	(*r.root.content)[ctx.Key().GetText()] = MakeStringValue(sd)
}

func (r *hocon) ExitNumber_data(ctx *parser.Number_dataContext) {
	(*r.root.content)[ctx.Key().GetText()] = MakeNumericValue(ctx.NUMBER().GetText())
}

func NewConfigObject(parent *ConfigObject) *ConfigObject {
	m := make(Content)
	co := ConfigObject{
		parent:  parent,
		content: &m,
	}
	return &co
}

func setObjectKey(path string, obj *ConfigObject) *ConfigObject {
	keys := strings.Split(path, ".")
	oldRoot := obj
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

		newObject := NewConfigObject(obj)
		(*obj.content)[key] = &Value{
			Type:     ObjectType,
			RefValue: newObject,
		}
		obj = newObject
	}
	obj.parent = oldRoot
	return obj
}

func ParseHocon(stream antlr.CharStream) (o *ConfigObject, err error) {
	h := new(hocon)
	h.updater = noop
	h.root = NewConfigObject(nil)
	ts := parser.NewHOCONLexer(stream)
	p := parser.NewHOCONParser(antlr.NewCommonTokenStream(ts, 0))
	p.AddParseListener(h)
	p.Hocon()
	o = h.root
	return o, err
}

func ParseHoconString(data string) (o *ConfigObject, err error) {
	h := new(hocon)
	h.updater = noop
	h.root = NewConfigObject(nil)
	ts := parser.NewHOCONLexer(antlr.NewInputStream(data))
	p := parser.NewHOCONParser(antlr.NewCommonTokenStream(ts, 0))
	p.AddParseListener(h)
	p.Hocon()
	o = h.root
	return o, err
}

func ParseHoconFile(filename string) (o *ConfigObject, err error) {
	h := new(hocon)
	h.updater = noop
	h.root = NewConfigObject(nil)
	if fs, err := antlr.NewFileStream(filename); err == nil {
		ts := parser.NewHOCONLexer(fs)
		p := parser.NewHOCONParser(antlr.NewCommonTokenStream(ts, 0))
		p.AddParseListener(h)
		p.Hocon()
		o = h.root
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
