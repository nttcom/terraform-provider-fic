package fic

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/nttcom/go-fic/fic/eri/v1/routers/components/firewalls"
)

func TestAccEriFirewallComponentV1WithFirewallConfigurations(t *testing.T) {
	var f firewalls.Firewall

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckArea(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEriFirewallComponentV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccConfigEriFirewallComponentV1Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEriFirewallComponentV1Exists("fic_eri_firewall_component_v1.firewall_1", &f),

					resource.TestCheckResourceAttr(
						"fic_eri_firewall_component_v1.firewall_1", "rules.0.from", "group_1"),
					resource.TestCheckResourceAttr(
						"fic_eri_firewall_component_v1.firewall_1", "rules.0.to", "group_2"),
					resource.TestCheckResourceAttr(
						"fic_eri_firewall_component_v1.firewall_1", "rules.0.entries.0.name", "rule-01"),
					resource.TestCheckResourceAttr(
						"fic_eri_firewall_component_v1.firewall_1", "rules.0.entries.0.action", "permit"),

					resource.TestCheckResourceAttr(
						"fic_eri_firewall_component_v1.firewall_1", "custom_applications.0.name", "google-drive-web"),
					resource.TestCheckResourceAttr(
						"fic_eri_firewall_component_v1.firewall_1", "custom_applications.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"fic_eri_firewall_component_v1.firewall_1", "custom_applications.0.destination_port", "443"),

					resource.TestCheckResourceAttr(
						"fic_eri_firewall_component_v1.firewall_1", "application_sets.0.name", "app_set_1"),
					resource.TestCheckResourceAttr(
						"fic_eri_firewall_component_v1.firewall_1", "application_sets.0.applications.0", "google-drive-web"),
					resource.TestCheckResourceAttr(
						"fic_eri_firewall_component_v1.firewall_1", "application_sets.0.applications.1", "pre-defined-ftp"),

					resource.TestCheckResourceAttr(
						"fic_eri_firewall_component_v1.firewall_1", "routing_group_settings.0.group_name", "group_1"),
					resource.TestCheckResourceAttr(
						"fic_eri_firewall_component_v1.firewall_1", "routing_group_settings.0.address_sets.0.name", "group1_addset_1"),
					resource.TestCheckResourceAttr(
						"fic_eri_firewall_component_v1.firewall_1", "routing_group_settings.0.address_sets.0.addresses.0", "172.18.1.0/24"),
				),
			},
		},
	})
}

func TestAccEriFirewallComponentV1ActivateThenUpdateFirewallRules(t *testing.T) {
	var f firewalls.Firewall

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckArea(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEriFirewallComponentV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccConfigEriFirewallComponentV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEriFirewallComponentV1Exists("fic_eri_firewall_component_v1.firewall_1", &f),
				),
			},
			resource.TestStep{
				Config: testAccConfigEriFirewallComponentV1Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEriFirewallComponentV1Exists("fic_eri_firewall_component_v1.firewall_1", &f),

					resource.TestCheckResourceAttr(
						"fic_eri_firewall_component_v1.firewall_1", "rules.0.from", "group_1"),
					resource.TestCheckResourceAttr(
						"fic_eri_firewall_component_v1.firewall_1", "rules.0.to", "group_2"),
					resource.TestCheckResourceAttr(
						"fic_eri_firewall_component_v1.firewall_1", "rules.0.entries.0.name", "rule-01"),
					resource.TestCheckResourceAttr(
						"fic_eri_firewall_component_v1.firewall_1", "rules.0.entries.0.action", "permit"),

					resource.TestCheckResourceAttr(
						"fic_eri_firewall_component_v1.firewall_1", "custom_applications.0.name", "google-drive-web"),
					resource.TestCheckResourceAttr(
						"fic_eri_firewall_component_v1.firewall_1", "custom_applications.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"fic_eri_firewall_component_v1.firewall_1", "custom_applications.0.destination_port", "443"),

					resource.TestCheckResourceAttr(
						"fic_eri_firewall_component_v1.firewall_1", "application_sets.0.name", "app_set_1"),
					resource.TestCheckResourceAttr(
						"fic_eri_firewall_component_v1.firewall_1", "application_sets.0.applications.0", "google-drive-web"),
					resource.TestCheckResourceAttr(
						"fic_eri_firewall_component_v1.firewall_1", "application_sets.0.applications.1", "pre-defined-ftp"),

					resource.TestCheckResourceAttr(
						"fic_eri_firewall_component_v1.firewall_1", "routing_group_settings.0.group_name", "group_1"),
					resource.TestCheckResourceAttr(
						"fic_eri_firewall_component_v1.firewall_1", "routing_group_settings.0.address_sets.0.name", "group1_addset_1"),
					resource.TestCheckResourceAttr(
						"fic_eri_firewall_component_v1.firewall_1", "routing_group_settings.0.address_sets.0.addresses.0", "172.18.1.0/24"),
				),
			},
		},
	})
}

func testAccCheckEriFirewallComponentV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.eriV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "fic_eri_firewall_component_v1" {
			continue
		}

		id := rs.Primary.ID
		routerID := strings.Split(id, "/")[0]
		firewallID := strings.Split(id, "/")[1]
		_, err := firewalls.Get(client, routerID, firewallID).Extract()
		if err == nil {
			return fmt.Errorf("Firewall Component still exists")
		}
	}

	return nil
}

func testAccCheckEriFirewallComponentV1Exists(n string, f *firewalls.Firewall) resource.TestCheckFunc {
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
		firewallID := strings.Split(id, "/")[1]
		found, err := firewalls.Get(client, routerID, firewallID).Extract()
		if err != nil {
			return err
		}

		if found.ID != firewallID {
			return fmt.Errorf("Firewall Component not found")
		}

		*f = *found

		return nil
	}
}

var testAccConfigEriFirewallComponentV1Router = fmt.Sprintf(`
resource "fic_eri_router_v1" "router_1" {
	name = "terraform_router_1"
	area = "%s"
	user_ip_address = "10.0.0.0/27"
	redundant = true
}`,
	OS_AREA_NAME,
)

var testAccConfigEriFirewallComponentV1Basic = fmt.Sprintf(`
%s

resource "fic_eri_firewall_component_v1" "firewall_1" {
	router_id = "${fic_eri_router_v1.router_1.id}"
	firewall_id = "${fic_eri_router_v1.router_1.firewall_id}"

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
}
`,
	testAccConfigEriFirewallComponentV1Router,
)

var testAccConfigEriFirewallComponentV1Update = fmt.Sprintf(`
%s

resource "fic_eri_firewall_component_v1" "firewall_1" {
    router_id = "${fic_eri_router_v1.router_1.id}"
    firewall_id = "${fic_eri_router_v1.router_1.firewall_id}"

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

	rules {
		from = "group_1"
        to = "group_2"
		entries {
			name = "rule-01"
			match_source_address_sets = [
				 "group1_addset_1"
			]
			match_destination_address_sets = [
				 "group2_addset_1"
			]
			match_application = "app_set_1"
			action = "permit"
		}
	}

	custom_applications {
    	name = "google-drive-web"
    	protocol = "tcp"
		destination_port = "443"
	}

	application_sets {
    	name = "app_set_1"
    	applications = [
			"google-drive-web",
            "pre-defined-ftp"
		]
	}

	routing_group_settings {
		group_name = "group_1"
		address_sets {
			name = "group1_addset_1"
			addresses = [
				"172.18.1.0/24"
			]
		}
	}

	routing_group_settings {
		group_name = "group_2"
		address_sets {
			name = "group2_addset_1"
			addresses = [
				"192.168.1.0/24"
			]
		}
	}
}
`,
	testAccConfigEriFirewallComponentV1Router,
)
