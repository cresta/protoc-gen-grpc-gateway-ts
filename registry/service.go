package registry

import (
	descriptorpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	log "github.com/sirupsen/logrus" // nolint: depguard
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"

	"github.com/grpc-ecosystem/protoc-gen-grpc-gateway-ts/data"
)

const (
	serviceDescriptorMethodFieldNumber = 2
)

type httpBinding struct {
	HTTPMethod string
	URL        string
	Body       *string
	Index      int
}

func getHTTPAnnotation(m *descriptorpb.MethodDescriptorProto) *annotations.HttpRule {
	option := proto.GetExtension(m.GetOptions(), annotations.E_Http)
	return option.(*annotations.HttpRule)
}

func hasHTTPAnnotation(m *descriptorpb.MethodDescriptorProto) bool {
	return getHTTPAnnotation(m) != nil
}

// extractBindingFromRule extracts HTTP method, path, and body from an HttpRule.
// Returns nil for unsupported patterns (e.g., CustomHttpPattern).
func extractBindingFromRule(rule *annotations.HttpRule, index int) *httpBinding {
	if rule == nil {
		return nil
	}

	var method, path string
	switch rule.Pattern.(type) {
	case *annotations.HttpRule_Get:
		method, path = "GET", rule.GetGet()
	case *annotations.HttpRule_Post:
		method, path = "POST", rule.GetPost()
	case *annotations.HttpRule_Put:
		method, path = "PUT", rule.GetPut()
	case *annotations.HttpRule_Patch:
		method, path = "PATCH", rule.GetPatch()
	case *annotations.HttpRule_Delete:
		method, path = "DELETE", rule.GetDelete()
	default:
		log.Warnf("unsupported HTTP pattern type %T, skipping", rule.Pattern)
		return nil
	}

	if method == "" || path == "" {
		return nil
	}

	// GET requests have empty body, others use the rule's body field
	var body *string
	empty := ""
	if method == "GET" {
		body = &empty
	} else {
		b := rule.GetBody()
		body = &b
	}

	return &httpBinding{
		HTTPMethod: method,
		URL:        path,
		Body:       body,
		Index:      index,
	}
}

// collectHTTPBindings returns all HTTP bindings (primary + additional) from an HttpRule.
// The primary binding has Index=0, additional bindings have Index=1, 2, etc.
func collectHTTPBindings(rule *annotations.HttpRule) []*httpBinding {
	var bindings []*httpBinding

	if primary := extractBindingFromRule(rule, 0); primary != nil {
		bindings = append(bindings, primary)
	}

	for i, additionalRule := range rule.GetAdditionalBindings() {
		if binding := extractBindingFromRule(additionalRule, i+1); binding != nil {
			bindings = append(bindings, binding)
		}
	}

	return bindings
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

		methodComment := commentInfo.GetSubComment(serviceDescriptorMethodFieldNumber, i).GetText()

		// Collect all HTTP bindings (primary + additional)
		var bindings []*httpBinding
		if hasHTTPAnnotation(method) {
			bindings = collectHTTPBindings(getHTTPAnnotation(method))
		}

		// If no HTTP annotation or no valid bindings, create a default binding
		if len(bindings) == 0 {
			bindings = []*httpBinding{{
				HTTPMethod: "POST",
				URL:        "/" + serviceURLPart + "/" + method.GetName(),
				Body:       nil,
				Index:      0,
			}}
		}

		// Generate method data for each binding
		for _, binding := range bindings {
			methodData := &data.Method{
				Name: method.GetName(),
				URL:  binding.URL,
				Input: &data.MethodArgument{
					Type:       inputTypeFQName,
					IsExternal: isInputTypeExternal,
				},
				Output: &data.MethodArgument{
					Type:       outputTypeFQName,
					IsExternal: isOutputTypeExternal,
				},
				ServerStreaming: method.GetServerStreaming(),
				ClientStreaming: method.GetClientStreaming(),
				HTTPMethod:      binding.HTTPMethod,
				HTTPRequestBody: binding.Body,
				Comment:         methodComment,
				BindingIndex:    binding.Index,
			}

			fileData.TrackPackageNonScalarType(methodData.Input)
			fileData.TrackPackageNonScalarType(methodData.Output)

			serviceData.Methods = append(serviceData.Methods, methodData)
		}
	}

	fileData.Services = append(fileData.Services, serviceData)
}
