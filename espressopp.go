/**
 * @begin 15-Feb-2020
 * @author <a href="mailto:giuseppe.greco@skeeterhealth.com">Giuseppe Greco</a>
 * @copyright 2020 <a href="skeeterhealth.com">Skeeter</a>
 */

package espressopp

// Espressopp is the Interpreter implementation that provides functionality
// for parsing Espresso++ expressions.
type Espressopp struct {
}

// Accept lets cg access the functionality provided by i. More precisely, cg
// invokes i to parse the Espresso++ expressions in r and get back the grammar,
// which is then used to produce the native query into w.
func (i *Espressopp) Accept(cg CodeGenerator, r io.Reader, w io.Writer) error {
	return cg.Visit(i, r, w)
}

// Parse parses the expressions in r and returns the resulting grammar.
func (i *Espressopp) Parse(r io.Reader) (error, Grammar) {
	return nil, nil
}
