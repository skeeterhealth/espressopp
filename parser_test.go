/**
 * @begin 4-Mar-2020
 * @author <a href="mailto:giuseppe.greco@skeeterhealth.com">Giuseppe Greco</a>
 * @copyright 2020 <a href="skeeterhealth.com">Skeeter</a>
 */

package espressopp

import (
	"strings"
	"testing"

	"gitlab.com/skeeterhealth/espressopp"
)

// Test case for sql generator.
func TestParse(t *testing.T) {

	dataItems := []testDataItem{
		{"ident eq 10", equality, false},
		/*
			{"ident eq 'test'", "ident = 'text'", false},
			{"ident neq 10", "ident <> 10", false},
			{"ident neq 'test'", "ident <> 'text'", false},

			{"ident is true", "ident = 1", false},
			{"ident is false", "ident = 0", false},
			{"ident is null", "ident IS NULL", false},
			{"ident is not null", "ident IS NOT NULL", false},
			{"is ident", "ident = 1", false},
			{"is not ident", "ident = 0", false},

			{"ident gt 10", "ident > 10", false},
			{"ident gte 10", "ident >= 10", false},
			{"ident lt 10", "ident < 10", false},
			{"ident lte 10", "ident <= 10", false},
			{"ident between 1 and 10", "ident BETWEEN 1 and 10", false},

			{"ident startswith 'text'", "ident LIKE 'text%'", false},
			{"ident endswith 'text'", "ident LIKE '%text'", false},
			{"ident contains 'text'", "ident LIKE '%text%'", false},

			{"ident contains 'text'", "ident LIKE '%text%'", false},
			{"ident contains 'text'", "ident LIKE '%text%'", false},
			{"ident contains 'text'", "ident LIKE '%text%'", false},
			{"ident contains 'text'", "ident LIKE '%text%'", false},

			{"ident1 startswith 'text' and (ident2 eq 1 or ident2 gt 10)", "ident LIKE 'text%' AND (ident2 = 1 OR ident2 > 10)", false},
			{"ident1 startswith 'text' or (ident2 gte 1 and ident2 lte 10)", "ident LIKE 'text%' OR (ident2 >= 1 AND ident2 <= 10)", false},
			{"ident1 startswith 'text' and not (ident2 eq 1 or ident2 gt 10)", "ident LIKE 'text%' AND NOT(ident2 = 1 OR ident2 > 10)", false},
			{"ident1 startswith 'text' or not (ident2 gte 1 and ident2 lte 10)", "ident LIKE 'text%' OR NOT (ident2 >= 1 AND ident2 <= 10)", false},

			{"ident eq #today", "ident >= TRUNC(SYSDATE) AND ident < TRUNC(SYSDATE) + 1)", false},
			{"ident lt (#now minus #duration('PT2H'))", "ident < (SYSDATE - INTERVAL) '2' HOUR", false},
			{"ident lt (#now plus #duration('PT2H'))", "ident < (SYSDATE + INTERVAL) '2' HOUR", false},
		*/
	}

	parser := espressopp.newParser()

	for _, item := range dataItems {
		r := strings.NewReader(item.input)
		err, grammar := parser.parse(r)

		result := item.result(grammar)

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

func equality(g *espressopp.Grammar) string {
	return "ident = 10"
}
