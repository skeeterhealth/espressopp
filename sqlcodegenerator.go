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
	// RenderingOptions is used to control the way native SQL is produced.
	RenderingOptions *RenderingOptions
}

// NewSqlCodeGenerator creates a new instance of SqlCodeGenerator.
func NewSqlCodeGenerator() *SqlCodeGenerator {
	return &SqlCodeGenerator{
		RenderingOptions: NewRenderingOptions(),
	}
}

// Visit lets cg access the functionality provided by i to parse the Espresso++
// expressions in r and get back the grammar, which is then used to produce native
// SQL into w.
func (cg *SqlCodeGenerator) Visit(i Interpreter, r io.Reader, w io.Writer) error {
	if i == nil {
		return errors.New("interpreter not specified")
	}

	grammar, err := i.Parse(r)
	if err != nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r)
		return errors.Wrapf(err, "error parsing %v", buf.String())
	}

	s, err := cg.emitGrammar(grammar)
	if err != nil {
		return errors.Wrapf(err, "error generating sql")
	}

	_, err = io.WriteString(w, s)
	return err
}

// emitGrammar renders g.
func (cg *SqlCodeGenerator) emitGrammar(g *Grammar) (string, error) {
	var err error
	var sb strings.Builder

	for _, e := range g.Expressions {
		s, err := cg.emitExpression(e)
		if err != nil {
			return "", err
		}
		sb.WriteString(s)
	}

	return sb.String(), err
}

// emitExpression renders e.
func (cg *SqlCodeGenerator) emitExpression(e *Expression) (string, error) {
	var err error
	var s string

	if e.Op != nil {
		s = fmt.Sprintf(" %s ", strings.ToUpper(*e.Op))
	} else if e.SubExpression != nil {
		s, err = cg.emitSubExpression(e.SubExpression)
	} else if e.Comparison != nil {
		s, err = cg.emitComparison(e.Comparison)
	} else if e.Equality != nil {
		s, err = cg.emitEquality(e.Equality)
	} else if e.Range != nil {
		s, err = cg.emitRange(e.Range)
	} else if e.Match != nil {
		s, err = cg.emitMatch(e.Match)
	} else if e.Is != nil {
		s, err = cg.emitIs(e.Is)
	}

	return s, err
}

// emitExpression renders se.
func (cg *SqlCodeGenerator) emitSubExpression(se *SubExpression) (string, error) {
	var err error
	var sb strings.Builder

	if se.Not {
		sb.WriteString("NOT ")
	}

	sb.WriteString("(")

	for _, e := range se.Expressions {
		s, err := cg.emitExpression(e)
		if err != nil {
			return "", err
		}
		sb.WriteString(s)
	}

	sb.WriteString(")")

	return sb.String(), err
}

// emitComparison renders c.
func (cg *SqlCodeGenerator) emitComparison(c *Comparison) (string, error) {
	t1, tt1, err := cg.emitTermOrMath(c.TermOrMath1)
	if err != nil {
		return "", err
	}

	t2, tt2, err := cg.emitTermOrMath(c.TermOrMath2)
	if err != nil {
		return "", err
	}

	tt, err := cg.validateTypes(tt1, tt2)
	if err != nil {
		return "", err
	} else if tt != intType && tt != decimalType && tt != dateType && tt != timeType && tt != dateTimeType {
		return "", errors.Errorf("cannot compare values of type %s", cg.toTypeName(tt))
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

	return fmt.Sprintf("%s %s %s", t1, op, t2), err
}

// emitEquality renders e.
func (cg *SqlCodeGenerator) emitEquality(e *Equality) (string, error) {
	t1, tt1, err := cg.emitTermOrMath(e.TermOrMath1)
	if err != nil {
		return "", err
	}

	t2, tt2, err := cg.emitTermOrMath(e.TermOrMath2)
	if err != nil {
		return "", err
	}

	_, err = cg.validateTypes(tt1, tt2)
	if err != nil {
		return "", err
	}

	var op string
	switch e.Op {
	case "eq":
		op = "="
	case "neq":
		op = "<>"
	}

	return fmt.Sprintf("%s %s %s", t1, op, t2), err
}

// emitRange renders r.
func (cg *SqlCodeGenerator) emitRange(r *Range) (string, error) {
	t1, tt1, err := cg.emitTermOrMath(r.TermOrMath1)
	if err != nil {
		return "", err
	}

	t2, tt2, err := cg.emitTermOrMath(r.TermOrMath2)
	if err != nil {
		return "", err
	}

	t3, tt3, err := cg.emitTermOrMath(r.TermOrMath3)
	if err != nil {
		return "", err
	}

	tt, err := cg.validateTypes(tt1, tt2)
	if err != nil {
		return "", err
	}
	tt, err = cg.validateTypes(tt, tt3)
	if err != nil {
		return "", err
	} else if tt != intType && tt != decimalType {
		return "", errors.Errorf("cannot range values of type %s", cg.toTypeName(tt))
	}

	return fmt.Sprintf("%s %s %s %s %s", t1,
		strings.ToUpper(r.Between), t2,
		strings.ToUpper(r.And), t3), err
}

// emitMatch renders m.
func (cg *SqlCodeGenerator) emitMatch(m *Match) (string, error) {
	t1, tt1, err := cg.emitTerm(m.Term1)
	if err != nil {
		return "", err
	}

	t2, tt2, err := cg.emitTerm(m.Term2)
	if err != nil {
		return "", err
	}

	tt, err := cg.validateTypes(tt1, tt2)
	if err != nil {
		return "", err
	} else if tt != stringType {
		return "", errors.Errorf("cannot match values of type %s", cg.toTypeName(tt))
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

	return fmt.Sprintf("%s %s %s", t1, "LIKE", t2), err
}

// emitIs renders i.
func (cg *SqlCodeGenerator) emitIs(i *Is) (string, error) {
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

	return sb.String(), nil
}

// emitTermOrMath renders tm.
func (cg *SqlCodeGenerator) emitTermOrMath(tm *TermOrMath) (string, termType, error) {
	var err error
	var s string
	var t termType

	if tm.Math != nil {
		s, t, err = cg.emitMath(tm.Math)
	} else if tm.SubMath != nil {
		if s, t, err = cg.emitMath(tm.SubMath); err == nil {
			s = fmt.Sprintf("(%s)", s)
		}
	} else if tm.Term != nil {
		s, t, err = cg.emitTerm(tm.Term)
	}

	return s, t, err
}

// emitTerm renders t.
func (cg *SqlCodeGenerator) emitTerm(t *Term) (string, termType, error) {
	var err error
	var s string
	var tt termType

	if t.Identifier != nil {
		tt = identType
		s, err = cg.applyRenderingOptions(*t.Identifier, tt)
	} else if t.Integer != nil {
		tt = intType
		s, err = cg.applyRenderingOptions(strconv.Itoa(*t.Integer), tt)
	} else if t.Decimal != nil {
		tt = decimalType
		s, err = cg.applyRenderingOptions(strconv.FormatFloat(*t.Decimal, 'f', -1, 64), tt)
	} else if t.String != nil {
		tt = stringType
		s, err = cg.applyRenderingOptions(fmt.Sprintf("'%s'", *t.String), tt)
	} else if t.Date != nil {
		tt = dateType
		s, err = cg.applyRenderingOptions(fmt.Sprintf("'%s'", *t.Date), tt)
	} else if t.Time != nil {
		tt = timeType
		s, err = cg.applyRenderingOptions(fmt.Sprintf("'%s'", *t.Time), tt)
	} else if t.DateTime != nil {
		tt = dateTimeType
		s, err = cg.applyRenderingOptions(fmt.Sprintf("'%s'", strings.Replace(*t.DateTime, "T", " ", -1)), tt)
	} else if t.Bool != nil {
		tt = boolType
		if *t.Bool == "true" {
			s = "1"
		} else {
			s = "0"
		}
		s, err = cg.applyRenderingOptions(s, tt)
	} else if t.Macro != nil {
		s, tt, err = cg.emitMacro(t.Macro)
	}

	return s, tt, err
}

// emitMath renders m.
func (cg *SqlCodeGenerator) emitMath(m *Math) (string, termType, error) {
	t1, tt1, err := cg.emitTerm(m.Term1)
	if err != nil {
		return "", undefType, err
	}

	t2, tt2, err := cg.emitTerm(m.Term2)
	if err != nil {
		return "", undefType, err
	}

	tt, err := cg.validateTypes(tt1, tt2)
	if err != nil {
		return "", undefType, err
	} else if tt != intType && tt != decimalType && tt != dateType && tt != timeType && tt != dateTimeType {
		return "", tt, errors.Errorf("cannot compute values of type %s", cg.toTypeName(tt))
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

	return fmt.Sprintf("%s %s %s", t1, op, t2), tt, err
}

// emitMacro renders m.
func (cg *SqlCodeGenerator) emitMacro(m *Macro) (string, termType, error) {
	var err error
	var s string
	var t termType

	switch m.Name {
	case "#now":
		s, t, err = cg.emitNowMacro(m)
	case "#duration":
		s, t, err = cg.emitDurationMacro(m)
	}

	return s, t, err
}

// emitNowMacro renders m.
func (cg *SqlCodeGenerator) emitNowMacro(m *Macro) (string, termType, error) {
	return "CURRENT_TIMESTAMP", dateTimeType, nil
}

// emitDurationMacro renders m.
func (cg *SqlCodeGenerator) emitDurationMacro(m *Macro) (string, termType, error) {
	if m.Args == nil {
		return "", undefType, errors.Errorf("%s: missing parameter: iso8601 interval", m.Name)
	}

	const interval = "INTERVAL '%s'"
	var err error
	var items []string
	p := pluralize.NewClient()

	for _, a := range m.Args {
		s, t, err := cg.emitTerm(a)
		if err != nil {
			return "", undefType, err
		} else if t != stringType {
			return "", undefType, errors.Errorf("iso8601 interval cannot be of type %s", cg.toTypeName(t))
		}
		d, err := duration.FromString(s)
		if err != nil {
			return "", undefType, err
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

	return sb.String(), dateTimeType, err
}

// validateTypes verifies whether or not t1 and t2 are compatible, and if they are,
// it returns the result type of the current expression.
func (cg *SqlCodeGenerator) validateTypes(t1 termType, t2 termType) (termType, error) {
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

	return t, err
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
func (cg *SqlCodeGenerator) applyRenderingOptions(f string, t termType) (string, error) {
	if cg.RenderingOptions != nil {
		if t == identType {
			if val := cg.RenderingOptions.GetFieldProps(f); val != nil {
				if !val.Filterable {
					return "", errors.Errorf("field %v is not filterable", f)
				}
				if len(val.NativeName) > 0 {
					return val.NativeName, nil
				}
			}
		} else if cg.RenderingOptions.NamedParamsEnabled() {
			v, _ := cg.RenderingOptions.GetNamedParamValues()
			paramName := fmt.Sprintf("%s%d", cg.RenderingOptions.GetNamedParamsPrefix(), len(v)+1)
			v[paramName] = f
			return ":" + paramName, nil
		}
	}

	return f, nil
}
