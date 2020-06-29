package fic

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/nttcom/go-fic"

	connections "github.com/nttcom/go-fic/fic/eri/v1/router_paired_to_gcp_connections"

	"github.com/hashicorp/terraform/terraform"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccEriRouterPairedToGCPConnectionV1_basic(t *testing.T) {
	var connection connections.Connection
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "fic_eri_router_paired_to_port_connection_v1.test"
	primaryParingKey := ""
	secondaryParingKey := ""

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRouterPairedToGCPConnectionV1Destroy,
		IDRefreshName: resourceName,
		Steps: []resource.TestStep{
			{
				Config: testAccEriRouterPairedToGCPConnectionConfig(rName, "10M", "noRoute", "privateRoute", 10, primaryParingKey, secondaryParingKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouterPairedToGCPConnectionV1Exists(resourceName, &connection),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "name", rName),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "bandwidth", "10M"),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "source.0.router_id", connection.Source.RouterID),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "source.0.group_name", "group_1"),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "source.0.route_filter.0.in", "noRoute"),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "source.0.route_filter.0.out", "privateRoute"),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "source.0.primary_med_out", "10"),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "source.0.secondary_med_out", "20"),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "destination.0.primary.0.interconnect", "Equinix-TY2-2"),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "destination.0.primary.0.paring_key", ""),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "destination.0.secondary.0.interconnect", "Equinix-TY2-2"),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "destination.0.secondary.0.paring_key", ""),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "destination.0.qos_type", "guarantee"),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "redundant", "true"),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "tenant_id", connection.TenantID),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "area", "JPEAST"),
					resource.TestCheckResourceAttrSet("fic_eri_router_paired_to_port_connection_v1.test", "operation_id"),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "opetation_status", "Completed"),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "primary_connected_network_address", connection.PrimaryConnectedNetworkAddress),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "secondary_connected_network_address", connection.SecondaryConnectedNetworkAddress),
				),
			},
			{
				Config: testAccEriRouterPairedToGCPConnectionConfig(rName, "50M", "fullRoute", "defaultRoute", 30, primaryParingKey, secondaryParingKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouterPairedToGCPConnectionV1Exists(resourceName, &connection),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "name", rName),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "bandwidth", "50M"),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "source.0.router_id", connection.Source.RouterID),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "source.0.group_name", "group_1"),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "source.0.route_filter.0.in", "fullRoute"),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "source.0.route_filter.0.out", "defaultRoute"),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "source.0.primary_med_out", "30"),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "source.0.secondary_med_out", "40"),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "destination.0.primary.0.interconnect", "Equinix-TY2-2"),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "destination.0.primary.0.paring_key", ""),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "destination.0.secondary.0.interconnect", "Equinix-TY2-2"),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "destination.0.secondary.0.paring_key", ""),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "destination.0.qos_type", "guarantee"),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "redundant", "true"),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "tenant_id", connection.TenantID),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "area", "JPEAST"),
					resource.TestCheckResourceAttrSet("fic_eri_router_paired_to_port_connection_v1.test", "operation_id"),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "opetation_status", "Completed"),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "primary_connected_network_address", connection.PrimaryConnectedNetworkAddress),
					resource.TestCheckResourceAttr("fic_eri_router_paired_to_port_connection_v1.test", "secondary_connected_network_address", connection.SecondaryConnectedNetworkAddress),
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

func testAccCheckRouterPairedToGCPConnectionV1Exists(resourceName string, connection *connections.Connection) resource.TestCheckFunc {
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

func testAccCheckRouterPairedToGCPConnectionV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.eriV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating FIC client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "fic_eri_router_paired_to_port_connection_v1" {
			continue
		}

		if result := connections.Get(client, rs.Primary.ID); result.Err != nil {
			var e *fic.ErrDefault404
			if errors.As(result.Err, &e) {
				return nil
			}
			return fmt.Errorf("error getting FIC paired router to GCP connection: %w", err)
		}

		return fmt.Errorf("connection (%s) still exists", rs.Primary.ID)
	}

	return nil
}

func testAccEriRouterPairedToGCPConnectionConfig(rName, bandwidth, routeFilterIn, routeFilterOut string, primaryMEDOut int, primaryParingKey, secondaryParingKey string) string {
	return fmt.Sprintf(`
resource "fic_eri_router_v1" "test" {
	name = %[1]q
	area = "JPEAST"
	user_ip_address = "10.0.0.0/27"
	redundant = true
}

resource "fic_eri_router_paired_to_port_connection_v1" "test" {
	name = %[1]q
	bandwidth = %[2]q
	source {
		router_id = fic_eri_router_v1.test.id
		group_name = "group_1"
		route_filter {
			in = %[3]q
			out = %[4]q
		}
		primary_med_out = %d
	}
	destination {
		primary {
			interconnect = "Equinix-TY2-2"
			paring_key = %[5]q
		}
		secondary {
			interconnect = "Equinix-TY2-2"
			paring_key = %[6]q
		}
	}
}
`, rName, bandwidth, routeFilterIn, routeFilterOut, primaryMEDOut, primaryParingKey, secondaryParingKey)
}
