/**
 * @begin 15-Feb-2020
 * @author <a href="mailto:giuseppe.greco@skeeterhealth.com">Giuseppe Greco</a>
 * @copyright 2020 <a href="skeeterhealth.com">Skeeter</a>
 */

package espressopp

// Interpreter is the interface implemented by any interpreter that provides
// functionality for parsing a specific language and getting back the resulting
// grammar.
type Interpreter struct {
    // Accept lets the specified code generator access the functionality
    // provided by the interpreter. More precisely, the code generator invokes the
    // interpreter to parse the expressions in the specified reader and get back
    // the grammar, which is then used to produce the native query into the
    // specified writer.
	Accept (CodeGenerator, io.Reader, io.Writer)
	
	// Parse parses the expressions in the specified reader and returns the
	// resulting grammar.
	Parse(io.Reader) Grammar
}
