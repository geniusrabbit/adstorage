package accountaccessor

import (
	"github.com/geniusrabbit/adcorelib/admodels"

	"github.com/geniusrabbit/adstorage/accessors/generalaccessor"
	"github.com/geniusrabbit/adstorage/loader"
)

type AccountConvertFunc[AT any] generalaccessor.ObjectConvertFunc[AT, *admodels.Account]

// AccountAccessor provides accessor to the admodel company type
type AccountAccessor[AT any] struct {
	generalaccessor.DataAccessor[*admodels.Account, AT]
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
func (acc *AccountAccessor[AT]) CompanyList() ([]*admodels.Account, error) {
	return acc.List()
}

// AccountByID returns account object with specific ID
func (acc *AccountAccessor[AT]) AccountByID(id uint64) (*admodels.Account, error) {
	return acc.ByID(id)
}
