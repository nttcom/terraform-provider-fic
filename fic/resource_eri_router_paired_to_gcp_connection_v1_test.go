package fic

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/nttcom/go-fic"

	connections "github.com/nttcom/go-fic/fic/eri/v1/router_paired_to_gcp_connections"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccPairedRouterToGCPConnection_basic(t *testing.T) {
	var connection connections.Connection
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "fic_eri_router_paired_to_gcp_connection_v1.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheckGCPConnection(t) },
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckPairedRouterToGCPConnectionDestroy,
		IDRefreshName: resourceName,
		Steps: []resource.TestStep{
			{
				Config: testAccPairedRouterToGCPConnectionConfig(rName, "10M", "noRoute", "privateRoute", 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPairedRouterToGCPConnectionExists(resourceName, &connection),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "bandwidth", "10M"),
					resource.TestCheckResourceAttrSet(resourceName, "source.0.router_id"),
					resource.TestCheckResourceAttr(resourceName, "source.0.group_name", "group_1"),
					resource.TestCheckResourceAttr(resourceName, "source.0.route_filter.0.in", "noRoute"),
					resource.TestCheckResourceAttr(resourceName, "source.0.route_filter.0.out", "privateRoute"),
					resource.TestCheckResourceAttr(resourceName, "source.0.primary_med_out", "10"),
					resource.TestCheckResourceAttr(resourceName, "source.0.secondary_med_out", "20"),
					resource.TestCheckResourceAttr(resourceName, "destination.0.primary.0.interconnect", "Equinix-TY2-2"),
					resource.TestCheckResourceAttrSet(resourceName, "destination.0.primary.0.pairing_key"),
					resource.TestCheckResourceAttr(resourceName, "destination.0.secondary.0.interconnect", "@Tokyo-CC2-2"),
					resource.TestCheckResourceAttrSet(resourceName, "destination.0.secondary.0.pairing_key"),
					resource.TestCheckResourceAttr(resourceName, "destination.0.qos_type", "guarantee"),
					resource.TestCheckResourceAttr(resourceName, "redundant", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "tenant_id"),
					resource.TestCheckResourceAttr(resourceName, "area", "JPEAST"),
					resource.TestCheckResourceAttrSet(resourceName, "operation_id"),
					resource.TestCheckResourceAttr(resourceName, "operation_status", "Completed"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_connected_network_address"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_connected_network_address"),
				),
			},
			{
				Config: testAccPairedRouterToGCPConnectionConfig(rName, "50M", "fullRoute", "defaultRoute", 30),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPairedRouterToGCPConnectionExists(resourceName, &connection),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "bandwidth", "50M"),
					resource.TestCheckResourceAttrSet(resourceName, "source.0.router_id"),
					resource.TestCheckResourceAttr(resourceName, "source.0.group_name", "group_1"),
					resource.TestCheckResourceAttr(resourceName, "source.0.route_filter.0.in", "fullRoute"),
					resource.TestCheckResourceAttr(resourceName, "source.0.route_filter.0.out", "defaultRoute"),
					resource.TestCheckResourceAttr(resourceName, "source.0.primary_med_out", "30"),
					resource.TestCheckResourceAttr(resourceName, "source.0.secondary_med_out", "40"),
					resource.TestCheckResourceAttr(resourceName, "destination.0.primary.0.interconnect", "Equinix-TY2-2"),
					resource.TestCheckResourceAttrSet(resourceName, "destination.0.primary.0.pairing_key"),
					resource.TestCheckResourceAttr(resourceName, "destination.0.secondary.0.interconnect", "@Tokyo-CC2-2"),
					resource.TestCheckResourceAttrSet(resourceName, "destination.0.secondary.0.pairing_key"),
					resource.TestCheckResourceAttr(resourceName, "destination.0.qos_type", "guarantee"),
					resource.TestCheckResourceAttr(resourceName, "redundant", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "tenant_id"),
					resource.TestCheckResourceAttr(resourceName, "area", "JPEAST"),
					resource.TestCheckResourceAttrSet(resourceName, "operation_id"),
					resource.TestCheckResourceAttr(resourceName, "operation_status", "Completed"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_connected_network_address"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_connected_network_address"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"operation_id",
				},
			},
		},
	})
}

func testAccCheckPairedRouterToGCPConnectionExists(resourceName string, connection *connections.Connection) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("id is not set")
		}

		config := testAccProvider.Meta().(*Config)
		client, err := config.eriV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating FIC client: %w", err)
		}

		actual, err := connections.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return fmt.Errorf("error getting FIC paired router to GCP connection: %w", err)
		}

		*connection = *actual

		return nil
	}
}

func testAccCheckPairedRouterToGCPConnectionDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.eriV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating FIC client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "fic_eri_router_paired_to_gcp_connection_v1" {
			continue
		}

		if result := connections.Get(client, rs.Primary.ID); result.Err != nil {
			var e fic.ErrDefault404
			if errors.As(result.Err, &e) {
				return nil
			}
			return fmt.Errorf("error getting FIC paired router to GCP connection: %w", err)
		}

		return fmt.Errorf("connection (%s) still exists", rs.Primary.ID)
	}

	return nil
}

func testAccPairedRouterToGCPConnectionConfig(rName, bandwidth, routeFilterIn, routeFilterOut string, primaryMEDOut int) string {
	return fmt.Sprintf(`
resource "google_compute_router" "test1" {
	name = "%[1]s1"
	network = "default"
	bgp {
		asn = 16550
	}
}

resource "google_compute_router" "test2" {
	name = "%[1]s2"
	network = "default"
	bgp {
		asn = 16550
	}
}

resource "google_compute_interconnect_attachment" "test1" {
	name = "%[1]s1"
	router = google_compute_router.test1.id
	type = "PARTNER"
	edge_availability_domain = "AVAILABILITY_DOMAIN_1"
}

resource "google_compute_interconnect_attachment" "test2" {
	name = "%[1]s2"
	router = google_compute_router.test2.id
	type = "PARTNER"
	edge_availability_domain = "AVAILABILITY_DOMAIN_2"
}

resource "fic_eri_router_v1" "test" {
	name = %[1]q
	area = "JPEAST"
	user_ip_address = "10.0.0.0/27"
	redundant = true
}

resource "fic_eri_router_paired_to_gcp_connection_v1" "test" {
	name = %[1]q
	bandwidth = %[2]q
	source {
		router_id = fic_eri_router_v1.test.id
		group_name = "group_1"
		route_filter {
			in = %[3]q
			out = %[4]q
		}
		primary_med_out = %[5]d
	}
	destination {
		primary {
			interconnect = "Equinix-TY2-2"
			pairing_key = google_compute_interconnect_attachment.test1.pairing_key
		}
		secondary {
			interconnect = "@Tokyo-CC2-2"
			pairing_key = google_compute_interconnect_attachment.test2.pairing_key
		}
	}
}
`, rName, bandwidth, routeFilterIn, routeFilterOut, primaryMEDOut)
}
