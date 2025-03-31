package adsourceaccessor

// Option is a function type for configuring the Accessor
type Option[AccType any] func(m *Accessor[AccType])

// WithCustomIterator sets a custom iterator function for the Accessor
func WithCustomIterator[AccType any](fnk CustomIteratorFnk) Option[AccType] {
	return func(acc *Accessor[AccType]) {
		acc.customIterator = fnk
	}
}
