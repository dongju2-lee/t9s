output "monitoring_config_path" {
  description = "Path to monitoring configuration file"
  value       = local_file.monitoring_config.filename
}

output "dashboard_url" {
  description = "Path to dashboard HTML file"
  value       = "file://${abspath(local_file.dashboard_config.filename)}"
}

output "environment" {
  description = "Current environment"
  value       = var.environment
}

output "alerts_enabled" {
  description = "Whether alerts are enabled"
  value       = var.enable_alerts
}

output "setup_timestamp" {
  description = "When the monitoring was set up"
  value       = null_resource.setup_monitoring.triggers.timestamp
}
