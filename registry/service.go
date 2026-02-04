package registry

import (
	"fmt"

	descriptorpb "github.com/golang/protobuf/protoc-gen-go/descriptor" // nolint: depguard
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"

	"github.com/grpc-ecosystem/protoc-gen-grpc-gateway-ts/data"
)

const (
	serviceDescriptorMethodFieldNumber = 2
)

func getHTTPAnnotation(m *descriptorpb.MethodDescriptorProto) *annotations.HttpRule {
	option := proto.GetExtension(m.GetOptions(), annotations.E_Http)
	return option.(*annotations.HttpRule)
}

func hasHTTPAnnotation(m *descriptorpb.MethodDescriptorProto) bool {
	return getHTTPAnnotation(m) != nil
}

// getHTTPMethodPathFromRule extracts HTTP method and path from an HttpRule.
// Returns empty strings for unsupported patterns (e.g., CustomHttpPattern).
func getHTTPMethodPathFromRule(rule *annotations.HttpRule) (method, path string) {
	if rule == nil {
		return "", ""
	}
	switch rule.Pattern.(type) {
	case *annotations.HttpRule_Get:
		return "GET", rule.GetGet()
	case *annotations.HttpRule_Post:
		return "POST", rule.GetPost()
	case *annotations.HttpRule_Put:
		return "PUT", rule.GetPut()
	case *annotations.HttpRule_Patch:
		return "PATCH", rule.GetPatch()
	case *annotations.HttpRule_Delete:
		return "DELETE", rule.GetDelete()
	default:
		panic(fmt.Sprintf("unsupported HTTP method %T", rule.Pattern))
	}
}

// getHTTPBodyFromRule extracts body from an HttpRule
func getHTTPBodyFromRule(rule *annotations.HttpRule) *string {
	if rule == nil {
		return nil
	}
	empty := ""
	switch rule.Pattern.(type) {
	case *annotations.HttpRule_Get:
		return &empty
	default:
		body := rule.GetBody()
		return &body
	}
}

// generateAdditionalBindingMethodName creates unique method name for additional bindings.
// Uses a simple index-based suffix (e.g., MethodName$1, MethodName$2) to avoid
// confusion when the base method name already contains an HTTP method reference.
func generateAdditionalBindingMethodName(baseName string, bindingIndex int) string {
	return fmt.Sprintf("%s$%d", baseName, bindingIndex+1)
}

// createMethodData is a helper to construct a data.Method with the given parameters.
// This reduces duplication between primary method and additional binding creation.
func createMethodData(name, url string, inputType, outputType string, isInputExternal, isOutputExternal bool, serverStreaming, clientStreaming bool, httpMethod string, body *string, comment string) *data.Method {
	return &data.Method{
		Name: name,
		URL:  url,
		Input: &data.MethodArgument{
			Type:       inputType,
			IsExternal: isInputExternal,
		},
		Output: &data.MethodArgument{
			Type:       outputType,
			IsExternal: isOutputExternal,
		},
		ServerStreaming: serverStreaming,
		ClientStreaming: clientStreaming,
		HTTPMethod:      httpMethod,
		HTTPRequestBody: body,
		Comment:         comment,
	}
}

func (r *Registry) analyseService(fileData *data.File, packageName string, fileName string, service *descriptorpb.ServiceDescriptorProto, commentInfo *CommentInfo) {
	packageIdentifier := service.GetName()
	fqName := "." + packageName + "." + packageIdentifier

	// register itself in the registry map
	r.Types[fqName] = &TypeInformation{
		FullyQualifiedName: fqName,
		Package:            packageName,
		File:               fileName,
		PackageIdentifier:  packageIdentifier,
		LocalIdentifier:    service.GetName(),
	}

	serviceData := data.NewService()
	serviceData.Name = service.GetName()
	serviceData.Comment = commentInfo.GetText()
	serviceURLPart := packageName + "." + serviceData.Name

	for i, method := range service.Method {
		// don't support client streaming, will ignore the client streaming method
		if method.GetClientStreaming() {
			continue
		}

		inputTypeFQName := *method.InputType
		isInputTypeExternal := r.isExternalDependenciesOutsidePackage(inputTypeFQName, packageName)

		if isInputTypeExternal {
			fileData.ExternalDependingTypes = append(fileData.ExternalDependingTypes, inputTypeFQName)
		}

		outputTypeFQName := *method.OutputType
		isOutputTypeExternal := r.isExternalDependenciesOutsidePackage(outputTypeFQName, packageName)

		if isOutputTypeExternal {
			fileData.ExternalDependingTypes = append(fileData.ExternalDependingTypes, outputTypeFQName)
		}

		httpMethod := "POST"
		url := "/" + serviceURLPart + "/" + method.GetName()
		var body *string
		var rule *annotations.HttpRule
		if hasHTTPAnnotation(method) {
			rule = getHTTPAnnotation(method)
			if hm, u := getHTTPMethodPathFromRule(rule); hm != "" && u != "" {
				httpMethod = hm
				url = u
			}
			body = getHTTPBodyFromRule(rule)
		}

		methodComment := commentInfo.GetSubComment(serviceDescriptorMethodFieldNumber, i).GetText()
		methodData := createMethodData(
			method.GetName(),
			url,
			inputTypeFQName,
			outputTypeFQName,
			isInputTypeExternal,
			isOutputTypeExternal,
			method.GetServerStreaming(),
			method.GetClientStreaming(),
			httpMethod,
			body,
			methodComment,
		)

		fileData.TrackPackageNonScalarType(methodData.Input)
		fileData.TrackPackageNonScalarType(methodData.Output)

		serviceData.Methods = append(serviceData.Methods, methodData)

		// Process additional bindings
		if rule != nil {
			for bindingIndex, binding := range rule.GetAdditionalBindings() {
				bindingHTTPMethod, bindingURL := getHTTPMethodPathFromRule(binding)
				if bindingHTTPMethod == "" || bindingURL == "" {
					continue
				}
				bindingBody := getHTTPBodyFromRule(binding)
				bindingMethodName := generateAdditionalBindingMethodName(method.GetName(), bindingIndex)

				additionalMethodData := createMethodData(
					bindingMethodName,
					bindingURL,
					inputTypeFQName,
					outputTypeFQName,
					isInputTypeExternal,
					isOutputTypeExternal,
					method.GetServerStreaming(),
					method.GetClientStreaming(),
					bindingHTTPMethod,
					bindingBody,
					methodComment,
				)

				fileData.TrackPackageNonScalarType(additionalMethodData.Input)
				fileData.TrackPackageNonScalarType(additionalMethodData.Output)

				serviceData.Methods = append(serviceData.Methods, additionalMethodData)
			}
		}
	}

	fileData.Services = append(fileData.Services, serviceData)
}
