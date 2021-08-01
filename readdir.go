package ftp

import (
	"context"

	types "github.com/beyondstorage/go-storage/v4/types"
)

func (s *Storage) listDirNext(ctx context.Context, page *types.ObjectPage) (err error) {
	input := page.Status.(*listDirInput)
	if input.objList == nil {
		input.objList, err = s.connection.List(input.rp)
	}
	n := len(input.objList)
	if input.counter >= n {
		return types.IterateDone
	}
	v := input.objList[input.counter]
	obj, err := s.formatFileObject(v, input.rp)
	if err != nil {
		return err
	}
	obj.GetID()
	page.Data = append(page.Data, obj)
	input.counter++
	return
}
