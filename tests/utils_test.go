package tests

import (
	"log"
	"testing"

	_ "github.com/beyondstorage/go-service-ftp"
	"github.com/beyondstorage/go-storage/v4/services"
	"github.com/beyondstorage/go-storage/v4/types"
)

func initTest(t *testing.T) (store types.Storager) {
	store, err := services.NewStoragerFromString("ftp://localhost:21/")
	if err != nil {
		log.Fatalf("service init failed: %v", err)
	}
	return
}
