package fic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	connections "github.com/nttcom/go-fic/fic/eri/v1/port_to_port_connections"
	"github.com/nttcom/go-fic/fic/eri/v1/ports"
)

func TestAccEriPortToPortConnectionV1Basic(t *testing.T) {
	var p1, p2 ports.Port
	var c connections.Connection

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckSwitchName(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEriPortToPortConnectionV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccConfigEriPortToPortConnectionV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEriPortToPortConnectionV1Exists("fic_eri_port_to_port_connection_v1.connection_1", &c),
					testAccCheckEriPortV1Exists("fic_eri_port_v1.port_1", &p1),
					testAccCheckEriPortV1Exists("fic_eri_port_v1.port_2", &p2),
					resource.TestCheckResourceAttr(
						"fic_eri_port_to_port_connection_v1.connection_1", "name", "terraform_connection_1"),
					resource.TestCheckResourceAttrPtr(
						"fic_eri_port_to_port_connection_v1.connection_1", "source_port_id", &p1.ID),
					resource.TestCheckResourceAttrPtr(
						"fic_eri_port_to_port_connection_v1.connection_1", "destination_port_id", &p2.ID),
				),
			},
		},
	})
}

func testAccCheckEriPortToPortConnectionV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.eriV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "fic_eri_port_to_port_connection_v1" {
			continue
		}

		_, err := connections.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Connection still exists")
		}
	}

	return nil
}

func testAccCheckEriPortToPortConnectionV1Exists(n string, c *connections.Connection) resource.TestCheckFunc {
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

var testAccConfigEriPortToPortConnectionV1Basic = fmt.Sprintf(`
resource "fic_eri_port_v1" "port_1" {
	name = "terraform_port_1"
	switch_name = "%s"
	port_type = "1G"
	is_activated = true

	vlan_ranges {
		start = 1137
		end = 1152
	}
}

resource "fic_eri_port_v1" "port_2" {
	name = "terraform_port_2"
	switch_name = "%s"
	port_type = "1G"
	is_activated = true
	depends_on = ["fic_eri_port_v1.port_1"]

	vlan_ranges {
		start = 1153
		end = 1168
	}
}

resource "fic_eri_port_to_port_connection_v1" "connection_1" {
	name = "terraform_connection_1"
	source_port_id = "${fic_eri_port_v1.port_1.id}"
	source_vlan = "${fic_eri_port_v1.port_1.vlan_ranges.0.start}"
	destination_port_id = "${fic_eri_port_v1.port_2.id}"
	destination_vlan = "${fic_eri_port_v1.port_2.vlan_ranges.0.start}"
	bandwidth = "10M"
}
`,
	OS_SWITCH_NAME,
	OS_SWITCH_NAME,
)
