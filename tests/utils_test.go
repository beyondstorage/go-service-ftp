package tests

import (
	"testing"

	_ "github.com/beyondstorage/go-service-ftp"
	"github.com/beyondstorage/go-storage/v4/services"
	"github.com/beyondstorage/go-storage/v4/types"
)

func initTest(t *testing.T) (store types.Storager) {
	t.Log("Setup test for ftp")

	store, err := services.NewStoragerFromString("ftp://?endpoint=tcp:127.0.0.1:2121&credential=basic:user:password")
	if err != nil {
		t.Errorf("create storager: %v", err)
	}

	return
}
