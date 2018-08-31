package dsl

import (
	"goa.design/goa/eval"
	"goa.design/goa/expr"
)

// Meta is a set of key/value pairs that can be assigned to an object. Each
// value consists of a slice of strings so that multiple invocation of the
// Meta function on the same target using the same key builds up the slice.
// Meta may be set on attributes, result types, endpoints, responses,
// services and API definitions.
//
// While keys can have any value the following names are handled explicitly by
// goa when set on attributes or types.
//
// type:generate:force forces the code generation for the type it is defined
// on. By default goa only generates types that are used explicitly by the
// service methods. This meta key makes it possible to override this
// behavior and forces goa to generate the corresponding struct. The value is
// a slice of strings that lists the names of the services for which to
// generate the struct. If left empty then the struct is generated for all
// services.
//
//        package design
//
//        var ExternalType = Type("ExternalType", func() {
//                Attribute("name", String)
//                Meta("type:generate:force", service1, service2)
//        })
//
//        var _ = Service("service1", func() {
//                ...
//        })
//
//        var _ = Service("service2", func() {
//                ...
//        })
//
// struct:error:name identifies the attribute of a result type used to select
// the returned error when multiple errors are defined on the same method.
// The value of the field corresponding to the attribute with the
// struct:error:name meta is matched against the names of the method
// errors as defined in the design. This makes it possible to define distinct
// transport mappings for the various errors (for example to return different
// HTTP status codes). There must be one and exactly one attribute with the
// struct:error:name meta defined on result types used to define error
// results.
//
//        var CustomErrorType = ResultType("application/vnd.goa.error", func() {
//                Attribute("message", String, "Error returned.", func() {
//                        Meta("struct:error:name")
//                })
//                Attribute("occurred_at", DateTime, "Time error occurred.")
//        })
//
//        var _ = Service("MyService", func() {
//                Error("internal_error", CustomErrorType)
//                Error("bad_request", CustomErrorType)
//        })
//
// `struct:field:name`: overrides the Go struct field name generated by default
// by goa.  Applicable to attributes only.
//
//        Meta("struct:field:name", "MyName")
//
// `struct:field:origin`: overrides the name of the value used to initialize an
// attribute value. For example if the attributes describes an HTTP header this
// meta specifies the name of the header in case it's different from the name
// of the attribute. Applicable to attributes only.
//
//        Meta("struct:field:origin", "X-API-Version")
//
// `struct:tag:xxx`: sets the struct field tag xxx on generated Go structs.
// Overrides tags that goa would otherwise set.  If the meta value is a
// slice then the strings are joined with the space character as separator.
// Applicable to attributes only.
//
//        Meta("struct:tag:json", "myName,omitempty")
//        Meta("struct:tag:xml", "myName,attr")
//
// `swagger:generate`: specifies whether Swagger specification should be
// generated. Defaults to true.
// Applicable to services, methods and file servers.
//
//        Meta("swagger:generate", "false")
//
// `swagger:summary`: sets the Swagger operation summary field.
// Applicable to endpoints.
//
//        Meta("swagger:summary", "Short summary of what endpoint does")
//
// `swagger:example`: specifies whether to generate random example. Defaults to
// true.
// Applicable to API (for global setting) or individual attributes.
//
//        Meta("swagger:example", "false")
//
// `swagger:tag:xxx`: sets the Swagger object field tag xxx.
// Applicable to services and endpoints.
//
//        Meta("swagger:tag:Backend")
//        Meta("swagger:tag:Backend:desc", "Description of Backend")
//        Meta("swagger:tag:Backend:url", "http://example.com")
//        Meta("swagger:tag:Backend:url:desc", "See more docs here")
//
// `swagger:extension:xxx`: sets the Swagger extensions xxx. It can have any
// valid JSON format value.
// Applicable to:
// api as within the info and tag object,
// service within the paths object,
// endpoint as within the path-item object,
// route as within the operation object,
// param as within the parameter object,
// response as within the response object
// and security as within the security-scheme object.
// See https://github.com/OAI/OpenAPI-Specification/blob/master/guidelines/EXTENSIONS.md.
//
//        Meta("swagger:extension:x-api", `{"foo":"bar"}`)
//
// The special key names listed above may be used as follows:
//
//        var Account = Type("Account", func() {
//                Attribute("service", String, "Name of service", func() {
//                        // Override default name
//                        Meta("struct:field:name", "ServiceName")
//                })
//        })
//
func Meta(name string, value ...string) {
	appendMeta := func(meta expr.MetaExpr, name string, value ...string) expr.MetaExpr {
		if meta == nil {
			meta = make(map[string][]string)
		}
		meta[name] = append(meta[name], value...)
		return meta
	}

	switch expr := eval.Current().(type) {
	case expr.CompositeExpr:
		att := expr.Attribute()
		att.Meta = appendMeta(att.Meta, name, value...)
	case *expr.AttributeExpr:
		expr.Meta = appendMeta(expr.Meta, name, value...)
	case *expr.ResultTypeExpr:
		expr.Meta = appendMeta(expr.Meta, name, value...)
	case *expr.MethodExpr:
		expr.Meta = appendMeta(expr.Meta, name, value...)
	case *expr.ServiceExpr:
		expr.Meta = appendMeta(expr.Meta, name, value...)
	case *expr.APIExpr:
		expr.Meta = appendMeta(expr.Meta, name, value...)
	default:
		eval.IncompatibleDSL()
	}
}
