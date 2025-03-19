package main

import (
	"github.com/raiich/protoc-gen-proto/internal/gen"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/gofeaturespb"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	protogen.Options{
		DefaultAPILevel: gofeaturespb.GoFeatures_API_LEVEL_UNSPECIFIED,
	}.Run(run)
}

func run(plugin *protogen.Plugin) error {
	plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL) |
		uint64(pluginpb.CodeGeneratorResponse_FEATURE_SUPPORTS_EDITIONS)
	plugin.SupportedEditionsMinimum = descriptorpb.Edition_EDITION_PROTO3
	plugin.SupportedEditionsMaximum = descriptorpb.Edition_EDITION_2023
	for _, file := range plugin.Files {
		if file.Generate {
			gen.GenerateFile(plugin, file)
		}
	}
	return nil
}
