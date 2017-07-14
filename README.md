# HOCON configuration for Golang :: [![Build Status](https://travis-ci.org/jdevelop/go-hocon.svg?branch=master)](https://travis-ci.org/jdevelop/go-hocon)

The basic implementation of [HOCON](https://github.com/typesafehub/config/blob/master/HOCON.md) for Golang, using [ANTLR](https://github.com/antlr/antlr4) grammar derived from [JSON](https://github.com/antlr/grammars-v4/blob/master/json/JSON.g4).

### Usage example

```Go
package main

import (
    "github.com/jdevelop/go-hocon"
    "fmt"
    )

func main() {
	res, _ := hocon.ParseHoconFile("reference.conf")
	
	res.getString("akka.persistence.view.auto-update")
	res.getString("akka.persistence.view.auto-update-replay-max")
	res.getInt("akka.persistence.view.auto-update-replay-min")

	obj := res.getObject("akka.persistence.snapshot-store.proxy")
	obj.getString("init-timeout")
}

```