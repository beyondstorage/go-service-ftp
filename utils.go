package ftp

import (
	"fmt"
	"net/textproto"
	"path/filepath"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"

	"github.com/beyondstorage/go-storage/v4/services"
	"github.com/beyondstorage/go-storage/v4/types"
)

const (
	PathSeparator = string(filepath.Separator)
)

// Storage is the example client.
type Storage struct {
	connection *ftp.ServerConn
	user       string
	password   string
	name       string
	workDir    string

	defaultPairs DefaultStoragePairs
	features     StorageFeatures

	types.UnimplementedStorager
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf("Storager ftp {Name: %s, User: %s, WorkDir: %s}", s.name, s.user, s.workDir)
}

// NewStorager will create Storager only.
func NewStorager(pairs ...types.Pair) (types.Storager, error) {
	return newStoragerWithFTPClient(pairs...)
}

func newStoragerWithFTPClient(pairs ...types.Pair) (store *Storage, err error) {
	defer func() {
		if err != nil {
			err = services.InitError{Op: "new_storager", Type: Type, Err: formatError(err), Pairs: pairs}
		}
	}()

	store = &Storage{
		connection: nil,
		user:       "anonymous",
		password:   "anonymous",
		name:       "localhost:21",
		workDir:    "/",
	}

	opt, err := parsePairStorageNew(pairs)
	if err != nil {
		return
	}
	if opt.HasName {
		store.name = opt.Name
	}
	if opt.HasWorkDir {
		store.workDir = opt.WorkDir
	}
	if opt.HasCredential {

	}
	err = store.connect()
	if err != nil {
		return nil, err
	}
	return
}

func (s *Storage) connect() error {
	c, err := ftp.Dial(s.name, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return err
	}

	err = c.Login(s.user, s.password)
	if err != nil {
		return err
	}
	err = c.ChangeDir(s.workDir)
	if err != nil {
		return err
	}
	s.connection = c
	return err
}

func (s *Storage) makeDir(path string) (err error) {
	rp := s.getAbsPath(path)
	err = s.connection.MakeDir(rp)
	return
}

func (s *Storage) getAbsPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	absPath := filepath.Join(s.workDir, path)

	// Join will clean the trailing "/", we need to append it back.
	if strings.HasSuffix(path, PathSeparator) {
		absPath += PathSeparator
	}
	return absPath
}

func (s *Storage) getNameList(path string) (namelist []string, err error) {
	namelist, err = s.connection.NameList(s.getAbsPath(path))
	if err != nil {
		return nil, err
	}
	return
}

func (s *Storage) newObject(done bool) *types.Object {
	return types.NewObject(s, done)
}

func (s *Storage) mapMode(fet ftp.EntryType) types.ObjectMode {
	switch fet {
	case ftp.EntryTypeFile:
		return types.ModeRead
	case ftp.EntryTypeFolder:
		return types.ModeDir
	case ftp.EntryTypeLink:
		return types.ModeLink
	}
	return types.ModeRead
}

func (s *Storage) formatFileObject(fe *ftp.Entry, parent string) (obj *types.Object, err error) {
	obj = types.NewObject(s, false)
	obj.ID = filepath.Join(parent, fe.Name)
	obj.Mode = s.mapMode(fe.Type)
	obj.Path = fe.Target
	return
}

func formatError(err error) error {
	if _, ok := err.(services.InternalError); ok {
		return err
	}
	switch errX := err.(type) {
	case *textproto.Error:
		switch errX.Code {
		case ftp.StatusInvalidCredentials,
			ftp.StatusLoginNeedAccount,
			ftp.StatusStorNeedAccount:
			return fmt.Errorf("%w, %v", services.ErrPermissionDenied, err)
		case ftp.StatusFileUnavailable,
			ftp.StatusFileActionIgnored:
			return fmt.Errorf("%w, %v", services.ErrObjectNotExist, err)
		default:
			return fmt.Errorf("%w, %v", services.ErrServiceInternal, err)
		}
	}
	return fmt.Errorf("%w, %v", services.ErrUnexpected, err)
}

func (s *Storage) formatError(op string, err error, path ...string) error {
	if err == nil {
		return nil
	}

	return services.StorageError{
		Op:       op,
		Err:      formatError(err),
		Storager: s,
		Path:     path,
	}
}
