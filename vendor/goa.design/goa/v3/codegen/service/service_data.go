package service

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/template"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/expr"
)

// Services holds the data computed from the design needed to generate the code
// of the services.
var Services = make(ServicesData)

var (
	// initTypeTmpl is the template used to render the code that initializes a
	// projected type or viewed result type or a result type.
	initTypeCodeTmpl = template.Must(
		template.New("initTypeCode").
			Funcs(template.FuncMap{"goify": codegen.Goify}).
			Parse(readTemplate("return_type_init")),
	)

	// validateTypeCodeTmpl is the template used to render the code to
	// validate a projected type or a viewed result type.
	validateTypeCodeTmpl = template.Must(
		template.New("validateType").
			Funcs(template.FuncMap{"goify": codegen.Goify}).
			Parse(readTemplate("type_validate")),
	)
)

type (
	// ServicesData encapsulates the data computed from the service designs.
	ServicesData map[string]*Data

	// Data contains the data used to render the code related to a single
	// service.
	Data struct {
		// Name is the service name.
		Name string
		// Description is the service description.
		Description string
		// APIName is the name of the API the service belongs to.
		APIName string
		// APIVersion is the API version.
		APIVersion string
		// StructName is the service struct name.
		StructName string
		// VarName is the service variable name (first letter in lowercase).
		VarName string
		// PathName is the service name as used in file and import paths.
		PathName string
		// PkgName is the name of the package containing the generated service
		// code.
		PkgName string
		// ViewsPkg is the name of the package containing the projected and viewed
		// result types.
		ViewsPkg string
		// Methods lists the service interface methods.
		Methods []*MethodData
		// Schemes is the list of security schemes required by the service methods.
		Schemes SchemesData
		// ServerInterceptors contains the data needed to render the server-side
		// interceptors code.
		ServerInterceptors []*InterceptorData
		// ClientInterceptors contains the data needed to render the client-side
		// interceptors code.
		ClientInterceptors []*InterceptorData
		// Scope initialized with all the service types.
		Scope *codegen.NameScope
		// ViewScope initialized with all the viewed types.
		ViewScope *codegen.NameScope
		// UserTypeImports lists the import specifications for the user
		// types used by the service.
		UserTypeImports []*codegen.ImportSpec
		// ProtoImports lists the import specifications for the custom
		// proto types used by the service.
		ProtoImports []*codegen.ImportSpec

		// userTypes lists the type definitions that the service depends on.
		userTypes []*UserTypeData
		// errorTypes lists the error type definitions that the service depends on.
		errorTypes []*UserTypeData
		// errorInits list the information required to generate error init
		// functions.
		errorInits []*ErrorInitData
		// projectedTypes lists the types which uses pointers for all fields to
		// define view specific validation logic.
		projectedTypes []*ProjectedTypeData
		// union methods that need to be defined in views package.
		viewedUnionMethods []*UnionValueMethodData
		// viewedResultTypes lists all the viewed method result types.
		viewedResultTypes []*ViewedResultTypeData
		// unionValueMethods lists the methods used to define union types.
		unionValueMethods []*UnionValueMethodData
	}

	// MethodData describes a single service method.
	MethodData struct {
		// Name is the method name.
		Name string
		// Description is the method description.
		Description string
		// VarName is the Go method name.
		VarName string
		// Payload is the name of the payload type if any,
		Payload string
		// PayloadLoc defines the file and Go package of the payload type
		// if overridden via Meta.
		PayloadLoc *codegen.Location
		// PayloadDef is the payload type definition if any.
		PayloadDef string
		// PayloadRef is a reference to the payload type if any,
		PayloadRef string
		// PayloadDesc is the payload type description if any.
		PayloadDesc string
		// PayloadEx is an example of a valid payload value.
		PayloadEx any
		// PayloadDefault is the default value of the payload if any.
		PayloadDefault any
		// StreamingPayload is the name of the streaming payload type if any.
		StreamingPayload string
		// StreamingPayloadDef is the streaming payload type definition if any.
		StreamingPayloadDef string
		// StreamingPayloadRef is a reference to the streaming payload type if any.
		StreamingPayloadRef string
		// StreamingPayloadDesc is the streaming payload type description if any.
		StreamingPayloadDesc string
		// StreamingPayloadEx is an example of a valid streaming payload value.
		StreamingPayloadEx any
		// Result is the name of the result type if any.
		Result string
		// ResultLoc defines the file and Go package of the result type
		// if overridden via Meta.
		ResultLoc *codegen.Location
		// ResultDef is the result type definition if any.
		ResultDef string
		// ResultRef is the reference to the result type if any.
		ResultRef string
		// ResultDesc is the result type description if any.
		ResultDesc string
		// ResultEx is an example of a valid result value.
		ResultEx any
		// Errors list the possible errors defined in the design if any.
		Errors []*ErrorInitData
		// ErrorLocs lists the file and Go package of the error type
		// if overridden via Meta indexed by error name.
		ErrorLocs map[string]*codegen.Location
		// Requirements contains the security requirements for the
		// method.
		Requirements RequirementsData
		// Schemes contains the security schemes types used by the
		// method.
		Schemes SchemesData
		// ServerInterceptors list the server interceptors that apply to this
		// method.
		ServerInterceptors []string
		// ClientInterceptors list the client interceptors that apply to this
		// method.
		ClientInterceptors []string
		// ViewedResult contains the data required to generate the code handling
		// views if any.
		ViewedResult *ViewedResultTypeData
		// ServerStream indicates that the service method receives a payload
		// stream or sends a result stream or both.
		ServerStream *StreamData
		// ClientStream indicates that the service method receives a result
		// stream or sends a payload result or both.
		ClientStream *StreamData
		// StreamKind is the kind of the stream (payload or result or
		// bidirectional).
		StreamKind expr.StreamKind
		// SkipRequestBodyEncodeDecode is true if the method payload includes
		// the raw HTTP request body reader.
		SkipRequestBodyEncodeDecode bool
		// SkipResponseBodyEncodeDecode is true if the method result includes
		// the raw HTTP response body reader.
		SkipResponseBodyEncodeDecode bool
		// RequestStruct is the name of the data structure containing the
		// payload and request body reader when SkipRequestBodyEncodeDecode is
		// used.
		RequestStruct string
		// ResponseStruct is the name of the data structure containing the
		// result and response body reader when SkipResponseBodyEncodeDecode is
		// used.
		ResponseStruct string
	}

	// StreamData is the data used to generate client and server interfaces that
	// a streaming endpoint implements. It is initialized if a method defines a
	// streaming payload or result or both.
	StreamData struct {
		// Interface is the name of the stream interface.
		Interface string
		// VarName is the name of the struct type that implements the stream
		// interface.
		VarName string
		// SendName is the name of the send function.
		SendName string
		// SendDesc is the description for the send function.
		SendDesc string
		// SendWithContextName is the name of the send function with context.
		SendWithContextName string
		// SendWithContextDesc is the description for the send function with context.
		SendWithContextDesc string
		// SendTypeName is the type name sent through the stream.
		SendTypeName string
		// SendTypeRef is the reference to the type sent through the stream.
		SendTypeRef string
		// RecvName is the name of the receive function.
		RecvName string
		// RecvDesc is the description for the recv function.
		RecvDesc string
		// RecvWithContextName is the name of the receive function with context.
		RecvWithContextName string
		// RecvWithContextDesc is the description for the recv function with context.
		RecvWithContextDesc string
		// RecvTypeName is the type name received from the stream.
		RecvTypeName string
		// RecvTypeRef is the reference to the type received from the stream.
		RecvTypeRef string
		// MustClose indicates whether the stream should implement the Close()
		// function.
		MustClose bool
		// EndpointStruct is the name of the endpoint struct that holds a payload
		// reference (if any) and the endpoint server stream.
		EndpointStruct string
		// Kind is the kind of the stream (payload, result or bidirectional).
		Kind expr.StreamKind
	}

	// ErrorInitData describes an error returned by a service method of type
	// ErrorResult.
	ErrorInitData struct {
		// Name is the name of the init function.
		Name string
		// Description is the error description.
		Description string
		// ErrName is the name of the error.
		ErrName string
		// TypeName is the error struct type name.
		TypeName string
		// TypeRef is the reference to the error type.
		TypeRef string
		// Temporary indicates whether the error is temporary.
		Temporary bool
		// Timeout indicates whether the error is due to timeouts.
		Timeout bool
		// Fault indicates whether the error is server-side fault.
		Fault bool
	}

	// InterceptorData contains the data required to render the service-level
	// interceptor code. interceptors.go.tpl
	InterceptorData struct {
		// Name is the name of the interceptor used in the generated code.
		Name string
		// DesignName is the name of the interceptor as defined in the design.
		DesignName string
		// Description is the description of the interceptor from the design.
		Description string
		// Methods
		Methods []*MethodInterceptorData
		// ReadPayload contains payload attributes that the interceptor can
		// read.
		ReadPayload []*AttributeData
		// WritePayload contains payload attributes that the interceptor can
		// write.
		WritePayload []*AttributeData
		// ReadResult contains result attributes that the interceptor can read.
		ReadResult []*AttributeData
		// WriteResult contains result attributes that the interceptor can
		// write.
		WriteResult []*AttributeData
		// ReadStreamingPayload contains streaming payload attributes that the interceptor can read.
		ReadStreamingPayload []*AttributeData
		// WriteStreamingPayload contains streaming payload attributes that the interceptor can write.
		WriteStreamingPayload []*AttributeData
		// ReadStreamingResult contains streaming result attributes that the interceptor can read.
		ReadStreamingResult []*AttributeData
		// WriteStreamingResult contains streaming result attributes that the interceptor can write.
		WriteStreamingResult []*AttributeData
		// HasPayloadAccess indicates that the interceptor info object has a
		// payload access interface.
		HasPayloadAccess bool
		// HasResultAccess indicates that the interceptor info object has a
		// result access interface.
		HasResultAccess bool
		// HasStreamingPayloadAccess indicates that the interceptor info object has a
		// streaming payload access interface.
		HasStreamingPayloadAccess bool
		// HasStreamingResultAccess indicates that the interceptor info object has a
		// streaming result access interface.
		HasStreamingResultAccess bool
	}

	// MethodInterceptorData contains the data required to render the
	// method-level interceptor code.
	MethodInterceptorData struct {
		// MethodName is the name of the method.
		MethodName string
		// PayloadAccess is the name of the payload access struct.
		PayloadAccess string
		// ResultAccess is the name of the result access struct.
		ResultAccess string
		// StreamingPayloadAccess is the name of the streaming payload access struct.
		StreamingPayloadAccess string
		// StreamingResultAccess is the name of the streaming result access struct.
		StreamingResultAccess string
		// PayloadRef is the reference to the method payload type.
		PayloadRef string
		// ResultRef is the reference to the method result type.
		ResultRef string
		// StreamingPayloadRef is the reference to the streaming payload type.
		StreamingPayloadRef string
		// StreamingResultRef is the reference to the streaming result type.
		StreamingResultRef string
		// ServerStream is the stream data if the endpoint defines a server stream.
		ServerStream *StreamInterceptorData
		// ClientStream is the stream data if the endpoint defines a client stream.
		ClientStream *StreamInterceptorData
	}

	// StreamInterceptorData is the stream data for an interceptor.
	StreamInterceptorData struct {
		// Interface is the name of the stream interface.
		Interface string
		// SendName is the name of the send function.
		SendName string
		// SendWithContextName is the name of the send function with context.
		SendWithContextName string
		// SendTypeRef is the reference to the type sent through the stream.
		SendTypeRef string
		// RecvName is the name of the recv function.
		RecvName string
		// RecvWithContextName is the name of the recv function with context.
		RecvWithContextName string
		// RecvTypeRef is the reference to the type received from the stream.
		RecvTypeRef string
		// MustClose indicates whether the stream should implement the Close()
		// function.
		MustClose bool
		// EndpointStruct is the name of the endpoint struct that holds a payload
		// reference (if any) and the endpoint server stream.
		EndpointStruct string
	}

	// AttributeData describes a single attribute.
	AttributeData struct {
		// Name is the name of the attribute.
		Name string
		// TypeRef is the reference to the attribute type.
		TypeRef string
		// Pointer is true if the attribute is a pointer.
		Pointer bool
	}

	// RequirementsData is the list of security requirements.
	RequirementsData []*RequirementData

	// SchemesData is the list of security schemes.
	SchemesData []*SchemeData

	// RequirementData lists the schemes and scopes defined by a single
	// security requirement.
	RequirementData struct {
		// Schemes list the requirement schemes.
		Schemes []*SchemeData
		// Scopes list the required scopes.
		Scopes []string
	}

	// UserTypeData contains the data describing a user-defined type.
	UserTypeData struct {
		// Name is the type name.
		Name string
		// VarName is the corresponding Go type name.
		VarName string
		// Description is the type human description.
		Description string
		// Def is the type definition Go code.
		Def string
		// Ref is the reference to the type.
		Ref string
		// Loc defines the file and Go package of the type if overridden
		// via Meta.
		Loc *codegen.Location
		// Type is the underlying type.
		Type expr.UserType
	}

	// UnionValueMethodData describes a method used on a union value type.
	UnionValueMethodData struct {
		// Name is the name of the function.
		Name string
		// TypeRef is a reference on the target union value type.
		TypeRef string
		// Loc defines the file and Go package of the method if
		// overridden in corresponding union type via Meta.
		Loc *codegen.Location
	}

	// SchemeData describes a single security scheme.
	SchemeData struct {
		// Kind is the type of scheme, one of "Basic", "APIKey", "JWT"
		// or "OAuth2".
		Type string
		// SchemeName is the name of the scheme.
		SchemeName string
		// Name refers to a header or parameter name, based on In's
		// value.
		Name string
		// UsernameField is the name of the payload field that should be
		// initialized with the basic auth username if any.
		UsernameField string
		// UsernamePointer is true if the username field is a pointer.
		UsernamePointer bool
		// UsernameAttr is the name of the attribute that contains the
		// username.
		UsernameAttr string
		// UsernameRequired specifies whether the attribute that
		// contains the username is required.
		UsernameRequired bool
		// PasswordField is the name of the payload field that should be
		// initialized with the basic auth password if any.
		PasswordField string
		// PasswordPointer is true if the password field is a pointer.
		PasswordPointer bool
		// PasswordAttr is the name of the attribute that contains the
		// password.
		PasswordAttr string
		// PasswordRequired specifies whether the attribute that
		// contains the password is required.
		PasswordRequired bool
		// CredField contains the name of the payload field that should
		// be initialized with the API key, the JWT token or the OAuth2
		// access token.
		CredField string
		// CredPointer is true if the credential field is a pointer.
		CredPointer bool
		// CredRequired specifies if the key is a required attribute.
		CredRequired bool
		// KeyAttr is the name of the attribute that contains
		// the security tag (for APIKey, OAuth2, and JWT schemes).
		KeyAttr string
		// Scopes lists the scopes that apply to the scheme.
		Scopes []string
		// Flows describes the OAuth2 flows.
		Flows []*expr.FlowExpr
		// In indicates the request element that holds the credential.
		In string
	}

	// ViewedResultTypeData contains the data used to generate a viewed result type
	// (i.e. a method result type with more than one view). The viewed result
	// type holds the projected type and a view based on which it creates the
	// projected type. It also contains the code to validate the viewed result
	// type and the functions to initialize a viewed result type from a result
	// type and vice versa.
	ViewedResultTypeData struct {
		// the viewed result type
		*UserTypeData
		// Views lists the views defined on the viewed result type.
		Views []*ViewData
		// Validate is the validation run on the viewed result type.
		Validate *ValidateData
		// Init is the constructor code to initialize a viewed result type from
		// a result type.
		Init *InitData
		// ResultInit is the constructor code to initialize a result type
		// from the viewed result type.
		ResultInit *InitData
		// FullName is the fully qualified name of the viewed result type.
		FullName string
		// FullRef is the complete reference to the viewed result type
		// (including views package name).
		FullRef string
		// IsCollection indicates whether the viewed result type is a collection.
		IsCollection bool
		// ViewName is the view name to use to render the result type. It is set
		// only if the result type has at most one view.
		ViewName string
		// ViewsPkg is the views package name.
		ViewsPkg string
	}

	// ViewData contains data about a result type view.
	ViewData struct {
		// Name is the view name.
		Name string
		// Description is the view description.
		Description string
		// Attributes is the list of attributes rendered in the view.
		Attributes []string
		// TypeVarName is the Go variable name of the type that defines the view.
		TypeVarName string
	}

	// ProjectedTypeData contains the data used to generate a projected type for
	// the corresponding user type or result type in the service package. The
	// generated type uses pointers for all fields. It also contains the data
	// to generate view-based validation logic and transformation functions to
	// convert a projected type to its corresponding service type and vice versa.
	ProjectedTypeData struct {
		// the projected type
		*UserTypeData
		// Validations lists the validation functions to run on the projected type.
		// If the projected type corresponds to a result type then a validation
		// function for each view is generated. For user types, only one validation
		// function is generated.
		Validations []*ValidateData
		// Projections contains the code to create a projected type based on
		// views. If the projected type corresponds to a result type, then a
		// function for each view is generated.
		Projections []*InitData
		// TypeInits contains the code to convert a projected type to its
		// corresponding service type. If the projected type corresponds to a
		// result type, then a function for each view is generated.
		TypeInits []*InitData
		// ViewsPkg is the views package name.
		ViewsPkg string
		// Views lists the views defined on the projected type.
		Views []*ViewData
	}

	// InitData contains the data to render a constructor to initialize service
	// types from viewed result types and vice versa.
	InitData struct {
		// Name is the name of the constructor function.
		Name string
		// Description is the function description.
		Description string
		// Args lists arguments to this function.
		Args []*InitArgData
		// ReturnTypeRef is the reference to the return type.
		ReturnTypeRef string
		// Code is the transformation code.
		Code string
		// Helpers contain the helpers used in the transformation code.
		Helpers []*codegen.TransformFunctionData
	}

	// InitArgData represents a single constructor argument.
	InitArgData struct {
		// Name is the argument name.
		Name string
		// Ref is the reference to the argument type.
		Ref string
	}

	// ValidateData contains data to render a validate function to validate a
	// projected type or a viewed result type based on views.
	ValidateData struct {
		// Name is the validation function name.
		Name string
		// Ref is the reference to the type on which the validation function
		// is defined.
		Ref string
		// Description is the description for the validation function.
		Description string
		// Validate is the validation code.
		Validate string
	}
)

// Get retrieves the data for the service with the given name computing it if
// needed. It returns nil if there is no service with the given name.
func (d ServicesData) Get(name string) *Data {
	if data, ok := d[name]; ok {
		return data
	}
	service := expr.Root.Service(name)
	if service == nil {
		return nil
	}
	d[name] = d.analyze(service)
	return d[name]
}

// Method returns the service method data for the method with the given name,
// nil if there isn't one.
func (d *Data) Method(name string) *MethodData {
	for _, m := range d.Methods {
		if m.Name == name {
			return m
		}
	}
	return nil
}

// initUserTypeImports sets the import paths for the user types defined in the
// service.  User types may be declared in multiple packages when defined with
// the Meta key "struct:pkg:path".
func (d *Data) initUserTypeImports(genpkg string) {
	importsByPath := make(map[string]*codegen.ImportSpec)

	initLoc := func(loc *codegen.Location) {
		if loc == nil {
			return
		}
		importsByPath[loc.FilePath] = &codegen.ImportSpec{Name: loc.PackageName(), Path: genpkg + "/" + loc.RelImportPath}
	}

	for _, m := range d.Methods {
		initLoc(m.PayloadLoc)
		initLoc(m.ResultLoc)
		for _, l := range m.ErrorLocs {
			initLoc(l)
		}
		for _, ut := range d.userTypes {
			initLoc(ut.Loc)
		}
		for _, et := range d.errorTypes {
			initLoc(et.Loc)
		}
	}

	imports := make([]*codegen.ImportSpec, len(importsByPath))
	i := 0
	for _, imp := range importsByPath { // Order does not matter, imports are sorted during formatting.
		imports[i] = imp
		i++
	}
	d.UserTypeImports = imports
}

// Scheme returns the scheme data with the given scheme name.
func (r RequirementsData) Scheme(name string) *SchemeData {
	for _, req := range r {
		for _, s := range req.Schemes {
			if s.SchemeName == name {
				return s
			}
		}
	}
	return nil
}

// Dup creates a copy of the scheme data.
func (s *SchemeData) Dup() *SchemeData {
	return &SchemeData{
		Type:             s.Type,
		SchemeName:       s.SchemeName,
		Name:             s.Name,
		UsernameField:    s.UsernameField,
		UsernamePointer:  s.UsernamePointer,
		UsernameAttr:     s.UsernameAttr,
		UsernameRequired: s.UsernameRequired,
		PasswordField:    s.PasswordField,
		PasswordPointer:  s.PasswordPointer,
		PasswordAttr:     s.PasswordAttr,
		PasswordRequired: s.PasswordRequired,
		CredField:        s.CredField,
		CredPointer:      s.CredPointer,
		CredRequired:     s.CredRequired,
		KeyAttr:          s.KeyAttr,
		Scopes:           s.Scopes,
		Flows:            s.Flows,
		In:               s.In,
	}
}

// Append appends a scheme data to schemes only if it doesn't exist.
func (s SchemesData) Append(d *SchemeData) SchemesData {
	found := false
	for _, se := range s {
		if se.SchemeName == d.SchemeName {
			found = true
			break
		}
	}
	if found {
		return s
	}
	return append(s, d)
}

// DedupeByType returns a new SchemesData slice that is deduplicated by scheme
// type.
func (s SchemesData) DedupeByType() SchemesData {
	seen := make(map[string]struct{})
	uniqueSchemes := SchemesData{}
	for _, s := range s {
		if _, ok := seen[s.Type]; !ok {
			seen[s.Type] = struct{}{}
			uniqueSchemes = append(uniqueSchemes, s)
		}
	}

	return uniqueSchemes
}

// analyze creates the data necessary to render the code of the given service.
// It records the user types needed by the service definition in userTypes.
func (d ServicesData) analyze(service *expr.ServiceExpr) *Data {
	var (
		types            []*UserTypeData
		errTypes         []*UserTypeData
		errorInits       []*ErrorInitData
		projTypes        []*ProjectedTypeData
		viewedUnionMeths []*UnionValueMethodData
		viewedRTs        []*ViewedResultTypeData
	)
	scope := codegen.NewNameScope()
	scope.Unique("Use") // Reserve "Use" for Endpoints struct Use method.
	viewScope := codegen.NewNameScope()
	pkgName := scope.HashedUnique(service, strings.ToLower(codegen.Goify(service.Name, false)), "svc")
	viewspkg := pkgName + "views"
	seen := make(map[string]struct{})
	seenErrors := make(map[string]struct{})
	seenProj := make(map[string]*ProjectedTypeData)
	seenViewed := make(map[string]*ViewedResultTypeData)

	// A function to collect user types from an error expression
	recordError := func(er *expr.ErrorExpr) {
		errTypes = append(errTypes, collectTypes(er.AttributeExpr, scope, seen)...)
		if er.Type == expr.ErrorResult {
			if _, ok := seenErrors[er.Name]; ok {
				return
			}
			seenErrors[er.Name] = struct{}{}
			errorInits = append(errorInits, buildErrorInitData(er, scope))
		}
	}
	for _, er := range service.Errors {
		recordError(er)
	}

	// A function to collect inner user types from an attribute expression
	collectUserTypes := func(att *expr.AttributeExpr) {
		if ut, ok := att.Type.(expr.UserType); ok {
			att = ut.Attribute()
		}
		types = append(types, collectTypes(att, scope, seen)...)
	}
	for _, m := range service.Methods {
		// collect inner user types
		collectUserTypes(m.Payload)
		collectUserTypes(m.StreamingPayload)
		collectUserTypes(m.Result)
		// Collect projected types
		if hasResultType(m.Result) {
			types, umeths := collectProjectedTypes(expr.DupAtt(m.Result), m.Result, viewspkg, scope, viewScope, seenProj)
			projTypes = append(projTypes, types...)
			viewedUnionMeths = append(viewedUnionMeths, umeths...)
		}
		for _, er := range m.Errors {
			recordError(er)
		}
	}

	// A function to convert raw object type to user type.
	wrapObject := func(att *expr.AttributeExpr, name, id string) {
		if _, ok := att.Type.(*expr.Object); ok {
			att.Type = &expr.UserTypeExpr{
				AttributeExpr: expr.DupAtt(att),
				TypeName:      scope.Name(name),
				UID:           id,
			}
		}
		if ut, ok := att.Type.(expr.UserType); ok {
			seen[ut.ID()] = struct{}{}
		}
	}

	for _, m := range service.Methods {
		name := codegen.Goify(m.Name, true)
		// Create user type for raw object payloads
		wrapObject(m.Payload, name+"Payload", service.Name+"#"+name+"Payload")
		// Create user type for raw object streaming payloads
		wrapObject(m.StreamingPayload, name+"StreamingPayload", service.Name+"#"+name+"StreamingPayload")
		// Create user type for raw object results
		wrapObject(m.Result, name+"Result", service.Name+"#"+name+"Result")
	}

	// Add forced types
	for _, t := range expr.Root.Types {
		svcs, ok := t.Attribute().Meta["type:generate:force"]
		if !ok {
			continue
		}
		att := &expr.AttributeExpr{Type: t}
		if len(svcs) > 0 {
			// Force generate type only in the specified services
			for _, svc := range svcs {
				if svc == service.Name {
					types = append(types, collectTypes(att, scope, seen)...)
					break
				}
			}
			continue
		}
		// Force generate type in all the services
		types = append(types, collectTypes(att, scope, seen)...)
	}

	var (
		methods []*MethodData
		schemes SchemesData
	)
	methods = make([]*MethodData, len(service.Methods))
	for i, e := range service.Methods {
		m := buildMethodData(e, scope)
		methods[i] = m
		for _, s := range m.Schemes {
			schemes = schemes.Append(s)
		}
		rt, ok := e.Result.Type.(*expr.ResultTypeExpr)
		if !ok {
			continue
		}
		var view string
		if v, ok := e.Result.Meta.Last(expr.ViewMetaKey); ok {
			view = v
		}
		if vrt, ok := seenViewed[m.Result+"::"+view]; ok {
			m.ViewedResult = vrt
			continue
		}
		projected := seenProj[rt.ID()]
		projAtt := &expr.AttributeExpr{Type: projected.Type}
		vrt := buildViewedResultType(e.Result, projAtt, viewspkg, scope, viewScope)
		found := false
		for _, rt := range viewedRTs {
			if rt.Type.ID() == vrt.Type.ID() {
				found = true
				break
			}
		}
		if !found {
			viewedRTs = append(viewedRTs, vrt)
		}
		m.ViewedResult = vrt
		seenViewed[vrt.Name+"::"+view] = vrt
	}

	var unionMethods []*UnionValueMethodData
	var ms []*UnionValueMethodData
	seen = make(map[string]struct{})
	for _, t := range types {
		ms = append(ms, collectUnionMethods(&expr.AttributeExpr{Type: t.Type}, scope, t.Loc, seen)...)
	}
	for _, t := range errTypes {
		ms = append(ms, collectUnionMethods(&expr.AttributeExpr{Type: t.Type}, scope, t.Loc, seen)...)
	}
	for _, m := range service.Methods {
		ms = append(ms, collectUnionMethods(m.Payload, scope, codegen.UserTypeLocation(m.Payload.Type), seen)...)
		ms = append(ms, collectUnionMethods(m.StreamingPayload, scope, codegen.UserTypeLocation(m.StreamingPayload.Type), seen)...)
		ms = append(ms, collectUnionMethods(m.Result, scope, codegen.UserTypeLocation(m.Result.Type), seen)...)
		for _, e := range m.Errors {
			ms = append(ms, collectUnionMethods(e.AttributeExpr, scope, codegen.UserTypeLocation(e.Type), seen)...)
		}
	}
	sort.Slice(ms, func(i, j int) bool {
		return ms[i].Name < ms[j].Name
	})
	pkgs := make(map[string]struct{})
	for _, m := range ms {
		key := m.TypeRef + "::" + m.Name + "::" + m.Loc.PackageName()
		if _, ok := pkgs[key]; ok {
			continue
		}
		pkgs[key] = struct{}{}
		unionMethods = append(unionMethods, m)
	}

	desc := service.Description
	if desc == "" {
		desc = fmt.Sprintf("Service is the %s service interface.", service.Name)
	}

	varName := codegen.Goify(service.Name, false)
	data := &Data{
		Name:               service.Name,
		Description:        desc,
		APIName:            expr.Root.API.Name,
		APIVersion:         expr.Root.API.Version,
		VarName:            varName,
		PathName:           codegen.SnakeCase(varName),
		StructName:         codegen.Goify(service.Name, true),
		PkgName:            pkgName,
		ViewsPkg:           viewspkg,
		Methods:            methods,
		Schemes:            schemes,
		ServerInterceptors: collectInterceptors(service, methods, scope, true),
		ClientInterceptors: collectInterceptors(service, methods, scope, false),
		Scope:              scope,
		ViewScope:          viewScope,
		errorTypes:         errTypes,
		errorInits:         errorInits,
		userTypes:          types,
		projectedTypes:     projTypes,
		viewedUnionMethods: viewedUnionMeths,
		viewedResultTypes:  viewedRTs,
		unionValueMethods:  unionMethods,
	}

	d[service.Name] = data

	return data
}

// collectInterceptors returns the set of interceptors defined on the given
// service including any interceptor defined on specific service methods or API.
func collectInterceptors(svc *expr.ServiceExpr, methods []*MethodData, scope *codegen.NameScope, server bool) []*InterceptorData {
	var ints []*expr.InterceptorExpr
	if server {
		ints = expr.Root.API.ServerInterceptors
		ints = append(ints, svc.ServerInterceptors...)
		for _, m := range svc.Methods {
			ints = append(ints, m.ServerInterceptors...)
		}
	} else {
		ints = expr.Root.API.ClientInterceptors
		ints = append(ints, svc.ClientInterceptors...)
		for _, m := range svc.Methods {
			ints = append(ints, m.ClientInterceptors...)
		}
	}
	// remove duplicate interceptors
	sort.Slice(ints, func(i, j int) bool {
		return ints[i].Name < ints[j].Name
	})
	for i := 1; i < len(ints); i++ {
		if ints[i-1].Name == ints[i].Name {
			ints = append(ints[:i], ints[i+1:]...)
			i--
		}
	}

	res := make([]*InterceptorData, 0, len(ints))
	for _, i := range ints {
		res = append(res, buildInterceptorData(svc, methods, i, scope, server))
	}
	return res
}

// typeContext returns a contextual attribute for service types. Service types
// are Go types and uses non-pointers to hold attributes having default values.
func typeContext(pkg string, scope *codegen.NameScope) *codegen.AttributeContext {
	return codegen.NewAttributeContext(false, false, true, pkg, scope)
}

// projectedTypeContext returns a contextual attribute for a projected type.
// Projected types are Go types that uses pointers for all attributes (even the
// required ones).
func projectedTypeContext(pkg string, ptr bool, scope *codegen.NameScope) *codegen.AttributeContext {
	return codegen.NewAttributeContext(ptr, false, true, pkg, scope)
}

// collectTypes recurses through the attribute to gather all user types and
// records them in userTypes.
func collectTypes(at *expr.AttributeExpr, scope *codegen.NameScope, seen map[string]struct{}) (data []*UserTypeData) {
	if at == nil || at.Type == expr.Empty {
		return
	}
	collect := func(at *expr.AttributeExpr) []*UserTypeData { return collectTypes(at, scope, seen) }
	switch dt := at.Type.(type) {
	case expr.UserType:
		if _, ok := seen[dt.ID()]; ok {
			return nil
		}
		data = append(data, &UserTypeData{
			Name:        dt.Name(),
			VarName:     scope.GoTypeName(at),
			Description: dt.Attribute().Description,
			Def:         scope.GoTypeDef(dt.Attribute(), false, true),
			Ref:         scope.GoTypeRef(at),
			Loc:         codegen.UserTypeLocation(dt),
			Type:        dt,
		})
		seen[dt.ID()] = struct{}{}
		data = append(data, collect(dt.Attribute())...)
	case *expr.Object:
		for _, nat := range *dt {
			data = append(data, collect(nat.Attribute)...)
		}
	case *expr.Array:
		data = append(data, collect(dt.ElemType)...)
	case *expr.Map:
		data = append(data, collect(dt.KeyType)...)
		data = append(data, collect(dt.ElemType)...)
	case *expr.Union:
		for _, nat := range dt.Values {
			data = append(data, collect(nat.Attribute)...)
		}
	}
	return
}

// collectUnionMethods traverses the attribute to gather all union value methods.
func collectUnionMethods(att *expr.AttributeExpr, scope *codegen.NameScope, loc *codegen.Location, seen map[string]struct{}) (data []*UnionValueMethodData) {
	if att == nil || att.Type == expr.Empty {
		return
	}
	collect := func(at *expr.AttributeExpr, loc *codegen.Location) []*UnionValueMethodData {
		return collectUnionMethods(at, scope, loc, seen)
	}
	switch dt := att.Type.(type) {
	case expr.UserType:
		if _, ok := seen[dt.ID()]; ok {
			return nil
		}
		seen[dt.ID()] = struct{}{}
		data = append(data, collect(dt.Attribute(), codegen.UserTypeLocation(dt))...)
	case *expr.Object:
		for _, nat := range *dt {
			data = append(data, collect(nat.Attribute, loc)...)
		}
	case *expr.Array:
		data = append(data, collect(dt.ElemType, loc)...)
	case *expr.Map:
		data = append(data, collect(dt.KeyType, loc)...)
		data = append(data, collect(dt.ElemType, loc)...)
	case *expr.Union:
		for _, nat := range dt.Values {
			data = append(data, &UnionValueMethodData{
				Name:    codegen.UnionValTypeName(dt.Name()),
				TypeRef: scope.GoTypeRef(nat.Attribute),
				Loc:     loc,
			})
		}
		for _, nat := range dt.Values {
			data = append(data, collect(nat.Attribute, loc)...)
		}
	}
	return
}

// buildErrorInitData creates the data needed to generate code around endpoint error return values.
func buildErrorInitData(er *expr.ErrorExpr, scope *codegen.NameScope) *ErrorInitData {
	_, temporary := er.Meta["goa:error:temporary"]
	_, timeout := er.Meta["goa:error:timeout"]
	_, fault := er.Meta["goa:error:fault"]
	var pkg string
	if ut, ok := er.Type.(expr.UserType); ok {
		pkg = codegen.UserTypeLocation(ut).PackageName()
	}
	return &ErrorInitData{
		Name:        fmt.Sprintf("Make%s", codegen.Goify(er.Name, true)),
		Description: er.Description,
		ErrName:     er.Name,
		TypeName:    scope.GoTypeName(er.AttributeExpr),
		TypeRef:     scope.GoFullTypeRef(er.AttributeExpr, pkg),
		Temporary:   temporary,
		Timeout:     timeout,
		Fault:       fault,
	}
}

// buildMethodData creates the data needed to render the given endpoint. It
// records the user types needed by the service definition in userTypes.
func buildMethodData(m *expr.MethodExpr, scope *codegen.NameScope) *MethodData {
	var (
		vname       string
		desc        string
		payloadName string
		payloadLoc  *codegen.Location
		payloadDef  string
		payloadRef  string
		payloadDesc string
		payloadEx   any
		rname       string
		resultLoc   *codegen.Location
		resultDef   string
		resultRef   string
		resultDesc  string
		resultEx    any
		errors      []*ErrorInitData
		errorLocs   map[string]*codegen.Location
		reqs        RequirementsData
		schemes     SchemesData
	)
	vname = scope.Unique(codegen.Goify(m.Name, true), "Endpoint")
	desc = m.Description
	if desc == "" {
		desc = codegen.Goify(m.Name, true) + " implements " + m.Name + "."
	}
	if m.Payload.Type != expr.Empty {
		payloadName = scope.GoTypeName(m.Payload)
		if dt, ok := m.Payload.Type.(expr.UserType); ok {
			payloadDef = scope.GoTypeDef(dt.Attribute(), false, true)
			payloadLoc = codegen.UserTypeLocation(dt)
		}
		payloadRef = scope.GoFullTypeRef(m.Payload, payloadLoc.PackageName())
		payloadDesc = m.Payload.Description
		if payloadDesc == "" {
			payloadDesc = fmt.Sprintf("%s is the payload type of the %s service %s method.",
				payloadName, m.Service.Name, m.Name)
		}
		payloadEx = m.Payload.Example(expr.Root.API.ExampleGenerator)
	}
	if m.Result.Type != expr.Empty {
		rname = scope.GoTypeName(m.Result)
		if dt, ok := m.Result.Type.(expr.UserType); ok {
			resultDef = scope.GoTypeDef(dt.Attribute(), false, true)
			resultLoc = codegen.UserTypeLocation(dt)
		}
		resultRef = scope.GoFullTypeRef(m.Result, resultLoc.PackageName())
		resultDesc = m.Result.Description
		if resultDesc == "" {
			resultDesc = fmt.Sprintf("%s is the result type of the %s service %s method.",
				rname, m.Service.Name, m.Name)
		}
		resultEx = m.Result.Example(expr.Root.API.ExampleGenerator)
	}
	if len(m.Errors) > 0 {
		errors = make([]*ErrorInitData, len(m.Errors))
		errorLocs = make(map[string]*codegen.Location, len(m.Errors))
		for i, er := range m.Errors {
			errors[i] = buildErrorInitData(er, scope)
			errorLocs[er.Name] = codegen.UserTypeLocation(er.Type)
		}
	}
	for _, req := range m.Requirements {
		var rs SchemesData
		for _, s := range req.Schemes {
			sch := BuildSchemeData(s, m)
			rs = rs.Append(sch)
			schemes = schemes.Append(sch)
		}
		reqs = append(reqs, &RequirementData{Schemes: rs, Scopes: req.Scopes})
	}
	var httpMet *expr.HTTPEndpointExpr
	if httpSvc := expr.Root.HTTPService(m.Service.Name); httpSvc != nil {
		httpMet = httpSvc.Endpoint(m.Name)
	}
	data := &MethodData{
		Name:                         m.Name,
		VarName:                      vname,
		Description:                  desc,
		Payload:                      payloadName,
		PayloadLoc:                   payloadLoc,
		PayloadDef:                   payloadDef,
		PayloadRef:                   payloadRef,
		PayloadDesc:                  payloadDesc,
		PayloadEx:                    payloadEx,
		PayloadDefault:               m.Payload.DefaultValue,
		Result:                       rname,
		ResultLoc:                    resultLoc,
		ResultDef:                    resultDef,
		ResultRef:                    resultRef,
		ResultDesc:                   resultDesc,
		ResultEx:                     resultEx,
		Errors:                       errors,
		ErrorLocs:                    errorLocs,
		Requirements:                 reqs,
		Schemes:                      schemes,
		StreamKind:                   m.Stream,
		SkipRequestBodyEncodeDecode:  httpMet != nil && httpMet.SkipRequestBodyEncodeDecode,
		SkipResponseBodyEncodeDecode: httpMet != nil && httpMet.SkipResponseBodyEncodeDecode,
		RequestStruct:                vname + "RequestData",
		ResponseStruct:               vname + "ResponseData",
	}
	initStreamData(data, m, vname, rname, resultRef, scope)
	return data
}

// initStreamData initializes the streaming payload data structures and methods.
func initStreamData(data *MethodData, m *expr.MethodExpr, vname, rname, resultRef string, scope *codegen.NameScope) {
	if !m.IsStreaming() {
		return
	}
	var (
		spayloadName string
		spayloadRef  string
		spayloadDef  string
		spayloadDesc string
		spayloadEx   any
	)
	if m.StreamingPayload.Type != expr.Empty {
		spayloadName = scope.GoTypeName(m.StreamingPayload)
		spayloadRef = scope.GoTypeRef(m.StreamingPayload)
		if dt, ok := m.StreamingPayload.Type.(expr.UserType); ok {
			spayloadDef = scope.GoTypeDef(dt.Attribute(), false, true)
		}
		spayloadDesc = m.StreamingPayload.Description
		if spayloadDesc == "" {
			spayloadDesc = fmt.Sprintf("%s is the streaming payload type of the %s service %s method.",
				spayloadName, m.Service.Name, m.Name)
		}
		spayloadEx = m.StreamingPayload.Example(expr.Root.API.ExampleGenerator)
	}
	svrStream := &StreamData{
		Interface:           vname + "ServerStream",
		VarName:             scope.Unique(codegen.Goify(m.Name, true), "ServerStream"),
		EndpointStruct:      vname + "EndpointInput",
		Kind:                m.Stream,
		SendName:            "Send",
		SendDesc:            fmt.Sprintf("Send streams instances of %q.", rname),
		SendWithContextName: "SendWithContext",
		SendWithContextDesc: fmt.Sprintf("SendWithContext streams instances of %q with context.", rname),
		SendTypeName:        rname,
		SendTypeRef:         resultRef,
		MustClose:           true,
	}
	cliStream := &StreamData{
		Interface:           vname + "ClientStream",
		VarName:             scope.Unique(codegen.Goify(m.Name, true), "ClientStream"),
		Kind:                m.Stream,
		RecvName:            "Recv",
		RecvDesc:            fmt.Sprintf("Recv reads instances of %q from the stream.", rname),
		RecvWithContextName: "RecvWithContext",
		RecvWithContextDesc: fmt.Sprintf("RecvWithContext reads instances of %q from the stream with context.", rname),
		RecvTypeName:        rname,
		RecvTypeRef:         resultRef,
	}
	if m.Stream == expr.ClientStreamKind || m.Stream == expr.BidirectionalStreamKind {
		switch m.Stream {
		case expr.ClientStreamKind:
			if resultRef != "" {
				svrStream.SendName = "SendAndClose"
				svrStream.SendDesc = fmt.Sprintf("SendAndClose streams instances of %q and closes the stream.", rname)
				svrStream.SendWithContextName = "SendAndCloseWithContext"
				svrStream.SendWithContextDesc = fmt.Sprintf("SendAndCloseWithContext streams instances of %q and closes the stream with context.", rname)
				svrStream.MustClose = false
				cliStream.RecvName = "CloseAndRecv"
				cliStream.RecvDesc = fmt.Sprintf("CloseAndRecv stops sending messages to the stream and reads instances of %q from the stream.", rname)
				cliStream.RecvWithContextName = "CloseAndRecvWithContext"
				cliStream.RecvWithContextDesc = fmt.Sprintf("CloseAndRecvWithContext stops sending messages to the stream and reads instances of %q from the stream with context.", rname)
			} else {
				cliStream.MustClose = true
			}
		case expr.BidirectionalStreamKind:
			cliStream.MustClose = true
		}
		svrStream.RecvName = "Recv"
		svrStream.RecvDesc = fmt.Sprintf("Recv reads instances of %q from the stream.", spayloadName)
		svrStream.RecvWithContextName = "RecvWithContext"
		svrStream.RecvWithContextDesc = fmt.Sprintf("RecvWithContext reads instances of %q from the stream with context.", spayloadName)
		svrStream.RecvTypeName = spayloadName
		svrStream.RecvTypeRef = spayloadRef
		cliStream.SendName = "Send"
		cliStream.SendDesc = fmt.Sprintf("Send streams instances of %q.", spayloadName)
		cliStream.SendWithContextName = "SendWithContext"
		cliStream.SendWithContextDesc = fmt.Sprintf("SendWithContext streams instances of %q with context.", spayloadName)
		cliStream.SendTypeName = spayloadName
		cliStream.SendTypeRef = spayloadRef
	}
	data.ClientStream = cliStream
	data.ServerStream = svrStream
	data.StreamingPayload = spayloadName
	data.StreamingPayloadDef = spayloadDef
	data.StreamingPayloadRef = spayloadRef
	data.StreamingPayloadDesc = spayloadDesc
	data.StreamingPayloadEx = spayloadEx
}

// buildInterceptorData creates the data needed to generate interceptor code.
func buildInterceptorData(svc *expr.ServiceExpr, methods []*MethodData, i *expr.InterceptorExpr, scope *codegen.NameScope, server bool) *InterceptorData {
	data := &InterceptorData{
		Name:        codegen.Goify(i.Name, true),
		DesignName:  i.Name,
		Description: i.Description,
	}
	if len(svc.Methods) == 0 {
		return data
	}
	attributesCollected := false
	for _, m := range svc.Methods {
		applies := false
		intExprs := m.ServerInterceptors
		if !server {
			intExprs = m.ClientInterceptors
		}
		for _, in := range intExprs {
			if in.Name == i.Name {
				if !attributesCollected {
					payload, result, streamingPayload := m.Payload, m.Result, m.StreamingPayload
					data.ReadPayload = collectAttributes(i.ReadPayload, payload, scope)
					data.WritePayload = collectAttributes(i.WritePayload, payload, scope)
					data.ReadResult = collectAttributes(i.ReadResult, result, scope)
					data.WriteResult = collectAttributes(i.WriteResult, result, scope)
					data.ReadStreamingPayload = collectAttributes(i.ReadStreamingPayload, streamingPayload, scope)
					data.WriteStreamingPayload = collectAttributes(i.WriteStreamingPayload, streamingPayload, scope)
					data.ReadStreamingResult = collectAttributes(i.ReadStreamingResult, result, scope)
					data.WriteStreamingResult = collectAttributes(i.WriteStreamingResult, result, scope)
					if len(data.ReadPayload) > 0 || len(data.WritePayload) > 0 {
						data.HasPayloadAccess = true
					}
					if len(data.ReadResult) > 0 || len(data.WriteResult) > 0 {
						data.HasResultAccess = true
					}
					if len(data.ReadStreamingPayload) > 0 || len(data.WriteStreamingPayload) > 0 {
						data.HasStreamingPayloadAccess = true
					}
					if len(data.ReadStreamingResult) > 0 || len(data.WriteStreamingResult) > 0 {
						data.HasStreamingResultAccess = true
					}
					attributesCollected = true
				}
				applies = true
				break
			}
		}
		if !applies {
			continue
		}
		var md *MethodData
		for _, mt := range methods {
			if m.Name == mt.Name {
				md = mt
				break
			}
		}
		data.Methods = append(data.Methods, buildInterceptorMethodData(i, md))
		if server {
			md.ServerInterceptors = append(md.ServerInterceptors, i.Name)
		} else {
			md.ClientInterceptors = append(md.ClientInterceptors, i.Name)
		}
	}
	return data
}

// buildInterceptorMethodData creates the data needed to generate interceptor
// method code.
func buildInterceptorMethodData(i *expr.InterceptorExpr, md *MethodData) *MethodInterceptorData {
	var serverStream, clientStream *StreamInterceptorData
	if md.ServerStream != nil {
		serverStream = &StreamInterceptorData{
			Interface:           md.ServerStream.Interface,
			SendName:            md.ServerStream.SendName,
			SendWithContextName: md.ServerStream.SendWithContextName,
			SendTypeRef:         md.ServerStream.SendTypeRef,
			RecvName:            md.ServerStream.RecvName,
			RecvWithContextName: md.ServerStream.RecvWithContextName,
			RecvTypeRef:         md.ServerStream.RecvTypeRef,
			MustClose:           md.ServerStream.MustClose,
			EndpointStruct:      md.ServerStream.EndpointStruct,
		}
	}
	if md.ClientStream != nil {
		clientStream = &StreamInterceptorData{
			Interface:           md.ClientStream.Interface,
			SendName:            md.ClientStream.SendName,
			SendWithContextName: md.ClientStream.SendWithContextName,
			SendTypeRef:         md.ClientStream.SendTypeRef,
			RecvName:            md.ClientStream.RecvName,
			RecvWithContextName: md.ClientStream.RecvWithContextName,
			RecvTypeRef:         md.ClientStream.RecvTypeRef,
			MustClose:           md.ClientStream.MustClose,
		}
	}
	var payloadAccess, resultAccess, streamingPayloadAccess, streamingResultAccess string
	if i.ReadPayload != nil || i.WritePayload != nil {
		payloadAccess = codegen.Goify(i.Name, false) + md.VarName + "Payload"
	}
	if i.ReadResult != nil || i.WriteResult != nil {
		resultAccess = codegen.Goify(i.Name, false) + md.VarName + "Result"
	}
	if i.ReadStreamingPayload != nil || i.WriteStreamingPayload != nil {
		streamingPayloadAccess = codegen.Goify(i.Name, false) + md.VarName + "StreamingPayload"
	}
	if i.ReadStreamingResult != nil || i.WriteStreamingResult != nil {
		streamingResultAccess = codegen.Goify(i.Name, false) + md.VarName + "StreamingResult"
	}
	return &MethodInterceptorData{
		MethodName:             md.VarName,
		PayloadAccess:          payloadAccess,
		ResultAccess:           resultAccess,
		PayloadRef:             md.PayloadRef,
		ResultRef:              md.ResultRef,
		StreamingPayloadAccess: streamingPayloadAccess,
		StreamingPayloadRef:    md.StreamingPayloadRef,
		StreamingResultAccess:  streamingResultAccess,
		StreamingResultRef:     md.ResultRef,
		ClientStream:           clientStream,
		ServerStream:           serverStream,
	}
}

// BuildSchemeData builds the scheme data for the given scheme and method expr.
func BuildSchemeData(s *expr.SchemeExpr, m *expr.MethodExpr) *SchemeData {
	if !expr.IsObject(m.Payload.Type) {
		return nil
	}
	switch s.Kind {
	case expr.BasicAuthKind:
		userAtt := expr.TaggedAttribute(m.Payload, "security:username")
		user := codegen.Goify(userAtt, true)
		passAtt := expr.TaggedAttribute(m.Payload, "security:password")
		pass := codegen.Goify(passAtt, true)
		var scopes []string
		if len(s.Scopes) > 0 {
			scopes = make([]string, len(s.Scopes))
			for i, s := range s.Scopes {
				scopes[i] = s.Name
			}
		}
		return &SchemeData{
			Type:             s.Kind.String(),
			SchemeName:       s.SchemeName,
			UsernameAttr:     userAtt,
			UsernameField:    user,
			UsernamePointer:  m.Payload.IsPrimitivePointer(userAtt, true),
			UsernameRequired: m.Payload.IsRequired(userAtt),
			PasswordAttr:     passAtt,
			PasswordField:    pass,
			PasswordPointer:  m.Payload.IsPrimitivePointer(passAtt, true),
			PasswordRequired: m.Payload.IsRequired(passAtt),
			Scopes:           scopes,
		}
	case expr.APIKeyKind:
		if keyAtt := expr.TaggedAttribute(m.Payload, "security:apikey:"+s.SchemeName); keyAtt != "" {
			key := codegen.Goify(keyAtt, true)
			var scopes []string
			if len(s.Scopes) > 0 {
				scopes = make([]string, len(s.Scopes))
				for i, s := range s.Scopes {
					scopes[i] = s.Name
				}
			}
			return &SchemeData{
				Type:         s.Kind.String(),
				Name:         s.Name,
				SchemeName:   s.SchemeName,
				CredField:    key,
				CredPointer:  m.Payload.IsPrimitivePointer(keyAtt, true),
				CredRequired: m.Payload.IsRequired(keyAtt),
				KeyAttr:      keyAtt,
				Scopes:       scopes,
				In:           s.In,
			}
		}
	case expr.JWTKind:
		if keyAtt := expr.TaggedAttribute(m.Payload, "security:token"); keyAtt != "" {
			key := codegen.Goify(keyAtt, true)
			var scopes []string
			if len(s.Scopes) > 0 {
				scopes = make([]string, len(s.Scopes))
				for i, s := range s.Scopes {
					scopes[i] = s.Name
				}
			}
			return &SchemeData{
				Type:         s.Kind.String(),
				Name:         s.Name,
				SchemeName:   s.SchemeName,
				CredField:    key,
				CredPointer:  m.Payload.IsPrimitivePointer(keyAtt, true),
				CredRequired: m.Payload.IsRequired(keyAtt),
				KeyAttr:      keyAtt,
				Scopes:       scopes,
				In:           s.In,
			}
		}
	case expr.OAuth2Kind:
		if keyAtt := expr.TaggedAttribute(m.Payload, "security:accesstoken"); keyAtt != "" {
			key := codegen.Goify(keyAtt, true)
			var scopes []string
			if len(s.Scopes) > 0 {
				scopes = make([]string, len(s.Scopes))
				for i, s := range s.Scopes {
					scopes[i] = s.Name
				}
			}
			return &SchemeData{
				Type:         s.Kind.String(),
				Name:         s.Name,
				SchemeName:   s.SchemeName,
				CredField:    key,
				CredPointer:  m.Payload.IsPrimitivePointer(keyAtt, true),
				CredRequired: m.Payload.IsRequired(keyAtt),
				KeyAttr:      keyAtt,
				Scopes:       scopes,
				Flows:        s.Flows,
				In:           s.In,
			}
		}
	}
	return nil
}

// collectAttributes builds AttributeData from an AttributeExpr
func collectAttributes(attrNames, parent *expr.AttributeExpr, scope *codegen.NameScope) []*AttributeData {
	if attrNames == nil {
		return nil
	}
	obj := expr.AsObject(attrNames.Type)
	if obj == nil {
		return nil
	}
	data := make([]*AttributeData, len(*obj))
	for i, nat := range *obj {
		parentAttr := parent.Find(nat.Name)
		if parentAttr == nil {
			continue
		}
		var pkg string
		if loc := codegen.UserTypeLocation(parentAttr.Type); loc != nil {
			pkg = loc.PackageName()
		}
		data[i] = &AttributeData{
			Name:    codegen.Goify(nat.Name, true),
			TypeRef: scope.GoFullTypeRef(parentAttr, pkg),
			Pointer: parent.IsPrimitivePointer(nat.Name, true),
		}
	}
	return data
}

// collectProjectedTypes builds a projected type for every user type found
// when recursing through the attributes. The projected types live in the views
// package and support the marshaling and unmarshalling of result types that
// make use of views. We need to build projected types for all user types - not
// just result types - because user types make contain result types and thus may
// need to be marshalled in different ways depending on the view being used.
func collectProjectedTypes(projected, att *expr.AttributeExpr, viewspkg string, scope, viewScope *codegen.NameScope, seen map[string]*ProjectedTypeData) (data []*ProjectedTypeData, umeths []*UnionValueMethodData) {
	collect := func(projected, att *expr.AttributeExpr) ([]*ProjectedTypeData, []*UnionValueMethodData) {
		return collectProjectedTypes(projected, att, viewspkg, scope, viewScope, seen)
	}
	switch pt := projected.Type.(type) {
	case expr.UserType:
		dt := att.Type.(expr.UserType)
		if pd, ok := seen[dt.ID()]; ok {
			// a projected type is already created for this user type. We change the
			// attribute type to this seen projected type. The seen projected type
			// can be nil if the attribute type has a circular type definition in
			// which case we don't change the attribute type until the projected type
			// is created during the recursion.
			if pd != nil {
				projected.Type = pd.Type
			}
			return
		}
		seen[dt.ID()] = nil
		pt.Rename(pt.Name() + "View")
		// We recurse before building the projected type so that user types within
		// a projected type is also converted to their respective projected types.
		types, ms := collect(pt.Attribute(), dt.Attribute())
		pd := buildProjectedType(projected, att, viewspkg, scope, viewScope)
		seen[dt.ID()] = pd
		data = append(data, pd)
		data = append(data, types...)
		umeths = append(umeths, ms...)
	case *expr.Array:
		dt := att.Type.(*expr.Array)
		types, ms := collect(pt.ElemType, dt.ElemType)
		data = append(data, types...)
		umeths = append(umeths, ms...)
	case *expr.Map:
		dt := att.Type.(*expr.Map)
		types, ms := collect(pt.KeyType, dt.KeyType)
		data = append(data, types...)
		umeths = append(umeths, ms...)
		types, ms = collect(pt.ElemType, dt.ElemType)
		data = append(data, types...)
		umeths = append(umeths, ms...)
	case *expr.Object:
		dt := att.Type.(*expr.Object)
		for _, n := range *pt {
			types, ms := collect(n.Attribute, dt.Attribute(n.Name))
			data = append(data, types...)
			umeths = append(umeths, ms...)
		}
	case *expr.Union:
		dt := att.Type.(*expr.Union)
		for i, n := range pt.Values {
			types, ms := collect(n.Attribute, dt.Values[i].Attribute)
			data = append(data, types...)
			umeths = append(umeths, ms...)
		}
		for _, nat := range pt.Values {
			umeths = append(umeths, &UnionValueMethodData{
				Name:    codegen.UnionValTypeName(pt.Name()),
				TypeRef: scope.GoTypeRef(nat.Attribute),
			})
		}
	}
	return
}

// hasResultType returns true if the given attribute has a result type recursively.
func hasResultType(att *expr.AttributeExpr, seens ...map[string]struct{}) bool {
	if _, ok := att.Type.(*expr.ResultTypeExpr); ok {
		return true
	}
	var seen map[string]struct{}
	if len(seens) > 0 {
		seen = seens[0]
	} else {
		seen = make(map[string]struct{})
	}
	switch a := att.Type.(type) {
	case expr.UserType:
		if _, ok := seen[a.ID()]; ok {
			return false
		}
		seen[a.ID()] = struct{}{}
		return hasResultType(a.Attribute(), seen)
	case *expr.Array:
		return hasResultType(a.ElemType, seen)
	case *expr.Map:
		return hasResultType(a.KeyType, seen) || hasResultType(a.ElemType, seen)
	case *expr.Object:
		for _, nat := range *a {
			if hasResultType(nat.Attribute, seen) {
				return true
			}
		}
	case *expr.Union:
		for _, nat := range a.Values {
			if hasResultType(nat.Attribute, seen) {
				return true
			}
		}
	}
	return false
}

// buildProjectedType builds projected type for the given user type.
//
// viewspkg is the name of the views package
func buildProjectedType(projected, att *expr.AttributeExpr, viewspkg string, scope, viewScope *codegen.NameScope) *ProjectedTypeData {
	var (
		projections []*InitData
		typeInits   []*InitData
		views       []*ViewData

		varname = viewScope.GoTypeName(projected)
		pt      = projected.Type.(expr.UserType)
	)
	if _, isrt := pt.(*expr.ResultTypeExpr); isrt {
		typeInits = buildTypeInits(projected, att, viewspkg, scope, viewScope)
		projections = buildProjections(projected, att, viewspkg, scope, viewScope)
		views = buildViews(att.Type.(*expr.ResultTypeExpr), viewScope)
	}
	validations := buildValidations(projected, viewScope)
	removeMeta(projected)
	return &ProjectedTypeData{
		UserTypeData: &UserTypeData{
			Name:        varname,
			Description: fmt.Sprintf("%s is a type that runs validations on a projected type.", varname),
			VarName:     varname,
			Def:         viewScope.GoTypeDef(pt.Attribute(), true, true),
			Ref:         viewScope.GoTypeRef(projected),
			Type:        pt,
		},
		Projections: projections,
		TypeInits:   typeInits,
		Validations: validations,
		ViewsPkg:    viewspkg,
		Views:       views,
	}
}

// buildViews builds the view data for all the views in the given result type.
func buildViews(rt *expr.ResultTypeExpr, viewScope *codegen.NameScope) []*ViewData {
	views := make([]*ViewData, len(rt.Views))
	for i, view := range rt.Views {
		vatt := expr.AsObject(view.Type)
		attrs := make([]string, len(*vatt))
		for j, nat := range *vatt {
			attrs[j] = nat.Name
		}
		views[i] = &ViewData{
			Name:        view.Name,
			Description: view.Description,
			Attributes:  attrs,
			TypeVarName: viewScope.GoTypeName(&expr.AttributeExpr{Type: rt}),
		}
	}
	return views
}

// buildViewedResultType builds a viewed result type from the given result type
// and projected type.
func buildViewedResultType(att, projected *expr.AttributeExpr, viewspkg string, scope, viewScope *codegen.NameScope) *ViewedResultTypeData {
	// collect result type views
	rt := att.Type.(*expr.ResultTypeExpr)
	isarr := expr.IsArray(att.Type)
	var viewName string
	if !rt.HasMultipleViews() {
		viewName = expr.DefaultView
	}
	if v, ok := att.Meta.Last(expr.ViewMetaKey); ok {
		viewName = v
	}
	views := buildViews(rt, viewScope)

	// build validation data
	resvar := scope.GoTypeName(att)
	resref := scope.GoTypeRef(att)
	data := map[string]any{
		"Projected": scope.GoTypeName(projected),
		"ArgVar":    "result",
		"Source":    "result",
		"Views":     views,
		"IsViewed":  true,
	}
	buf := &bytes.Buffer{}
	if err := validateTypeCodeTmpl.Execute(buf, data); err != nil {
		panic(err) // bug
	}
	name := "Validate" + resvar
	validate := &ValidateData{
		Name:        name,
		Description: fmt.Sprintf("%s runs the validations defined on the viewed result type %s.", name, resvar),
		Ref:         resref,
		Validate:    buf.String(),
	}

	// build constructor to initialize viewed result type from result type
	vresref := viewScope.GoFullTypeRef(att, viewspkg)
	data = map[string]any{
		"ToViewed":      true,
		"ArgVar":        "res",
		"ReturnVar":     "vres",
		"Views":         views,
		"ReturnTypeRef": vresref,
		"IsCollection":  isarr,
		"TargetType":    scope.GoFullTypeName(att, viewspkg),
		"InitName":      "new" + viewScope.GoTypeName(projected),
	}
	buf = &bytes.Buffer{}
	if err := initTypeCodeTmpl.Execute(buf, data); err != nil {
		panic(err) // bug
	}
	pkg := ""
	if loc := codegen.UserTypeLocation(att.Type); loc != nil {
		pkg = loc.PackageName()
	}
	name = "NewViewed" + resvar
	init := &InitData{
		Name:        name,
		Description: fmt.Sprintf("%s initializes viewed result type %s from result type %s using the given view.", name, resvar, resvar),
		Args: []*InitArgData{
			{Name: "res", Ref: scope.GoFullTypeRef(att, pkg)},
			{Name: "view", Ref: "string"},
		},
		ReturnTypeRef: vresref,
		Code:          buf.String(),
	}

	// build constructor to initialize result type from viewed result type
	if loc := codegen.UserTypeLocation(att.Type); loc != nil {
		resref = scope.GoFullTypeRef(att, loc.PackageName())
	}
	data = map[string]any{
		"ToResult":      true,
		"ArgVar":        "vres",
		"ReturnVar":     "res",
		"Views":         views,
		"ReturnTypeRef": resref,
		"InitName":      "new" + scope.GoTypeName(att),
	}
	buf = &bytes.Buffer{}
	if err := initTypeCodeTmpl.Execute(buf, data); err != nil {
		panic(err) // bug
	}
	name = "New" + resvar
	resinit := &InitData{
		Name:          name,
		Description:   fmt.Sprintf("%s initializes result type %s from viewed result type %s.", name, resvar, resvar),
		Args:          []*InitArgData{{Name: "vres", Ref: scope.GoFullTypeRef(att, viewspkg)}},
		ReturnTypeRef: resref,
		Code:          buf.String(),
	}

	projT := wrapProjected(projected.Type.(expr.UserType))
	return &ViewedResultTypeData{
		UserTypeData: &UserTypeData{
			Name:        resvar,
			Description: fmt.Sprintf("%s is the viewed result type that is projected based on a view.", resvar),
			VarName:     resvar,
			Def:         viewScope.GoTypeDef(projT.Attribute(), false, true),
			Ref:         resref,
			Type:        projT,
		},
		FullName:     scope.GoFullTypeName(att, viewspkg),
		FullRef:      vresref,
		ResultInit:   resinit,
		Init:         init,
		Views:        views,
		Validate:     validate,
		IsCollection: isarr,
		ViewName:     viewName,
		ViewsPkg:     viewspkg,
	}
}

// wrapProjected builds a viewed result type by wrapping the given projected
// in a result type with "projected" and "view" attributes.
func wrapProjected(projected expr.UserType) expr.UserType {
	rt := projected.(*expr.ResultTypeExpr)
	pratt := &expr.NamedAttributeExpr{
		Name:      "projected",
		Attribute: &expr.AttributeExpr{Type: rt, Description: "Type to project"},
	}
	prview := &expr.NamedAttributeExpr{
		Name:      "view",
		Attribute: &expr.AttributeExpr{Type: expr.String, Description: "View to render"},
	}
	return &expr.ResultTypeExpr{
		UserTypeExpr: &expr.UserTypeExpr{
			AttributeExpr: &expr.AttributeExpr{
				Type:       &expr.Object{pratt, prview},
				Validation: &expr.ValidationExpr{Required: []string{"projected", "view"}},
			},
			TypeName: rt.TypeName,
		},
		Identifier: rt.Identifier,
		Views:      rt.Views,
	}
}

// buildTypeInits builds the data to generate the constructor code to
// initialize a result type from a projected type.
func buildTypeInits(projected, att *expr.AttributeExpr, viewspkg string, scope, viewScope *codegen.NameScope) []*InitData {
	prt := projected.Type.(*expr.ResultTypeExpr)
	pobj := expr.AsObject(projected.Type)
	parr := expr.AsArray(projected.Type)
	if parr != nil {
		// result type collection
		pobj = expr.AsObject(parr.ElemType.Type)
	}

	// For every view defined in the result type, build a constructor function
	// to create the result type from a projected type based on the view.
	init := make([]*InitData, 0, len(prt.Views))
	for _, view := range prt.Views {
		var typ expr.DataType
		obj := &expr.Object{}
		walkViewAttrs(pobj, view, func(name string, att, _ *expr.AttributeExpr) {
			obj.Set(name, att)
		})
		typ = obj
		if parr != nil {
			typ = &expr.Array{ElemType: &expr.AttributeExpr{
				Type: &expr.ResultTypeExpr{
					UserTypeExpr: &expr.UserTypeExpr{
						AttributeExpr: &expr.AttributeExpr{Type: obj},
						TypeName:      scope.GoTypeName(parr.ElemType),
					},
				},
			}}
		}
		src := &expr.AttributeExpr{
			Type: &expr.ResultTypeExpr{
				UserTypeExpr: &expr.UserTypeExpr{
					AttributeExpr: &expr.AttributeExpr{Type: typ},
					TypeName:      scope.GoTypeName(projected),
				},
				Views:      prt.Views,
				Identifier: prt.Identifier,
			},
		}

		srcCtx := projectedTypeContext(viewspkg, true, viewScope)
		tgtCtx := typeContext("", scope)
		resvar := scope.GoTypeName(att)
		name := "new" + resvar
		if view.Name != expr.DefaultView {
			name += codegen.Goify(view.Name, true)
		}
		code, helpers := buildConstructorCode(src, att, "vres", "res", srcCtx, tgtCtx, view.Name)

		pkg := ""
		if loc := codegen.UserTypeLocation(att.Type); loc != nil {
			pkg = loc.PackageName()
		}
		init = append(init, &InitData{
			Name:          name,
			Description:   fmt.Sprintf("%s converts projected type %s to service type %s.", name, resvar, resvar),
			Args:          []*InitArgData{{Name: "vres", Ref: viewScope.GoFullTypeRef(projected, viewspkg)}},
			ReturnTypeRef: scope.GoFullTypeRef(att, pkg),
			Code:          code,
			Helpers:       helpers,
		})
	}
	return init
}

// buildProjections builds the data to generate the constructor code to
// project a result type to a projected type based on a view.
func buildProjections(projected, att *expr.AttributeExpr, viewspkg string, scope, viewScope *codegen.NameScope) []*InitData {
	rt := att.Type.(*expr.ResultTypeExpr)
	projections := make([]*InitData, 0, len(rt.Views))
	for _, view := range rt.Views {
		var typ expr.DataType
		obj := &expr.Object{}
		pobj := expr.AsObject(projected.Type)
		parr := expr.AsArray(projected.Type)
		if parr != nil {
			// result type collection
			pobj = expr.AsObject(parr.ElemType.Type)
		}
		walkViewAttrs(pobj, view, func(name string, att, _ *expr.AttributeExpr) {
			obj.Set(name, att)
		})
		typ = obj
		if parr != nil {
			typ = &expr.Array{ElemType: &expr.AttributeExpr{
				Type: &expr.ResultTypeExpr{
					UserTypeExpr: &expr.UserTypeExpr{
						AttributeExpr: &expr.AttributeExpr{Type: obj},
						TypeName:      parr.ElemType.Type.Name(),
					},
				},
			}}
		}
		tgt := &expr.AttributeExpr{
			Type: &expr.ResultTypeExpr{
				UserTypeExpr: &expr.UserTypeExpr{
					AttributeExpr: &expr.AttributeExpr{Type: typ},
					TypeName:      projected.Type.Name(),
				},
				Views:      rt.Views,
				Identifier: rt.Identifier,
			},
		}

		srcCtx := typeContext("", scope)
		tgtCtx := projectedTypeContext(viewspkg, true, viewScope)
		tname := scope.GoTypeName(projected)
		name := "new" + tname
		if view.Name != expr.DefaultView {
			name += codegen.Goify(view.Name, true)
		}
		code, helpers := buildConstructorCode(att, tgt, "res", "vres", srcCtx, tgtCtx, view.Name)

		pkg := ""
		if loc := codegen.UserTypeLocation(att.Type); loc != nil {
			pkg = loc.PackageName()
		}
		projections = append(projections, &InitData{
			Name:          name,
			Description:   fmt.Sprintf("%s projects result type %s to projected type %s using the %q view.", name, scope.GoTypeName(att), tname, view.Name),
			Args:          []*InitArgData{{Name: "res", Ref: scope.GoFullTypeRef(att, pkg)}},
			ReturnTypeRef: viewScope.GoFullTypeRef(projected, viewspkg),
			Code:          code,
			Helpers:       helpers,
		})
	}
	return projections
}

// buildValidations builds the data required to generate validations for the
// projected types.
func buildValidations(projected *expr.AttributeExpr, scope *codegen.NameScope) []*ValidateData {
	ut := projected.Type.(expr.UserType)
	tname := scope.GoTypeName(projected)
	var validations []*ValidateData
	if rt, isrt := ut.(*expr.ResultTypeExpr); isrt {
		// for result types we create a validation function containing view
		// specific validation logic for each view
		arr := expr.AsArray(projected.Type)
		for _, view := range rt.Views {
			data := map[string]any{
				"Projected":    tname,
				"ArgVar":       "result",
				"Source":       "result",
				"IsCollection": arr != nil,
			}
			var vn string
			name := "Validate" + tname
			if view.Name != expr.DefaultView {
				vn = codegen.Goify(view.Name, true)
				name += vn
			}

			if arr != nil {
				// dealing with an array type
				data["Source"] = "item"
				data["ValidateVar"] = "Validate" + scope.GoTypeName(arr.ElemType) + vn
			} else {
				var fields []map[string]any
				o := &expr.Object{}
				walkViewAttrs(expr.AsObject(projected.Type), view, func(name string, attr, vatt *expr.AttributeExpr) {
					if rt, ok := attr.Type.(*expr.ResultTypeExpr); ok {
						// use explicitly specified view (if any) for the attribute,
						// otherwise use default
						vw := ""
						if v, ok := vatt.Meta.Last(expr.ViewMetaKey); ok && v != expr.DefaultView {
							vw = v
						}
						fields = append(fields, map[string]any{
							"Name":        name,
							"ValidateVar": "Validate" + scope.GoTypeName(attr) + codegen.Goify(vw, true),
							"IsRequired":  rt.Attribute().IsRequired(name),
						})
					} else {
						o.Set(name, attr)
					}
				})
				ctx := projectedTypeContext("", !expr.IsPrimitive(projected.Type), scope)
				data["Validate"] = codegen.ValidationCode(&expr.AttributeExpr{Type: o, Validation: rt.Validation}, rt, ctx, true, false, true, "result")
				data["Fields"] = fields
			}

			buf := &bytes.Buffer{}
			if err := validateTypeCodeTmpl.Execute(buf, data); err != nil {
				panic(err) // bug
			}

			validations = append(validations, &ValidateData{
				Name:        name,
				Description: fmt.Sprintf("%s runs the validations defined on %s using the %q view.", name, tname, view.Name),
				Ref:         scope.GoTypeRef(projected),
				Validate:    buf.String(),
			})
		}
	} else {
		// for a user type or a result type with single view, we generate only one validation
		// function containing the validation logic
		name := "Validate" + tname
		ctx := projectedTypeContext("", !expr.IsPrimitive(projected.Type), scope)
		validations = append(validations, &ValidateData{
			Name:        name,
			Description: fmt.Sprintf("%s runs the validations defined on %s.", name, tname),
			Ref:         scope.GoTypeRef(projected),
			Validate:    codegen.ValidationCode(ut.Attribute(), ut, ctx, true, expr.IsAlias(ut), true, "result"),
		})
	}
	return validations
}

// buildConstructorCode builds the transformation code to create a projected
// type from a service type and vice versa.
//
// source and target contains the projected/service contextual attributes
//
// sourceVar and targetVar contains the variable name that holds the source and
// target data structures in the transformation code.
//
// view is used to generate the constructor function name.
func buildConstructorCode(src, tgt *expr.AttributeExpr, sourceVar, targetVar string, sourceCtx, targetCtx *codegen.AttributeContext, view string) (string, []*codegen.TransformFunctionData) {
	var (
		helpers []*codegen.TransformFunctionData
		buf     bytes.Buffer
	)
	rt := src.Type.(*expr.ResultTypeExpr)
	arr := expr.AsArray(tgt.Type)

	data := map[string]any{
		"ArgVar":       sourceVar,
		"ReturnVar":    targetVar,
		"IsCollection": arr != nil,
		"TargetType":   targetCtx.Scope.Name(tgt, targetCtx.Pkg(tgt), targetCtx.Pointer, targetCtx.UseDefault),
	}

	if arr != nil {
		// result type collection
		init := "new" + targetCtx.Scope.Name(arr.ElemType, "", targetCtx.Pointer, targetCtx.UseDefault)
		if view != "" && view != expr.DefaultView {
			init += codegen.Goify(view, true)
		}
		data["InitName"] = init
		if err := initTypeCodeTmpl.Execute(&buf, data); err != nil {
			panic(err) // bug
		}
		return buf.String(), helpers
	}

	// service type to projected type (or vice versa)
	targetRTs := &expr.Object{}
	tatt := expr.DupAtt(tgt)
	tobj := expr.AsObject(tatt.Type)
	for _, nat := range *tobj {
		if _, ok := nat.Attribute.Type.(*expr.ResultTypeExpr); ok {
			targetRTs.Set(nat.Name, nat.Attribute)
			tobj.Delete(nat.Name)
		}
	}
	data["Source"] = sourceVar
	data["Target"] = targetVar

	// build code for target with no result types
	code, helpers, err := codegen.GoTransform(src, tatt, sourceVar, targetVar, sourceCtx, targetCtx, "transform", true)
	if err != nil {
		panic(err) // bug
	}
	data["Code"] = code

	if view != "" {
		data["InitName"] = targetCtx.Scope.Name(src, "", targetCtx.Pointer, targetCtx.UseDefault)
	}
	fields := make([]map[string]any, 0, len(*targetRTs))
	// iterate through the result types found in the target and add the
	// code to initialize them
	for _, nat := range *targetRTs {
		finit := "new" + targetCtx.Scope.Name(nat.Attribute, "", targetCtx.Pointer, targetCtx.UseDefault)
		if view != "" {
			v := ""
			if vatt := rt.View(view).Find(nat.Name); vatt != nil {
				if attv, ok := vatt.Meta.Last(expr.ViewMetaKey); ok && attv != expr.DefaultView {
					// view is explicitly set for the result type on the attribute
					v = attv
				}
			}
			finit += codegen.Goify(v, true)
		}
		fields = append(fields, map[string]any{
			"VarName":   codegen.Goify(nat.Name, true),
			"FieldInit": finit,
		})
	}
	data["Fields"] = fields

	if err := initTypeCodeTmpl.Execute(&buf, data); err != nil {
		panic(err) // bug
	}
	return buf.String(), helpers
}

// walkViewAttrs iterates through the attributes in att that are found in the
// given view and executes the walker function.
func walkViewAttrs(obj *expr.Object, view *expr.ViewExpr, walker func(name string, attr, vatt *expr.AttributeExpr)) {
	for _, nat := range *expr.AsObject(view.Type) {
		if attr := obj.Attribute(nat.Name); attr != nil {
			walker(nat.Name, attr, nat.Attribute)
		}
	}
}

// removeMeta removes the meta attributes from the given attribute. This is
// needed to make sure that any field name overriding is removed when
// generating protobuf types (as protogen itself won't honor these overrides).
func removeMeta(att *expr.AttributeExpr) {
	_ = codegen.Walk(att, func(a *expr.AttributeExpr) error {
		delete(a.Meta, "struct:pkg:path")
		return nil
	})
}
