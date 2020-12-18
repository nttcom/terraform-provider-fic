
resource "fic_eri_router_v1" "router_01" {
    name            = var.fic_router_name
    area            = var.fic_router_area
    user_ip_address = var.fic_router_user_ip
    redundant       = true
}

resource "fic_eri_nat_component_v1" "nat_1" {
  router_id = fic_eri_router_v1.router_01.id
  nat_id = fic_eri_router_v1.router_01.nat_id

  user_ip_addresses = var.fic_nat_user_ip

	    global_ip_address_sets  {
    name = var.fic_nat_gip_set_name
    type = "sourceNapt"
    number_of_addresses = 1
  }

        source_napt_rules {
    from = var.fic_nat_rule_from
    to   = var.fic_nat_rule_to

    entries {
        then = var.fic_nat_entry
        }
    }

}

resource "fic_eri_router_to_ecl_connection_v1" "connection_ecl_01" {
    name    = var.fic_to_ecl_name

    source_router_id    = fic_eri_router_v1.router_01.id

    source_group_name   = var.fic_to_ecl_source_group_name
    source_route_filter_in  = var.fic_to_ecl_source_bgp_filter_in
    source_route_filter_out = var.fic_to_ecl_source_bgp_filter_out

    destination_interconnect    = var.fic_to_ecl_destination
    destination_qos_type    = "guarantee"
    destination_ecl_tenant_id   = var.ecl_tenant
    destination_ecl_api_key = var.ecl_key
    destination_ecl_api_secret_key  = var.ecl_secret

    bandwidth   = var.fic_to_ecl_bandwidth

    primary_connected_network_address = var.fic_to_ecl_primary_adress
	secondary_connected_network_address = var.fic_to_ecl_secondary_adress

    depends_on = [fic_eri_nat_component_v1.nat_1]
}
