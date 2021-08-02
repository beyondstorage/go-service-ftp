package tests

import (
	"testing"

	ftp "github.com/beyondstorage/go-service-ftp"
	"github.com/beyondstorage/go-storage/v4/types"
)

func initTest(t *testing.T) (store types.Storager) {
	t.Log("Setup test for ftp")

	store, err := ftp.NewStorager()
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
