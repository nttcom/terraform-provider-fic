---
layout: "fic"
page_title: "Flexible InterConnect: fic_eri_nat_global_ip_address_set_v1"
sidebar_current: "docs-fic-resource-eri-nat_global_ip_address_set-v1"
description: |-
  Manages a V1 NAT Global IP Address Set resource within Flexible InterConnect.
---

# fic\_eri\_nat\_global\_ip\_address\_set\_v1

Manages a V1 NAT Global IP Address Set resource within Flexible InterConnect.

## Example Usage

### Add Global IP Address Set under activated NAT Component

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

resource "fic_eri_nat_global_ip_address_set_v1" "gip_1" {
  router_id  = fic_eri_router_v1.router_1.id
  nat_id     = fic_eri_router_v1.router_1.nat_id
  depends_on = [fic_eri_nat_component_v1.nat_1]

  name                = "src-set-02"
  type                = "sourceNapt"
  number_of_addresses = 5
}
```

## Argument Reference

The following arguments are supported:

* `router_id` - (Required) The router ID this global ip address 
  set belongs to.

* `nat_id` - (Required) The nat component ID this global ip address 
  set belongs to.

* `name` - (Required) Name of the global ip address set.

* `type` - (Required) Address type of the global ip address set.
  "sourceNapt" or "destinationNat" can be specified.

* `number_of_addresses` - (Required) Number of the IP addresses.


## Attributes Reference

The following attributes are exported:

* `addresses` - Created global IP addresses.

