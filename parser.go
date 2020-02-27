/**
 * @begin 15-Feb-2020
 * @author <a href="mailto:giuseppe.greco@skeeterhealth.com">Giuseppe Greco</a>
 * @copyright 2020 <a href="skeeterhealth.com">Skeeter</a>
 */

package espressopp

// Grammar is the set of structural rules that govern the composition of an
// Espesso++ expression.
type Grammar struct {
}

// Parser is the part of an interpreter that attaches meaning by classifying strings
// of tokens from the input Espresso++ expression as particular non-terminals
// and by building the parse tree.
type Parser struct {
}

// parse parses the Espresso++ expressions in r and returns the resulting grammar.
func (p *Parser) parse(r io.Reader) Grammar {
}
