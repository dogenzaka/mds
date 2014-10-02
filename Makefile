__go=$(shell which go)
__goconver=$(shell which goconver)

__PROG_NAME=mds
__PROG=$(__PROG_NAME).go

__GOPATH=$(shell pwd)


test:
	$(__go) test


.PHONY: test