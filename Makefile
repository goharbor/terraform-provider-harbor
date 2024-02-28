ACCTEST_PARALLELISM ?= 20
GO_VER              ?= go

default: tools

.PHONY: tools docs default

all: tools docs

tools:
	cd tools && $(GO_VER) install github.com/golangci/golangci-lint/cmd/golangci-lint
	cd tools && $(GO_VER) install github.com/pavius/impi/cmd/impi
	cd tools && $(GO_VER) install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

docs:
	rm -f docs/data-sources/*.md
	rm -f docs/resources/*.md
	@tfplugindocs generate
