package fic

import (
	connections "github.com/nttcom/go-fic/fic/eri/v1/router_paired_to_gcp_connections"
)

func expandSource(in []interface{}) connections.Source {
	m := in[0].(map[string]interface{})
	primaryMEDOut := m["primary_med_out"].(int)

	return connections.Source{
		RouterID:    m["router_id"].(string),
		GroupName:   m["group_name"].(string),
		RouteFilter: expandRouteFilter(m["route_filter"].([]interface{})),
		Primary: connections.SourceHAInfo{
			MED: connections.MED{
				Out: primaryMEDOut,
			},
		},
		Secondary: connections.SourceHAInfo{
			MED: connections.MED{
				Out: primaryMEDOut + 10,
			},
		},
	}
}

func expandRouteFilter(in []interface{}) connections.RouteFilter {
	m := in[0].(map[string]interface{})

	return connections.RouteFilter{
		In:  m["in"].(string),
		Out: m["out"].(string),
	}
}

func expandDestination(in []interface{}) connections.Destination {
	m := in[0].(map[string]interface{})

	return connections.Destination{
		QosType:   "guarantee",
		Primary:   expandInterconnect(m["primary"].([]interface{})),
		Secondary: expandInterconnect(m["secondary"].([]interface{})),
	}
}

func expandInterconnect(in []interface{}) connections.DestinationHAInfo {
	m := in[0].(map[string]interface{})

	return connections.DestinationHAInfo{
		Interconnect: m["interconnect"].(string),
		PairingKey:   m["pairing_key"].(string),
	}
}

func flattenSource(in connections.Source) []interface{} {
	var out []interface{}
	m := make(map[string]interface{})

	m["router_id"] = in.RouterID
	m["group_name"] = in.GroupName
	m["route_filter"] = flattenRouteFilter(in.RouteFilter)
	m["primary_med_out"] = in.Primary.MED.Out
	m["secondary_med_out"] = in.Secondary.MED.Out

	out = append(out, m)
	return out
}

func flattenRouteFilter(in connections.RouteFilter) []interface{} {
	var out []interface{}
	m := make(map[string]interface{})

	m["in"] = in.In
	m["out"] = in.Out

	out = append(out, m)
	return out
}

func flattenDestination(in connections.Destination) []interface{} {
	var out []interface{}
	m := make(map[string]interface{})

	m["primary"] = flattenInterconnect(in.Primary)
	m["secondary"] = flattenInterconnect(in.Secondary)
	m["qos_type"] = in.QosType

	out = append(out, m)
	return out
}

func flattenInterconnect(in connections.DestinationHAInfo) []interface{} {
	var out []interface{}
	m := make(map[string]interface{})

	m["interconnect"] = in.Interconnect
	m["pairing_key"] = in.PairingKey

	out = append(out, m)
	return out
}
