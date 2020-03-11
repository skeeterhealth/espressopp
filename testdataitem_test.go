/**
 * @begin 2020-03-05
 * @author <a href="mailto:giuseppe.greco@skeeterhealth.com">Giuseppe Greco</a>
 * @copyright 2020 <a href="skeeterhealth.com">Skeeter</a>
 */

package espressopp

// testDataItem defines test data.
type testDataItem struct {
	input    string // input data
	result   string // test result
	hasError bool   // whether the test returned an error
}

// getTestDataItems returns an array of testDataItem structs with predefined
// test data.
func getTestDataItems() []testDataItem {
	return []testDataItem{
		{"ident eq 10", "ident = 10", false},
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

		{"ident1 eq (ident2 plus 1)", "ident1 = (ident2 + 1)", false},
		{"ident1 eq (ident2 minus 1)", "ident1 = (ident2 - 1)", false},
		{"ident1 eq (ident2 mul 1)", "ident1 = (ident2 * 1)", false},
		{"ident1 eq (ident2 div 1)", "ident1 = (ident2 / 1)", false},

		{"ident1 eq ident2 plus 1", "ident1 = ident2 + 1", false},
		{"ident1 eq ident2 minus 1", "ident1 = ident2 - 1", false},
		{"ident1 eq ident2 mul 1", "ident1 = ident2 * 1", false},
		{"ident1 eq ident2 div 1", "ident1 = ident2 / 1", false},

		{"ident eq #today", "ident >= TRUNC(SYSDATE) AND ident < TRUNC(SYSDATE) + 1)", false},
		{"ident lt (#now minus #duration('PT2H'))", "ident < (SYSDATE - INTERVAL) '2' HOUR", false},
		{"ident lt (#now plus #duration('PT2H'))", "ident < (SYSDATE + INTERVAL) '2' HOUR", false},
	}
}
