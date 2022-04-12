---
layout: "fic"
page_title: "Flexible InterConnect: fic_eri_router_single_to_port_connection_v1"
sidebar_current: "docs-fic-resource-eri-router-single-to-port-connection-v1"
description: |-
  Manages a V1 Router(Single) to Port Connection resource within Flexible InterConnect.
---

# fic\_eri\_router\_single\_to\_port\_connection\_v1
1
Manages a V1 Router to Port(Single) Connection resource within Flexible InterConnect.

## Example Usage

### Basic Usage

```hcl
resource "fic_eri_router_v1" "router_1" {
  name            = "terraform_router_1"
  area            = "JPEAST"
  user_ip_address = "10.0.0.0/27"
  redundant       = false
}

resource "fic_eri_port_v1" "port_1" {
  name        = "terraform_port_1"
  switch_name = "lxea01comnw1"
  port_type   = "10G"

  vlan_ranges {
    start = 1137
    end   = 1152
  }
}

resource "fic_eri_router_single_to_port_connection_v1" "connection_1" {
  name              = "terraform_connection_1"
  source_router_id  = fic_eri_router_v1.router_1.id
  source_group_name = "group_1"

  source_information {
    ip_address          = "10.0.1.1/30"
    as_path_prepend_in  = "4"
    as_path_prepend_out = "4"
  }

  source_route_filter_in  = "fullRoute"
  source_route_filter_out = "fullRouteWithDefaultRoute"

  destination_information {
    port_id    = fic_eri_port_v1.port_1.id
    vlan       = fic_eri_port_v1.port_1.vlan_ranges[0].start
    ip_address = "10.0.1.2/30"
    asn        = "65000"
  }

  bandwidth = "10M"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the resource.

* `source_router_id` - (Required) Source router ID of the connection.

* `source_group_name` - (Required) Source group name of the connection.
  Allowed values are "group_1", "group_2", "group_3", "group_4",
"group_5", "group_6", "group_7" and "group_8".

* `destination_port_id` - (Required) Destination port ID of the connection.

* `destination_vlan` - (Required) Destination VLAN ID of the connection.

* `bandwidth` - (Optional) Bandwidth of the connection. 
  Allowed values are "10M", "20M", "30M", "40M", "50M", "100M", "200M", "300M", "400M", "500M",
  "1G", "2G", "3G", "4G", "5G" and "10G" .

* `source_route_filter_in` - (Required) Ingress value of BGP Filter. 
  Either "fullRoute" or "noRoute" is allowed.

* `source_route_filter_out` - (Required) Egress value of BGP Filter. 
  Allowed values are "fullRoute", "noRoute" and "fullRouteWithDefaultRoute" .

* `source_information` - (Required) List of source information. 
  Length of list must be 1(means primary and secondary).

* `destination_information` - (Required) List of destination information. 
  Length of list must be 1(means primary and secondary).

* `bandwidth` - (Required) Bandwidth of the connection.
  Allowed values are "10M", "20M", "30M", "40M", "50M", "100M",
  "200M", "300M", "400M", "500M", "1G", "2G", "3G", "4G", 
  "5G" and "10G" .

The `source_information` block supports:

* `ip_address` - (Required) Source IP Address.
* `as_path_prepend_in` - (Required) Source AS Path Prepend for ingress.
* `as_path_prepend_out` - (Required) Source AS Path Prepend for Egress.

The `destination_information` block supports:

* `ip_address` - (Required) Destination port ID.
* `vlan` - (Required) Destination VLAN ID.
* `ip_address` - (Required) Destination IP Address.
* `asn` - (Required) Destination ASN.

## Attributes Reference

The following attributes are exported:

* `redundant` - Redundancy of the connection.
* `tenant_id` - Tenant ID of the connection.
* `area` - Area name of the connection.

