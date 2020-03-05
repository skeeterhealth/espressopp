/**
 * @begin 15-Feb-2020
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
	Identifier string   `  @Ident`
	Integer    *int     `| @Int`
	Decimal    *float64 `| @Float`
	String     *string  `| @String`
	Bool       *bool    `| @("true" | "false")`
}

type NumericTerm struct {
	Identifier string   `  @Ident`
	Integer    *int     `| @Int`
	Decimal    *float64 `| @Float`
}

type TextualTerm struct {
	Identifier string  `  @Ident`
	String     *string `| @String`
}

type Equality struct {
	Term1 *Term  `@@`
	Op    string `@("eq" | "neq")`
	Term2 *Term  `@@`
}

type Comparison struct {
	Term1 *NumericTerm `@@`
	Op    string       `@("gt" | "gte" | "lt" | "lte")`
	Term2 *NumericTerm `@@`
}

type NumericRange struct {
	Term1 *NumericTerm `@@ "between"`
	Term2 *NumericTerm `@@ "and"`
	Term3 *NumericTerm `@@`
}

type TextualMatching struct {
	Term1 *TextualTerm `@@`
	Op    string       `@("startswith" | "endswith" | "contains")`
	Term2 *TextualTerm `@@`
}

type Mathematics struct {
	Term1 NumericTerm `@@`
	Op    string      `@("plus" | "minus" | "mul" | "div")`
	Term2 NumericTerm `@@`
}

type Is struct {
	Ident string `@Ident?`
	Op    string `"is" @("not")?`
	Value string `@("true" | "false" | "null") | @Ident`
}

type Expression struct {
	Equality        *Equality        `  @@`
	Comparison      *Comparison      `| @@`
	NumericRange    *NumericRange    `| @@`
	TextualMatching *TextualMatching `| @@`
	Mathematics     *Mathematics     `| @@`
	Is              *Is              `| @@`
	ParenExpression *ParenExpression `| @@`
}

type ParenExpression struct {
	Op1        string      `@("and" | "or")`
	Op2        string      `@("not")?`
	Expression *Expression `"(" @@ ")"`
}

// Grammar is the set of structural rules that govern the composition of an
// Espesso++ expression.
type Grammar struct {
	Expressions []*Expression `@@+`
}

// Parser is the part of an interpreter that attaches meaning by classifying strings
// of tokens from the input Espresso++ expression as particular non-terminals
// and by building the parse tree.
type Parser struct {
}

// NewParser creates a new instance of Parser..
func NewParser() *Parser {
	return &Parser{}
}

// Parse parses the Espresso++ expressions in r and returns the resulting grammar.
func (p *Parser) Parse(r io.Reader) (error, *Grammar) {
	parser := participle.MustBuild(&Grammar{}, participle.UseLookahead(2))

	grammar := &Grammar{}
	err := parser.Parse(r, grammar)

	return err, grammar
}

// String returns a string representation of g.
func (p *Parser) String(g *Grammar) string {
	return repr.String(g, repr.Hide(&lexer.Position{}))
}
