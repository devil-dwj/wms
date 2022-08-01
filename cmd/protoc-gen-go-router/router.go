package main

import (
	"fmt"
	"strconv"
	"strings"

	annotations "google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

const (
	contextPackage = protogen.GoImportPath("context")
	runtimePackage = protogen.GoImportPath("github.com/devil-dwj/wms/runtime/http")
)

func generateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	if len(file.Services) == 0 {
		return nil
	}

	filename := file.GeneratedFilenamePrefix + ".router.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("// Code generated by protoc-gen-go-router. DO NOT EDIT.")

	generateFileContent(gen, file, g)

	return g
}

func generateFileContent(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile) {
	if len(file.Services) == 0 {
		return
	}

	g.P("package ", file.GoPackageName)
	g.P()

	for _, service := range file.Services {
		genService(gen, file, g, service)
	}
}

func genService(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, service *protogen.Service) {
	serverType := service.GoName + "ServerHandler"
	g.Annotate(serverType, service.Location)

	// interface.
	g.P("type ", serverType, " interface {")
	for _, method := range service.Methods {
		g.P(method.Comments.Leading, // 注释
			serverSignature(g, method))
	}
	g.P("}")
	g.P()

	// registration.
	serviceDescVar := service.GoName + "Router" + "_Desc"
	g.P("func Register", service.GoName, "Router(s ", runtimePackage.Ident("ServiceRegistrar"), ", srv ", serverType, ") {")
	g.P("s.RegisterService(&", serviceDescVar, `, srv)`)
	g.P("}")
	g.P()

	// impl.
	handlerNames := make([]string, 0, len(service.Methods))
	for _, method := range service.Methods {
		hname := genServerMethod(gen, file, g, method)
		handlerNames = append(handlerNames, hname)
	}

	// descriptor.
	g.P("var ", serviceDescVar, " = ", runtimePackage.Ident("RouterDesc"), " {")
	g.P("ServiceName: ", strconv.Quote(string(service.Desc.FullName())), ",")
	g.P("Methods: []", runtimePackage.Ident("MethodDesc"), "{")
	for i, method := range service.Methods {
		meth, path := getHttpRule(method)
		g.P("{")
		g.P("Name: ", strconv.Quote(string(method.Desc.Name())), ",")
		g.P("Method: ", strconv.Quote(meth), ",")
		g.P("Path: ", strconv.Quote(path), ",")
		g.P("Handler: ", handlerNames[i], ",")
		g.P("},")
	}
	g.P("},")
	g.P("}")
}

func serverSignature(g *protogen.GeneratedFile, method *protogen.Method) string {
	var reqArgs []string
	ret := "error"

	reqArgs = append(reqArgs, g.QualifiedGoIdent(contextPackage.Ident("Context")))
	reqArgs = append(reqArgs, "*"+g.QualifiedGoIdent(method.Input.GoIdent))
	ret = "(*" + g.QualifiedGoIdent(method.Output.GoIdent) + ", error)"
	return method.GoName + "(" + strings.Join(reqArgs, ", ") + ") " + ret
}

func genServerMethod(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, method *protogen.Method) string {
	service := method.Parent
	hname := fmt.Sprintf("_%sRouter_%s_Handler", service.GoName, method.GoName)

	g.P("func ", hname, "(srv interface{}, ctx ", contextPackage.Ident("Context"), ", dec func(interface{}) error, ", "interceptor ", runtimePackage.Ident("ServerInterceptor)"), " (interface{}, error) {")
	g.P("in := new(", method.Input.GoIdent, ")")
	g.P("if err := dec(in); err != nil { return nil, err }")
	g.P("h := func(ctx ", contextPackage.Ident("Context, "), "req interface{}) (interface{}, error) {")
	g.P("return srv.(", service.GoName, "ServerHandler).", method.GoName, "(ctx, req.(*", method.Input.GoIdent, "))")
	g.P("}")
	g.P("return interceptor(ctx, in, h)")
	g.P("}")
	return hname
}

func getHttpRule(method *protogen.Method) (string, string) {
	if method.Desc.Options() == nil || !proto.HasExtension(method.Desc.Options(), annotations.E_Http) {
		return "", ""
	}

	// http rules
	r := proto.GetExtension(method.Desc.Options(), annotations.E_Http)

	rule := r.(*annotations.HttpRule)
	var meth string
	var path string
	switch {
	case len(rule.GetDelete()) > 0:
		meth = "DELETE"
		path = rule.GetDelete()
	case len(rule.GetGet()) > 0:
		meth = "GET"
		path = rule.GetGet()
	case len(rule.GetPatch()) > 0:
		meth = "PATCH"
		path = rule.GetPatch()
	case len(rule.GetPost()) > 0:
		meth = "POST"
		path = rule.GetPost()
	case len(rule.GetPut()) > 0:
		meth = "PUT"
		path = rule.GetPut()
	}

	if len(meth) == 0 || len(path) == 0 {
		return "", ""
	}

	return meth, path
}
