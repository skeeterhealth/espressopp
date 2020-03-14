/**
 * @begin 2020-03-04
 * @author <a href="mailto:giuseppe.greco@skeeterhealth.com">Giuseppe Greco</a>
 * @copyright 2020 <a href="skeeterhealth.com">Skeeter</a>
 */

package espressopp

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

// TestParse tests the parsing of Espresso++ expressions.
func TestParse(t *testing.T) {
	parser := newParser()

	for _, item := range getTestDataItems() {
		r := strings.NewReader(item.input)
		err, grammar := parser.parse(r)

		result := emitGrammar(grammar)

		if item.hasError {
			if err == nil {
				t.Errorf("Parser with input '%v' : FAILED, expected an error but got '%v'", item.input, result)
			} else {
				t.Logf("Parser with input '%v' : PASSED, expected an error and got '%v'", item.input, err)
			}
		} else {
			if result != item.input {
				t.Errorf("Parser with input '%v' : FAILED, expected '%v' but got '%v'", item.input, item.input, result)
			} else {
				t.Logf("Parser with input '%v' : PASSED, expected '%v' and got '%v'", item.input, item.input, result)
			}
		}
	}
}

// emitGrammars renders the expressions in g.
func emitGrammar(g *Grammar) string {
	var sb strings.Builder

	for _, e := range g.Expressions {
		sb.WriteString(emitExpression(e))
	}

	return sb.String()
}

// emitExpression renders e.
func emitExpression(e *Expression) string {
	var s string

	if e.Op != nil {
		s = fmt.Sprintf(" %s ", *e.Op)
	} else if e.SubExpression != nil {
		s = emitSubExpression(e.SubExpression)
	} else if e.Comparison != nil {
		s = emitComparison(e.Comparison)
	} else if e.Equality != nil {
		s = emitEquality(e.Equality)
	} else if e.Range != nil {
		s = emitRange(e.Range)
	} else if e.Match != nil {
		s = emitMatch(e.Match)
	} else if e.Is != nil {
		s = emitIs(e.Is)
	}

	return s
}

// emitSubExpression renders pe.
func emitSubExpression(pe *SubExpression) string {
	var sb strings.Builder

	if pe.Not {
		sb.WriteString("not ")
	}

	sb.WriteString("(")

	for _, e := range pe.Expressions {
		sb.WriteString(emitExpression(e))
	}

	sb.WriteString(")")

	return sb.String()
}

// emitComparison renders c.
func emitComparison(c *Comparison) string {
	t1 := emitTermOrMath(c.TermOrMath1)
	t2 := emitTermOrMath(c.TermOrMath2)

	return fmt.Sprintf("%s %s %s", t1, c.Op, t2)
}

// emitEquality renders e.
func emitEquality(e *Equality) string {
	t1 := emitTermOrMath(e.TermOrMath1)
	t2 := emitTermOrMath(e.TermOrMath2)

	return fmt.Sprintf("%s %s %s", t1, e.Op, t2)
}

// emitRange renders r.
func emitRange(r *Range) string {
	t1 := emitTermOrMath(r.TermOrMath1)
	t2 := emitTermOrMath(r.TermOrMath2)
	t3 := emitTermOrMath(r.TermOrMath3)

	return fmt.Sprintf("%s between %s and %s", t1, t2, t3)
}

// emitMatch renders m.
func emitMatch(m *Match) string {
	t1 := emitTerm(m.Term1)
	t2 := emitTerm(m.Term2)

	return fmt.Sprintf("%s %s %s", t1, m.Op, t2)
}

// emitIs renders i.
func emitIs(i *Is) string {
	var sb strings.Builder

	if i.IsWithExplicitValue != nil {
		sb.WriteString(i.IsWithExplicitValue.Ident)
		sb.WriteString(" is ")
		if i.IsWithExplicitValue.Not {
			sb.WriteString("not ")
		}
		sb.WriteString(i.IsWithExplicitValue.Value)
	} else if i.IsWithImplicitValue != nil {
		sb.WriteString("is ")
		if i.IsWithImplicitValue.Not {
			sb.WriteString("not ")
		}
		sb.WriteString(i.IsWithImplicitValue.Ident)
	}

	return sb.String()
}

// emitMacro renders m.
func emitMacro(m *Macro) string {
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
			} else if a.Bool != nil {
				s = strconv.FormatBool(*a.Bool)
			} else if a.Macro != nil {
				s = emitMacro(a.Macro)
			}

			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(s)
		}

		sb.WriteString(")")
	}

	return sb.String()
}

// emitMath renders m.
func emitMath(m *Math) string {
	t1 := emitTerm(m.Term1)
	t2 := emitTerm(m.Term2)

	return fmt.Sprintf("%s %s %s", t1, m.Op, t2)
}

// emitTerm renders t.
func emitTerm(t *Term) string {
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
		s = fmt.Sprintf("'%s'", *t.DateTime)
	} else if t.Bool != nil {
		s = strconv.FormatBool(*t.Bool)
	} else if t.Macro != nil {
		s = emitMacro(t.Macro)
	}

	return s
}

// emitTermOrMath renders tm.
func emitTermOrMath(tm *TermOrMath) string {
	var s string

	if tm.Math != nil {
		s = emitMath(tm.Math)
	} else if tm.SubMath != nil {
		s = fmt.Sprintf("(%s)", emitMath(tm.SubMath))
	} else if tm.Term != nil {
		s = emitTerm(tm.Term)
	}

	return s
}
