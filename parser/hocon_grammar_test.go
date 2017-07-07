package parser

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"testing"
)

type Err struct {
	antlr.DefaultErrorListener
	t *testing.T
}

func (r *Err) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	r.t.Error(msg)
}

func TestSimpleGrammar(t *testing.T) {
	is, _ := antlr.NewFileStream("../test/simple1.conf")

	lex := NewHOCONLexer(is)
	p := NewHOCONParser(antlr.NewCommonTokenStream(lex, 0))
	p.BuildParseTrees = true

	e := Err{t: t}

	p.AddErrorListener(&e)
	p.Hocon()

}
