
/** Taken from "The Definitive ANTLR 4 Reference" by Terence Parr */

// Derived from http://json.org
grammar HOCON;


COMMENT
   : '#' ~( '\r' | '\n' )* -> skip
   ;

NUMBER
   : '-'? INT '.' [0-9] + EXP? | '-'? INT EXP | '-'? INT
   ;

STRING
   : '"' (ESC | ~ ["\\])* '"'
   | '\'' (ESC | ~ ['\\])* '\''
   ;

RAWSTRING
   : ( '.' | ALPHANUM | '-')+
   ;

REFERENCE
   : '${' (RAWSTRING)+ '}'
   ;

KV : [=:] -> skip
   ;

WS
   : [ \t\n\r] + -> skip
   ;

fragment ESC
   : '\\' (["\\/bfnrt] | UNICODE)
   ;


fragment UNICODE
   : 'u' HEX HEX HEX HEX
   ;

fragment ALPHANUM
   : ('0' .. '9') | ('a'..'z') | ('A' .. 'Z')
   ;

fragment HEX
   : [0-9a-fA-F]
   ;

fragment INT
   : '0' | [1-9] [0-9]*
   ;

// no leading zeros

fragment EXP
   : [Ee] [+\-]? INT
   ;

// \- since - means "range" inside [...]

//======================================================================================

hocon
   : value
   | property*
   ;

obj
   : object_begin property (','? property)* object_end
   | object_begin object_end
   ;

property
   : object_data
   | array_data
   | string_data
   | reference_data
   | number_data
   | rawstring_data
   ;

object_data
   : key KV? obj
   ;

array_data
   : key KV? array
   ;

string_data
   : key KV? STRING
   ;

reference_data
   : key KV? REFERENCE
   ;

number_data
   : key KV? NUMBER
   ;

rawstring_data
   : key KV? RAWSTRING
   ;

key
   : NAME=(STRING | RAWSTRING)
   ;

array
   : array_begin value (',' value)* array_end
   | array_begin array_end
   ;

value
   : STRING
   | REFERENCE
   | NUMBER
   | obj
   | array
   | RAWSTRING
   ;

array_begin
   : '['
   ;

array_end
   : ']'
   ;

object_begin
   : '{'
   ;

object_end
   : '}'
   ;