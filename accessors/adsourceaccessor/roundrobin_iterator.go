// thread-safeless

package adsourceaccessor

import (
	"iter"
	"math/rand"

	"github.com/geniusrabbit/adcorelib/adtype"
)

// NewRoundrobinIterator from request and source
func NewRoundrobinIterator(request *adtype.BidRequest, sources []adtype.Source) iter.Seq2[float32, adtype.Source] {
	startIndex := rand.Int() % len(sources)
	return (&linearIterator{
		started:  false,
		index:    startIndex,
		endIndex: startIndex,
	}).setSourcesExt(request, sources).seq()
}
