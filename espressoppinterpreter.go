/**
 * @begin 2020-02-15
 * @author <a href="mailto:giuseppe.greco@skeeterhealth.com">Giuseppe Greco</a>
 * @copyright 2020 <a href="skeeterhealth.com">Skeeter</a>
 */

package espressopp

import (
	"io"

	"github.com/pkg/errors"
)

// EspressoppInterpreter is the Interpreter implementation that provides
// functionality for parsing Espresso++ expressions.
type EspressoppInterpreter struct {
	parser *parser
}

// NewEspressoppInterpreter creates a new instance of EspressoppInterpreter.
func NewEspressoppInterpreter() *EspressoppInterpreter {
	return &EspressoppInterpreter{
		parser: newParser(),
	}
}

// Accept lets cg access the functionality provided by i. More precisely, cg
// invokes i to parse the Espresso++ expressions in r and get back the grammar,
// which is then used to produce the native query into w.
func (i *EspressoppInterpreter) Accept(cg CodeGenerator, r io.Reader, w io.Writer) error {
	if cg == nil {
		return errors.New("code generator not specified")
	}

	return cg.Visit(i, r, w)
}

// Parse parses the expressions in r and returns the resulting grammar.
func (i *EspressoppInterpreter) Parse(r io.Reader) (*Grammar, error) {
	return i.parser.parse(r)
}
