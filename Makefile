all: build

deps:
	go get -u github.com/antlr/antlr4/runtime/Go/antlr/...

gen:	deps
	java -jar ~/java/antlr-4.7-complete.jar -o parser -Dlanguage=Go parser/HOCON.g4

build:	gen
	go build ./parser/

test:	build
	go test ./parser/
    
