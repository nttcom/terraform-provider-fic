provider "fic" {
  auth_url          = var.fic_auth_url
  user_name         = var.fic_api_key
  password          = var.fic_api_secret
  tenant_id         = var.fic_tenant_id
  user_domain_id    = "default"
  project_domain_id = "default"
}


module "fic-router" {
  source = "./fic-router"

  fic_router_name    = "fic_router_01"
  fic_router_area    = "JPEAST"
  fic_router_user_ip = "10.0.0.0/27"
  fic_nat_user_ip = [
    "10.0.1.0/30",
    "10.0.1.4/30",
    "10.0.1.8/30",
    "10.0.1.12/30",
    "10.0.1.16/30",
    "10.0.1.20/30",
    "10.0.1.24/30",
    "10.0.1.28/30"
  ]
  fic_nat_gip_set_name = "for_Wasabi_GIP"
  fic_nat_rule_from = [
    "group_1"
  ]
  fic_nat_rule_to = "group_2"
  fic_nat_entry   = ["for_Wasabi_GIP"]

  fic_to_ecl_name                  = "connection_ecl_01"
  fic_to_ecl_source_group_name     = "group_1"
  fic_to_ecl_source_bgp_filter_in  = "noRoute"
  fic_to_ecl_source_bgp_filter_out = "noRoute"
  fic_to_ecl_destination           = "JP4-1"
  fic_to_ecl_bandwidth             = "10M"
  fic_to_ecl_primary_address       = "10.0.3.0/30"
  fic_to_ecl_secondary_address     = "10.0.3.4/30"

  ecl_key    = var.ecl_api_key
  ecl_secret = var.ecl_api_secret
  ecl_tenant = var.ecl_tenant_id
}
