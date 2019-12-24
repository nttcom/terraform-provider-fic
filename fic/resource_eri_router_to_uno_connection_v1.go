package fic

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/nttcom/go-fic"
	connections "github.com/nttcom/go-fic/fic/eri/v1/router_to_uno_connections"
)

func resourceEriRouterToUNOConnectionV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceEriRouterToUNOConnectionV1Create,
		Read:   resourceEriRouterToUNOConnectionV1Read,
		Update: resourceEriRouterToUNOConnectionV1Update,
		Delete: resourceEriRouterToUNOConnectionV1Delete,
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
				ValidateFunc: validation.StringInSlice(
					[]string{"fullRoute", "noRoute"}, false),
			},

			"source_route_filter_out": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"fullRoute", "fullRouteWithDefaultRoute", "noRoute"}, false),
			},

			"destination_interconnect": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"Interconnect-Tokyo-1", "Interconnect-Osaka-1"}, false),
			},

			"destination_c_number": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"destination_parent_contract_number": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"destination_vpn_number": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"destination_qos_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"guarantee"}, false),
			},

			"destination_route_filter_out": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"fullRoute", "fullRouteWithDefaultRoute",
						"defaultRoute", "privateRoute"}, false),
			},

			"connected_network_address": &schema.Schema{
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

			"destination_contract_number": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
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

func resourceEriRouterToUNOConnectionV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	sourceRouteFilter := connections.SourceRouteFilter{
		In:  d.Get("source_route_filter_in").(string),
		Out: d.Get("source_route_filter_out").(string),
	}

	source := connections.Source{
		RouterID:    d.Get("source_router_id").(string),
		GroupName:   d.Get("source_group_name").(string),
		RouteFilter: sourceRouteFilter,
	}

	destinationRouteFilter := connections.DestinationRouteFilter{
		Out: d.Get("destination_route_filter_out").(string),
	}

	destination := connections.Destination{
		Interconnect:         d.Get("destination_interconnect").(string),
		CNumber:              d.Get("destination_c_number").(string),
		ParentContractNumber: d.Get("destination_parent_contract_number").(string),
		VPNNumber:            d.Get("destination_vpn_number").(string),
		QosType:              d.Get("destination_qos_type").(string),
		RouteFilter:          destinationRouteFilter,
	}

	createOpts := &connections.CreateOpts{
		Name:                    d.Get("name").(string),
		Source:                  source,
		Destination:             destination,
		Bandwidth:               d.Get("bandwidth").(string),
		ConnectedNetworkAddress: d.Get("connected_network_address").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	r, err := connections.Create(client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI connection(router to uno): %s", err)
	}

	d.SetId(r.ID)

	log.Printf("[INFO] Connection ID: %s", r.ID)

	log.Printf(
		"[DEBUG] Waiting for connection (%s) to become available", r.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Processing"},
		Target:     []string{"Completed"},
		Refresh:    RouterToUNOConnectionV1StateRefreshFunc(client, r.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for connection (%s) to become ready: %s", r.ID, err)
	}

	// Even CNumber parameter is required, response of UNO connection does not have
	// this value, so this point is only change to store c number into state.
	log.Printf("[DEBUG] Preserve destination_c_number as: %#v", d.Get("destination_c_number").(string))
	d.Set("destination_c_number", d.Get("destination_c_number").(string))
	return resourceEriRouterToUNOConnectionV1Read(d, meta)
}

func resourceEriRouterToUNOConnectionV1Read(d *schema.ResourceData, meta interface{}) error {
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

	d.Set("destination_interconnect", r.Destination.Interconnect)

	// Response does not have c number so set this value by using stored value.
	d.Set("destination_c_number", d.Get("destination_c_number"))
	log.Printf("[DEBUG] Re-set destination_c_number as: %#v", d.Get("destination_c_number"))

	d.Set("destination_contract_number", r.Destination.ContractNumber)
	d.Set("destination_parent_contract_number", r.Destination.ParentContractNumber)
	d.Set("destination_vpn_number", r.Destination.VPNNumber)
	d.Set("destination_route_filter_out", r.Destination.RouteFilter.Out)

	d.Set("connected_network_address", r.ConnectedNetworkAddress)

	d.Set("bandwidth", r.Bandwidth)
	d.Set("redundant", r.Redundant)
	d.Set("tenant_id", r.TenantID)

	return nil
}

func resourceEriRouterToUNOConnectionV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	if d.HasChange("source_route_filter_in") || d.HasChange("source_route_filter_out") || d.HasChange("destination_route_filter_out") {

		source := connections.SourceForUpdate{
			RouteFilter: connections.SourceRouteFilter{
				In:  d.Get("source_route_filter_in").(string),
				Out: d.Get("source_route_filter_out").(string),
			},
		}
		destination := connections.DestinationForUpdate{
			RouteFilter: connections.DestinationRouteFilter{
				Out: d.Get("destination_route_filter_out").(string),
			},
		}
		updateOpts := connections.UpdateOpts{
			Source:      source,
			Destination: destination,
		}

		log.Printf("[DEBUG] Update Options: %#v", updateOpts)
		_, err := connections.Update(client, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error activating FIC ERI connection: %s", err)
		}

		log.Printf(
			"[DEBUG] Waiting for connection (%s) to become complete", d.Id())

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"Processing"},
			Target:     []string{"Completed"},
			Refresh:    RouterToUNOConnectionV1StateRefreshFunc(client, d.Id()),
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

	return resourceEriRouterToUNOConnectionV1Read(d, meta)
}

func resourceEriRouterToUNOConnectionV1Delete(d *schema.ResourceData, meta interface{}) error {
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
		Refresh:    RouterToUNOConnectionV1StateRefreshFunc(client, d.Id()),
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

func RouterToUNOConnectionV1StateRefreshFunc(client *fic.ServiceClient, connectionID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := connections.Get(client, connectionID).Extract()
		if err != nil {
			if _, ok := err.(fic.ErrDefault404); ok {
				return v, "Deleted", nil
			}
			return nil, "", err
		}

		if v.OperationStatus == "Error" {
			return v, v.OperationStatus, fmt.Errorf("There was an error retrieving the connection(router to uno) information.")
		}

		return v, v.OperationStatus, nil
	}
}
