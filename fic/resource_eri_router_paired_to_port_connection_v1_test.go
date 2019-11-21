package fic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	connections "github.com/nttcom/go-fic/fic/eri/v1/router_paired_to_port_connections"
)

func TestAccEriRouterPairedToPortConnectionV1Basic(t *testing.T) {
	var c connections.Connection

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckSwitchName(t)
			testAccPreCheckArea(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEriRouterPairedToPortConnectionV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccConfigEriRouterPairedToPortConnectionV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEriRouterPairedToPortConnectionV1Exists("fic_eri_router_paired_to_port_connection_v1.connection_1", &c),
					resource.TestCheckResourceAttr(
						"fic_eri_router_paired_to_port_connection_v1.connection_1", "name", "terraform_connection_1"),
				),
			},
			resource.TestStep{
				Config: testAccConfigEriRouterPairedToPortConnectionV1Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEriRouterPairedToPortConnectionV1Exists("fic_eri_router_paired_to_port_connection_v1.connection_1", &c),
					resource.TestCheckResourceAttr(
						"fic_eri_router_paired_to_port_connection_v1.connection_1", "name", "terraform_connection_1"),
				),
			},
		},
	})
}

func testAccCheckEriRouterPairedToPortConnectionV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.eriV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "fic_eri_router_paired_to_port_connection_v1" {
			continue
		}

		_, err := connections.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Connection still exists")
		}
	}

	return nil
}

func testAccCheckEriRouterPairedToPortConnectionV1Exists(n string, c *connections.Connection) resource.TestCheckFunc {
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

var testAccConfigEriRouterPairedToPortConnectionV1Basic = fmt.Sprintf(`
resource "fic_eri_router_v1" "router_1" {
	name = "terraform_router_1"
	area = "%s"
	user_ip_address = "10.0.0.0/27"
	redundant = true
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

resource "fic_eri_port_v1" "port_2" {
	name = "terraform_port_2"
	switch_name = "%s"
	port_type = "10G"
	depends_on = ["fic_eri_port_v1.port_1"]

	vlan_ranges {
		start = 1153
		end = 1168
	}
}

resource "fic_eri_router_paired_to_port_connection_v1" "connection_1" {
	name = "terraform_connection_1"
	source_router_id = "${fic_eri_router_v1.router_1.id}"
	source_group_name = "group_1"

	source_information {
		ip_address = "10.0.1.1/30"
		as_path_prepend_in = "4"
		as_path_prepend_out = "4"
	}

	source_information {
		ip_address = "10.0.1.5/30"
		as_path_prepend_in = "2"
		as_path_prepend_out = "1"
	}

	source_route_filter_in = "fullRoute"
	source_route_filter_out = "fullRouteWithDefaultRoute"

	destination_information {
		port_id = "${fic_eri_port_v1.port_1.id}"
		vlan = "${fic_eri_port_v1.port_1.vlan_ranges.0.start}"
		ip_address = "10.0.1.2/30"
		asn = "65000"
	}

	destination_information {
		port_id = "${fic_eri_port_v1.port_2.id}"
		vlan = "${fic_eri_port_v1.port_2.vlan_ranges.0.start}"
		ip_address = "10.0.1.6/30"
		asn = "65000"
	}

	bandwidth = "10M"
}
`,
	OS_AREA_NAME,
	OS_SWITCH_NAME,
	OS_SWITCH_NAME,
)

var testAccConfigEriRouterPairedToPortConnectionV1Update = fmt.Sprintf(`
resource "fic_eri_router_v1" "router_1" {
	name = "terraform_router_1"
	area = "%s"
	user_ip_address = "10.0.0.0/27"
	redundant = true
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

resource "fic_eri_port_v1" "port_2" {
	name = "terraform_port_2"
	switch_name = "%s"
	port_type = "10G"
	depends_on = ["fic_eri_port_v1.port_1"]

	vlan_ranges {
		start = 1153
		end = 1168
	}
}

resource "fic_eri_router_paired_to_port_connection_v1" "connection_1" {
	name = "terraform_connection_1"
	source_router_id = "${fic_eri_router_v1.router_1.id}"
	source_group_name = "group_1"

	source_information {
		ip_address = "10.0.1.1/30"
		as_path_prepend_in = "OFF"
		as_path_prepend_out = "2"
	}

	source_information {
		ip_address = "10.0.1.5/30"
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

	destination_information {
		port_id = "${fic_eri_port_v1.port_2.id}"
		vlan = "${fic_eri_port_v1.port_2.vlan_ranges.0.start}"
		ip_address = "10.0.1.6/30"
		asn = "65000"
	}

	bandwidth = "10M"
}
`,
	OS_AREA_NAME,
	OS_SWITCH_NAME,
	OS_SWITCH_NAME,
)
