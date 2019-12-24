package fic

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/nttcom/go-fic"
	"github.com/nttcom/go-fic/fic/eri/v1/routers"
)

func resourceEriRouterV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceEriRouterV1Create,
		Read:   resourceEriRouterV1Read,
		Delete: resourceEriRouterV1Delete,
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

			"area": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"user_ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"redundant": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"firewalls": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"is_activated": &schema.Schema{
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},

			"nats": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_activated": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"routing_groups": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"firewall_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"nat_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceEriRouterV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	redundant := d.Get("redundant").(bool)
	createOpts := &routers.CreateOpts{
		Name:          d.Get("name").(string),
		Area:          d.Get("area").(string),
		UserIPAddress: d.Get("user_ip_address").(string),
		Redundant:     &redundant,
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	r, err := routers.Create(client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI router: %s", err)
	}

	d.SetId(r.ID)

	log.Printf("[INFO] Router ID: %s", r.ID)

	log.Printf(
		"[DEBUG] Waiting for router (%s) to become available", r.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Processing"},
		Target:     []string{"Completed"},
		Refresh:    RouterV1StateRefreshFunc(client, r.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for router (%s) to become ready: %s", r.ID, err)
	}

	return resourceEriRouterV1Read(d, meta)
}

func getRouterFirewallForState(r *routers.Router) []map[string]interface{} {
	var result []map[string]interface{}
	for _, v := range r.Firewalls {
		m := map[string]interface{}{
			"id":           v.ID,
			"is_activated": v.IsActivated,
		}
		result = append(result, m)
	}
	return result
}

func getRouterNATForState(r *routers.Router) []map[string]interface{} {
	var result []map[string]interface{}
	for _, v := range r.NATs {
		m := map[string]interface{}{
			"id":           v.ID,
			"is_activated": v.IsActivated,
		}
		result = append(result, m)
	}
	return result
}

func getRoutingGroupForState(r *routers.Router) []map[string]interface{} {
	var result []map[string]interface{}
	for _, v := range r.RoutingGroups {
		m := map[string]interface{}{
			"name": v.Name,
		}
		result = append(result, m)
	}
	return result
}

func resourceEriRouterV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	r, err := routers.Get(client, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "router")
	}

	log.Printf("[DEBUG] Retrieved router %s: %+v", d.Id(), r)

	d.Set("name", r.Name)
	d.Set("area", r.Area)
	d.Set("user_ip_address", r.UserIPAddress)
	d.Set("redundant", r.Redundant)
	d.Set("tenant_id", r.TenantID)
	d.Set("firewalls", getRouterFirewallForState(r))
	d.Set("nats", getRouterNATForState(r))
	d.Set("routing_groups", getRoutingGroupForState(r))

	firewallID := r.Firewalls[0].ID
	natID := r.NATs[0].ID

	d.Set("firewall_id", firewallID)
	d.Set("nat_id", natID)

	return nil
}

func resourceEriRouterV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	if err := routers.Delete(client, d.Id()).ExtractErr(); err != nil {
		return CheckDeleted(d, err, "router")
	}

	log.Printf("[DEBUG] Waiting for router (%s) to delete", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Processing", "Completed"},
		Target:     []string{"Deleted"},
		Refresh:    RouterV1StateRefreshFunc(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for router (%s) to delete: %s",
			d.Id(), err)
	}

	d.SetId("")
	return nil
}

func RouterV1StateRefreshFunc(client *fic.ServiceClient, portID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := routers.Get(client, portID).Extract()
		if err != nil {
			if _, ok := err.(fic.ErrDefault404); ok {
				return v, "Deleted", nil
			}
			return nil, "", err
		}

		if v.OperationStatus == "Error" {
			return v, v.OperationStatus, fmt.Errorf("There was an error retrieving the router information.")
		}

		return v, v.OperationStatus, nil
	}
}
