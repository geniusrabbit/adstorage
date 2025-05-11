package zoneaccessor

import (
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/models"

	"github.com/geniusrabbit/adstorage/accessors/accountaccessor"
	"github.com/geniusrabbit/adstorage/accessors/generalaccessor"
	"github.com/geniusrabbit/adstorage/admodels"
	"github.com/geniusrabbit/adstorage/loader"
)

// ZoneAccessor provides accessor to the admodel company type
type ZoneAccessor struct {
	generalaccessor.DataAccessor[adtype.Target, string, models.Zone]
}

// NewZoneAccessor from dataAccessor
func NewZoneAccessor[AccType any](dataAccessor loader.DataAccessor[models.Zone], accountAccessor *accountaccessor.AccountAccessor[AccType]) *ZoneAccessor {
	return &ZoneAccessor{
		DataAccessor: *generalaccessor.NewDataAccessor(
			dataAccessor,
			func(st *models.Zone) (adtype.Target, bool) {
				acc, _ := accountAccessor.AccountByID(st.AccountID)
				trg := targetFromModel(st, acc)
				return trg, true
			},
		),
	}
}

// ZoneList returns list of prepared data
func (acc *ZoneAccessor) ZoneList() ([]adtype.Target, error) {
	return acc.List()
}

// TargetByCodename returns campaign object with specific codename
func (acc *ZoneAccessor) TargetByCodename(codename string) (adtype.Target, error) {
	return acc.ByKey(codename)
}

// TargetFromModel convert datavase model specified model
// which implements Target interface
func targetFromModel(zone *models.Zone, acc adtype.Account) adtype.Target {
	if zone.Type.IsSmartlink() {
		return admodels.SmartlinkFromModel(zone, acc)
	}
	return admodels.AdUnitFromModel(zone, acc)
}
