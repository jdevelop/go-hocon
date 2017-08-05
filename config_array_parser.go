package hocon

import (
	"github.com/jdevelop/go-hocon/parser"
)

func (hc *hocon) ExitArray_data(ctx *parser.Array_dataContext) {
	current, _ := hc.stack.Pop()
	parent, _ := hc.stack.Peek()
	parent.setArray(ctx.Key().GetText(), current.(*ConfigArray))
}

func (hc *hocon) EnterArray_data(ctx *parser.Array_dataContext) {
	hc.stack.Push(NewConfigArray(hc))
}

func (hc *hocon) EnterArray_array(ctx *parser.Array_arrayContext) {
	hc.stack.Push(NewConfigArray(hc))
}

func (hc *hocon) ExitArray_array(ctx *parser.Array_arrayContext) {
	obj, _ := hc.stack.Pop()
	p, _ := hc.stack.Peek()
	p.setArray("", obj.(*ConfigArray))
}

func (hc *hocon) EnterArray_string(ctx *parser.Array_stringContext) {
	hc.compoundRef = new(CompoundString)
	hc.compoundRef.Value = make([]*Value, 0)
}

func (hc *hocon) ExitArray_string(ctx *parser.Array_stringContext) {
	if hc.compoundRef != nil {
		if v, err := hc.stack.Peek(); err == nil {
			v.setCompoundString("", hc.compoundRef)
		}
		hc.compoundRef = nil // cleanup
	} else {
		if v, err := hc.stack.Peek(); err == nil {
			v.setString("", stripStringQuotas(ctx.String_value().GetText()))
		}
	}
}

func (hc *hocon) ExitArray_number(ctx *parser.Array_numberContext) {
	res, _ := hc.stack.Peek()
	res.setInt("", ctx.GetText())
}

func (hc *hocon) EnterArray_obj(ctx *parser.Array_objContext) {
	hc.stack.Push(NewConfigObject())
}

func (hc *hocon) ExitArray_obj(ctx *parser.Array_objContext) {
	obj, _ := hc.stack.Pop()
	p, _ := hc.stack.Peek()
	p.setObject("", obj.(*ConfigObject))
}
