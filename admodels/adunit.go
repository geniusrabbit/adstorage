package admodels

import (
	"github.com/geniusrabbit/gosql/v2"

	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/billing"
	"github.com/geniusrabbit/adcorelib/models"
)

// AdUnit model represents ad unit spot for ad placement
type AdUnit struct {
	id            uint64
	CodenameValue string

	Acc   adtype.Account
	AccID uint64

	FixedPurchasePrice billing.Money // The cost of single view
	MinECPM            float64

	AllowedTypes      gosql.NullableOrderedNumberArray[uint64]
	AllowedSources    gosql.NullableOrderedNumberArray[uint64]
	DisallowedSources gosql.NullableOrderedNumberArray[uint64]
	DefaultCode       map[string]string
}

// AdUnitFromModel convert database model to specified model
func AdUnitFromModel(zone *models.Zone, account adtype.Account) *AdUnit {
	return &AdUnit{
		id:                 zone.ID,
		CodenameValue:      zone.Codename,
		FixedPurchasePrice: billing.MoneyFloat(zone.FixedPurchasePrice),
		Acc:                account,
		AccID:              zone.AccountID,
		MinECPM:            zone.MinECPM,
		AllowedTypes:       zone.AllowedTypes,
		AllowedSources:     zone.AllowedSources,
		DisallowedSources:  zone.DisallowedSources,
		DefaultCode:        zone.DefaultCode.DataOr(nil),
	}
}

// ID of object
func (z *AdUnit) ID() uint64 {
	return z.id
}

// Codename of the target (equal to tagid)
func (z *AdUnit) Codename() string {
	return z.CodenameValue
}

// ObjectKey of the target
func (z *AdUnit) ObjectKey() string {
	return z.CodenameValue
}

// PricingModel of the target
func (z *AdUnit) PricingModel() types.PricingModel {
	return types.PricingModelUndefined
}

// AlternativeAdCode returns URL or any code (HTML, XML, etc)
func (z *AdUnit) AlternativeAdCode(key string) string {
	if z.DefaultCode == nil {
		return ""
	}
	return z.DefaultCode[key]
}

// PurchasePrice gives the price of view from external resource
func (z *AdUnit) PurchasePrice(action adtype.Action) billing.Money {
	if action.IsImpression() {
		return z.FixedPurchasePrice
	}
	return 0
}

// Account object
func (z *AdUnit) Account() adtype.Account {
	return z.Acc
}

// AccountID of current target
func (z *AdUnit) AccountID() uint64 {
	return z.AccID
}

// SetAccount for target
func (z *AdUnit) SetAccount(acc adtype.Account) {
	z.Acc = acc
}

// CommissionShareFactor which system get from publisher
func (z *AdUnit) CommissionShareFactor() float64 {
	return z.Acc.CommissionShareFactor()
}

// RevenueShareReduceFactor correction factor to reduce target proce of the access point to avoid descrepancy
// Returns percent from 0 to 1 for reducing of the value
// If there is 10% of price correction, it means that 10% of the final price must be ignored
func (z *AdUnit) RevenueShareReduceFactor() float64 { return 0 }

// IsAllowedSource for targeting
func (z *AdUnit) IsAllowedSource(id uint64, types []int) bool {
	if len(z.AllowedSources) > 0 {
		return z.AllowedSources.IndexOf(id) >= 0
	}
	if len(z.DisallowedSources) > 0 {
		return z.DisallowedSources.IndexOf(id) < 0
	}
	return true
}
