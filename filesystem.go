// package storage provides access to the data loader for filesystem storage
//
// FileSystem storage is used to store data in the filesystem.
// Structure of the data is stored in the following way:
//
// - format*.json - list of formats
//   - {format: {...single object...}, formats: [{...list of objects...}]
//
// - account*.json - list of accounts
//   - {account: {...single object...}, accounts: [{...list of objects...}]
//
// - campaign*.json - list of campaigns
//   - {campaign: {...single object...}, campaigns: [{...list of objects...}]
//
// - zone*.json - list of zones
//   - {zone: {...single object...}, zones: [{...list of objects...}]
//
// - rtb_source*.json - list of RTB sources
//   - {rtb_source: {...single object...}, rtb_sources: [{...list of objects...}]
//
// - rtb_access_point*.json - list of RTB access points
//   - {rtb_access_point: {...single object...}, rtb_access_points: [{...list of objects...}]
//
// Each type of objects from the list of different files will be merged into a single list
package adstorage

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/demdxx/gocast/v2"
	"github.com/geniusrabbit/adcorelib/models"

	"github.com/geniusrabbit/adstorage/loader"
	"github.com/geniusrabbit/adstorage/loader/fsloader"
	"github.com/geniusrabbit/adstorage/types/fstypes"
)

// FSDataAccessor provides access to the data loader for filesystem storage
func FSDataAccessor[AccType any](ctx context.Context, u *url.URL) any {
	rootDir := strings.Split(u.String()[5:], "?")[0]
	period, _ := time.ParseDuration(u.Query().Get("interval"))
	return &fsAccessor[AccType]{
		rootDir: rootDir,
		period:  gocast.IfThen(period == 0, time.Minute*5, period),
	}
}

type fsAccessor[AccType any] struct {
	rootDir string
	period  time.Duration
}

func (acc *fsAccessor[AccType]) Formats() (loader.DataAccessor[models.Format], error) {
	return loader.NewPeriodicReloader(
		&fstypes.FormatData{},
		fsloader.PatternLoader(acc.rootDir, "format*"),
		acc.period, acc.period*10, "loader_formats"), nil
}

func (acc *fsAccessor[AccType]) Accounts() (loader.DataAccessor[AccType], error) {
	return loader.NewPeriodicReloader(
		&fstypes.AccountData[*AccType]{},
		fsloader.PatternLoader(acc.rootDir, "account*"),
		acc.period, acc.period*10, "loader_accounts"), nil
}

func (acc *fsAccessor[AccType]) Campaigns() (loader.DataAccessor[models.Campaign], error) {
	return loader.NewPeriodicReloader(
		&fstypes.CampaignData{},
		fsloader.PatternLoader(acc.rootDir, "campaign*"),
		acc.period, acc.period*10, "loader_campaigns"), nil
}

func (acc *fsAccessor[AccType]) Apps() (loader.DataAccessor[models.Application], error) {
	return loader.NewPeriodicReloader(
		&fstypes.ApplicationData{},
		fsloader.PatternLoader(acc.rootDir, "app*"),
		acc.period, acc.period*10, "loader_apps"), nil
}

func (acc *fsAccessor[AccType]) Zones() (loader.DataAccessor[models.Zone], error) {
	return loader.NewPeriodicReloader(
		&fstypes.ZoneData{},
		fsloader.PatternLoader(acc.rootDir, "zone*"),
		acc.period, acc.period*10, "loader_zones"), nil
}

func (acc *fsAccessor[AccType]) RTBSources() (loader.DataAccessor[models.RTBSource], error) {
	return loader.NewPeriodicReloader(
		&fstypes.RTBSourceData{},
		fsloader.PatternLoader(acc.rootDir, "rtb_source*"),
		acc.period, acc.period*10, "loader_sources"), nil
}

func (acc *fsAccessor[AccType]) RTBAccessPoints() (loader.DataAccessor[models.RTBAccessPoint], error) {
	return loader.NewPeriodicReloader(
		&fstypes.RTBAccessPointData{},
		fsloader.PatternLoader(acc.rootDir, "rtb_access_point*"),
		acc.period, acc.period*10, "loader_access_points"), nil
}

func (acc *fsAccessor[AccType]) TrafficRouters() (loader.DataAccessor[models.TrafficRouter], error) {
	return loader.NewPeriodicReloader(
		&fstypes.TrafficRouterData{},
		fsloader.PatternLoader(acc.rootDir, "traffic_router*"),
		acc.period, acc.period*10, "loader_traffic_routers"), nil
}
