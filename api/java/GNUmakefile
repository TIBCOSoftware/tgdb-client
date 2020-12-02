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

.SUFFIXES: .class .java

vpath %.java src/com/tibco/tgdb/channel
vpath %.java src/com/tibco/tgdb/channel/impl
vpath %.java src/com/tibco/tgdb/connection
vpath %.java src/com/tibco/tgdb/connection/impl
vpath %.java src/com/tibco/tgdb/exception
vpath %.java src/com/tibco/tgdb/log
vpath %.java src/com/tibco/tgdb/model
vpath %.java src/com/tibco/tgdb/model/impl
vpath %.java src/com/tibco/tgdb/pdu
vpath %.java src/com/tibco/tgdb/pdu/impl
vpath %.java src/com/tibco/tgdb/utils
vpath %.java src/com/tibco/tgdb

CC=gcc
JAVA_VERSION := 1.8
JAVAC = javac
JFLAGS = -g
JCLASSPATH =
JAR	= jar
JAVA_SRC_DIR := src
CLASSES_DIR  := $(BUILDPATH)/tgdb-client/classes

JAVA_SRCLIST := $(BUILDPATH)/javasrc.list
shellstatus  := $(shell find $(JAVA_SRC_DIR) -name '*.java' | tee $(JAVA_SRCLIST))
JAVA_CLASSES := $(patsubst %.java, %.class, $(JAVA_SRCLIST))
JAVA_RESOURCES := 
LIBDIR	 := $(BUILDPATH)/lib
JAVADOC = javadoc
DOC_DIR := $(BUILDPATH)/doc/api/javadoc
#DOC_PACKAGES := com.tibco.tgdb com.tibco.tgdb.channel com.tibco.tgdb.connection com.tibco.tgdb.exception com.tibco.tgdb.log com.tibco.tgdb.model com.tibco.tgdb.pdu com.tibco.tgdb.utils
DOC_PACKAGES := com.tibco.tgdb com.tibco.tgdb.connection com.tibco.tgdb.model com.tibco.tgdb.exception com.tibco.tgdb.log com.tibco.tgdb.utils com.tibco.tgdb.query com.tibco.tgdb.admin

VERTAGFILE = src/com/tibco/tgdb/TGVersion.tag
VERJAVAFILE = src/com/tibco/tgdb/TGVersion.java
DEFINES:= -DVERS_BUILDTYPE_STR='"'$(BUILD_TYPE)'"' -DVERS_EDITION_STR='"'$(BUILD_EDITION)'"' -DVERS_MAJOR=$(BUILD_MAJOR) \
          -DVERS_MINOR=$(BUILD_MINOR) -DVERS_UPDATE=$(BUILD_UPDATE) -DVERS_HFNO=$(BUILD_HFNO) \
          -DVERS_BUILDNO=$(BUILD_NO) -DVERS_REV=$(BUILD_REV)

default : tgdb-java-client.jar

all : default

clean:
	rm -r -f $(CLASSES_DIR)
	rm -r -f $(DOC_DIR)

createdir:
	@mkdir -p $(CLASSES_DIR)

tag2java:
ifeq ($(BUILD_TAGS), 1)
	-chmod 0644 $(VERTAGFILE)
	-chmod 0644 $(VERJAVAFILE)
	$(CC) $(DEFINES) -traditional -P -C -o $(VERJAVAFILE) -E - < $(VERTAGFILE)
endif

	
tgdb-java-client.jar : createdir tag2java
	@echo building java classes
	@echo $(JAVA_SRCLIST)
	$(JAVAC) $(JFLAGS) -source $(JAVA_VERSION) -target $(JAVA_VERSION) -d $(CLASSES_DIR) -cp $(LIBDIR)/tgdb-gremlin-mincore.jar @$(JAVA_SRCLIST)
	rm $(JAVA_SRCLIST)
	@echo building Jar
	cd $(CLASSES_DIR) ; $(JAR) cf $(LIBDIR)/$@ .
	@echo building Javadoc
	$(JAVADOC) -d $(DOC_DIR) -source $(JAVA_VERSION) -sourcepath $(JAVA_SRC_DIR) $(DOC_PACKAGES)

		
.PHONY:clean createdir
