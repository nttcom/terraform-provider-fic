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
	"github.com/nttcom/go-fic/fic/eri/v1/routers/components/firewalls"
)

func resourceEriFirewallComponentV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceEriFirewallComponentV1Activate,
		Read:   resourceEriFirewallComponentV1Read,
		Update: resourceEriFirewallComponentV1Update,
		Delete: resourceEriFirewallComponentV1Deactivate,
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

			"firewall_id": &schema.Schema{
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

			"rules": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"group_1", "group_2", "group_3", "group_4",
							}, false),
						},
						"to": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"group_1", "group_2", "group_3", "group_4",
							}, false),
						},
						"entries": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"match_source_address_sets": &schema.Schema{
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 10,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"match_destination_address_sets": &schema.Schema{
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 10,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"match_application": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"action": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},

			"custom_applications": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"protocol": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"tcp", "udp",
							}, false),
						},
						"destination_port": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"application_sets": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"applications": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 10,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

			"routing_group_settings": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"group_1", "group_2", "group_3", "group_4",
							}, false),
						},
						"address_sets": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 5,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"addresses": &schema.Schema{
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 10,
										Elem:     &schema.Schema{Type: schema.TypeString},
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

func resourceEriFirewallComponentV1Activate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	routerID := d.Get("router_id").(string)
	firewallID := d.Get("firewall_id").(string)

	activateOpts := &firewalls.ActivateOpts{
		UserIPAddresses: getUserIPAddresses(d),
	}

	log.Printf("[DEBUG] Activate Options: %#v", activateOpts)
	r, err := firewalls.Activate(client, routerID, firewallID, activateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error activating FIC ERI firewall component: %s", err)
	}

	id := fmt.Sprintf("%s/%s", routerID, firewallID)
	d.SetId(id)

	log.Printf("[INFO] firewall component ID: %s", r.ID)

	log.Printf(
		"[DEBUG] Waiting for firewall component (%s) to become available", r.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Processing"},
		Target:     []string{"Completed"},
		Refresh:    FirewallComponentV1StateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for firewall component (%s) to become ready: %s", r.ID, err)
	}

	rules := d.Get("rules").([]interface{})
	customApplications := d.Get("custom_applications").([]interface{})
	applicationSets := d.Get("application_sets").([]interface{})
	routingGroupSettings := d.Get("routing_group_settings").([]interface{})

	log.Printf("[DEBUG] Rules are set as: %#v", rules)
	log.Printf("[DEBUG] Custom Applications are set as: %#v", customApplications)
	log.Printf("[DEBUG] Application Sets are set as: %#v", applicationSets)
	log.Printf("[DEBUG] RoutingGroupSettings  are set as: %#v", routingGroupSettings)

	if len(rules) > 0 || len(customApplications) > 0 || len(applicationSets) > 0 || len(routingGroupSettings) > 0 {
		err := updateFirewall(d, meta)
		if err != nil {
			return fmt.Errorf("Error updating firewall component: %s", err)
		}
	}

	return resourceEriFirewallComponentV1Read(d, meta)
}

func resourceEriFirewallComponentV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	id := d.Id()
	routerID := strings.Split(id, "/")[0]
	firewallID := strings.Split(id, "/")[1]
	r, err := firewalls.Get(client, routerID, firewallID).Extract()
	if err != nil {
		return CheckDeleted(d, err, "firewall")
	}

	log.Printf("[DEBUG] Retrieved firewall component %s: %+v", d.Id(), r)

	d.Set("router_id", d.Get("router_id").(string))
	d.Set("firewall_id", d.Get("firewall_id").(string))
	d.Set("user_ip_addresses", r.UserIPAddresses)
	d.Set("redundant", r.Redundant)
	d.Set("is_activated", r.IsActivated)

	d.Set("rules", getRulesForState(r))
	d.Set("custom_applications", getCustomApplicationsForState(r))
	d.Set("application_sets", getApplicationSetsForState(r))
	d.Set("routing_group_settings", getRoutingGroupSettingsForState(r))

	return nil
}

func resourceEriFirewallComponentV1Update(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("rules") || d.HasChange("custom_applications") ||
		d.HasChange("application_sets") || d.HasChange("routing_group_settings") {

		log.Printf("[DEBUG] Firewall is going to update...")
		err := updateFirewall(d, meta)
		if err != nil {
			return fmt.Errorf("Error updating firewall component: %s", err)
		}
	}
	return resourceEriFirewallComponentV1Read(d, meta)
}

func updateFirewall(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	id := d.Id()
	routerID := strings.Split(id, "/")[0]
	firewallID := strings.Split(id, "/")[1]

	if d.HasChange("rules") || d.HasChange("custom_applications") ||
		d.HasChange("application_sets") || d.HasChange("routing_group_settings") {

		updateOpts := firewalls.UpdateOpts{
			Rules:                getRules(d),
			CustomApplications:   getCustomApplications(d),
			ApplicationSets:      getApplicationSets(d),
			RoutingGroupSettings: getRoutingGroupSettings(d),
		}
		_, err := firewalls.Update(client, routerID, firewallID, updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating FIC ERI nat component: %s", err)
		}

		log.Printf(
			"[DEBUG] Waiting for firewall component (%s) to become complete", d.Id())

		activateStateConf := &resource.StateChangeConf{
			Pending:    []string{"Processing"},
			Target:     []string{"Completed"},
			Refresh:    FirewallComponentV1StateRefreshFunc(client, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		log.Printf("[DEBUG] Waiting for firewall component (%s) to become complete", d.Id())
		_, err = activateStateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for firewall component (%s) to become complete: %s", d.Id(), err)
		}
	}
	return nil
}

func resourceEriFirewallComponentV1Deactivate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.eriV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FIC ERI client: %s", err)
	}

	id := d.Id()
	routerID := strings.Split(id, "/")[0]
	natID := strings.Split(id, "/")[1]
	_, err = firewalls.Deactivate(client, routerID, natID).Extract()

	log.Printf("[DEBUG] Waiting for firewall component (%s) to delete", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Processing"},
		Target:     []string{"Completed"},
		Refresh:    FirewallComponentV1StateRefreshFunc(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for firewall component (%s) to delete: %s",
			d.Id(), err)
	}

	d.SetId("")
	return nil
}

func FirewallComponentV1StateRefreshFunc(client *fic.ServiceClient, id string) resource.StateRefreshFunc {
	routerID := strings.Split(id, "/")[0]
	firewallID := strings.Split(id, "/")[1]

	return func() (interface{}, string, error) {
		v, err := firewalls.Get(client, routerID, firewallID).Extract()
		if err != nil {
			if _, ok := err.(fic.ErrDefault404); ok {
				return v, "Deleted", nil
			}
			return nil, "", err
		}

		if v.OperationStatus == "Error" {
			return v, v.OperationStatus, fmt.Errorf("There was an error retrieving the firewall component information.")
		}

		return v, v.OperationStatus, nil
	}
}

func getRules(d *schema.ResourceData) []firewalls.Rule {
	result := make([]firewalls.Rule, 0)

	rawRules := d.Get("rules").([]interface{})
	for _, r := range rawRules {
		from := r.(map[string]interface{})["from"].(string)
		to := r.(map[string]interface{})["to"].(string)

		var entries []firewalls.Entry
		tmpEntries := r.(map[string]interface{})["entries"].([]interface{})
		for _, e := range tmpEntries {
			tmpEntry := e.(map[string]interface{})
			name := tmpEntry["name"].(string)
			action := tmpEntry["action"].(string)

			var sourceAddressSets []string
			var destinationAddressSets []string

			// tmpMatch := e.(map[string]interface{})["match"].(map[string]interface{})
			tmpSourceAddressSets := tmpEntry["match_source_address_sets"].([]interface{})
			tmpDestinationAddressSets := tmpEntry["match_destination_address_sets"].([]interface{})
			application := tmpEntry["match_application"].(string)
			for _, s := range tmpSourceAddressSets {
				sourceAddressSets = append(sourceAddressSets, s.(string))
			}
			for _, s := range tmpDestinationAddressSets {
				destinationAddressSets = append(destinationAddressSets, s.(string))
			}
			match := firewalls.Match{
				SourceAddressSets:      sourceAddressSets,
				DestinationAddressSets: destinationAddressSets,
				Application:            application,
			}

			entry := firewalls.Entry{
				Name:   name,
				Match:  match,
				Action: action,
			}
			entries = append(entries, entry)
		}
		thisRule := firewalls.Rule{
			From:    from,
			To:      to,
			Entries: entries,
		}
		result = append(result, thisRule)
	}

	return result
}

func getCustomApplications(d *schema.ResourceData) []firewalls.CustomApplication {
	result := make([]firewalls.CustomApplication, 0)
	rawCustomApplications := d.Get("custom_applications").([]interface{})
	for _, r := range rawCustomApplications {
		name := r.(map[string]interface{})["name"].(string)
		protocol := r.(map[string]interface{})["protocol"].(string)
		destinationPort := r.(map[string]interface{})["destination_port"].(string)

		c := firewalls.CustomApplication{
			Name:            name,
			Protocol:        protocol,
			DestinationPort: destinationPort,
		}
		result = append(result, c)
	}
	return result
}

func getApplicationSets(d *schema.ResourceData) []firewalls.ApplicationSet {
	result := make([]firewalls.ApplicationSet, 0)

	rawApplicationSets := d.Get("application_sets").([]interface{})
	for _, r := range rawApplicationSets {
		name := r.(map[string]interface{})["name"].(string)
		tmpApplications := r.(map[string]interface{})["applications"].([]interface{})

		var applications []string
		for _, a := range tmpApplications {
			applications = append(applications, a.(string))
		}

		thisApplication := firewalls.ApplicationSet{
			Name:         name,
			Applications: applications,
		}
		result = append(result, thisApplication)
	}
	return result
}

func getRoutingGroupSettings(d *schema.ResourceData) []firewalls.RoutingGroupSetting {
	result := make([]firewalls.RoutingGroupSetting, 0)
	rawRoutingGroupSettings := d.Get("routing_group_settings").([]interface{})
	for _, r := range rawRoutingGroupSettings {
		groupName := r.(map[string]interface{})["group_name"].(string)

		tmpAddressSets := r.(map[string]interface{})["address_sets"].([]interface{})

		var addressSets []firewalls.AddressSet
		for _, a := range tmpAddressSets {
			name := a.(map[string]interface{})["name"].(string)
			tmpAddresses := a.(map[string]interface{})["addresses"].([]interface{})

			var addresses []string
			for _, as := range tmpAddresses {
				addresses = append(addresses, as.(string))
			}

			addressSet := firewalls.AddressSet{
				Name:      name,
				Addresses: addresses,
			}
			addressSets = append(addressSets, addressSet)
		}
		rg := firewalls.RoutingGroupSetting{
			GroupName:   groupName,
			AddressSets: addressSets,
		}
		result = append(result, rg)
	}

	return result
}

func getRoutingGroupSettingsForState(f *firewalls.Firewall) []map[string]interface{} {
	var result []map[string]interface{}
	for _, rg := range f.RoutingGroupSettings {

		var addressSets []map[string]interface{}

		for _, as := range rg.AddressSets {
			tmpAddressSet := map[string]interface{}{
				"name":      as.Name,
				"addresses": as.Addresses,
			}
			addressSets = append(addressSets, tmpAddressSet)
		}
		tmp := map[string]interface{}{
			"group_name":   rg.GroupName,
			"address_sets": addressSets,
		}
		result = append(result, tmp)
	}
	return result
}

func getApplicationSetsForState(f *firewalls.Firewall) []map[string]interface{} {
	var result []map[string]interface{}
	for _, a := range f.ApplicationSets {
		tmp := map[string]interface{}{
			"name":         a.Name,
			"applications": a.Applications,
		}
		result = append(result, tmp)
	}
	return result
}

func getRulesForState(f *firewalls.Firewall) []map[string]interface{} {
	var result []map[string]interface{}
	for _, r := range f.Rules {

		var entries []map[string]interface{}

		for _, e := range r.Entries {
			tmpEntry := map[string]interface{}{
				"name":                           e.Name,
				"match_source_address_sets":      e.Match.SourceAddressSets,
				"match_destination_address_sets": e.Match.DestinationAddressSets,
				"match_application":              e.Match.Application,
				"action":                         e.Action,
			}
			entries = append(entries, tmpEntry)
		}
		rule := map[string]interface{}{
			"from":    r.From,
			"to":      r.To,
			"entries": entries,
		}
		result = append(result, rule)
	}

	return result
}

func getCustomApplicationsForState(f *firewalls.Firewall) []map[string]string {
	var result []map[string]string

	for _, c := range f.CustomApplications {
		tmp := map[string]string{
			"name":             c.Name,
			"protocol":         c.Protocol,
			"destination_port": c.DestinationPort,
		}
		result = append(result, tmp)
	}
	return result
}
