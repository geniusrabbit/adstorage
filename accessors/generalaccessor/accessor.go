package generalaccessor

import (
	"sort"

	"github.com/geniusrabbit/adstorage/loader"
)

type ObjectConvertFunc[S, T any] func(*S) (T, bool)

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

	convertFunc ObjectConvertFunc[SrcType, T]
	dataList    []T
}

// NewDataAccessor from dataAccessor
func NewDataAccessor[T TargetObjectType[KT], KT KeyType, SrcType any](dataAccessor loader.DataAccessor[SrcType], convertFunc ObjectConvertFunc[SrcType, T]) *DataAccessor[T, KT, SrcType] {
	return &DataAccessor[T, KT, SrcType]{
		DataAccessor: dataAccessor,
		convertFunc:  convertFunc,
	}
}

// List returns list of prepared data
func (acc *DataAccessor[T, KT, ST]) List() ([]T, error) {
	if acc.dataList != nil && !acc.NeedUpdate() {
		return acc.dataList, nil
	}

	data, err := acc.Data()
	if err != nil {
		return nil, err
	}

	dataList := make([]T, 0, len(data))
	for _, it := range data {
		if obj, ok := acc.convertFunc(it); ok {
			dataList = append(dataList, obj)
		}
	}
	sort.Slice(dataList, func(i, j int) bool { return dataList[i].ObjectKey() < dataList[j].ObjectKey() })
	acc.dataList = dataList

	return acc.dataList, nil
}

// ByKey returns object with specific codename
func (acc *DataAccessor[T, KT, ST]) ByKey(key KT) (T, error) {
	var (
		nilT      T
		list, err = acc.List()
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
