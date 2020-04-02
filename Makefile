# ********************************************************************
# * Makefile                                                         *
# *                                                                  *
# * 2020-03-17 First Version, JR                                     * 
# *                                                                  *
# * File with instructions associated to build and test the project. *
# *                                                                  *
# * Usage:                                                           *
# * $ make build                                                     *
# * $ make unit-test                                                 *
# * $ make test                                                      *
# * $ make cover                                                     *
# * $ make all                                                       *
# ********************************************************************

GO := go
MAKE := make
GOPKS := github.com/stretchr/testify/assert
FILE := authorizer.go

all: build

build:
	@$(GO) get ./...
	@$(GO) get $(GOPKS)
	@$(GO) build $(FILE)

unit-test:
	@$(GO) test -v ./...

test:
	@$(MAKE) -C ./test all

cover:
	@$(GO) test ./... -coverprofile coverage

.PHONY: all build unit-test test cover
