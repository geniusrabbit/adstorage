package adstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/geniusrabbit/adcorelib/admodels/types"

	"github.com/geniusrabbit/adstorage/accessors/accountaccessor"
	"github.com/geniusrabbit/adstorage/accessors/adsourceaccessor"
	"github.com/geniusrabbit/adstorage/accessors/appaccessor"
	"github.com/geniusrabbit/adstorage/accessors/formataccessor"
	"github.com/geniusrabbit/adstorage/accessors/trafficrouteraccessor"
	"github.com/geniusrabbit/adstorage/accessors/zoneaccessor"
	"github.com/geniusrabbit/adstorage/database"
	"github.com/geniusrabbit/adstorage/loader"
	"github.com/geniusrabbit/adstorage/loader/dbloader"
	"github.com/geniusrabbit/adstorage/loader/fsloader"
)

type AllAccessor[AccType any] struct {
	ctx       context.Context
	accessors AllDataAccessor[AccType]

	formats        types.FormatsAccessor
	accountCast    accountaccessor.AccountConvertFunc[AccType]
	accounts       *accountaccessor.AccountAccessor[AccType]
	apps           *appaccessor.AppAccessor
	zones          *zoneaccessor.ZoneAccessor
	trafficRouters *trafficrouteraccessor.TrafficRouterAccessor
}

func NewDataAccessor[T any, AccType any](conn any, dbtarget, fstarget loader.LoaderTarget[T], pattern, metricName string) (loader.DataAccessor[T], error) {
	type iDB interface {
		__isDB()
		getDB() *database.DB
		gerPeriod() time.Duration
	}
	type iFS interface {
		__isFS()
		getRoot() string
		gerPeriod() time.Duration
	}
	switch conn := conn.(type) {
	case iDB:
		return loader.NewPeriodicReloader(
			dbtarget,
			dbloader.Loader(conn.getDB()),
			conn.gerPeriod(),
			conn.gerPeriod()*10,
			metricName,
		), nil
	case iFS:
		return loader.NewPeriodicReloader(
			fstarget,
			fsloader.PatternLoader(conn.getRoot(), pattern),
			conn.gerPeriod(),
			conn.gerPeriod()*10,
			metricName,
		), nil
	}
	return nil, fmt.Errorf("unsupported connection type %T", conn)
}

func NewAllAccessor[AccType any](ctx context.Context, accessors AllDataAccessor[AccType], accountCast accountaccessor.AccountConvertFunc[AccType]) *AllAccessor[AccType] {
	return &AllAccessor[AccType]{
		ctx:         ctx,
		accessors:   accessors,
		accountCast: accountCast,
	}
}

func (acc *AllAccessor[AccType]) Formats() (types.FormatsAccessor, error) {
	if acc.formats != nil {
		return acc.formats, nil
	}
	formatDataAccessor, err := acc.accessors.Formats()
	if err != nil {
		return nil, err
	}
	acc.formats = formataccessor.NewFormatAccessor(formatDataAccessor)
	return acc.formats, nil
}

func (acc *AllAccessor[AccType]) Accounts() (*accountaccessor.AccountAccessor[AccType], error) {
	if acc.accounts != nil {
		return acc.accounts, nil
	}
	accountDataAccessor, err := acc.accessors.Accounts()
	if err != nil {
		return nil, err
	}
	acc.accounts = accountaccessor.NewAccessor(accountDataAccessor, acc.accountCast)
	return acc.accounts, nil
}

func (acc *AllAccessor[AccType]) Apps() (*appaccessor.AppAccessor, error) {
	if acc.apps != nil {
		return acc.apps, nil
	}
	accounts, err := acc.Accounts()
	if err != nil {
		return nil, err
	}
	appDataAccessor, err := acc.accessors.Apps()
	if err != nil {
		return nil, err
	}
	acc.apps = appaccessor.NewAppAccessor(appDataAccessor, accounts)
	return acc.apps, nil
}

func (acc *AllAccessor[AccType]) Zones() (*zoneaccessor.ZoneAccessor, error) {
	if acc.zones != nil {
		return acc.zones, nil
	}
	accounts, err := acc.Accounts()
	if err != nil {
		return nil, err
	}
	zoneDataAccessor, err := acc.accessors.Zones()
	if err != nil {
		return nil, err
	}
	acc.zones = zoneaccessor.NewZoneAccessor(zoneDataAccessor, accounts)
	return acc.zones, nil
}

func (acc *AllAccessor[AccType]) Sources(factories []adsourceaccessor.SourceFactory, opts ...adsourceaccessor.Option[AccType]) (*adsourceaccessor.Accessor[AccType], error) {
	accounts, err := acc.Accounts()
	if err != nil {
		return nil, err
	}
	sourceDataAccessor, err := acc.accessors.RTBSources()
	if err != nil {
		return nil, err
	}
	return adsourceaccessor.NewAccessor(acc.ctx, sourceDataAccessor, accounts, factories, opts...)
}

func (acc *AllAccessor[AccType]) TrafficRouters() (*trafficrouteraccessor.TrafficRouterAccessor, error) {
	if acc.trafficRouters != nil {
		return acc.trafficRouters, nil
	}
	trafficRouterDataAccessor, err := acc.accessors.TrafficRouters()
	if err != nil {
		return nil, err
	}
	acc.trafficRouters = trafficrouteraccessor.NewTrafficRouterAccessor(trafficRouterDataAccessor)
	return acc.trafficRouters, nil
}

func (acc *AllAccessor[AccType]) Conn() any {
	return acc.accessors
}
