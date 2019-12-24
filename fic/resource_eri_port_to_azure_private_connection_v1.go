package fic

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/nttcom/go-fic"
	connections "github.com/nttcom/go-fic/fic/eri/v1/port_to_azure_private_connections"
)

func resourceEriPortToAzurePrivateConnectionV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceEriPortToAzurePrivateConnectionV1Create,
		Read:   resourceEriPortToAzurePrivateConnectionV1Read,
		Delete: resourceEriPortToAzurePrivateConnectionV1Delete,
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

			"source_primary_port_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"source_primary_vlan": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"source_secondary_port_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"source_secondary_vlan": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"source_asn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"destination_interconnect": &schema.Schema{
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

			"destination_service_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"destination_shared_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
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
					"10M", "20M", "30M", "40M", "50M",
					"100M", "200M", "300M", "400M", "500M",
					"1G", "2G", "3G", "4G", "5G",
					"10G",
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

			"operation_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"operation_status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceEriPortToAzurePrivateConnectionV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	primary := connections.Primary{
		PortID: d.Get("source_primary_port_id").(string),
		VLAN:   d.Get("source_primary_vlan").(int),
	}

	secondary := connections.Secondary{
		PortID: d.Get("source_secondary_port_id").(string),
		VLAN:   d.Get("source_secondary_vlan").(int),
	}

	source := connections.Source{
		Primary:   primary,
		Secondary: secondary,
		ASN:       d.Get("source_asn").(string),
	}

	destination := connections.Destination{
		Interconnect: d.Get("destination_interconnect").(string),
		QosType:      d.Get("destination_qos_type").(string),
		ServiceKey:   d.Get("destination_service_key").(string),
		SharedKey:    d.Get("destination_shared_key").(string),
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
		return fmt.Errorf("Error creating FIC ERI port to azure private connection: %s", err)
	}

	d.SetId(r.ID)

	log.Printf("[INFO] Connection ID: %s", r.ID)

	log.Printf(
		"[DEBUG] Waiting for port to azure private connection (%s) to become available", r.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Processing"},
		Target:     []string{"Completed"},
		Refresh:    resourcePortToAzurePrivateConnectionV1StateRefreshFunc(client, r.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for port to azure private connection (%s) to become ready: %s", r.ID, err)
	}

	return resourceEriPortToAzurePrivateConnectionV1Read(d, meta)
}

func resourceEriPortToAzurePrivateConnectionV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	r, err := connections.Get(client, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "connection")
	}

	log.Printf("[DEBUG] Retrieved port to azure private connection %s: %+v", d.Id(), r)

	d.Set("name", r.Name)

	d.Set("source_primary_port_id", r.Source.Primary.PortID)
	d.Set("source_primary_vlan", r.Source.Primary.VLAN)

	d.Set("source_secondary_port_id", r.Source.Secondary.PortID)
	d.Set("source_secondary_vlan", r.Source.Secondary.VLAN)

	d.Set("source_asn", r.Source.ASN)

	d.Set("destination_interconnect", r.Destination.Interconnect)
	d.Set("destination_qos_type", r.Destination.QosType)
	d.Set("destination_service_key", r.Destination.ServiceKey)
	d.Set("destination_shared_key", r.Destination.SharedKey)

	d.Set("primary_connected_network_address", r.PrimaryConnectedNetworkAddress)
	d.Set("secondary_connected_network_address", r.SecondaryConnectedNetworkAddress)

	d.Set("bandwidth", r.Bandwidth)
	d.Set("redundant", r.Redundant)
	d.Set("tenant_id", r.TenantID)
	d.Set("area", r.Area)

	return nil
}

func resourceEriPortToAzurePrivateConnectionV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	if err := connections.Delete(client, d.Id()).ExtractErr(); err != nil {
		return CheckDeleted(d, err, "connection")
	}

	log.Printf("[DEBUG] Waiting for port to azure private connection (%s) to delete", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Processing", "Completed"},
		Target:     []string{"Deleted"},
		Refresh:    resourcePortToAzurePrivateConnectionV1StateRefreshFunc(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for port to azure private connection (%s) to delete: %s",
			d.Id(), err)
	}

	d.SetId("")
	return nil
}

func resourcePortToAzurePrivateConnectionV1StateRefreshFunc(client *fic.ServiceClient, connectionID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := connections.Get(client, connectionID).Extract()
		if err != nil {
			if _, ok := err.(fic.ErrDefault404); ok {
				return v, "Deleted", nil
			}
			return nil, "", err
		}

		if v.OperationStatus == "Error" {
			return v, v.OperationStatus, fmt.Errorf("there was an error retrieving the port to azure private connection information")
		}

		return v, v.OperationStatus, nil
	}
}
