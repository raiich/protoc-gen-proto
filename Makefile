include ./resources/make/proto.mk

.PHONY: clean
clean:
	rm -rf .local
	rm -rf generated
