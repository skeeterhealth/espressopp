/**
 * @begin 2020-02-15
 * @author <a href="mailto:giuseppe.greco@skeeterhealth.com">Giuseppe Greco</a>
 * @copyright 2020 <a href="skeeterhealth.com">Skeeter</a>
 */

package espressopp

import (
	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/lexer/ebnf"
	"github.com/alecthomas/repr"
	"io"
)

type Term struct {
	Identifier *string  `  @Ident`
	Integer    *int     `| @Int`
	Decimal    *float64 `| @Float`
	String     *string  `| @String`
	Date       *string  `| @Date`
	Time       *string  `| @Time`
	DateTime   *string  `| @DateTime`
	Bool       *bool    `| @("true" | "false")`
	Macro      *Macro   `| @@`
}

type Macro struct {
	Name string  `@Macro`
	Args []*Term `("(" (@@ ("," @@)*)? ")")?`
}

type Math struct {
	Term1 *Term  `@@`
	Op    string `@("plus" | "minus" | "mul" | "div")`
	Term2 *Term  `@@`
}

type TermOrMath struct {
	Math    *Math `  @@`
	SubMath *Math `| "(" @@ ")"`
	Term    *Term `| @@`
}

type Equality struct {
	TermOrMath1 *TermOrMath `@@`
	Op          string      `@("eq" | "neq")`
	TermOrMath2 *TermOrMath `@@`
}

type Comparison struct {
	TermOrMath1 *TermOrMath `@@`
	Op          string      `@("gt" | "gte" | "lt" | "lte")`
	TermOrMath2 *TermOrMath `@@`
}

type Range struct {
	TermOrMath1 *TermOrMath `@@ "between"`
	TermOrMath2 *TermOrMath `@@ "and"`
	TermOrMath3 *TermOrMath `@@`
}

type Match struct {
	Term1 *Term  `@@`
	Op    string `@("startswith" | "endswith" | "contains")`
	Term2 *Term  `@@`
}

type Is struct {
	IsWithExplicitValue *IsWithExplicitValue `  @@`
	IsWithImplicitValue *IsWithImplicitValue `| @@`
}

type IsWithExplicitValue struct {
	Ident string `@Ident`
	Not   bool   `"is" @("not")?`
	Value string `@("true" | "false" | "null")`
}

type IsWithImplicitValue struct {
	Not   bool   `"is" @("not")?`
	Ident string `@Ident`
}

type SubExpression struct {
	Not         bool          `@("not")?`
	Expressions []*Expression `"(" @@+ ")"`
}

type Expression struct {
	Op            *string        `  @("and" | "or")`
	SubExpression *SubExpression `| @@`
	Comparison    *Comparison    `| @@`
	Equality      *Equality      `| @@`
	Range         *Range         `| @@`
	Match         *Match         `| @@`
	Is            *Is            `| @@`
}

// Grammar is the set of structural rules that govern the composition of an
// Espesso++ expression.
type Grammar struct {
	Expressions []*Expression `@@+`
}

// parser is the part of an interpreter that attaches meaning by classifying strings
// of tokens from the input Espresso++ expression as particular non-terminals
// and by building the parse tree.
type parser struct {
	espressoppParser *participle.Parser
}

var (
	espressoppLexer = lexer.Must(ebnf.New(`
		Comment = "//" { "\u0000"…"\uffff"-"\n" } .
		Date = "\"" date "\"" | "'" date "'" .
		Time = "\"" time "\"" | "'" time "'" .
		DateTime = "\"" date "T" time [ "+" digit digit ] "\"" | "'" date "T" time [ "+" digit digit  ] "'" .
		Ident = ident .
		Macro = "#" ident .
		String = "\"" { "\u0000"…"\uffff"-"\""-"\\" | "\\" any } "\"" | "'" { "\u0000"…"\uffff"-"'"-"\\" | "\\" any } "'" .
		Int = [ "-" | "+" ] digit { digit } .
		Float = ("." | digit) {"." | digit} .
		Punct = "!"…"/" | ":"…"@" | "["…` + "\"`\"" + ` | "{"…"~" .
		Whitespace = " " | "\t" | "\n" | "\r" .

		alpha = "a"…"z" | "A"…"Z" .
		digit = "0"…"9" .
		any = "\u0000"…"\uffff" .
		ident = (alpha | "_") { "_" | alpha | digit } .
		date = digit digit digit digit "-" digit digit "-" digit digit .
		time = digit digit ":" digit digit ":" digit digit [ "." { digit } ] .
	`))
)

// newParser creates a new instance of parser.
func newParser() *parser {
	return &parser{
		espressoppParser: participle.MustBuild(&Grammar{},
			participle.Lexer(espressoppLexer),
			participle.Unquote("String", "Date", "Time", "DateTime"),
			participle.Elide("Whitespace", "Comment"),
			participle.UseLookahead(2)),
	}
}

// parse parses the Espresso++ expressions in r and returns the resulting grammar.
func (p *parser) parse(r io.Reader) (error, *Grammar) {
	grammar := &Grammar{}
	err := p.espressoppParser.Parse(r, grammar)

	return err, grammar
}

// string returns a string representation of g.
func (p *parser) string(g *Grammar) string {
	return repr.String(g, repr.Hide(&lexer.Position{}))
}
