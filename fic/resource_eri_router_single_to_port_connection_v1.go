package fic

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	connections "github.com/nttcom/go-fic/fic/eri/v1/router_single_to_port_connections"
)

func resourceEriRouterSingleToPortConnectionV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceEriRouterSingleToPortConnectionV1Create,
		Read:   resourceEriRouterSingleToPortConnectionV1Read,
		Update: resourceEriRouterSingleToPortConnectionV1Update,
		Delete: resourceEriRouterSingleToPortConnectionV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"source_router_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"source_group_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"group_1", "group_2", "group_3", "group_4",
					"group_5", "group_6", "group_7", "group_8",
				}, false),
			},

			"source_route_filter_in": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				// ForceNew: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"fullRoute", "noRoute"}, false),
			},

			"source_route_filter_out": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				// ForceNew: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"fullRoute", "fullRouteWithDefaultRoute", "noRoute"}, false),
			},

			"source_information": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				// ForceNew: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"as_path_prepend_in": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice(
								[]string{"OFF", "1", "2", "3", "4", "5"}, false),
						},
						"as_path_prepend_out": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice(
								[]string{"OFF", "1", "2", "3", "4", "5"}, false),
						},
					},
				},
			},

			"destination_information": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				// ForceNew: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port_id": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"vlan": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"ip_address": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"asn": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"bandwidth": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"10M", "20M", "30M", "40M", "50M", "100M", "200M", "300M", "400M", "500M",
					"1G", "2G", "3G", "4G", "5G", "10G",
				}, false),
			},

			"redundant": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},

			"tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"area": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func getSourceHAInfoOfRouterSingleToPortConnection(i map[string]interface{}) connections.SourceHAInfo {
	ipAddress := i["ip_address"].(string)
	asPathPrependIn := i["as_path_prepend_in"]
	asPathPrependOut := i["as_path_prepend_out"]

	if asPathPrependIn.(string) == "OFF" || asPathPrependIn.(string) == "" {
		asPathPrependIn = interface{}(nil)
	} else {
		tmp, _ := strconv.Atoi(asPathPrependIn.(string))
		asPathPrependIn = interface{}(tmp)
	}

	if asPathPrependOut.(string) == "OFF" || asPathPrependOut.(string) == "" {
		asPathPrependOut = interface{}(nil)
	} else {
		tmp, _ := strconv.Atoi(asPathPrependOut.(string))
		asPathPrependOut = interface{}(tmp)
	}

	primaryASPathPrepend := connections.ASPathPrepend{
		In:  &asPathPrependIn,
		Out: &asPathPrependOut,
	}
	haInfo := connections.SourceHAInfo{
		IPAddress:     ipAddress,
		ASPathPrepend: primaryASPathPrepend,
	}
	log.Printf("[DEBUG] ASPathPrepend: %#v", haInfo)
	return haInfo
}

func getSourceHAInfoOfRouterSingleToPortConnectionForUpdate(i map[string]interface{}) connections.SourceHAInfoForUpdate {
	asPathPrependIn := i["as_path_prepend_in"]
	asPathPrependOut := i["as_path_prepend_out"]

	if asPathPrependIn.(string) == "OFF" || asPathPrependIn.(string) == "" {
		asPathPrependIn = interface{}(nil)
	} else {
		tmp, _ := strconv.Atoi(asPathPrependIn.(string))
		asPathPrependIn = interface{}(tmp)
	}

	if asPathPrependOut.(string) == "OFF" || asPathPrependOut.(string) == "" {
		asPathPrependOut = interface{}(nil)
	} else {
		tmp, _ := strconv.Atoi(asPathPrependOut.(string))
		asPathPrependOut = interface{}(tmp)
	}

	asPathPrepend := connections.ASPathPrepend{
		In:  &asPathPrependIn,
		Out: &asPathPrependOut,
	}
	haInfo := connections.SourceHAInfoForUpdate{
		ASPathPrepend: asPathPrepend,
	}
	log.Printf("[DEBUG] ASPathPrepend: %#v", haInfo)
	return haInfo
}

func getSourceOfRouterSingleToPortConnection(d *schema.ResourceData) connections.Source {
	tmpSource := d.Get("source_information").([]interface{})

	tmpPrimary := tmpSource[0].(map[string]interface{})
	primary := getSourceHAInfoOfRouterSingleToPortConnection(tmpPrimary)
	source := connections.Source{
		RouterID:  d.Get("source_router_id").(string),
		GroupName: d.Get("source_group_name").(string),
		Primary:   primary,
	}

	// tmpSecondary := tmpSource[1].(map[string]interface{})
	// secondary := getSourceHAInfoOfRouterSingleToPortConnection(tmpSecondary)
	// source.Secondary = secondary

	routeFilter := connections.RouteFilter{
		In:  d.Get("source_route_filter_in").(string),
		Out: d.Get("source_route_filter_out").(string),
	}
	source.RouteFilter = routeFilter
	return source
}

func getSourceOfRouterSingleToPortConnectionForUpdate(d *schema.ResourceData) connections.SourceForUpdate {
	tmpSource := d.Get("source_information").([]interface{})

	tmpPrimary := tmpSource[0].(map[string]interface{})
	primary := getSourceHAInfoOfRouterSingleToPortConnectionForUpdate(tmpPrimary)
	source := connections.SourceForUpdate{
		Primary: primary,
	}

	// tmpSecondary := tmpSource[1].(map[string]interface{})
	// secondary := getSourceHAInfoOfRouterSingleToPortConnectionForUpdate(tmpSecondary)
	// source.Secondary = secondary

	routeFilter := connections.RouteFilter{
		In:  d.Get("source_route_filter_in").(string),
		Out: d.Get("source_route_filter_out").(string),
	}
	source.RouteFilter = routeFilter
	return source
}

func getDestinationOfRouterSingleToPortConnection(d *schema.ResourceData) connections.Destination {
	tmpDestination := d.Get("destination_information").([]interface{})

	tmpPrimary := tmpDestination[0]
	primary := connections.DestinationHAInfo{
		PortID:    tmpPrimary.(map[string]interface{})["port_id"].(string),
		VLAN:      tmpPrimary.(map[string]interface{})["vlan"].(int),
		IPAddress: tmpPrimary.(map[string]interface{})["ip_address"].(string),
		ASN:       tmpPrimary.(map[string]interface{})["asn"].(string),
	}
	destination := connections.Destination{
		Primary: primary,
	}

	// tmpSecondary := tmpDestination[1]
	// secondary := connections.DestinationHAInfo{
	// 	PortID:    tmpSecondary.(map[string]interface{})["port_id"].(string),
	// 	VLAN:      tmpSecondary.(map[string]interface{})["vlan"].(int),
	// 	IPAddress: tmpSecondary.(map[string]interface{})["ip_address"].(string),
	// 	ASN:       tmpSecondary.(map[string]interface{})["asn"].(string),
	// }
	// destination.Secondary = secondary

	return destination
}

func resourceEriRouterSingleToPortConnectionV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	createOpts := &connections.CreateOpts{
		Name:        d.Get("name").(string),
		Source:      getSourceOfRouterSingleToPortConnection(d),
		Destination: getDestinationOfRouterSingleToPortConnection(d),
		Bandwidth:   d.Get("bandwidth").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	r, err := connections.Create(client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI connection(router to port): %s", err)
	}

	d.SetId(r.ID)

	log.Printf("[INFO] Connection ID: %s", r.ID)

	log.Printf(
		"[DEBUG] Waiting for connection (%s) to become available", r.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Processing"},
		Target:     []string{"Completed"},
		Refresh:    RouterToPortConnectionV1StateRefreshFunc(client, r.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for connection (%s) to become ready: %s", r.ID, err)
	}

	return resourceEriRouterSingleToPortConnectionV1Read(d, meta)
}

func getSourceInformationOfRouterSingleToPortConnectionForState(r *connections.Connection) []map[string]interface{} {
	primary := map[string]interface{}{
		"ip_address":          r.Source.Primary.IPAddress,
		"as_path_prepend_in":  r.Source.Primary.ASPathPrepend.In,
		"as_path_prepend_out": r.Source.Primary.ASPathPrepend.Out,
	}
	// secondary := map[string]interface{}{
	// 	"ip_address":          r.Source.Secondary.IPAddress,
	// 	"as_path_prepend_in":  r.Source.Secondary.ASPathPrepend.In,
	// 	"as_path_prepend_out": r.Source.Secondary.ASPathPrepend.Out,
	// }
	return []map[string]interface{}{
		primary,
		// secondary,
	}
}

func getDestinationOfRouterSingleToPortConnectionInformationForState(r *connections.Connection) []map[string]interface{} {
	primary := map[string]interface{}{
		"port_id":    r.Destination.Primary.PortID,
		"vlan":       r.Destination.Primary.VLAN,
		"ip_address": r.Destination.Primary.IPAddress,
		"asn":        r.Destination.Primary.ASN,
	}
	// secondary := map[string]interface{}{
	// 	"port_id":    r.Destination.Secondary.PortID,
	// 	"vlan":       r.Destination.Secondary.VLAN,
	// 	"ip_address": r.Destination.Secondary.IPAddress,
	// 	"asn":        r.Destination.Secondary.ASN,
	// }
	return []map[string]interface{}{
		primary,
		// secondary,
	}
}

func resourceEriRouterSingleToPortConnectionV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	r, err := connections.Get(client, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "connection")
	}

	log.Printf("[DEBUG] Retrieved connection %s: %+v", d.Id(), r)

	d.Set("name", r.Name)

	d.Set("source_router_id", r.Source.RouterID)
	d.Set("source_group_name", r.Source.GroupName)

	d.Set("source_route_filter_in", r.Source.RouteFilter.In)
	d.Set("source_route_filter_out", r.Source.RouteFilter.Out)

	d.Set("source_information", getSourceInformationOfRouterSingleToPortConnectionForState(r))
	d.Set("destination_information", getDestinationOfRouterSingleToPortConnectionInformationForState(r))

	d.Set("bandwidth", r.Bandwidth)
	d.Set("redundant", r.Redundant)
	d.Set("tenant_id", r.TenantID)
	d.Set("area", r.Area)

	return nil
}

func resourceEriRouterSingleToPortConnectionV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	if d.HasChange("source_information") {
		updateOpts := connections.UpdateOpts{
			Source: getSourceOfRouterSingleToPortConnectionForUpdate(d),
		}
		_, err := connections.Update(client, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error activating FIC ERI connection: %s", err)
		}

		log.Printf(
			"[DEBUG] Waiting for connection (%s) to become complete", d.Id())

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"Processing"},
			Target:     []string{"Completed"},
			Refresh:    RouterToPortConnectionV1StateRefreshFunc(client, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		log.Printf("[DEBUG] Waiting for connection (%s) to become complete", d.Id())
		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for connection (%s) to become complete: %s", d.Id(), err)
		}
	}

	return resourceEriRouterSingleToPortConnectionV1Read(d, meta)
}

func resourceEriRouterSingleToPortConnectionV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	if err := connections.Delete(client, d.Id()).ExtractErr(); err != nil {
		return CheckDeleted(d, err, "connection")
	}

	log.Printf("[DEBUG] Waiting for connection (%s) to delete", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Processing", "Completed"},
		Target:     []string{"Deleted"},
		Refresh:    RouterToPortConnectionV1StateRefreshFunc(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for connection (%s) to delete: %s",
			d.Id(), err)
	}

	d.SetId("")
	return nil
}

// func RouterToPortConnectionV1StateRefreshFunc(client *fic.ServiceClient, connectionID string) resource.StateRefreshFunc {
// 	return func() (interface{}, string, error) {
// 		v, err := connections.Get(client, connectionID).Extract()
// 		if err != nil {
// 			if _, ok := err.(fic.ErrDefault404); ok {
// 				return v, "Deleted", nil
// 			}
// 			return nil, "", err
// 		}

// 		if v.OperationStatus == "Error" {
// 			return v, v.OperationStatus, fmt.Errorf("There was an error retrieving the connection(router to port) information.")
// 		}

// 		return v, v.OperationStatus, nil
// 	}
// }
