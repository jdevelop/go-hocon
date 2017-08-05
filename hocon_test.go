package hocon

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestReferenceListener(t *testing.T) {
	res, _ := ParseHoconFile("test/reference.conf")
	dumpConfig(1, res)
}

func TestSimpleListener(t *testing.T) {
	res, _ := ParseHoconFile("test/simple1.conf")
	dumpConfig(1, res)
	assert.Equal(t, "on", res.GetString("akka.persistence.view.auto-update"))
	assert.Equal(t, "off", res.GetString("akka.persistence.view.auto-update-replay-max"))
	assert.Equal(t, -1, res.GetInt("akka.persistence.view.auto-update-replay-min"))
	obj := res.GetObject("akka.persistence.snapshot-store.proxy")
	assert.Equal(t, "10s", obj.GetString("init-timeout"))
}

func TestSimpleArrayListener(t *testing.T) {
	res, _ := ParseHoconFile("test/simple2.conf")
	dumpConfig(1, res)
	arr := res.GetArray("akka.persistence.view.arrays.array")
	assert.Equal(t, "11000", arr.GetString(0))
	assert.Equal(t, 100500, arr.GetInt(5))
	assert.Equal(t, 1, arr.GetArray(6).GetInt(0))
	assert.Equal(t, 2, arr.GetArray(6).GetInt(1))
	assert.Equal(t, 3, arr.GetObject(4).GetArray("test.passed").GetInt(0))
}

func TestReferencesListener(t *testing.T) {
	res, _ := ParseHoconFile("test/references.conf")
	assert.Equal(t, "11005002", res.GetString("test_string"))
	assert.Equal(t, "100500", res.GetString("test"))
	assert.Equal(t, "hello world", res.GetString("another.sentence"))
	dumpConfig(1, res)
}

func TestMerge(t *testing.T) {
	res1, _ := ParseHoconFile("test/simple1.conf")
	res2, _ := ParseHoconFile("test/simple2.conf")
	res1.Merge(res2)
	dumpConfig(1, res1)
}

func dumpConfig(level int, conf *ConfigObject) {
	prefix := strings.Repeat("-", level)
	for k, v := range *conf.content {
		switch v.Type {
		case NumericType:
			fmt.Println(prefix, k, "=", v.RefValue.(int))
		case StringType:
			fmt.Println(prefix, k, "=", v.RefValue.(string))
		case ArrayType:
			fmt.Println(prefix, k, "= [")
			dumpArray(level, v.RefValue.(*ConfigArray))
			fmt.Println(prefix, "]")
		case ObjectType:
			fmt.Println(prefix, k, "{")
			dumpConfig(level+1, v.RefValue.(*ConfigObject))
			fmt.Println(prefix, "}")
		case CompoundStringType:
			fmt.Print(prefix, k, "=")
			for _, v := range v.RefValue.(*CompoundString).Value {
				fmt.Print(v.RefValue, ",")
			}
			fmt.Println()
			fmt.Println(prefix, "}")
		}
	}
}

func dumpArray(level int, arr *ConfigArray) {
	prefix := strings.Repeat("-", level)
	for i := 0; i < arr.idx; i++ {
		v := arr.content[i]
		switch v.Type {
		case NumericType:
			fmt.Println(prefix, v.RefValue.(int))
		case StringType:
			fmt.Println(prefix, v.RefValue.(string))
		case ArrayType:
			fmt.Println(prefix, "[")
			dumpArray(level+1, v.RefValue.(*ConfigArray))
			fmt.Println(prefix, "]")
		case ObjectType:
			fmt.Println(prefix, "{")
			dumpConfig(level+1, v.RefValue.(*ConfigObject))
			fmt.Println(prefix, "}")
		}
	}
}

func path(path string) []string {
	return strings.Split(path, ".")
}

func TestKeyTreeBuild(t *testing.T) {
	config := NewConfigObject()

	setObjectKey(path("test1"), config)
	setObjectKey(path("test2"), config)
	setObjectKey(path("test1.passed1"), config)
	p3 := setObjectKey(path("test2.passed2.passed3"), config)
	setObjectKey(path("nested1"), p3)
	setObjectKey(path("nested2"), p3)
	setObjectKey(path("nested3"), p3)

	dumpConfig(1, config)

	assert.NotNil(t, config.GetObject("test2.passed2").GetObject("passed3"))
	assert.Nil(t, config.GetObject("test2.passed2").GetObject("passed5"))

}
