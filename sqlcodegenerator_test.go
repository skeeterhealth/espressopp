/**
 * @begin 2020-02-29
 * @author <a href="mailto:giuseppe.greco@skeeterhealth.com">Giuseppe Greco</a>
 * @copyright 2020 <a href="skeeterhealth.com">Skeeter</a>
 */

package espressopp

import (
	"bytes"
	"strings"
	"testing"
)

// TestGenerateSql tests the generation of SQL from Espresso++ expressions.
func TestGenerateSql(t *testing.T) {
	interpreter := NewEspressoppInterpreter()
	codeGenerator := NewSqlCodeGenerator()

	for _, item := range getTestDataItems() {
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
