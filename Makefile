__go=$(shell which go)
__goconvey=_vendor/bin/goconvey

__PROG_NAME=mds
__PROG=$(__PROG_NAME).go

__pwd=$(shell pwd)
__GOPATH=$(__pwd)/_vendor:$(__pwd)

init:
	gom install

init-test:
	gom -test install


test:
	GOPATH=$(__GOPATH) $(__go) test

cover:
	GOPATH=$(__GOPATH) $(__goconvey) -depth=0

.PHONY: test
