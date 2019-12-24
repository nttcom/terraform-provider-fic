package fic

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/nttcom/go-fic"
	"github.com/nttcom/go-fic/fic/eri/v1/routers/components/nats"
)

func resourceEriNATComponentV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceEriNATComponentV1Activate,
		Read:   resourceEriNATComponentV1Read,
		Update: resourceEriNATComponentV1Update,
		Delete: resourceEriNATComponentV1Deactivate,
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

			"user_ip_addresses": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"global_ip_address_sets": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"sourceNapt", "destinationNat",
							}, false),
						},
						"number_of_addresses": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},

			"source_napt_rules": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									"group_1", "group_2", "group_3", "group_4",
								}, false),
							},
						},
						"to": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"group_1", "group_2", "group_3", "group_4",
							}, false),
						},
						"entries": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"then": &schema.Schema{
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 8,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},

			"destination_nat_rules": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"group_1", "group_2", "group_3", "group_4",
							}, false),
						},
						"to": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"group_1", "group_2", "group_3", "group_4",
							}, false),
							// Elem: &schema.Schema{
							// 	Type: schema.TypeString,
							// },
						},
						"entries": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_destination_address": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"then": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

			"redundant": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},

			"is_activated": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceEriNATComponentV1Activate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	routerID := d.Get("router_id").(string)
	natID := d.Get("nat_id").(string)

	activateOpts := &nats.ActivateOpts{
		UserIPAddresses:     getUserIPAddresses(d),
		GlobalIPAddressSets: getGlobalIPAddressSets(d),
	}

	log.Printf("[DEBUG] Activate Options: %#v", activateOpts)
	r, err := nats.Activate(client, routerID, natID, activateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error activating FIC ERI nat component: %s", err)
	}

	id := fmt.Sprintf("%s/%s", routerID, natID)
	d.SetId(id)

	log.Printf("[INFO] NAT Component ID: %s", r.ID)

	log.Printf(
		"[DEBUG] Waiting for nat component (%s) to become available", r.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Processing"},
		Target:     []string{"Completed"},
		Refresh:    NATComponentV1StateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for nat component (%s) to become ready: %s", r.ID, err)
	}

	sourceNAPTRules := d.Get("source_napt_rules").([]interface{})
	destinationNATRules := d.Get("destination_nat_rules").([]interface{})

	log.Printf("[DEBUG] Source NAPT Rule is set as: %#v", sourceNAPTRules)
	log.Printf("[DEBUG] Destination NAT Rule is set as: %#v", destinationNATRules)

	if len(sourceNAPTRules) > 0 || len(destinationNATRules) > 0 {
		updateSourceNAPTORDestinationNAT(d, meta)
	}

	return resourceEriNATComponentV1Read(d, meta)
}

func resourceEriNATComponentV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	id := d.Id()
	routerID := strings.Split(id, "/")[0]
	natID := strings.Split(id, "/")[1]
	r, err := nats.Get(client, routerID, natID).Extract()
	if err != nil {
		return CheckDeleted(d, err, "nat")
	}

	log.Printf("[DEBUG] Retrieved nat component %s: %+v", d.Id(), r)

	d.Set("router_id", d.Get("router_id").(string))
	d.Set("user_ip_address", r.UserIPAddresses)
	d.Set("source_napt_rules", getSourceNAPTRuleForState(r))
	d.Set("destination_nat_rules", getDestinationNATRuleForState(r))
	d.Set("redundant", r.Redundant)
	d.Set("is_activated", r.IsActivated)

	return nil
}

func resourceEriNATComponentV1Update(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] d.HasChange('source_napt_rules'): %#v", d.HasChange("source_napt_rules"))
	log.Printf("[DEBUG] d.HasChange('destination_nat_rules'): %#v", d.HasChange("source_napt_rules"))
	if d.HasChange("source_napt_rules") || d.HasChange("destination_nat_rules") {
		log.Printf("[DEBUG] Either Source NAPT or Destination NAT is going to update...")
		updateSourceNAPTORDestinationNAT(d, meta)
	}
	return resourceEriNATComponentV1Read(d, meta)
}

func resourceEriNATComponentV1Deactivate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	id := d.Id()
	routerID := strings.Split(id, "/")[0]
	natID := strings.Split(id, "/")[1]
	_, err = nats.Deactivate(client, routerID, natID).Extract()

	log.Printf("[DEBUG] Waiting for nat component (%s) to delete", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Processing"},
		Target:     []string{"Completed"},
		Refresh:    NATComponentV1StateRefreshFunc(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for nat component (%s) to delete: %s",
			d.Id(), err)
	}

	d.SetId("")
	return nil
}

func updateSourceNAPTORDestinationNAT(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	id := d.Id()
	routerID := strings.Split(id, "/")[0]
	natID := strings.Split(id, "/")[1]

	if d.HasChange("source_napt_rules") || d.HasChange("destination_nat_rules") {

		updateOpts := nats.UpdateOpts{
			SourceNAPTRules:     getSourceNAPTRules(d),
			DestinationNATRules: getDestinationNATRules(d),
		}
		_, err := nats.Update(client, routerID, natID, updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating FIC ERI nat component: %s", err)
		}

		log.Printf(
			"[DEBUG] Waiting for nat component (%s) to become complete", d.Id())

		activateStateConf := &resource.StateChangeConf{
			Pending:    []string{"Processing"},
			Target:     []string{"Completed"},
			Refresh:    NATComponentV1StateRefreshFunc(client, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		log.Printf("[DEBUG] Waiting for nat component (%s) to become complete", d.Id())
		_, err = activateStateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for nat component (%s) to become complete: %s", d.Id(), err)
		}
	}
	return nil
}

func NATComponentV1StateRefreshFunc(client *fic.ServiceClient, id string) resource.StateRefreshFunc {
	routerID := strings.Split(id, "/")[0]
	natID := strings.Split(id, "/")[1]

	return func() (interface{}, string, error) {
		v, err := nats.Get(client, routerID, natID).Extract()
		if err != nil {
			if _, ok := err.(fic.ErrDefault404); ok {
				return v, "Deleted", nil
			}
			return nil, "", err
		}

		if v.OperationStatus == "Error" {
			return v, v.OperationStatus, fmt.Errorf("There was an error retrieving the nat component information.")
		}

		return v, v.OperationStatus, nil
	}
}

func getSourceNAPTRuleForState(r *nats.NAT) []map[string]interface{} {
	var result []map[string]interface{}

	for _, n := range r.SourceNAPTRules {
		var entries []map[string]interface{}
		for _, e := range n.Entries {
			entry := map[string]interface{}{
				"then": e.Then,
			}
			entries = append(entries, entry)
		}

		tmp := map[string]interface{}{
			"from":    n.From,
			"to":      n.To,
			"entries": entries,
		}
		result = append(result, tmp)
	}

	log.Printf("[DEBUG] State of source NAPT rule is set like: %s", result)
	return result
}

func getDestinationNATRuleForState(r *nats.NAT) []map[string]interface{} {
	var result []map[string]interface{}

	for _, n := range r.DestinationNATRules {
		var entries []map[string]string
		for _, e := range n.Entries {

			entry := map[string]string{
				"match_destination_address": e.Match.DestinationAddress,
				"then":                      e.Then,
			}
			entries = append(entries, entry)
		}

		tmp := map[string]interface{}{
			"from":    n.From,
			"to":      n.To,
			"entries": entries,
		}
		result = append(result, tmp)
	}

	log.Printf("[DEBUG] State of destination NAT rule is set like: %s", result)
	return result
}

func getGlobalIPAddressSets(d *schema.ResourceData) []nats.GlobalIPAddressSet {
	var result []nats.GlobalIPAddressSet

	rawAddresses := d.Get("global_ip_address_sets").([]interface{})
	for _, r := range rawAddresses {
		name := r.(map[string]interface{})["name"].(string)
		addressType := r.(map[string]interface{})["type"].(string)
		numberOfAddresses := r.(map[string]interface{})["number_of_addresses"].(int)
		temp := nats.GlobalIPAddressSet{
			Name:              name,
			Type:              addressType,
			NumberOfAddresses: numberOfAddresses,
		}
		result = append(result, temp)
	}
	return result
}

func getUserIPAddresses(d *schema.ResourceData) []string {
	var result []string
	rawUserIPAddresses := d.Get("user_ip_addresses").([]interface{})
	for _, r := range rawUserIPAddresses {
		ip := r.(string)
		result = append(result, ip)
	}
	return result
}

func getSourceNAPTRules(d *schema.ResourceData) []nats.SourceNAPTRule {
	var result []nats.SourceNAPTRule
	rawSourceNAPTRules := d.Get("source_napt_rules").([]interface{})
	for _, r := range rawSourceNAPTRules {

		var from []string
		tmpFrom := r.(map[string]interface{})["from"].([]interface{})
		for _, f := range tmpFrom {
			from = append(from, f.(string))
		}

		to := r.(map[string]interface{})["to"].(string)

		var entries []nats.EntryInSourceNAPTRule
		tmpEntries := r.(map[string]interface{})["entries"].([]interface{})

		for _, e := range tmpEntries {
			var then []string
			tmpThen := e.(map[string]interface{})["then"].([]interface{})
			for _, t := range tmpThen {
				then = append(then, t.(string))
			}

			entry := nats.EntryInSourceNAPTRule{
				Then: then,
			}

			entries = append(entries, entry)
		}

		tmp := nats.SourceNAPTRule{
			From:    from,
			To:      to,
			Entries: entries,
		}

		result = append(result, tmp)
	}

	return result
}

func getDestinationNATRules(d *schema.ResourceData) []nats.DestinationNATRule {
	var result []nats.DestinationNATRule
	rawDestinatonNATRules := d.Get("destination_nat_rules").([]interface{})
	for _, r := range rawDestinatonNATRules {

		from := r.(map[string]interface{})["from"].(string)
		to := r.(map[string]interface{})["to"].(string)

		var entries []nats.EntryInDestinationNATRule
		tmpEntries := r.(map[string]interface{})["entries"].([]interface{})

		for _, e := range tmpEntries {
			then := e.(map[string]interface{})["then"].(string)
			matchDestinationAddress := e.(map[string]interface{})["match_destination_address"].(string)

			entry := nats.EntryInDestinationNATRule{
				Then: then,
				Match: nats.Match{
					DestinationAddress: matchDestinationAddress,
				},
			}

			entries = append(entries, entry)
		}

		tmp := nats.DestinationNATRule{
			From:    from,
			To:      to,
			Entries: entries,
		}

		result = append(result, tmp)
	}

	return result
}
