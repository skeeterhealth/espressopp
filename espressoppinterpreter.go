/**
 * @begin 15-Feb-2020
 * @author <a href="mailto:giuseppe.greco@skeeterhealth.com">Giuseppe Greco</a>
 * @copyright 2020 <a href="skeeterhealth.com">Skeeter</a>
 */

package espressopp

import "io"

// EspressoppInterpreter is the Interpreter implementation that provides
// functionality for parsing Espresso++ expressions.
type EspressoppInterpreter struct {
	parser *parser
}

// NewEspressoppInterpreter creates a new instance of EspressoppInterpreter.
func NewEspressoppInterpreter() Interpreter {
	return &EspressoppInterpreter{
		parser: &parser{},
	}
}

// Accept lets cg access the functionality provided by i. More precisely, cg
// invokes i to parse the Espresso++ expressions in r and get back the grammar,
// which is then used to produce the native query into w.
func (i *EspressoppInterpreter) Accept(cg CodeGenerator, r io.Reader, w io.Writer) error {
	return cg.Visit(i, r, w)
}

// Parse parses the expressions in r and returns the resulting grammar.
func (i *EspressoppInterpreter) Parse(r io.Reader) (error, *Grammar) {
	return nil, nil
}
