/**
 * @begin 15-Feb-2020
 * @author <a href="mailto:giuseppe.greco@skeeterhealth.com">Giuseppe Greco</a>
 * @copyright 2020 <a href="skeeterhealth.com">Skeeter</a>
 */

package espressopp

import (
	"github.com/alecthomas/participle/lexer"
	"io"
)

type Term struct {
	Pos        lexer.Position
	Identifier string     `  @Ident`
	Integer    *int       `| @Int`
	Decimal    *float64   `| @Float`
	String     *string    `| @String`
	Bool       *bool      `| @("true" | "false")`
}

type NumericTerm struct {
	Pos        lexer.Position
	Identifier string   `  @Ident`
	Integer    *int     `| @Int`
	Decimal    *float64 `| @Float`
}

type TextualTerm struct {
	Pos        lexer.Position
	Identifier string  `  @Ident`
	String     *string `| @String`
}

type Equality struct {
	Pos   lexer.Position
	Term1 *Term  `@@`
	Op    string `@("eq" | "neq")`
	Term2 *Term  `@@`
}

type Comparison struct {
	Pos   lexer.Position
	Term1 *NumericTerm `@@`
	Op    string       `@("gt" | "gte" | "lt" | "lte")`
	Term2 *NumericTerm `@@`
}

type NumericRange struct {
	Pos   lexer.Position
	Term1 *NumericTerm `@@ "between"`
	Term2 *NumericTerm `@@ "and"`
	Term3 *NumericTerm `@@`
}

type TextualMatching struct {
	Pos   lexer.Position
	Term1 *TextualTerm `@@`
	Op    string       `@("startswith" | "endswith" | "contains")`
	Term2 *TextualTerm `@@`
}

type MathOperation struct {
    Pos   lexer.Position
    Term1 NumericTerm `@@`
    Op    string      `@("plus" | "minus" | "mul" | "div")`
    Term2 NumericTerm `@@`
}

type Is struct {
	Pos   lexer.Position
	Ident string `@Ident?`
	Op    string `"is" @("not")?`
	Value string `@("true" | "false" | "null") | @Ident`
}

type Expression struct {
	Pos             lexer.Position
	Equality        *Equality        `  @@`
	Comparison      *Comparison      `| @@`
	NumericRange    *NumericRange    `| @@`
	TextualMatching *TextualMatching `| @@`
	MathOperation   *MathOperation   `| @@`
	Is              *Is              `| @@`
	ParenExpression *ParenExpression `| @@`
}

type ParenExpression struct {
	Pos        lexer.Position
	Op1        string      `@("and" | "or")`
	Op2        string      `@("not")?`
	Expression *Expression `(" @@ ")"`
}

// Grammar is the set of structural rules that govern the composition of an
// Espesso++ expression.
type Grammar struct {
	Pos         lexer.Position
	Expressions []*Expression `@@+`
}

// Parser is the part of an interpreter that attaches meaning by classifying strings
// of tokens from the input Espresso++ expression as particular non-terminals
// and by building the parse tree.
type Parser struct {
}

// parse parses the Espresso++ expressions in r and returns the resulting grammar.
func (p *Parser) parse(r io.Reader) (error, *Grammar) {
	return nil, nil
}
