package dbtypes

import (
	"sort"

	"github.com/geniusrabbit/adcorelib/models"
)

// ApplicationList represented in DB
type ApplicationList []*models.Application

// Result of data as a list
func (d ApplicationList) Result() []*models.Application {
	return d
}

// Reset stored data
func (d *ApplicationList) Reset() {
	if d != nil {
		*d = (*d)[:0]
	}
}

// Target real of the list
func (d *ApplicationList) Target() any {
	return (*[]*models.Application)(d)
}

// Merge loaded data
func (d *ApplicationList) Merge(l any) {
	newData := make([]*models.Application, 0, len(*d))
	for _, it := range *d {
		if it.Status.IsApproved() && it.Active.IsActive() {
			newData = append(newData, it)
		}
	}
	for _, it := range l.([]*models.Application) {
		if !it.Status.IsApproved() || !it.Active.IsActive() {
			continue
		}
		i := sort.Search(len(newData), func(i int) bool { return newData[i].ID >= it.ID })
		if i >= 0 && i < len(newData) && newData[i].ID == it.ID {
			newData[i] = it
		} else {
			newData = append(newData, it)
		}
	}
	sort.Slice(newData, func(i, j int) bool { return newData[i].ID < newData[j].ID })
	*d = newData
}
