package fic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccEriV1SwitchDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEriV1SwitchDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.fic_eri_switch_v1.sw1", "id"),
					resource.TestCheckResourceAttr("data.fic_eri_switch_v1.sw1", "name", OS_SWITCH_NAME),
					resource.TestCheckResourceAttr("data.fic_eri_switch_v1.sw1", "area", OS_AREA_NAME),
					resource.TestCheckResourceAttr("data.fic_eri_switch_v1.sw1", "location", "NTTComTokyo(NW1)"),
					resource.TestCheckResourceAttrSet("data.fic_eri_switch_v1.sw1", "port_types.0.port_type"),
					resource.TestCheckResourceAttrSet("data.fic_eri_switch_v1.sw1", "port_types.0.available"),
					resource.TestCheckResourceAttrSet("data.fic_eri_switch_v1.sw1", "number_of_available_vlans"),
					resource.TestCheckResourceAttrSet("data.fic_eri_switch_v1.sw1", "vlan_ranges.0.vlan_range"),
					resource.TestCheckResourceAttrSet("data.fic_eri_switch_v1.sw1", "vlan_ranges.0.available"),
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
}
`,
	OS_SWITCH_NAME,
	OS_AREA_NAME,
)
