#TIBCO Software Inc. Copyright 2015

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
DOC_DIR := $(BUILDPATH)/doc/javadoc
#DOC_PACKAGES := com.tibco.tgdb com.tibco.tgdb.channel com.tibco.tgdb.connection com.tibco.tgdb.exception com.tibco.tgdb.log com.tibco.tgdb.model com.tibco.tgdb.pdu com.tibco.tgdb.utils
DOC_PACKAGES := com.tibco.tgdb com.tibco.tgdb.connection com.tibco.tgdb.model com.tibco.tgdb.exception com.tibco.tgdb.log com.tibco.tgdb.utils com.tibco.tgdb.query

default : tgdb-client.jar

all : default

clean:
	rm -r -f $(CLASSES_DIR)
	rm -r -f $(DOC_DIR)

createdir:
	@mkdir -p $(CLASSES_DIR)	
	
tgdb-client.jar : createdir
	@echo building java classes
	@echo $(JAVA_SRCLIST)
	$(JAVAC) $(JFLAGS) -source $(JAVA_VERSION) -target $(JAVA_VERSION) -d $(CLASSES_DIR) -cp $(LIBDIR)/tgdb-gremlin-mincore.jar @$(JAVA_SRCLIST)
	rm $(JAVA_SRCLIST)
	@echo building Jar
	cd $(CLASSES_DIR) ; $(JAR) cf $(LIBDIR)/$@ .
	@echo building Javadoc
	$(JAVADOC) -d $(DOC_DIR) -source $(JAVA_VERSION) -sourcepath $(JAVA_SRC_DIR) $(DOC_PACKAGES)

		
