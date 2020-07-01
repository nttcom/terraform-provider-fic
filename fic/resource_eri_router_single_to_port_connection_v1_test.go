package fic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	connections "github.com/nttcom/go-fic/fic/eri/v1/router_single_to_port_connections"
)

func TestAccEriRouterSingleToPortConnectionV1Basic(t *testing.T) {
	// var p1, p2 ports.Port
	var c connections.Connection

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckSwitchName(t)
			testAccPreCheckArea(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEriRouterSingleToPortConnectionV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccConfigEriRouterSingleToPortConnectionV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEriRouterSingleToPortConnectionV1Exists("fic_eri_router_single_to_port_connection_v1.connection_1", &c),
					resource.TestCheckResourceAttr(
						"fic_eri_router_single_to_port_connection_v1.connection_1", "name", "terraform_connection_1"),
				),
			},
			resource.TestStep{
				Config: testAccConfigEriRouterSingleToPortConnectionV1Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEriRouterSingleToPortConnectionV1Exists("fic_eri_router_single_to_port_connection_v1.connection_1", &c),
					resource.TestCheckResourceAttr(
						"fic_eri_router_single_to_port_connection_v1.connection_1", "name", "terraform_connection_1"),
				),
			},
		},
	})
}

func testAccCheckEriRouterSingleToPortConnectionV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.eriV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "fic_eri_router_single_to_port_connection_v1" {
			continue
		}

		_, err := connections.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Connection still exists")
		}
	}

	return nil
}
func testAccCheckEriRouterSingleToPortConnectionV1Exists(n string, c *connections.Connection) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		client, err := config.eriV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FIC ERI client: %s", err)
		}

		found, err := connections.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Connection not found")
		}

		*c = *found

		return nil
	}
}

var testAccConfigEriRouterSingleToPortConnectionV1Basic = fmt.Sprintf(`
resource "fic_eri_router_v1" "router_1" {
	name = "terraform_router_1"
	area = "%s"
	user_ip_address = "10.0.0.0/27"
	redundant = false
}

resource "fic_eri_port_v1" "port_1" {
	name = "terraform_port_1"
	switch_name = "%s"
	port_type = "10G"

	vlan_ranges {
		start = 1137
		end = 1152
	}
}

resource "fic_eri_router_single_to_port_connection_v1" "connection_1" {
	name = "terraform_connection_1"
	source_router_id = "${fic_eri_router_v1.router_1.id}"
	source_group_name = "group_1"

	source_information {
		ip_address = "10.0.1.1/30"
		as_path_prepend_in = "4"
		as_path_prepend_out = "4"
	}

	source_route_filter_in = "fullRoute"
	source_route_filter_out = "fullRouteWithDefaultRoute"

	destination_information {
		port_id = "${fic_eri_port_v1.port_1.id}"
		vlan = "${fic_eri_port_v1.port_1.vlan_ranges.0.start}"
		ip_address = "10.0.1.2/30"
		asn = "65000"
	}

	bandwidth = "10M"
}
`,
	OS_AREA_NAME,
	OS_SWITCH_NAME,
)

var testAccConfigEriRouterSingleToPortConnectionV1Update = fmt.Sprintf(`
resource "fic_eri_router_v1" "router_1" {
	name = "terraform_router_1"
	area = "%s"
	user_ip_address = "10.0.0.0/27"
	redundant = false
}

resource "fic_eri_port_v1" "port_1" {
	name = "terraform_port_1"
	switch_name = "%s"
	port_type = "10G"

	vlan_ranges {
		start = 1137
		end = 1152
	}
}

resource "fic_eri_router_single_to_port_connection_v1" "connection_1" {
	name = "terraform_connection_1"
	source_router_id = "${fic_eri_router_v1.router_1.id}"
	source_group_name = "group_1"

	source_information {
		ip_address = "10.0.1.1/30"
		as_path_prepend_in = "OFF"
		as_path_prepend_out = "2"
	}

	source_route_filter_in = "noRoute"
	source_route_filter_out = "fullRoute"

	destination_information {
		port_id = "${fic_eri_port_v1.port_1.id}"
		vlan = "${fic_eri_port_v1.port_1.vlan_ranges.0.start}"
		ip_address = "10.0.1.2/30"
		asn = "65000"
	}

	bandwidth = "10M"
}
`,
	OS_AREA_NAME,
	OS_SWITCH_NAME,
)
