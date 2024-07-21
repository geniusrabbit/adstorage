package dbtypes

import (
	"sort"

	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/models"
	"gorm.io/gorm"
)

// RTBAccessPointList represented in DB
type RTBAccessPointList []*models.RTBAccessPoint

// Result of data as a list
func (d RTBAccessPointList) Result() []*models.RTBAccessPoint {
	return d
}

// Reset stored data
func (d *RTBAccessPointList) Reset() {
	if d != nil {
		*d = (*d)[:0]
	}
}

// Target real of the list
func (d *RTBAccessPointList) Target() any {
	return (*[]*models.RTBAccessPoint)(d)
}

// Merge loaded data
func (d *RTBAccessPointList) Merge(l any) {
	newData := make([]*models.RTBAccessPoint, 0, len(*d))
	for _, it := range *d {
		if it.Status.IsApproved() && it.Active.IsActive() {
			newData = append(newData, it)
		}
	}
	for _, it := range l.([]*models.RTBAccessPoint) {
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

func (d *RTBAccessPointList) PrepareQuery(db *gorm.DB) *gorm.DB {
	return db.Where("status = ? AND active = ?", types.StatusApproved, types.StatusActive)
}
