[![Build Status](https://github.com/beyondstorage/go-service-ftp/workflows/Unit%20Test/badge.svg?branch=master)](https://github.com/beyondstorage/go-service-ftp/actions?query=workflow%3A%22Unit+Test%22)
[![License](https://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/Xuanwo/storage/blob/master/LICENSE)
[![](https://img.shields.io/matrix/beyondstorage@go-service-ftp:matrix.org.svg?logo=matrix)](https://matrix.to/#/#beyondstorage@go-service-ftp:matrix.org)

# go-service-ftp

[FTP](https://datatracker.ietf.org/doc/html/rfc959) service support for [go-storage](https://github.com/beyondstorage/go-storage).

## Install

```go
go get github.com/beyondstorage/go-service-ftp
```

## Usage

```go
import (
	"log"

	_ "github.com/beyondstorage/go-service-ftp"
	"github.com/beyondstorage/go-storage/v4/services"
)

func main() {
	store, err := services.NewStoragerFromString("ftp:///path/to/workdir?credential=basic:<user>:<password>&endpoint=tcp:<host>:<port>")
	if err != nil {
		log.Fatal(err)
	}

	// Write data from io.Reader into hello.txt
	n, err := store.Write("hello.txt", r, length)
}
```

- See more examples in [go-storage-example](https://github.com/beyondstorage/go-storage-example).
- Read [more docs](https://beyondstorage.io/docs/go-storage/services/ftp) about go-service-ftp. 
