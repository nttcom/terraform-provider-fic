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

	destination_interconnect = "ECL-OSA-GU-Z1-01"
	destination_qos_type = "guarantee"
	destination_ecl_tenant_id = "%s"
	destination_ecl_api_key = "%s"
	destination_ecl_api_secret_key = "%s"

	bandwidth = "10M"

	primary_connected_network_address = "192.168.0.0/30"
	secondary_connected_network_address = "192.168.1.0/30"
}
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

* `source_route_filter_in` - (Required) Egress value of BGP Filter. 
  Allowed values are "fullRoute", "noRoute" and "fullRouteWithDefaultRoute" .

* `destination_port_id` - (Required) Destination port ID of the connection.

* `destination_vlan` - (Required) Destination VLAN ID of the connection.

* `bandwidth` - (Optional) Bandwidth of the connection. 
  Allowed values are "10M", "20M", "30M", "40M", "50M", "100M", "200M", "300M", "400M", "500M",
  "1G", "2G", "3G", "4G", "5G" and "10G" .


* `destination_interconnect` - (Required) Target cloud of the connection. 
  Either "JP3-1" or "JP5-1" is allowed.

* `destination_qos_type` - (Required) QOS Type of the connection.
  Currently only "guarantee" is supported.

* `destination_uno_tenant_id` - (Required) Target tenant ID of Enterprise Cloud.

* `destination_uno_api_key` - (Required) Your API Key in Enterprise Cloud service.

* `destination_uno_api_secret_key` - (Required) Your API Secret Key in Enterprise Cloud service.

* `primary_connected_network_address` - (Required) Primary network address of the connection.

* `secondary_connected_network_address` - (Required) Secondary network address of the connection.

* `bandwidth` - (Required) Bandwidth of the connection.
  Allowed values are "10M", "20M", "30M", "40M", "50M", "100M",
  "200M", "300M", "400M", "500M" and "1G" .


## Attributes Reference

The following attributes are exported:

* `redundant` - Redundancy of the connection.
* `tenant_id` - Tenant ID of the connection.
* `area` - Area name of the connection.

