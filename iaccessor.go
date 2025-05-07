package adstorage

import (
	"github.com/geniusrabbit/adcorelib/models"

	"github.com/geniusrabbit/adstorage/loader"
)

// AllDataAccessor provides access to the data loader
type AllDataAccessor[AccType any] interface {
	Formats() (loader.DataAccessor[models.Format], error)
	Accounts() (loader.DataAccessor[AccType], error)
	Apps() (loader.DataAccessor[models.Application], error)
	Zones() (loader.DataAccessor[models.Zone], error)
	RTBSources() (loader.DataAccessor[models.RTBSource], error)
	TrafficRouters() (loader.DataAccessor[models.TrafficRouter], error)
}
