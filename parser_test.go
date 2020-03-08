/**
 * @begin 4-Mar-2020
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
	} else if e.ParenExpression != nil {
		s = emitParenExpression(e.ParenExpression)
	} else if e.Equality != nil {
		s = emitEquality(e.Equality)
	} else if e.Comparison != nil {
		s = emitComparison(e.Comparison)
	} else if e.Range != nil {
		s = emitRange(e.Range)
	} else if e.Matching != nil {
		s = emitMatching(e.Matching)
	} else if e.Is != nil {
		s = emitIs(e.Is)
	}

	return s
}

// emitParenExpression renders pe.
func emitParenExpression(pe *ParenExpression) string {
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

// emitEquality renders e.
func emitEquality(e *Equality) string {
	t1 := emitTermOrMaths(e.TermOrMaths1)
	t2 := emitTermOrMaths(e.TermOrMaths2)

	return fmt.Sprintf("%s %s %s", t1, e.Op, t2)
}

// emitComparison renders c.
func emitComparison(c *Comparison) string {
	fmt.Println(fmt.Sprintf("Term1 = %v", c.TermOrMaths1))
	fmt.Println(fmt.Sprintf("Term2 = %v", c.TermOrMaths2))
	t1 := emitTermOrMaths(c.TermOrMaths1)
	t2 := emitTermOrMaths(c.TermOrMaths2)
	return fmt.Sprintf("%s %s %s", t1, c.Op, t2)
}

// emitRange renders r.
func emitRange(r *Range) string {
	t1 := emitTermOrMaths(r.TermOrMaths1)
	t2 := emitTermOrMaths(r.TermOrMaths2)
	t3 := emitTermOrMaths(r.TermOrMaths3)

	return fmt.Sprintf("%s between %s and %s", t1, t2, t3)
}

// emitMatching renders m.
func emitMatching(m *Matching) string {
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
			} else if a.Macro != nil {
				s = emitMacro(a.Macro)
			} else if a.Integer != nil {
				s = strconv.Itoa(*a.Integer)
			} else if a.Decimal != nil {
				s = strconv.FormatFloat(*a.Decimal, 'f', -1, 64)
			} else if a.String != nil {
				s = fmt.Sprintf("'%s'", *a.String)
			} else if a.Bool != nil {
				s = strconv.FormatBool(*a.Bool)
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

// emitMaths renders m.
func emitMaths(m *Maths) string {
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
	} else if t.Bool != nil {
		s = strconv.FormatBool(*t.Bool)
	} else if t.Macro != nil {
		s = emitMacro(t.Macro)
	}

	return s
}

// emitTermOrMaths renders tm.
func emitTermOrMaths(tm *TermOrMaths) string {
	var s string

	if tm.Maths != nil {
		s = emitMaths(tm.Maths)
	} else if tm.ParenMaths != nil {
		s = fmt.Sprintf("(%s)", emitMaths(tm.ParenMaths))
	} else if tm.Term != nil {
		s = emitTerm(tm.Term)
	}

	return s
}
