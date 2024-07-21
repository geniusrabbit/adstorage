package adstorage

import (
	"context"
	"fmt"
	"net/url"

	"github.com/geniusrabbit/adstorage/accessors/accountaccessor"
)

// DataLoaderConnectFnk returns general type data accessor
type DataLoaderConnectFnk func(ctx context.Context, u *url.URL) any

var dataLoaderAccessor = map[string]DataLoaderConnectFnk{}

// Connect to the dataLoader accessor
func Connect[AccType any](ctx context.Context, urlConnect string) (AllDataAccessor[AccType], error) {
	u, err := url.Parse(urlConnect)
	if err != nil {
		return nil, err
	}
	accessor := dataLoaderAccessor[u.Scheme]
	if accessor == nil {
		return nil, fmt.Errorf("unsupported data accessor type [%s]", u.Scheme)
	}
	return accessor(ctx, u).(AllDataAccessor[AccType]), nil
}

// ConnectAllAccessors to the dataLoader accessor
func ConnectAllAccessors[AccType any](ctx context.Context, urlConnect string, accountCast accountaccessor.AccountConvertFunc[AccType]) (*AllAccessor[AccType], error) {
	allDataAccessors, err := Connect[AccType](ctx, urlConnect)
	if err != nil {
		return nil, err
	}
	return NewAllAccessor(ctx, allDataAccessors, accountCast), nil
}

// Register new data accessor
func Register[AccType any](scheme string, fnk DataLoaderConnectFnk) {
	dataLoaderAccessor[scheme] = fnk
}
