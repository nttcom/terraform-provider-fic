---
layout: "fic"
page_title: "Flexible InterConnect: fic_eri_firewall_component_v1"
sidebar_current: "docs-fic-resource-eri-firewall_component-v1"
description: |-
  Manages a V1 Firewall Component resource within Flexible InterConnect.
---

# fic\_eri\_firewall\_component\_v1

Manages a V1 Firewall Component resource within Flexible InterConnect.

## Example Usage

### Activate Firewall Component under router

```hcl
resource "fic_eri_router_v1" "router_1" {
  name            = "terraform_router_1"
  area            = "JPEAST"
  user_ip_address = "10.0.0.0/27"
  redundant       = true
}

resource "fic_eri_firewall_component_v1" "firewall_1" {
  router_id   = fic_eri_router_v1.router_1.id
  firewall_id = fic_eri_router_v1.router_1.firewall_id

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

  rules {
    from = "group_1"
    to   = "group_2"
    entries {
      name = "rule-01"
      match_source_address_sets = [
        "group1_addset_1"
      ]
      match_destination_address_sets = [
        "group2_addset_1"
      ]
      match_application = "app_set_1"
      action            = "permit"
    }
  }

  custom_applications {
    name             = "google-drive-web"
    protocol         = "tcp"
    destination_port = "443"
  }

  application_sets {
    name = "app_set_1"
    applications = [
      "google-drive-web",
      "pre-defined-ftp"
    ]
  }

  routing_group_settings {
    group_name = "group_1"
    address_sets {
      name = "group1_addset_1"
      addresses = [
        "172.18.1.0/24"
      ]
    }
  }

  routing_group_settings {
    group_name = "group_2"
    address_sets {
      name = "group2_addset_1"
      addresses = [
        "192.168.1.0/24"
      ]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `router_id` - (Required) The router ID this Firewall Component belongs to.

* `firewall_id` - (Required) ID of this Firewall Component.
  You can get this parameter only from parent router response body.
  So you have to add this to define Firewall Component in Terraform configurations.

* `user_ip_addresses` - (Required) List of user IP address.

* `rules` - (Optional) List of firewall rules.
* `custom_applications` - (Optional) List of Firewall custom applications.
* `application_sets` - (Optional) List of Firewall application sets.
* `routing_group_settings` - (Optional) List of Firewall routing group settings.

The `rules` block supports:

* `from` - (Required) Name of the group as "from" parameter of this rule.
* `to` - (Required) Name of the group as "to" parameter of this rule.
* `entries` - (Required) List of details of this rule.

The `custom_applications` block supports:

* `name` - (Required) Custom application name
* `protocol` - (Required) Protocol of the custom application. Either "tcp" or "udp" is allowed.
* `destination_port` - (Required) Destination port of the custom application.

The `application_sets` block supports:

* `name` - (Required) Name of the application set.
* `applications` - (Required) List of applications.

The `routing_group_settings` block supports:

* `group_name` - (Required) Name of the routing group set.
* `address_sets` - (Required) List of routing group setting details.


## Attributes Reference

The following attributes are exported:

* `redundant` - Redundancy of the Firewall Component.
* `is_activated` - Activation status of the Firewall Component.

