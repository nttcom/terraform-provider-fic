package fic

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/nttcom/go-fic"
	"github.com/nttcom/go-fic/fic/eri/v1/routers/components/nat_global_ip_address_sets"
)

func resourceEriNATGlobalIPAddressSetV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceEriNATGlobalIPAddressSetV1Create,
		Read:   resourceEriNATGlobalIPAddressSetV1Read,
		Delete: resourceEriNATGlobalIPAddressSetV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			"router_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"nat_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"sourceNapt", "destinationNat",
				}, false),
			},

			"number_of_addresses": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"addresses": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceEriNATGlobalIPAddressSetV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	routerID := d.Get("router_id").(string)
	natID := d.Get("nat_id").(string)

	createOpts := &nat_global_ip_address_sets.CreateOpts{
		Name:              d.Get("name").(string),
		Type:              d.Get("type").(string),
		NumberOfAddresses: d.Get("number_of_addresses").(int),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	r, err := nat_global_ip_address_sets.Create(client, routerID, natID, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error activating FIC ERI global ip address set: %s", err)
	}

	globalIPAddressSetID := r.ID
	id := fmt.Sprintf("%s/%s/%s", routerID, natID, globalIPAddressSetID)
	d.SetId(id)

	log.Printf("[INFO] Global IP Address Set ID: %s", r.ID)

	log.Printf(
		"[DEBUG] Waiting for global ip address set (%s) to become available", r.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Processing"},
		Target:     []string{"Completed"},
		Refresh:    NATGlobalIPAddressSetV1StateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for global ip address set (%s) to become ready: %s", r.ID, err)
	}
	return resourceEriNATGlobalIPAddressSetV1Read(d, meta)
}

func resourceEriNATGlobalIPAddressSetV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	id := d.Id()
	routerID := strings.Split(id, "/")[0]
	natID := strings.Split(id, "/")[1]
	globalIPAddressSetID := strings.Split(id, "/")[2]
	r, err := nat_global_ip_address_sets.Get(
		client, routerID, natID, globalIPAddressSetID).Extract()
	if err != nil {
		return CheckDeleted(d, err, "nat_global_ip_address_set")
	}

	log.Printf("[DEBUG] Retrieved global ip address set %s: %+v", d.Id(), r)

	d.Set("router_id", d.Get("router_id").(string))
	d.Set("nat_id", d.Get("nat_id").(string))

	d.Set("name", r.Name)
	d.Set("type", r.Type)
	d.Set("number_of_addresses", r.NumberOfAddresses)
	d.Set("number_of_addresses", r.NumberOfAddresses)
	d.Set("addresses", r.Addresses)
	d.Set("operation_status", r.OperationStatus)

	return nil
}

func resourceEriNATGlobalIPAddressSetV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	id := d.Id()
	routerID := strings.Split(id, "/")[0]
	natID := strings.Split(id, "/")[1]
	globalIPAddressSetID := strings.Split(id, "/")[2]
	_, err = nat_global_ip_address_sets.Delete(
		client, routerID, natID, globalIPAddressSetID).Extract()

	log.Printf("[DEBUG] Waiting for global ip address set (%s) to delete", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Processing"},
		Target:     []string{"Deleted"},
		Refresh:    NATGlobalIPAddressSetV1StateRefreshFunc(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for global ip address set (%s) to delete: %s",
			d.Id(), err)
	}

	d.SetId("")
	return nil
}

func NATGlobalIPAddressSetV1StateRefreshFunc(client *fic.ServiceClient, id string) resource.StateRefreshFunc {
	routerID := strings.Split(id, "/")[0]
	natID := strings.Split(id, "/")[1]
	globalIPAddressSetID := strings.Split(id, "/")[2]

	return func() (interface{}, string, error) {
		v, err := nat_global_ip_address_sets.Get(client, routerID, natID, globalIPAddressSetID).Extract()
		if err != nil {
			if _, ok := err.(fic.ErrDefault404); ok {
				return v, "Deleted", nil
			}
			return nil, "", err
		}

		if v.OperationStatus == "Error" {
			return v, v.OperationStatus, fmt.Errorf("There was an error retrieving the global ip address set information.")
		}

		return v, v.OperationStatus, nil
	}
}
