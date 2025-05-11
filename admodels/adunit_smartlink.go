package admodels

import (
	"github.com/geniusrabbit/gosql/v2"

	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/billing"
	"github.com/geniusrabbit/adcorelib/models"
)

// Smartlink model
type Smartlink struct {
	id            uint64
	CodenameValue string

	Acc   adtype.Account
	AccID uint64

	// The cost of single view
	FixedPurchasePrice billing.Money

	// Filtering
	AllowedTypes      gosql.NullableOrderedNumberArray[uint64]
	AllowedSources    gosql.NullableOrderedNumberArray[uint64]
	DisallowedSources gosql.NullableOrderedNumberArray[uint64]
	Campaigns         gosql.NullableOrderedNumberArray[uint64]

	// DefaultCode for the target for the specified format (banner, video, direct, etc)
	DefaultCode map[string]string
}

// SmartlinkFromModel convert database model to specified model
func SmartlinkFromModel(zone *models.Zone, account adtype.Account) *Smartlink {
	return &Smartlink{
		id:                 zone.ID,
		CodenameValue:      zone.Codename,
		FixedPurchasePrice: billing.MoneyFloat(zone.FixedPurchasePrice),
		Acc:                account,
		AccID:              zone.AccountID,
		AllowedTypes:       zone.AllowedTypes,
		AllowedSources:     zone.AllowedSources,
		DisallowedSources:  zone.DisallowedSources,
		Campaigns:          zone.Campaigns,
		DefaultCode:        *zone.DefaultCode.Data,
	}
}

// ID of object
func (l *Smartlink) ID() uint64 {
	return l.id
}

// Codename of the target (equal to tagid)
func (l *Smartlink) Codename() string {
	return l.CodenameValue
}

// ObjectKey of the target
func (l *Smartlink) ObjectKey() string {
	return l.CodenameValue
}

// PricingModel of the target
func (l *Smartlink) PricingModel() types.PricingModel {
	return types.PricingModelUndefined
}

// AlternativeAdCode returns URL or any code (HTML, XML, etc)
func (l *Smartlink) AlternativeAdCode(key string) string {
	if l.DefaultCode == nil {
		return ""
	}
	return l.DefaultCode[key]
}

// PurchasePrice gives the price of view from external resource
func (l *Smartlink) PurchasePrice(action adtype.Action) billing.Money {
	if action.IsImpression() {
		return l.FixedPurchasePrice
	}
	return 0
}

// Account object
func (l *Smartlink) Account() adtype.Account {
	return l.Acc
}

// AccountID of current target
func (l *Smartlink) AccountID() uint64 {
	return l.AccID
}

// SetAccount for target
func (l *Smartlink) SetAccount(acc adtype.Account) {
	l.Acc = acc
}

// CommissionShareFactor which system get from publisher
func (l *Smartlink) CommissionShareFactor() float64 {
	return l.Acc.CommissionShareFactor()
}

// RevenueShareReduceFactor correction factor to reduce target proce of the access point to avoid descrepancy
// Returns percent from 0 to 1 for reducing of the value
// If there is 10% of price correction, it means that 10% of the final price must be ignored
func (l *Smartlink) RevenueShareReduceFactor() float64 { return 0 }

// IsAllowedSource for targeting
func (l *Smartlink) IsAllowedSource(id uint64, types []int) bool {
	if len(l.AllowedSources) > 0 {
		return l.AllowedSources.IndexOf(id) >= 0
	}
	if len(l.DisallowedSources) > 0 {
		return l.DisallowedSources.IndexOf(id) < 0
	}
	return true
}
