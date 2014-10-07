__go=$(shell which go)
__goconvey=$(shell which goconvey)

__PROG_NAME=mds
__PROG=$(__PROG_NAME).go

__pwd=$(shell pwd)
__GOPATH=$$GOPATH:$(__pwd)/_vendor:$(__pwd)

init:
	gom install

test:
	GOPATH=$(__GOPATH) $(__go) test

cover:
	GOPATH=$(__GOPATH) $(__goconvey)

.PHONY: test
