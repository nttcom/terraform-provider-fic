package fic

import (
	"fmt"
	"regexp"

	connections "github.com/nttcom/go-fic/fic/eri/v1/router_paired_to_gcp_connections"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceEriRouterPairedToGCPConnectionV1() *schema.Resource {
	validInterconnects := []string{
		"Equinix-TY2-2", "@Tokyo-CC2-2", "Equinix-TY2-3", "@Tokyo-CC2-3", "Equinix-OS1-1", "NTT-Dojima2-1",
	}
	interconnectSchema := &schema.Resource{
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
		Read: resourceEriRouterPairedToGCPConnectionV1Read,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bandwidth": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"10M", "50M", "100M", "200M", "300M", "400M", "500M", "1G", "2G", "5G", "10G"}, false),
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
							Elem:     interconnectSchema,
						},
						"secondary": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem:     interconnectSchema,
						},
						"qos_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"redundant": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"area": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operation_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operation_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"primary_connected_network_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secondary_connected_network_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceEriRouterPairedToGCPConnectionV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating FIC ERI client: %s", err)
	}

	conn, err := connections.Get(client, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "connection")
	}

	d.Set("name", conn.Name)
	d.Set("bandwidth", conn.Bandwidth)
	d.Set("source", flattenSource(conn.Source))
	d.Set("destination", flattenDestination(conn.Destination))
	d.Set("redundant", conn.Redundant)
	d.Set("tenant_id", conn.TenantID)
	//d.Set("area", conn.Area)
	d.Set("operation_id", conn.OperationID)
	d.Set("operation_status", conn.OperationStatus)
	d.Set("primary_connected_network_address", conn.PrimaryConnectedNetworkAddress)
	d.Set("secondary_connected_network_address", conn.SecondaryConnectedNetworkAddress)

	return nil
}

func flattenSource(in connections.Source) []interface{} {
	var out []interface{}
	m := make(map[string]interface{})

	m["router_id"] = in.RouterID
	m["group_name"] = in.GroupName
	m["route_filter"] = flattenRouteFilter(in.RouteFilter)

	out = append(out, m)
	return out
}

func flattenRouteFilter(in connections.RouteFilter) []interface{} {
	var out []interface{}
	m := make(map[string]interface{})

	m["in"] = in.In
	m["out"] = in.Out

	out = append(out, m)
	return out
}

func flattenDestination(in connections.Destination) []interface{} {
	var out []interface{}
	m := make(map[string]interface{})

	m["primary"] = flattenInterconnect(in.Primary)
	m["secondary"] = flattenInterconnect(in.Secondary)
	m["qos_type"] = in.QosType

	out = append(out, m)
	return out
}

func flattenInterconnect(in connections.DestinationHAInfo) []interface{} {
	var out []interface{}
	m := make(map[string]interface{})

	m["interconnect"] = in.Interconnect
	m["paring_key"] = in.PairingKey

	out = append(out, m)
	return out
}
