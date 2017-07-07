all: test

antlr/antlr-4.7-complete.jar : 
	mkdir antlr
	wget "http://www.antlr.org/download/antlr-4.7-complete.jar" -O antlr/antlr-4.7-complete.jar

deps: antlr/antlr-4.7-complete.jar
	go get -u github.com/antlr/antlr4/runtime/Go/antlr/...

gen:	deps
	java -jar antlr/antlr-4.7-complete.jar -o . -Dlanguage=Go parser/HOCON.g4

build:	gen
	go build ./...

test:	build
	go test ./...
    
