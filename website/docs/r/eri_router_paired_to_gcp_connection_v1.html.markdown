---
layout: "fic"
page_title: "Flexible InterConnect: fic_eri_router_paired_to_gcp_connection_v1"
sidebar_current: "docs-fic-resource-eri-router-paired-to-gcp-connection-v1"
description: |-
  Manages a V1 Router(Paired) to GCP Connection resource within Flexible InterConnect.
---

# fic\_eri\_router\_paired\_to\_gcp\_connection\_v1

Manages a V1 Router(Paired) to GCP Connection resource within Flexible InterConnect.

## Example Usage

### Basic Usage

```hcl
resource "google_compute_router" "router1" {
	name = "tf-router1"
	network = "default"
	bgp {
		asn = 16550
	}
}

resource "google_compute_router" "router2" {
	name = "tf-router2"
	network = "default"
	bgp {
		asn = 16550
	}
}

resource "google_compute_interconnect_attachment" "interconnect1" {
	name = "tf-interconnect1"
	router = google_compute_router.router1.id
	type = "PARTNER"
	edge_availability_domain = "AVAILABILITY_DOMAIN_1"
}

resource "google_compute_interconnect_attachment" "interconnect2" {
	name = "tf-interconnect2"
	router = google_compute_router.router2.id
	type = "PARTNER"
	edge_availability_domain = "AVAILABILITY_DOMAIN_2"
}

resource "fic_eri_router_v1" "router" {
	name = "tf-router"
	area = "JPEAST"
	user_ip_address = "10.0.0.0/27"
	redundant = true
}

resource "fic_eri_router_paired_to_gcp_connection_v1" "connection" {
	name = "tf-connection"
	bandwidth = "10M"
	source {
		router_id = fic_eri_router_v1.router.id
		group_name = "group_1"
		route_filter {
			in = "noRoute"
			out = "privateRoute"
		}
		primary_med_out = 10
	}
	destination {
		primary {
			interconnect = "Equinix-TY2-2"
			pairing_key = google_compute_interconnect_attachment.interconnect1.pairing_key
		}
		secondary {
			interconnect = "@Tokyo-CC2-2"
			pairing_key = google_compute_interconnect_attachment.interconnect2.pairing_key
		}
	}
}
```

## Argument Reference

The following arguments supported:

* `name` - (Required) Name of the connection.
  It must be less than 64 characters in half-width alphanumeric characters and some symbols &()-_.

* `bandwidth` - (Required) Bandwidth of the connection.
  Either "10M", "50M", "100M", "200M", "300M", "400M", "500M", "1G", "2G", "5G" or "10G".

* `source` - (Required) Source of the connection. Structure is documented below.

* `destination` - (Required) Destination of the connection. Structure is documented below.

The `source` block supports:

* `router_id` - (Required) Router ID. It must be a F + 12-digit number.

* `group_name` - (Required) Group name.
  Either "group_1", "group_2", "group_3", "group_4", "group_5", "group_6", "group_7" or "group_8".

* `route_fileter` - (Required) Route filter. Structure is documented below.

* `primary_med_out` - (Required) MED egress value of primary. Either 10 or 30.

The `route_filter` block supports:

* `in` - (Required) BGP filter ingress value. Either "fullRoute" or "noRoute".

* `out` - (Required) BGP filter egress value.
  Either "fullRoute", "fullRouteWithDefaultRoute", "defaultRoute", "privateRoute" or "noRoute".

The `destination` block supports:

* `primary` - (Required) Primary interconnect of destination.

* `secondary` - (Required) Secondary interconnect of destination.

The `primary` and `secondary` blocks support:

* `interconnect` - (Required) Connecting point.
  See "1.3. FIC-Connection Google Cloud" in [FIC-Connection Google Cloud](https://sdpf.ntt.com/services/docs/fic/service-descriptions/connection-gcp/connection-gcp.html#id4).

* `pairing_key` - (Required) Paring key of google hybrid interconnect.

## Attributes Reference

The following attributes are exported:

* `source.0.secondary_med_out` - MED egress value of secondary. It would be source.primary_med_out plus 10.
* `destination.0.qos_type` - QoS type. It would be "guarantee".
* `redundant` - Redundant flag of the connection. It would be true.
* `tenant_id` - Tenant ID where the connection belongs.
* `area` - Area name of the connection.
* `operation_id` - ID of the last operation.
* `operation_status` - Status of the last operation.
* `primary_connected_network_address` - Primary connected network address. It would be "<network_address>/29".
* `secondary_connected_network_address` - Secondary connected network address. It would be "<network_address>/29".

## Timeouts

This resource provides the following Timeout configuration options:

- `create` - Default is 10 minutes.
- `update` - Default is 10 minutes.
- `delete` - Default is 10 minutes.

## Import

Connections can be imported using the ID:

```
$ terraform import fic_eri_router_paired_to_gcp_connection_v1.connection F030123456789
```
