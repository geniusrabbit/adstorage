package accountaccessor

import (
	"context"

	"github.com/geniusrabbit/adcorelib/admodels"

	"github.com/geniusrabbit/adstorage/accessors/generalaccessor"
	"github.com/geniusrabbit/adstorage/loader"
)

type AccountConvertFunc[AT any] generalaccessor.ObjectConvertFunc[AT, *admodels.Account]

// AccountAccessor provides accessor to the admodel company type
type AccountAccessor[AT any] struct {
	generalaccessor.DataAccessor[*admodels.Account, uint64, AT]
}

// NewAccessor from dataAccessor
func NewAccessor[AT any](dataAccessor loader.DataAccessor[AT], accountConvert AccountConvertFunc[AT]) *AccountAccessor[AT] {
	return &AccountAccessor[AT]{
		DataAccessor: *generalaccessor.NewDataAccessor(
			dataAccessor,
			generalaccessor.ObjectConvertFunc[AT, *admodels.Account](accountConvert),
		),
	}
}

// CompanyList returns list of prepared data
func (acc *AccountAccessor[AT]) CompanyList(ctx context.Context) ([]*admodels.Account, error) {
	return acc.List(ctx)
}

// AccountByID returns account object with specific ID
func (acc *AccountAccessor[AT]) AccountByID(ctx context.Context, id uint64) (*admodels.Account, error) {
	return acc.ByKey(ctx, id)
}
