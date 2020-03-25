/**
 * @begin 2020-02-15
 * @author <a href="mailto:giuseppe.greco@skeeterhealth.com">Giuseppe Greco</a>
 * @copyright 2020 <a href="skeeterhealth.com">Skeeter</a>
 */

package espressopp

import "io"

// Interpreter is the interface implemented by any interpreter that provides
// functionality for parsing a specific language and getting back the resulting
// grammar.
type Interpreter interface {
	// Accept lets the specified code generator access the functionality
	// provided by the interpreter. More precisely, the code generator invokes the
	// interpreter to parse the expressions in the specified reader and get back
	// the grammar, which is then used to produce the native query into the
	// specified writer.
	Accept(CodeGenerator, io.Reader, io.Writer) error

	// Parse parses the expressions in the specified reader and returns the
	// resulting grammar.
	Parse(io.Reader) (*Grammar, error)
}
