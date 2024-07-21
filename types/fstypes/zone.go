package fstypes

import "github.com/geniusrabbit/adcorelib/models"

// ZoneData represented in files
type ZoneData struct {
	Zones []*models.Zone `json:"zones" yaml:"zones"`
	Zone  *models.Zone   `json:"zone" yaml:"zone"`
}

// Merge two data items
func (d *ZoneData) Merge(it any) {
	v := it.(*ZoneData)
	if v.Zone != nil {
		d.Zones = append(d.Zones, v.Zone)
	}
	d.Zones = append(d.Zones, v.Zones...)
}

// Result of data as a list
func (d *ZoneData) Result() []*models.Zone {
	return d.Zones
}

// Reset stored data
func (d *ZoneData) Reset() {
	d.Zone = nil
	d.Zones = d.Zones[:0]
}
