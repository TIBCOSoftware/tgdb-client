#
# Copyright (c) 2019 TIBCO Software Inc.
# All rights reserved.
#
#  File name : GNUmakefile
#  Created on: May 27, 2015
#  Created by: suresh
#
#  SVN Id: $Id: GNUmakefile 3524 2019-11-13 23:10:22Z ssubrama $
#  Makefile for TIBCO Graph Database Gremlin Minimal Core Subsystem for Java
#
#

.SUFFIXES: .class .java

vpath %.java src/org/apache/commons/configuration
vpath %.java src/org/apache/tinkerpop/gremlin/process/computer
vpath %.java src/org/apache/tinkerpop/gremlin/process/remote/traversal
vpath %.java src/org/apache/tinkerpop/gremlin/process/traversal/dsl/graph
vpath %.java src/org/apache/tinkerpop/gremlin/process/traversal/step
vpath %.java src/org/apache/tinkerpop/gremlin/process/traversal/step/util
vpath %.java src/org/apache/tinkerpop/gremlin/process/traversal/traverser
vpath %.java src/org/apache/tinkerpop/gremlin/process/traversal/util
vpath %.java src/org/apache/tinkerpop/gremlin/structure
vpath %.java src/org/apache/tinkerpop/gremlin/structure/io
vpath %.java src/org/apache/tinkerpop/gremlin/structure/util
vpath %.java src/org/apache/tinkerpop/gremlin/structure/util/empty
vpath %.java src/org/apache/tinkerpop/gremlin/util



JAVA_VERSION := 1.8
JAVAC = javac
JFLAGS = 
JCLASSPATH =
JAR	= jar
JAVA_SRC_DIR := src
CLASSES_DIR  := $(BUILDPATH)/gremlin/classes

JAVA_SRCLIST := $(BUILDPATH)/gremlin-src.list
shellstatus  := $(shell find $(JAVA_SRC_DIR) -name '*.java' | tee $(JAVA_SRCLIST))
JAVA_CLASSES := $(patsubst %.java, %.class, $(JAVA_SRCLIST))
JAVA_RESOURCES := 
LIBDIR	 := $(BUILDPATH)/lib

default : tgdb-gremlin-mincore.jar

all : default

clean:
	rm -r -f $(CLASSES_DIR)
	rm -r -f $(DOC_DIR)

createdir:
	@mkdir -p $(CLASSES_DIR)	
	
tgdb-gremlin-mincore.jar : createdir
	@echo building java classes
	@echo $(JAVA_SRCLIST)
	$(JAVAC) $(JFLAGS) -source $(JAVA_VERSION) -target $(JAVA_VERSION) -d $(CLASSES_DIR) @$(JAVA_SRCLIST)
	rm $(JAVA_SRCLIST)
	@echo building Jar
	cd $(CLASSES_DIR) ; $(JAR) cf $(LIBDIR)/$@ .


		
