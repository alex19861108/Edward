# This how we want to name the binary output
BINARY=edward
# These are the values we want to pass for VERSION  and BUILD
VERSION=1.0.0
BUILD=`date +%FT%T%z`
# Setup the -Idflags options for go build here,interpolate the variable values
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

# Default target: builds the project
build:
	go build ${LDFLAGS} -o ${BINARY}
# Installs our project: copies binaries
install:
	go install ${LDFLAGS}
# Cleans our project: deletes binaries
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
.PHONY:  clean install
