package hocon

import (
	"github.com/jdevelop/go-hocon/parser"
)

func (hc *hocon) ExitObject_data(ctx *parser.Object_dataContext) {
	current, _ := hc.stack.Pop()
	parent, _ := hc.stack.Peek()
	parent.setObject(ctx.Key().GetText(), current.(*ConfigObject))
}

func (hc *hocon) EnterObject_data(ctx *parser.Object_dataContext) {
	hc.stack.Push(NewConfigObject())
}

func (hc *hocon) ExitV_string(context *parser.V_stringContext) {
	if hc.compoundRef != nil {
		hc.compoundRef.Value = append(hc.compoundRef.Value,
			MakeStringValue(stripStringQuotas(context.STRING().GetText())))
	}
}

func (hc *hocon) ExitV_rawstring(context *parser.V_rawstringContext) {
	if hc.compoundRef != nil {
		hc.compoundRef.Value = append(hc.compoundRef.Value, MakeStringValue(
			stripStringQuotas(context.Rawstring().GetText())))
	}
}

func (hc *hocon) ExitV_reference(context *parser.V_referenceContext) {
	if hc.compoundRef != nil {
		hc.compoundRef.Value = append(hc.compoundRef.Value,
			MakeReferenceValue(stripStringQuotas(context.REFERENCE().GetText())))
	}
}

func (hc *hocon) EnterString_data(ctx *parser.String_dataContext) {
	hc.compoundRef = new(CompoundString)
	hc.compoundRef.Value = make([]*Value, 0)
}

func (hc *hocon) ExitString_data(ctx *parser.String_dataContext) {
	if hc.compoundRef != nil {
		if v, err := hc.stack.Peek(); err == nil {
			v.setCompoundString(ctx.Key().GetText(), hc.compoundRef)
		}
		hc.compoundRef = nil // cleanup
	} else {
		if v, err := hc.stack.Peek(); err == nil {
			v.setString(ctx.Key().GetText(), stripStringQuotas(ctx.String_value().GetText()))
		}
	}
}

func (hc *hocon) ExitNumber_data(ctx *parser.Number_dataContext) {
	if v, err := hc.stack.Peek(); err == nil {
		v.setInt(ctx.Key().GetText(), ctx.NUMBER().GetText())
	}
}

func (hc *hocon) ExitReference_data(ctx *parser.Reference_dataContext) {
	if v, err := hc.stack.Peek(); err == nil {
		v.setInt(ctx.Key().GetText(), ctx.REFERENCE().GetText())
	}
}
