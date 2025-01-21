include ./resources/make/proto.mk
include ./resources/make/diff.mk

.PHONY: clean
clean:
	rm -rf .local
	rm -rf generated
