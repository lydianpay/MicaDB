<div align="center">

# MicaDB

![micadb.png](micadb.png)


</div>

Written in Go ("Golang" for search engines) with zero external dependencies, this package implements an in-memory, 
key-value NoSQL database. Optimized specifically for blazingly fast read/write access, with optional disk storage 
backups.

---

## Installation
1. Once confirming you have [Go](https://go.dev/doc/install) installed, the command below will add
   `micadb` as a dependency to your Go program.
```shell
go get -u github.com/tether-payments/micadb
```
2. Import the package into your code
```go
package main

import (
   "github.com/Tether-Payments/micadb"
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
