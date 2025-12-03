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

# Variables from tfvars
variable "environment" {
  description = "Environment name"
  type        = string
  default     = "dev"
}

variable "project_name" {
  description = "Project name"
  type        = string
  default     = "monitoring"
}

variable "enable_alerts" {
  description = "Enable monitoring alerts"
  type        = bool
  default     = true
}

# Create a local file with monitoring configuration
resource "local_file" "monitoring_config" {
  filename = "${path.module}/output/monitoring-${var.environment}.json"
  content = jsonencode({
    environment   = var.environment
    project       = var.project_name
    alerts_enabled = var.enable_alerts
    timestamp     = timestamp()
    metrics = {
      cpu_threshold    = 80
      memory_threshold = 85
      disk_threshold   = 90
    }
  })
  file_permission = "0644"
}

# Create a dashboard configuration file
resource "local_file" "dashboard_config" {
  filename = "${path.module}/output/dashboard-${var.environment}.html"
  content  = <<-EOT
    <!DOCTYPE html>
    <html>
    <head>
      <title>${var.project_name} - ${var.environment}</title>
    </head>
    <body>
      <h1>Monitoring Dashboard</h1>
      <p>Environment: ${var.environment}</p>
      <p>Project: ${var.project_name}</p>
      <p>Alerts: ${var.enable_alerts ? "Enabled" : "Disabled"}</p>
    </body>
    </html>
  EOT
  file_permission = "0644"
}

# Simulate a provisioning step with null_resource
resource "null_resource" "setup_monitoring" {
  triggers = {
    environment = var.environment
    timestamp   = timestamp()
  }

  provisioner "local-exec" {
    command = "echo 'Setting up monitoring for ${var.environment}...'"
  }

  provisioner "local-exec" {
    command = "mkdir -p ${path.module}/output/logs"
  }

  provisioner "local-exec" {
    command = "echo '[${timestamp()}] Monitoring setup completed for ${var.environment}' >> ${path.module}/output/logs/setup.log"
  }
}

# Create alert rules file
resource "local_file" "alert_rules" {
  count = var.enable_alerts ? 1 : 0
  
  filename = "${path.module}/output/alerts-${var.environment}.yaml"
  content  = <<-EOT
    alerts:
      - name: high_cpu
        threshold: 80
        severity: warning
      - name: high_memory
        threshold: 85
        severity: warning
      - name: disk_full
        threshold: 90
        severity: critical
    environment: ${var.environment}
  EOT
  file_permission = "0644"
}
