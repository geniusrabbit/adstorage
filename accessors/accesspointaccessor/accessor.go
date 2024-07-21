package accesspointaccessor

import (
	"context"
	"sort"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/geniusrabbit/adcorelib/accesspoint"
	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/context/ctxlogger"
	"github.com/geniusrabbit/adcorelib/models"

	"github.com/geniusrabbit/adstorage/accessors/accountaccessor"
	"github.com/geniusrabbit/adstorage/accessors/generalaccessor"
	"github.com/geniusrabbit/adstorage/loader"
)

var (
	errUnsupportedAccessPointProtocol = errors.New("unsupported DSP protocol")
)

// Accessor object ad reloader
type Accessor struct {
	generalaccessor.DataAccessor[accesspoint.Platformer, models.RTBAccessPoint]

	factories   map[string]accesspoint.Factory
	factoryList []accesspoint.Factory
}

// NewAccessor object
func NewAccessor[AccType any](
	ctx context.Context,
	dataAccessor loader.DataAccessor[models.RTBAccessPoint],
	accountAccessor *accountaccessor.AccountAccessor[AccType],
	factoryList ...accesspoint.Factory,
) (*Accessor, error) {
	if dataAccessor == nil {
		return nil, errors.New("data accessor is required")
	}
	if accountAccessor == nil {
		return nil, errors.New("account accessor is required")
	}
	factories := make(map[string]accesspoint.Factory, len(factoryList))
	for _, fact := range factoryList {
		factories[fact.Info().Protocol] = fact
	}
	acc := &Accessor{
		factories:   factories,
		factoryList: factoryList,
	}

	acc.DataAccessor = *generalaccessor.NewDataAccessor(
		dataAccessor,
		func(src *models.RTBAccessPoint) (accesspoint.Platformer, bool) {
			if src.AccountID == 0 {
				ctxlogger.Get(ctx).Error("access point without account",
					zap.Uint64("access_point_id", src.ID))
				return nil, false
			}
			accObj, _ := accountAccessor.AccountByID(src.AccountID)
			accessPointModel := admodels.RTBAccessPointFromModel(src, accObj)
			accessPoint, err := acc.newAccessPoint(ctx, accessPointModel)
			if err == nil {
				return accessPoint, true
			} else {
				ctxlogger.Get(ctx).Error("create DSP accessor",
					zap.Uint64("access_point_id", accessPointModel.ID),
					zap.String("access_point_protocol", accessPointModel.Protocol),
					zap.Error(err))
			}
			return nil, false
		},
	)
	return acc, nil
}

// AccesspointList returns list of accesspoints
func (acc *Accessor) AccesspointList() ([]accesspoint.Platformer, error) {
	return acc.List()
}

// AccessPointByID returns source instance
func (acc *Accessor) AccessPointByID(id uint64) (accesspoint.Platformer, error) {
	list, err := acc.AccesspointList()
	if err != nil {
		return nil, err
	}
	idx := sort.Search(len(list), func(i int) bool { return list[i].ID() >= id })
	if idx >= 0 && idx < len(list) && list[idx].ID() == id {
		return list[idx], nil
	}
	return nil, nil
}

// ListFactories of platforms
func (acc *Accessor) ListFactories() []accesspoint.Factory {
	return acc.factoryList
}

// PlatformByProtocol returns platform by codename and protocol
func (acc *Accessor) PlatformByProtocol(protocol, codename string) (accesspoint.Platformer, error) {
	list, err := acc.AccesspointList()
	if err != nil {
		return nil, err
	}
	for _, plt := range list {
		if plt.Codename() == codename {
			if protocol != "" && plt.Protocol() != protocol {
				return nil, errors.Wrap(errUnsupportedAccessPointProtocol, protocol)
			}
			return plt, nil
		}
	}
	return nil, nil
}

func (acc *Accessor) newAccessPoint(ctx context.Context, accessPoint *admodels.RTBAccessPoint) (accesspoint.Platformer, error) {
	if acc.factories == nil {
		return nil, errors.Wrap(errUnsupportedAccessPointProtocol, accessPoint.Protocol)
	}
	fact := acc.factories[accessPoint.Protocol]
	if fact == nil {
		return nil, errors.Wrap(errUnsupportedAccessPointProtocol, accessPoint.Protocol)
	}
	return fact.New(ctx, accessPoint)
}

var _ accesspoint.DSPPlatformAccessor = &Accessor{}
