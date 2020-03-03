/**
 * @begin 15-Feb-2020
 * @author <a href="mailto:giuseppe.greco@skeeterhealth.com">Giuseppe Greco</a>
 * @copyright 2020 <a href="skeeterhealth.com">Skeeter</a>
 */

package espressopp

import (
	//	"github.com/alecthomas/participle"
	"io"
)

type Term struct {
	Identifier string   `  @Ident`
	Number     *float64 `| @Float | @Int`
	String     *string  `| @String`
	Bool       *bool    `| @("true" | "false")`
}

type NumericTerm struct {
	Identifier string   `  @Ident`
	Number     *float64 `| @Float | @Int`
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
	Term1   *NumericTerm `@@`
	Between string       `"between"`
	Term2   *NumericTerm `@@`
	And     string       `"and"`
	Term3   *NumericTerm `@@`
}

type TextualMatching struct {
	Term1 *TextualTerm `@@`
	Op    string       `@("startswith" | "endswith" | "contains")`
	Term2 *TextualTerm `@@`
}

type Is struct {
	IdentIs string `  @Ident "is" @("not")? @("true" | "false" | "null")`
	IsIdent string `| "is" @("not")? @Ident`
}

type Expression struct {
	Equality        *Equality        `  @@`
	Comparison      *Comparison      `| @@`
	NumericRange    *NumericRange    `| @@`
	TextualMatching *TextualMatching `| @@`
	Is              *Is              `| @@`
	ParenExpression *ParenExpression `| @("and" | "or")? @@`
}

type ParenExpression struct {
	Expression *Expression `@("not")? "(" @@ ")"`
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

// parse parses the Espresso++ expressions in r and returns the resulting grammar.
func (p *Parser) parse(r io.Reader) (error, *Grammar) {
	return nil, nil
}
