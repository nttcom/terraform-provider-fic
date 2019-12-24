package fic

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/nttcom/go-fic"
	connections "github.com/nttcom/go-fic/fic/eri/v1/router_to_ecl_connections"
)

func resourceEriRouterToECLConnectionV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceEriRouterToECLConnectionV1Create,
		Read:   resourceEriRouterToECLConnectionV1Read,
		Update: resourceEriRouterToECLConnectionV1Update,
		Delete: resourceEriRouterToECLConnectionV1Delete,
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

			"destination_interconnect": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				// ValidateFunc: validation.StringInSlice(
				// 	[]string{"JP3-1", "JP5-1"}, false),
			},

			"destination_qos_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"guarantee"}, false),
			},

			"destination_ecl_tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"destination_ecl_api_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"destination_ecl_api_secret_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"primary_connected_network_address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"secondary_connected_network_address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"bandwidth": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"10M", "20M", "30M", "40M", "50M", "100M",
					"200M", "300M", "400M", "500M",
					"1G",
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

func resourceEriRouterToECLConnectionV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	routeFilter := connections.RouteFilter{
		In:  d.Get("source_route_filter_in").(string),
		Out: d.Get("source_route_filter_out").(string),
	}

	source := connections.Source{
		RouterID:    d.Get("source_router_id").(string),
		GroupName:   d.Get("source_group_name").(string),
		RouteFilter: routeFilter,
	}

	destination := connections.Destination{
		Interconnect:    d.Get("destination_interconnect").(string),
		QosType:         d.Get("destination_qos_type").(string),
		ECLTenantID:     d.Get("destination_ecl_tenant_id").(string),
		ECLAPIKey:       d.Get("destination_ecl_api_key").(string),
		ECLAPISecretKey: d.Get("destination_ecl_api_secret_key").(string),
	}

	createOpts := &connections.CreateOpts{
		Name:                             d.Get("name").(string),
		Source:                           source,
		Destination:                      destination,
		Bandwidth:                        d.Get("bandwidth").(string),
		PrimaryConnectedNetworkAddress:   d.Get("primary_connected_network_address").(string),
		SecondaryConnectedNetworkAddress: d.Get("secondary_connected_network_address").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	r, err := connections.Create(client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI connection(router to ecl): %s", err)
	}

	d.SetId(r.ID)

	log.Printf("[INFO] Connection ID: %s", r.ID)

	log.Printf(
		"[DEBUG] Waiting for connection (%s) to become available", r.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Processing"},
		Target:     []string{"Completed"},
		Refresh:    RouterToECLConnectionV1StateRefreshFunc(client, r.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for connection (%s) to become ready: %s", r.ID, err)
	}

	return resourceEriRouterToECLConnectionV1Read(d, meta)
}

func resourceEriRouterToECLConnectionV1Read(d *schema.ResourceData, meta interface{}) error {
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

	d.Set("destination_inter_connect", r.Destination.Interconnect)
	d.Set("destination_qos_type", r.Destination.QosType)

	d.Set("bandwidth", r.Bandwidth)
	d.Set("redundant", r.Redundant)
	d.Set("tenant_id", r.TenantID)
	d.Set("area", r.Area)

	return nil
}

func resourceEriRouterToECLConnectionV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	if d.HasChange("source_route_filter_in") || d.HasChange("source_route_filter_out") {
		routeFilter := connections.RouteFilter{
			In:  d.Get("source_route_filter_in").(string),
			Out: d.Get("source_route_filter_out").(string),
		}
		source := connections.SourceForUpdate{
			RouteFilter: routeFilter,
		}
		updateOpts := connections.UpdateOpts{
			Source: source,
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
			Refresh:    RouterToECLConnectionV1StateRefreshFunc(client, d.Id()),
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

	return resourceEriRouterToECLConnectionV1Read(d, meta)
}

func resourceEriRouterToECLConnectionV1Delete(d *schema.ResourceData, meta interface{}) error {
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
		Refresh:    RouterToECLConnectionV1StateRefreshFunc(client, d.Id()),
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

func RouterToECLConnectionV1StateRefreshFunc(client *fic.ServiceClient, connectionID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := connections.Get(client, connectionID).Extract()
		if err != nil {
			if _, ok := err.(fic.ErrDefault404); ok {
				return v, "Deleted", nil
			}
			return nil, "", err
		}

		if v.OperationStatus == "Error" {
			return v, v.OperationStatus, fmt.Errorf("There was an error retrieving the connection(router to ecl) information.")
		}

		return v, v.OperationStatus, nil
	}
}
