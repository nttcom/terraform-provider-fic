## Create Router and ECL-connection

This example provides the FIC router and the connection settings to the ECL.

The default configuration is to do source NAT from "group1" to "group2". This will be used for future connections to Wasabi Object Storage and has no significance at this time.
ECL connections will be made to "group1".
Be very careful with the values you set for IP addresses, etc., as they can affect your existing environment.

After setting the variables according to your environment, run `terraform plan`, and if everything is ok, run `terraform apply`.
