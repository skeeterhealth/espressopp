/**
 * @begin 15-Feb-2020
 * @author <a href="mailto:giuseppe.greco@skeeterhealth.com">Giuseppe Greco</a>
 * @copyright 2020 <a href="skeeterhealth.com">Skeeter</a>
 */

package espressopp

type Espressopp struct {
}

func (i *Espressopp) Accept(cg CodeGenerator, r io.Reader, w io.Writer) {
	cg.Visit(i, r, w)
}

func (i *Espressopp) Parse(r io.Reader) Grammar {
	return nil
}
