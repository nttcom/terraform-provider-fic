---
layout: "fic"
page_title: "Provider: Flexible InterConnect"
sidebar_current: "docs-fic-index"
description: |-
  The Flexible InterConnect(FIC) provider is used to interact with the many resources supported by Flexible InterConnect. The provider needs to be configured with the proper credentials before it can be used.
---

# Flexible InterConnect Provider

The Flexible InterConnect provider is used to interact with the
many resources supported by Flexible InterConnect.
The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the Flexible InterConnect Provider
terraform {
  required_providers {
    fic = {
      source  = "nttcom/fic"
      version = "x.y.z"
    }
  }
}

provider "fic" {
  user_name         = "my-api-key"
  password          = "my-api-secret-key"
  tenant_name       = "my-tenant-id"
  auth_url          = "https://keystone-myregion-fic.api.ntt.com/v3/"
  user_domain_id    = "default"
  project_domain_id = "default"
}

# Create a port
resource "fic_eri_port_v1" "test-port" {
  # ...
}
```

## Configuration Reference

The following arguments are supported:

* `auth_url` - (Optional; required if `cloud` is not specified) The Identity
  authentication URL. If omitted, the `OS_AUTH_URL` environment variable is used.

* `cloud` - (Optional; required if `auth_url` is not specified) An entry in a
  `clouds.yaml` file. See the FIC CLI `os-client-config`
  [documentation](https://ecl.ntt.com/en/documents/tutorials/eclc/rsts/installation.html)
  for more information about `clouds.yaml` files. If omitted, the `OS_CLOUD`
  environment variable is used.

* `region` - (Optional) The region of the Flexible InterConnect to use. If omitted,
  the `OS_REGION_NAME` environment variable is used. If `OS_REGION_NAME` is
  not set, then no region will be used. It should be possible to omit the
  region in single-region Flexible InterConnect environments, but this behavior
  may vary depending on the Flexible InterConnect environment being used.

* `user_name` - (Optional) The Username to login with. If omitted, the
  `OS_USERNAME` environment variable is used.

* `user_id` - (Optional) The User ID to login with. If omitted, the
  `OS_USER_ID` environment variable is used.

* `tenant_id` - (Optional) The ID of the Tenant (Identity v2) or Project
  (Identity v3) to login with. If omitted, the `OS_TENANT_ID` or
  `OS_PROJECT_ID` environment variables are used.

* `tenant_name` - (Optional) The Name of the Tenant to login with.
  If omitted, the `OS_TENANT_NAME` or `OS_PROJECT_NAME` environment
  variable are used.

* `password` - (Optional) The Password to login with. If omitted, the
  `OS_PASSWORD` environment variable is used.

* `token` - (Optional; Required if not using `user_name` and `password`)
  A token is an expiring, temporary means of access issued via the Keystone
  service. By specifying a token, you do not have to specify a username/password
  combination, since the token was already created by a username/password out of
  band of Terraform. If omitted, the `OS_TOKEN` or `OS_AUTH_TOKEN` environment
  variables are used.

* `user_domain_name` - (Optional) The domain name where the user is located. If
  omitted, the `OS_USER_DOMAIN_NAME` environment variable is checked.

* `user_domain_id` - (Optional) The domain ID where the user is located. If
  omitted, the `OS_USER_DOMAIN_ID` environment variable is checked.

* `project_domain_name` - (Optional) The domain name where the project is
  located. If omitted, the `OS_PROJECT_DOMAIN_NAME` environment variable is
  checked.

* `project_domain_id` - (Optional) The domain ID where the project is located
  If omitted, the `OS_PROJECT_DOMAIN_ID` environment variable is checked.

* `domain_id` - (Optional) The ID of the Domain to scope to (Identity v3). If
  omitted, the `OS_DOMAIN_ID` environment variable is checked.

* `domain_name` - (Optional) The Name of the Domain to scope to (Identity v3).
  If omitted, the following environment variables are checked (in this order):
  `OS_DOMAIN_NAME`.

* `default_domain` - (Optional) The ID of the Domain to scope to if no other
  domain is specified (Identity v3). If omitted, the environment variable
  `OS_DEFAULT_DOMAIN` is checked or a default value of "default" will be
  used.

* `insecure` - (Optional) Trust self-signed SSL certificates. If omitted, the
  `OS_INSECURE` environment variable is used.

* `cacert_file` - (Optional) Specify a custom CA certificate when communicating
  over SSL. You can specify either a path to the file or the contents of the
  certificate. If omitted, the `OS_CACERT` environment variable is used.

* `cert` - (Optional) Specify client certificate file for SSL client
  authentication. You can specify either a path to the file or the contents of
  the certificate. If omitted the `OS_CERT` environment variable is used.

* `key` - (Optional) Specify client private key file for SSL client
  authentication. You can specify either a path to the file or the contents of
  the key. If omitted the `OS_KEY` environment variable is used.

* `endpoint_type` - (Optional) Specify which type of endpoint to use from the
  service catalog. It can be set using the OS_ENDPOINT_TYPE environment
  variable. If not set, public endpoints is used.

## Additional Logging

This provider has the ability to log all HTTP requests and responses between
Terraform and the Flexible InterConnect which is useful for troubleshooting and
debugging.

To enable these logs, set the `OS_DEBUG` environment variable to `1` along
with the usual `TF_LOG=DEBUG` environment variable:

```shell
$ OS_DEBUG=1 TF_LOG=DEBUG terraform apply
```

If you submit these logs with a bug report, please ensure any sensitive
information has been scrubbed first!

## Testing and Development

Thank you for your interest in further developing the Flexible InterConnect provider!
Here are a few notes which should help you get started. If you have any questions or
feel these notes need further details, please open an Issue and let us know.

### Coding and Style

This provider aims to adhere to the coding style and practices used in the
other major Terraform Providers. However, this is not a strict rule.

We're very mindful that not everyone is a full-time developer (most of the
Flexible InterConnect Provider contributors aren't!) and we're happy to provide 
guidance. Don't be afraid if this is your first contribution to the
Flexible InterConnect provider or even your first contribution to an
open source project!

### Testing Environment

In order to start fixing bugs or adding features, you need access to an
Flexible InterConnect environment.
You can use a production Flexible InterConnect which you have access to.

### go-fic

This provider uses [go-fic](https://github.com/nttcom/go-fic)
as the Go Flexible InterConnect SDK.
All API interaction between this provider and an Flexible InterConnect is done exclusively with go-fic.

### Adding a Feature

If you'd like to add a new feature to this provider, it must first be supported in go-fic. If go-fic is missing the feature, then it'll first have to be added there before you can start working on the feature in Terraform.
Fortunately, most of the regular Flexible InterConnect Provider contributors
also work on go-fic, so we can try to get the feature added quickly.

If the feature is already included in go-fic, then you can begin work
directly in the Flexible InterConnect provider.

If you have any questions about whether go-fic currently supports a
certain feature, please feel free to ask and we can verify for you.

### Fixing a Bug

Similarly, if you find a bug in this provider, the bug might actually be a
go-fic bug. If this is the case, then we'll need to get the bug fixed in
go-fic first.

However, if the bug is with Terraform itself, then you can begin work directly
in the Flexible InterConnect provider.

Again, if you have any questions about whether the bug you're trying to fix is a go-fic but, please ask.

### Acceptance Tests

Acceptance Tests are a crucial part of adding features or fixing a bug. Please
make sure to read the core [testing](https://www.terraform.io/docs/extend/testing/index.html)
documentation for more information about how Acceptance Tests work.

In order to run the Acceptance Tests, you'll need to set the following
environment variables:

* `OS_SWITCH_NAME` - a Name of switch you create a port.

We recommend only running the acceptance tests related to the feature or bug
you're working on. To do this, run:

```shell
$ cd path/to/terraform-provider-fic
$ make testacc TEST=./fic TESTARGS="-run=<keyword> -count=1"
```

Where `<keyword>` is the full name or partial name of a test. For example:

```shell
$ make testacc TEST=./fic TESTARGS="-run=TestAccEriPortV1Basic -count=1"
```

We recommend running tests with logging set to `DEBUG`:

```shell
$ TF_LOG=DEBUG make testacc TEST=./fic TESTARGS="-run=TestAccEriPortV1Basic -count=1"
```

And you can even enable Flexible InterConnect debugging to see the actual HTTP API requests:

```shell
$ TF_LOG=DEBUG OS_DEBUG=1 make testacc TEST=./fic TESTARGS="-run=TestAccEriPortV1Basic -count=1"
```

### Creating a Pull Request

When you're ready to submit a Pull Request, create a branch, commit your code,
and push to your forked version of `terraform-provider-fic`:

```shell
$ git remote add my-github-username https://github.com/my-github-username/terraform-provider-fic
$ git checkout -b my-feature
$ git add .
$ git commit
$ git push -u my-github-username my-feature
```

Then navigate to https://github.com/nttcom/terraform-provider-fic
and create a Pull Request.
