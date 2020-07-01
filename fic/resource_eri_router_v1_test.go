package fic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/nttcom/go-fic/fic/eri/v1/ports"
	"github.com/nttcom/go-fic/fic/eri/v1/routers"
)

func TestAccEriRouterV1Basic(t *testing.T) {
	var router routers.Router

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckArea(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEriRouterV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccConfigEriRouterV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEriRouterV1Exists("fic_eri_router_v1.router_1", &router),
					resource.TestCheckResourceAttr(
						"fic_eri_router_v1.router_1", "name", "terraform_router_1"),
					resource.TestCheckResourceAttr(
						"fic_eri_router_v1.router_1", "redundant", "true"),
					resource.TestCheckResourceAttr(
						"fic_eri_router_v1.router_1", "area", OS_AREA_NAME),
					resource.TestCheckResourceAttr(
						"fic_eri_router_v1.router_1", "user_ip_address", "10.0.0.0/27"),
					resource.TestCheckResourceAttr(
						"fic_eri_router_v1.router_1", "firewalls.0.is_activated", "false"),
					resource.TestCheckResourceAttr(
						"fic_eri_router_v1.router_1", "nats.0.is_activated", "false"),
				),
			},
		},
	})
}

func TestAccEriRouterV1RedundantFalse(t *testing.T) {
	var router routers.Router

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckArea(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEriRouterV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccConfigEriRouterV1RedundantFalse,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEriRouterV1Exists("fic_eri_router_v1.router_1", &router),
					resource.TestCheckResourceAttr(
						"fic_eri_router_v1.router_1", "name", "terraform_router_1"),
					resource.TestCheckResourceAttr(
						"fic_eri_router_v1.router_1", "redundant", "false"),
					resource.TestCheckResourceAttr(
						"fic_eri_router_v1.router_1", "area", OS_AREA_NAME),
					resource.TestCheckResourceAttr(
						"fic_eri_router_v1.router_1", "user_ip_address", "10.0.0.0/27"),
					resource.TestCheckResourceAttr(
						"fic_eri_router_v1.router_1", "firewalls.0.is_activated", "false"),
					resource.TestCheckResourceAttr(
						"fic_eri_router_v1.router_1", "nats.0.is_activated", "false"),
				),
			},
		},
	})
}

func testAccCheckEriRouterV1TagLength(port *ports.Port, length int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(port.VLANs) != length {
			return fmt.Errorf(
				"Tag length is different: expected %d, actual %d",
				length, len(port.VLANs))
		}
		return nil
	}
}

func testAccCheckEriRouterV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.eriV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "fic_eri_port_v1" {
			continue
		}

		_, err := ports.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Router still exists")
		}
	}

	return nil
}

func testAccCheckEriRouterV1Exists(n string, router *routers.Router) resource.TestCheckFunc {
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

		found, err := routers.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Router not found")
		}

		*router = *found

		return nil
	}
}

var testAccConfigEriRouterV1Basic = fmt.Sprintf(`
resource "fic_eri_router_v1" "router_1" {
	name = "terraform_router_1"
	area = "%s"
	user_ip_address = "10.0.0.0/27"
	redundant = true
}
`,
	OS_AREA_NAME,
)

var testAccConfigEriRouterV1RedundantFalse = fmt.Sprintf(`
resource "fic_eri_router_v1" "router_1" {
	name = "terraform_router_1"
	area = "%s"
	user_ip_address = "10.0.0.0/27"
	redundant = false
}
`,
	OS_AREA_NAME,
)
