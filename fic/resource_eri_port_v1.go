package fic

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/nttcom/go-fic"
	"github.com/nttcom/go-fic/fic/eri/v1/ports"
)

func resourceEriPortV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceEriPortV1Create,
		Read:   resourceEriPortV1Read,
		Update: resourceEriPortV1Update,
		Delete: resourceEriPortV1Delete,
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

			"switch_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"number_of_vlans": &schema.Schema{
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"vlan_ranges"},
				ValidateFunc: IntInSlice([]int{
					16, 32, 48, 64, 80, 96, 112, 128, 144, 160, 176, 192, 208, 224, 240, 256, 272,
					288, 304, 320, 336, 352, 368, 384, 400, 416, 432, 448, 464, 480, 496, 512,
				}),
			},

			"vlan_ranges": &schema.Schema{
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"number_of_vlans"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"end": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},

			"port_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"10G", "1G"}, false),
			},

			"is_activated": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
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

			"location": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"vlans": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vid": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceEriPortV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	createOpts := &ports.CreateOpts{
		Name:       d.Get("name").(string),
		SwitchName: d.Get("switch_name").(string),
		PortType:   d.Get("port_type").(string),
	}

	if numberOfVLANs, ok := d.Get("number_of_vlans").(int); ok && numberOfVLANs != 0 {
		createOpts.NumberOfVLANs = numberOfVLANs
	} else {
		createOpts.VLANRanges = getVLANRanges(d)
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	r, err := ports.Create(client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI port: %s", err)
	}

	d.SetId(r.ID)

	log.Printf("[INFO] Port ID: %s", r.ID)

	log.Printf(
		"[DEBUG] Waiting for port (%s) to become available", r.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Processing"},
		Target:     []string{"Completed"},
		Refresh:    PortV1StateRefreshFunc(client, r.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for port (%s) to become ready: %s", r.ID, err)
	}

	// Activate port if is_activated parameter is specified in
	// CREATE: phase of Terraform
	isActivated := d.Get("is_activated").(bool)
	if isActivated {
		r, err = ports.Activate(client, r.ID).Extract()
		if err != nil {
			return fmt.Errorf("Error activating FIC ERI port: %s", err)
		}

		log.Printf(
			"[DEBUG] Waiting for port (%s) to become active", r.ID)

		activateStateConf := &resource.StateChangeConf{
			Pending:    []string{"Processing"},
			Target:     []string{"Completed"},
			Refresh:    PortV1StateRefreshFunc(client, r.ID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		log.Printf("[DEBUG] Waiting for port (%s) to become active", r.ID)
		_, err = activateStateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for port (%s) to become active: %s", r.ID, err)
		}
	}

	return resourceEriPortV1Read(d, meta)
}

func getVLANsForState(r *ports.Port) []map[string]interface{} {
	var result []map[string]interface{}
	for _, v := range r.VLANs {
		vid := v.VID
		status := v.Status
		m := map[string]interface{}{
			"vid":    vid,
			"status": status,
		}
		result = append(result, m)
	}
	return result
}

func resourceEriPortV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	r, err := ports.Get(client, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "port")
	}

	log.Printf("[DEBUG] Retrieved port %s: %+v", d.Id(), r)

	d.Set("name", r.Name)
	d.Set("switch_name", r.SwitchName)
	d.Set("vlan_ranges", r.VLANRanges)
	d.Set("is_activated", r.IsActivated)
	d.Set("tenant_id", r.TenantID)
	d.Set("area", r.Area)
	d.Set("location", r.Location)
	d.Set("vlans", getVLANsForState(r))

	return nil
}

func resourceEriPortV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	if d.HasChange("is_activated") {
		_, err := ports.Activate(client, d.Id()).Extract()
		if err != nil {
			return fmt.Errorf("Error activating FIC ERI port: %s", err)
		}

		log.Printf(
			"[DEBUG] Waiting for port (%s) to become active", d.Id())

		activateStateConf := &resource.StateChangeConf{
			Pending:    []string{"Processing"},
			Target:     []string{"Completed"},
			Refresh:    PortV1StateRefreshFunc(client, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		log.Printf("[DEBUG] Waiting for port (%s) to become active", d.Id())
		_, err = activateStateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for port (%s) to become active: %s", d.Id(), err)
		}
	}

	return resourceEriPortV1Read(d, meta)
}

func resourceEriPortV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	if err := ports.Delete(client, d.Id()).ExtractErr(); err != nil {
		return CheckDeleted(d, err, "port")
	}

	log.Printf("[DEBUG] Waiting for port (%s) to delete", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Processing", "Completed"},
		Target:     []string{"Deleted"},
		Refresh:    PortV1StateRefreshFunc(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for port (%s) to delete: %s",
			d.Id(), err)
	}

	d.SetId("")
	return nil
}

func getVLANRanges(d *schema.ResourceData) []string {
	var result []string
	rawRanges := d.Get("vlan_ranges").([]interface{})
	for _, r := range rawRanges {
		start := r.(map[string]interface{})["start"].(int)
		end := r.(map[string]interface{})["end"].(int)
		v := fmt.Sprintf("%d-%d", start, end)
		result = append(result, v)
	}
	return result
}

func PortV1StateRefreshFunc(client *fic.ServiceClient, portID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := ports.Get(client, portID).Extract()
		if err != nil {
			if _, ok := err.(fic.ErrDefault404); ok {
				return v, "Deleted", nil
			}
			return nil, "", err
		}

		if v.OperationStatus == "Error" {
			return v, v.OperationStatus, fmt.Errorf("There was an error retrieving the port information.")
		}

		return v, v.OperationStatus, nil
	}
}
