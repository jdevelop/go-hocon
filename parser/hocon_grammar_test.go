package parser

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"testing"
)

type Err struct {
	antlr.DefaultErrorListener
	path string
	t    *testing.T
}

func (r *Err) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	r.t.Error(msg + r.path)
}

func TestSimpleGrammar(t *testing.T) {
	for _, path := range []string{"../test/simple1.conf", "../test/reference.conf"} {
		is, _ := antlr.NewFileStream(path)

		lex := NewHOCONLexer(is)
		p := NewHOCONParser(antlr.NewCommonTokenStream(lex, 0))
		p.BuildParseTrees = true

		e := Err{t: t, path: path}

		p.AddErrorListener(&e)
		p.Hocon()
	}

}
