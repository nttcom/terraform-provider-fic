---
layout: "fic"
page_title: "Flexible InterConnect: fic_eri_router_to_azure_private_connection_v1"
sidebar_current: "docs-fic-resource-eri-router-to-azure-private-connection-v1"
description: |-
  Manages a V1 Router to Azure Private Connection resource within Flexible InterConnect.
---

# fic\_eri\_router\_to\_azure\_private\_connection\_v1

Manages a V1 Router to Azure Private Connection resource within Flexible InterConnect.

## Example Usage

### Basic Usage

```hcl
resource "fic_eri_router_v1" "router_1" {
    name = "terraform_router_1"
    area = "JPWEST"
    user_ip_address = "10.0.0.0/27"
    redundant = false
}

resource "fic_eri_router_to_azure_private_connection_v1" "connection_1" {
    name = "terraform_connection_1"

    source_router_id = "${fic_eri_router_v1.router_1.id}"
    source_group_name = "group_1"
    source_route_filter_in = "fullRoute"
    source_route_filter_out = "fullRoute"

    destination_interconnect = "Osaka-1"
    destination_qos_type = "guarantee"
    destination_service_key = "service_key"

    primary_connected_network_address = "10.10.0.0/30"
    secondary_connected_network_address = "10.20.0.0/30"

    bandwidth = "40M"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the connection.

* `source_router_id` - (Required) Source router ID of the connection.

* `source_group_name` - (Required) Source group name of the connection.
  Allowed values are: "group_1", "group_2", "group_3", "group_4",
  "group_5", "group_6", "group_7" and "group_8"

* `source_route_filter_in` - (Required) Ingress value of BGP Filter. 
  Allowed values are: "fullRoute", "noRoute"

* `source_route_filter_out` - (Required) Egress value of BGP Filter. 
  Allowed values are: "fullRoute", "fullRouteWithDefaultRoute",
  "defaultRoute", "privateRoute", "noRoute"

* `destination_interconnect` - (Required) Target cloud name of the connection.

* `destination_qos_type` - (Required) QOS Type of the connection.
  Currently only "guarantee" is supported.

* `destination_service_key` - (Required) Service key of the target cloud.

* `primary_connected_network_address` - (Required) Primary network address of the connection.

* `secondary_connected_network_address` - (Required) Secondary network address of the connection.

* `bandwidth` - (Required) Bandwidth of the connection. Allowed values are:
  "10M", "20M", "30M", "40M", "50M",
  "100M", "200M", "300M", "400M", "500M",
  "1G", "2G", "3G", "4G", "5G",
  "10G"

## Attributes Reference

The following attributes are exported:

* `redundant` - Redundancy of the connection.

* `tenant_id` - Tenant ID of the connection.

* `area` - Area name of the connection.
