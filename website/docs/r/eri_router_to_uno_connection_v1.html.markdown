---
layout: "fic"
page_title: "Flexible InterConnect: fic_eri_router_to_uno_connection_v1"
sidebar_current: "docs-fic-resource-eri-router-to-uno-connection-v1"
description: |-
  Manages a V1 Router to ECL Connection resource within Flexible InterConnect.
---

# fic\_eri\_router\_to\_uno\_connection\_v1

Manages a V1 Router to UNO Connection resource within Flexible InterConnect.

## Example Usage

### Basic Usage

```hcl
resource "fic_eri_router_v1" "router_1" {
  name            = "terraform_router_1"
  area            = "JPEAST"
  user_ip_address = "10.0.0.0/27"
  redundant       = false
}

resource "fic_eri_router_to_uno_connection_v1" "connection_1" {
  name = "terraform_connection_1"

  source_router_id        = fic_eri_router_v1.router_1.id
  source_group_name       = "group_1"
  source_route_filter_in  = "noRoute"
  source_route_filter_out = "fullRouteWithDefaultRoute"

  destination_interconnect           = "Interconnect-Osaka-1"
  destination_c_number               = "C0250124868"
  destination_parent_contract_number = "N190005036"
  destination_vpn_number             = "V19000708"
  destination_qos_type               = "guarantee"
  destination_route_filter_out       = "fullRoute"

  connected_network_address = "192.168.0.0/29"

  bandwidth = "10M"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the resource.

* `source_router_id` - (Required) Source router ID of the connection.

* `source_group_name` - (Required) Source group name of the connection.
  Allowed values are "group_1", "group_2", "group_3", "group_4", "group_5", "group_6", "group_7" and "group_8".

* `source_route_filter_in` - (Required) Source ingress value of BGP Filter.
  Either "fullRoute" or "noRoute" is allowed.

* `source_route_filter_out` - (Required) Source egress value of BGP Filter.
  Allowed values are "fullRoute", "fullRouteWithDefaultRoute" and "noRoute".

* `destination_interconnect` - (Required) Target cloud of the connection.
  Either "Interconnect-Tokyo-1" or "Interconnect-Osaka-1" is allowed.

* `destination_c_number` - (Optional) Destination C number of the connection.

* `destination_parent_contract_number` - (Required) Destination parent contract number of the connection.

* `destination_vpn_number` - (Required) Destination VPN number of the connection.

* `destination_qos_type` - (Required) QOS Type of the connection.
  Currently only "guarantee" is supported.

* `destination_route_filter_out` - (Required) Destination egress value of BGP Filter.
  Allowed values are "fullRoute", "fullRouteWithDefaultRoute", "defaultRoute" and "privateRoute".

* `connected_network_address` - (Required) Network address of the connection.

* `bandwidth` - (Required) Bandwidth of the connection.
  Allowed values are "10M", "20M", "30M", "40M", "50M", "100M", "200M", "300M", "400M", "500M" and "1G" .

## Attributes Reference

The following attributes are exported:

* `destination_contract_number` - Destination contract number of the connection.

* `redundant` - Redundancy of the connection.

* `tenant_id` - Tenant ID of the connection.

* `area` - Area name of the connection.

