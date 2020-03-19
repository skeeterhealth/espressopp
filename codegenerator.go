/**
 * @begin 2020-02-15
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

	// GetRenderingOptions gets the rendering options used by the CodeGenerator
	// implementation.
	GetRenderingOptions() *RenderingOptions
}
