package go_hocon

import (
	"github.com/jdevelop/go-hocon/parser"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"strings"
)

type ValueType int

const (
	StringType ValueType = iota
	IntType
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

type Content map[string]*Value

type ConfigObject struct {
	parent  *ConfigObject
	content *Content
}

type updateKey func(string)

func noop(string) {}

type Hocon struct {
	*parser.BaseHOCONListener
	root    *ConfigObject
	updater updateKey
}

func (r *Hocon) ExitObject_data(ctx *parser.Object_dataContext) {
	r.root = r.root.parent
	r.updater = noop
}

func (r *Hocon) EnterObject_data(ctx *parser.Object_dataContext) {
	r.updater = func(path string) {
		r.root = setObjectKey(path, r.root)
		r.updater = noop
	}
}

func (r *Hocon) ExitKey(ctx *parser.KeyContext) {
	r.updater(ctx.NAME.GetText())
}

func (r *Hocon) ExitString_data(ctx *parser.String_dataContext) {
	(*r.root.content)[ctx.Key().GetText()] = MakeStringValue(ctx.STRING().GetText())
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

func ParseHocon(stream antlr.CharStream) (err error, h *Hocon) {
	h = new(Hocon)
	h.updater = noop
	h.root = NewConfigObject(nil)
	ts := parser.NewHOCONLexer(stream)
	p := parser.NewHOCONParser(antlr.NewCommonTokenStream(ts, 0))
	p.AddParseListener(h)
	p.Hocon()
	return err, h
}
