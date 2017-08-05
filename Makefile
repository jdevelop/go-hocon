ANTLRJAR=antlr/antlr-4.7-complete.jar


all: test

$(ANTLRJAR) : 
	mkdir antlr
	wget "http://www.antlr.org/download/antlr-4.7-complete.jar" -O antlr/antlr-4.7-complete.jar

deps: $(ANTLRJAR)
	go get -u github.com/antlr/antlr4/runtime/Go/antlr/...
	go get -u github.com/stretchr/testify/...

gen:	parser/HOCON.g4
	java -jar $(ANTLRJAR) -o . -listener -Dlanguage=Go parser/HOCON.g4

gen-java:	java/parser/*.java
	java -jar $(ANTLRJAR) -o java/ -listener -Dlanguage=Java parser/HOCON.g4

compile-java:	gen-java java/parser/*.java
	javac -cp $(ANTLRJAR) java/parser/*.java

grun:	compile-java
	java -cp $(ANTLRJAR):java/parser/ org.antlr.v4.gui.TestRig HOCON hocon ${GARGS}

build:	gen
	go build ./...

test:	build
	go test ./...
