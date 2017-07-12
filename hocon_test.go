package go_hocon

import (
	"testing"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"fmt"
	"strings"
)

func TestReferenceListener(t *testing.T) {
	is, _ := antlr.NewFileStream("test/reference.conf")
	_, res := ParseHocon(is)
	dumpConfig(1, res.root)
}

func TestSimpleListener(t *testing.T) {
	is, _ := antlr.NewFileStream("test/simple1.conf")
	_, res := ParseHocon(is)
	dumpConfig(1, res.root)
}

func dumpConfig(level int, conf *ConfigObject) {
	prefix := strings.Repeat("+", level)
	for k, v := range *conf.content {
		switch v.Type {
		case StringType:
			fmt.Println(prefix, k, "=", v.RefValue.(string))
		case ObjectType:
			fmt.Println(prefix, k)
			dumpConfig(level+1, v.RefValue.(*ConfigObject))
		}
	}
}

func TestKeyTreeBuild(t *testing.T) {
	config := NewConfigObject(nil)

	setObjectKey("test1", config)
	setObjectKey("test2", config)
	setObjectKey("test1.passed1", config)
	p3 := setObjectKey("test2.passed2.passed3", config)
	setObjectKey("nested1", p3)
	setObjectKey("nested2", p3)
	setObjectKey("nested3", p3)

	dumpConfig(1, config)

}
