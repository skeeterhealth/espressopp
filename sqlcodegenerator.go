/**
 * @begin 2020-02-15
 * @author <a href="mailto:giuseppe.greco@skeeterhealth.com">Giuseppe Greco</a>
 * @copyright 2020 <a href="skeeterhealth.com">Skeeter</a>
 */

package espressopp

import "io"

// SqlCodeGenerator is the CodeGenerator implementation that produces native SQL
// from Espresso++ expressions.
type SqlCodeGenerator struct {
	fieldNames map[string]string
}

// NewSqlCodeGenerator creates a new instance of SqlCodeGenerator.
func NewSqlCodeGenerator() CodeGenerator {
	return &SqlCodeGenerator{
		fieldNames: make(map[string]string),
	}
}

// Visit lets cg access the functionality provided by i to parse the Espresso++
// expressions in r and get back the grammar, which is then used to produce native
// SQL into w.
func (cg *SqlCodeGenerator) Visit(i Interpreter, r io.Reader, w io.Writer) error {
	return nil
}

// MapFieldNames lets the client application map each field in the input expression
// with a different name in the output SQL query. This Mapping is only necessary for
// those fields in the input expression that do not match the field names of the
// underlying database. If m is nil, then no mapping is applied.
func (cg *SqlCodeGenerator) MapFieldNames(m map[string]string) {
	if m == nil {
		if len(cg.fieldNames) > 0 {
			for fn := range cg.fieldNames {
				delete(cg.fieldNames, fn)
			}
		}
	} else {
		cg.fieldNames = m
	}
}
