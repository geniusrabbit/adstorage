package zoneaccessor

import (
	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/models"

	"github.com/geniusrabbit/adstorage/accessors/accountaccessor"
	"github.com/geniusrabbit/adstorage/accessors/generalaccessor"
	"github.com/geniusrabbit/adstorage/loader"
)

// ZoneAccessor provides accessor to the admodel company type
type ZoneAccessor struct {
	generalaccessor.DataAccessor[admodels.Target, models.Zone]
}

// NewZoneAccessor from dataAccessor
func NewZoneAccessor[AccType any](dataAccessor loader.DataAccessor[models.Zone], accountAccessor *accountaccessor.AccountAccessor[AccType]) *ZoneAccessor {
	return &ZoneAccessor{
		DataAccessor: *generalaccessor.NewDataAccessor(
			dataAccessor,
			func(st *models.Zone) (admodels.Target, bool) {
				acc, _ := accountAccessor.AccountByID(st.AccountID)
				trg := admodels.TargetFromModel(st)
				trg.SetAccount(acc)
				return trg, true
			},
		),
	}
}

// ZoneList returns list of prepared data
func (acc *ZoneAccessor) ZoneList() ([]admodels.Target, error) {
	return acc.List()
}

// TargetByID returns campaign object with specific ID
func (acc *ZoneAccessor) TargetByID(id uint64) (admodels.Target, error) {
	return acc.ByID(id)
}
