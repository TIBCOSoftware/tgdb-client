# Don't search for anything!
.SUFFIXES:

#
# Make commands for subdirectories common to all platforms
#

#Initialize variables
export ROOTDIR		:=  $(CURDIR)
export BUILDPATH	:=	$(CURDIR)/build/

all:	default

default:
	mkdir -p $(BUILDPATH)
	mkdir -p $(BUILDPATH)/lib
	cd gremlin-mincore && $(MAKE)
	cd java && $(MAKE)
