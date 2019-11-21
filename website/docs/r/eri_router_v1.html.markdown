---
layout: "fic"
page_title: "Flexible InterConnect: fic_eri_router_v1"
sidebar_current: "docs-fic-resource-eri-router-v1"
description: |-
  Manages a V1 Router resource within Flexible InterConnect.
---

# fic\_eri\_router\_v1

Manages a V1 Router resource within Flexible InterConnect.

## Example Usage

### Basic Usage

```hcl
resource "fic_eri_router_v1" "router_1" {
	name = "terraform_router_1"
	area = "JPEAST"
	user_ip_address = "10.0.0.0/27"
	redundant = true
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the resource.

* `area` - (Required) Area name you create a router.

* `user_ip_address` - (Required; Required if `vlan_ranges` is empy) The IP Address of the router.
  It must have prefix 27.

* `redundant` - (Required) The redundant option of the router.


## Attributes Reference

The following attributes are exported:

* `name` - See Argument Reference above.
* `area` - See Argument Reference above.
* `user_ip_address` - See Argument Reference above.
* `redundant` - See Argument Reference above.
* `tenant_id` - Tenant ID the router belongs to.
* `firewalls/id` - Firewall ID.
* `firewalls/is_activated` - Activate status of the Firewall.
* `nats/id` - NAT component ID.
* `nats/is_activated` - Activate status of the NAT.
* `routing_groups/name` - Routing group name of the router.
