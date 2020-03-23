/**
 * @begin 2020-03-18
 * @author <a href="mailto:giuseppe.greco@skeeterhealth.com">Giuseppe Greco</a>
 * @copyright 2020 <a href="skeeterhealth.com">Skeeter</a>
 */

package espressopp

import "github.com/pkg/errors"

// FieldProps is the set of properties associated with a field.
type FieldProps struct {
	// Filterable specifies whether or not the field can be used in a query. If
	// not and an expression actually contains it, the underlying CodeGenerator
	// implementation shall raise an error.
	Filterable bool

	// NativeName is used to map those fields in an input expression that do
	// not match the field names of the underlying database.
	NativeName string
}

// namedParams lets code generators render named parameters and set aside their
// values for client code.
type namedParams struct {
	// enabled specifies whether or not named parameters are enabled.
	enabled bool

	// prefix specifies the string that is prepended to parameter names.
	prefix string

	// values contains the values of the named parameters present in rendered
	// code.
	values map[string]string
}

// RenderingOptions is the set of options used by CodeGenerator implementations
// to control the way target code is generated.
type RenderingOptions struct {
	fields      map[string]*FieldProps
	namedParams *namedParams
}

const (
	defaultPrefix = "P"
)

// NewRenderingOptions creates a new instance of RenderingOptions.
func NewRenderingOptions() *RenderingOptions {
	return &RenderingOptions{
		fields: make(map[string]*FieldProps),
		namedParams: &namedParams{
			prefix: defaultPrefix,
		},
	}
}

// Clone performs a shallow copy of read-only data and a deep copy of
// read-write data. Read-only data includes field properties whereas
// read-write data includes named parameters.
func (ro *RenderingOptions) Clone() *RenderingOptions {
	var m map[string]string

	if ro.namedParams.values != nil {
		m = make(map[string]string)
		for k, v := range ro.namedParams.values {
			m[k] = v
		}
	}

	return &RenderingOptions{
		fields: ro.fields,
		namedParams: &namedParams{
			enabled: ro.namedParams.enabled,
			prefix:  ro.namedParams.prefix,
			values:  m,
		},
	}
}

// Fields initializes the fields rendering options with the specified map of
// fieldName:fieldProps items. Previous fields rendering options are lost. If
// m is nil then all fields rendering options are removed.
func (ro *RenderingOptions) Fields(m map[string]*FieldProps) *RenderingOptions {
	for k := range ro.fields {
		delete(ro.fields, k)
	}

	if m != nil {
		for k, v := range m {
			if len(v.NativeName) == 0 {
				v.NativeName = k
			}
			ro.fields[k] = &FieldProps{
				Filterable: v.Filterable,
				NativeName: v.NativeName,
			}
		}
	}

	return ro
}

// FieldsWithDefault initializes the fields rendering options with the specified
// map of fieldName:nativeFieldName items. Previous fields rendering options are
// lost and Filterable is default to true for each field. If m is nil then all
// fields rendering options are removed.
func (ro *RenderingOptions) FieldsWithDefault(m map[string]string) *RenderingOptions {
	for k := range ro.fields {
		delete(ro.fields, k)
	}

	if m != nil {
		for k, v := range m {
			if len(v) == 0 {
				v = k
			}
			ro.fields[k] = &FieldProps{
				Filterable: true,
				NativeName: v,
			}
		}
	}

	return ro
}

// AddFieldProps adds the specified field properties to the rendering options.
func (ro *RenderingOptions) AddFieldProps(fieldName string, fp *FieldProps) error {
	if len(fieldName) == 0 {
		return errors.New("field name not specified")
	}

	if fp == nil {
		return errors.Errorf("properties for field %v not specified", fieldName)
	}

	if len(fp.NativeName) == 0 {
		fp.NativeName = fieldName
	}

	ro.fields[fieldName] = fp
	return nil
}

// RemoveFieldProps removes the properties of the specified field from the
// rendering options.
func (ro *RenderingOptions) RemoveFieldProps(fieldName string) *FieldProps {
	fp := ro.fields[fieldName]
	if fp != nil {
		delete(ro.fields, fieldName)
	}

	return fp
}

// GetFieldProps retrieves the properties of the specified field from the
// rendering options.
func (ro *RenderingOptions) GetFieldProps(fieldName string) *FieldProps {
	return ro.fields[fieldName]
}

// EnableNamedParams enables named parameters in rendered code.
func (ro *RenderingOptions) EnableNamedParams() {
	if !ro.namedParams.enabled {
		ro.namedParams.enabled = true
		ro.namedParams.values = make(map[string]string)
	}
}

// DisableNamedParams disables named parameters in rendered code.
func (ro *RenderingOptions) DisableNamedParams() {
	ro.namedParams.enabled = false
	ro.namedParams.values = nil
}

// NamedParamsEnabled returns a Boolean value indicating whether or not
// named parameters are enabled.
func (ro *RenderingOptions) NamedParamsEnabled() bool {
	return ro.namedParams.enabled
}

// GetNamedParamValues returns a map containing the values of the named
// parameters present in rendered code.
func (ro *RenderingOptions) GetNamedParamValues() (error, map[string]string) {
	if !ro.namedParams.enabled {
		return errors.New("named parameters not enabled"), nil
	}

	return nil, ro.namedParams.values
}

// SetNamedParamsPrefix sets the string that is prepended to parameter names.
func (ro *RenderingOptions) SetNamedParamsPrefix(p string) {
	prefix := p
	if prefix == "" {
		prefix = defaultPrefix
	}

	ro.namedParams.prefix = prefix
}

// getNamedParamsPrefix gets the string that is prepended to parameter names.
func (ro *RenderingOptions) GetNamedParamsPrefix() string {
	return ro.namedParams.prefix
}
