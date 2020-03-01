/**
 * @begin 15-Feb-2020
 * @author <a href="mailto:giuseppe.greco@skeeterhealth.com">Giuseppe Greco</a>
 * @copyright 2020 <a href="skeeterhealth.com">Skeeter</a>
 */

package espressopp

import "io"

// CodeGenerator is the interface implemented by any code generator that produces
// native queries from expressions written in a language supported by a given
// interpreter.
type CodeGenerator interface {
	// Visit lets the code generator access the functionality provided by the
	// specified interpreter. More precisely, the code generator invokes the
	// interpreter to parse the expressions in the specified reader and get
	// back the grammar, which is then used to produce the native query into
	// the specified writer.
	Visit(Interpreter, io.Reader, io.Writer) error

	// MapFieldNames lets the client application map each field in the input
	// expression with a different name in the output native query. This
	// mapping is only necessary for those fields in the input expression
	// that do not match the field names of the underlying data management
	// system.
	MapFieldNames(map[string]string)
}
