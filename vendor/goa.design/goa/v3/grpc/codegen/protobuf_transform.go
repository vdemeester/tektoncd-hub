package codegen

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/expr"
)

type transformAttrs struct {
	*codegen.TransformAttrs
	// proto if true indicates that the transformation code is used to initialize
	// a protocol buffer type from a service type. If false, the transformation
	// code is used to initialize a service type from a protocol buffer type.
	proto bool
	// targetInit is the initialization code for the target type for nested
	// map and array types.
	targetInit string
	// wrapped indicates whether the source or target is in a wrapped state.
	// See grpc/docs/FAQ.md. `wrapped` is true and `proto` is true indicates
	// the target attribute is in wrapped state. `wrapped` is true and `proto`
	// is false indicates the source attribute is in wrapped state.
	wrapped bool
	// message is the reference to the Go type generated by protoc for the
	// source message when proto is false.  (protoc builds union struct type
	// names from the parent message name).
	message string
}

// unionData is used by both transformUnion methods.
type unionData struct {
	Source              *expr.Union
	Target              *expr.Union
	SourceValues        []*expr.NamedAttributeExpr
	TargetValues        []*expr.NamedAttributeExpr
	SourceValueTypeRefs []string
}

var (
	// transformGoArrayT is the template to generate Go array transformation
	// code.
	transformGoArrayT *template.Template
	// transformGoMapT is the template to generate Go map transformation code.
	transformGoMapT *template.Template
	// transformGoUnionT is the template to generate Go union transformation
	// code to protobuf.
	transformGoUnionToProtoT *template.Template
	// transformGoUnionT is the template to generate Go union transformation
	// code from protobuf.
	transformGoUnionFromProtoT *template.Template
)

// NOTE: can't initialize inline because https://github.com/golang/go/issues/1817
func init() {
	fm := template.FuncMap{"transformAttribute": transformAttribute, "convertType": convertType}
	transformGoArrayT = template.Must(template.New("transformGoArray").Funcs(fm).Parse(readTemplate("transform_go_array")))
	transformGoMapT = template.Must(template.New("transformGoMap").Funcs(fm).Parse(readTemplate("transform_go_map")))
	transformGoUnionToProtoT = template.Must(template.New("transformGoUnionToProto").Funcs(fm).Parse(readTemplate("transform_go_union_to_proto")))
	transformGoUnionFromProtoT = template.Must(template.New("transformGoUnionFromProto").Funcs(fm).Parse(readTemplate("transform_go_union_from_proto")))
}

// protoBufTransform produces Go code to initialize a data structure defined
// by target from an instance of data structure defined by source. The source
// or target is a protocol buffer type.
//
// source, target are the source and target attributes used in transformation
//
// sourceVar, targetVar are the source and target variables
//
// sourceCtx, targetCtx are the source and target attribute contexts
//
// `proto` param if true indicates that the target is a protocol buffer type
//
// newVar if true initializes a target variable with the generated Go code
// using `:=` operator. If false, it assigns Go code to the target variable
// using `=`.
func protoBufTransform(source, target *expr.AttributeExpr, sourceVar, targetVar string, sourceCtx, targetCtx *codegen.AttributeContext, proto, newVar bool) (string, []*codegen.TransformFunctionData, error) {
	ta := &transformAttrs{TransformAttrs: &codegen.TransformAttrs{
		SourceCtx: sourceCtx,
		TargetCtx: targetCtx,
	}}
	if proto {
		target = expr.DupAtt(target)
		removeMeta(target)
		ta.Prefix = "svc"
		ta.proto = true
	} else {
		source = expr.DupAtt(source)
		removeMeta(source)
		ta.Prefix = "protobuf"
		ta.proto = false
	}

	code, err := transformAttribute(source, target, sourceVar, targetVar, newVar, ta)
	if err != nil {
		return "", nil, err
	}

	funcs, err := transformAttributeHelpers(source, target, ta, make(map[string]*codegen.TransformFunctionData))
	if err != nil {
		return "", nil, err
	}

	return strings.TrimRight(code, "\n"), funcs, nil
}

// removeMeta removes meta attributes from the given attribute that cannot be
// honored. This is needed to make sure that any field name overridding is
// removed when generating protobuf types (as protogen itself won't honor these
// overrides).
func removeMeta(att *expr.AttributeExpr) {
	_ = codegen.Walk(att, func(a *expr.AttributeExpr) error {
		delete(a.Meta, "struct:field:name")
		delete(a.Meta, "struct:field:external")
		delete(a.Meta, "struct.field.external") // Deprecated syntax. Only present for backward compatibility.
		return nil
	})
}

// transformAttribute returns the code to initialize a target data structure
// from an instance of source data structure. It returns an error if source and
// target are not compatible for transformation (different types, fields of
// different type).
func transformAttribute(source, target *expr.AttributeExpr, sourceVar, targetVar string, newVar bool, ta *transformAttrs) (string, error) {
	var (
		initCode string
		err      error
	)

	if err := codegen.IsCompatible(source.Type, target.Type, sourceVar, targetVar); err != nil {
		if ta.proto {
			name := ta.TargetCtx.Scope.Name(target, ta.TargetCtx.Pkg(target), ta.TargetCtx.Pointer, ta.TargetCtx.UseDefault)
			initCode += fmt.Sprintf("%s := &%s{}\n", targetVar, name)
			targetVar += ".Field"
			newVar = false
			target = unwrapAttr(expr.DupAtt(target))
		} else {
			source = unwrapAttr(expr.DupAtt(source))
			sourceVar += ".Field"
		}
		if err = codegen.IsCompatible(source.Type, target.Type, sourceVar, targetVar); err != nil {
			return "", err
		}
	}

	if ta.proto {
		if isUnionMessage(target) {
			ta = dupTransformAttrs(ta)
			ta.message = ta.TargetCtx.Scope.Name(target, ta.TargetCtx.Pkg(target), false, false)
		}
	} else {
		if isUnionMessage(source) {
			ta = dupTransformAttrs(ta)
			ta.message = ta.SourceCtx.Scope.Ref(source, ta.SourceCtx.Pkg(source))
		}
	}

	var code string
	{
		switch {
		case expr.IsArray(source.Type):
			code, err = transformArray(expr.AsArray(source.Type), expr.AsArray(target.Type), sourceVar, targetVar, newVar, ta)
		case expr.IsMap(source.Type):
			code, err = transformMap(expr.AsMap(source.Type), expr.AsMap(target.Type), sourceVar, targetVar, newVar, ta)
		case expr.IsObject(source.Type):
			code, err = transformObject(source, target, sourceVar, targetVar, newVar, ta)
		case expr.IsUnion(source.Type):
			if ta.proto {
				code, err = transformUnionToProto(source, target, sourceVar, targetVar, ta)
			} else {
				code, err = transformUnionFromProto(source, target, sourceVar, targetVar, ta)
			}
		default:
			assign := "="
			if newVar {
				assign = ":="
			}
			srcField := convertType(source, target, false, false, sourceVar, ta)
			code = fmt.Sprintf("%s %s %s\n", targetVar, assign, srcField)
		}
	}
	if err != nil {
		return "", err
	}
	return initCode + code, nil
}

// transformObject returns the code to transform source attribute of object
// type to target attribute of object type. It returns an error if source
// and target are not compatible for transformation.
func transformObject(source, target *expr.AttributeExpr, sourceVar, targetVar string, newVar bool, ta *transformAttrs) (string, error) {
	var (
		initCode     string
		postInitCode string
	)
	{
		// iterate through primitive attributes to initialize the struct
		walkMatches(source, target, func(srcMatt, tgtMatt *expr.MappedAttributeExpr, srcc, tgtc *expr.AttributeExpr, n string) {
			if !expr.IsPrimitive(srcc.Type) {
				return
			}
			var (
				exp          string
				srcField     = sourceVar + "." + ta.SourceCtx.Scope.Field(srcc, srcMatt.ElemName(n), true)
				tgtField     = ta.TargetCtx.Scope.Field(tgtc, tgtMatt.ElemName(n), true)
				srcPtr       = ta.SourceCtx.IsPrimitivePointer(n, srcMatt.AttributeExpr)
				tgtPtr       = ta.TargetCtx.IsPrimitivePointer(n, tgtMatt.AttributeExpr)
				srcFieldConv = convertType(srcc, tgtc, srcPtr, tgtPtr, srcField, ta)
				_, isSrcUT   = srcc.Type.(expr.UserType)
				_, isTgtUT   = tgtc.Type.(expr.UserType)
			)
			switch {
			case isSrcUT || isTgtUT || (srcField != srcFieldConv):
				var deref string
				if srcPtr {
					deref = "*"
				}
				exp = srcFieldConv
				if isSrcUT && !ta.proto {
					// If the source is an alias type and the code is initializing a service
					// type then we must cast to the alias type.
					exp = fmt.Sprintf("%s(%s%s)", ta.TargetCtx.Scope.Ref(tgtc, ta.TargetCtx.Pkg(tgtc)), deref, srcField)
				}
				if srcPtr && !srcMatt.IsRequired(n) {
					postInitCode += fmt.Sprintf("if %s != nil {\n", srcField)
					if tgtPtr {
						tmp := codegen.Goify(tgtMatt.ElemName(n), false)
						postInitCode += fmt.Sprintf("%s := %s\n%s.%s = &%s\n", tmp, exp, targetVar, tgtField, tmp)
					} else {
						postInitCode += fmt.Sprintf("%s.%s = %s\n", targetVar, tgtField, exp)
					}
					postInitCode += "}\n"
					return
				} else if tgtPtr {
					tmp := codegen.Goify(tgtMatt.ElemName(n), false)
					postInitCode += fmt.Sprintf("%s := %s\n%s.%s = &%s\n", tmp, exp, targetVar, tgtField, tmp)
					return
				}
			case srcPtr && !tgtPtr:
				exp = "*" + srcField
				if !srcMatt.IsRequired(n) {
					postInitCode += fmt.Sprintf("if %s != nil {\n\t%s.%s = %s\n}\n", srcField, targetVar, tgtField, exp)
					return
				}
			case !srcPtr && tgtPtr:
				exp = "&" + srcField
			default:
				exp = srcField
			}
			initCode += fmt.Sprintf("\n%s: %s,", tgtField, exp)
		})
		if initCode != "" {
			initCode += "\n"
		}
	}

	buffer := &bytes.Buffer{}
	deref := "&"
	// if the target is a raw struct no need to return a pointer
	if _, ok := target.Type.(*expr.Object); ok {
		deref = ""
	}
	assign := "="
	if newVar {
		assign = ":="
	}
	tname := ta.TargetCtx.Scope.Name(target, ta.TargetCtx.Pkg(target), ta.TargetCtx.Pointer, ta.TargetCtx.UseDefault)
	fmt.Fprintf(buffer, "%s %s %s%s{%s}\n", targetVar, assign, deref, tname, initCode)
	buffer.WriteString(postInitCode)

	// iterate through attributes to initialize rest of the struct fields and
	// handle default values
	var err error
	walkMatches(source, target, func(srcMatt, tgtMatt *expr.MappedAttributeExpr, srcc, tgtc *expr.AttributeExpr, n string) {
		srcc = unAlias(srcc)
		tgtc = unAlias(tgtc)
		var (
			code string

			srcVar = sourceVar + "." + ta.SourceCtx.Scope.Field(srcc, srcMatt.ElemName(n), true)
			tgtVar = targetVar + "." + ta.TargetCtx.Scope.Field(tgtc, tgtMatt.ElemName(n), true)
		)
		{
			if err = codegen.IsCompatible(srcc.Type, tgtc.Type, "", ""); err != nil {
				if ta.proto {
					ta.targetInit = ta.TargetCtx.Scope.Name(tgtc, ta.TargetCtx.Pkg(tgtc), ta.TargetCtx.Pointer, ta.TargetCtx.UseDefault)
					tgtc = unwrapAttr(tgtc)
				} else {
					srcc = unwrapAttr(srcc)
				}
				ta.wrapped = true
				if err = codegen.IsCompatible(srcc.Type, tgtc.Type, "", ""); err != nil {
					return
				}
			}
			_, isUserType := srcc.Type.(expr.UserType)
			switch {
			case expr.IsArray(srcc.Type):
				code, err = transformArray(expr.AsArray(srcc.Type), expr.AsArray(tgtc.Type), srcVar, tgtVar, false, ta)
			case expr.IsMap(srcc.Type):
				code, err = transformMap(expr.AsMap(srcc.Type), expr.AsMap(tgtc.Type), srcVar, tgtVar, false, ta)
			case isUserType:
				if ta.TargetCtx.IsInterface {
					ref := ta.TargetCtx.Scope.Ref(target, ta.TargetCtx.Pkg(target))
					tgtVar = targetVar + ".(" + ref + ")." + codegen.GoifyAtt(tgtc, tgtMatt.ElemName(n), true)
				}
				if !expr.IsPrimitive(srcc.Type) {
					code = fmt.Sprintf("%s = %s(%s)\n", tgtVar, transformHelperName(srcc, tgtc, ta), srcVar)
				}
			case expr.IsObject(srcc.Type):
				code, err = transformAttribute(srcc, tgtc, srcVar, tgtVar, false, ta)
			case expr.IsUnion(srcc.Type):
				code, err = transformAttribute(srcc, tgtc, srcVar, tgtVar, false, ta)
			}
		}
		if err != nil {
			return
		}

		// We need to check for a nil source if it holds a reference (pointer to
		// primitive or an object, array or map) and is not required. We also want
		// to always check nil if the attribute is not a primitive; it's a
		// 1) user type and we want to avoid calling transform helper functions
		// with nil value
		// 2) it's an object, map or array to avoid making empty arrays and maps
		// and to avoid derefencing nil.
		var checkNil bool
		{
			isRef := !expr.IsPrimitive(srcc.Type) && !srcMatt.IsRequired(n) || ta.SourceCtx.IsPrimitivePointer(n, srcMatt.AttributeExpr) && expr.IsPrimitive(srcc.Type)
			marshalNonPrimitive := !expr.IsPrimitive(srcc.Type)
			checkNil = isRef || marshalNonPrimitive
		}
		if code != "" && checkNil {
			code = fmt.Sprintf("if %s != nil {\n\t%s}\n", srcVar, code)
		}

		// Default value handling. We need to handle default values if the target
		// type uses default values (i.e. attributes with default values are
		// non-pointers) and has a default value set.
		if tdef := tgtMatt.GetDefault(n); tdef != nil && ta.TargetCtx.UseDefault && !ta.TargetCtx.Pointer && !srcMatt.IsRequired(n) {
			switch {
			case ta.SourceCtx.IsPrimitivePointer(n, srcMatt.AttributeExpr) || !expr.IsPrimitive(srcc.Type):
				// source attribute is a primitive pointer or not a primitive
				code += fmt.Sprintf("if %s == nil {\n\t", srcVar)
				if ta.TargetCtx.IsPrimitivePointer(n, tgtMatt.AttributeExpr) && expr.IsPrimitive(tgtc.Type) {
					nativeTypeName := codegen.GoNativeTypeName(tgtc.Type)
					if ta.proto {
						nativeTypeName = protoBufNativeGoTypeName(tgtc.Type)
					}
					code += fmt.Sprintf("var tmp %s = %#v\n\t%s = &tmp\n", nativeTypeName, tdef, tgtVar)
				} else {
					code += fmt.Sprintf("%s = %#v\n", tgtVar, tdef)
				}
				code += "}\n"
			case expr.IsPrimitive(srcc.Type) && srcMatt.HasDefaultValue(n) && ta.SourceCtx.UseDefault:
				// source attribute is a primitive with default value
				// (the field is not a pointer in this case)
				code += "{\n\t"
				if ta.proto {
					code += fmt.Sprintf("var zero %s\n\t", protoBufNativeGoTypeName(tgtc.Type))
				} else {
					if typeName, _ := codegen.GetMetaType(tgtc); typeName != "" {
						code += fmt.Sprintf("var zero %s\n\t", typeName)
					} else if _, ok := tgtc.Type.(expr.UserType); ok {
						// aliased primitive
						code += fmt.Sprintf("var zero %s\n\t", ta.TargetCtx.Scope.Ref(tgtc, ta.TargetCtx.Pkg(tgtc)))
					} else {
						code += fmt.Sprintf("var zero %s\n\t", codegen.GoNativeTypeName(tgtc.Type))
					}
				}
				code += fmt.Sprintf("if %s == zero {\n\t%s = %#v\n}\n", tgtVar, tgtVar, tdef)
				code += "}\n"
			}
		}
		buffer.WriteString(code)
	})
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

// transformArray returns the code to transform source attribute of array
// type to target attribute of array type. It returns an error if source
// and target are not compatible for transformation.
func transformArray(source, target *expr.Array, sourceVar, targetVar string, newVar bool, ta *transformAttrs) (string, error) {
	elem := target.ElemType
	if ta.proto {
		elem = unAlias(elem)
	}
	targetRef := ta.TargetCtx.Scope.Ref(elem, ta.TargetCtx.Pkg(elem))

	var (
		code string
		err  error
	)

	// If targetInit is set, the target array element is in a nested state.
	// See grpc/docs/FAQ.md.
	if ta.targetInit != "" {
		assign := "="
		if newVar {
			assign = ":="
		}
		code = fmt.Sprintf("%s %s &%s{}\n", targetVar, assign, ta.targetInit)
		ta.targetInit = ""
	}
	if ta.wrapped {
		if ta.proto {
			targetVar += ".Field"
			newVar = false
		} else {
			sourceVar += ".Field"
		}
		ta.wrapped = false
	}

	src := source.ElemType
	tgt := elem
	if err = codegen.IsCompatible(src.Type, tgt.Type, "[0]", "[0]"); err != nil {
		if ta.proto {
			ta.targetInit = ta.TargetCtx.Scope.Name(tgt, ta.TargetCtx.Pkg(tgt), ta.TargetCtx.Pointer, ta.TargetCtx.UseDefault)
			tgt = unwrapAttr(expr.DupAtt(tgt))
		} else {
			src = unwrapAttr(expr.DupAtt(src))
		}
		ta.wrapped = true
		if err = codegen.IsCompatible(src.Type, tgt.Type, "[0]", "[0]"); err != nil {
			return "", err
		}
	}
	valVar := "val"
	if obj := expr.AsObject(src.Type); obj != nil {
		if len(*obj) == 0 {
			valVar = ""
		}
	}

	data := map[string]any{
		"ElemTypeRef":    targetRef,
		"SourceElem":     src,
		"TargetElem":     tgt,
		"SourceVar":      sourceVar,
		"TargetVar":      targetVar,
		"NewVar":         newVar,
		"TransformAttrs": ta,
		"LoopVar":        string(rune(105 + strings.Count(targetVar, "["))),
		"ValVar":         valVar,
	}
	var buf bytes.Buffer
	if err := transformGoArrayT.Execute(&buf, data); err != nil {
		return "", err
	}
	return code + buf.String(), nil
}

// transformMap returns the code to transform source attribute of map
// type to target attribute of map type. It returns an error if source
// and target are not compatible for transformation.
func transformMap(source, target *expr.Map, sourceVar, targetVar string, newVar bool, ta *transformAttrs) (string, error) {
	// Target map key cannot be nested in protocol buffers. So no need to worry
	// about unwrapping.
	if err := codegen.IsCompatible(source.KeyType.Type, target.KeyType.Type, sourceVar+"[key]", targetVar+"[key]"); err != nil {
		return "", err
	}
	kt := target.KeyType
	if ta.proto {
		kt = unAlias(kt)
	}
	et := target.ElemType
	if ta.proto {
		et = unAlias(et)
	}
	targetKeyRef := ta.TargetCtx.Scope.Ref(kt, ta.TargetCtx.Pkg(kt))
	targetElemRef := ta.TargetCtx.Scope.Ref(et, ta.TargetCtx.Pkg(et))

	var (
		code string
		err  error
	)

	// If targetInit is set, the target map element is in a nested state.
	// See grpc/docs/FAQ.md.
	if ta.targetInit != "" {
		assign := "="
		if newVar {
			assign = ":="
		}
		code = fmt.Sprintf("%s %s &%s{}\n", targetVar, assign, ta.targetInit)
		ta.targetInit = ""
	}
	if ta.wrapped {
		if ta.proto {
			targetVar += ".Field"
			newVar = false
		} else {
			sourceVar += ".Field"
		}
		ta.wrapped = false
	}

	src := source.ElemType
	tgt := target.ElemType
	if err = codegen.IsCompatible(src.Type, tgt.Type, "[*]", "[*]"); err != nil {
		if ta.proto {
			ta.targetInit = ta.TargetCtx.Scope.Name(tgt, ta.TargetCtx.Pkg(tgt), ta.TargetCtx.Pointer, ta.TargetCtx.UseDefault)
			tgt = unwrapAttr(expr.DupAtt(tgt))
		} else {
			src = unwrapAttr(expr.DupAtt(src))
		}
		ta.wrapped = true
		if err = codegen.IsCompatible(src.Type, tgt.Type, "[*]", "[*]"); err != nil {
			return "", err
		}
	}
	data := map[string]any{
		"KeyTypeRef":     targetKeyRef,
		"ElemTypeRef":    targetElemRef,
		"SourceKey":      source.KeyType,
		"TargetKey":      target.KeyType,
		"SourceElem":     src,
		"TargetElem":     tgt,
		"SourceVar":      sourceVar,
		"TargetVar":      targetVar,
		"NewVar":         newVar,
		"TransformAttrs": ta,
		"LoopVar":        "",
	}
	if depth := codegen.MapDepth(target); depth > 0 {
		data["LoopVar"] = string(rune(97 + depth))
	}
	var buf bytes.Buffer
	if err := transformGoMapT.Execute(&buf, data); err != nil {
		return "", err
	}
	return code + buf.String(), nil
}

// transformUnionToProto returns the code to transform an attribute of type
// union from Goa to protobuf. It returns an error if source and target are not
// compatible for transformation.
func transformUnionToProto(source, target *expr.AttributeExpr, sourceVar, targetVar string, ta *transformAttrs) (string, error) {
	if err := codegen.IsCompatible(source.Type, target.Type, sourceVar, targetVar); err != nil {
		return "", err
	}
	tdata := transformUnionData(source, target, ta)
	targetValueTypeNames := make([]string, len(tdata.TargetValues))
	targetFieldNames := make([]string, len(tdata.TargetValues))
	for i, v := range tdata.TargetValues {
		fieldName := ta.TargetCtx.Scope.Field(v.Attribute, v.Name, true)
		targetFieldNames[i] = fieldName
		targetValueTypeNames[i] = ta.message + "_" + fieldName
	}

	data := map[string]any{
		"SourceVar":            sourceVar,
		"TargetVar":            targetVar,
		"SourceValueTypeRefs":  tdata.SourceValueTypeRefs,
		"SourceValues":         tdata.SourceValues,
		"TargetValues":         tdata.TargetValues,
		"TransformAttrs":       ta,
		"TargetValueTypeNames": targetValueTypeNames,
		"TargetFieldNames":     targetFieldNames,
	}
	var buf bytes.Buffer
	if err := transformGoUnionToProtoT.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// transformUnionFromProto returns the code to transform an attribute of type
// union from Goa to protobuf. It returns an error if source and target are not
// compatible for transformation.
func transformUnionFromProto(source, target *expr.AttributeExpr, sourceVar, targetVar string, ta *transformAttrs) (string, error) {
	if err := codegen.IsCompatible(source.Type, target.Type, sourceVar, targetVar); err != nil {
		return "", err
	}
	tdata := transformUnionData(source, target, ta)
	sourceFieldNames := make([]string, len(tdata.SourceValues))
	for i, v := range tdata.SourceValues {
		fieldName := ta.SourceCtx.Scope.Field(v.Attribute, v.Name, true)
		sourceFieldNames[i] = fieldName
	}

	data := map[string]any{
		"SourceVar":           sourceVar,
		"TargetVar":           targetVar,
		"SourceValueTypeRefs": tdata.SourceValueTypeRefs,
		"SourceFieldNames":    sourceFieldNames,
		"SourceValues":        tdata.SourceValues,
		"TargetValues":        tdata.TargetValues,
		"TransformAttrs":      ta,
	}
	var buf bytes.Buffer
	if err := transformGoUnionFromProtoT.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// convertType produces code to initialize a target type from a source type
// held by sourceVar.
func convertType(src, tgt *expr.AttributeExpr, srcPtr bool, tgtPtr bool, srcVar string, ta *transformAttrs) string {
	if expr.IsAlias(src.Type) || expr.IsAlias(tgt.Type) {
		srcp, tgtp := unAlias(src), unAlias(tgt)
		if srcp.Type == tgtp.Type {
			if ta.proto {
				return convertPrimitiveToProto(src, tgtp, srcPtr, tgtPtr, srcVar, ta)
			}
			return convertPrimitiveFromProto(srcp, tgt, srcPtr, tgtPtr, srcVar, ta)
		}
		return fmt.Sprintf("%s(%s)", transformHelperName(src, tgt, ta), srcVar)
	}

	if _, ok := src.Type.(expr.UserType); ok {
		return fmt.Sprintf("%s(%s)", transformHelperName(src, tgt, ta), srcVar)
	}

	srcType, _ := codegen.GetMetaType(src)
	tgtType, _ := codegen.GetMetaType(tgt)
	if srcType == "" && tgtType == "" && (src.Type != expr.Int) && (src.Type != expr.UInt) {
		// Nothing to do
		return srcVar
	}

	if ta.proto {
		return convertPrimitiveToProto(src, tgt, srcPtr, tgtPtr, srcVar, ta)
	}
	return convertPrimitiveFromProto(src, tgt, srcPtr, tgtPtr, srcVar, ta)
}

// convertPrimitive returns the code to convert a primitive type from one
// representation to another.
// NOTE: For Int and UInt kinds, protocol buffer Go compiler generates
// int32 and uint32 respectively whereas Goa generates int and uint.
func convertPrimitiveToProto(_, tgt *expr.AttributeExpr, srcPtr, _ bool, srcVar string, _ *transformAttrs) string {
	tgtType := protoBufNativeGoTypeName(tgt.Type)
	if srcPtr {
		srcVar = "*" + srcVar
	}
	return fmt.Sprintf("%s(%s)", tgtType, srcVar)
}

func convertPrimitiveFromProto(_, tgt *expr.AttributeExpr, srcPtr, _ bool, srcVar string, ta *transformAttrs) string {
	tgtType, _ := codegen.GetMetaType(tgt)
	if tgtType == "" {
		tgtType = ta.TargetCtx.Scope.Ref(tgt, ta.TargetCtx.Pkg(tgt))
	}
	if srcPtr {
		srcVar = "*" + srcVar
	}
	return fmt.Sprintf("%s(%s)", tgtType, srcVar)
}

// transformUnionData returns data needed by both transformUnion functions.
func transformUnionData(source, target *expr.AttributeExpr, ta *transformAttrs) *unionData {
	src := expr.AsUnion(source.Type)
	tgt := expr.AsUnion(target.Type)
	srcValues := make([]*expr.NamedAttributeExpr, len(src.Values))
	copy(srcValues, src.Values)
	tgtValues := make([]*expr.NamedAttributeExpr, len(tgt.Values))
	copy(tgtValues, tgt.Values)

	sourceValueTypeRefs := make([]string, len(src.Values))
	for i, v := range src.Values {
		if ta.proto {
			sourceValueTypeRefs[i] = ta.SourceCtx.Scope.Ref(v.Attribute, ta.SourceCtx.Pkg(v.Attribute))
		} else {
			fieldName := ta.SourceCtx.Scope.Field(v.Attribute, v.Name, true)
			sourceValueTypeRefs[i] = ta.message + "_" + fieldName
		}
	}
	return &unionData{
		Source:              src,
		Target:              tgt,
		SourceValues:        srcValues,
		TargetValues:        tgtValues,
		SourceValueTypeRefs: sourceValueTypeRefs,
	}
}

// transformAttributeHelpers returns the Go transform functions and their definitions
// that may be used in code produced by Transform. It returns an error if source and
// target are incompatible (different types, fields of different type etc).
//
// source, target are the source and target attributes used in transformation
//
// ta is the transform attributes
//
// seen keeps track of generated transform functions to avoid recursion
func transformAttributeHelpers(source, target *expr.AttributeExpr, ta *transformAttrs, seen map[string]*codegen.TransformFunctionData) ([]*codegen.TransformFunctionData, error) {
	var (
		helpers []*codegen.TransformFunctionData
		err     error
	)
	{
		if err = codegen.IsCompatible(source.Type, target.Type, "", ""); err != nil {
			if ta.proto {
				target = unwrapAttr(expr.DupAtt(target))
			} else {
				source = unwrapAttr(expr.DupAtt(source))
			}
			if err = codegen.IsCompatible(source.Type, target.Type, "", ""); err != nil {
				return nil, err
			}
		}
		// Do not generate a transform function for the top most user type.
		switch {
		case expr.IsArray(source.Type):
			source = expr.AsArray(source.Type).ElemType
			target = expr.AsArray(target.Type).ElemType
			helpers, err = transformAttributeHelpers(source, target, ta, seen)
		case expr.IsMap(source.Type):
			sm := expr.AsMap(source.Type)
			tm := expr.AsMap(target.Type)
			helpers, err = transformAttributeHelpers(sm.ElemType, tm.ElemType, ta, seen)
			if err == nil {
				var other []*codegen.TransformFunctionData
				other, err = transformAttributeHelpers(sm.KeyType, tm.KeyType, ta, seen)
				helpers = append(helpers, other...)
			}
		case expr.IsUnion(source.Type):
			srcAttrs := expr.AsUnion(source.Type)
			tgtAttrs := expr.AsUnion(target.Type)
			if len(srcAttrs.Values) != len(tgtAttrs.Values) {
				return nil, fmt.Errorf("cannot transform union attribute %s with %d types to union attribute %s with %d types",
					source.Type.Name(), len(srcAttrs.Values), target.Type.Name(), len(tgtAttrs.Values))
			}
			for i, srcAtt := range srcAttrs.Values {
				tgtAtt := tgtAttrs.Values[i]
				h, err := collectHelpers(srcAtt.Attribute, tgtAtt.Attribute, true, ta, seen)
				if err == nil {
					helpers = append(helpers, h...)
				}
			}
		case expr.IsObject(source.Type):
			walkMatches(source, target, func(srcMatt, _ *expr.MappedAttributeExpr, srcc, tgtc *expr.AttributeExpr, n string) {
				if err != nil {
					return
				}
				if err = codegen.IsCompatible(srcc.Type, tgtc.Type, "", ""); err != nil {
					if ta.proto {
						tgtc = unwrapAttr(tgtc)
					} else {
						srcc = unwrapAttr(srcc)
					}
					if err = codegen.IsCompatible(srcc.Type, tgtc.Type, "", ""); err != nil {
						return
					}
				}
				h, err2 := collectHelpers(srcc, tgtc, srcMatt.IsRequired(n), ta, seen)
				if err2 != nil {
					err = err2
					return
				}
				helpers = append(helpers, h...)
			})
		}
	}
	if err != nil {
		return nil, err
	}
	return helpers, nil
}

// collectHelpers recursively traverses the given attributes and return the
// transform helper functions required to generate the transform code.
func collectHelpers(source, target *expr.AttributeExpr, req bool, ta *transformAttrs, seen map[string]*codegen.TransformFunctionData) ([]*codegen.TransformFunctionData, error) {
	var data []*codegen.TransformFunctionData
	switch {
	case expr.IsArray(source.Type):
		helpers, err := transformAttributeHelpers(
			expr.AsArray(source.Type).ElemType,
			expr.AsArray(target.Type).ElemType,
			ta, seen)
		if err != nil {
			return nil, err
		}
		data = append(data, helpers...)
	case expr.IsMap(source.Type):
		helpers, err := transformAttributeHelpers(
			expr.AsMap(source.Type).KeyType,
			expr.AsMap(target.Type).KeyType,
			ta, seen)
		if err != nil {
			return nil, err
		}
		data = append(data, helpers...)
		helpers, err = transformAttributeHelpers(
			expr.AsMap(source.Type).ElemType,
			expr.AsMap(target.Type).ElemType,
			ta, seen)
		if err != nil {
			return nil, err
		}
		data = append(data, helpers...)
	case expr.IsUnion(source.Type):
		srcAttrs := expr.AsUnion(source.Type)
		tgtAttrs := expr.AsUnion(target.Type)
		if len(srcAttrs.Values) != len(tgtAttrs.Values) {
			return nil, fmt.Errorf("cannot transform union attribute %s with %d types to union attribute %s with %d types",
				source.Type.Name(), len(srcAttrs.Values), target.Type.Name(), len(tgtAttrs.Values))
		}
		for i, srcVal := range srcAttrs.Values {
			src := srcVal.Attribute
			tgt := tgtAttrs.Values[i].Attribute
			if ta.proto {
				tgt = unwrapAttr(tgt)
			} else {
				src = unwrapAttr(src)
			}
			helpers, err := collectHelpers(src, tgt, true, ta, seen)
			if err != nil {
				return nil, err
			}
			data = append(data, helpers...)
		}
	case expr.IsObject(source.Type):
		if ut, ok := source.Type.(expr.UserType); ok {
			if ta.proto {
				if isUnionMessage(target) {
					ta = dupTransformAttrs(ta)
					ta.message = ta.TargetCtx.Scope.Name(target, ta.TargetCtx.Pkg(target), false, false)
				}
			} else {
				if isUnionMessage(source) {
					ta = dupTransformAttrs(ta)
					ta.message = ta.SourceCtx.Scope.Ref(source, ta.SourceCtx.Pkg(source))
				}
			}
			name := transformHelperName(source, target, ta)
			if _, ok := seen[name]; ok {
				return nil, nil
			}
			code, err := transformAttribute(ut.Attribute(), target, "v", "res", true, ta)
			if err != nil {
				return nil, err
			}
			if !req {
				code = "if v == nil {\n\treturn nil\n}\n" + code
			}
			tfd := &codegen.TransformFunctionData{
				Name:          name,
				ParamTypeRef:  ta.SourceCtx.Scope.Ref(source, ta.SourceCtx.Pkg(source)),
				ResultTypeRef: ta.TargetCtx.Scope.Ref(target, ta.TargetCtx.Pkg(target)),
				Code:          code,
			}
			seen[name] = tfd
			data = append(data, tfd)
		}

		// collect helpers
		var err error
		{
			walkMatches(source, target, func(srcMatt, _ *expr.MappedAttributeExpr, srcc, tgtc *expr.AttributeExpr, n string) {
				if err = codegen.IsCompatible(srcc.Type, tgtc.Type, "", ""); err != nil {
					if ta.proto {
						tgtc = unwrapAttr(tgtc)
					} else {
						srcc = unwrapAttr(srcc)
					}
					if err = codegen.IsCompatible(srcc.Type, tgtc.Type, "", ""); err != nil {
						return
					}
				}
				var helpers []*codegen.TransformFunctionData
				helpers, err = collectHelpers(srcc, tgtc, srcMatt.IsRequired(n), ta, seen)
				if err != nil {
					return
				}
				data = append(data, helpers...)
			})
		}
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

// walkMatches iterates through the source attribute expression and executes
// the walker function.
func walkMatches(source, target *expr.AttributeExpr, walker func(src, tgt *expr.MappedAttributeExpr, srcc, tgtc *expr.AttributeExpr, n string)) {
	srcMatt := expr.NewMappedAttributeExpr(source)
	tgtMatt := expr.NewMappedAttributeExpr(target)
	srcObj := expr.AsObject(srcMatt.Type)
	tgtObj := expr.AsObject(tgtMatt.Type)
	for _, nat := range *srcObj {
		if att := tgtObj.Attribute(nat.Name); att != nil {
			walker(srcMatt, tgtMatt, nat.Attribute, att, nat.Name)
		}
	}
}

// transformHelperName returns the transformation function name to initialize a
// target user type from an instance of a source user type.
func transformHelperName(source, target *expr.AttributeExpr, ta *transformAttrs) string {
	var (
		sname  string
		tname  string
		prefix string
	)
	{
		// Do not consider package overrides for protogen generated types
		if ta.proto {
			target = expr.DupAtt(target)
			codegen.Walk(target, func(att *expr.AttributeExpr) error { // nolint: errcheck
				delete(att.Meta, "struct:pkg:path")
				return nil
			})
		} else {
			source = expr.DupAtt(source)
			codegen.Walk(source, func(att *expr.AttributeExpr) error { // nolint: errcheck
				delete(att.Meta, "struct:pkg:path")
				return nil
			})
		}
		sname = codegen.Goify(ta.SourceCtx.Scope.Name(source, ta.SourceCtx.Pkg(source), ta.TargetCtx.Pointer, ta.TargetCtx.UseDefault), true)
		tname = codegen.Goify(ta.TargetCtx.Scope.Name(target, ta.TargetCtx.Pkg(target), ta.TargetCtx.Pointer, ta.TargetCtx.UseDefault), true)
		prefix = ta.Prefix
	}
	return codegen.Goify(prefix+sname+"To"+tname, false)
}

// unAlias returns the base AttributeExpr of an aliased one.
func unAlias(at *expr.AttributeExpr) *expr.AttributeExpr {
	if prim := getPrimitive(at); prim != nil {
		return prim
	}
	return at
}

// isUnionMessage returns true if the given attribute is a union message.
func isUnionMessage(at *expr.AttributeExpr) bool {
	ut, ok := at.Type.(expr.UserType)
	if !ok {
		return false
	}
	obj := expr.AsObject(ut.Attribute().Type)
	if obj == nil {
		return false
	}
	for _, nat := range *obj {
		if expr.IsUnion(nat.Attribute.Type) {
			return true
		}
	}
	return false
}

// dupTransformAttrs returns a shallow copy of the given transformAttrs.
func dupTransformAttrs(ta *transformAttrs) *transformAttrs {
	return &transformAttrs{
		TransformAttrs: ta.TransformAttrs,
		proto:          ta.proto,
		targetInit:     ta.targetInit,
		wrapped:        ta.wrapped,
		message:        ta.message,
	}
}
