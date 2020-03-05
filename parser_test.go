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

		result := emitExpressions(grammar)

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

// emitExpressions renders the expressions in g.
func emitExpressions(g *Grammar) string {
	var sb strings.Builder

	for _, expression := range g.Expressions {
		if expression.Op1 != nil {
			sb.WriteString(*expression.Op1)
			sb.WriteString(" ")
		}

		if expression.Op2 != nil {
			sb.WriteString(*expression.Op2)
			sb.WriteString(" ")
		}

		if expression.Equality != nil {
			sb.WriteString(emitEquality(expression.Equality))
		} else if expression.Comparison != nil {
			sb.WriteString(emitComparison(expression.Comparison))
		} else if expression.NumericRange != nil {
			sb.WriteString(emitNumericRange(expression.NumericRange))
		} else if expression.TextualMatching != nil {
			sb.WriteString(emitTextualMatching(expression.TextualMatching))
		} else if expression.Mathematics != nil {
			sb.WriteString(emitMathematics(expression.Mathematics))
		} else if expression.Is != nil {
			sb.WriteString(emitIs(expression.Is))
		} else if expression.ParenExpression != nil {
			sb.WriteString(emitParenExpression(expression.ParenExpression))
		}

		sb.WriteString(" ")
	}

	return strings.Trim(sb.String(), " ")
}

// emitEquality renders e.
func emitEquality(e *Equality) string {
	var t1 string
	var t2 string

	if e.Term1.Identifier != nil {
		t1 = *e.Term1.Identifier
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
	} else if c.Term1.Integer != nil {
		t1 = strconv.Itoa(*c.Term1.Integer)
	} else if c.Term1.Decimal != nil {
		t1 = strconv.FormatFloat(*c.Term1.Decimal, 'f', -1, 64)
	}

	if c.Term2.Identifier != nil {
		t2 = *c.Term2.Identifier
	} else if c.Term2.Integer != nil {
		t2 = strconv.Itoa(*c.Term2.Integer)
	} else if c.Term2.Decimal != nil {
		t2 = strconv.FormatFloat(*c.Term2.Decimal, 'f', -1, 64)
	}

	return fmt.Sprintf("%s %s %s", t1, c.Op, t2)
}

// emitNumericRange renders n.
func emitNumericRange(n *NumericRange) string {
	return ""
}

// emitTextualMatching renders t.
func emitTextualMatching(t *TextualMatching) string {
	return ""
}

// emitMathematics renders m.
func emitMathematics(m *Mathematics) string {
	return ""
}

// emitIs renders i.
func emitIs(i *Is) string {
	var sb strings.Builder

	if i.Ident != nil {
		sb.WriteString(*i.Ident)
	}

	sb.WriteString(" is")

	if i.Op != nil {
		sb.WriteString(" ")
		sb.WriteString(*i.Op)
	}

	sb.WriteString(" ")
	sb.WriteString(i.Value)

	return sb.String()
}

// emitParenExpression renders p.
func emitParenExpression(p *ParenExpression) string {
	return ""
}
