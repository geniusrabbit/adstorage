package adstorage

import (
	"context"
	"net/url"
	"time"

	"github.com/demdxx/gocast/v2"

	"github.com/geniusrabbit/adcorelib/models"

	"github.com/geniusrabbit/adstorage/database"
	"github.com/geniusrabbit/adstorage/loader"
	"github.com/geniusrabbit/adstorage/loader/dbloader"
	"github.com/geniusrabbit/adstorage/types/dbtypes"
)

// RegisterAllSchemas for all supported dialects
func RegisterAllSchemas[AccType any]() {
	for _, dialect := range database.ListOfDialects() {
		Register(dialect, DBDataAccessor[AccType])
	}
}

// DBDataAccessor provides access to the data loader for database storage
func DBDataAccessor[AccType any](ctx context.Context, u *url.URL) any {
	db, err := database.Connect(ctx, u.String(), gocast.Bool(u.Query().Get("debug")))
	if err != nil {
		panic(err)
	}
	period, _ := time.ParseDuration(u.Query().Get("interval"))
	if period == 0 {
		period = time.Minute * 5
	}
	return &dbAccessor[AccType]{db: db, period: period}
}

type dbAccessor[AccType any] struct {
	db     *database.DB
	period time.Duration
}

func (acc *dbAccessor[AccType]) Formats() (loader.DataAccessor[models.Format], error) {
	return loader.NewPeriodicReloader(&dbtypes.FormatList{},
		dbloader.Loader(acc.db), acc.period, acc.period*10, "loader_formats"), nil
}

func (acc *dbAccessor[AccType]) Accounts() (loader.DataAccessor[AccType], error) {
	return loader.NewPeriodicReloader(&dbtypes.AccountList[*AccType]{}, dbloader.Loader(acc.db),
		acc.period, acc.period*10, "loader_accounts"), nil
}

func (acc *dbAccessor[AccType]) Apps() (loader.DataAccessor[models.Application], error) {
	return loader.NewPeriodicReloader(&dbtypes.ApplicationList{}, dbloader.Loader(acc.db),
		acc.period, acc.period*10, "loader_apps"), nil
}

func (acc *dbAccessor[AccType]) Zones() (loader.DataAccessor[models.Zone], error) {
	return loader.NewPeriodicReloader(&dbtypes.ZoneList{}, dbloader.Loader(acc.db),
		acc.period, acc.period*10, "loader_zones"), nil
}

func (acc *dbAccessor[AccType]) RTBSources() (loader.DataAccessor[models.RTBSource], error) {
	return loader.NewPeriodicReloader(&dbtypes.RTBSourceList{}, dbloader.Loader(acc.db),
		acc.period, acc.period*10, "loader_sources"), nil
}

func (acc *dbAccessor[AccType]) TrafficRouters() (loader.DataAccessor[models.TrafficRouter], error) {
	return loader.NewPeriodicReloader(&dbtypes.TrafficRouterList{}, dbloader.Loader(acc.db),
		acc.period, acc.period*10, "loader_traffic_routers"), nil
}

//lint:ignore U1000 This method is used to ensure dbAccessor implements the required interface
func (*dbAccessor[AccType]) __isDB() {}

func (acc *dbAccessor[AccType]) getDB() *database.DB {
	return acc.db
}

func (acc *dbAccessor[AccType]) gerPeriod() time.Duration {
	return acc.period
}
