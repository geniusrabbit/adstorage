package fstypes

import "github.com/geniusrabbit/adcorelib/models"

// TrafficRouterData represents the data structure for traffic routers in files.
type TrafficRouterData struct {
	Routers []*models.TrafficRouter `json:"traffic_routers" yaml:"traffic_routers"`
	Router  *models.TrafficRouter   `json:"traffic_router" yaml:"traffic_router"`
}

// Merge combines the current data with another TrafficRouterData instance.
func (d *TrafficRouterData) Merge(it any) {
	v := it.(*TrafficRouterData)
	if v.Router != nil {
		d.Routers = append(d.Routers, v.Router)
	}
	d.Routers = append(d.Routers, v.Routers...)
}

// Result returns the list of traffic routers.
func (d *TrafficRouterData) Result() []*models.TrafficRouter {
	return d.Routers
}

// Reset clears the stored data.
func (d *TrafficRouterData) Reset() {
	d.Router = nil
	d.Routers = d.Routers[:0]
}
