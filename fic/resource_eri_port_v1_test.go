package fic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/nttcom/go-fic/fic/eri/v1/ports"
)

func TestAccEriPortV1Basic(t *testing.T) {
	var port ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckSwitchName(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEriPortV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccConfigEriPortV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEriPortV1Exists("fic_eri_port_v1.port_1", &port),
					resource.TestCheckResourceAttr(
						"fic_eri_port_v1.port_1", "name", "terraform_port_1"),
					resource.TestCheckResourceAttr(
						"fic_eri_port_v1.port_1", "switch_name", OS_SWITCH_NAME),
					resource.TestCheckResourceAttr(
						"fic_eri_port_v1.port_1", "port_type", "1G"),
					testAccCheckEriPortV1TagLength(&port, 16),
				),
			},
		},
	})
}

func TestAccEriPortV1CreateWithActiveState(t *testing.T) {
	var port ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckSwitchName(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEriPortV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccConfigEriPortV1WithActiveState,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEriPortV1Exists("fic_eri_port_v1.port_1", &port),
					resource.TestCheckResourceAttr(
						"fic_eri_port_v1.port_1", "name", "terraform_port_1"),
					resource.TestCheckResourceAttr(
						"fic_eri_port_v1.port_1", "switch_name", OS_SWITCH_NAME),
					resource.TestCheckResourceAttr(
						"fic_eri_port_v1.port_1", "port_type", "1G"),
					resource.TestCheckResourceAttr(
						"fic_eri_port_v1.port_1", "is_activated", "true"),
					testAccCheckEriPortV1TagLength(&port, 16),
				),
			},
		},
	})
}

func TestAccEriPortV1ActiveAfterward(t *testing.T) {
	var port ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckSwitchName(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEriPortV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccConfigEriPortV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEriPortV1Exists("fic_eri_port_v1.port_1", &port),
					resource.TestCheckResourceAttr(
						"fic_eri_port_v1.port_1", "is_activated", "false"),
				),
			},
			resource.TestStep{
				Config: testAccConfigEriPortV1WithActiveState,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEriPortV1Exists("fic_eri_port_v1.port_1", &port),
					resource.TestCheckResourceAttr(
						"fic_eri_port_v1.port_1", "is_activated", "true"),
				),
			},
		},
	})
}

func TestAccEriPortV1CreateWithVLANRanges(t *testing.T) {
	var port ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckSwitchName(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEriPortV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccConfigEriPortV1CreateWithVLANRanges,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEriPortV1Exists("fic_eri_port_v1.port_1", &port),
					resource.TestCheckResourceAttr(
						"fic_eri_port_v1.port_1", "is_activated", "false"),
					testAccCheckEriPortV1TagLength(&port, 32),
				),
			},
			resource.TestStep{
				Config: testAccConfigEriPortV1CreateWithVLANRangesActivateAfterward,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEriPortV1Exists("fic_eri_port_v1.port_1", &port),
					resource.TestCheckResourceAttr(
						"fic_eri_port_v1.port_1", "is_activated", "true"),
				),
			},
		},
	})
}

func testAccCheckEriPortV1TagLength(port *ports.Port, length int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(port.VLANs) != length {
			return fmt.Errorf(
				"Tag length is different: expected %d, actual %d",
				length, len(port.VLANs))
		}
		return nil
	}
}

func testAccCheckEriPortV1Destroy(s *terraform.State) error {
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
			return fmt.Errorf("Port still exists")
		}
	}

	return nil
}

func testAccCheckEriPortV1Exists(n string, port *ports.Port) resource.TestCheckFunc {
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

		found, err := ports.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Port not found")
		}

		*port = *found

		return nil
	}
}

var testAccConfigEriPortV1Basic = fmt.Sprintf(`
resource "fic_eri_port_v1" "port_1" {
	name = "terraform_port_1"
	switch_name = "%s"
	port_type = "1G"
	number_of_vlans = 16
}
`,
	OS_SWITCH_NAME,
)

var testAccConfigEriPortV1WithActiveState = fmt.Sprintf(`
resource "fic_eri_port_v1" "port_1" {
	name = "terraform_port_1"
	switch_name = "%s"
	port_type = "1G"
	number_of_vlans = 16
	is_activated = true
}
`,
	OS_SWITCH_NAME,
)

var testAccConfigEriPortV1CreateWithVLANRanges = fmt.Sprintf(`
resource "fic_eri_port_v1" "port_1" {
	name = "terraform_port_1"
	switch_name = "%s"
	port_type = "1G"

	vlan_ranges {
		start = 1137
		end = 1152
	}

	vlan_ranges {
		start = 1153
		end = 1168
	}
}
`,
	OS_SWITCH_NAME,
)

var testAccConfigEriPortV1CreateWithVLANRangesActivateAfterward = fmt.Sprintf(`
resource "fic_eri_port_v1" "port_1" {
	name = "terraform_port_1"
	switch_name = "%s"
	port_type = "1G"
	is_activated = true

	vlan_ranges {
		start = 1137
		end = 1152
	}

	vlan_ranges {
		start = 1153
		end = 1168
	}
}
`,
	OS_SWITCH_NAME,
)
