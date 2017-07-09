package go_hocon

import (
	"testing"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"fmt"
	"strings"
)

func TestSimpleListener(t *testing.T) {
	is, _ := antlr.NewFileStream("test/reference.conf")
	res, _ := ParseHocon(is)
	var f func(int, *NamedValue)

	f = func(level int, n *NamedValue) {
		if n == nil {
			return
		}
		fmt.Print(strings.Repeat("+", level), "Node ", n.name)
		switch n.value.(type) {
		case []*NamedValue:
			fmt.Println()
			nvp, _ := n.value.([]*NamedValue)
			for _, v := range nvp {
				f(level+2, v)
			}
		default:
			fmt.Println("=", n.value)
		}
	}

	f(0, &res.NamedValue)
}
