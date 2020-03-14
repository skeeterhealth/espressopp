/**
 * @begin 2020-02-15
 * @author <a href="mailto:giuseppe.greco@skeeterhealth.com">Giuseppe Greco</a>
 * @copyright 2020 <a href="skeeterhealth.com">Skeeter</a>
 */

package espressopp

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"strings"
)

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
	err, grammar := i.Parse(r)
	if err != nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r)
		return errors.Wrapf(err, "error parsing %v", buf.String())
	}

	err, s := cg.emitGrammar(grammar)
	if err != nil {
		return errors.Wrapf(err, "error generating sql")
	}

	_, err = io.WriteString(w, s)
	return err
}

// MapFieldNames lets the client application map each field in the input expression
// with a different name in the output SQL query. This Mapping is only necessary for
// those fields in the input expression that do not match the field names of the
// underlying database. If m is nil, then no mapping is applied.
func (cg *SqlCodeGenerator) MapFieldNames(m map[string]string) {
	cg.fieldNames = m
}

// emitGrammar renders g.
func (cg *SqlCodeGenerator) emitGrammar(g *Grammar) (error, string) {
	var err error
	var sb strings.Builder

	for _, e := range g.Expressions {
		err, s := cg.emitExpression(e)
		if err != nil {
			return err, ""
		}
		sb.WriteString(s)
	}

	return err, sb.String()
}

// emitExpression renders e.
func (cg *SqlCodeGenerator) emitExpression(e *Expression) (error, string) {
	var err error
	var s string

	if e.Op != nil {
		s = fmt.Sprintf(" %s ", strings.ToUpper(*e.Op))
	} else if e.SubExpression != nil {
		err, s = cg.emitSubExpression(e.SubExpression)
	} else if e.Comparison != nil {
		err, s = cg.emitComparison(e.Comparison)
	} else if e.Equality != nil {
		err, s = cg.emitEquality(e.Equality)
	} else if e.Range != nil {
		err, s = cg.emitRange(e.Range)
	} else if e.Match != nil {
		err, s = cg.emitMatch(e.Match)
	} else if e.Is != nil {
		err, s = cg.emitIs(e.Is)
	}

	return err, s
}

// emitExpression renders se.
func (cg *SqlCodeGenerator) emitSubExpression(se *SubExpression) (error, string) {
	return nil, ""
}

// emitComparison renders c.
func (cg *SqlCodeGenerator) emitComparison(c *Comparison) (error, string) {
	return nil, ""
}

// emitEquality renders e.
func (cg *SqlCodeGenerator) emitEquality(e *Equality) (error, string) {
	return nil, ""
}

// emitRange renders r.
func (cg *SqlCodeGenerator) emitRange(r *Range) (error, string) {
	return nil, ""
}

// emitMatch renders m.
func (cg *SqlCodeGenerator) emitMatch(m *Match) (error, string) {
	return nil, ""
}

// emitIs renders i.
func (cg *SqlCodeGenerator) emitIs(i *Is) (error, string) {
	return nil, ""
}

// lookupFieldName gets the native field name associated with f. If no mapping is
// found, then f is returned.
func (cg *SqlCodeGenerator) lookupFieldName(f string) string {
	if val, ok := cg.fieldNames[f]; ok {
		return val
	}

	return f
}
