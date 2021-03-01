output "fic_router_id" {
  value = fic_eri_router_v1.router_01.id
}

output "fic_router_nat_id" {
  value = fic_eri_router_v1.router_01.nat_id
}

output "fic_gateway_id" {
  value = fic_eri_router_to_ecl_connection_v1.connection_ecl_01.id
}