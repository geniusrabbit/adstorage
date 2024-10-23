package appaccessor

import (
	"strings"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/models"

	"github.com/geniusrabbit/adstorage/accessors/accountaccessor"
	"github.com/geniusrabbit/adstorage/accessors/generalaccessor"
	"github.com/geniusrabbit/adstorage/loader"
)

// AppAccessor provides accessor to the admodel app type
type AppAccessor struct {
	generalaccessor.DataAccessor[*admodels.Application, string, models.Application]
}

// NewAppAccessor from dataAccessor
func NewAppAccessor[AccType any](dataAccessor loader.DataAccessor[models.Application], accountAccessor *accountaccessor.AccountAccessor[AccType]) *AppAccessor {
	return &AppAccessor{
		DataAccessor: *generalaccessor.NewDataAccessor(
			dataAccessor,
			func(app *models.Application) (*admodels.Application, bool) {
				acc, _ := accountAccessor.AccountByID(app.AccountID)
				napp := admodels.ApplicationFromModel(app)
				napp.Account = acc
				return &napp, true
			},
		),
	}
}

// AppList returns list of prepared data
func (acc *AppAccessor) AppList() ([]*admodels.Application, error) {
	return acc.List()
}

// AppByURI returns application object with specific URI
// If not found, try to find by part of URI
func (acc *AppAccessor) AppByURI(uri string) (*admodels.Application, error) {
	app, err := acc.ByKey(uri)
	if err != nil {
		return nil, err
	}
	if app != nil {
		return app, nil
	}
	for {
		idx := strings.IndexByte(uri, '.')
		if idx < 0 {
			break
		}
		uri = uri[idx+1:]
		if app, err = acc.ByKey(uri); err != nil {
			return nil, err
		}
		if app != nil && app.Type != types.ApplicationSite {
			return app, nil
		}
	}
	return nil, nil
}
