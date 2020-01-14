---
layout: "fic"
page_title: "Flexible InterConnect: fic_eri_switch_v1"
sidebar_current: "docs-fic-datasource-eri-switch-v1"
description: |-
  Get a V1 Switch information within Flexible InterConnect.
---

# fic\_eri\_switch\_v1

Use this data source to get the ID, the VLAN ranges and Details within Flexible InterConnect.

## Example Usage

### Basic Usage

```hcl
data "fic_eri_switch_v1" "switch_1" {
	name = "lxea03comnw1"
	area = "JPEAST"
    location = "NTTComTokyo(NW1)"
    port_type = "1G"
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Optional) Alias name of switch.

* `area` - (Optional) Area name.

* `location` - (Optional) Location(Data center) name.

* `port_type` - (Required) Port type, 1G or 10G.


## Attributes Reference

The following attributes are exported:

* `name` - See Argument Reference above.
* `area` - See Argument Reference above.
* `location` - See Argument Reference above.
* `port_type` - See Argument Reference above.
* `id` - ID of switch.
* `number_of_available_vlans` - Number of available VLANs.
* `vlan_ranges` - List of available VLAN ranges.
* `vlan_ranges/start` - Start number of VLAN range.
* `vlan_ranges/end` - End number of VLAN range.
