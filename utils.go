package ftp

import (
	"fmt"
	"time"

	"github.com/jlaffaye/ftp"

	"github.com/beyondstorage/go-storage/v4/services"
	"github.com/beyondstorage/go-storage/v4/types"
)

// Storage is the example client.
type Storage struct {
	connection *ftp.ServerConn
	user       string

	workDir string

	defaultPairs DefaultStoragePairs
	features     StorageFeatures

	types.UnimplementedStorager
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf("Storager ftp {User: %s, Password: %s, WorkDir: %s}", s.user, s.password, s.workDir)
}

// NewStorager will create Storager only.
func NewStorager(pairs ...types.Pair) (types.Storager, error) {
	return newStoragerWithFTPClient(pairs...)
}

func newStoragerWithFTPClient(pairs ...types.Pair) (store *Storage, err error) {
	defer func() {
		if err != nil {
			err = services.InitError{Op: "new_storager", Type: Type, Err: formatErr(err), Pairs: pairs}
		}
	}()
	opt, err = parsePairStorageNew(pairs)
	if err != nil {
		return
	}

	c, err := ftp.Dial("ftp.example.org:21", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return
	}
	store = &Storage{
		connection: c,
		workDir:    "",
	}

	return
}

func formatErr(err error) error {
	if _, ok := err.(services.InternalError); ok {
		return err
	}
	panic("implement me")
}

func (s *Storage) newObject(done bool) *types.Object {
	return types.NewObject(s, done)
}

func (s *Storage) formatError(op string, err error, path ...string) error {
	panic("implement me")
}
