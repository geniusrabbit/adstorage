package generalaccessor

import (
	"context"
	"sort"

	"github.com/geniusrabbit/adstorage/loader"
)

type ObjectConvertFunc[S, T any] func(context.Context, *S) (T, bool)
type ObjectConvertObjectSliceFunc[S, T any] func(context.Context, *S) []T

type KeyType interface {
	~string | ~int | ~int64 | ~uint | ~uint64
}

// TargetObjectType interface
type TargetObjectType[KT KeyType] interface {
	ObjectKey() KT
}

// DataAccessor provides accessor to the admodel company type
type DataAccessor[T TargetObjectType[KT], KT KeyType, SrcType any] struct {
	loader.DataAccessor[SrcType]

	convertFunc      ObjectConvertFunc[SrcType, T]
	convertSliceFunc ObjectConvertObjectSliceFunc[SrcType, T]
	dataList         []T
}

// NewDataAccessor from dataAccessor
func NewDataAccessor[T TargetObjectType[KT], KT KeyType, SrcType any](
	dataAccessor loader.DataAccessor[SrcType],
	convertFunc ObjectConvertFunc[SrcType, T],
) *DataAccessor[T, KT, SrcType] {
	return &DataAccessor[T, KT, SrcType]{
		DataAccessor: dataAccessor,
		convertFunc:  convertFunc,
	}
}

func NewDataAccessorList[T TargetObjectType[KT], KT KeyType, SrcType any](
	dataAccessor loader.DataAccessor[SrcType],
	convertSliceFunc ObjectConvertObjectSliceFunc[SrcType, T],
) *DataAccessor[T, KT, SrcType] {
	return &DataAccessor[T, KT, SrcType]{
		DataAccessor:     dataAccessor,
		convertSliceFunc: convertSliceFunc,
	}
}

// List returns list of prepared data
func (acc *DataAccessor[T, KT, ST]) List(ctx context.Context) ([]T, error) {
	if acc.dataList != nil && !acc.NeedUpdate() {
		return acc.dataList, nil
	}

	data, err := acc.Data()
	if err != nil {
		return nil, err
	}

	dataList := make([]T, 0, len(data))
	for _, it := range data {
		if acc.convertFunc != nil {
			if obj, ok := acc.convertFunc(ctx, it); ok {
				dataList = append(dataList, obj)
			}
		} else if acc.convertSliceFunc != nil {
			if objs := acc.convertSliceFunc(ctx, it); len(objs) > 0 {
				dataList = append(dataList, objs...)
			}
		}
	}
	sort.Slice(dataList, func(i, j int) bool { return dataList[i].ObjectKey() < dataList[j].ObjectKey() })
	acc.dataList = dataList

	return acc.dataList, nil
}

// ByKey returns object with specific codename
func (acc *DataAccessor[T, KT, ST]) ByKey(ctx context.Context, key KT) (T, error) {
	var (
		nilT      T
		list, err = acc.List(ctx)
	)
	if err != nil {
		return nilT, err
	}
	idx := sort.Search(len(list), func(i int) bool { return list[i].ObjectKey() >= key })
	if idx >= 0 && idx < len(list) && list[idx].ObjectKey() == key {
		return list[idx], nil
	}
	return nilT, nil
}
