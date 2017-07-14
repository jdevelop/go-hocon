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
	assert.Equal(t, "on", res.getString("akka.persistence.view.auto-update"))
	assert.Equal(t, "off", res.getString("akka.persistence.view.auto-update-replay-max"))
	assert.Equal(t, -1, res.getInt("akka.persistence.view.auto-update-replay-min"))
	obj := res.getObject("akka.persistence.snapshot-store.proxy")
	assert.Equal(t, "10s", obj.getString("init-timeout"))
	dumpConfig(1, res)
}

func dumpConfig(level int, conf *ConfigObject) {
	prefix := strings.Repeat("-", level)
	for k, v := range *conf.content {
		switch v.Type {
		case NumericType:
			fmt.Println(prefix, k, "=", v.RefValue.(int))
		case StringType:
			fmt.Println(prefix, k, "=", v.RefValue.(string))
		case ObjectType:
			fmt.Println(prefix, k, "{")
			dumpConfig(level+1, v.RefValue.(*ConfigObject))
			fmt.Println(prefix, "}")
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

	assert.NotNil(t, config.getObject("test2.passed2").getObject("passed3"))
	assert.Nil(t, config.getObject("test2.passed2").getObject("passed5"))

}
