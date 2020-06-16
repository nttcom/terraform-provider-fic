package fic

import (
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceEriRouterPairedToGCPConnectionV1() *schema.Resource {
	validInterconnects := []string{
		"Equinix-TY2-2", "@Tokyo-CC2-2", "Equinix-TY2-3", "@Tokyo-CC2-3", "Equinix-OS1-1", "NTT-Dojima2-1",
	}
	destinationSchema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"interconnect": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(validInterconnects, false),
			},
			"pairing_key": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-fA-F0-9]{8}(-[a-fA-F0-9]{4}){3}-[a-fA-F0-9]{12}/[a-zA-Z0-9-]*/[1,2]$`), ""),
			},
		},
	}

	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"router_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"group_name": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"group_1", "group_2", "group_3", "group_4", "group_5", "group_6", "group_7", "group_8"}, false),
						},
						"route_filter": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"in": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"fullRoute", "noRoute"}, false),
									},
									"out": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"fullRoute", "fullRouteWithDefaultRoute", "defaultRoute", "privateRoute", "noRoute"}, false),
									},
								},
							},
						},
						"primary_med_out": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: IntInSlice([]int{10, 30}),
						},
						"secondary_med_out": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: IntInSlice([]int{20, 40}),
						},
					},
				},
			},
			"destination": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"primary": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem:     destinationSchema,
						},
						"secondary": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem:     destinationSchema,
						},
					},
				},
			},
			"bandwidth": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"10M", "50M", "100M", "200M", "300M", "400M", "500M", "1G", "2G", "5G", "10G"}, false),
			},
		},
	}
}
