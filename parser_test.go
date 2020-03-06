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
	} else if e.Not {
		s = "not "
	} else if e.Paren != nil {
		s = *e.Paren
	} else if e.Equality != nil {
		s = emitEquality(e.Equality)
	} else if e.Comparison != nil {
		s = emitComparison(e.Comparison)
	} else if e.NumericRange != nil {
		s = emitNumericRange(e.NumericRange)
	} else if e.TextualMatching != nil {
		s = emitTextualMatching(e.TextualMatching)
	} else if e.Mathematics != nil {
		s = emitMathematics(e.Mathematics)
	} else if e.Is != nil {
		s = emitIs(e.Is)
	}

	return s
}

// emitEquality renders e.
func emitEquality(e *Equality) string {
	var t1 string
	var t2 string

	if e.Term1.Identifier != nil {
		t1 = *e.Term1.Identifier
	} else if e.Term1.Macro != nil {
		t1 = *e.Term1.Macro
	} else if e.Term1.Integer != nil {
		t1 = strconv.Itoa(*e.Term1.Integer)
	} else if e.Term1.Decimal != nil {
		t1 = strconv.FormatFloat(*e.Term1.Decimal, 'f', -1, 64)
	} else if e.Term1.String != nil {
		t1 = fmt.Sprintf("'%s'", *e.Term1.String)
	} else if e.Term1.Bool != nil {
		t1 = strconv.FormatBool(*e.Term1.Bool)
	}

	if e.Term2.Identifier != nil {
		t2 = *e.Term2.Identifier
	} else if e.Term2.Macro != nil {
		t2 = *e.Term2.Macro
	} else if e.Term2.Integer != nil {
		t2 = strconv.Itoa(*e.Term2.Integer)
	} else if e.Term2.Decimal != nil {
		t2 = strconv.FormatFloat(*e.Term2.Decimal, 'f', -1, 64)
	} else if e.Term2.String != nil {
		t2 = fmt.Sprintf("'%s'", *e.Term2.String)
	} else if e.Term2.Bool != nil {
		t2 = strconv.FormatBool(*e.Term2.Bool)
	}

	return fmt.Sprintf("%s %s %s", t1, e.Op, t2)
}

// emitComparison renders c.
func emitComparison(c *Comparison) string {
	var t1 string
	var t2 string

	if c.Term1.Identifier != nil {
		t1 = *c.Term1.Identifier
	} else if c.Term1.Macro != nil {
		t1 = *c.Term1.Macro
	} else if c.Term1.Integer != nil {
		t1 = strconv.Itoa(*c.Term1.Integer)
	} else if c.Term1.Decimal != nil {
		t1 = strconv.FormatFloat(*c.Term1.Decimal, 'f', -1, 64)
	}

	if c.Term2.Identifier != nil {
		t2 = *c.Term2.Identifier
	} else if c.Term2.Macro != nil {
		t2 = *c.Term2.Macro
	} else if c.Term2.Integer != nil {
		t2 = strconv.Itoa(*c.Term2.Integer)
	} else if c.Term2.Decimal != nil {
		t2 = strconv.FormatFloat(*c.Term2.Decimal, 'f', -1, 64)
	}

	return fmt.Sprintf("%s %s %s", t1, c.Op, t2)
}

// emitNumericRange renders n.
func emitNumericRange(n *NumericRange) string {
	var t1 string
	var t2 string
	var t3 string

	if n.Term1.Identifier != nil {
		t1 = *n.Term1.Identifier
	} else if n.Term1.Macro != nil {
		t1 = *n.Term1.Macro
	} else if n.Term1.Integer != nil {
		t1 = strconv.Itoa(*n.Term1.Integer)
	} else if n.Term1.Decimal != nil {
		t1 = strconv.FormatFloat(*n.Term1.Decimal, 'f', -1, 64)
	}

	if n.Term2.Identifier != nil {
		t2 = *n.Term2.Identifier
	} else if n.Term2.Macro != nil {
		t2 = *n.Term2.Macro
	} else if n.Term2.Integer != nil {
		t2 = strconv.Itoa(*n.Term2.Integer)
	} else if n.Term2.Decimal != nil {
		t2 = strconv.FormatFloat(*n.Term2.Decimal, 'f', -1, 64)
	}

	if n.Term3.Identifier != nil {
		t3 = *n.Term3.Identifier
	} else if n.Term3.Macro != nil {
		t3 = *n.Term3.Macro
	} else if n.Term3.Integer != nil {
		t3 = strconv.Itoa(*n.Term3.Integer)
	} else if n.Term3.Decimal != nil {
		t3 = strconv.FormatFloat(*n.Term3.Decimal, 'f', -1, 64)
	}

	return fmt.Sprintf("%s between %s and %s", t1, t2, t3)
}

// emitTextualMatching renders t.
func emitTextualMatching(t *TextualMatching) string {
	var t1 string
	var t2 string

	if t.Term1.Identifier != nil {
		t1 = *t.Term1.Identifier
	} else if t.Term1.Macro != nil {
		t1 = *t.Term1.Macro
	} else if t.Term1.String != nil {
		t1 = fmt.Sprintf("'%s'", *t.Term1.String)
	}

	if t.Term2.Identifier != nil {
		t2 = *t.Term2.Identifier
	} else if t.Term2.Macro != nil {
		t2 = *t.Term2.Macro
	} else if t.Term2.String != nil {
		t2 = fmt.Sprintf("'%s'", *t.Term2.String)
	}

	return fmt.Sprintf("%s %s %s", t1, t.Op, t2)
}

// emitMathematics renders m.
func emitMathematics(m *Mathematics) string {
	var t1 string
	var t2 string

	if m.Term1.Identifier != nil {
		t1 = *m.Term1.Identifier
	} else if m.Term1.Macro != nil {
		t1 = *m.Term1.Macro
	} else if m.Term1.Integer != nil {
		t1 = strconv.Itoa(*m.Term1.Integer)
	} else if m.Term1.Decimal != nil {
		t1 = strconv.FormatFloat(*m.Term1.Decimal, 'f', -1, 64)
	}

	if m.Term2.Identifier != nil {
		t2 = *m.Term2.Identifier
	} else if m.Term2.Macro != nil {
		t2 = *m.Term2.Macro
	} else if m.Term2.Integer != nil {
		t2 = strconv.Itoa(*m.Term2.Integer)
	} else if m.Term2.Decimal != nil {
		t2 = strconv.FormatFloat(*m.Term2.Decimal, 'f', -1, 64)
	}

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
