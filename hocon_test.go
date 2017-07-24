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
	assert.Equal(t, "on", res.getString("akka.persistence.view.auto-update"))
	assert.Equal(t, "off", res.getString("akka.persistence.view.auto-update-replay-max"))
	assert.Equal(t, -1, res.getInt("akka.persistence.view.auto-update-replay-min"))
	obj := res.getObject("akka.persistence.snapshot-store.proxy")
	assert.Equal(t, "10s", obj.getString("init-timeout"))
}

func TestSimpleArrayListener(t *testing.T) {
	res, _ := ParseHoconFile("test/simple2.conf")
	dumpConfig(1, res)
	arr := res.getArray("akka.persistence.view.arrays.array")
	assert.Equal(t, "1", arr.getString(0))
	assert.Equal(t, 100500, arr.getInt(4))
	assert.Equal(t, 1, arr.getArray(5).getInt(0))
	assert.Equal(t, 2, arr.getArray(5).getInt(1))
	assert.Equal(t, 3, arr.getObject(3).getArray("test.passed").getInt(0))
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

	assert.NotNil(t, config.getObject("test2.passed2").getObject("passed3"))
	assert.Nil(t, config.getObject("test2.passed2").getObject("passed5"))

}