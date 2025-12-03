output "app_info_path" {
  description = "Path to app info file"
  value       = local_file.app_info.filename
}

output "config_path" {
  description = "Path to config JSON"
  value       = local_file.config_json.filename
}

output "log_config_path" {
  description = "Path to logging config"
  value       = var.enable_logging ? local_file.log_config[0].filename : "N/A (logging disabled)"
}

output "summary" {
  description = "Configuration summary"
  value = {
    app_name    = var.app_name
    environment = var.environment
    version     = var.app_version
    logging     = var.enable_logging
  }
}
