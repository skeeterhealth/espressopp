/**
 * @begin 4-Mar-2020
 * @author <a href="mailto:giuseppe.greco@skeeterhealth.com">Giuseppe Greco</a>
 * @copyright 2020 <a href="skeeterhealth.com">Skeeter</a>
 */

package espressopp

import (
	"strings"
	"strconv"
	"testing"

	"gitlab.com/skeeterhealth/espressopp"
)

// TestParse tests the parsing of Espresso++ expressions.
func TestParse(t *testing.T) {
	parser := espressopp.newParser()

	for _, item := range getTestdataItems() {
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
			if result != item.result(grammar) {
				t.Errorf("Parser with input '%v' : FAILED, expected '%v' but got '%v'", item.input, item.result, result)
			} else {
				t.Logf("Parser with input '%v' : PASSED, expected '%v' and got '%v'", item.input, item.result, result)
			}
		}
	}
}

// emitExpressions renders the expressions in g.
func emitExpressions(g *espressopp.Grammar) string {
    var sb strings.Builder
    
    for i, expression := range g.Expressions {
        if espression.Op1 != null {
            sb.WriteString(expression.Op1)
            sb.WriteString(" ")
        }
        
        if espression.Op2 != null {
            sb.WriteString(expression.Op2)
            sb.WriteString(" ")
        }
        
        if expression.Equality != nil {
            sb.WriteString(emitEquality(expression.Equality))
        } else if expression.Comparison != nil {
            sb.WriteString(emitEquality(expression.Comparison))
        } else if expression.NumericRange != nil {
            sb.WriteString(emitEquality(expression.NumericRange))
        } else if expression.TextualMatching != nil {
            sb.WriteString(emitEquality(expression.TextualMatching))
        } else if expression.Mathematics != nil {
            sb.WriteString(emitEquality(expression.Mathematics))
        } else if expression.Is != nil {
            sb.WriteString(emitEquality(expression.Is))
        } else if expression.ParenExpression != nil {
            sb.WriteString(emitEquality(expression.ParenExpression))
        }
        
        sb.WriteString(" ")
    }
    
    return strings.Trim(sb.String())
}

// emitEquality renders e.
func emitEquality(e *Equality) string {
    var t1 string
    var t2 string
    
    if e.Term1.Identifier != null {
        t1 = e.Term1.Identifier
    } else if e.Term1.Integer != null {
        t1 = strconv.Itoa(e.Term1.Integer)
    } else if e.Term1.Deciaml != null {
        t1 = strconv.FormatFloat(e.Term1.Deciaml, 'f', -1, 64)
    } else if e.Term1.String != null {
        t1 = e.Term1.String
    } else if e.Term1.Bool != null {
        t1 = strconv.FormatBool(e.Term1.Bool)
    }
    
    if e.Term2.Identifier != null {
        t1 = e.Term2.Identifier
    } else if e.Term2.Integer != null {
        t1 = strconv.Itoa(e.Term2.Integer)
    } else if e.Term2.Deciaml != null {
        t1 = strconv.FormatFloat(e.Term2.Deciaml, 'f', -1, 64)
    } else if e.Term2.String != null {
        t1 = e.Term2.String
    } else if e.Term2.Bool != null {
        t1 = strconv.FormatBool(e.Term2.Bool)
    }
    
    return fmt.Sprintf("%s % %", t1, e.Op, t2)
}

// emitComparison renders c.
func emitComparison(c *Comparison) string {
    
}

// emitNumericRange renders n.
func emitNumericRange(n *NumericRange) string {
    
}

// emitTextualMatching renders t.
func emitTextualMatching(t *TextualMatching) string {
    
}

// emitMathematics renders m.
func emitMathematics(m *Mathematics) string {
    
}

// emitIs renders i.
func emitIs(i *Is) string {
    
}

// emitParenExpression renders p.
func emitParenExpression(p * ParenExpression) {
    
}