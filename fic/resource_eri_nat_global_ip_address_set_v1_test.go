package fic

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/nttcom/go-fic/fic/eri/v1/routers/components/nat_global_ip_address_sets"
)

func TestAccEriNATGlobalIPAddressSetV1Basic(t *testing.T) {
	var gip nat_global_ip_address_sets.GlobalIPAddressSet

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckArea(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEriNATGlobalIPAddressSetV1Destroy,
		Steps: []resource.TestStep{

			// Create Global IP Address Set under router nad NAT component.
			resource.TestStep{
				Config: testAccConfigEriNATGlobalIPAddressSetV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEriNATGlobalIPAddressSetV1Exists("fic_eri_nat_global_ip_address_set_v1.gip_1", &gip),
					resource.TestCheckResourceAttr(
						"fic_eri_nat_global_ip_address_set_v1.gip_1", "name", "src-set-02"),
					resource.TestCheckResourceAttr(
						"fic_eri_nat_global_ip_address_set_v1.gip_1", "type", "sourceNapt"),
					resource.TestCheckResourceAttr(
						"fic_eri_nat_global_ip_address_set_v1.gip_1", "number_of_addresses", "5"),
				),
			},

			// Update source NAPT rule of NAT Component by using created global ip address.
			resource.TestStep{
				Config: testAccConfigEriNATGlobalIPAddressSetV1UpdateNATByUsingNewGlobalIPAddressSet,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEriNATGlobalIPAddressSetV1Exists("fic_eri_nat_global_ip_address_set_v1.gip_1", &gip),
				),
			},

			// Update source NAPT rule not to use global ip address previously created.
			resource.TestStep{
				Config: testAccConfigEriNATGlobalIPAddressSetV1UpdateRule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEriNATGlobalIPAddressSetV1Exists("fic_eri_nat_global_ip_address_set_v1.gip_1", &gip),
				),
			},
		},
	})
}

func testAccCheckEriNATGlobalIPAddressSetV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.eriV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "fic_eri_nat_global_ip_address_set_v1" {
			continue
		}

		id := rs.Primary.ID
		routerID := strings.Split(id, "/")[0]
		natID := strings.Split(id, "/")[1]
		globalIPAddressSetID := strings.Split(id, "/")[1]
		_, err := nat_global_ip_address_sets.Get(
			client, routerID, natID, globalIPAddressSetID).Extract()
		if err == nil {
			return fmt.Errorf("NAT Global IP Address Set still exists")
		}
	}

	return nil
}

func testAccCheckEriNATGlobalIPAddressSetV1Exists(n string, gip *nat_global_ip_address_sets.GlobalIPAddressSet) resource.TestCheckFunc {
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

		id := rs.Primary.ID
		routerID := strings.Split(id, "/")[0]
		natID := strings.Split(id, "/")[1]
		globalIPAddressSetID := strings.Split(id, "/")[2]
		found, err := nat_global_ip_address_sets.Get(
			client, routerID, natID, globalIPAddressSetID).Extract()
		if err != nil {
			return err
		}

		if found.ID != globalIPAddressSetID {
			return fmt.Errorf("NAT Global IP Address Set not found")
		}

		*gip = *found

		return nil
	}
}

var testAccConfigEriNATGlobalIPAddressSetV1Basic = fmt.Sprintf(`
resource "fic_eri_router_v1" "router_1" {
	name = "terraform_router_1"
	area = "%s"
	user_ip_address = "10.0.0.0/27"
	redundant = true
}

resource "fic_eri_nat_component_v1" "nat_1" {
	router_id = "${fic_eri_router_v1.router_1.id}"
	nat_id = "${fic_eri_router_v1.router_1.nat_id}"

	user_ip_addresses = [
		"192.168.0.0/30",
        "192.168.4.0/30",
        "192.168.8.0/30",
        "192.168.12.0/30",
		"192.168.16.0/30",
        "192.168.20.0/30",
        "192.168.24.0/30",
        "192.168.28.0/30"
    ]

	global_ip_address_sets  {
        name = "src-set-01"
        type = "sourceNapt"
        number_of_addresses = 5
	}

	global_ip_address_sets  {
        name = "dst-set-01"
        type = "destinationNat"
        number_of_addresses = 1
	}
}

resource "fic_eri_nat_global_ip_address_set_v1" "gip_1" {
	router_id = "${fic_eri_router_v1.router_1.id}"
	nat_id = "${fic_eri_router_v1.router_1.nat_id}"
	depends_on = ["fic_eri_nat_component_v1.nat_1"]

	name = "src-set-02"
	type = "sourceNapt"
	number_of_addresses = 5
}
`,
	OS_AREA_NAME,
)

var testAccConfigEriNATGlobalIPAddressSetV1UpdateNATByUsingNewGlobalIPAddressSet = fmt.Sprintf(`
resource "fic_eri_router_v1" "router_1" {
	name = "terraform_router_1"
	area = "%s"
	user_ip_address = "10.0.0.0/27"
	redundant = true
}

resource "fic_eri_nat_component_v1" "nat_1" {
	router_id = "${fic_eri_router_v1.router_1.id}"
	nat_id = "${fic_eri_router_v1.router_1.nat_id}"

	user_ip_addresses = [
		"192.168.0.0/30",
        "192.168.4.0/30",
        "192.168.8.0/30",
        "192.168.12.0/30",
		"192.168.16.0/30",
        "192.168.20.0/30",
        "192.168.24.0/30",
        "192.168.28.0/30"
    ]

	global_ip_address_sets  {
        name = "src-set-01"
        type = "sourceNapt"
        number_of_addresses = 5
	}

	global_ip_address_sets  {
        name = "dst-set-01"
        type = "destinationNat"
        number_of_addresses = 1
	}

	source_napt_rules {
		from = [
			"group_1"
		]

		to = "group_2"

		entries {
			then = [
            	"src-set-02"
            ]
		}
	}

	destination_nat_rules {
		from = "group_1"
		to = "group_2"
		entries {
			match_destination_address = "dst-set-01"
			then = "192.168.0.1/32"
		}
	}
}

resource "fic_eri_nat_global_ip_address_set_v1" "gip_1" {
	router_id = "${fic_eri_router_v1.router_1.id}"
	nat_id = "${fic_eri_router_v1.router_1.nat_id}"
	depends_on = ["fic_eri_nat_component_v1.nat_1"]

	name = "src-set-02"
	type = "sourceNapt"
	number_of_addresses = 5
}
`, OS_AREA_NAME,
)

var testAccConfigEriNATGlobalIPAddressSetV1UpdateRule = fmt.Sprintf(`
resource "fic_eri_router_v1" "router_1" {
	name = "terraform_router_1"
	area = "%s"
	user_ip_address = "10.0.0.0/27"
	redundant = true
}

resource "fic_eri_nat_component_v1" "nat_1" {
	router_id = "${fic_eri_router_v1.router_1.id}"
	nat_id = "${fic_eri_router_v1.router_1.nat_id}"

	user_ip_addresses = [
		"192.168.0.0/30",
        "192.168.4.0/30",
        "192.168.8.0/30",
        "192.168.12.0/30",
		"192.168.16.0/30",
        "192.168.20.0/30",
        "192.168.24.0/30",
        "192.168.28.0/30"
    ]

	global_ip_address_sets  {
        name = "src-set-01"
        type = "sourceNapt"
        number_of_addresses = 5
	}

	global_ip_address_sets  {
        name = "dst-set-01"
        type = "destinationNat"
        number_of_addresses = 1
	}

	source_napt_rules {
		from = [
			"group_1"
		]

		to = "group_2"

		entries {
			then = [
            	"src-set-01"
            ]
		}
	}

	destination_nat_rules {
		from = "group_1"
		to = "group_2"
		entries {
			match_destination_address = "dst-set-01"
			then = "192.168.0.1/32"
		}
	}


}

resource "fic_eri_nat_global_ip_address_set_v1" "gip_1" {
	router_id = "${fic_eri_router_v1.router_1.id}"
	nat_id = "${fic_eri_router_v1.router_1.nat_id}"
	depends_on = ["fic_eri_nat_component_v1.nat_1"]

	name = "src-set-02"
	type = "sourceNapt"
	number_of_addresses = 5
}
`,
	OS_AREA_NAME,
)
