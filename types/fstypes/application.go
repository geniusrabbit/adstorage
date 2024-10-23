package fstypes

import "github.com/geniusrabbit/adcorelib/models"

// ApplicationData represented in files
type ApplicationData struct {
	Apps []*models.Application `json:"applications" yaml:"applications"`
	App  *models.Application   `json:"application" yaml:"application"`
}

// Merge two data items
func (d *ApplicationData) Merge(it any) {
	v := it.(*ApplicationData)
	if v.App != nil {
		d.Apps = append(d.Apps, v.App)
	}
	d.Apps = append(d.Apps, v.Apps...)
}

// Result of data as a list
func (d *ApplicationData) Result() []*models.Application {
	return d.Apps
}

// Reset stored data
func (d *ApplicationData) Reset() {
	d.App = nil
	d.Apps = d.Apps[:0]
}
