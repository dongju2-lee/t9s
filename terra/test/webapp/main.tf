terraform {
  required_version = ">= 1.0"
  required_providers {
    local = {
      source  = "hashicorp/local"
      version = "~> 2.0"
    }
    null = {
      source  = "hashicorp/null"
      version = "~> 3.0"
    }
  }
}

# Variables
variable "app_name" {
  description = "Application name"
  type        = string
  default     = "myapp"
}

variable "environment" {
  description = "Environment (dev, staging, prod)"
  type        = string
  default     = "dev"
}

variable "app_version" {
  description = "Application version"
  type        = string
  default     = "1.0.0"
}

variable "enable_logging" {
  description = "Enable logging"
  type        = bool
  default     = true
}

# Simple text file with app info
resource "local_file" "app_info" {
  filename = "${path.module}/output/app-${var.environment}.txt"
  content  = <<-EOT
    ========================================
    Application: ${var.app_name}
    Environment: ${var.environment}
    Version: ${var.app_version}
    Logging: ${var.enable_logging ? "Enabled" : "Disabled"}
    Created: ${timestamp()}
    ========================================
  EOT
  file_permission = "0644"
}

# Configuration JSON
resource "local_file" "config_json" {
  filename = "${path.module}/output/config-${var.environment}.json"
  content = jsonencode({
    app_name    = var.app_name
    environment = var.environment
    version     = var.app_version
    logging     = var.enable_logging
    timestamp   = timestamp()
  })
  file_permission = "0644"
}

# Log file (if enabled)
resource "local_file" "log_config" {
  count = var.enable_logging ? 1 : 0
  
  filename = "${path.module}/output/logging-${var.environment}.conf"
  content  = <<-EOT
    # Logging Configuration for ${var.app_name}
    level: ${var.environment == "prod" ? "error" : "debug"}
    output: /var/log/${var.app_name}/${var.environment}.log
    format: json
    timestamp: true
  EOT
  file_permission = "0644"
}

# Simple echo commands - no actual installation
resource "null_resource" "setup" {
  triggers = {
    environment = var.environment
    version     = var.app_version
  }

  provisioner "local-exec" {
    command = "echo '🚀 Setting up ${var.app_name} v${var.app_version} for ${var.environment}'"
  }

  provisioner "local-exec" {
    command = "echo '� Creating output directory...'"
  }

  provisioner "local-exec" {
    command = "mkdir -p ${path.module}/output/logs"
  }

  provisioner "local-exec" {
    command = "echo '✅ Setup complete!' >> ${path.module}/output/logs/setup.log"
  }
}
