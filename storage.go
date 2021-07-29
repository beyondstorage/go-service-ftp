package ftp

import (
	"context"
	"io"
	"path/filepath"

	"github.com/aos-dev/go-storage/v2/pkg/iowrap"
	. "github.com/beyondstorage/go-storage/v4/types"
)

func (s *Storage) create(path string, opt pairStorageCreate) (o *Object) {
	if opt.HasObjectMode && opt.ObjectMode.IsDir() {
		o = s.newObject(false)
		o.Mode = ModeDir
	} else {
		o = s.newObject(false)
		o.Mode = ModeRead
	}

	o.ID = filepath.Join(s.workDir, path)
	o.Path = path
	return o
}
func (s *Storage) createDir(ctx context.Context, path string) (o *Object, err error) {
	rp := s.getAbsPath(path)
	s.connection.ChangeDir(s.workDir)
	err = s.connection.MakeDir(rp)
	if err != nil {
		return
	}

	o = s.newObject(true)
	o.ID = rp
	o.Path = path
	o.Mode |= ModeDir
	return
}

func (s *Storage) delete(ctx context.Context, path string, opt pairStorageDelete) (err error) {
	rp := s.getAbsPath(path)
	s.connection.ChangeDir(filepath.Dir(rp))
	err = s.connection.Delete(filepath.Base(rp))
	s.connection.ChangeDir(s.workDir)
	if err != nil {
		return err
	}
	return
}

func (s *Storage) list(ctx context.Context, path string, opt pairStorageList) (oi *ObjectIterator, err error) {
	s.connection.ChangeDir(path)

	panic("not implemented")
}

func (s *Storage) metadata(opt pairStorageMetadata) (meta *StorageMeta) {
	panic("not implemented")
}

func (s *Storage) read(ctx context.Context, path string, w io.Writer, opt pairStorageRead) (n int64, err error) {
	rp := s.getAbsPath(path)

	if err != nil {
		return 0, err
	}

	r, err := s.connection.Retr(rp)

	defer func() {
		closeErr := r.Close()
		if err == nil {
			err = closeErr
		}
	}()

	if opt.HasOffset {
		_, err = r.Seek(opt.Offset, 0)
		if err != nil {
			return n, err
		}
	}
	if opt.HasSize {
		r = iowrap.LimitReadCloser(r, opt.Size)
	}
	if opt.HasIoCallback {
		r = iowrap.CallbackReadCloser(r, opt.IoCallback)
	}

	return io.Copy(w, r)
}

func (s *Storage) stat(ctx context.Context, path string, opt pairStorageStat) (o *Object, err error) {
	panic("not implemented")
}

func (s *Storage) write(ctx context.Context, path string, r io.Reader, size int64, opt pairStorageWrite) (n int64, err error) {

	lr := io.LimitReader(r, n)
	if opt.HasIoCallback {
		lr = iowrap.CallbackReader(lr, opt.IoCallback)
	}
	rp := s.getAbsPath(path)
	err = s.makeDir(filepath.Dir(rp))
	if err != nil {
		return
	}
	err = s.connection.Stor(rp, lr)
	if err != nil {
		return
	}
	return
}
