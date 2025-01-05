package main

import (
	"bytes"
	"errors"
	"go/format"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	dbtemplate "my-proto-plugin/cmd/protoc-gen-myplugin/template"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"google.golang.org/protobuf/proto"
)

func run() error {
	req, err := parseReq(os.Stdin)
	if err != nil {
		return err
	}
	resp, err := process(req)
	if err != nil {
		return err
	}
	return emitResp(resp)
}

func parseReq(r io.Reader) (*plugin.CodeGeneratorRequest, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var req plugin.CodeGeneratorRequest
	if err = proto.Unmarshal(buf, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

func emitResp(resp *plugin.CodeGeneratorResponse) error {
	buf, err := proto.Marshal(resp)
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(buf)
	return nil
}

func process(req *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {
	files := make(map[string]*descriptor.FileDescriptorProto, 0)
	for _, f := range req.ProtoFile {
		files[f.GetName()] = f
	}
	var resp plugin.CodeGeneratorResponse
	for _, fname := range req.FileToGenerate {
		files, err := generateFiles(files[fname])
		if err != nil {
			return nil, err
		}
		resp.File = append(resp.File, files...)
	}
	return &resp, nil
}

func generateFiles(protoFile *descriptor.FileDescriptorProto) ([]*plugin.CodeGeneratorResponse_File, error) {
	var goPkgName string
	if options := protoFile.GetOptions(); options != nil {
		goPkgName = filepath.Base(options.GetGoPackage())
	}
	if goPkgName == "" {
		return nil, errors.New("package name not found")
	}
	files := make([]*plugin.CodeGeneratorResponse_File, 0, len(protoFile.GetMessageType()))
	for _, msg := range protoFile.GetMessageType() {
		filename, code, err := generateCode(msg)
		if err != nil {
			return nil, err
		}
		files = append(files, &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(filename),
			Content: proto.String(code),
		})
	}
	return files, nil
}

type templateData struct {
	GoPackage string
	Entity    string // initial is uppercase
	Table     string
	Fields    []field
	PK        string
}

type field struct {
	Name string
}

// output is filename and its code
func generateCode(message *descriptor.DescriptorProto) (string, string, error) {
	table := strings.ToLower(message.GetName())
	fields := make([]field, 0, len(message.GetField()))
	for _, f := range message.GetField() {
		name := snakeCaseToPascalCase(f.GetName())
		fields = append(fields, field{Name: name})
	}
	data := templateData{
		GoPackage: table,
		Entity:    message.GetName(),
		Table:     table,
		Fields:    fields,
		PK:        "",
	}
	tmpl, err := template.New("goCode").Parse(dbtemplate.DBTemplate)
	if err != nil {
		return "", "", err
	}
	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, data); err != nil {
		return "", "", err
	}
	formated, err := format.Source(buf.Bytes())
	if err != nil {
		return "", "", err
	}
	return table + ".gen.go", string(formated), nil
}

func snakeCaseToPascalCase(s string) string {
	return strings.Title(strings.ReplaceAll(s, "_", ""))
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
