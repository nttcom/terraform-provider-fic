package fic

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/nttcom/go-fic/fic/eri/v1/switches"
)

func dataSourceEriSwitchV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEriSwitchV1Read,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"area": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"location": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"port_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"available": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"number_of_available_vlans": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"vlan_ranges": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vlan_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"available": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceEriSwitchV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return err
	}

	pages, err := switches.List(client, nil).AllPages()
	if err != nil {
		return fmt.Errorf("unable to retrieve switches: %s", err)
	}

	sws, err := switches.ExtractSwitches(pages)
	if err != nil {
		return fmt.Errorf("unable to extract switches: %s", err)
	}

	opts := struct {
		switchName string
		area       string
		location   string
	}{}

	if v, ok := d.GetOk("name"); ok {
		opts.switchName = v.(string)
	}

	if v, ok := d.GetOk("area"); ok {
		opts.area = v.(string)
	}

	if v, ok := d.GetOk("location"); ok {
		opts.location = v.(string)
	}

	var matches []switches.Switch
	for _, sw := range sws {
		if opts.switchName != "" && opts.switchName != sw.SwitchName {
			continue
		}

		if opts.area != "" && opts.area != sw.Area {
			continue
		}

		if opts.location != "" && opts.location != sw.Location {
			continue
		}

		matches = append(matches, sw)
	}

	if len(matches) == 0 {
		return fmt.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	if len(matches) >= 2 {
		return fmt.Errorf("your query returned more than one result. Please try a more specific search criteria")
	}

	match := matches[0]

	log.Printf("[DEBUG] Retrieved Eri Switch %s: %+v", match.ID, match)
	d.SetId(match.ID)

	d.Set("name", match.SwitchName)
	d.Set("area", match.Area)
	d.Set("location", match.Location)
	d.Set("number_of_available_vlans", match.NumberOfAvailableVLANs)

	var portTypes []map[string]interface{}
	for _, pt := range match.PortTypes {
		portTypes = append(portTypes, map[string]interface{}{
			"port_type": pt.Type,
			"available": pt.Available,
		})
	}
	d.Set("port_types", portTypes)

	var vlanRanges []map[string]interface{}
	for _, vr := range match.VLANRanges {
		vlanRanges = append(vlanRanges, map[string]interface{}{
			"vlan_range": vr.Range,
			"available":  vr.Available,
		})
	}
	d.Set("vlan_ranges", vlanRanges)

	return nil
}
