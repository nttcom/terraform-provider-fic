---
layout: "fic"
page_title: "Flexible InterConnect: fic_eri_port_to_azure_microsoft_connection_v1"
sidebar_current: "docs-fic-resource-eri-port-to-azure-microsoft-connection-v1"
description: |-
  Manages a V1 Port to Azure Microsoft Connection resource within Flexible InterConnect.
---

# fic\_eri\_port\_to\_azure\_microsoft\_connection\_v1

Manages a V1 Port to Azure Microsoft Connection resource within Flexible InterConnect.

## Example Usage

### Basic Usage

```hcl
resource "fic_eri_port_v1" "port_1" {
  name        = "terraform_port_1"
  switch_name = "switch_name"
  port_type   = "1G"

  vlan_ranges {
    start = 497
    end   = 512
  }
}

resource "fic_eri_port_to_azure_microsoft_connection_v1" "connection_1" {
  name = "terraform_connection_1"

  source_primary_port_id   = fic_eri_port_v1.port_1.id
  source_primary_vlan      = 497
  source_secondary_port_id = fic_eri_port_v1.port_1.id
  source_secondary_vlan    = 498
  source_asn               = "65530"

  destination_interconnect = "Osaka-1"
  destination_qos_type     = "guarantee"
  destination_service_key  = "service_key"
  destination_shared_key   = "shared_key"
  destination_advertised_public_prefixes = [
    "100.100.1.1/32",
    "100.100.1.2/32",
    "100.100.1.3/32"
  ]
  destination_routing_registry_name = "ARIN"

  primary_connected_network_address   = "10.10.0.0/30"
  secondary_connected_network_address = "10.20.0.0/30"

  bandwidth = "40M"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the connection.

* `source_primary_port_id` - (Required) Primary source port's ID of the connection.

* `source_primary_vlan` - (Required) Primary source VLAN ID of the connection.

* `source_secondary_port_id` - (Required) Secondary source port's ID of the connection.

* `source_secondary_vlan` - (Required) Secondary source VLAN ID of the connection.

* `source_asn` - (Required) AS Number of the connection.

* `destination_interconnect` - (Required) Target cloud name of the connection.

* `destination_qos_type` - (Required) QOS Type of the connection.
  Currently only "guarantee" is supported.

* `destination_service_key` - (Required) Service key of the target cloud.

* `destination_shared_key` - (Optional) BGP MD5 key.

* `destination_advertised_public_prefixes` - (Required) Advertised Public Prefixes.

* `destination_routing_registry_name` - (Optional) Routing Registry Name. Allowed values are:
  "ARIN", "APNIC", "AFRINIC", "LACNIC", "RIPE",
  "NCC", "RADB", "ALTDB"

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
