package fic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	connections "github.com/nttcom/go-fic/fic/eri/v1/router_to_azure_microsoft_connections"
)

func TestAccEriRouterToAzureMicrosoftConnectionV1Basic(t *testing.T) {
	var c connections.Connection

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckArea(t)
			testAccPreCheckAzureConnection(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEriRouterToAzureMicrosoftConnectionV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccConfigEriRouterToAzureMicrosoftConnectionV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEriRouterToAzureMicrosoftConnectionV1Exists("fic_eri_router_to_azure_microsoft_connection_v1.connection_1", &c),
				),
			},
			resource.TestStep{
				Config: testAccConfigEriRouterToAzureMicrosoftConnectionV1Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEriRouterToAzureMicrosoftConnectionV1Exists("fic_eri_router_to_azure_microsoft_connection_v1.connection_1", &c),
				),
			},
		},
	})
}

func testAccCheckEriRouterToAzureMicrosoftConnectionV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.eriV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "fic_eri_router_to_azure_microsoft_connection_v1" {
			continue
		}

		_, err := connections.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Connection still exists")
		}
	}

	return nil
}

func testAccCheckEriRouterToAzureMicrosoftConnectionV1Exists(n string, c *connections.Connection) resource.TestCheckFunc {
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

var testAccConfigEriRouterToAzureMicrosoftConnectionV1Basic = fmt.Sprintf(`
resource "fic_eri_router_v1" "router_1" {
    name = "terraform_router_1"
    area = "%s"
    user_ip_address = "10.0.0.0/27"
    redundant = false
}

resource "fic_eri_nat_component_v1" "nat_1" {
    router_id = "${fic_eri_router_v1.router_1.id}"
    nat_id = "${fic_eri_router_v1.router_1.nat_id}"

    user_ip_addresses = [
        "192.168.0.0/30",
        "192.168.0.4/30",
        "192.168.0.8/30",
        "192.168.0.12/30",
    ]

    global_ip_address_sets  {
        name = "src-set-01"
        type = "sourceNapt"
        number_of_addresses = 5
    }
}

resource "fic_eri_nat_global_ip_address_set_v1" "gip_1" {
    depends_on = ["fic_eri_nat_component_v1.nat_1"]

    router_id = "${fic_eri_router_v1.router_1.id}"
    nat_id = "${fic_eri_router_v1.router_1.nat_id}"

    name = "src-set-02"
    type = "sourceNapt"
    number_of_addresses = 5
}

resource "fic_eri_router_to_azure_microsoft_connection_v1" "connection_1" {
    depends_on = ["fic_eri_nat_global_ip_address_set_v1.gip_1"]

    name = "terraform_connection_1"

    source_router_id = "${fic_eri_router_v1.router_1.id}"
    source_group_name = "group_1"
    source_route_filter_in = "fullRoute"
    source_route_filter_out = "natRoute"

    destination_interconnect = "Osaka-1"
    destination_qos_type = "guarantee"
    destination_service_key = "%s"
    destination_advertised_public_prefixes = [
        "${fic_eri_nat_global_ip_address_set_v1.gip_1.addresses.0}/32"
    ]

    bandwidth = "40M"
}
`,
	OS_AREA_NAME,
	OS_AZURE_SERVICE_KEY,
)

var testAccConfigEriRouterToAzureMicrosoftConnectionV1Update = fmt.Sprintf(`
resource "fic_eri_router_v1" "router_1" {
    name = "terraform_router_1"
    area = "%s"
    user_ip_address = "10.0.0.0/27"
    redundant = false
}

resource "fic_eri_nat_component_v1" "nat_1" {
    router_id = "${fic_eri_router_v1.router_1.id}"
    nat_id = "${fic_eri_router_v1.router_1.nat_id}"

    user_ip_addresses = [
        "192.168.0.0/30",
        "192.168.0.4/30",
        "192.168.0.8/30",
        "192.168.0.12/30",
    ]

    global_ip_address_sets  {
        name = "src-set-01"
        type = "sourceNapt"
        number_of_addresses = 5
    }
}

resource "fic_eri_nat_global_ip_address_set_v1" "gip_1" {
    depends_on = ["fic_eri_nat_component_v1.nat_1"]

    router_id = "${fic_eri_router_v1.router_1.id}"
    nat_id = "${fic_eri_router_v1.router_1.nat_id}"

    name = "src-set-02"
    type = "sourceNapt"
    number_of_addresses = 5
}

resource "fic_eri_router_to_azure_microsoft_connection_v1" "connection_1" {
    depends_on = ["fic_eri_nat_component_v1.nat_1"]

    name = "terraform_connection_1"

    source_router_id = "${fic_eri_router_v1.router_1.id}"
    source_group_name = "group_1"
    source_route_filter_in = "noRoute"
    source_route_filter_out = "noRoute"

    destination_interconnect = "Osaka-1"
    destination_qos_type = "guarantee"
    destination_service_key = "%s"
    destination_advertised_public_prefixes = [
        "${fic_eri_nat_global_ip_address_set_v1.gip_1.addresses.1}/32"
    ]

    bandwidth = "40M"
}
`,
	OS_AREA_NAME,
	OS_AZURE_SERVICE_KEY,
)
