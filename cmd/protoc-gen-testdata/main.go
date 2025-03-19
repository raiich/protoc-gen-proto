package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/raiich/protoc-gen-proto/lib/must"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/gofeaturespb"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	backupStdin := os.Stdin
	defer func() { os.Stdin = backupStdin }()

	stdinBytes := must.Must(io.ReadAll(os.Stdin))
	stdinMock := must.Must(os.CreateTemp("", "Stdin"))
	defer func() { must.NoError(stdinMock.Close()) }()
	must.Must(stdinMock.Write(stdinBytes))
	must.Must(stdinMock.Seek(0, io.SeekStart))

	os.Stdin = stdinMock
	protogen.Options{
		DefaultAPILevel: gofeaturespb.GoFeatures_API_LEVEL_UNSPECIFIED,
	}.Run(func(plugin *protogen.Plugin) error {
		plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL) |
			uint64(pluginpb.CodeGeneratorResponse_FEATURE_SUPPORTS_EDITIONS)
		plugin.SupportedEditionsMinimum = descriptorpb.Edition_EDITION_PROTO3
		plugin.SupportedEditionsMaximum = descriptorpb.Edition_EDITION_2023
		g := plugin.NewGeneratedFile(
			filepath.Join("stdin.bin.pb"),
			"",
		)
		_, err := g.Write(stdinBytes)
		return err
	})
}
