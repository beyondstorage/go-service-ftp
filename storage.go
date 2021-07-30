package ftp

import (
	"context"
	"io"
	"path/filepath"

	"github.com/beyondstorage/go-storage/v4/pkg/iowrap"
	. "github.com/beyondstorage/go-storage/v4/types"
	"github.com/jlaffaye/ftp"
	mime "github.com/qingstor/go-mime"
)

type listDirInput struct {
	rp  string
	dir string

	started           bool
	continuationToken string
	objChan           []*ftp.Entry
	counter           int
	buf               []byte
}

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
func (input *listDirInput) ContinuationToken() string {
	return input.continuationToken
}

func (s *Storage) list(ctx context.Context, path string, opt pairStorageList) (oi *ObjectIterator, err error) {
	//rp := s.getAbsPath(path)
	input := listDirInput{
		// Always keep service original name as rp.
		rp: s.getAbsPath(path),
		// Then convert the dir to slash separator.
		dir: filepath.ToSlash(path),

		// if HasContinuationToken, we should start after we scanned this token.
		// else, we can start directly.
		started:           !opt.HasContinuationToken,
		continuationToken: opt.ContinuationToken,

		buf: make([]byte, 8192),
	}

	return NewObjectIterator(ctx, s.listDirNext, &input), nil
}

func (s *Storage) metadata(opt pairStorageMetadata) (meta *StorageMeta) {
	meta = NewStorageMeta()
	meta.WorkDir = s.workDir
	return meta

}

func (s *Storage) read(ctx context.Context, path string, w io.Writer, opt pairStorageRead) (n int64, err error) {
	rp := s.getAbsPath(path)

	if err != nil {
		return 0, err
	}

	r, err := s.connection.Retr(rp)
	if opt.HasOffset {
		r, err = s.connection.RetrFrom(rp, uint64(opt.Offset))
		if err != nil {
			return n, err
		}
	}
	defer func() {
		closeErr := r.Close()
		if err == nil {
			err = closeErr
		}
	}()

	// if opt.HasIoCallback {
	// 	r = iowrap.CallbackReadCloser(*r, opt.IoCallback)
	// }
	if opt.HasSize {
		return io.CopyN(w, r, opt.Size)
	}
	return io.Copy(w, r)
}

func (s *Storage) stat(ctx context.Context, path string, opt pairStorageStat) (o *Object, err error) {
	rp := s.getAbsPath(path)

	fl, err := s.connection.List(filepath.Dir(rp))
	if err != nil {
		return nil, err
	}
	var fe *ftp.Entry = nil
	for i := range fl {
		if fl[i].Name == filepath.Base(rp) {
			fe = fl[i]
		}
	}
	if fe == nil {
		return nil, err
	}
	o = s.newObject(true)
	o.ID = rp
	o.Path = path

	if fe.Type == ftp.EntryTypeFolder {
		o.Mode |= ModeDir
		return
	}

	if fe.Type != ftp.EntryTypeLink {
		o.Mode |= ModeRead | ModePage | ModeAppend

		o.SetContentLength(int64(fe.Size))
		o.SetLastModified(fe.Time)

		if v := mime.DetectFilePath(path); v != "" {
			o.SetContentType(v)
		}
	}

	// Check if this file is a link.
	if fe.Type == ftp.EntryTypeLink {
		o.Mode |= ModeLink

		target, err := filepath.EvalSymlinks(rp)
		if err != nil {
			return nil, err
		}
		o.SetLinkTarget(target)
	}

	return o, nil
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
