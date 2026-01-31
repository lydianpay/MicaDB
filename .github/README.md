<div align="center">

# MicaDB

![micadb.png](micadb.png)
[![Go Report Card](https://goreportcard.com/badge/github.com/lydianpay/micadb)](https://goreportcard.com/report/github.com/lydianpay/micadb)
[![Maintainability](https://qlty.sh/gh/lydianpay/projects/MicaDB/maintainability.svg)](https://qlty.sh/gh/lydianpay/projects/MicaDB)
[![Code Coverage](https://qlty.sh/gh/lydianpay/projects/MicaDB/coverage.svg)](https://qlty.sh/gh/lydianpay/projects/MicaDB)
[![CodeQL](https://github.com/lydianpay/MicaDB/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/lydianpay/MicaDB/actions/workflows/github-code-scanning/codeql)

</div>

Written in Go ("Golang" for search engines) with zero external dependencies, this package implements an in-memory, 
key-value NoSQL database. Optimized specifically for blazingly fast read/write access, with optional disk storage 
backups.

---

## Installation
1. Once confirming you have [Go](https://go.dev/doc/install) installed, the command below will add
   `micadb` as a dependency to your Go program.
```shell
go get -u github.com/lydianpay/micadb
```
2. Import the package into your code
```go
package main

import (
   "github.com/lydianpay/micadb"
)
```

## Usage

### Basic Creation
```go
db, err := micadb.New(micadb.Options{}).Start()
```

### Creation With Options
```go
db, err := micadb.New(micadb.Options{
		Filename:        ".customfilename.bin",
		IsTest:          false,
		BackupFrequency: 10,
	}).Start()
```

### Creation With Custom Types
```go
db, err := micadb.New(micadb.Options{
		Filename:        "./customfilename.bin",
		IsTest:          false,
		BackupFrequency: 10,
	}).WithCustomTypes(
		tests.TestingStruct2{},
		tests.TestingStruct1{},
	).Start()
```

### Get A Value
```go
val := db.Get("some-key")
```

### Set A Value
```go
db.Set(key, value)
```

### Delete A Value
```go
db.Delete(key, value)``
```
