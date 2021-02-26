variable "fic_router_name" {
    type = string
}
variable "fic_router_area" {
    type = string
}
variable "fic_router_user_ip" {
    type = string
}
variable "fic_nat_user_ip" {
    type = list(string)
}
variable "fic_nat_gip_set_name" {
    type = string
}
variable "fic_nat_rule_from" {
    type = list(string)
}
variable "fic_nat_rule_to" {
    type = string
}
variable "fic_nat_entry" {
    type = list(string)
}

variable "fic_to_ecl_name" {
    type = string
}
variable "fic_to_ecl_source_group_name" {
    type = string
}
variable "fic_to_ecl_source_bgp_filter_in" {
    type = string
}
variable "fic_to_ecl_source_bgp_filter_out" {
    type = string
}
variable "fic_to_ecl_destination" {
    type = string
}
variable "fic_to_ecl_bandwidth" {
    type = string
}
variable "fic_to_ecl_primary_address" {
    type = string
}
variable "fic_to_ecl_secondary_address" {
    type = string
}

variable "ecl_key" {}
variable "ecl_secret" {}
variable "ecl_tenant" {}
