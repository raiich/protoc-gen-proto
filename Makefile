include ./resources/make/proto.mk

.PHONY: clean
clean:
	rm -rf .local
	rm -rf generated

.PHONY: test
test: generated testdata
	diff -ru resources/proto/ generated/proto/
	go test ./...
