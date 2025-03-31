package adsourceaccessor

import (
	"iter"

	"github.com/demdxx/xtypes"
	"github.com/geniusrabbit/adcorelib/adtype"
)

// PriorityEval is a function type for evaluating the priority of a source
// based on the request and source.
// It returns a float32 value representing the priority.
// A higher value indicates a higher priority.
// The function should return 0 if the source is not applicable for the request.
type PriorityEval func(request *adtype.BidRequest, src adtype.Source) float32

// NewPriorityIterator from request and source
func NewPriorityIterator(request *adtype.BidRequest, sources []adtype.Source, priorEval PriorityEval) iter.Seq2[float32, adtype.Source] {
	return (&linearIterator{
		started:  false,
		index:    0,
		endIndex: 0,
		sources: xtypes.SliceApply(sources, func(src adtype.Source) sourceItem {
			return sourceItem{
				priority: priorEval(request, src),
				source:   src,
			}
		}).Filter(func(src sourceItem) bool {
			return src.priority > 0 && src.source != nil &&
				request.SourceFilterCheck(src.source.ID()) && src.source.Test(request)
		}),
	}).seq()
}
