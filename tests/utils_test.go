package tests

import (
	"os"
	"testing"

	_ "github.com/beyondstorage/go-service-ftp"
	ps "github.com/beyondstorage/go-storage/v4/pairs"
	"github.com/beyondstorage/go-storage/v4/services"
	"github.com/beyondstorage/go-storage/v4/types"
)

func initTest(t *testing.T) (store types.Storager) {
	t.Log("Setup test for ftp")

	store, err := services.NewStorager("ftp",
		ps.WithCredential(os.Getenv("STORAGE_FTP_CREDENTIAL")),
		ps.WithEndpoint(os.Getenv("STORAGE_FTP_ENDPOINT")),
	)
	if err != nil {
		t.Errorf("create storager: %v", err)
	}

	return
}
