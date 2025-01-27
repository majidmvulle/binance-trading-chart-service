variable "k8s_host" {
  type        = string
  description = "K8s host server"
}

variable "k8s_client_certificate" {
  type        = string
  description = "K8s client certificate for Kubernetes API"
}

variable "k8s_client_key" {
  type        = string
  description = "K8s client key for Kubernetes API"
}

variable "k8s_cluster_ca_certificate" {
  type        = string
  description = "K8s cluster client ca certificate for Kubernetes API"
}

variable "db_write_dsn" {
  type        = string
  description = "PostgreSQL write database connection string for persistor service"
  sensitive   = true
}

variable "db_read_dsn" {
  type        = string
  description = "PostgreSQL read database connection string for persistor service"
  sensitive   = true
}

variable "binance_symbols" {
  type        = string
  description = "Symbols to track"
  default     = "BTCUSDT ETHUSDT PEPEUSDT"
}

variable "app_grpc_port" {
  type        = string
  description = "GRPC port"
  default     = "50051"
}

variable "app_env" {
  type        = string
  description = "Environment"
  default     = "development"
}

variable "app_debug" {
  type        = string
  description = "Whether to debug or not"
  default     = "true"
}
