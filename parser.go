/**
 * @begin 2020-02-15
 * @author <a href="mailto:giuseppe.greco@skeeterhealth.com">Giuseppe Greco</a>
 * @copyright 2020 <a href="skeeterhealth.com">Skeeter</a>
 */

package espressopp

import (
	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/repr"
	"io"
)

type Term struct {
	Identifier *string  `  @Ident`
	Integer    *int     `| @Int`
	Decimal    *float64 `| @Float`
	String     *string  `| @String`
	Bool       *bool    `| @("true" | "false")`
	Macro      *Macro   `| @@`
}

type Macro struct {
	Name string  `@("#") @Ident`
	Args []*Term `("(" (@@ ("," @@)*)? ")")?`
}

type Math struct {
	Term1 *Term  `@@`
	Op    string `@("plus" | "minus" | "mul" | "div")`
	Term2 *Term  `@@`
}

type TermOrMath struct {
	Math      *Math `  @@`
	ParenMath *Math `| "(" @@ ")"`
	Term      *Term `| @@`
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

type ParenExpression struct {
	Not         bool          `@("not")?`
	Expressions []*Expression `"(" @@+ ")"`
}

type Expression struct {
	Op              *string          `  @("and" | "or")`
	ParenExpression *ParenExpression `| @@`
	Comparison      *Comparison      `| @@`
	Equality        *Equality        `| @@`
	Match           *Match           `| @@`
	Range           *Range           `| @@`
	Is              *Is              `| @@`
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
	p *participle.Parser
}

// newParser creates a new instance of parser.
func newParser() *parser {
	return &parser{
		p: participle.MustBuild(&Grammar{}, participle.UseLookahead(2)),
	}
}

// parse parses the Espresso++ expressions in r and returns the resulting grammar.
func (p *parser) parse(r io.Reader) (error, *Grammar) {
	grammar := &Grammar{}
	err := p.p.Parse(r, grammar)

	return err, grammar
}

// string returns a string representation of g.
func (p *parser) string(g *Grammar) string {
	return repr.String(g, repr.Hide(&lexer.Position{}))
}
