package trafficrouteraccessor

import (
	"context"
	"slices"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/models"

	"github.com/geniusrabbit/adstorage/accessors/generalaccessor"
	"github.com/geniusrabbit/adstorage/loader"
)

// TrafficRouterAccessor is a data accessor for traffic routers.
type TrafficRouterAccessor struct {
	generalaccessor.DataAccessor[*admodels.TrafficRouter, uint64, models.TrafficRouter]
}

// NewTrafficRouterAccessor creates a new TrafficRouterAccessor with the given data accessor.
func NewTrafficRouterAccessor(dataAccessor loader.DataAccessor[models.TrafficRouter]) *TrafficRouterAccessor {
	return &TrafficRouterAccessor{
		DataAccessor: *generalaccessor.NewDataAccessor(
			dataAccessor,
			func(ctx context.Context, rt *models.TrafficRouter) (*admodels.TrafficRouter, bool) {
				return &admodels.TrafficRouter{
					ID:           rt.ID,
					RTBSourceIDs: slices.Clone(rt.RTBSourceIDs),
					Percent:      float32(rt.Percent),
					Filter:       trafficRouterFilter(rt),
				}, true
			},
		),
	}
}

// TrafficRouterList returns a list of all traffic routers.
func (tra *TrafficRouterAccessor) TrafficRouterList(ctx context.Context) ([]*admodels.TrafficRouter, error) {
	return tra.List(ctx)
}

// TrafficRouterByID returns a traffic router by its ID.
func (tra *TrafficRouterAccessor) TrafficRouterByID(ctx context.Context, id uint64) (*admodels.TrafficRouter, error) {
	return tra.ByKey(ctx, id)
}

func trafficRouterFilter(rt *models.TrafficRouter) types.BaseFilter {
	fl := types.BaseFilter{
		Secure:          int8(rt.Secure),
		AdBlock:         int8(rt.AdBlock),
		PrivateBrowsing: int8(rt.PrivateBrowsing),
		IP:              int8(rt.IP),
	}

	fl.SetFormats(rt.Formats)
	// fl.SetInterstitialFormats(rt.InterstitialFormats)
	fl.SetDeviceTypes(rt.DeviceTypes)
	fl.SetDevices(rt.Devices)
	fl.SetOS(rt.OS)
	fl.SetBrowsers(rt.Browsers)
	// fl.SetCarriers(rt.Carriers)
	fl.SetCategories(rt.Categories)
	fl.SetCountries(rt.Countries)
	fl.SetLanguages(rt.Languages)

	fl.SetDomains(rt.Domains)
	fl.SetAppIDs(rt.Applications)
	fl.SetZoneIDs(rt.Zones)

	return fl
}
