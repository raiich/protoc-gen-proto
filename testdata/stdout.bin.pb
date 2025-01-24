Á Ëz¿
Rgithub.com/raiich/protoc-plugin-template-go/generated/go/extension/extension.protozÈsyntax = "proto3";

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

// None is a message that has no field.
message None {
  reserved 1 to max;
}
z¢
Pgithub.com/raiich/protoc-plugin-template-go/generated/go/app/v1/deprecated.protozÕsyntax = "proto3";

package app.v1;

import "extension/extension.proto";

option deprecated = true;
option go_package = "github.com/raiich/protoc-plugin-template-go/generated/go/app/v1";

service DeprecatedService {
  option deprecated = true;

  rpc DeprecatedMethod1(DeprecatedMessage) returns (extension.None) {
    option deprecated = true;
  }
}

message DeprecatedMessage {
  option deprecated = true;
  optional string message = 1 [deprecated = true];
}
zó
Mgithub.com/raiich/protoc-plugin-template-go/generated/go/app/v1/service.protoz≈syntax = "proto3";

package app.v1;

import "extension/extension.proto";

option go_package = "github.com/raiich/protoc-plugin-template-go/generated/go/app/v1";

// Service is a template service.
service Service {
  option (extension.service_opts) = {id: 123};

  // method1 is a template method.
  rpc Method1(Request) returns (extension.None) {
    option (extension.method_opts) = {id: 123};
  }
}

message Request {
  optional string message = 1;
}
