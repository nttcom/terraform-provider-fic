package fic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	connections "github.com/nttcom/go-fic/fic/eri/v1/router_to_ecl_connections"
)

func TestAccEriRouterToECLConnectionV1Basic(t *testing.T) {
	var c connections.Connection

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckArea(t)
			testAccPreCheckRouterToECLConnection(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEriRouterToECLConnectionV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccConfigEriRouterToECLConnectionV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEriRouterToECLConnectionV1Exists("fic_eri_router_to_ecl_connection_v1.connection_1", &c),
				),
			},
			resource.TestStep{
				Config: testAccConfigEriRouterToECLConnectionV1Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEriRouterToECLConnectionV1Exists("fic_eri_router_to_ecl_connection_v1.connection_1", &c),
				),
			},
		},
	})
}

func testAccCheckEriRouterToECLConnectionV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.eriV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "fic_eri_router_to_ecl_connection_v1" {
			continue
		}

		_, err := connections.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Connection still exists")
		}
	}

	return nil
}

func testAccCheckEriRouterToECLConnectionV1Exists(n string, c *connections.Connection) resource.TestCheckFunc {
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

var testAccConfigEriRouterToECLConnectionV1Basic = fmt.Sprintf(`
resource "fic_eri_router_v1" "router_1" {
	name = "terraform_router_1"
	area = "%s"
	user_ip_address = "10.0.0.0/27"
	redundant = false
}

resource "fic_eri_router_to_ecl_connection_v1" "connection_1" {
	name = "terraform_connection_1"
	source_router_id = "${fic_eri_router_v1.router_1.id}"

	source_group_name = "group_1"
	source_route_filter_in = "noRoute"
	source_route_filter_out = "noRoute"

	destination_interconnect = "ECL-OSA-GU-Z1-01"
	destination_qos_type = "guarantee"
	destination_ecl_tenant_id = "%s"
	destination_ecl_api_key = "%s"
	destination_ecl_api_secret_key = "%s"

	bandwidth = "10M"

	primary_connected_network_address = "192.168.0.0/30"
	secondary_connected_network_address = "192.168.1.0/30"
}
`,
	OS_AREA_NAME,
	OS_ECL_TENANT_ID,
	OS_ECL_API_KEY,
	OS_ECL_API_SECRET_KEY,
)

var testAccConfigEriRouterToECLConnectionV1Update = fmt.Sprintf(`
resource "fic_eri_router_v1" "router_1" {
	name = "terraform_router_1"
	area = "%s"
	user_ip_address = "10.0.0.0/27"
	redundant = false
}

resource "fic_eri_router_to_ecl_connection_v1" "connection_1" {
	name = "terraform_connection_1"
	source_router_id = "${fic_eri_router_v1.router_1.id}"

	source_group_name = "group_1"
	source_route_filter_in = "fullRoute"
	source_route_filter_out = "fullRouteWithDefaultRoute"

	destination_interconnect = "ECL-OSA-GU-Z1-01"
	destination_qos_type = "guarantee"
	destination_ecl_tenant_id = "%s"
	destination_ecl_api_key = "%s"
	destination_ecl_api_secret_key = "%s"

	bandwidth = "10M"

	primary_connected_network_address = "192.168.0.0/30"
	secondary_connected_network_address = "192.168.1.0/30"
}
`,
	OS_AREA_NAME,
	OS_ECL_TENANT_ID,
	OS_ECL_API_KEY,
	OS_ECL_API_SECRET_KEY,
)
