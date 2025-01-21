package main

import (
	"bytes"
	"io"
	"testing"

	"github.com/raiich/protoc-plugin-template-go/testdata"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/gofeaturespb"
	"google.golang.org/protobuf/types/pluginpb"
)

func TestGenerate(t *testing.T) {
	in, err := io.ReadAll(bytes.NewReader(testdata.PluginStdin))
	assert.NoError(t, err)
	req := &pluginpb.CodeGeneratorRequest{}
	assert.NoError(t, proto.Unmarshal(in, req))

	gen, err := protogen.Options{
		DefaultAPILevel: gofeaturespb.GoFeatures_API_LEVEL_UNSPECIFIED,
	}.New(req)
	assert.NoError(t, err)
	assert.NoError(t, run(gen))

	resp := gen.Response()
	assert.Nil(t, resp.Error)
	out, err := proto.Marshal(resp)
	assert.NoError(t, err)
	assert.Equal(t, testdata.ExpectedStdout, out)
}
