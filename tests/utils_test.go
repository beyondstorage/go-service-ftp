package tests

import (
	"os"
	"testing"

	_ "github.com/beyondstorage/go-service-ftp"
	"github.com/beyondstorage/go-storage/v4/services"
	"github.com/beyondstorage/go-storage/v4/types"
	"github.com/google/uuid"
)

func initTest(t *testing.T) (store types.Storager) {
	t.Log("Setup test for ftp")

	hostName := os.Getenv("STORAGE_FTP_NAME")

	store, err := services.NewStoragerFromString(("ftp://") + hostName + "/" + uuid.New().String())
	if err != nil {
		t.Errorf("create storager: %v", err)
	}

	t.Cleanup(func() {
		err = store.Delete("")
		if err != nil {
			t.Errorf("cleanup: %v", err)
		}
	})
	return
}
