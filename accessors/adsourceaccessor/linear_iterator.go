package adsourceaccessor

import (
	"iter"

	"github.com/demdxx/xtypes"
	"github.com/geniusrabbit/adcorelib/adtype"
)

type sourceItem struct {
	priority float32
	source   adtype.Source
}

type linearIterator struct {
	started  bool
	index    int
	endIndex int
	sources  []sourceItem
}

// NewLinearIterator from request and source
func NewLinearIterator(request *adtype.BidRequest, sources []adtype.Source) iter.Seq2[float32, adtype.Source] {
	return (&linearIterator{
		started:  false,
		index:    0,
		endIndex: 0,
	}).setSourcesExt(request, sources).seq()
}

func (iter *linearIterator) setSourcesExt(request *adtype.BidRequest, sources []adtype.Source) *linearIterator {
	iter.sources = xtypes.SliceApply(
		xtypes.Slice[adtype.Source](sources).Filter(func(src adtype.Source) bool {
			return src != nil && request.SourceFilterCheck(src.ID()) && src.Test(request)
		}),
		func(src adtype.Source) sourceItem {
			return sourceItem{
				priority: 1,
				source:   src,
			}
		})
	iter.endIndex = len(iter.sources)
	return iter
}

func (iter *linearIterator) seq() iter.Seq2[float32, adtype.Source] {
	return func(yield func(float32, adtype.Source) bool) {
		for {
			prior, src := iter.Next()
			if src == nil {
				break
			}
			if !yield(prior, src) {
				return
			}
		}
	}
}

// Next returns next source
func (iter *linearIterator) Next() (float32, adtype.Source) {
	if len(iter.sources) == 0 {
		return 0, nil
	}
	if iter.index >= len(iter.sources) {
		iter.index = 0
	}
	if iter.index == iter.endIndex && iter.started {
		return 0, nil
	}
	src := &iter.sources[iter.index]
	iter.index++
	iter.started = true
	return src.priority, src.source
}
