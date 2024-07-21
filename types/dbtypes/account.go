package dbtypes

import (
	"sort"
)

type AccountType interface {
	GetID() uint64
	IsApproved() bool
}

// AccountList represented in DB
type AccountList[T any] []T

// Result of data as a list
func (d AccountList[T]) Result() []T {
	return d
}

// Reset stored data
func (d *AccountList[T]) Reset() {
	if d != nil {
		*d = (*d)[:0]
	}
}

// Target real of the list
func (d *AccountList[T]) Target() any {
	return (*[]T)(d)
}

// Merge loaded data
func (d *AccountList[T]) Merge(l any) {
	newData := make([]T, 0, len(*d))
	for _, it := range *d {
		if any(it).(AccountType).IsApproved() {
			newData = append(newData, it)
		}
	}
	for _, it := range l.([]T) {
		if !any(it).(AccountType).IsApproved() {
			continue
		}
		i := sort.Search(len(newData), func(i int) bool { return compareAccount(newData[i], it) >= 0 })
		if i >= 0 && i < len(newData) && compareAccount(newData[i], it) == 0 {
			newData[i] = it
		} else {
			newData = append(newData, it)
		}
	}
	sort.Slice(newData, func(i, j int) bool { return compareAccount(newData[i], newData[j]) < 0 })
	*d = newData
}

func compareAccount(a, b any) int {
	id1 := a.(AccountType).GetID()
	id2 := b.(AccountType).GetID()
	if id1 == id2 {
		return 0
	}
	if id1 < id2 {
		return -1
	}
	return 1
}
