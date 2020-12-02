#
# Copyright (c) 2020 TIBCO Software Inc.
# All rights reserved.
#
#  File name : GNUmakefile
#  Created on: 12/09/2019
#  Created by: nimish
#
#  Makefile for TIBCO Graph Database REST API Subsystem.
#
#

.SUFFIXES: .go

vpath %.go src/tgdbrest


GO 	 := go
REST_HOME := $(ROOTDIR)/api/rest
API_HOME := $(ROOTDIR)/api/go
SRC  := $(REST_HOME)/src
BINDIR := $(BUILDPATH)/bin

ifeq ($(OSNAME), Darwin)
	BINEXT=
	GORUN= export GOPATH=$(REST_HOME):$(API_HOME) && export GOOS=darwin && export GOARCH=amd64 && $(GO)
endif
ifeq ($(OSNAME), Linux)
	BINEXT=
	GORUN= export GOPATH=$(REST_HOME):$(API_HOME) && export GOOS=linux && export GOARCH=amd64 && $(GO)
endif
ifeq ($(OSNAME), Win64)
	BINEXT=.exe
	GORUN= export GOPATH=$(REST_HOME):$(API_HOME) && export GOOS=windows && export GOARCH=amd64 && $(GO)
endif

LDFLAGS=
ifeq ($(BUILD_DEBUG), 0)
	LDFLAGS=-ldflags="-s -w"
endif

default : tgdb-rest

all : default

clean:
	rm -r -f $(BINDIR)/tgdb-rest$(BINEXT)

tgdb-rest:
	@echo Building $@
	$(GORUN) build $(LDFLAGS) -o $(BINDIR)/$@$(BINEXT) $(SRC)/restmain.go
	@echo Done $@

.PHONY: all clean default tgdb-rest
