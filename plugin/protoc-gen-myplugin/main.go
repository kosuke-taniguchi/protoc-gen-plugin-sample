package main

import (
	"bytes"
	"go/format"
	"html/template"
	"io"
	"log"
	"os"
	"strings"

	dbtemplate "my-proto-plugin/plugin/protoc-gen-myplugin/template"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
	files := make([]*plugin.CodeGeneratorResponse_File, 0, len(protoFile.GetMessageType()))
	protoFilename := protoFile.GetPackage()
	for _, msg := range protoFile.GetMessageType() {
		filename, code, err := generateCode(protoFilename, msg)
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
	ProtoPackage string
	GoPackage    string
	Entity       string // initial is uppercase
	Table        string
	Fields       []field
	PK           string
}

type field struct {
	Name string
}

// output is filename and its code
func generateCode(profoFilename string, message *descriptor.DescriptorProto) (string, string, error) {
	table := strings.ToLower(message.GetName())
	fields := make([]field, 0, len(message.GetField()))
	for _, f := range message.GetField() {
		name := snakeToPascalCase(f.GetName())
		fields = append(fields, field{Name: name})
	}
	data := templateData{
		ProtoPackage: profoFilename,
		GoPackage:    table,
		Entity:       message.GetName(),
		Fields:       fields,
		PK:           "",
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

func snakeToPascalCase(s string) string {
	words := strings.Split(s, "_")
	for i, word := range words {
		words[i] = cases.Title(language.AmericanEnglish).String(word)
	}
	return strings.Join(words, "")
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
