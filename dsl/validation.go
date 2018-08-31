package dsl

import (
	"reflect"
	"regexp"
	"strconv"

	"goa.design/goa/expr"
	"goa.design/goa/eval"
)

// Enum adds a "enum" validation to the attribute.
// See http://json-schema.org/latest/json-schema-validation.html#anchor76.
func Enum(vals ...interface{}) {
	if a, ok := eval.Current().(*expr.AttributeExpr); ok {
		for i, v := range vals {
			// When can a.Type be nil? glad you asked
			// There are two ways to write an Attribute declaration with the DSL that
			// don't set the type: with one argument - just the name - in which case the type
			// is set to String or with two arguments - the name and DSL. In this latter form
			// the type can end up being either String - if the DSL does not define any
			// attribute - or object if it does.
			// Why allowing this? because it's not always possible to specify the type of an
			// object - an object may just be declared inline to represent a substructure.
			// OK then why not assuming object and not allowing for string? because the DSL
			// where there's only one argument and the type is string implicitly is very
			// useful and common, for example to list attributes that refer to other attributes
			// such as responses that refer to responses defined at the API level or links that
			// refer to the result type attributes. So if the form that takes a DSL always ended
			// up defining an object we'd have a weird situation where one arg is string and
			// two args is object. Breaks the least surprise principle. Soooo long story
			// short the lesser evil seems to be to allow the ambiguity. Also tests like the
			// one below are really a convenience to the user and not a fundamental feature
			// - not checking in the case the type is not known yet is OK.
			if a.Type != nil && !a.Type.IsCompatible(v) {
				eval.ReportError("value %#v at index %d is incompatible with attribute of type %s",
					v, i, a.Type.Name())
				ok = false
			}
		}
		if ok {
			if a.Validation == nil {
				a.Validation = &expr.ValidationExpr{}
			}
			a.Validation.Values = make([]interface{}, len(vals))
			for i, v := range vals {
				switch actual := v.(type) {
				case expr.MapVal:
					a.Validation.Values[i] = actual.ToMap()
				case expr.ArrayVal:
					a.Validation.Values[i] = actual.ToSlice()
				default:
					a.Validation.Values[i] = actual
				}
			}
		}
	}
}

// Format adds a "format" validation to the attribute.
// See http://json-schema.org/latest/json-schema-validation.html#anchor104.
// The formats supported by goa are:
//
// FormatDate: RFC3339 date
//
// FormatDateTime: RFC3339 date time
//
// FormatUUID: RFC4122 uuid
//
// FormatEmail: RFC5322 email address
//
// FormatHostname: RFC1035 internet host name
//
// FormatIPv4, FormatIPv6, FormatIP: RFC2373 IPv4, IPv6 address or either
//
// FormatURI: RFC3986 URI
//
// FormatMAC: IEEE 802 MAC-48, EUI-48 or EUI-64 MAC address
//
// FormatCIDR: RFC4632 or RFC4291 CIDR notation IP address
//
// FormatRegexp: RE2 regular expression
//
// FormatJSON: JSON text
//
// FormatRFC1123: RFC1123 date time
//
func Format(f expr.ValidationFormat) {
	if a, ok := eval.Current().(*expr.AttributeExpr); ok {
		if !a.IsSupportedValidationFormat(f) {
			eval.ReportError("invalid validation format %q", f)
		}
		if a.Type != nil && a.Type.Kind() != expr.StringKind {
			incompatibleAttributeType("format", a.Type.Name(), "a string")
		} else {
			if a.Validation == nil {
				a.Validation = &expr.ValidationExpr{}
			}
			a.Validation.Format = expr.ValidationFormat(f)
		}
	}
}

// Pattern adds a "pattern" validation to the attribute.
// See http://json-schema.org/latest/json-schema-validation.html#anchor33.
func Pattern(p string) {
	if a, ok := eval.Current().(*expr.AttributeExpr); ok {
		if a.Type != nil && a.Type.Kind() != expr.StringKind {
			incompatibleAttributeType("pattern", a.Type.Name(), "a string")
		} else {
			_, err := regexp.Compile(p)
			if err != nil {
				eval.ReportError("invalid pattern %#v, %s", p, err)
			} else {
				if a.Validation == nil {
					a.Validation = &expr.ValidationExpr{}
				}
				a.Validation.Pattern = p
			}
		}
	}
}

// Minimum adds a "minimum" validation to the attribute.
// See http://json-schema.org/latest/json-schema-validation.html#anchor21.
func Minimum(val interface{}) {
	if a, ok := eval.Current().(*expr.AttributeExpr); ok {
		if a.Type != nil &&
			a.Type.Kind() != expr.IntKind && a.Type.Kind() != expr.UIntKind &&
			a.Type.Kind() != expr.Int32Kind && a.Type.Kind() != expr.UInt32Kind &&
			a.Type.Kind() != expr.Int64Kind && a.Type.Kind() != expr.UInt64Kind &&
			a.Type.Kind() != expr.Float32Kind && a.Type.Kind() != expr.Float64Kind {

			incompatibleAttributeType("minimum", a.Type.Name(), "an integer or a number")
		} else {
			var f float64
			switch v := val.(type) {
			case float32, float64, int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
				f = reflect.ValueOf(v).Convert(reflect.TypeOf(float64(0.0))).Float()
			case string:
				var err error
				f, err = strconv.ParseFloat(v, 64)
				if err != nil {
					eval.ReportError("invalid number value %#v", v)
					return
				}
			default:
				eval.ReportError("invalid number value %#v", v)
				return
			}
			if a.Validation == nil {
				a.Validation = &expr.ValidationExpr{}
			}
			a.Validation.Minimum = &f
		}
	}
}

// Maximum adds a "maximum" validation to the attribute.
// See http://json-schema.org/latest/json-schema-validation.html#anchor17.
func Maximum(val interface{}) {
	if a, ok := eval.Current().(*expr.AttributeExpr); ok {
		if a.Type != nil &&
			a.Type.Kind() != expr.IntKind && a.Type.Kind() != expr.UIntKind &&
			a.Type.Kind() != expr.Int32Kind && a.Type.Kind() != expr.UInt32Kind &&
			a.Type.Kind() != expr.Int64Kind && a.Type.Kind() != expr.UInt64Kind &&
			a.Type.Kind() != expr.Float32Kind && a.Type.Kind() != expr.Float64Kind {

			incompatibleAttributeType("maximum", a.Type.Name(), "an integer or a number")
		} else {
			var f float64
			switch v := val.(type) {
			case float32, float64, int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
				f = reflect.ValueOf(v).Convert(reflect.TypeOf(float64(0.0))).Float()
			case string:
				var err error
				f, err = strconv.ParseFloat(v, 64)
				if err != nil {
					eval.ReportError("invalid number value %#v", v)
					return
				}
			default:
				eval.ReportError("invalid number value %#v", v)
				return
			}
			if a.Validation == nil {
				a.Validation = &expr.ValidationExpr{}
			}
			a.Validation.Maximum = &f
		}
	}
}

// MinLength adds a "minItems" validation to the attribute.
// See http://json-schema.org/latest/json-schema-validation.html#anchor45.
func MinLength(val int) {
	if a, ok := eval.Current().(*expr.AttributeExpr); ok {
		if a.Type != nil {
			kind := a.Type.Kind()
			if kind != expr.BytesKind &&
				kind != expr.StringKind &&
				kind != expr.ArrayKind &&
				kind != expr.MapKind {

				incompatibleAttributeType("minimum length", a.Type.Name(), "a string or an array")
				return
			}
		}
		if a.Validation == nil {
			a.Validation = &expr.ValidationExpr{}
		}
		a.Validation.MinLength = &val
	}
}

// MaxLength adds a "maxItems" validation to the attribute.
// See http://json-schema.org/latest/json-schema-validation.html#anchor42.
func MaxLength(val int) {
	if a, ok := eval.Current().(*expr.AttributeExpr); ok {
		if a.Type != nil {
			kind := a.Type.Kind()
			if kind != expr.BytesKind &&
				kind != expr.StringKind &&
				kind != expr.ArrayKind &&
				kind != expr.MapKind {

				incompatibleAttributeType("maximum length", a.Type.Name(), "a string or an array")
				return
			}
		}
		if a.Validation == nil {
			a.Validation = &expr.ValidationExpr{}
		}
		a.Validation.MaxLength = &val
	}
}

// Required adds a "required" validation to the attribute.
// See http://json-schema.org/latest/json-schema-validation.html#anchor61.
func Required(names ...string) {
	var at *expr.AttributeExpr

	switch def := eval.Current().(type) {
	case *expr.AttributeExpr:
		at = def
	case *expr.ResultTypeExpr:
		at = def.AttributeExpr
	default:
		eval.IncompatibleDSL()
		return
	}

	if at.Type != nil && at.Type.Kind() != expr.ObjectKind {
		incompatibleAttributeType("required", at.Type.Name(), "an object")
	} else {
		if at.Validation == nil {
			at.Validation = &expr.ValidationExpr{}
		}
		at.Validation.AddRequired(names...)
	}
}

// incompatibleAttributeType reports an error for validations defined on
// incompatible attributes (e.g. max value on string).
func incompatibleAttributeType(validation, actual, expected string) {
	eval.ReportError("invalid %s validation definition: attribute must be %s (but type is %s)",
		validation, expected, actual)
}
