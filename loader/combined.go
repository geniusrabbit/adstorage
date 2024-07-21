package loader

type MergeFnk[T any] func(datas ...[]*T) []*T

// CombinedLoader loads and merge data into other type
type CombinedLoader[T any] struct {
	loaders []DataAccessor[T]
	merge   MergeFnk[T]
}

// NewCombinedLoader returns combined implementation of dataloader
func NewCombinedLoader[T any](merge MergeFnk[T], loaders ...DataAccessor[T]) *CombinedLoader[T] {
	return &CombinedLoader[T]{
		loaders: loaders,
		merge:   merge,
	}
}

func (l *CombinedLoader[T]) NeedUpdate() bool {
	for _, lr := range l.loaders {
		if lr.NeedUpdate() {
			return true
		}
	}
	return false
}

// Data returns loaded data and reload if necessary
func (l *CombinedLoader[T]) Data() ([]*T, error) {
	datas := make([][]*T, 0, len(l.loaders))
	for _, lr := range l.loaders {
		dr, err := lr.Data()
		if err != nil {
			return nil, err
		}
		datas = append(datas, dr)
	}
	return l.merge(datas...), nil
}
