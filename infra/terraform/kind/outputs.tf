output "host" {
  value = kind_cluster.default.endpoint
}

output "client_certificate" {
  value = base64encode(kind_cluster.default.client_certificate)
}

output "client_key" {
  value = base64encode(kind_cluster.default.client_key)
}

output "cluster_ca_certificate" {
  value = base64encode(kind_cluster.default.cluster_ca_certificate)
}
