
mds During the development :)
=====

[![Build Status](https://travis-ci.org/dogenzaka/mds.svg?branch=master)](https://travis-ci.org/dogenzaka/mds)
[![codecov.io](https://codecov.io/github/dogenzaka/mds/coverage.svg?branch=master)](https://codecov.io/github/dogenzaka/mds?branch=master)

msd is a library for managing multiple databases.

# Features

- [x] Another name on the can be managed as multiple databases.
- Support database
  - [x] MongoDB
  - [ ] Redis
  - [ ] RethinkDB


# Requirements

- [mgo](https://github.com/go-mgo/mgo)
- [mapstructure](https://github.com/mitchellh/mapstructure)


## Test

- [goconvey](https://github.com/smartystreets/goconvey)

# Getting started

```sh
$ go get github.com/dogenzaka/mds
```


# Example

Please read the test code.

# Developer

## Started

```sh
$ make init-test
```

## Test

```sh
$ make test
or
$ make cover
```

# License

MIT License
