#
# Copyright (c) 2020 TIBCO Software Inc.
# All rights reserved.
#
#  File name : GNUmakefile
#  Created on: May 27, 2015
#  Created by: suresh
#
#  Makefile for TIBCO Graph Database Java Client Library Subsystem.
#
#

.SUFFIXES: .go

vpath %.go src/tgdb
vpath %.go src/tgdb/impl

GO 	 := go
HOME := $(ROOTDIR)/api/go
SRC  := $(HOME)/src
LIBDIR := $(BUILDPATH)/lib

VERTAGFILE = src/tgdb/impl/clientversionimpl.tag
VERGOFILE = src/tgdb/impl/clientversionimpl.go
DEFINES:= -DVERS_BUILDTYPE_STR='"'$(BUILD_TYPE)'"' -DVERS_EDITION_STR='"'$(BUILD_EDITION)'"' -DVERS_MAJOR=$(BUILD_MAJOR) \
          -DVERS_MINOR=$(BUILD_MINOR) -DVERS_UPDATE=$(BUILD_UPDATE) -DVERS_HFNO=$(BUILD_HFNO) \
          -DVERS_BUILDNO=$(BUILD_NO) -DVERS_REV=$(BUILD_REV)

ifeq ($(OSNAME), Darwin)
	LIB=.a
	BLDMODE=archive
	PKGDIR=$(LIBDIR)/go/pkg/darwin_amd64
	GORUN= export GOPATH=$(HOME) && export GOOS=darwin && export GOARCH=amd64 && $(GO)
endif
ifeq ($(OSNAME), Linux)
	LIB=.a
	BLDMODE=archive
	PKGDIR=$(LIBDIR)/go/pkg/linux_amd64
	GORUN= export GOPATH=$(HOME) && export GOOS=linux && export GOARCH=amd64 && $(GO)
endif
ifeq ($(OSNAME), Win64)
	LIB=.a
	BLDMODE=archive
	PKGDIR=$(LIBDIR)/go/pkg/windows_amd64
	GORUN= export GOPATH=$(HOME) && export GOOS=windows && export GOARCH=amd64 && $(GO)
endif


default : tgdb

all : default

clean:
	rm -r -f $(LIB_DIR)/tgdb-go-client$(SO)


createdir:
	mkdir -p $(PKGDIR)

tag2go:
ifeq ($(BUILD_TAGS), 1)
	-chmod 0644 $(VERTAGFILE)
	-chmod 0644 $(VERGOFILE)
	$(CC) $(DEFINES) -traditional -P -C -o $(VERGOFILE) -E - < $(VERTAGFILE)
endif
	
tgdb: createdir tag2go
	@echo Building $@
#	$(GORUN) list -f '{{.GoFiles}}' $(SRC)/tgdb/...
	$(GORUN) build -buildmode=$(BLDMODE) -a -o $(PKGDIR)/$@$(LIB) $(SRC)/tgdb
	+- rm -r $(BUILDPATH)/doc/api/godoc ||:
	cp -r godoc $(BUILDPATH)/doc/api
	@echo Done $@

.PHONY: all clean createdir default tgdb
