package fic

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/httpclient"

	"github.com/nttcom/go-fic"
	"github.com/nttcom/go-fic/fic/utils"

	"github.com/nttcom/terraform-provider-fic/fic/clientconfig"

	"log"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/pathorcontents"
)

type Config struct {
	CACertFile        string
	ClientCertFile    string
	ClientKeyFile     string
	Cloud             string
	DefaultDomain     string
	DomainID          string
	DomainName        string
	EndpointType      string
	ForceSSSEndpoint  string
	IdentityEndpoint  string
	Insecure          *bool
	Password          string
	ProjectDomainName string
	ProjectDomainID   string
	Region            string
	TenantID          string
	TenantName        string
	Token             string
	UserDomainName    string
	UserDomainID      string
	Username          string
	UserID            string
	terraformVersion  string

	OsClient *fic.ProviderClient
}

func (c *Config) LoadAndValidate() error {
	// Make sure at least one of auth_url or cloud was specified.
	if c.IdentityEndpoint == "" && c.Cloud == "" {
		return fmt.Errorf("One of 'auth_url' or 'cloud' must be specified")
	}

	validEndpoint := false
	validEndpoints := []string{
		"internal", "internalURL",
		"admin", "adminURL",
		"public", "publicURL",
		"",
	}

	for _, endpoint := range validEndpoints {
		if c.EndpointType == endpoint {
			validEndpoint = true
		}
	}

	if !validEndpoint {
		return fmt.Errorf("Invalid endpoint type provided")
	}

	clientOpts := new(clientconfig.ClientOpts)

	// If a cloud entry was given, base AuthOptions on a clouds.yaml file.
	if c.Cloud != "" {
		clientOpts.Cloud = c.Cloud

		cloud, err := clientconfig.GetCloudFromYAML(clientOpts)
		if err != nil {
			return err
		}

		if c.Region == "" && cloud.RegionName != "" {
			c.Region = cloud.RegionName
		}

		if c.CACertFile == "" && cloud.CACertFile != "" {
			c.CACertFile = cloud.CACertFile
		}

		if c.ClientCertFile == "" && cloud.ClientCertFile != "" {
			c.ClientCertFile = cloud.ClientCertFile
		}

		if c.ClientKeyFile == "" && cloud.ClientKeyFile != "" {
			c.ClientKeyFile = cloud.ClientKeyFile
		}

		if c.Insecure == nil && cloud.Verify != nil {
			v := (!*cloud.Verify)
			c.Insecure = &v
		}
	} else {
		authInfo := &clientconfig.AuthInfo{
			AuthURL:           c.IdentityEndpoint,
			DefaultDomain:     c.DefaultDomain,
			DomainID:          c.DomainID,
			DomainName:        c.DomainName,
			Password:          c.Password,
			ProjectDomainID:   c.ProjectDomainID,
			ProjectDomainName: c.ProjectDomainName,
			ProjectID:         c.TenantID,
			ProjectName:       c.TenantName,
			Token:             c.Token,
			UserDomainID:      c.UserDomainID,
			UserDomainName:    c.UserDomainName,
			Username:          c.Username,
			UserID:            c.UserID,
		}
		clientOpts.AuthInfo = authInfo
	}

	ao, err := clientconfig.AuthOptions(clientOpts)
	if err != nil {
		return err
	}

	client, err := utils.NewClient(ao.IdentityEndpoint)
	if err != nil {
		return err
	}

	// Set UserAgent
	client.UserAgent.Prepend(httpclient.TerraformUserAgent(c.terraformVersion))

	config := &tls.Config{}
	if c.CACertFile != "" {
		caCert, _, err := pathorcontents.Read(c.CACertFile)
		if err != nil {
			return fmt.Errorf("Error reading CA Cert: %s", err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(caCert))
		config.RootCAs = caCertPool
	}

	if c.Insecure == nil {
		config.InsecureSkipVerify = false
	} else {
		config.InsecureSkipVerify = *c.Insecure
	}

	if c.ClientCertFile != "" && c.ClientKeyFile != "" {
		clientCert, _, err := pathorcontents.Read(c.ClientCertFile)
		if err != nil {
			return fmt.Errorf("Error reading Client Cert: %s", err)
		}
		clientKey, _, err := pathorcontents.Read(c.ClientKeyFile)
		if err != nil {
			return fmt.Errorf("Error reading Client Key: %s", err)
		}

		cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
		if err != nil {
			return err
		}

		config.Certificates = []tls.Certificate{cert}
		config.BuildNameToCertificate()
	}

	// if OS_DEBUG is set, log the requests and responses
	var osDebug bool
	if os.Getenv("OS_DEBUG") != "" {
		osDebug = true
	}

	transport := &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: config}
	client.HTTPClient = http.Client{
		Transport: &LogRoundTripper{
			Rt:      transport,
			OsDebug: osDebug,
		},
	}

	err = utils.Authenticate(client, *ao)
	if err != nil {
		return err
	}
	c.OsClient = client

	return nil
}

func (c *Config) determineRegion(region string) string {
	// If a resource-level region was not specified, and a provider-level region was set,
	// use the provider-level region.
	if region == "" && c.Region != "" {
		region = c.Region
	}

	log.Printf("[DEBUG] FIC Region is: %s", region)
	return region
}

func (c *Config) getEndpointType() fic.Availability {
	return fic.AvailabilityPublic
}

func (c *Config) eriV1Client(region string) (*fic.ServiceClient, error) {
	return utils.NewEriV1(c.OsClient, fic.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getEndpointType(),
	})
}
