package fic

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/nttcom/go-fic"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	connections "github.com/nttcom/go-fic/fic/eri/v1/router_paired_to_gcp_connections"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourcePairedRouterToGCPConnection() *schema.Resource {
	validInterconnects := []string{
		"Equinix-TY2-1", "Equinix-TY2-2", "Equinix-TY2-3", "Equinix-TY2-4", "@Tokyo-CC2-1", "@Tokyo-CC2-2", "@Tokyo-CC2-3", "@Tokyo-CC2-4", "Equinix-OS1-1", "Equinix-OS1-2", "Equinix-OS1-3", "NTT-Dojima2-1", "NTT-Dojima2-2", "NTT-Dojima2-3",
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
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-fA-F\d]{8}(-[a-fA-F\d]{4}){3}-[a-fA-F\d]{12}/[a-zA-Z\d-]*/[1,2]$`), "see https://cloud.google.com/network-connectivity/docs/interconnect/concepts/terminology?_ga=2.264742223.-1966628098.1560150466#pairingkey"),
			},
		},
	}

	return &schema.Resource{
		Create: resourcePairedRouterToGCPConnectionCreate,
		Read:   resourcePairedRouterToGCPConnectionRead,
		Update: resourcePairedRouterToGCPConnectionUpdate,
		Delete: resourcePairedRouterToGCPConnectionDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[\w&()-]{1,64}$`), "must be less than 64 characters in half-width alphanumeric characters and some symbols &()-_"),
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
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile(`^F\d{12}$`), "must be a F + 12-digit number"),
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
							Type:     schema.TypeInt,
							Computed: true,
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

func resourcePairedRouterToGCPConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating FIC client: %w", err)
	}

	opts := &connections.CreateOpts{
		Name:        d.Get("name").(string),
		Source:      expandSource(d.Get("source").([]interface{})),
		Destination: expandDestination(d.Get("destination").([]interface{})),
		Bandwidth:   d.Get("bandwidth").(string),
	}

	conn, err := connections.Create(client, opts).Extract()
	if err != nil {
		return fmt.Errorf("error creating FIC paired router to GCP connection: %w", err)
	}

	d.SetId(conn.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Processing"},
		Target:     []string{"Completed"},
		Refresh:    routerToGCPConnectionRefresh(client, conn.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for connection (%s) to become ready: %w", conn.ID, err)
	}

	d.Set("operation_id", conn.OperationID)

	return resourcePairedRouterToGCPConnectionRead(d, meta)
}

func routerToGCPConnectionRefresh(c *fic.ServiceClient, id string) func() (interface{}, string, error) {
	return func() (interface{}, string, error) {
		conn, err := connections.Get(c, id).Extract()
		if err != nil {
			var e fic.ErrDefault404
			if errors.As(err, &e) {
				return nil, "", nil
			}
			return nil, "", err
		}

		return conn, conn.OperationStatus, nil
	}
}

func resourcePairedRouterToGCPConnectionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating FIC client: %w", err)
	}

	conn, err := connections.Get(client, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "error getting FIC paired router to GCP connection")
	}

	d.Set("name", conn.Name)
	d.Set("bandwidth", conn.Bandwidth)
	d.Set("source", flattenSource(conn.Source))
	d.Set("destination", flattenDestination(conn.Destination))
	d.Set("redundant", conn.Redundant)
	d.Set("tenant_id", conn.TenantID)
	d.Set("area", conn.Area)
	d.Set("operation_status", conn.OperationStatus)
	d.Set("primary_connected_network_address", conn.PrimaryConnectedNetworkAddress)
	d.Set("secondary_connected_network_address", conn.SecondaryConnectedNetworkAddress)

	return nil
}

func resourcePairedRouterToGCPConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating FIC client: %w", err)
	}

	source := expandSource(d.Get("source").([]interface{}))
	opts := connections.UpdateOpts{
		Source: connections.SourceForUpdate{
			RouteFilter: source.RouteFilter,
			Primary:     source.Primary,
			Secondary:   source.Secondary,
		},
		Bandwidth: d.Get("bandwidth").(string),
	}

	conn, err := connections.Update(client, d.Id(), opts).Extract()
	if err != nil {
		return fmt.Errorf("error updating FIC paired router to GCP connection: %w", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Processing"},
		Target:     []string{"Completed"},
		Refresh:    routerToGCPConnectionRefresh(client, conn.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for connection (%s) to become ready: %w", conn.ID, err)
	}

	d.Set("operation_id", conn.OperationID)

	return resourcePairedRouterToGCPConnectionRead(d, meta)
}

func resourcePairedRouterToGCPConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating FIC client: %w", err)
	}

	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		if err := connections.Delete(client, d.Id()).ExtractErr(); err != nil {
			var e404 fic.ErrDefault404
			if errors.As(err, &e404) {
				return nil
			}

			var e409 fic.ErrDefault409
			if errors.As(err, &e409) {
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error deleting FIC paired router to GCP connection: %w", err)
	}

	d.SetId("")

	return nil
}
