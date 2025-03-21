package gen

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/raiich/protoc-gen-proto/generated/go/extension"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

const (
	ProgramName = "protoc-gen-proto"

	GeneratedFilenameSuffix = "_custom.pb.proto"
)

var (
	FileDescriptorProtoSyntaxFieldNumber  protoreflect.FieldNumber
	FileDescriptorProtoPackageFieldNumber protoreflect.FieldNumber
)

func init() {
	fields := (*descriptorpb.FileDescriptorProto)(nil).ProtoReflect().Descriptor().Fields()
	FileDescriptorProtoSyntaxFieldNumber = fields.ByName("syntax").Number()
	FileDescriptorProtoPackageFieldNumber = fields.ByName("package").Number()
}

func GenerateFile(plugin *protogen.Plugin, file *protogen.File) {
	var suffix string
	if strings.HasSuffix(GeneratedFilenameSuffix, ".proto") {
		suffix = ".proto"
	}
	g := &GeneratedFile{
		Base: plugin.NewGeneratedFile(
			file.GeneratedFilenamePrefix+suffix,
			file.GoImportPath,
		),
	}
	d := file.Desc
	if suffix != ".proto" {
		g.P("// Code generated by ", ProgramName, ". DO NOT EDIT.")
		g.P()
	}
	g.P("syntax = ", strconv.Quote(d.Syntax().String()), ";")
	g.P()
	g.P("package ", d.Package(), ";")
	for i := 0; i < d.Imports().Len(); i++ {
		if i == 0 {
			g.P()
		}
		g.P("import ", strconv.Quote(d.Imports().Get(i).Path()), ";")
	}
	opts := d.Options().(*descriptorpb.FileOptions)
	if opts != nil {
		g.P()
		if opts.Deprecated != nil {
			g.P("option deprecated = ", strconv.FormatBool(*opts.Deprecated), ";")
		}
		if opts.GoPackage != nil {
			g.P("option go_package = ", strconv.Quote(*opts.GoPackage), ";")
		}
	}

	info := &Info{
		File: file,
	}
	for _, service := range file.Services {
		g.P()
		generateService(g, info, service)
	}
	for _, enum := range file.Enums {
		g.P()
		generateEnum(g, info, enum)
	}
	for _, ext := range file.Extensions {
		g.P()
		generateExtension(g, info, ext)
	}
	for _, message := range file.Messages {
		g.P()
		generateMessage(g, info, message)
	}
}

func generateService(g *GeneratedFile, info *Info, s *protogen.Service) {
	comments := s.Comments
	for _, cs := range comments.LeadingDetached {
		for _, c := range getComments(cs) {
			g.P(c)
		}
		g.P()
	}

	for _, c := range getComments(comments.Leading) {
		g.P(c)
	}
	d := s.Desc
	g.P("service ", d.Name(), " {", getTrailingComment(comments.Trailing))
	defer g.P("}")
	g.WithIndent("  ", func() {
		if opts := d.Options().(*descriptorpb.ServiceOptions); opts != nil {
			generateServiceOptions(g, info, opts)
		}
		for _, m := range s.Methods {
			g.P()
			generateMethod(g, info, m)
		}
	})
}

func generateServiceOptions(g *GeneratedFile, info *Info, opts *descriptorpb.ServiceOptions) {
	if opts.Deprecated != nil {
		g.P("option deprecated = ", strconv.FormatBool(*opts.Deprecated), ";")
	}

	extInfo := extension.E_ServiceOpts
	ext := proto.GetExtension(opts, extInfo).(*extension.ServiceOpts)
	if ext == nil {
		return
	}
	protoReflect := ext.ProtoReflect()
	fields := protoReflect.Descriptor().Fields()
	switch fields.Len() {
	case 1:
		i := 0
		field := fields.Get(i)
		g.P("option (", info.Ident(extInfo.TypeDescriptor()), ") = {",
			field.Name(), ": ", protoReflect.Get(field),
			"};",
		)
	default:
		panic("implement me")
	}
}

func generateMethod(g *GeneratedFile, info *Info, m *protogen.Method) {
	comments := m.Comments
	for _, cs := range comments.LeadingDetached {
		for _, c := range getComments(cs) {
			g.P(c)
		}
		g.P()
	}

	for _, c := range getComments(comments.Leading) {
		g.P(c)
	}
	d := m.Desc
	g.P("rpc ", d.Name(), "(", info.Ident(m.Input.Desc), ") returns (", info.Ident(m.Output.Desc), ") {",
		getTrailingComment(comments.Trailing))
	defer g.P("}")
	g.WithIndent("  ", func() {
		opts := d.Options().(*descriptorpb.MethodOptions)
		if opts != nil {
			generateMethodOptions(g, info, opts)
		}
	})
}

func generateMethodOptions(g *GeneratedFile, info *Info, opts *descriptorpb.MethodOptions) {
	if opts.Deprecated != nil {
		g.P("option deprecated = ", strconv.FormatBool(*opts.Deprecated), ";")
	}
	extInfo := extension.E_MethodOpts
	ext := proto.GetExtension(opts, extInfo).(*extension.MethodOpts)
	if ext == nil {
		return
	}
	protoReflect := ext.ProtoReflect()
	fields := protoReflect.Descriptor().Fields()
	switch fields.Len() {
	case 1:
		i := 0
		field := fields.Get(i)
		g.P("option (", info.Ident(extInfo.TypeDescriptor()), ") = {",
			field.Name(), ": ", protoReflect.Get(field),
			"};",
		)
	default:
		panic("implement me")
	}
}

func generateEnum(g *GeneratedFile, info *Info, enum *protogen.Enum) {
	comments := enum.Comments
	for _, cs := range comments.LeadingDetached {
		for _, c := range getComments(cs) {
			g.P(c)
		}
		g.P()
	}

	for _, c := range getComments(comments.Leading) {
		g.P(c)
	}
	d := enum.Desc
	g.P("enum ", d.Name(), " {", getTrailingComment(comments.Trailing))
	defer g.P("}")
	opts := d.Options().(*descriptorpb.EnumOptions)
	if opts != nil {
		if opts.Deprecated != nil {
			g.P("option deprecated = ", strconv.FormatBool(*opts.Deprecated), ";")
		}
	}
	for _, v := range enum.Values {
		vd := v.Desc
		g.P(vd.Name(), " = ", vd.Number())
	}
}

func generateExtension(g *GeneratedFile, info *Info, ext *protogen.Extension) {
	comments := ext.Comments
	for _, cs := range comments.LeadingDetached {
		for _, c := range getComments(cs) {
			g.P(c)
		}
		g.P()
	}

	for _, c := range getComments(comments.Leading) {
		g.P(c)
	}
	g.P("extend ", info.Ident(ext.Extendee.Desc), " {", getTrailingComment(comments.Trailing))
	defer g.P("}")
	g.WithIndent("  ", func() {
		d := ext.Desc
		g.P(info.Ident(ext.Message.Desc), " ", d.Name(), " = ", d.Number(), ";")
	})
}

func generateMessage(g *GeneratedFile, info *Info, message *protogen.Message) {
	comments := message.Comments
	for _, cs := range comments.LeadingDetached {
		for _, c := range getComments(cs) {
			g.P(c)
		}
		g.P()
	}

	for _, c := range getComments(comments.Leading) {
		g.P(c)
	}
	d := message.Desc
	g.P("message ", d.Name(), " {", getTrailingComment(comments.Trailing))
	defer g.P("}")
	g.WithIndent("  ", func() {
		opts := d.Options().(*descriptorpb.MessageOptions)
		if opts != nil {
			if opts.Deprecated != nil {
				g.P("option deprecated = ", strconv.FormatBool(*opts.Deprecated), ";")
			}
			generateMessageOptions(g, info, opts)
		}
		r := d.ReservedRanges()
		for i := 0; i < r.Len(); i++ {
			ns := r.Get(i)
			from, to := ns[0], ns[1]
			g.P("reserved ", num(from), " to ", num(to), ";")
		}
		for _, field := range message.Fields {
			generateMessageField(g, info, field)
		}
	})
}

func generateMessageOptions(g *GeneratedFile, info *Info, opts *descriptorpb.MessageOptions) {
	extInfo := extension.E_MessageOpts
	ext := proto.GetExtension(opts, extInfo).(*extension.MessageOpts)
	if ext == nil {
		return
	}
	protoReflect := ext.ProtoReflect()
	fields := protoReflect.Descriptor().Fields()
	switch fields.Len() {
	case 1:
		i := 0
		field := fields.Get(i)
		g.P("option (", info.Ident(extInfo.TypeDescriptor()), ") = {",
			field.Name(), ": ", protoReflect.Get(field),
			"};",
		)
	default:
		panic("implement me")
	}
}

func generateMessageField(g *GeneratedFile, info *Info, field *protogen.Field) {
	d := field.Desc
	optional := ""
	if d.HasOptionalKeyword() {
		optional = "optional "
	}
	var typeName string
	switch d.Kind() {
	case protoreflect.MessageKind:
		typeName = info.Ident(d.Message())
	default:
		typeName = d.Kind().String()
	}

	opts := d.Options().(*descriptorpb.FieldOptions)
	g.P(optional, typeName, " ", d.Name(), " = ", d.Number(), getFieldOpts(opts), ";")
}

type GeneratedFile struct {
	Base   *protogen.GeneratedFile
	indent string
}

func (g *GeneratedFile) WithIndent(indent string, f func()) {
	g.indent += indent
	defer func() { g.indent = strings.TrimPrefix(g.indent, indent) }()
	f()
}

func (g *GeneratedFile) P(v ...any) {
	if len(v) == 0 {
		g.Base.P()
		return
	}
	g.Base.P(append([]any{g.indent}, v...)...)
}

func num(n protoreflect.FieldNumber) string {
	if n == 536870912 {
		return "max"
	}
	return fmt.Sprintf("%d", n)
}

type Info struct {
	*protogen.File
}

func (i *Info) Ident(d protoreflect.Descriptor) string {
	if d.FullName().Parent() == i.File.Desc.Package() {
		return string(d.Name())
	}
	return string(d.FullName())
}

func getComments(comments protogen.Comments) []string {
	if len(comments) == 0 {
		return nil
	}
	return []string{"// " + strings.TrimSpace(string(comments))}
}

func getTrailingComment(comments protogen.Comments) string {
	if len(comments) == 0 {
		return ""
	}
	return " " + "// " + strings.TrimSpace(string(comments))
}

func getFieldOpts(opts *descriptorpb.FieldOptions) string {
	if opts != nil && opts.Deprecated != nil {
		return " [deprecated = " + strconv.FormatBool(*opts.Deprecated) + "]"
	}
	return ""
}
