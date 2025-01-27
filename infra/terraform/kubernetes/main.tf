terraform {
  required_providers {
    kubernetes = {
      source = "hashicorp/kubernetes"
    }
  }
}

provider "kubernetes" {
  host = var.k8s_host

  client_certificate     = base64decode(var.k8s_client_certificate)
  client_key             = base64decode(var.k8s_client_key)
  cluster_ca_certificate = base64decode(var.k8s_cluster_ca_certificate)
}

# --- Ingestor Service Deployment ---
resource "kubernetes_config_map" "ingestor_config_map" {
  provider = kubernetes
  metadata {
    name      = "ingestor-config"
    namespace = "default"
  }
  data = {
    app_name : "binance-trading-chart-service/ingestor"
    app_grpc_port : var.app_grpc_port
    app_env : var.app_env
    app_debug : var.app_debug
    binance_websocket_base_url : var.binance_ws_base_url
    binance_symbols : var.binance_symbols
  }
}

resource "kubernetes_manifest" "ingestor_deployment" {
  provider = kubernetes
  manifest = yamldecode(file("./../../k8s/ingestor/ingestor-deployment.yaml"))
}

resource "kubernetes_manifest" "ingestor_service" {
  provider   = kubernetes
  manifest   = yamldecode(file("./../../k8s/ingestor/ingestor-service.yaml"))
  depends_on = [kubernetes_manifest.ingestor_deployment]
}

resource "kubernetes_manifest" "ingestor_allow_egress_internet" {
  provider   = kubernetes
  manifest   = yamldecode(file("./../../k8s/ingestor/ingestor-egress-internet-access.yaml"))
  depends_on = [kubernetes_manifest.ingestor_service]
}

# --- Persistor Service Deployment ---
resource "kubernetes_manifest" "persistor_deployment" {
  provider   = kubernetes
  manifest   = yamldecode(file("./../../k8s/persistor/persistor-deployment.yaml"))
  depends_on = [kubernetes_manifest.ingestor_service]
}

resource "kubernetes_secret" "persistor_secrets" {
  metadata {
    name      = "persistor-secrets"
    namespace = "default"
  }
  type = "Opaque"
  data = {
    db_write_dsn = base64decode(var.db_write_dsn)
    db_read_dsn  = base64decode(var.db_read_dsn)
  }
}

resource "kubernetes_config_map" "persistor_config_map" {
  provider = kubernetes
  metadata {
    name      = "persistor-config"
    namespace = "default"
  }
  data = {
    app_name : "binance-trading-chart-service/persistor"
    app_env : var.app_env
    app_debug : var.app_debug
    server_address : format("ingestor-service.default.svc.cluster.local:%s", var.app_grpc_port)
  }
}
