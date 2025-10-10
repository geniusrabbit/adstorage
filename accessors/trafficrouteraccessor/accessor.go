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

func (tra *TrafficRouterAccessor) TrafficRouterList(ctx context.Context) ([]*admodels.TrafficRouter, error) {
	return tra.List(ctx)
}

func (tra *TrafficRouterAccessor) TrafficRouterByID(ctx context.Context, id uint64) (*admodels.TrafficRouter, error) {
	return tra.ByKey(ctx, id)
}

func trafficRouterFilter(rt *models.TrafficRouter) types.BaseFilter {
	filter := types.BaseFilter{
		Secure:          int8(rt.Secure),
		AdBlock:         int8(rt.AdBlock),
		PrivateBrowsing: int8(rt.PrivateBrowsing),
		IP:              int8(rt.IP),
	}

	filter.Set(types.FieldFormat, rt.Formats)
	filter.Set(types.FieldDeviceTypes, rt.DeviceTypes)
	filter.Set(types.FieldDevices, rt.Devices)
	filter.Set(types.FieldOS, rt.OS)
	filter.Set(types.FieldBrowsers, rt.Browsers)
	filter.Set(types.FieldCategories, rt.Categories)
	filter.Set(types.FieldCountries, rt.Countries)
	filter.Set(types.FieldLanguages, rt.Languages)

	filter.Set(types.FieldDomains, rt.Domains)
	filter.Set(types.FieldApps, rt.Applications)
	filter.Set(types.FieldZones, rt.Zones)

	return filter
}
