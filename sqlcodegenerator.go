/**
 * @begin 2020-02-15
 * @author <a href="mailto:giuseppe.greco@skeeterhealth.com">Giuseppe Greco</a>
 * @copyright 2020 <a href="skeeterhealth.com">Skeeter</a>
 */

package espressopp

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	duration "github.com/channelmeter/iso8601duration"
	pluralize "github.com/gertd/go-pluralize"
	"github.com/pkg/errors"
)

type termType int

const (
	undefType termType = iota
	identType
	intType
	decimalType
	stringType
	dateType
	timeType
	dateTimeType
	boolType
)

// SqlCodeGenerator is the CodeGenerator implementation that produces native SQL
// from Espresso++ expressions.
type SqlCodeGenerator struct {
	renderingOptions *RenderingOptions
}

// NewSqlCodeGenerator creates a new instance of SqlCodeGenerator.
func NewSqlCodeGenerator() CodeGenerator {
	return &SqlCodeGenerator{
		renderingOptions: NewRenderingOptions(),
	}
}

// GetRenderingOptions gets the rendering options used by the SqlCodeGenerator.
func (cg *SqlCodeGenerator) GetRenderingOptions() *RenderingOptions {
	return cg.renderingOptions
}

// Visit lets cg access the functionality provided by i to parse the Espresso++
// expressions in r and get back the grammar, which is then used to produce native
// SQL into w.
func (cg *SqlCodeGenerator) Visit(i Interpreter, r io.Reader, w io.Writer) error {
	if i == nil {
		return errors.New("interpreter not specified")
	}

	err, grammar := i.Parse(r)
	if err != nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r)
		return errors.Wrapf(err, "error parsing %v", buf.String())
	}

	err, s := cg.emitGrammar(grammar)
	if err != nil {
		return errors.Wrapf(err, "error generating sql")
	}

	_, err = io.WriteString(w, s)
	return err
}

// emitGrammar renders g.
func (cg *SqlCodeGenerator) emitGrammar(g *Grammar) (error, string) {
	var err error
	var sb strings.Builder

	for _, e := range g.Expressions {
		err, s := cg.emitExpression(e)
		if err != nil {
			return err, ""
		}
		sb.WriteString(s)
	}

	return err, sb.String()
}

// emitExpression renders e.
func (cg *SqlCodeGenerator) emitExpression(e *Expression) (error, string) {
	var err error
	var s string

	if e.Op != nil {
		s = fmt.Sprintf(" %s ", strings.ToUpper(*e.Op))
	} else if e.SubExpression != nil {
		err, s = cg.emitSubExpression(e.SubExpression)
	} else if e.Comparison != nil {
		err, s = cg.emitComparison(e.Comparison)
	} else if e.Equality != nil {
		err, s = cg.emitEquality(e.Equality)
	} else if e.Range != nil {
		err, s = cg.emitRange(e.Range)
	} else if e.Match != nil {
		err, s = cg.emitMatch(e.Match)
	} else if e.Is != nil {
		err, s = cg.emitIs(e.Is)
	}

	return err, s
}

// emitExpression renders se.
func (cg *SqlCodeGenerator) emitSubExpression(se *SubExpression) (error, string) {
	var err error
	var sb strings.Builder

	if se.Not {
		sb.WriteString("NOT ")
	}

	sb.WriteString("(")

	for _, e := range se.Expressions {
		err, s := cg.emitExpression(e)
		if err != nil {
			return err, ""
		}
		sb.WriteString(s)
	}

	sb.WriteString(")")

	return err, sb.String()
}

// emitComparison renders c.
func (cg *SqlCodeGenerator) emitComparison(c *Comparison) (error, string) {
	err, t1, tt1 := cg.emitTermOrMath(c.TermOrMath1)
	if err != nil {
		return err, ""
	}

	err, t2, tt2 := cg.emitTermOrMath(c.TermOrMath2)
	if err != nil {
		return err, ""
	}

	err, tt := cg.validateTypes(tt1, tt2)
	if err != nil {
		return err, ""
	} else if tt != intType && tt != decimalType && tt != dateType && tt != timeType && tt != dateTimeType {
		return errors.Errorf("cannot compare values of type %s", cg.toTypeName(tt)), ""
	}

	var op string
	switch c.Op {
	case "gt":
		op = ">"
	case "gte":
		op = ">="
	case "lt":
		op = "<"
	case "lte":
		op = "<="
	}

	return err, fmt.Sprintf("%s %s %s", t1, op, t2)
}

// emitEquality renders e.
func (cg *SqlCodeGenerator) emitEquality(e *Equality) (error, string) {
	err, t1, tt1 := cg.emitTermOrMath(e.TermOrMath1)
	if err != nil {
		return err, ""
	}

	err, t2, tt2 := cg.emitTermOrMath(e.TermOrMath2)
	if err != nil {
		return err, ""
	}

	err, _ = cg.validateTypes(tt1, tt2)
	if err != nil {
		return err, ""
	}

	var op string
	switch e.Op {
	case "eq":
		op = "="
	case "neq":
		op = "<>"
	}

	return err, fmt.Sprintf("%s %s %s", t1, op, t2)
}

// emitRange renders r.
func (cg *SqlCodeGenerator) emitRange(r *Range) (error, string) {
	err, t1, tt1 := cg.emitTermOrMath(r.TermOrMath1)
	if err != nil {
		return err, ""
	}

	err, t2, tt2 := cg.emitTermOrMath(r.TermOrMath2)
	if err != nil {
		return err, ""
	}

	err, t3, tt3 := cg.emitTermOrMath(r.TermOrMath3)
	if err != nil {
		return err, ""
	}

	err, tt := cg.validateTypes(tt1, tt2)
	if err != nil {
		return err, ""
	}
	err, tt = cg.validateTypes(tt, tt3)
	if err != nil {
		return err, ""
	} else if tt != intType && tt != decimalType {
		return errors.Errorf("cannot range values of type %s", cg.toTypeName(tt)), ""
	}

	return err, fmt.Sprintf("%s %s %s %s %s", t1,
		strings.ToUpper(r.Between), t2,
		strings.ToUpper(r.And), t3)
}

// emitMatch renders m.
func (cg *SqlCodeGenerator) emitMatch(m *Match) (error, string) {
	err, t1, tt1 := cg.emitTerm(m.Term1)
	if err != nil {
		return err, ""
	}

	err, t2, tt2 := cg.emitTerm(m.Term2)
	if err != nil {
		return err, ""
	}

	err, tt := cg.validateTypes(tt1, tt2)
	if err != nil {
		return err, ""
	} else if tt != stringType {
		return errors.Errorf("cannot match values of type %s", cg.toTypeName(tt)), ""
	}

	t2 = strings.ReplaceAll(strings.ReplaceAll(t2, "'", ""), "\"", "")

	switch m.Op {
	case "startswith":
		t2 = fmt.Sprintf("'%s%%'", t2)
	case "endswith":
		t2 = fmt.Sprintf("'%%%s'", t2)
	case "contains":
		t2 = fmt.Sprintf("'%%%s%%'", t2)
	}

	return err, fmt.Sprintf("%s %s %s", t1, "LIKE", t2)
}

// emitIs renders i.
func (cg *SqlCodeGenerator) emitIs(i *Is) (error, string) {
	var sb strings.Builder

	if i.IsWithExplicitValue != nil {
		sb.WriteString(i.IsWithExplicitValue.Ident)
		if i.IsWithExplicitValue.Value == "null" {
			var not string
			if i.IsWithExplicitValue.Not {
				not = "NOT "
			}
			sb.WriteString(fmt.Sprintf(" IS %s%s", not, strings.ToUpper(i.IsWithExplicitValue.Value)))
		} else {
			var not string
			if i.IsWithExplicitValue.Not {
				not = "!"
			}
			boolean := "0"
			if i.IsWithExplicitValue.Value == "true" {
				boolean = "1"
			}
			sb.WriteString(fmt.Sprintf(" %s= %s", not, boolean))
		}
	} else if i.IsWithImplicitValue != nil {
		sb.WriteString(i.IsWithImplicitValue.Ident)
		boolean := "1"
		if i.IsWithImplicitValue.Not {
			boolean = "0"
		}
		sb.WriteString(fmt.Sprintf(" = %s", boolean))
	}

	return nil, sb.String()
}

// emitTermOrMath renders tm.
func (cg *SqlCodeGenerator) emitTermOrMath(tm *TermOrMath) (error, string, termType) {
	var err error
	var s string
	var t termType

	if tm.Math != nil {
		err, s, t = cg.emitMath(tm.Math)
	} else if tm.SubMath != nil {
		if err, s, t = cg.emitMath(tm.SubMath); err == nil {
			s = fmt.Sprintf("(%s)", s)
		}
	} else if tm.Term != nil {
		err, s, t = cg.emitTerm(tm.Term)
	}

	return err, s, t
}

// emitTerm renders t.
func (cg *SqlCodeGenerator) emitTerm(t *Term) (error, string, termType) {
	var err error
	var s string
	var tt termType

	if t.Identifier != nil {
		s = cg.applyRenderingOptions(*t.Identifier)
		tt = identType
	} else if t.Integer != nil {
		s = strconv.Itoa(*t.Integer)
		tt = intType
	} else if t.Decimal != nil {
		s = strconv.FormatFloat(*t.Decimal, 'f', -1, 64)
		tt = decimalType
	} else if t.String != nil {
		s = fmt.Sprintf("'%s'", *t.String)
		tt = stringType
	} else if t.Date != nil {
		s = fmt.Sprintf("'%s'", *t.Date)
		tt = dateType
	} else if t.Time != nil {
		s = fmt.Sprintf("'%s'", *t.Time)
		tt = timeType
	} else if t.DateTime != nil {
		s = fmt.Sprintf("'%s'", strings.Replace(*t.DateTime, "T", " ", -1))
		tt = dateTimeType
	} else if t.Bool != nil {
		s = strconv.FormatBool(*t.Bool)
		tt = boolType
	} else if t.Macro != nil {
		err, s, tt = cg.emitMacro(t.Macro)
	}

	return err, s, tt
}

// emitMath renders m.
func (cg *SqlCodeGenerator) emitMath(m *Math) (error, string, termType) {
	err, t1, tt1 := cg.emitTerm(m.Term1)
	if err != nil {
		return err, "", undefType
	}

	err, t2, tt2 := cg.emitTerm(m.Term2)
	if err != nil {
		return err, "", undefType
	}

	err, tt := cg.validateTypes(tt1, tt2)
	if err != nil {
		return err, "", undefType
	} else if tt != intType && tt != decimalType && tt != dateType && tt != timeType && tt != dateTimeType {
		return errors.Errorf("cannot compute values of type %s", cg.toTypeName(tt)), "", tt
	}

	if tt == dateType || tt == timeType || tt == dateTimeType {
		var f string
		switch tt {
		case dateType:
			f = "DATE"
		case timeType:
			f = "TIME"
		case dateTimeType:
			f = "TIMESTAMP"
		}

		if strings.HasPrefix(t1, "'") && strings.HasSuffix(t1, "'") {
			t1 = fmt.Sprintf("%s %s", f, t1)
		}
		if strings.HasPrefix(t2, "'") && strings.HasSuffix(t2, "'") {
			t2 = fmt.Sprintf("%s %s", f, t2)
		}
	}

	var op string
	switch m.Op {
	case "add":
		op = "+"
	case "sub":
		op = "-"
	case "mul":
		op = "*"
	case "div":
		op = "/"
	}

	return err, fmt.Sprintf("%s %s %s", t1, op, t2), tt
}

// emitMacro renders m.
func (cg *SqlCodeGenerator) emitMacro(m *Macro) (error, string, termType) {
	var err error
	var s string
	var t termType

	switch m.Name {
	case "#now":
		err, s, t = cg.emitNowMacro(m)
	case "#duration":
		err, s, t = cg.emitDurationMacro(m)
	}

	return err, s, t
}

// emitNowMacro renders m.
func (cg *SqlCodeGenerator) emitNowMacro(m *Macro) (error, string, termType) {
	return nil, "CURRENT_TIMESTAMP", dateTimeType
}

// emitDurationMacro renders m.
func (cg *SqlCodeGenerator) emitDurationMacro(m *Macro) (error, string, termType) {
	if m.Args == nil {
		return errors.Errorf("%s: missing parameter: iso8601 interval", m.Name), "", undefType
	}

	const interval = "INTERVAL '%s'"
	var err error
	var items []string
	p := pluralize.NewClient()

	for _, a := range m.Args {
		err, s, t := cg.emitTerm(a)
		if err != nil {
			return err, "", undefType
		} else if t != stringType {
			return errors.Errorf("iso8601 interval cannot be of type %s", cg.toTypeName(t)), "", undefType
		}
		d, err := duration.FromString(s)
		if err != nil {
			return err, "", undefType
		}
		if d.Years > 0 {
			items = append(items, fmt.Sprintf(interval, p.Pluralize("YEAR", d.Years, true)))
		}
		if d.Weeks > 0 {
			items = append(items, fmt.Sprintf(interval, p.Pluralize("WEEK", d.Weeks, true)))
		}
		if d.Days > 0 {
			items = append(items, fmt.Sprintf(interval, p.Pluralize("DAY", d.Days, true)))
		}
		if d.Hours > 0 {
			items = append(items, fmt.Sprintf(interval, p.Pluralize("HOUR", d.Hours, true)))
		}
		if d.Minutes > 0 {
			items = append(items, fmt.Sprintf(interval, p.Pluralize("MINUTE", d.Minutes, true)))
		}
		if d.Seconds > 0 {
			items = append(items, fmt.Sprintf(interval, p.Pluralize("SECOND", d.Seconds, true)))
		}
	}

	var sb strings.Builder
	for i, s := range items {
		sb.WriteString(s)
		if i > 0 {
			sb.WriteString(" ")
		}
	}

	return err, sb.String(), dateTimeType
}

// validateTypes verifies whether or not t1 and t2 are compatible, and if they are,
// it returns the result type of the current expression.
func (cg *SqlCodeGenerator) validateTypes(t1 termType, t2 termType) (error, termType) {
	var err error
	var t termType

	if t1 == t2 {
		t = t1
	} else if t1 == identType {
		t = t2
	} else if t2 == identType {
		t = t1
	} else {
		err = errors.Errorf("type %s is not compatible with type %s", cg.toTypeName(t1), cg.toTypeName(t2))
	}

	return err, t
}

// toTypeName returns the name of t.
func (cg *SqlCodeGenerator) toTypeName(t termType) string {
	var n string

	switch t {
	case stringType:
		n = "string"
	case intType:
		n = "int"
	case decimalType:
		n = "decimal"
	case dateType:
		n = "date"
	case timeType:
		n = "time"
	case dateTimeType:
		n = "datetime"
	case boolType:
		n = "bool"
	}

	return n
}

// applyRenderingOptions applies the rendering options to f.
func (cg *SqlCodeGenerator) applyRenderingOptions(f string) string {
	if val, ok := cg.renderingOptions.fields[f]; ok {
		return val.NativeName
	}

	return f
}
