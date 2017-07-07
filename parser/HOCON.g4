
/** Taken from "The Definitive ANTLR 4 Reference" by Terence Parr */

// Derived from http://json.org
grammar HOCON;

hocon
   : value
   | pair
   ;

obj
   : object_begin pair (','? pair)* object_end
   | object_begin object_end
   ;

pair
   : STRING KV? value {fmt.Println("pairstr",$STRING.GetText())}
   | KEY KV? value {fmt.Println("pairkey",$KEY.GetText())}
   ;

array
   : array_begin value (',' value)* array_end
   | array_begin array_end
   ;

value
   : STRING {fmt.Println($STRING.GetText())}
   | REFERENCE {fmt.Println($REFERENCE.GetText())}
   | RAWSTRING {fmt.Println($RAWSTRING.GetText())}
   | NUMBER {fmt.Println($NUMBER.GetText())}
   | obj
   | array
   | 'true'
   | 'false'
   | 'null'
   ;

COMMENT
   : '#' ~( '\r' | '\n' )* -> skip
   ;

STRING
   : '"' (ESC | ~ ["\\])* '"'
   | '\'' (ESC | ~ ['\\])* '\''
   ;

RAWSTRING
   : (ESC | ALPHANUM)+
   ;

KEY
   : ( '.' | ALPHANUM | '-')+
   ;

REFERENCE
   : '${' (ALPHANUM|'.')+ '}'
   ;

fragment ESC
   : '\\' (["\\/bfnrt] | UNICODE)
   ;


fragment UNICODE
   : 'u' HEX HEX HEX HEX
   ;

fragment ALPHANUM
   : [a-zA-Z]
   ;

fragment HEX
   : [0-9a-fA-F]
   ;

KV
   : [=:]
   ;

array_begin
   : '[' { fmt.Println("BEGIN [") }
   ;

array_end
   : ']' { fmt.Println("] END") }
   ;

object_begin
   : '{' { fmt.Println("OBJ {") }
   ;

object_end
   : '}' { fmt.Println("} OBJ") }
   ;

NUMBER
   : '-'? INT '.' [0-9] + EXP? | '-'? INT EXP | '-'? INT
   ;

fragment INT
   : '0' | [1-9] [0-9]*
   ;

// no leading zeros

fragment EXP
   : [Ee] [+\-]? INT
   ;

// \- since - means "range" inside [...]

WS
   : [ \t\n\r] + -> skip
   ;
