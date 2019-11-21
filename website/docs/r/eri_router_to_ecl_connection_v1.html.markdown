---
layout: "fic"
page_title: "Flexible InterConnect: fic_eri_router_to_ecl_connection_v1"
sidebar_current: "docs-fic-resource-eri-router-to-ecl-connection-v1"
description: |-
  Manages a V1 Router to ECL Connection resource within Flexible InterConnect.
---

# fic\_eri\_router\_to\_ecl\_connection\_v1

Manages a V1 Router to ECL Connection resource within Flexible InterConnect.

## Example Usage

### Basic Usage

```hcl
resource "fic_eri_router_v1" "router_1" {
	name = "terraform_router_1"
	area = "JPEAST"
	user_ip_address = "10.0.0.0/27"
	redundant = false
}

resource "fic_eri_router_to_ecl_connection_v1" "connection_1" {
	name = "terraform_connection_1"
	source_router_id = "${fic_eri_router_v1.router_1.id}"

	source_group_name = "group_1"
	source_route_filter_in = "noRoute"
	source_route_filter_out = "noRoute"

	destination_interconnect = "JP5-1"
	destination_qos_type = "guarantee"
	destination_ecl_tenant_id = "215169c6e4fe4fb9b2c0ae12545a422c"
	destination_ecl_api_key = "62490c125f374dc484ea027c3fc9141b"
	destination_ecl_api_secret_key = "484ea027c3fc9141"

	bandwidth = "100M"

	primary_connected_network_address = "10.0.0.0/30"
	secondary_connected_network_address = "10.0.0.0/30"
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the resource.

* `source_router_id` - (Required) Source router ID of the connection.

* `source_group_name` - (Required) Source group name of the connection.
  Allowed values are "group_1", "group_2", "group_3", "group_4",
"group_5", "group_6", "group_7" and "group_8".

* `source_route_filter_in` - (Required) Ingress value of BGP Filter. 
  Either "fullRoute" or "noRoute" is allowed.

* `source_route_filter_out` - (Required) Egress value of BGP Filter. 
  Allowed values are "fullRoute", "noRoute" and "fullRouteWithDefaultRoute" .

* `destination_interconnect` - (Required) Target cloud of the connection. 
  Either "Interconnect-Tokyo-1" or "Interconnect-Osaka-1"" is allowed.

* `destination_c_number` - (Required) Destination C number of the connection.

* `destination_parent_contract_number` - (Required) 
  Destination parent contract number of the connection.

* `destination_vpn_number` - (Required) 
  Destination VPN number of the connection.

* `destination_qos_type` - (Required) QOS Type of the connection.
  Currently only "guarantee" is supported.

* `source_route_filter_out` - (Required) Egress value of BGP Filter. 
  Allowed values are "fullRoute", "fullRouteWithDefaultRoute", "defaultRoute" and "privateRoute".

* `connected_network_address` - (Required) Network address of the connection.

* `bandwidth` - (Optional) Bandwidth of the connection. 
  Allowed values are "10M", "20M", "30M", "40M", "50M", "100M", "200M", "300M", "400M", "500M" and "1G" .


## Attributes Reference

The following attributes are exported:

* `destination_contract_number` - 
  Destination contract number of the connection.
* `redundant` - Redundancy of the connection.
* `tenant_id` - Tenant ID of the connection.
* `area` - Area name of the connection.

