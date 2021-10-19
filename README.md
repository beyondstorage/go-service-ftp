# go-service-ftp

[FTP](https://datatracker.ietf.org/doc/html/rfc959) service support for [go-storage](https://github.com/beyondstorage/go-storage).

## Notes

**This package has been moved to [go-storage](https://github.com/beyondstorage/go-storage/tree/master/services/ftp).**

```shell
go get go.beyondstorage.io/services/ftp
```

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
