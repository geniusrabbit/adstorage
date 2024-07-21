package fstypes

import (
	"github.com/geniusrabbit/adcorelib/models"
)

// CampaignData represented in files
type CampaignData struct {
	Campaigns []*models.Campaign `json:"campaigns" yaml:"campaigns"`
	Campaign  *models.Campaign   `json:"campaign" yaml:"campaign"`
}

// Merge two data items
func (d *CampaignData) Merge(it any) {
	v := it.(*CampaignData)
	if v.Campaign != nil {
		d.Campaigns = append(d.Campaigns, v.Campaign)
	}
	d.Campaigns = append(d.Campaigns, v.Campaigns...)
}

// Result of data as a list
func (d *CampaignData) Result() []*models.Campaign {
	return d.Campaigns
}

// Reset stored data
func (d *CampaignData) Reset() {
	d.Campaign = nil
	d.Campaigns = d.Campaigns[:0]
}
