/**
 * @begin 29-Feb-2020
 * @author <a href="mailto:giuseppe.greco@skeeterhealth.com">Giuseppe Greco</a>
 * @copyright 2020 <a href="skeeterhealth.com">Skeeter</a>
 */

package espressopp_test

import (
	"bytes"
	"strings"
	"testing"

	"gitlab.com/skeeterhealth/espressopp"
)

// Input-result data.
type TestDataItem struct {
	input    string // Espresso++ expression
	result   string // generated sql
	hasError bool   // code generator returned error
}

// Test case for sql generator.
func TestGenerateSql(t *testing.T) {

	dataItems := []TestDataItem{
		{"surname eq 'Walker' and name startswith 'J'", "surname = 'Walker' AND name LIKE 'J%'", false},
		{"age between 20 and 40", "age BETWEEN 20 AND 40", false},
		{"create_time eq #today and (status neq 'processed' or size gt 3000)", "create_time >= TRUNC(SYSDATE) AND create_date < TRUNC(SYSDATE) + 1 AND (status != 'processed' OR size > 3000) ", false},
		{"size gte 2000 and not(create_time lt #now - #duration('PT2H'))", "size >= 2000 AND NOT (create_time < SYSDATE - INTERVAL '2' HOUR)", false},
		{"employee_id eq 110110 and internal is true", "employee_id = 110110 AND internal = 1", false},
		{"employee_id eq 110110 and is internal", "employee_id = 110110 AND internal = 1", false},
		{"customer_note is not null", "customer_note IS NOT NULL", false},
	}

	interpreter := espressopp.NewEspressopp()
	codeGenerator := espressopp.NewSqlCodeGenerator()

	for _, item := range dataItems {
		r := strings.NewReader(item.input)
		w := new(bytes.Buffer)
		err := interpreter.Accept(codeGenerator, r, w)

		result := w.String()

		if item.hasError {
			if err == nil {
				t.Errorf("Interpreter with input '%v' : FAILED, expected an error but got '%v'", item.input, result)
			} else {
				t.Logf("Interpreter with input '%v' : PASSED, expected an error and got '%v'", item.input, err)
			}
		} else {
			if result != item.result {
				t.Errorf("Interpreter with input '%v' : FAILED, expected '%v' but got '%v'", item.input, item.result, result)
			} else {
				t.Logf("Interpreter with input '%v' : PASSED, expected '%v' and got '%v'", item.input, item.result, result)
			}
		}
	}
}
