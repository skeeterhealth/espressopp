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
	"strconv"
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
	var err error
	var sb strings.Builder

	if se.Not {
		sb.WriteString("NOT ")
	}

	sb.WriteString("(")

	for _, e := range se.Expressions {
		err, s := cg.emitExpression(e)
		if err != nil {
			return err, ""
		}
		sb.WriteString(s)
	}

	sb.WriteString(")")

	return err, sb.String()
}

// emitComparison renders c.
func (cg *SqlCodeGenerator) emitComparison(c *Comparison) (error, string) {
	err, t1 := cg.emitTermOrMath(c.TermOrMath1)
	if err != nil {
		return err, ""
	}

	err, t2 := cg.emitTermOrMath(c.TermOrMath2)
	if err != nil {
		return err, ""
	}

	var op string
	switch c.Op {
	case "gt":
		op = ">"
	case "gte":
		op = ">="
	case "lt":
		op = "<"
	case "lte":
		op = "<="
	}

	return err, fmt.Sprintf("%s %s %s", t1, op, t2)
}

// emitEquality renders e.
func (cg *SqlCodeGenerator) emitEquality(e *Equality) (error, string) {
	err, t1 := cg.emitTermOrMath(e.TermOrMath1)
	if err != nil {
		return err, ""
	}

	err, t2 := cg.emitTermOrMath(e.TermOrMath2)
	if err != nil {
		return err, ""
	}

	var op string
	switch e.Op {
	case "eq":
		op = "="
	case "neq":
		op = "<>"
	}

	return err, fmt.Sprintf("%s %s %s", t1, op, t2)
}

// emitRange renders r.
func (cg *SqlCodeGenerator) emitRange(r *Range) (error, string) {
	err, t1 := cg.emitTermOrMath(r.TermOrMath1)
	if err != nil {
		return err, ""
	}

	err, t2 := cg.emitTermOrMath(r.TermOrMath2)
	if err != nil {
		return err, ""
	}

	err, t3 := cg.emitTermOrMath(r.TermOrMath3)
	if err != nil {
		return err, ""
	}

	return err, fmt.Sprintf("%s %s %s %s %s", t1,
		strings.ToUpper(r.Between), t2,
		strings.ToUpper(r.And), t3)
}

// emitMatch renders m.
func (cg *SqlCodeGenerator) emitMatch(m *Match) (error, string) {
	err, t1 := cg.emitTerm(m.Term1)
	if err != nil {
		return err, ""
	}

	err, t2 := cg.emitTerm(m.Term2)
	if err != nil {
		return err, ""
	}

	t2 = strings.ReplaceAll(strings.ReplaceAll(t2, "'", ""), "\"", "")

	switch m.Op {
	case "startswith":
		t2 = fmt.Sprintf("'%s%%'", t2)
	case "endswith":
		t2 = fmt.Sprintf("'%%%s'", t2)
	case "contains":
		t2 = fmt.Sprintf("'%%%s%%'", t2)
	}

	return err, fmt.Sprintf("%s %s %s", t1, "LIKE", t2)
}

// emitIs renders i.
func (cg *SqlCodeGenerator) emitIs(i *Is) (error, string) {
	var sb strings.Builder

	if i.IsWithExplicitValue != nil {
		sb.WriteString(i.IsWithExplicitValue.Ident)
		if i.IsWithExplicitValue.Value == "null" {
			var not string
			if i.IsWithExplicitValue.Not {
				not = "NOT "
			}
			sb.WriteString(fmt.Sprintf(" IS %s%s", not, strings.ToUpper(i.IsWithExplicitValue.Value)))
		} else {
			var not string
			if i.IsWithExplicitValue.Not {
				not = "!"
			}
			boolean := "0"
			if i.IsWithExplicitValue.Value == "true" {
				boolean = "1"
			}
			sb.WriteString(fmt.Sprintf(" %s= %s", not, boolean))
		}
	} else if i.IsWithImplicitValue != nil {
		sb.WriteString(i.IsWithImplicitValue.Ident)
		boolean := "1"
		if i.IsWithImplicitValue.Not {
			boolean = "0"
		}
		sb.WriteString(fmt.Sprintf(" = %s", boolean))
	}

	return nil, sb.String()
}

// emitTermOrMath renders tm.
func (cg *SqlCodeGenerator) emitTermOrMath(tm *TermOrMath) (error, string) {
	var err error
	var s string

	if tm.Math != nil {
		err, s = cg.emitMath(tm.Math)
	} else if tm.SubMath != nil {
		if err, s = cg.emitMath(tm.SubMath); err == nil {
			s = fmt.Sprintf("(%s)", s)
		}
	} else if tm.Term != nil {
		err, s = cg.emitTerm(tm.Term)
	}

	return err, s
}

// emitTerm renders t.
func (cg *SqlCodeGenerator) emitTerm(t *Term) (error, string) {
	var err error
	var s string

	if t.Identifier != nil {
		s = *t.Identifier
	} else if t.Integer != nil {
		s = strconv.Itoa(*t.Integer)
	} else if t.Decimal != nil {
		s = strconv.FormatFloat(*t.Decimal, 'f', -1, 64)
	} else if t.String != nil {
		s = fmt.Sprintf("'%s'", *t.String)
	} else if t.Date != nil {
		s = fmt.Sprintf("'%s'", *t.Date)
	} else if t.Time != nil {
		s = fmt.Sprintf("'%s'", *t.Time)
	} else if t.DateTime != nil {
		s = fmt.Sprintf("'%s'", strings.Replace(*t.DateTime, "T", " ", -1))
	} else if t.Bool != nil {
		s = strconv.FormatBool(*t.Bool)
	} else if t.Macro != nil {
		err, s = cg.emitMacro(t.Macro)
	}

	return err, s
}

// emitMath renders m.
func (cg *SqlCodeGenerator) emitMath(m *Math) (error, string) {
	err, t1 := cg.emitTerm(m.Term1)
	if err != nil {
		return err, ""
	}

	err, t2 := cg.emitTerm(m.Term2)
	if err != nil {
		return err, ""
	}

	var op string
	switch m.Op {
	case "add":
		op = "+"
	case "sub":
		op = "-"
	case "mul":
		op = "*"
	case "div":
		op = "/"
	}

	return err, fmt.Sprintf("%s %s %s", t1, op, t2)
}

// emitMacro renders m.
func (cg *SqlCodeGenerator) emitMacro(m *Macro) (error, string) {
	var err error
	var sb strings.Builder

	sb.WriteString(m.Name)

	if m.Args != nil {
		sb.WriteString("(")
		var s string
		for i, a := range m.Args {
			if a.Identifier != nil {
				s = *a.Identifier
			} else if a.Integer != nil {
				s = strconv.Itoa(*a.Integer)
			} else if a.Decimal != nil {
				s = strconv.FormatFloat(*a.Decimal, 'f', -1, 64)
			} else if a.String != nil {
				s = fmt.Sprintf("'%s'", *a.String)
			} else if a.Date != nil {
				s = fmt.Sprintf("'%s'", *a.Date)
			} else if a.Time != nil {
				s = fmt.Sprintf("'%s'", *a.Time)
			} else if a.DateTime != nil {
				s = fmt.Sprintf("'%s'", strings.Replace(*a.DateTime, "T", " ", -1))
			} else if a.Bool != nil {
				s = strconv.FormatBool(*a.Bool)
			} else if a.Macro != nil {
				err, s = cg.emitMacro(a.Macro)
			}

			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(s)
		}

		sb.WriteString(")")
	}

	return err, sb.String()
}

// lookupFieldName gets the native field name associated with f. If no mapping is
// found, then f is returned.
func (cg *SqlCodeGenerator) lookupFieldName(f string) string {
	if val, ok := cg.fieldNames[f]; ok {
		return val
	}

	return f
}
