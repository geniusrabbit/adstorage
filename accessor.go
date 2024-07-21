package adstorage

import (
	"context"

	"github.com/geniusrabbit/adcorelib/accesspoint"
	"github.com/geniusrabbit/adcorelib/admodels/types"

	"github.com/geniusrabbit/adstorage/accessors/accesspointaccessor"
	"github.com/geniusrabbit/adstorage/accessors/accountaccessor"
	"github.com/geniusrabbit/adstorage/accessors/adsourceaccessor"
	"github.com/geniusrabbit/adstorage/accessors/campaignaccessor"
	"github.com/geniusrabbit/adstorage/accessors/formataccessor"
	"github.com/geniusrabbit/adstorage/accessors/zoneaccessor"
)

type AllAccessor[AccType any] struct {
	ctx       context.Context
	accessors AllDataAccessor[AccType]

	formats     types.FormatsAccessor
	accountCast accountaccessor.AccountConvertFunc[AccType]
	accounts    *accountaccessor.AccountAccessor[AccType]
	zones       *zoneaccessor.ZoneAccessor
	campaigns   *campaignaccessor.CampaignAccessor
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

func (acc *AllAccessor[AccType]) Campaigns() (*campaignaccessor.CampaignAccessor, error) {
	if acc.campaigns != nil {
		return acc.campaigns, nil
	}
	formats, err := acc.Formats()
	if err != nil {
		return nil, err
	}
	accounts, err := acc.Accounts()
	if err != nil {
		return nil, err
	}
	campaignDataAccessor, err := acc.accessors.Campaigns()
	if err != nil {
		return nil, err
	}
	acc.campaigns = campaignaccessor.NewCampaignAccessor(campaignDataAccessor, accounts, formats)
	return acc.campaigns, nil
}

func (acc *AllAccessor[AccType]) Sources(factories ...adsourceaccessor.SourceFactory) (*adsourceaccessor.Accessor[AccType], error) {
	accounts, err := acc.Accounts()
	if err != nil {
		return nil, err
	}
	sourceDataAccessor, err := acc.accessors.RTBSources()
	if err != nil {
		return nil, err
	}
	return adsourceaccessor.NewAccessor(acc.ctx, sourceDataAccessor, accounts, factories...)
}

func (acc *AllAccessor[AccType]) AccessPoints(factoryList ...accesspoint.Factory) (*accesspointaccessor.Accessor, error) {
	accounts, err := acc.Accounts()
	if err != nil {
		return nil, err
	}
	accesspointDataAccessor, err := acc.accessors.RTBAccessPoints()
	if err != nil {
		return nil, err
	}
	return accesspointaccessor.NewAccessor(
		acc.ctx,
		accesspointDataAccessor,
		accounts,
		factoryList...,
	)
}
