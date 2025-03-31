package dbtypes

import (
	"sort"

	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/models"
	"gorm.io/gorm"
)

// TrafficRouterList represents a list of traffic routers in the database.
type TrafficRouterList []*models.TrafficRouter

// Result returns the list of traffic routers.
func (d TrafficRouterList) Result() []*models.TrafficRouter {
	return d
}

// Reset clears the stored data.
func (d *TrafficRouterList) Reset() {
	if d != nil {
		*d = (*d)[:0]
	}
}

// Target returns the underlying slice of traffic routers.
func (d *TrafficRouterList) Target() any {
	return (*[]*models.TrafficRouter)(d)
}

// Merge combines the loaded data with the existing data in the list.
func (d *TrafficRouterList) Merge(l any) {
	newData := make([]*models.TrafficRouter, 0, len(*d))
	newData = append(newData, (*d)...)
	for _, it := range l.([]*models.TrafficRouter) {
		if !it.Active.IsActive() {
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

func (d *TrafficRouterList) PrepareQuery(db *gorm.DB) *gorm.DB {
	return db.Where("active = ?", types.StatusActive)
}
