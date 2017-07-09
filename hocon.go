package go_hocon

import (
	"github.com/jdevelop/go-hocon/parser"
	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type NamedValue struct {
	name  string
	value interface{}
}

type Hocon struct {
	NamedValue
	parent          *Hocon
	currentChildIdx int
}

type hoconParser struct {
	*parser.BaseHOCONListener
	current *Hocon
}

func NewHocon() *Hocon {
	h := Hocon{
		currentChildIdx: 0,
		parent:          nil,
	}
	h.value = make([]*NamedValue, 10)
	return &h
}

func (h *Hocon) values() []*NamedValue {
	if res, err := h.value.([]*NamedValue); err {
		return res
	} else {
		return nil
	}
}

func (h *Hocon) ensureSpace() {
	if h.currentChildIdx > len(h.values())-2 {
		newVals := make([]*NamedValue, len(h.values())*2)
		copy(newVals, h.values())
		h.value = newVals
	}
}

func (h *Hocon) addNamedValueChild(name string) {
	h.ensureSpace()
	h.values()[h.currentChildIdx] = &NamedValue{
		name:  name,
		value: nil,
	}
}

func (h *Hocon) currentProperty() *NamedValue {
	return h.values()[h.currentChildIdx]
}

func (p *hoconParser) ExitKey(ctx *parser.KeyContext) {
	p.current.addNamedValueChild(ctx.NAME.GetText())
}

func (p *hoconParser) EnterObj(ctx *parser.ObjContext) {
	h := NewHocon()
	h.parent = p.current
	p.current = h
}

func (p *hoconParser) ExitObj(ctx *parser.ObjContext) {
	parent := p.current.parent
	p.current.name = parent.values()[parent.currentChildIdx].name
	parent.values()[parent.currentChildIdx] = &p.current.NamedValue
	parent.currentChildIdx++
	p.current = p.current.parent
}

func (p *hoconParser) ExitL_string(ctx *parser.L_stringContext) {
	p.current.currentProperty().value = ctx.STRING().GetText()
	p.current.currentChildIdx++
}

func (p *hoconParser) ExitL_rawstring(ctx *parser.L_rawstringContext) {
	p.current.currentProperty().value = ctx.RAWSTRING().GetText()
	p.current.currentChildIdx++
}

func (p *hoconParser) ExitL_number(ctx *parser.L_numberContext) {
	p.current.currentProperty().value = ctx.NUMBER().GetText()
	p.current.currentChildIdx++
}

func (p *hoconParser) ExitL_reference(ctx *parser.L_referenceContext) {
	p.current.currentProperty().value = ctx.REFERENCE().GetText()
	p.current.currentChildIdx++
}

func ParseHocon(is antlr.CharStream) (*Hocon, error) {
	hocon := NewHocon()
	hocon.name = "root"
	l := hoconParser{current: hocon}

	lex := parser.NewHOCONLexer(is)
	p := parser.NewHOCONParser(antlr.NewCommonTokenStream(lex, 0))
	p.BuildParseTrees = true
	p.AddParseListener(&l)

	p.Hocon()

	return hocon, nil
}
