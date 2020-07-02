package fic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/mutexkv"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// This is a global MutexKV for use within this plugin.
var osMutexKV = mutexkv.NewMutexKV()

// Provider returns a schema.Provider for Flexible InterConnect.
func Provider() terraform.ResourceProvider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"auth_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_AUTH_URL", ""),
				Description: descriptions["auth_url"],
			},

			"region": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["region"],
				DefaultFunc: schema.EnvDefaultFunc("OS_REGION_NAME", ""),
			},

			"user_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USERNAME", ""),
				Description: descriptions["user_name"],
			},

			"user_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USER_ID", ""),
				Description: descriptions["user_name"],
			},

			"tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_TENANT_ID",
					"OS_PROJECT_ID",
				}, ""),
				Description: descriptions["tenant_id"],
			},

			"tenant_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_TENANT_NAME",
					"OS_PROJECT_NAME",
				}, ""),
				Description: descriptions["tenant_name"],
			},

			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("OS_PASSWORD", ""),
				Description: descriptions["password"],
			},

			"token": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_TOKEN",
					"OS_AUTH_TOKEN",
				}, ""),
				Description: descriptions["token"],
			},

			"user_domain_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USER_DOMAIN_NAME", ""),
				Description: descriptions["user_domain_name"],
			},

			"user_domain_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USER_DOMAIN_ID", ""),
				Description: descriptions["user_domain_id"],
			},

			"project_domain_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_PROJECT_DOMAIN_NAME", ""),
				Description: descriptions["project_domain_name"],
			},

			"project_domain_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_PROJECT_DOMAIN_ID", ""),
				Description: descriptions["project_domain_id"],
			},

			"domain_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_DOMAIN_ID", ""),
				Description: descriptions["domain_id"],
			},

			"domain_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_DOMAIN_NAME", ""),
				Description: descriptions["domain_name"],
			},

			"default_domain": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_DEFAULT_DOMAIN", "default"),
				Description: descriptions["default_domain"],
			},

			"insecure": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_INSECURE", nil),
				Description: descriptions["insecure"],
			},

			"endpoint_type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_ENDPOINT_TYPE", ""),
			},

			"cacert_file": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CACERT", ""),
				Description: descriptions["cacert_file"],
			},

			"cert": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CERT", ""),
				Description: descriptions["cert"],
			},

			"key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_KEY", ""),
				Description: descriptions["key"],
			},

			"cloud": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CLOUD", ""),
				Description: descriptions["cloud"],
			},

			"force_sss_endpoint": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_FORCE_SSS_ENDPOINT", ""),
				Description: descriptions["force_sss_endpoint"],
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"fic_eri_switch_v1": dataSourceEriSwitchV1(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"fic_eri_firewall_component_v1":                   resourceEriFirewallComponentV1(),
			"fic_eri_nat_component_v1":                        resourceEriNATComponentV1(),
			"fic_eri_nat_global_ip_address_set_v1":            resourceEriNATGlobalIPAddressSetV1(),
			"fic_eri_port_to_azure_microsoft_connection_v1":   resourceEriPortToAzureMicrosoftConnectionV1(),
			"fic_eri_port_to_azure_private_connection_v1":     resourceEriPortToAzurePrivateConnectionV1(),
			"fic_eri_port_to_port_connection_v1":              resourceEriPortToPortConnectionV1(),
			"fic_eri_port_v1":                                 resourceEriPortV1(),
			"fic_eri_router_to_ecl_connection_v1":             resourceEriRouterToECLConnectionV1(),
			"fic_eri_router_paired_to_gcp_connection_v1":      resourcePairedRouterToGCPConnection(),
			"fic_eri_router_paired_to_port_connection_v1":     resourceEriRouterPairedToPortConnectionV1(),
			"fic_eri_router_single_to_port_connection_v1":     resourceEriRouterSingleToPortConnectionV1(),
			"fic_eri_router_to_azure_microsoft_connection_v1": resourceEriRouterToAzureMicrosoftConnectionV1(),
			"fic_eri_router_to_azure_private_connection_v1":   resourceEriRouterToAzurePrivateConnectionV1(),
			"fic_eri_router_to_uno_connection_v1":             resourceEriRouterToUNOConnectionV1(),
			"fic_eri_router_v1":                               resourceEriRouterV1(),
		},
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return configureProvider(d, terraformVersion)
	}

	return provider
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"auth_url": "The Identity authentication URL.",

		"region": "The Flexible InterConnect region to connect to.",

		"user_name": "Username to login with.",

		"user_id": "User ID to login with.",

		"tenant_id": "The ID of the Tenant (Identity v2) or Project (Identity v3)\n" +
			"to login with.",

		"tenant_name": "The name of the Tenant (Identity v2) or Project (Identity v3)\n" +
			"to login with.",

		"password": "Password to login with.",

		"token": "Authentication token to use as an alternative to username/password.",

		"user_domain_name": "The name of the domain where the user resides (Identity v3).",

		"user_domain_id": "The ID of the domain where the user resides (Identity v3).",

		"project_domain_name": "The name of the domain where the project resides (Identity v3).",

		"project_domain_id": "The ID of the domain where the project resides (Identity v3).",

		"domain_id": "The ID of the Domain to scope to (Identity v3).",

		"domain_name": "The name of the Domain to scope to (Identity v3).",

		"default_domain": "The name of the Domain ID to scope to if no other domain is specified. Defaults to `default` (Identity v3).",

		"insecure": "Trust self-signed certificates.",

		"cacert_file": "A Custom CA certificate.",

		"endpoint_type": "The catalog endpoint type to use.",

		"cert": "A client certificate to authenticate with.",

		"key": "A client private key to authenticate with.",

		"cloud": "An entry in a `clouds.yaml` file to use.",
	}
}

func configureProvider(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	config := Config{
		CACertFile:        d.Get("cacert_file").(string),
		ClientCertFile:    d.Get("cert").(string),
		ClientKeyFile:     d.Get("key").(string),
		Cloud:             d.Get("cloud").(string),
		DefaultDomain:     d.Get("default_domain").(string),
		DomainID:          d.Get("domain_id").(string),
		DomainName:        d.Get("domain_name").(string),
		EndpointType:      d.Get("endpoint_type").(string),
		IdentityEndpoint:  d.Get("auth_url").(string),
		Password:          d.Get("password").(string),
		ProjectDomainID:   d.Get("project_domain_id").(string),
		ProjectDomainName: d.Get("project_domain_name").(string),
		Region:            d.Get("region").(string),
		Token:             d.Get("token").(string),
		TenantID:          d.Get("tenant_id").(string),
		TenantName:        d.Get("tenant_name").(string),
		UserDomainID:      d.Get("user_domain_id").(string),
		UserDomainName:    d.Get("user_domain_name").(string),
		Username:          d.Get("user_name").(string),
		UserID:            d.Get("user_id").(string),
		terraformVersion:  terraformVersion,
	}

	v, ok := d.GetOkExists("insecure")
	if ok {
		insecure := v.(bool)
		config.Insecure = &insecure
	}

	if err := config.LoadAndValidate(); err != nil {
		return nil, err
	}

	return &config, nil
}
