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
		path := stripStringQuotas(context.REFERENCE().GetText())
		refPath := referencePath(path)
		var value *Value
		if v, ok := hc.refs[refPath]; ok {
			value = v
		} else {
			value = MakeReferenceValue(path)
			hc.refs[refPath] = value
		}
		hc.compoundRef.Value = append(hc.compoundRef.Value, value)
	}
}

func (hc *hocon) EnterString_data(ctx *parser.String_dataContext) {
	hc.compoundRef = new(CompoundString)
	hc.compoundRef.Value = make([]*Value, 0)
}

func (hc *hocon) ExitString_data(ctx *parser.String_dataContext) {
	if hc.compoundRef != nil {
		refExists := false
		for _, v := range hc.compoundRef.Value {
			if v.Type == ReferenceType {
				refExists = true
			}
		}
		if refExists {
			if v, err := hc.stack.Peek(); err == nil {
				v.setCompoundString(ctx.Key().GetText(), hc.compoundRef)
			}
		} else {
			str := ""
			for _, v := range hc.compoundRef.Value {
				str = v.RefValue.(string) + str
			}
			if v, err := hc.stack.Peek(); err == nil {
				v.setString(ctx.Key().GetText(), str)
			}
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
