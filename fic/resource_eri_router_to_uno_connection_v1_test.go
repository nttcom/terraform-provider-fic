package fic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	connections "github.com/nttcom/go-fic/fic/eri/v1/router_to_uno_connections"
)

func TestAccEriRouterToUNOConnectionV1Basic(t *testing.T) {
	var c connections.Connection

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckArea(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEriRouterToUNOConnectionV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccConfigEriRouterToUNOConnectionV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEriRouterToUNOConnectionV1Exists("fic_eri_router_to_uno_connection_v1.connection_1", &c),
					resource.TestCheckResourceAttr(
						"fic_eri_router_to_uno_connection_v1.connection_1", "name", "terraform_connection_1"),
				),
			},
			resource.TestStep{
				Config: testAccConfigEriRouterToUNOConnectionV1Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEriRouterToUNOConnectionV1Exists("fic_eri_router_to_uno_connection_v1.connection_1", &c),
					resource.TestCheckResourceAttr(
						"fic_eri_router_to_uno_connection_v1.connection_1", "name", "terraform_connection_1"),
				),
			},
		},
	})
}

func testAccCheckEriRouterToUNOConnectionV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.eriV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "fic_eri_router_to_uno_connection_v1" {
			continue
		}

		_, err := connections.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Connection still exists")
		}
	}

	return nil
}
func testAccCheckEriRouterToUNOConnectionV1Exists(n string, c *connections.Connection) resource.TestCheckFunc {
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

var testAccConfigEriRouterToUNOConnectionV1Basic = fmt.Sprintf(`
resource "fic_eri_router_v1" "router_1" {
	name = "terraform_router_1"
	area = "%s"
	user_ip_address = "10.0.0.0/27"
	redundant = false
}

resource "fic_eri_router_to_uno_connection_v1" "connection_1" {
	name = "terraform_connection_1"
	source_router_id = "${fic_eri_router_v1.router_1.id}"
	source_group_name = "group_1"

	source_route_filter_in = "noRoute"
	source_route_filter_out = "fullRouteWithDefaultRoute"

	destination_interconnect = "Interconnect-Osaka-1"
	destination_c_number = "C0250124868"
	destination_parent_contract_number = "N190005036"
	destination_vpn_number = "V19000708"
	destination_qos_type = "guarantee"
	destination_route_filter_out = "fullRoute"

	connected_network_address = "192.168.0.0/29"

	bandwidth = "10M"
}
`,
	OS_AREA_NAME,
)

var testAccConfigEriRouterToUNOConnectionV1Update = fmt.Sprintf(`
resource "fic_eri_router_v1" "router_1" {
	name = "terraform_router_1"
	area = "%s"
	user_ip_address = "10.0.0.0/27"
	redundant = false
}

resource "fic_eri_router_to_uno_connection_v1" "connection_1" {
	name = "terraform_connection_1"
	source_router_id = "${fic_eri_router_v1.router_1.id}"
	source_group_name = "group_1"

	source_route_filter_in = "fullRoute"
	source_route_filter_out = "fullRoute"

	destination_interconnect = "Interconnect-Osaka-1"
	destination_c_number = "C0250124868"
	destination_parent_contract_number = "N190005036"
	destination_vpn_number = "V19000708"
	destination_qos_type = "guarantee"
	destination_route_filter_out = "defaultRoute"

	connected_network_address = "192.168.0.0/29"

	bandwidth = "10M"
}
`,
	OS_AREA_NAME,
)
