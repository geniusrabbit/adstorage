package adsourceaccessor

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/context/ctxlogger"
	"github.com/geniusrabbit/adcorelib/models"
	"github.com/geniusrabbit/adcorelib/platform/info"

	"github.com/geniusrabbit/adstorage/accessors/accountaccessor"
	"github.com/geniusrabbit/adstorage/accessors/generalaccessor"
	"github.com/geniusrabbit/adstorage/loader"
)

var errUnsupportedSourceProtocol = errors.New("unsupported source protocol")

type SourceFactory interface {
	New(ctx context.Context, source *admodels.RTBSource, opts ...any) (adtype.SourceTester, error)
	Info() info.Platform
	Protocols() []string
}

// Accessor object ad reloader
type Accessor[AccType any] struct {
	generalaccessor.DataAccessor[adtype.Source, models.RTBSource]

	factories   map[string]SourceFactory
	factoryList []SourceFactory
}

// NewAccessor object
func NewAccessor[AccType any](ctx context.Context, dataAccessor loader.DataAccessor[models.RTBSource], accountAccessor *accountaccessor.AccountAccessor[AccType], factories ...SourceFactory) (*Accessor[AccType], error) {
	if dataAccessor == nil {
		return nil, errors.New("data accessor is required")
	}
	if accountAccessor == nil {
		return nil, errors.New("account accessor is required")
	}
	mapFactory := map[string]SourceFactory{}
	for _, fact := range factories {
		for _, protoName := range fact.Protocols() {
			mapFactory[protoName] = fact
		}
	}
	accessor := &Accessor[AccType]{
		factories:   mapFactory,
		factoryList: factories,
	}
	accessor.DataAccessor = *generalaccessor.NewDataAccessor(
		dataAccessor,
		func(src *models.RTBSource) (adtype.Source, bool) {
			if src.AccountID == 0 {
				ctxlogger.Get(ctx).Error("source without account",
					zap.Uint64("source_id", src.ID))
				return nil, false
			}
			acc, _ := accountAccessor.AccountByID(src.AccountID)
			rtbSrc := admodels.RTBSourceFromModel(src, acc)
			if src, err := accessor.newSource(ctx, rtbSrc); err == nil {
				return src, true
			} else {
				ctxlogger.Get(ctx).Error("create RTB source",
					zap.Uint64("source_id", rtbSrc.ID),
					zap.String("source_protocol", rtbSrc.Protocol),
					zap.Error(err))
			}
			return nil, false
		},
	)
	return accessor, nil
}

// FactoryList returns list of source factories
func (acc *Accessor[AT]) FactoryList() []SourceFactory {
	return acc.factoryList
}

// SourceList returns list of sources
func (acc *Accessor[AT]) SourceList() ([]adtype.Source, error) {
	return acc.List()
}

// Iterator returns the configured queue acc
func (acc *Accessor[AT]) Iterator(request *adtype.BidRequest) adtype.SourceIterator {
	list, _ := acc.SourceList()
	return NewPriorityIterator(request, list)
}

// SourceByID returns source instance
func (acc *Accessor[AT]) SourceByID(id uint64) (adtype.Source, error) {
	return acc.ByID(id)
}

func (acc *Accessor[AT]) newSource(ctx context.Context, src *admodels.RTBSource) (adtype.Source, error) {
	if acc.factories == nil {
		return nil, errors.Wrap(errUnsupportedSourceProtocol, src.Protocol)
	}
	fact := acc.factories[src.Protocol]
	if fact == nil {
		return nil, errors.Wrap(errUnsupportedSourceProtocol, src.Protocol)
	}
	return fact.New(ctx, src)
}

// SetTimeout for sourcer
func (acc *Accessor[AT]) SetTimeout(timeout time.Duration) {
	list, _ := acc.SourceList()
	for _, src := range list {
		if srcSetTM, _ := src.(adtype.SourceTimeoutSetter); srcSetTM != nil {
			srcSetTM.SetTimeout(timeout)
		}
	}
}

var _ adtype.SourceAccessor = &Accessor[any]{}
