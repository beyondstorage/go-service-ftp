package tests

import (
	"testing"

	tests "github.com/beyondstorage/go-integration-test/v4"
)

func TestStorger(t *testing.T) {
	store := initTest(t)

	tests.TestStorager(t, store)
}
