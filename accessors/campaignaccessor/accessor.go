package campaignaccessor

import (
	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/models"

	"github.com/geniusrabbit/adstorage/accessors/accountaccessor"
	"github.com/geniusrabbit/adstorage/accessors/generalaccessor"
	"github.com/geniusrabbit/adstorage/loader"
)

type CampaignPrepareFunc func(c *admodels.Campaign)

// CampaignAccessor provides accessor to the admodel company type
type CampaignAccessor struct {
	generalaccessor.DataAccessor[*admodels.Campaign, uint64, models.Campaign]
}

// NewCampaignAccessor from dataAccessor
func NewCampaignAccessor[AccType any](dataAccessor loader.DataAccessor[models.Campaign], accountAccessor *accountaccessor.AccountAccessor[AccType], formatAccessor types.FormatsAccessor, prepareFunc CampaignPrepareFunc) *CampaignAccessor {
	return &CampaignAccessor{
		DataAccessor: *generalaccessor.NewDataAccessor(
			dataAccessor,
			func(c *models.Campaign) (*admodels.Campaign, bool) {
				acc, _ := accountAccessor.AccountByID(c.AccountID)
				camp := admodels.CampaignFromModel(c, formatAccessor)
				camp.SetAccount(acc)
				if prepareFunc != nil {
					prepareFunc(camp)
				}
				return camp, true
			},
		),
	}
}

// CampaignList returns list of prepared data
func (acc *CampaignAccessor) CampaignList() ([]*admodels.Campaign, error) {
	return acc.List()
}

// CampaignByID returns campaign object with specific ID
func (acc *CampaignAccessor) CampaignByID(id uint64) (*admodels.Campaign, error) {
	return acc.ByKey(id)
}
