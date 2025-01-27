terraform {
  required_providers {
    kind = {
      source  = "tehcyx/kind"
      version = "0.4.0"
    }
  }
}
resource "kind_cluster" "default" {
  name            = var.cluster_name
  node_image      = "kindest/node:v1.27.1"
  kubeconfig_path = pathexpand("/tmp/kube.config")
  wait_for_ready  = true

  kind_config {
    kind        = "Cluster"
    api_version = "kind.x-k8s.io/v1alpha4"

    node {
      role = "control-plane"
      extra_port_mappings {
        container_port = 8080
        host_port      = 8080
      }
    }

    node {
      role = "worker"
    }

    node {
      role = "worker"
    }
  }
}
