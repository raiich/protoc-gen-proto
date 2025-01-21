Á Ëzò
Rgithub.com/raiich/protoc-plugin-template-go/generated/go/extension/extension.protoz¡syntax = "proto3";

package extension;

import "google/protobuf/descriptor.proto";

option go_package = "github.com/raiich/protoc-plugin-template-go/generated/go/extension";

extend google.protobuf.ServiceOptions {
  ServiceOpts service_opts = 50001;
}

extend google.protobuf.MethodOptions {
  MethodOpts method_opts = 50001;
}

message ServiceOpts {
  int32 id = 1;
}

message MethodOpts {
  int32 id = 1;
}

message None {
  reserved 1 to max;
}
z“
Mgithub.com/raiich/protoc-plugin-template-go/generated/go/app/v1/service.protozÄsyntax = "proto3";

package app.v1;

import "extension/extension.proto";

option go_package = "github.com/raiich/protoc-plugin-template-go/generated/go/app/v1";

service Service {
  option (extension.service_opts) = {id: 123};

  rpc Method1(Request) returns (extension.None) {
    option (extension.method_opts) = {id: 123};
  }
}

message Request {
  optional string message = 1;
}
