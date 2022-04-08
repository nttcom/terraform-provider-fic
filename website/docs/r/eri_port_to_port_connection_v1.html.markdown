---
layout: "fic"
page_title: "Flexible InterConnect: fic_eri_port_to_port_connection_v1"
sidebar_current: "docs-fic-resource-eri-port-to-port-connection-v1"
description: |-
  Manages a V1 Port to Port Connection resource within Flexible InterConnect.
---

# fic\_eri\_port\_to\_port\_connection\_v1

Manages a V1 Port to Port Connection resource within Flexible InterConnect.

## Example Usage

### Basic Usage

```hcl
resource "fic_eri_port_v1" "port_1" {
  name         = "terraform_port_1"
  switch_name  = "lxea01comnw1"
  port_type    = "1G"
  is_activated = true

  vlan_ranges {
    start = 1137
    end   = 1152
  }
}

resource "fic_eri_port_v1" "port_2" {
  name         = "terraform_port_2"
  switch_name  = "lxea01comnw1"
  port_type    = "1G"
  is_activated = true

  vlan_ranges {
    start = 1153
    end   = 1168
  }
}

resource "fic_eri_port_to_port_connection_v1" "connection_1" {
  name                = "terraform_connection_1"
  source_port_id      = fic_eri_port_v1.port_1.id
  source_vlan         = fic_eri_port_v1.port_1.vlan_ranges[0].start
  destination_port_id = fic_eri_port_v1.port_2.id
  destination_vlan    = fic_eri_port_v1.port_2.vlan_ranges[0].start
  bandwidth           = "10M"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the resource.

* `source_port_id` - (Required) Source port ID of the connection.

* `source_vlan` - (Required) Source VLAN ID of the connection.

* `destination_port_id` - (Required) Destination port ID of the connection.

* `destination_vlan` - (Required) Destination VLAN ID of the connection.

* `bandwidth` - (Optional) Bandwidth of the connection. 
  Allowed values are "10M", "20M", "30M", "40M", "50M", "100M", "200M", "300M", "400M", "500M",
					"1G", "2G", "3G", "4G", "5G" and "10G" .

## Attributes Reference

The following attributes are exported:

* `redundant` - Redundancy of the connection.
* `tenant_id` - Tenant ID of the connection.
* `area` - Area name of the connection.

