package fic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	connections "github.com/nttcom/go-fic/fic/eri/v1/port_to_azure_microsoft_connections"
)

func TestAccEriPortToAzureMicrosoftConnectionV1Basic(t *testing.T) {
	var c connections.Connection

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckArea(t)
			testAccPreCheckAzureConnection(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEriPortToAzureMicrosoftConnectionV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccConfigEriPortToAzureMicrosoftConnectionV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEriPortToAzureMicrosoftConnectionV1Exists("fic_eri_port_to_azure_microsoft_connection_v1.connection_1", &c),
				),
			},
			resource.TestStep{
				Config: testAccConfigEriPortToAzureMicrosoftConnectionV1Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEriPortToAzureMicrosoftConnectionV1Exists("fic_eri_port_to_azure_microsoft_connection_v1.connection_1", &c),
				),
			},
		},
	})
}

func testAccCheckEriPortToAzureMicrosoftConnectionV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.eriV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "fic_eri_port_to_azure_microsoft_connection_v1" {
			continue
		}

		_, err := connections.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Connection still exists")
		}
	}

	return nil
}

func testAccCheckEriPortToAzureMicrosoftConnectionV1Exists(n string, c *connections.Connection) resource.TestCheckFunc {
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

var testAccConfigEriPortToAzureMicrosoftConnectionV1Basic = fmt.Sprintf(`
resource "fic_eri_port_v1" "port_1" {
    name = "terraform_port_1"
    switch_name = "%s"
    port_type = "1G"
    number_of_vlans = 16
}

resource "fic_eri_port_to_azure_microsoft_connection_v1" "connection_1" {
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
    destination_advertised_public_prefixes = [
        "100.100.1.1/32",
        "100.100.1.2/32",
        "100.100.1.3/32"
    ]
    destination_routing_registry_name = "ARIN"

    primary_connected_network_address = "10.10.0.0/30"
    secondary_connected_network_address = "10.20.0.0/30"

    bandwidth = "40M"
}
`,
	OS_SWITCH_NAME,
	OS_AZURE_SERVICE_KEY,
	OS_AZURE_SHARED_KEY,
)

var testAccConfigEriPortToAzureMicrosoftConnectionV1Update = fmt.Sprintf(`
resource "fic_eri_port_v1" "port_1" {
    name = "terraform_port_1"
    switch_name = "%s"
    port_type = "1G"
    number_of_vlans = 16
}

resource "fic_eri_port_to_azure_microsoft_connection_v1" "connection_1" {
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
    destination_advertised_public_prefixes = [
        "100.100.1.4/32",
        "100.100.1.5/32",
        "100.100.1.6/32"
    ]
    destination_routing_registry_name = "APNIC"

    primary_connected_network_address = "10.10.0.0/30"
    secondary_connected_network_address = "10.20.0.0/30"

    bandwidth = "40M"
}
`,
	OS_SWITCH_NAME,
	OS_AZURE_SERVICE_KEY,
	OS_AZURE_SHARED_KEY,
)
