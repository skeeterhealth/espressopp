/**
 * @begin 15-Feb-2020
 * @author <a href="mailto:giuseppe.greco@skeeterhealth.com">Giuseppe Greco</a>
 * @copyright 2020 <a href="skeeterhealth.com">Skeeter</a>
 */

package espressopp

import "io"

// SqlCodeGenerator is the CodeGenerator implementation that produces native SQL
// from Espresso++ expressions.
type SqlCodeGenerator struct {
}

// Visit lets cg access the functionality provided by i to parse the Espresso++
// expressions in r and get back the grammar, which is then used to produce native
// SQL into w.
func (cg *SqlCodeGenerator) Visit(i Interpreter, r io.Reader, w io.Writer) error {
	return nil
}
