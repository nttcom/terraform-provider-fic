package fic

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/nttcom/go-fic/fic/eri/v1/routers/components/nats"
)

func TestAccEriNATComponentV1Basic(t *testing.T) {
	var nat nats.NAT

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckArea(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEriNATComponentV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccConfigEriNATComponentV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEriNATComponentV1Exists("fic_eri_nat_component_v1.nat_1", &nat),
					resource.TestCheckResourceAttr(
						"fic_eri_nat_component_v1.nat_1", "user_ip_addresses.0", "192.168.0.0/30"),
					resource.TestCheckResourceAttr(
						"fic_eri_nat_component_v1.nat_1", "user_ip_addresses.7", "192.168.28.0/30"),
				),
			},
		},
	})
}

func TestAccEriNATComponentV1WithNATAndNAPTRules(t *testing.T) {
	var nat nats.NAT

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckArea(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEriNATComponentV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccConfigEriNATComponentV1WithNATAndNAPTRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEriNATComponentV1Exists("fic_eri_nat_component_v1.nat_1", &nat),

					resource.TestCheckResourceAttr(
						"fic_eri_nat_component_v1.nat_1", "source_napt_rules.0.from.0", "group_1"),
					resource.TestCheckResourceAttr(
						"fic_eri_nat_component_v1.nat_1", "source_napt_rules.0.to", "group_2"),
					resource.TestCheckResourceAttr(
						"fic_eri_nat_component_v1.nat_1",
						"source_napt_rules.0.entries.0.then.0", "src-set-01"),

					resource.TestCheckResourceAttr(
						"fic_eri_nat_component_v1.nat_1", "destination_nat_rules.0.from", "group_1"),
					resource.TestCheckResourceAttr(
						"fic_eri_nat_component_v1.nat_1", "destination_nat_rules.0.to", "group_2"),
					resource.TestCheckResourceAttr(
						"fic_eri_nat_component_v1.nat_1",
						"destination_nat_rules.0.entries.0.match_destination_address", "dst-set-01"),
					resource.TestCheckResourceAttr(
						"fic_eri_nat_component_v1.nat_1",
						"destination_nat_rules.0.entries.0.then", "192.168.0.1/32"),
				),
			},
		},
	})
}

func TestAccEriNATComponentV1ActivateThenUpdateNAPTAndNATRules(t *testing.T) {
	var nat nats.NAT

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckArea(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEriNATComponentV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccConfigEriNATComponentV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEriNATComponentV1Exists("fic_eri_nat_component_v1.nat_1", &nat),
				),
			},
			resource.TestStep{
				Config: testAccConfigEriNATComponentV1WithNATAndNAPTRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEriNATComponentV1Exists("fic_eri_nat_component_v1.nat_1", &nat),

					resource.TestCheckResourceAttr(
						"fic_eri_nat_component_v1.nat_1", "source_napt_rules.0.from.0", "group_1"),
					resource.TestCheckResourceAttr(
						"fic_eri_nat_component_v1.nat_1", "source_napt_rules.0.to", "group_2"),
					resource.TestCheckResourceAttr(
						"fic_eri_nat_component_v1.nat_1",
						"source_napt_rules.0.entries.0.then.0", "src-set-01"),

					resource.TestCheckResourceAttr(
						"fic_eri_nat_component_v1.nat_1", "destination_nat_rules.0.from", "group_1"),
					resource.TestCheckResourceAttr(
						"fic_eri_nat_component_v1.nat_1", "destination_nat_rules.0.to", "group_2"),
					resource.TestCheckResourceAttr(
						"fic_eri_nat_component_v1.nat_1",
						"destination_nat_rules.0.entries.0.match_destination_address", "dst-set-01"),
					resource.TestCheckResourceAttr(
						"fic_eri_nat_component_v1.nat_1",
						"destination_nat_rules.0.entries.0.then", "192.168.0.1/32"),
				),
			},
		},
	})
}

func testAccCheckEriNATComponentV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.eriV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "fic_eri_nat_component_v1" {
			continue
		}

		id := rs.Primary.ID
		routerID := strings.Split(id, "/")[0]
		natID := strings.Split(id, "/")[1]
		_, err := nats.Get(client, routerID, natID).Extract()
		if err == nil {
			return fmt.Errorf("NAT Component still exists")
		}
	}

	return nil
}

func testAccCheckEriNATComponentV1Exists(n string, nat *nats.NAT) resource.TestCheckFunc {
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
		found, err := nats.Get(client, routerID, natID).Extract()
		if err != nil {
			return err
		}

		if found.ID != natID {
			return fmt.Errorf("NAT Component not found")
		}

		*nat = *found

		return nil
	}
}

var testAccConfigEriNATComponentV1Router = fmt.Sprintf(`
resource "fic_eri_router_v1" "router_1" {
	name = "terraform_router_1"
	area = "%s"
	user_ip_address = "10.0.0.0/27"
	redundant = true
}`,
	OS_AREA_NAME,
)

var testAccConfigEriNATComponentV1Basic = fmt.Sprintf(`
%s

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
`,
	testAccConfigEriNATComponentV1Router,
)

var testAccConfigEriNATComponentV1WithNATAndNAPTRules = fmt.Sprintf(`
%s

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
`,
	testAccConfigEriNATComponentV1Router,
)
