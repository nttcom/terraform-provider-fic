package fic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccEriV1SwitchDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEriV1SwitchDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.fic_eri_switch_v1.sw1", "id"),
					resource.TestCheckResourceAttr(
						"data.fic_eri_switch_v1.sw1", "name", OS_SWITCH_NAME),
					resource.TestCheckResourceAttr(
						"data.fic_eri_switch_v1.sw1", "area", OS_AREA_NAME),
					resource.TestCheckResourceAttr(
						"data.fic_eri_switch_v1.sw1", "location", "NTTComTokyo(NW1)"),
					resource.TestCheckResourceAttr(
						"data.fic_eri_switch_v1.sw1", "port_type", "1G"),
					resource.TestCheckResourceAttrSet(
						"data.fic_eri_switch_v1.sw1", "number_of_available_vlans"),
					resource.TestCheckResourceAttrSet(
						"data.fic_eri_switch_v1.sw1", "vlan_ranges.0.start"),
					resource.TestCheckResourceAttrSet(
						"data.fic_eri_switch_v1.sw1", "vlan_ranges.0.end"),
				),
			},
		},
	})
}

var testAccEriV1SwitchDataSourceBasic = fmt.Sprintf(`
data "fic_eri_switch_v1" "sw1" {
	name = "%s"
	area = "%s"
	location = "NTTComTokyo(NW1)"
	port_type = "1G"
}
`,
	OS_SWITCH_NAME,
	OS_AREA_NAME,
)
