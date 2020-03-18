/**
 * @begin 2020-03-18
 * @author <a href="mailto:giuseppe.greco@skeeterhealth.com">Giuseppe Greco</a>
 * @copyright 2020 <a href="skeeterhealth.com">Skeeter</a>
 */

package espressopp

import "github.com/pkg/errors"

// FieldProps is the set of properties associated with a field.
type FieldProps struct {
	Filterable bool
	NativeName string
}

// RenderingOptions is the set of options used by CodeGenerator implementations
// to control the way target code is generated.
type RenderingOptions struct {
	fields map[string]*FieldProps
}

// NewRenderingOptions creates a new instance of RenderingOptions.
func NewRenderingOptions() *RenderingOptions {
	return &RenderingOptions{
		fields: make(map[string]*FieldProps),
	}
}

// NewRenderingOptionsFromFieldNames creates a new instance of RenderingOptions
// from the specified map.
func NewRenderingOptionsFromFieldNames(fn map[string]string) *RenderingOptions {
	ro := NewRenderingOptions()

	for k, v := range fn {
		ro.AddFieldProps(k, &FieldProps{
			Filterable: true,
			NativeName: v,
		})
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
