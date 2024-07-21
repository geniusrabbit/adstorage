package generalaccessor

import (
	"sort"

	"github.com/geniusrabbit/adstorage/loader"
)

type ObjectConvertFunc[S, T any] func(*S) (T, bool)

// TargetObjectType interface
type TargetObjectType interface {
	ID() uint64
}

// DataAccessor provides accessor to the admodel company type
type DataAccessor[T TargetObjectType, SrcType any] struct {
	loader.DataAccessor[SrcType]

	convertFunc ObjectConvertFunc[SrcType, T]
	dataList    []T
}

// NewDataAccessor from dataAccessor
func NewDataAccessor[T TargetObjectType, SrcType any](dataAccessor loader.DataAccessor[SrcType], convertFunc ObjectConvertFunc[SrcType, T]) *DataAccessor[T, SrcType] {
	return &DataAccessor[T, SrcType]{
		DataAccessor: dataAccessor,
		convertFunc:  convertFunc,
	}
}

// List returns list of prepared data
func (acc *DataAccessor[T, ST]) List() ([]T, error) {
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
	sort.Slice(dataList, func(i, j int) bool { return dataList[i].ID() < dataList[j].ID() })
	acc.dataList = dataList

	return acc.dataList, nil
}

// ByID returns object with specific ID
func (acc *DataAccessor[T, ST]) ByID(id uint64) (T, error) {
	var (
		nilT      T
		list, err = acc.List()
	)
	if err != nil {
		return nilT, err
	}
	idx := sort.Search(len(list), func(i int) bool { return list[i].ID() >= id })
	if idx >= 0 && idx < len(list) && list[idx].ID() == id {
		return list[idx], nil
	}
	return nilT, nil
}
