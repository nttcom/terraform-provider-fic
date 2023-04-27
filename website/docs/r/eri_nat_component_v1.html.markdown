---
layout: "fic"
page_title: "Flexible InterConnect: fic_eri_nat_component_v1"
sidebar_current: "docs-fic-resource-eri-nat_component-v1"
description: |-
  Manages a V1 NAT Component resource within Flexible InterConnect.
---

# fic\_eri\_nat\_component\_v1

Manages a V1 NAT Component resource within Flexible InterConnect.

## Example Usage

### Activate NAT Component under router

```hcl
resource "fic_eri_router_v1" "router_1" {
  name            = "terraform_router_1"
  area            = "JPEAST"
  user_ip_address = "10.0.0.0/27"
  redundant       = true
}

resource "fic_eri_nat_component_v1" "nat_1" {
  router_id = fic_eri_router_v1.router_1.id
  nat_id    = fic_eri_router_v1.router_1.nat_id

  user_ip_addresses = [
    "192.168.0.0/30",
    "192.168.4.0/30",
    "192.168.8.0/30",
    "192.168.12.0/30",
    "192.168.16.0/30",
    "192.168.20.0/30",
    "192.168.24.0/30",
    "192.168.28.0/30"
  ]

  global_ip_address_sets {
    name                = "src-set-01"
    type                = "sourceNapt"
    number_of_addresses = 5
  }

  global_ip_address_sets {
    name                = "dst-set-01"
    type                = "destinationNat"
    number_of_addresses = 1
  }
}
```

## Argument Reference

The following arguments are supported:

* `router_id` - (Required) The router ID this NAT component belongs to.

* `nat_id` - (Required) ID of this NAT component.
  You can get this parameter only from parent router response body.
  So you have to add this to define NAT component in Terraform configurations.

* `user_ip_addresses` - (Required) List of user IP address.

* `global_ip_address_sets` - (Required) Global IP Address Set
  definition in activating NAT component.

* `source_napt_rules` - (Required) Source NAPT Rules
  of the NAT component.

* `destination_nat_rules` - (Required) Destination NAT Rules 
  of the NAT component.

The `source_napt_rules` block supports:

* `from` - (Required) List of source group names.
* `to` - (Required) Destination group name.
* `entries` - (Required) Conversion rules of the NAPT.

The `destination_nat_rules` block supports:

* `from` - (Required) Source group name.
* `to` - (Required) Destination group name.
* `entries` - (Required) Conversion rules of the NAT.

## Attributes Reference

The following attributes are exported:

* `source_napt_rules` - See Argument Reference above.
* `destination_nat_rules` - See Argument Reference above.
* `redundant` - Redundancy of the NAT component.
* `is_activated` - Activation status of the NAT component.

