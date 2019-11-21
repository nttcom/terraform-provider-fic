---
layout: "fic"
page_title: "Flexible InterConnect: fic_eri_port_v1"
sidebar_current: "docs-fic-resource-eri-port-v1"
description: |-
  Manages a V1 Port resource within Flexible InterConnect.
---

# fic\_eri\_port\_v1

Manages a V1 Port resource within Flexible InterConnect.

## Example Usage

### Basic Usage

```hcl
resource "fic_eri_port_v1" "port_1" {
  name = "terraform_port_1"
  switch_name = "%s"
  port_type = "1G"
  number_of_vlans = 16
}
```

### Create and Activate

```hcl
resource "fic_eri_port_v1" "port_1" {
  name = "terraform_port_1"
  switch_name = "%s"
  port_type = "1G"
  number_of_vlans = 16
  is_activated = true
}
```

### Create with VLAN Ranges

```hcl
resource "fic_eri_port_v1" "port_1" {
	name = "terraform_port_1"
	switch_name = "%s"
	port_type = "1G"

  vlan_ranges {
    start = 1137
    end = 1152
  }

  vlan_ranges {
    start = 1153
    end = 1168  
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the resource.

* `switch_name` - (Required) Switch name you create a port.

* `number_of_vlans` - (Optional; Required if `vlan_ranges` is empty) The number of VLANs used by port.

* `vlan_ranges` - (Optional; Required if `number_of_vlans` is empty) The list of VLAN ranges object.

* `port_type` - (Optional) Type of port either "1G" or "10G".

* `is_acivated` - (Optional) Activate status of the port.


## Attributes Reference

The following attributes are exported:

* `name` - See Argument Reference above.
* `switch_name` - See Argument Reference above.
* `port_type` - See Argument Reference above.
* `is_activated` - See Argument Reference above.
* `tenant_id` - Tenant ID the port belongs to.
* `area` - Area name the port belongs to.
* `location` - Location name the port belongs to.
* `vlans/vid` - VLAN ID of the router.
* `vlans/status` - VLAN status of the port.
