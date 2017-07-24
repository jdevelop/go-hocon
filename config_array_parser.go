package hocon

import (
	"github.com/jdevelop/go-hocon/parser"
)

func (r *hocon) ExitArray_data(ctx *parser.Array_dataContext) {
	current, _ := r.stack.Pop()
	parent, _ := r.stack.Peek()
	parent.setArray(ctx.Key().GetText(), current.(*ConfigArray))
}

func (r *hocon) EnterArray_data(ctx *parser.Array_dataContext) {
	r.stack.Push(NewConfigArray())
}

func (r *hocon) EnterArray_array(ctx *parser.Array_arrayContext) {
	r.stack.Push(NewConfigArray())
}

func (r *hocon) ExitArray_array(ctx *parser.Array_arrayContext) {
	obj, _ := r.stack.Pop()
	p, _ := r.stack.Peek()
	p.setArray("", obj.(*ConfigArray))
}

func (r *hocon) ExitArray_string(ctx *parser.Array_stringContext) {
	res, _ := r.stack.Peek()
	res.setString("", stripStringQuotas(ctx.GetText()))
}

func (r *hocon) ExitArray_number(ctx *parser.Array_numberContext) {
	res, _ := r.stack.Peek()
	res.setInt("", ctx.GetText())
}

func (r *hocon) EnterArray_obj(ctx *parser.Array_objContext) {
	r.stack.Push(NewConfigObject())
}

func (r *hocon) ExitArray_obj(ctx *parser.Array_objContext) {
	obj, _ := r.stack.Pop()
	p, _ := r.stack.Peek()
	p.setObject("", obj.(*ConfigObject))
}
