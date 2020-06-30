package fic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	connections "github.com/nttcom/go-fic/fic/eri/v1/port_to_azure_private_connections"
)

func TestAccEriPortToAzurePrivateConnectionV1Basic(t *testing.T) {
	var c connections.Connection

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckArea(t)
			testAccPreCheckAzureConnection(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEriPortToAzurePrivateConnectionV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccConfigEriPortToAzurePrivateConnectionV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEriPortToAzurePrivateConnectionV1Exists("fic_eri_port_to_azure_private_connection_v1.connection_1", &c),
				),
			},
		},
	})
}

func testAccCheckEriPortToAzurePrivateConnectionV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.eriV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "fic_eri_port_to_azure_private_connection_v1" {
			continue
		}

		_, err := connections.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Connection still exists")
		}
	}

	return nil
}

func testAccCheckEriPortToAzurePrivateConnectionV1Exists(n string, c *connections.Connection) resource.TestCheckFunc {
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

var testAccConfigEriPortToAzurePrivateConnectionV1Basic = fmt.Sprintf(`
resource "fic_eri_port_v1" "port_1" {
    name = "terraform_port_1"
    switch_name = "%s"
    port_type = "1G"

    number_of_vlans = 16
}

resource "fic_eri_port_to_azure_private_connection_v1" "connection_1" {
    name = "terraform_connection_1"

    source_primary_port_id = "${fic_eri_port_v1.port_1.id}"
    source_primary_vlan = "${fic_eri_port_v1.port_1.vlans.0.vid}"
    source_secondary_port_id = "${fic_eri_port_v1.port_1.id}"
    source_secondary_vlan = "${fic_eri_port_v1.port_1.vlans.1.vid}"
    source_asn = "65530"

    destination_interconnect = "Osaka-1"
    destination_qos_type = "guarantee"
    destination_service_key = "%s"
    destination_shared_key = "%s"

    primary_connected_network_address = "10.10.0.0/30"
    secondary_connected_network_address = "10.20.0.0/30"

    bandwidth = "40M"
}
`,
	OS_SWITCH_NAME,
	OS_AZURE_SERVICE_KEY,
	OS_AZURE_SHARED_KEY,
)
