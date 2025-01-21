PROTOC_URL := https://github.com/protocolbuffers/protobuf/releases/download/v29.3/protoc-29.3-osx-aarch_64.zip
PROTODEP_VERSION := v0.1.8
BUF_VERSION := v1.49.0

BIN    := .local/bin
GO_BIN := .local/go/bin
TEMP   := .local/temp

THIRD_PARTY := .local/protodep
GO_INSTALL := GOBIN=$(CURDIR)/$(GO_BIN) go install

PATH := $(BIN):$(GO_BIN):$(PATH)
PROTOC := $(BIN)/protoc
PROTODEP := $(GO_BIN)/protodep
PROTOC_DEP := $(PROTOC) \
  $(GO_BIN)/protoc-gen-go \
  $(GO_BIN)/protoc-gen-template-go \
  $(GO_BIN)/protoc-gen-testdata \
  protodep.lock
BUF := $(GO_BIN)/buf

GO_FILES := $(shell find . -path "./generated" -prune -or -name \*.go -print)
PROTO := $(shell find resources/proto -name \*.proto)

generated: $(PROTO) $(PROTOC_DEP)
	mkdir -p generated
	touch generated
	mkdir -p generated/go
	mkdir -p generated/proto
	$(PROTOC) -I=resources/proto -I=$(THIRD_PARTY) \
	  --go_out=generated/go --go_opt=paths=source_relative \
	  --template-go_out=generated/proto --template-go_opt=paths=source_relative \
	  $(PROTO)

testdata: $(PROTO) $(PROTOC_DEP)
	mkdir -p testdata
	touch testdata
	$(PROTOC) -I=resources/proto -I=$(THIRD_PARTY) \
	  --testdata_out=testdata \
	  $(PROTO)

protodep.lock: protodep.toml $(PROTODEP)
	$(PROTODEP) up --force --use-https
	@mkdir -p $(THIRD_PARTY)  # `--force` may remove dir...

protodep.toml:
	# ok

.PHONY: format-proto
format-proto: $(BUF)
	$(BUF) format resources/proto -w

$(BUF):
	mkdir -p $(GO_BIN)
	$(GO_INSTALL) github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION)

$(PROTOC):
	mkdir -p $(TEMP)
	curl -sSL -o $(TEMP)/protoc.zip "$(PROTOC_URL)"
	unzip -q $(TEMP)/protoc.zip -d .local

$(GO_BIN)/protoc-gen-template-go: $(GO_FILES)
	mkdir -p $(GO_BIN)
	$(GO_INSTALL) ./cmd/protoc-gen-template-go

$(GO_BIN)/protoc-gen-testdata: $(GO_FILES)
	mkdir -p $(GO_BIN)
	$(GO_INSTALL) ./cmd/protoc-gen-testdata

$(GO_BIN)/protoc-gen-go:
	mkdir -p $(GO_BIN)
	$(GO_INSTALL) google.golang.org/protobuf/cmd/protoc-gen-go@latest

$(PROTODEP):
	mkdir -p $(GO_BIN)
	$(GO_INSTALL) github.com/stormcat24/protodep@$(PROTODEP_VERSION)
