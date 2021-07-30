package ftp

import (
	"context"

	types "github.com/beyondstorage/go-storage/v4/types"
)

func (s *Storage) listDirNext(ctx context.Context, page *types.ObjectPage) error {
	input := page.Status.(*listDirInput)
	var err error = nil
	if input.objChan == nil {
		input.objChan, err = s.connection.List(input.rp)
	}
	n := len(input.objChan)
	for i := 0; i <= n; i++ {
		if n == i {
			return types.IterateDone
		}
		v := input.objChan[i]

		obj, err := s.formatFileObject(v, input.rp)
		if err != nil {
			return err
		}
		obj.GetID()
		page.Data = append(page.Data, obj)
		input.counter++
	}
	return err
}
