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
		{"ident eq 'text'", "ident = 'text'", false},
		{"ident neq 10", "ident <> 10", false},
		{"ident neq 'text'", "ident <> 'text'", false},

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
		{"ident between 1 and 10", "ident BETWEEN 1 AND 10", false},

		{"ident startswith 'text'", "ident LIKE 'text%'", false},
		{"ident endswith 'text'", "ident LIKE '%text'", false},
		{"ident contains 'text'", "ident LIKE '%text%'", false},

		{"ident startswith 1", "ident LIKE %1", true},
		{"ident endswith '2020-03-15'", "ident LIKE '%2020-03-15'", true},
		{"ident contains ident", "ident LIKE %ident", true},

		{"ident1 startswith 'text' and (ident2 eq 1 or ident2 gt 10)", "ident1 LIKE 'text%' AND (ident2 = 1 OR ident2 > 10)", false},
		{"ident1 startswith 'text' or (ident2 gte 1 and ident2 lte 10)", "ident1 LIKE 'text%' OR (ident2 >= 1 AND ident2 <= 10)", false},
		{"ident1 startswith 'text' and not (ident2 eq 1 or ident2 gt 10)", "ident1 LIKE 'text%' AND NOT (ident2 = 1 OR ident2 > 10)", false},
		{"ident1 startswith 'text' or not (ident2 gte 1 and ident2 lte 10)", "ident1 LIKE 'text%' OR NOT (ident2 >= 1 AND ident2 <= 10)", false},

		{"ident1 eq (ident2 add 1)", "ident1 = (ident2 + 1)", false},
		{"ident1 eq (ident2 sub 1)", "ident1 = (ident2 - 1)", false},
		{"ident1 eq (ident2 mul 1)", "ident1 = (ident2 * 1)", false},
		{"ident1 eq (ident2 div 1)", "ident1 = (ident2 / 1)", false},

		{"ident1 eq ident2 add 1", "ident1 = ident2 + 1", false},
		{"ident1 eq ident2 sub 1", "ident1 = ident2 - 1", false},
		{"ident1 eq ident2 mul 1", "ident1 = ident2 * 1", false},
		{"ident1 eq ident2 div 1", "ident1 = ident2 / 1", false},

		{"ident eq '2020-03-15'", "ident = '2020-03-15'", false},
		{"ident eq '15:30:55'", "ident = '15:30:55'", false},
		{"ident eq '2020-03-15T14:10:25+02'", "ident = '2020-03-15 14:10:25+02'", false},

		{"ident eq #now", "ident = CURRENT_TIMESTAMP", false},
		{"ident gt #now // this is a comment", "ident > CURRENT_TIMESTAMP", false},
		{"ident lt (#now sub #duration('PT1H'))", "ident < (CURRENT_TIMESTAMP - INTERVAL '1 HOUR')", false},
		{"ident lt (#now add #duration('PT2H'))", "ident < (CURRENT_TIMESTAMP + INTERVAL '2 HOURS')", false},
		{"ident lt (#now add #duration)", "ident < (CURRENT_TIMESTAMP + INTERVAL)", true},
	}
}
