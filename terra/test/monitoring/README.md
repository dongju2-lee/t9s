# Monitoring Test Environment

This is a test Terraform configuration that creates local files instead of cloud resources.
Perfect for testing T9s features like history, apply, destroy without actual infrastructure costs.

## What it does

- Creates monitoring configuration files (JSON)
- Generates a simple HTML dashboard
- Creates alert rules (YAML) if alerts are enabled
- Simulates provisioning with local-exec commands
- Logs setup activities

## Files Created

All files are created in the `output/` directory:
- `monitoring-{env}.json` - Monitoring configuration
- `dashboard-{env}.html` - Dashboard HTML file
- `alerts-{env}.yaml` - Alert rules (if enabled)
- `logs/setup.log` - Setup activity log

## Testing Scenarios

### 1. Development Environment
```bash
terraform init
terraform plan -var-file=config/env.tfvars
terraform apply -var-file=config/env.tfvars
```

### 2. Production Environment
```bash
terraform apply -var-file=config/prod.tfvars
```

### 3. Staging (No Alerts)
```bash
terraform apply -var-file=config/staging.tfvars
```

### 4. Destroy
```bash
terraform destroy -var-file=config/env.tfvars
```

## Testing T9s Features

1. **History**: Apply with different tfvars files to see history entries
2. **Tab Focus**: Apply to see long output, then press Tab to scroll
3. **File Selection**: Select different .tfvars files in tree to use different configs
4. **Plan vs Apply**: Compare plan output with apply output

## Clean Up

```bash
rm -rf output/
rm -rf .terraform/
rm terraform.tfstate*
```
