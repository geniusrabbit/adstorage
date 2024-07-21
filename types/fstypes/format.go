package fstypes

import (
	"github.com/geniusrabbit/adcorelib/models"
)

// FormatData represented in files
type FormatData struct {
	Formats []*models.Format `json:"formats"`
	Format  *models.Format   `json:"format"`
}

// Merge two data items
func (d *FormatData) Merge(it any) {
	v := it.(*FormatData)
	if v.Format != nil {
		d.Formats = append(d.Formats, v.Format)
	}
	d.Formats = append(d.Formats, v.Formats...)
}

// Result of data as a list
func (d *FormatData) Result() []*models.Format {
	return d.Formats
}

// Reset stored data
func (d *FormatData) Reset() {
	d.Format = nil
	d.Formats = d.Formats[:0]
}
